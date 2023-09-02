package games

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (g *GameServer) UpdateGameScores() error {
	coll := g.client.Database("games-db").Collection("games")
	cursor, err := coll.Find(context.TODO(), bson.M{
		"commencetime": bson.M{"$lt": time.Now(), "$gt": time.Now().AddDate(0, 0, -3)},
	})
	if err != nil {
		return err
	}

	games := []Game{}
	for err := cursor.All(context.TODO(), &games); err != nil; {
		return err
	}

	sportGameIdMap := map[string][]string{}

	for _, game := range games {

		if _, ok := sportGameIdMap[game.SportKey]; !ok {
			sportGameIdMap[game.SportKey] = []string{game.GameId}
		} else {
			sportGameIdMap[game.SportKey] = append(sportGameIdMap[game.SportKey], game.GameId)
		}

	}

	for sport, ids := range sportGameIdMap {
		joinedGameIds := strings.Join(ids, ",")
		reqUrl := fmt.Sprintf("https://api.the-odds-api.com/v4/sports/%s/scores/?daysFrom=3&apiKey=e0ae2e9cd2c145da9659ce53ddbc4442&eventIds=%s", sport, joinedGameIds)
		fmt.Println(reqUrl)

		res, err := http.Get(reqUrl)

		if err != nil {
			return err
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		scoredGames := []ScoreAPIReturn{}
		json.Unmarshal(body, &scoredGames)

		// update completed, scores, and lastupdate of game id
		for _, game := range scoredGames {

			filter := bson.D{{Key: "gameid", Value: game.GameId}}
			update := bson.M{"$set": bson.M{"scores": game.Scores, "lastupdate": game.LastUpdate, "completed": game.Completed}}

			result, err := coll.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				return err
			}

			fmt.Printf("Modified %v documents\n", result.ModifiedCount)

		}
	}

	return nil
}
