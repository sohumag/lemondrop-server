package scores

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// private only used by other services
type ScoreServer struct {
	client *mongo.Client
}

func NewScoreServer() *ScoreServer {
	return &ScoreServer{client: ConnectDB()}
}

func ConnectDB() *mongo.Client {
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client
}

func (s *ScoreServer) StartScoreServerAPI(api fiber.Router) error {
	log.Println("Adding score server endpoints to API")

	scores := api.Group("scores")
	scores.Get("/", func(c *fiber.Ctx) error {
		c.SendStatus(http.StatusOK)
		return nil
	})

	return nil
}

func (s *ScoreServer) Start(api fiber.Router) error {
	go func() {
		s.StartScoresUpdates()
	}()
	s.StartScoreServerAPI(api)
	return nil
}
