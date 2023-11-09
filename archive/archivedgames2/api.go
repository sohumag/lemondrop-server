package games

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*

NOTES: should automatically send whatever is in cache, and then async fetch new data and send on next request. reduces latency significatnly
- bloom filter?
*/

func ReturnAllLeagues() []string {
	return validLeagues
}

var validLeagues = []string{
	"americanfootball_nfl",
	"americanfootball_ncaaf",

	"basketball_nba",
	"basketball_ncaab",

	"icehockey_nhl",

	"soccer_uefa_champs_league",
}

func NewGameServer() *GameServer {
	log := map[string]time.Time{}
	cache := map[string][]ParsedGame{}

	for _, leagueName := range validLeagues {
		log[leagueName] = time.Now()
		//? initializing empty cache to start. maybe change later
		cache[leagueName] = []ParsedGame{}
	}

	return &GameServer{
		client: ConnectDB(),
		cache: Cache{
			updateLog: log,
			gameCache: cache,
		},
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

func (s *GameServer) StartGameServerAPI(api fiber.Router) error {
	log.Println("Adding game caching server endpoints to API")

	// Games
	gamesApi := api.Group("/games")

	// /props/:propname
	// needs to use straight mongo bc cache will not be large enough
	gamesApi.Get("/league/:league/game/:id/props/:prop", func(c *fiber.Ctx) error {
		return s.ReturnPropsByIdAndPropName(c, c.Params("league"), c.Params("id"), c.Params("prop"))
	})

	gamesApi.Get("/game/:id", func(c *fiber.Ctx) error {
		return s.ReturnGameById(c, c.Params("id"))
	})

	// TODO
	gamesApi.Get("/league/:league", func(c *fiber.Ctx) error {
		return s.CacheAndReturnGamesByLeague(c, c.Params("league"))
	})

	// TODO: Add /all endpoint

	gamesApi.Get("/sport/:sport", func(c *fiber.Ctx) error {
		return s.CacheAndReturnGamesBySport(c, c.Params("sport"))
	})

	return nil
}

func (s *GameServer) CacheAndReturnGamesBySport(c *fiber.Ctx, sport string) error {
	sportLeagueMap := map[string][]string{
		"football":   {"americanfootball_nfl", "americanfootball_ncaaf"},
		"basketball": {"basketball_nba", "basketball_ncaaf"},
		"hockey":     {"icehockey_nhl"},
		"soccer":     {"soccer_uefa_champions_leauge"},
	}

	if _, ok := sportLeagueMap[sport]; !ok {
		// sport doesnt exist
		c.SendStatus(http.StatusBadRequest)
		return nil
	}

	leagues := sportLeagueMap[sport]
	allGames := []ParsedGame{}
	for _, league := range leagues {
		rawGames := s.cache.gameCache[league]
		for _, game := range rawGames {
			allGames = append(allGames, game)
		}
	}

	c.JSON(allGames)

	return nil
}

func (s *GameServer) CacheAndReturnGamesByLeague(c *fiber.Ctx, league string) error {
	maxTimeBeforeUpdate := 24 * time.Hour
	if err := ValidateLeagueExists(league); err != nil {
		c.Send([]byte(err.Error()))
	}

	// return all games from cache immediately
	curGames := s.cache.gameCache[league]
	c.JSON(curGames)

	// check if too much time has passed
	// -> return if it has not
	// -> update if has not
	timePassed, err := s.cache.TimeSinceLastUpdate(league)
	if err != nil {
		return err
	}

	// does not need to update
	if timePassed < maxTimeBeforeUpdate {
		return nil
	}

	fmt.Println("Max time has passed. updating cache and logs")
	games, err := s.GetAllGamesByLeague(league)
	if err != nil {
		return err
	}
	s.AddNewGamesToDB(games)
	s.UpdateGameCacheAndLog(league, games)

	return nil
}

func (s *GameServer) Start(api fiber.Router) error {
	// s.ClearDatabase()
	// s.GetAllGamesAndLogs()

	// implement auto updates for mongo fetches

	s.InitGamesAndLogs()
	return s.StartGameServerAPI(api)
}

/*
- Start: Get all games from mongo and put into cache
- Get Games: get all games from api and test if already exists. Update if it does. Add if it doesnt
- Update cache(games)
endpoints:
	- be able to get games by league(simple)
	- get by sport (all leagues...)
	- get all games .. sort by league
*/
