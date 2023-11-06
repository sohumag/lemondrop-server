package bets

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewBetServer() *BetServer {
	return &BetServer{client: ConnectDB()}
}

type BetServer struct{ client *mongo.Client }

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
		return s.AddBetToDB(c)
	})

	betsApi.Get("/bet/user/:user", func(c *fiber.Ctx) error {
		return s.GetAllBetsByUserId(c, c.Params("user"))
	})

	betsApi.Get("/all", func(c *fiber.Ctx) error {
		return s.GetAllBets(c)
	})

	return nil
}

func (s *BetServer) Start(api fiber.Router) error {
	s.StartBetServerAPI(api)
	return nil
}

type Bet struct {
	// user information
	UserId           string  `json:"user_id"`
	UserEmail        string  `json:"user_email"`
	UserBalance      float64 `json:"user_balance"`
	UserAvailability float64 `json:"user_availability"`
	UserPending      float64 `json:"user_pending"`
	UserFreePlay     float64 `json:"user_free_play"`

	// game information
	GameId        string    `json:"game_id"`
	HomeTeam      string    `json:"home_team"`
	AwayTeam      string    `json:"away_team"`
	GameStartTime time.Time `json:"game_start_time"`

	// bet information
	BetId       string  `json:"bet_id"`
	BetType     string  `json:"bet_type"`
	BetOnTeam   string  `json:"bet_on_team"`
	BetCategory string  `json:"bet_category"` // for player props etc: really bet description
	BetPoint    float64 `json:"bet_point"`
	BetPrice    float64 `json:"bet_price"`
	BetAmount   float64 `json:"bet_amount"`
	BetVerified bool    `json:"bet_verified"`
	BetCashed   bool    `json:"bet_cashed"`
	BetStatus   string  `json:"bet_status"`
}
