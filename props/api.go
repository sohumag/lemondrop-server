package props

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PropServer struct {
	client *mongo.Client
}

func NewPropServer() *PropServer {
	return &PropServer{
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

func (s *PropServer) StartPropServerAPI(api fiber.Router) error {
	log.Println("Adding prop server endpoints to API")

	return nil
}

func (s *PropServer) Start(api fiber.Router) error {
	s.StartPropServerAPI(api)
	return nil
}

// type AllGamesByLeague struct {
// 	League string
// 	Date   string
// 	Games  []*Game
// }

// type Game struct {
// 	League       string   `json:"league"`
// 	Date         string   `json:"date"`
// 	GameId       string   `json:"game_id"`
// 	AwayTeam     string   `json:"away_team"`
// 	HomeTeam     string   `json:"home_team"`
// 	StartTime    string   `json:"start_timestamp"`
// 	Participants []string `json:"participants"`
// 	Markets      Markets  `json:"markets"`
// }

// type Market struct {
// 	Name       string     `json:"name"`
// 	MarketInfo MarketInfo `json:"market_info"`
// }

// type Markets struct {
// 	GameId  string   `json:"game_id"`
// 	Markets []Market `json:"markets"`
// }

// type MarketInfo struct {
// 	GameId      string     `json:"game_id"`
// 	Sportsbooks []BookInfo `json:"sportsbooks"`
// }

// type BookInfo struct {
// 	BookieKey string        `json:"bookie_key"`
// 	Market    MarketOptions `json:"market"`
// }

// type MarketOptions struct {
// 	MarketKey string    `json:"market_key"`
// 	Outcomes  []Outcome `json:"outcomes"`
// }

// type Outcome struct {
// 	Timestamp   string  `json:"timestamp"`
// 	Handicap    float64 `json:"handicap"`
// 	Odds        int64   `json:"odds"`
// 	Participant int64   `json:"participant"`
// 	Name        string  `json:"name"`
// 	Description string  `json:"description"`
// }

// func (s *PropServer) GetAllPropsGames() error {
// 	// sports := []string{"nba", "nfl", "nhl", "mlb", "ncaaf", "wnba", "atp", "wta", "itf"}
// 	sports := []string{"nfl"}
// 	for _, sport := range sports {
// 		reqEndpoint := fmt.Sprintf("https://api.prop-odds.com/beta/games/%s", sport)
// 		apiKey := "d2RRm5vxSTGMNmzzoDeshq7KC33Gs2Igby5070s9tI"
// 		date := "2023-10-22"
// 		reqUrl := fmt.Sprintf("%v?api_key=%v&date=%v", reqEndpoint, apiKey, date)

// 		// get all games
// 		res, err := http.Get(reqUrl)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		}

// 		defer res.Body.Close()

// 		bytes, err := io.ReadAll(res.Body)

// 		allGames := AllGamesByLeague{}
// 		json.Unmarshal(bytes, &allGames)

// 		for _, game := range allGames.Games {
// 			// for i := 0; i < 1; i++ {
// 			// 	game := allGames.Games[i]
// 			gameId := game.GameId
// 			game.League = allGames.League
// 			game.Date = allGames.Date
// 			reqEndpoint = fmt.Sprintf("https://api.prop-odds.com/beta/markets/%s", gameId)
// 			reqUrl = fmt.Sprintf("%v?api_key=%v", reqEndpoint, apiKey)

// 			// fmt.Println(reqUrl)

// 			res, err = http.Get(reqUrl)
// 			if err != nil {
// 				fmt.Println(err.Error())
// 			}
// 			defer res.Body.Close()

// 			bytes, _ := io.ReadAll(res.Body)

// 			markets := Markets{}
// 			json.Unmarshal(bytes, &markets)

// 			game.Markets = markets

// 			for _, market := range markets.Markets {
// 				reqUrl = fmt.Sprintf("https://api.prop-odds.com/beta/odds/%s/%s?api_key=%s", gameId, market.Name, apiKey)
// 				res, err := http.Get(reqUrl)
// 				if err != nil {
// 					fmt.Println(err.Error())
// 				}
// 				defer res.Body.Close()

// 				bytes, _ := io.ReadAll(res.Body)
// 				// fmt.Println(string(bytes))

// 				// cannot unmarshal data into bookinfo
// 				marketInfo := MarketInfo{}
// 				json.Unmarshal(bytes, &marketInfo)
// 				fmt.Println(marketInfo)
// 				market.MarketInfo = marketInfo
// 			}
// 		}

// 		for _, game := range allGames.Games {
// 			conn := s.client.Database("props-db").Collection(game.League)

// 			result, err := conn.InsertOne(context.TODO(), game)
// 			if err != nil {
// 				log.Fatal(err)
// 			}

// 			fmt.Printf("Inserted game(propped) with id %v\n", result.InsertedID)
// 		}
// 	}
// 	return nil
// }
