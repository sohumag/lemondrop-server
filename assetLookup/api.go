package assetLookup

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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
		url := fmt.Sprintf("https://assets.sportsbook.fanduel.com/images/team/%s/%s.png", leagueParsed, teamName)
		c.Write([]byte(url))
		return nil
	})
	return nil
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

// func (s *AssetServer) ConvertToUrl(league, team string) string {
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
