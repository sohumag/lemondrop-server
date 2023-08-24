package games

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (g *GameServer) GetUpcomingGamesBySport(sport string) ([]Game, error) {
	coll := g.client.Database("games-db").Collection("games")
	cursor, err := coll.Find(context.TODO(), bson.D{{"sporttitle", sport}})
	if err != nil {
		return nil, err
	}

	games := []Game{}
	if err = cursor.All(context.TODO(), &games); err != nil {
		return nil, err
	}

	if len(games) == 0 {
		return nil, fmt.Errorf("invalid sport")
	}

	return games, nil
}

func (g *GameServer) GetAllUpcomingGames(max int) ([]Game, error) {
	internalMaxNum := 25
	coll := g.client.Database("games-db").Collection("games")

	if max > internalMaxNum {
		max = internalMaxNum
	}

	opts := options.Find().SetLimit(int64(max))
	cursor, err := coll.Find(context.TODO(), bson.D{{}}, opts)
	if err != nil {
		return nil, err
	}

	games := []Game{}
	if err = cursor.All(context.TODO(), &games); err != nil {
		return nil, err
	}

	return games, nil
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
