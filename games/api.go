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

	gamesApi.Get("/leagues", func(c *fiber.Ctx) error {
		return s.GetAllLeagues(c)
	})

	gamesApi.Get("/sports", func(c *fiber.Ctx) error {
		return s.GetAllSports(c)
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

	picksApi := api.Group("/picks")
	picksApi.Get("/", func(c *fiber.Ctx) error {
		c.Send([]byte("picks api is working"))
		return nil
	})

	picksApi.Get("/:league", func(c *fiber.Ctx) error {
		return s.GetPicksByLeagueId(c)
	})

	picksApi.Get("/:league/markets", func(c *fiber.Ctx) error {
		return s.GetMarketsByLeagueId(c)
	})

	return nil
}

func (s *GameServer) GetPicksByLeagueId(c *fiber.Ctx) error {
	coll := s.client.Database("games-db").Collection("scraped-picks")
	// filter := bson.D{{Key: "league_id", Value: strings.ToLower(c.Params("league"))}}
	currentDate := time.Now()
	maxDate := currentDate.Add(time.Hour * 24 * 7)
	filter := bson.M{"start_date": bson.M{"$gt": currentDate, "$lt": maxDate}, "league_id": strings.ToLower(c.Params("league"))}
	opts := options.Find().SetSort(bson.D{{Key: "market", Value: -1}})
	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return err
	}

	picks := []Pick{}
	pick := Pick{}
	for cursor.Next(context.TODO()) {
		cursor.Decode(&pick)
		picks = append(picks, pick)
	}

	allPicks := map[string][]Pick{}
	for _, p := range picks {
		if _, ok := allPicks[p.Market]; !ok {
			allPicks[p.Market] = []Pick{}
		}
		allPicks[p.Market] = append(allPicks[p.Market], p)
	}

	c.JSON(allPicks)

	return nil
}

func (s *GameServer) GetMarketsByLeagueId(c *fiber.Ctx) error {
	coll := s.client.Database("games-db").Collection("scraped-picks")
	filter := bson.D{{Key: "league_id", Value: strings.ToLower(c.Params("league"))}}
	results, err := coll.Distinct(context.TODO(), "market", filter)
	if err != nil {
		return err
	}

	c.JSON(results)

	return nil
}

func (s *GameServer) GetGameById(c *fiber.Ctx) error {
	coll := s.client.Database("games-db").Collection("scraped-games")
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	game := Game{}
	coll.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&game)

	c.JSON(game)

	return nil
}

func (s *GameServer) GetAllSports(c *fiber.Ctx) error {
	sports := []Sport{
		{Name: "Basketball"},
		{Name: "Football"},
		{Name: "Ice Hockey"},
		{Name: "Soccer"},
	}

	c.JSON(sports)

	return nil
}

func (s *GameServer) GetAllLeagues(c *fiber.Ctx) error {
	coll := s.client.Database("games-db").Collection("leagues")
	cursor, err := coll.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	leagues := []League{}
	league := League{}
	for cursor.Next(context.TODO()) {
		cursor.Decode(&league)
		leagues = append(leagues, league)
	}

	c.JSON(leagues)

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

func (s *GameServer) Start(api fiber.Router) error {
	s.StartGameServerAPI(api)
	return nil
}
