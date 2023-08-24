package games

import (
	"context"
	"fmt"
)

func (g *GameServer) AddSportToDB(sport *Sport) error {
	coll := g.client.Database("games-db").Collection("sports")
	result, err := coll.InsertOne(context.TODO(), *sport)
	if err != nil {
		return err
	}

	fmt.Printf("Added sport with id: %v\n", result.InsertedID)

	return nil
}

func (g *GameServer) AddGameToDB(game *Game) error {
	coll := g.client.Database("games-db").Collection("games")
	result, err := coll.InsertOne(context.TODO(), *game)
	if err != nil {
		return err
	}

	fmt.Printf("Added game with id: %v\n", result.InsertedID)
	return nil
}
