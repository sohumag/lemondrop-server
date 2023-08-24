package games

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewGameServer(listenAddr string) *GameServer {
	// g := &GameServer{port: listenAddr}

	return &GameServer{
		port:   listenAddr,
		client: ConnectDB(),
	}
}

func ConnectDB() *mongo.Client {
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	// defer func() {
	// 	if err := client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	return client
}

func (g *GameServer) Start() error {
	go g.StartAPI()

	// auto update all new games and scores for games
	for {
		for i := 0; i < 4; i++ {
			time.Sleep(time.Hour * 6)
			g.UpdateGameScores()

		}

		g.MigrateAllGames()

	}

}
