package payments

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
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

	return nil
}

func (s *PaymentServer) Start(api fiber.Router) error {
	s.StartPaymentServerAPI(api)
	return nil
}
