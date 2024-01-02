package payments

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rlvgl/bookie-server/users"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/price"
	"github.com/stripe/stripe-go/v76/product"
	"github.com/stripe/stripe-go/v76/webhook"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewPaymentServer() *PaymentServer {
	return &PaymentServer{
		client: ConnectDB(),
	}
}

type PaymentServer struct {
	client *mongo.Client
}

func ConnectDB() *mongo.Client {
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client
}

func (s *PaymentServer) StartPaymentServerAPI(api fiber.Router) error {
	log.Println("Adding Payment server endpoints to API")
	// Payment
	paymentsApi := api.Group("/payments")
	paymentsApi.Get("/", func(c *fiber.Ctx) error {
		c.Send([]byte("hello world"))
		return nil
	})

	paymentsApi.Post("/checkout/:id", func(c *fiber.Ctx) error {
		return s.HandleCheckoutRequest(c)
	})

	// paymentsApi.Post("/payout", func(c *fiber.Ctx) error {
	// 	return s.HandlePayoutRequest(c)
	// })

	paymentsApi.Post("/webhook", func(c *fiber.Ctx) error {
		return s.HandleWebhook(c)
	})

	paymentsApi.Post("/payout", func(c *fiber.Ctx) error {
		return s.HandlePayout(c)
	})
	return nil
}

type Payment struct {
	StripeUserId string    `json:"stripe_id" bson:"stripe_id"`
	UserEmail    string    `json:"user_email" bson:"user_email"`
	Amount       float64   `json:"amount" bson:"amount"`
	DatePlaced   time.Time `json:"date_placed" bson:"date_placed"`
}

func (s *PaymentServer) HandleWebhook(c *fiber.Ctx) error {
	fmt.Println("handling webhook")
	const MaxBodyBytes = int64(65536)
	body := string(c.Request().Body())
	// fmt.Fprintf(os.Stdout, "Got body: %s\n", body)

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SIGNING_SECRET_TEST")
	// event, err := webhook.ConstructEvent([]byte(body), c.GetReqHeaders()["Stripe-Signature"], endpointSecret)
	event, err := webhook.ConstructEventWithOptions([]byte(body), c.GetReqHeaders()["Stripe-Signature"], endpointSecret, webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		c.SendStatus(http.StatusBadRequest)
		return err
	}

	if event.Type == "checkout.session.completed" {
		var sesh stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &sesh)
		if err != nil {
			return err
		}
		params := &stripe.CheckoutSessionParams{}
		params.AddExpand("line_items")
		if sesh.PaymentStatus == "paid" {
			// add payment doc to db
			// update user with new balance
			subtotal := float64(sesh.AmountSubtotal) / 100

			// fmt.Printf("user email: %v\n", sesh.CustomerEmail)
			fmt.Printf("user id: %v\n", sesh.ClientReferenceID)
			coll := s.client.Database("finances-db").Collection("payments")
			payment := Payment{
				UserEmail:  sesh.CustomerEmail,
				Amount:     subtotal,
				DatePlaced: time.Now(),
			}
			coll.InsertOne(context.TODO(), payment)

			// find better way to do in one query
			// rn have to get user then update
			user := users.User{}
			userColl := s.client.Database("users-db").Collection("users")
			id, _ := primitive.ObjectIDFromHex(sesh.ClientReferenceID)
			filter := bson.M{"_id": id}
			userColl.FindOne(context.TODO(), filter).Decode(&user)

			newBalance := user.CurrentBalance + subtotal
			update := bson.M{"$set": bson.M{"current_balance": newBalance}}
			_, err := userColl.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Finished processing payment. Updated user balance")
		}
	}

	c.SendStatus(http.StatusOK)

	return nil
}

type Payout struct {
	Email      string    `json:"email" bson:"email"`
	Name       string    `json:"name" bson:"name"`
	UserId     string    `json:"user_id" bson:"user_id"`
	Amount     float64   `json:"amount" bson:"amount"`
	TimePlaced time.Time `json:"time_placed" bson:"time_placed"`
}

func (s *PaymentServer) HandlePayoutRequest(c *fiber.Ctx) error {
	payout := Payout{}
	c.BodyParser(&payout)

	if payout.Email == "" || payout.Name == "" || payout.UserId == "" {
		// theres actually an error here tbh
		c.SendStatus(http.StatusBadRequest)
		return fmt.Errorf("Invalid request")
	}

	// get amount from mongo
	coll := s.client.Database("users-db").Collection("users")
	id, _ := primitive.ObjectIDFromHex(payout.UserId)
	filter := bson.M{"_id": id}
	user := users.User{}
	coll.FindOne(context.TODO(), filter).Decode(&user)

	if user.CurrentBalance < 0.50 {
		return fmt.Errorf("user is a broke sack of shit")
	}
	payout.Amount = user.CurrentBalance

	// update user balance to back to 0
	update := bson.M{"$set": bson.M{"current_balance": 0}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Printf("Error updating user balance. User of id %v has %v in balance: %n", payout.UserId, payout.Amount, err)
	}

	url := "https://testflight.tremendous.com/api/v2/orders"

	str := fmt.Sprintf("{\"payment\":{\"funding_source_id\":\"BALANCE\",\"channel\":\"UI\"},\"reward\":{\"campaign_id\":\"T2WXQSQDTNQP\",\"products\":[],\"value\":{\"denomination\":%v,\"currency_code\":\"USD\"},\"recipient\":{\"name\":\"%v\",\"email\":\"%v\",\"phone\":\"123-456-7890\"},\"delivery\":{\"method\":\"EMAIL\"},\"language\":\"en\"}}", payout.Amount, payout.Name, payout.Email)
	payload := strings.NewReader(str)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Bearer TEST_RJ6njjMQo--lBLO7lD56Z7K8PaqoQntw5ABsNy38TMV")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	payout.TimePlaced = time.Now()
	coll = s.client.Database("finances-db").Collection("payouts")
	coll.InsertOne(context.TODO(), payout)

	c.Send([]byte("Successfully sent payout to user"))
	return nil
}

func (s *PaymentServer) HandleCheckoutRequest(c *fiber.Ctx) error {
	// domain := "http://localhost:8080"
	// domain := "http://localhost:5173/dashboard"
	// only need stripe customer id and user id
	user := users.User{}
	c.BodyParser(&user)

	productParams := &stripe.ProductParams{Name: stripe.String("Add Funds to Lemondrop")}
	result, err := product.New(productParams)

	priceParams := &stripe.PriceParams{
		Currency:         stripe.String(string(stripe.CurrencyUSD)),
		CustomUnitAmount: &stripe.PriceCustomUnitAmountParams{Enabled: stripe.Bool(true)},
		Product:          stripe.String(result.ID),
	}
	prresult, err := price.New(priceParams)
	// fmt.Println(prresult.ID)
	if err != nil {
		fmt.Println(err)
	}

	if prresult.ID == "" {
		fmt.Println("no id")
	}
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(prresult.ID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:              stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:        stripe.String("https://lemondrop.ag/bets"),
		CancelURL:         stripe.String("https://lemondrop.ag/bets"),
		ClientReferenceID: stripe.String(c.Params("id")),
	}

	res, err := session.New(params)

	if err != nil {
		log.Printf("session.New: %v", err)
	}

	c.Redirect(res.URL)

	return nil
}

func (s *PaymentServer) Start(api fiber.Router) error {
	stripe.Key = os.Getenv("STRIPE_SECRET_TEST_KEY")
	s.StartPaymentServerAPI(api)
	return nil
}
