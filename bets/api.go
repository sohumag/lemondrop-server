package bets

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		return s.HandleBetRequest(c)
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
	s.RunBetCheckingRepeater()
	s.StartBetServerAPI(api)
	return nil
}

type Bet struct {
	// user information
	UserId      string  `json:"user_id" bson:"user_id"`
	UserEmail   string  `json:"user_email" bson:"user_email"`
	UserBalance float64 `json:"user_balance" bson:"user_balance"`
	// UserAvailability float64 `json:"user_availability" bson:"user_availability"`
	UserPending  float64 `json:"user_pending" bson:"user_pending"`
	UserFreePlay float64 `json:"user_free_play" bson:"user_free_play"`
	TotalProfit  float64 `json:"total_profit" bson:"total_profit"`

	// game information
	GameId        string    `json:"game_id" bson:"game_id"`
	GameHash      string    `json:"game_hash" bson:"game_hash"`
	HomeTeam      string    `json:"home_team" bson:"home_team"`
	AwayTeam      string    `json:"away_team" bson:"away_team"`
	GameStartTime time.Time `json:"game_start_time" bson:"game_start_time"`
	GameErr       bool      `json:"game_err" bson:"game_err"`

	// bet information
	BetId          primitive.ObjectID `json:"bet_id" bson:"_id"`
	BetType        string             `json:"bet_type" bson:"bet_type"`
	BetOnTeam      string             `json:"bet_on_team" bson:"bet_on_team"`
	BetCategory    string             `json:"bet_category" bson:"bet_category"` // for player props etc: really bet description
	BetPoint       string             `json:"bet_point" bson:"bet_point"`
	BetPrice       string             `json:"bet_price" bson:"bet_price"`
	BetAmount      string             `json:"bet_amount" bson:"bet_amount"`
	BetVerified    bool               `json:"bet_verified" bson:"bet_verified"`
	BetCashed      bool               `json:"bet_cashed" bson:"bet_cashed"`
	BetStatus      string             `json:"bet_status" bson:"bet_status"` // Pending, Won, Lost, Pushed
	IsParlay       bool               `json:"is_parlay" bson:"is_parlay"`
	ParlayFinished bool               `json:"parlay_finished" bson:"parlay_finished"`
	Bets           []Bet              `json:"bets" bson:"bets"`

	BetPlacedTime time.Time `json:"bet_placed_time" bson:"bet_placed_time"`
}
