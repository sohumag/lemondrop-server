package payments

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rlvgl/bookie-server/users"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/transfer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RequestBody struct {
	UserID string `json:"user_id"`
	// Add other fields as needed
}

func (s *PaymentServer) HandlePayout(c *fiber.Ctx) error {
	// get user id from body
	// check amount in balance of user, set to 0.
	// add payout intent to personal mongo database
	// get stripe express id from user struct
	// transfer .98 * balance to user
	// return

	var requestBody RequestBody
	if err := c.BodyParser(&requestBody); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON body: %v\n", err)
		c.Status(fiber.StatusBadRequest)
		return err
	}

	// Access the user_id field
	userId := requestBody.UserID

	coll := s.client.Database("users-db").Collection("users")
	userIdBin, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: userIdBin}}

	user := users.User{}
	coll.FindOne(context.Background(), filter).Decode(&user)

	balance := user.CurrentBalance
	if balance < 5 {
		return fmt.Errorf("Insufficient balance to cash out.")
	}
	connectedAccountId := user.StripeExpressId

	fmt.Println(balance)
	fmt.Println(connectedAccountId)

	// update := bson.D{{"current_balance", 0}}
	// coll.UpdateOne(context.Background(), filter, update)

	coll = s.client.Database("payments").Collection("payouts")
	payout := MongoPayout{
		UserId: userId,
		Amount: balance,
		Time:   time.Now(),
	}
	coll.InsertOne(context.Background(), payout)

	stripe.Key = os.Getenv("STRIPE_SECRET_TEST_KEY")

	params := &stripe.TransferParams{
		Amount:      stripe.Int64(balance),
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
		Destination: stripe.String(connectedAccountId),
	}

	tr, err := transfer.New(params)
	if err != nil {
		return err
	}

	fmt.Println(tr)

	return nil
}

type MongoPayout struct {
	UserId string    `json:"user_id" bson:"user_id"`
	Amount float64   `json:"amount" bson:"amount"`
	Time   time.Time `json:"time" bson:"time"`
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
