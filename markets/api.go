package markets

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MarketServer struct {
	client *mongo.Client
}

func NewMarketServer() *MarketServer {
	return &MarketServer{
		client: ConnectDB(),
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

func (s *MarketServer) StartMarketServerAPI(api fiber.Router) error {
	log.Println("Adding market server endpoints to API")
	markets := map[string][]string{
		"americanfootball_nfl":   {"player_pass_tds", "player_pass_yds", "player_pass_completions", "player_pass_attempts", "player_pass_interceptions", "player_pass_longest_completion", "player_rush_yds", "player_rush_attempts", "player_rush_longest", "player_receptions", "player_reception_yds", "player_reception_longest", "player_kicking_points", "player_field_goals", "player_tackels_assists", "player_1st_td", "player_last_td", "player_anytime_td"},
		"americanfootball_ncaaf": {"player_pass_tds", "player_pass_yds", "player_pass_completions", "player_pass_attempts", "player_pass_interceptions", "player_pass_longest_completion", "player_rush_yds", "player_rush_attempts", "player_rush_longest", "player_receptions", "player_reception_yds", "player_reception_longest", "player_kicking_points", "player_field_goals", "player_tackels_assists", "player_1st_td", "player_last_td", "player_anytime_td"},

		"basketball_nba":   {"player_points", "player_rebounds", "player_assists", "player_threes", "player_double_double", "player_blocks", "player_steals", "player_turnovers", "player_points_rebounds_assists", "player_points_rebounds", "player_points_assists", "player_rebounds_assists"},
		"basketball_ncaab": {"player_points", "player_rebounds", "player_assists", "player_threes", "player_double_double", "player_blocks", "player_steals", "player_turnovers", "player_points_rebounds_assists", "player_points_rebounds", "player_points_assists", "player_rebounds_assists"},

		"baseball_mlb": {"batter_home_runs", "batter_hits", "batter_total_bases", "batter_rbis", "batter_runs_scored", "batter_hits_runs_rbis", "batter_singles", "batter_doubles", "batter_triples", "batter_walks", "batter_stolen_bases", "pitcher_strikeouts", "pitcher_record_a_win", "pitcher_hits_allowed", "pitcher_walks", "pitcher_earned_runs", "pitcher_outs"},

		"icehockey_nhl": {"player_points", "player_power_play_points", "player_assists", "player_blocked_shots", "player_shots_on_goal"},
	}

	marketsApi := api.Group("/markets")
	marketsApi.Get("/:market", func(c *fiber.Ctx) error {
		c.JSON(markets[c.Params("market")])
		return nil
	})

	return nil
}

func (s *MarketServer) Start(api fiber.Router) error {
	s.StartMarketServerAPI(api)
	return nil
}
