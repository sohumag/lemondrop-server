package games

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (c *Cache) TimeSinceLastUpdate(league string) (time.Duration, error) {
	if err := ValidateLeagueExists(league); err != nil {
		return time.Since(time.Now()), err
	}
	lastUpdate := c.updateLog[league]

	timeSince := time.Since(lastUpdate)
	return timeSince, nil
}

func (s *GameServer) UpdateGameCacheAndLog(league string, games []Game) error {
	// only need max 3 bookmakers
	maxBookmakers := 1
	for i, game := range games {
		newGame := game
		if len(game.Bookmakers) > maxBookmakers {
			newGame.Bookmakers = game.Bookmakers[0:maxBookmakers]
		}

		games[i] = newGame
	}

	s.cache.gameCache[league] = games
	s.cache.updateLog[league] = time.Now()
	return nil
}

func (s *GameServer) GetAllGamesAndLogs() error {
	for _, leagueName := range validLeagues {
		fmt.Printf("getting %v\n", leagueName)
		games, err := s.GetAllGamesByLeague(leagueName)
		if err != nil {
			return err
		}

		s.AddNewGamesToDB(games)
		s.UpdateGameCacheAndLog(leagueName, games)
	}

	return nil
}

func (s *GameServer) GetAllGamesByLeague(league string) ([]Game, error) {
	apiKey := os.Getenv("ODDS_API_KEY")
	reqUrl := fmt.Sprintf("https://api.the-odds-api.com/v4/sports/%s/odds?api_key=%s&regions=us&markets=h2h,totals,spreads", league, apiKey)
	res, err := http.Get(reqUrl)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	games := []Game{}
	json.Unmarshal(bytes, &games)

	return games, nil
}

func (s *GameServer) AddNewGamesToDB(games []Game) error {
	// duplication error
	coll := s.client.Database("backup").Collection("games")

	// if speed is issue, can do concurrently
	for _, game := range games {
		res, err := coll.InsertOne(context.TODO(), game)
		if err != nil {
			return err
		}

		fmt.Printf("Inserted game with id: %v\n", res.InsertedID)
	}

	return nil
}

// use mongo to init in memory store and update on request to api
func (s *GameServer) InitGamesAndLogs() error {
	fmt.Println("Adding all games from database to in memory store and updating logs...")
	coll := s.client.Database("backup").Collection("games")

	totalGames := 0
	for _, league := range validLeagues {
		filter := bson.M{"sport_key": league, "commence_time": bson.M{"$gt": time.Now()}}

		cursor, err := coll.Find(context.TODO(), filter)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		games := []Game{}
		if err = cursor.All(context.TODO(), &games); err != nil {
			fmt.Println(err.Error())
			return err
		}

		s.UpdateGameCacheAndLog(league, games)

		totalGames += len(games)
	}

	fmt.Printf("Added %d games to in memory store\n", totalGames)
	return nil
}

// !!!! DANGER
func (s *GameServer) ClearDatabase() error {
	// coll := s.client.Database("backup").Collection("games")

	// filter := bson.D{{}}
	// res, err := coll.DeleteMany(context.TODO(), filter)

	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// fmt.Printf("Deleted %d games from database\n", res.DeletedCount)
	// return err

	return nil
}
