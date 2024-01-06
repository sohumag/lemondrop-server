package payments

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rlvgl/bookie-server/users"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/price"
	"github.com/stripe/stripe-go/v76/product"
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

	paymentsApi.Post("/webhook/:source", func(c *fiber.Ctx) error {
		return s.HandleWebhook(c)
	})

	paymentsApi.Post("/payout", func(c *fiber.Ctx) error {
		return s.HandlePayout(c)
	})

	// paymentsApi.Post("/payout", func(c *fiber.Ctx) error {
	// 	return s.HandlePayout(c)
	// })
	return nil
}

type Payment struct {
	StripeUserId string    `json:"stripe_id" bson:"stripe_id"`
	UserEmail    string    `json:"user_email" bson:"user_email"`
	Amount       float64   `json:"amount" bson:"amount"`
	DatePlaced   time.Time `json:"date_placed" bson:"date_placed"`
}

type Payout struct {
	Email      string    `json:"email" bson:"email"`
	Name       string    `json:"name" bson:"name"`
	UserId     string    `json:"user_id" bson:"user_id"`
	Amount     float64   `json:"amount" bson:"amount"`
	TimePlaced time.Time `json:"time_placed" bson:"time_placed"`
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
