package games

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (s *GameServer) StartGameServerAPI(api fiber.Router) error {
	log.Println("Adding Game server endpoints to API")
	gamesApi := api.Group("/games")
	gamesApi.Get("/", func(c *fiber.Ctx) error {
		c.Send([]byte("games api is working"))
		return nil
	})

	gamesApi.Get("/all", func(c *fiber.Ctx) error {
		return s.GetAllGamesFromDB(c)
	})

	gamesApi.Get("/sport/:sport", func(c *fiber.Ctx) error {
		return s.GetGamesBySport(c)
	})

	gamesApi.Get("/league/:league", func(c *fiber.Ctx) error {
		return s.GetGamesByLeagueId(c)
	})

	return nil
}

func (s *GameServer) GetAllGamesFromDB(c *fiber.Ctx) error {
	coll := s.client.Database("games-db").Collection("scraped-games")
	currentDate := time.Now()
	maxDate := currentDate.Add(time.Hour * 24 * 7)
	cursor, err := coll.Find(context.TODO(), bson.M{"start_date": bson.M{"$gt": currentDate, "$lt": maxDate}})
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	games, err := s.DecodeCursorIntoGames(cursor)
	if err != nil {
		return err
	}

	c.JSON(*games)

	return nil
}

func (s *GameServer) GetGamesBySport(c *fiber.Ctx) error {
	coll := s.client.Database("games-db").Collection("scraped-games")
	currentDate := time.Now()
	maxDate := currentDate.Add(time.Hour * 24 * 7)

	parsedSport := strings.Replace(c.Params("sport"), "%20", " ", -1)
	// parsedSport = strings.Replace(c.Params("sport"), "-", " ", -1)
	parsedSport = strings.Title(parsedSport)

	cursor, err := coll.Find(context.TODO(), bson.M{"start_date": bson.M{"$gt": currentDate, "$lt": maxDate}, "sport": parsedSport})
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	games, err := s.DecodeCursorIntoGames(cursor)
	if err != nil {
		return err
	}

	c.JSON(*games)

	return nil
}
func (s *GameServer) GetGamesByLeagueId(c *fiber.Ctx) error {
	coll := s.client.Database("games-db").Collection("scraped-games")
	currentDate := time.Now()
	maxDate := currentDate.Add(time.Hour * 24 * 7)

	parsedLeague := strings.Replace(c.Params("league"), "%20", " ", -1)
	parsedLeague = strings.ToLower(parsedLeague)

	cursor, err := coll.Find(context.TODO(), bson.M{"start_date": bson.M{"$gt": currentDate, "$lt": maxDate}, "league_id": parsedLeague})
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	games, err := s.DecodeCursorIntoGames(cursor)
	if err != nil {
		return err
	}

	c.JSON(*games)

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

type Game struct {
	Id              primitive.ObjectID `json:"_id" bson:"_id"`
	GameType        string             `json:"game_type" bson:"game_type"`
	League          string             `json:"league" bson:"league" `
	LeagueId        string             `json:"league_id" bson:"league_id"`
	Sport           string             `json:"sport" bson:"sport"`
	StartDate       time.Time          `json:"start_date" bson:"start_date"`
	LastUpdated     time.Time          `json:"last_updated" bson:"last_updated"`
	Hash            string             `json:"hash" bson:"hash"`
	AwayTeamName    string             `json:"away_team_name" bson:"away_team_name"`
	HomeTeamName    string             `json:"home_team_name" bson:"home_team_name"`
	AwayMoneyline   string             `json:"away_moneyline" bson:"away_moneyline"`
	HomeMoneyline   string             `json:"home_moneyline" bson:"home_moneyline"`
	DrawMoneyline   string             `json:"draw_moneyline" bson:"draw_moneyline"`
	AwaySpreadPoint string             `json:"away_spread_point" bson:"away_spread_point"`
	AwaySpreadPrice string             `json:"away_spread_price" bson:"away_spread_price"`
	HomeSpreadPoint string             `json:"home_spread_point" bson:"home_spread_point"`
	HomeSpreadPrice string             `json:"home_spread_price" bson:"home_spread_price"`
	UnderPoint      string             `json:"under_point" bson:"under_point"`
	UnderPrice      string             `json:"under_price" bson:"under_price"`
	OverPoint       string             `json:"over_point" bson:"over_point"`
	OverPrice       string             `json:"over_price" bson:"over_price"`
}

func (s *GameServer) Start(api fiber.Router) error {
	s.StartGameServerAPI(api)
	return nil
}
