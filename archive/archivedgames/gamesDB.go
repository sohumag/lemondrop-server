package archivedGames

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (g *GameServer) GetUpcomingGamesBySport(sport string) ([]Game, error) {
	coll := g.client.Database("games-db").Collection("games")
	filter := bson.M{"commencetime": bson.M{"$gt": time.Now(), "$lt": time.Now().Add(time.Hour * 36)}, "sportkey": sport}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	games := []Game{}
	if err = cursor.All(context.TODO(), &games); err != nil {
		return nil, err
	}

	if len(games) == 0 {
		return []Game{}, nil
	}

	return games, nil
}

func (g *GameServer) GetAllUpcomingGames(maxNum int) ([]Game, error) {
	internalMaxNum := 60
	coll := g.client.Database("games-db").Collection("games")

	if maxNum > internalMaxNum {
		maxNum = internalMaxNum
	}

	filter := bson.M{
		"commencetime": bson.M{"$gt": time.Now(), "$lt": time.Now().Add(time.Hour * 36)},
	}

	opts := options.Find().SetLimit(int64(maxNum))
	cursor, err := coll.Find(context.TODO(), filter, opts)
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

	// find if game already exists, if does then update
	var needsUpdateGame Game
	err := coll.FindOne(context.TODO(), bson.D{{Key: "gameid", Value: game.GameId}}).Decode(&needsUpdateGame)

	// if game doesn't exist
	if err != nil {
		result, err := coll.InsertOne(context.TODO(), *game)
		if err != nil {
			return err
		}

		fmt.Printf("Added game with id: %v\n", result.InsertedID)
		return nil
	}

	// if game does exist, update
	filter := bson.D{{Key: "gameid", Value: game.GameId}}
	update := bson.M{"$set": bson.M{"bookmakers": game.Bookmakers, "scores": game.Scores, "completed": game.Completed, "lastupdate": game.LastUpdate}}

	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	fmt.Printf("Modified %v games\n", result.ModifiedCount)

	return nil
}
