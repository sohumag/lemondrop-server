package games

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewGameServer() *GameServer {
	return &GameServer{
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

func (g *GameServer) Start(api fiber.Router) error {
	// go func() {
	// 	for {
	// 		time.Sleep(time.Hour * 4)
	// 		g.MigrateAllGames()
	// 		g.UpdateGameScores()
	// 	}
	// }()

	// go g.MigrateAllGames()

	g.StartGameServerAPI(api)

	return nil
}
