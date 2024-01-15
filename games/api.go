package games

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GameServer struct {
	client *mongo.Client
}

func NewGameServer() *GameServer {
	return &GameServer{client: ConnectDB()}
}

func ConnectDB() *mongo.Client {
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client
}

func (s *GameServer) DeleteAllGames() error {
	coll := s.client.Database("games-db").Collection("scraped-games")
	_, err := coll.DeleteMany(context.Background(), bson.D{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *GameServer) StartGameServerAPI(api fiber.Router) error {
	log.Println("Adding Game server endpoints to API")
	gamesApi := api.Group("/games")
	gamesApi.Get("/", func(c *fiber.Ctx) error {
		c.Send([]byte("games api is working"))
		return nil
	})

	gamesApi.Get("/sport/:sport", func(c *fiber.Ctx) error {
		return s.GetGamesBySport(c)
	})

	gamesApi.Get("/league/:league", func(c *fiber.Ctx) error {
		return s.GetGamesByLeagueId(c)
	})

	gamesApi.Get("/game/:id", func(c *fiber.Ctx) error {
		return s.GetGameById(c)
	})

	leaguesApi := api.Group("/leagues")

	leaguesApi.Get("/", func(c *fiber.Ctx) error {
		return s.GetAllLeagues(c)
	})

	leaguesApi.Get("/all", func(c *fiber.Ctx) error {
		return s.GetAllSportsAndLeagues(c)
	})

	leaguesApi.Get("/sports", func(c *fiber.Ctx) error {
		return s.GetAllSports(c)
	})

	leaguesApi.Get("/:sport", func(c *fiber.Ctx) error {
		return s.GetAllLeaguesBySport(c)
	})

	return nil
}

func (s *GameServer) DecodeCursorIntoGames(cursor *mongo.Cursor) (*[]Game, error) {
	games := []Game{}
	game := &Game{}
	for cursor.Next(context.TODO()) {
		cursor.Decode(&game)
		games = append(games, *game)
	}
	return &games, nil
}

func (s *GameServer) Start(api fiber.Router) error {
	s.StartGameServerAPI(api)
	return nil
}
