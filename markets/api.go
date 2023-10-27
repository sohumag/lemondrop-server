package markets

import (
	"context"
	"log"
	"net/http"
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

type Markets struct {
	Name       string   `json:"name"`
	AllMarkets []string `json:"all_markets"`
}

func (s *MarketServer) StartMarketServerAPI(api fiber.Router) error {
	alternates := []string{"alternate_spreads", "alternate_totals"}

	footballPlayerProps := []string{"player_pass_tds", "player_pass_yds", "player_pass_completions", "player_pass_attempts", "player_pass_interceptions", "player_pass_longest_completion", "player_rush_yds", "player_rush_attempts", "player_rush_longest", "player_receptions", "player_reception_yds", "player_reception_longest", "player_kicking_points", "player_field_goals", "player_tackels_assists", "player_1st_td", "player_last_td", "player_anytime_td"}
	basketballPlayerProps := []string{"player_points", "player_rebounds", "player_assists", "player_threes", "player_double_double", "player_blocks", "player_steals", "player_turnovers", "player_points_rebounds_assists", "player_points_rebounds", "player_points_assists", "player_rebounds_assists"}
	hockeyPlayerProps := []string{"player_points", "player_power_play_points", "player_assists", "player_blocked_shots", "player_shots_on_goal"}

	footballMlProps := []string{"h2h_q1", "h2h_q2", "h2h_q3", "h2h_q4", "h2h_h1", "h2h_h2"}
	footballSpreadProps := []string{"spreads_q1", "spreads_q2", "spreads_q3", "spreads_q4", "spreads_h1", "spreads_h2"}
	footballTotalsProps := []string{"totals_q1", "totals_q2", "totals_q3", "totals_q4", "totals_h1", "totals_h2"}

	hockeyMlProps := []string{"h2h_p1", "h2h_p2", "h2h_p3"}
	hockeySpreadsProps := []string{"spreads_p1", "spreads_p2", "spreads_p3"}
	hockeyTotalsProps := []string{"totals_p1", "totals_p2", "totals_p3"}

	footballMarket := []Markets{
		{
			Name:       "Alternates",
			AllMarkets: alternates,
		},

		{
			Name:       "Player Props",
			AllMarkets: footballPlayerProps,
		},
		{
			Name:       "Moneyline Props",
			AllMarkets: footballMlProps,
		},
		{
			Name:       "Spread Props",
			AllMarkets: footballSpreadProps,
		},
		{
			Name:       "Totals Props",
			AllMarkets: footballTotalsProps,
		},
	}

	basketballMarket := []Markets{
		{
			Name:       "Alternates",
			AllMarkets: alternates,
		},

		{
			Name:       "Player Props",
			AllMarkets: basketballPlayerProps,
		},
		{
			Name:       "Moneyline Props",
			AllMarkets: footballMlProps,
		},
		{
			Name:       "Spread Props",
			AllMarkets: footballSpreadProps,
		},
		{
			Name:       "Totals Props",
			AllMarkets: footballTotalsProps,
		},
	}

	hockeyMarket := []Markets{
		{
			Name:       "Alternates",
			AllMarkets: alternates,
		},

		{
			Name:       "Player Props",
			AllMarkets: hockeyPlayerProps,
		},
		{
			Name:       "Moneyline Props",
			AllMarkets: hockeyMlProps,
		},
		{
			Name:       "Spread Props",
			AllMarkets: hockeySpreadsProps,
		},
		{
			Name:       "Totals Props",
			AllMarkets: hockeyTotalsProps,
		},
	}

	soccerMarket := []Markets{
		{
			Name:       "Alternates",
			AllMarkets: alternates,
		},
	}

	log.Println("Adding market server endpoints to API")
	markets := map[string][]Markets{
		"americanfootball_nfl":   footballMarket,
		"americanfootball_ncaaf": footballMarket,

		"basketball_nba":   basketballMarket,
		"basketball_ncaab": basketballMarket,

		"icehockey_nhl":             hockeyMarket,
		"soccer_uefa_champs_league": soccerMarket,
	}

	marketsApi := api.Group("/markets")
	marketsApi.Get("/:market", func(c *fiber.Ctx) error {
		if _, ok := markets[c.Params("market")]; !ok {
			c.SendStatus(http.StatusNotFound)
		}
		c.JSON(markets[c.Params("market")])
		return nil
	})

	return nil
}

func (s *MarketServer) Start(api fiber.Router) error {
	s.StartMarketServerAPI(api)
	return nil
}
