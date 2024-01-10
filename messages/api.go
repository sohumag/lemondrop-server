package messages

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMessageServer() *MessageServer {
	return &MessageServer{
		client: ConnectDB(),
	}
}

type MessageServer struct {
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

func (s *MessageServer) StartMessageServerAPI(api fiber.Router) error {
	log.Println("Adding Message server endpoints to API")
	// Message
	mapi := api.Group("/messages")

	mapi.Post("/contact", func(c *fiber.Ctx) error {
		return s.HandleContactMessageRequest(c)
	})

	return nil
}

func (s *MessageServer) Start(api fiber.Router) error {
	s.StartMessageServerAPI(api)
	return nil
}

func (s *MessageServer) HandleContactMessageRequest(c *fiber.Ctx) error {
	// Parse the JSON request body into a ContactMessage struct
	var contactMessage ContactMessage
	if err := c.BodyParser(&contactMessage); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON body: %v\n", err)
		c.Status(fiber.StatusBadRequest)
		return err
	}

	// Access the "messages" database and create the "contact" collection
	db := s.client.Database("messages")
	contactCollection := db.Collection("contact")

	// Insert the ContactMessage document into the "contact" collection
	_, err := contactCollection.InsertOne(context.Background(), contactMessage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error inserting document into 'contact' collection: %v\n", err)
		c.Status(fiber.StatusInternalServerError)
		return err
	}

	return nil
}

type ContactMessage struct {
	Id      string `json:"id" bson:"_id"`
	Name    string `json:"name" bson:"name"`
	Email   string `json:"email" bson:"email"`
	Message string `json:"message" bson:"message"`
}
