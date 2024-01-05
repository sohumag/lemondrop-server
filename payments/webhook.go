package payments

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rlvgl/bookie-server/users"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *PaymentServer) HandleWebhook(c *fiber.Ctx) error {
	// webhookSource := c.Params("source")
	fmt.Println()
	fmt.Println("Webhook:")

	const MaxBodyBytes = int64(65536)
	body := c.Request().Body()

	// Read the body into a byte slice
	bodyBytes := []byte(body)
	// Create an io.Reader from the byte slice
	bodyReader := bytes.NewReader(bodyBytes)

	// Parse the JSON payload to get the event type
	var stripeEvent stripe.Event
	if err := json.NewDecoder(bodyReader).Decode(&stripeEvent); err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding JSON payload: %v\n", err)
		c.SendStatus(http.StatusBadRequest)
		return err
	}

	// Access the event type from the parsed event
	eventType := stripeEvent.Type
	fmt.Println("Event Type:", eventType)

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SIGNING_SECRET_TEST")
	_, err := webhook.ConstructEventWithOptions(bodyBytes, c.GetReqHeaders()["Stripe-Signature"], endpointSecret, webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true})

	//! Perform signature verification (commented out for testing)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
	// 	c.SendStatus(http.StatusBadRequest)
	// 	return err
	// }

	// Handle the webhook event based on its type
	err = s.HandleStripeWebhookEvent(stripeEvent)
	if err != nil {
		return err
	}

	c.SendStatus(http.StatusOK)

	return nil
}

func (s *PaymentServer) HandleStripeWebhookEvent(event stripe.Event) error {
	if event.Type == "checkout.session.completed" {
		err := s.HandleStripeCheckoutSessionCompletedEvent(event)
		if err != nil {
			return err
		}
	}

	if event.Type == "account.updated" {
		err := s.HandleStripeConnectAccountUpdated(event)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PaymentServer) HandleStripeConnectAccountUpdated(event stripe.Event) error {
	data := event.Data.Object
	detailsSubmitted := data["details_submitted"]
	accountId := data["id"]

	fmt.Println(detailsSubmitted)
	fmt.Println(accountId)

	coll := s.client.Database("users-db").Collection("users")
	filter := bson.D{{Key: "stripe_express_id", Value: accountId}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "details_submitted", Value: detailsSubmitted},
		}},
	}
	coll.ReplaceOne(context.Background(), filter, update)

	return nil
}

func (s *PaymentServer) HandleStripeCheckoutSessionCompletedEvent(event stripe.Event) error {
	var sesh stripe.CheckoutSession
	err := json.Unmarshal(event.Data.Raw, &sesh)
	if err != nil {
		return err
	}
	params := &stripe.CheckoutSessionParams{}
	params.AddExpand("line_items")
	if sesh.PaymentStatus != "paid" {
		return nil
	}

	// user has paid
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
	_, err = userColl.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Finished processing payment. Updated user balance")

	return nil
}
