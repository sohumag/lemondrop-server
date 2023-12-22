package bets

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewBetServer() *BetServer {
	return &BetServer{client: ConnectDB(), queue: make(chan *Bet, 500)}
}

type BetServer struct {
	client *mongo.Client
	queue  chan *Bet
	wg     sync.WaitGroup
}

func ConnectDB() *mongo.Client {
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client
}

func (s *BetServer) StartBetServerAPI(api fiber.Router) error {
	log.Println("Adding bet server endpoints to API")

	betsApi := api.Group("/bets")
	betsApi.Get("/", func(c *fiber.Ctx) error {
		c.Send([]byte("Bets api is active"))
		return nil
	})

	betsApi.Post("/bet", func(c *fiber.Ctx) error {
		return s.HandleBetRequest(c)
	})

	// betsApi.Get("/bet/user/:user", func(c *fiber.Ctx) error {
	// 	return s.GetAllBetsByUserId(c, c.Params("user"))
	// })

	// betsApi.Get("/all", func(c *fiber.Ctx) error {
	// 	return s.GetAllBets(c)
	// })

	return nil
}

func (s *BetServer) Start(api fiber.Router) error {
	// s.BetChecker()
	// go s.ProcessBets()
	s.StartBetServerAPI(api)
	return nil
}

func (s *BetServer) DeleteAllBets() {
	coll := s.client.Database("bets-db").Collection("bets")

	// Define a filter to match all documents (empty filter)
	filter := bson.D{}

	// Delete all documents in the collection
	_, err := coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
}
