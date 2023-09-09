package games

import (
	"context"
	"os"

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

	return client
}

func (g *GameServer) Start() error {
	// go func() {
	// 	for {
	// 		time.Sleep(time.Hour * 4)
	// 		g.MigrateAllGames()
	// 		g.UpdateGameScores()
	// 	}
	// }()

	// go g.MigrateAllGames()
	g.StartAPI()
	return nil
}
