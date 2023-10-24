package bets

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewBetServer() *BetServer {
	return &BetServer{
		client: ConnectDB(),
	}
}

type BetServer struct {
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

func (s *BetServer) StartBetServerAPI(api fiber.Router) error {
	log.Println("Adding bet server endpoints to API")

	betsApi := api.Group("/bets")
	betsApi.Get("/", func(c *fiber.Ctx) error {
		c.Send([]byte("Bets api is active"))
		return nil
	})

	betsApi.Post("/add", func(c *fiber.Ctx) error {
		return s.AddBetToDB(c)
	})

	return nil
}

func (s *BetServer) Start(api fiber.Router) error {
	s.StartBetServerAPI(api)
	return nil
}

type Bet struct {
	UserId           string  `json:"user_id"`
	UserEmail        string  `json:"user_email"`
	UserBalance      float64 `json:"user_balance"`
	UserAvailability float64 `json:"user_availability"`
	UserPending      float64 `json:"user_pending"`
	UserFreePlay     float64 `json:"user_free_play"`

	GameId        string    `json:"game_id"`
	HomeTeam      string    `json:"home_team"`
	AwayTeam      string    `json:"away_team"`
	GameStartTime time.Time `json:"game_start_time"`

	BetType   string  `json:"bet_type"`
	BetOnTeam string  `json:"bet_on_team"`
	BetPoint  float64 `json:"bet_point"`
	BetPrice  float64 `json:"bet_price"`
	BetAmount float64 `json:"bet_amount"`

	BetCashed bool `json:"bet_cashed"`
}

func (s *BetServer) AddBetToDB(c *fiber.Ctx) error {
	bet := &Bet{}
	if err := c.BodyParser(&bet); err != nil {
		fmt.Println(err)
	}

	coll := s.client.Database("bets-db").Collection("bets")
	result, err := coll.InsertOne(context.TODO(), &bet)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Bet placed with id: %v\n", result.InsertedID)
	return nil
}
