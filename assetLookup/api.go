package assetLookup

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AssetServer struct {
	client *mongo.Client
}

func NewAssetServer() *AssetServer {
	return &AssetServer{
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
func (s *AssetServer) StartAssetServerAPI(api fiber.Router) error {
	log.Println("Adding asset lookup server endpoints to API")

	assetsApi := api.Group("/assets")
	assetsApi.Get("/:league/:teamName", func(c *fiber.Ctx) error {
		leagueParsed := ConvertRawString(c.Params("league"))
		teamName := ConvertRawString(c.Params("teamName"))

		if leagueParsed == "nfl" || leagueParsed == "nba" || leagueParsed == "nhl" {
			url := fmt.Sprintf("https://assets.sportsbook.fanduel.com/images/team/%s/%s.png", leagueParsed, teamName)
			c.Write([]byte(url))
		} else if leagueParsed == "ncaaf" {
			url := s.GetCollegeFootballLogo()
			c.Write([]byte(url))
		} else if leagueParsed == "ncaab" {
			url := s.GetCollegeBasketballLogo()
			c.Write([]byte(url))
		} else {
			url := s.GetSoccerLogo()
			c.Write([]byte(url))
		}

		return nil
	})
	return nil
}

func (s *AssetServer) GetSoccerLogo() string {
	urls := []string{
		"https://assets.sportsbook.fanduel.com/images/team/soccer/epl/west_ham.png",
		"https://assets.sportsbook.fanduel.com/images/team/soccer/epl/everton.png",
		"https://assets.sportsbook.fanduel.com/images/team/soccer/epl/brighton.png",
		"https://assets.sportsbook.fanduel.com/images/team/soccer/epl/fulham.png",
		"https://assets.sportsbook.fanduel.com/images/team/soccer/epl/aston_villa.png",
		"https://assets.sportsbook.fanduel.com/images/team/soccer/epl/luton.png",
		"https://assets.sportsbook.fanduel.com/images/team/soccer/epl/liverpool.png",
		"https://assets.sportsbook.fanduel.com/images/team/soccer/epl/nottingham_forest.png",
		"https://assets.sportsbook.fanduel.com/images/team/soccer/epl/man_utd.png",
		"https://assets.sportsbook.fanduel.com/images/team/soccer/epl/man_city.png",
	}

	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(urls)
	return urls[n]
}

func (s *AssetServer) GetCollegeFootballLogo() string {
	urls := []string{
		"https://assets.sportsbook.fanduel.com/images/team/ncaaf/san_jose_state.png",
		"https://assets.sportsbook.fanduel.com/images/team/ncaaf/hawaii.png",
		"https://assets.sportsbook.fanduel.com/images/team/ncaaf/northern_illinois_huskies.png",
		"https://assets.sportsbook.fanduel.com/images/team/ncaaf/central_michigan.png",
		"https://assets.sportsbook.fanduel.com/images/team/ncaaf/buffalo.png",
		"https://assets.sportsbook.fanduel.com/images/team/ncaaf/toledo_rockets.png",
		"https://assets.sportsbook.fanduel.com/images/team/ncaaf/kent_state.png",
		"https://assets.sportsbook.fanduel.com/images/team/ncaaf/akron_zips.png",
		"https://assets.sportsbook.fanduel.com/images/team/ncaaf/penn_state_lions.png",
		"https://assets.sportsbook.fanduel.com/images/team/ncaaf/maryland.png",
	}

	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(urls)
	return urls[n]
}

func (s *AssetServer) GetCollegeBasketballLogo() string {
	urls := []string{
		"https://assets.sportsbook.fanduel.com/images/team/ncaab/lsu.png",
		"https://assets.sportsbook.fanduel.com/images/team/ncaab/connecticut.png",
		"https://assets.sportsbook.fanduel.com/images/team/ncaab/iowa.png",
		"https://assets.sportsbook.fanduel.com/images/team/ncaab/virginia_tech.png",
		"https://assets.sportsbook.fanduel.com/images/team/ncaab/utah.png",
		"https://assets.sportsbook.fanduel.com/images/team/ncaab/indiana.png",
		"https://assets.sportsbook.fanduel.com/images/team/ncaab/south_carolina.png",
		"https://assets.sportsbook.fanduel.com/images/team/ncaab/stanford.png",
		"https://assets.sportsbook.fanduel.com/images/team/ncaab/north_carolina.png",
		"https://assets.sportsbook.fanduel.com/images/team/ncaab/duke.png",
	}

	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(urls)
	return urls[n]
}

func (s *AssetServer) Start(api fiber.Router) error {
	s.StartAssetServerAPI(api)
	return nil
}

func ConvertRawString(raw string) string {
	// trim, lowercase, replace
	new := strings.Trim(raw, " ")
	new = strings.ToLower(new)
	new = strings.Replace(new, " ", "_", -1)
	new = strings.Replace(new, "%20", "_", -1)
	return new
}

// func (s *AssetServer) ConvertToleague, team string) string {
// 	convertedLeague := ConvertRawString(league)
// 	convertedTeam := ConvertRawString(team)
// 	if convertedLeague == "nba" || convertedTeam == "nba" {
// 		return GetUrlForProe
// 	}
// }

// func (s *AssetServer) GetUrlForProfessional(league, team string ) string {
// 	url := fmt.Sprintf("https://assets.sportsbook.fanduel.com/images/team/%s/%s.png", league, team)
// 	return url
// }
