package games

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

/*

NOTES: should automatically send whatever is in cache, and then async fetch new data and send on next request. reduces latency significatnly
- bloom filter?
*/

var validLeagues = []string{
	"americanfootball_nfl",
	"americanfootball_ncaaf",

	"basketball_nba",
	"basketball_ncaab",

	"icehockey_nhl",
}

func NewGameServer() *GameServer {
	log := map[string]time.Time{}
	cache := map[string][]Game{}

	for _, leagueName := range validLeagues {
		log[leagueName] = time.Now()
		//? initializing empty cache to start. maybe change later
		cache[leagueName] = []Game{}
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
	gamesApi.Get("/:league", func(c *fiber.Ctx) error {
		return s.CacheAndReturnGamesByLeague(c, c.Params("league"))
	})

	return nil
}

func (s *GameServer) CacheAndReturnGamesByLeague(c *fiber.Ctx, league string) error {
	maxTimeBeforeUpdate := 8 * time.Hour
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
	s.InitGamesAndLogs()
	return s.StartGameServerAPI(api)
}
