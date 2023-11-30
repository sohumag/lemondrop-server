package mailing

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MailingServer struct {
	client *mongo.Client
}

func NewMailingServer() *MailingServer {
	return &MailingServer{
		client: ConnectDB(),
	}
}

func ConnectDB() *mongo.Client {
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client
}

func (s *MailingServer) StartMailingServerAPI(api fiber.Router) error {
	log.Println("Adding mailing server endpoints to API")

	g := api.Group("/mailing")
	g.Post("/add", func(c *fiber.Ctx) error {
		return s.AddUserToMailingList(c)
	})

	return nil
}

type Email struct {
	Email string `json:"email" bson:"email"`
}

func (s *MailingServer) AddUserToMailingList(c *fiber.Ctx) error {
	coll := s.client.Database("mailing").Collection("emails")
	email := Email{}
	c.BodyParser(&email)
	if email.Email == "" {
		c.SendStatus(http.StatusBadRequest)
		return fmt.Errorf("Missing required field email")
	}
	coll.InsertOne(context.TODO(), email)
	return nil
}

func (s *MailingServer) Start(api fiber.Router) error {
	s.StartMailingServerAPI(api)
	return nil
}
