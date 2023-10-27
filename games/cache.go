package games

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func (s *GameServer) ParseGame(game Game) (ParsedGame, error) {
	//! IMPROVEMENT: IF ONE MARKET DOESNT EXIST, CHECK OTHER BOOKS FOR IT
	// get draftkings bookmaker if exists
	// if doesnt exist, take first bookmaker

	// check what markets exist
	// Fill in fields based on markets that exist

	parsedGame := ParsedGame{
		Id:                  game.Id,
		SportKey:            game.SportKey,
		SportTitle:          game.SportTitle,
		CommenceTime:        game.CommenceTime,
		HomeTeam:            game.HomeTeam,
		AwayTeam:            game.AwayTeam,
		MoneylinesExist:     false,
		SpreadsExist:        false,
		TotalsExist:         false,
		DrawMoneylineExists: false,
	}

	// choosing book to use
	maxMarkets := 0
	usedBook := Bookmaker{}
	for _, book := range game.Bookmakers {
		if len(book.Markets) > maxMarkets {
			maxMarkets = len(book.Markets)
			usedBook = book
		}
	}

	// parsing book
	for _, market := range usedBook.Markets {
		parsedGame.LastUpdate = market.LastUpdate
		if market.Key == "h2h" {
			parsedGame.MoneylinesExist = true

			if market.Outcomes[0].Name == game.HomeTeam {
				parsedGame.HomeMoneylinePrice = market.Outcomes[0].Price
				parsedGame.AwayMoneylinePrice = market.Outcomes[1].Price
			} else {
				parsedGame.HomeMoneylinePrice = market.Outcomes[1].Price
				parsedGame.AwayMoneylinePrice = market.Outcomes[0].Price
			}

			if len(market.Outcomes) >= 3 {
				parsedGame.DrawMoneylineExists = true
				parsedGame.DrawMoneylinePrice = market.Outcomes[2].Price
			}

		} else if market.Key == "spreads" {
			parsedGame.SpreadsExist = true

			if market.Outcomes[0].Name == game.HomeTeam {
				parsedGame.HomeSpreadPoint = market.Outcomes[0].Point
				parsedGame.HomeSpreadPrice = market.Outcomes[0].Price
				parsedGame.AwaySpreadPoint = market.Outcomes[1].Point
				parsedGame.AwaySpreadPrice = market.Outcomes[1].Price
			} else {
				parsedGame.HomeSpreadPoint = market.Outcomes[1].Point
				parsedGame.HomeSpreadPrice = market.Outcomes[1].Price
				parsedGame.AwaySpreadPoint = market.Outcomes[0].Point
				parsedGame.AwaySpreadPrice = market.Outcomes[0].Price
			}

		} else if market.Key == "totals" {
			parsedGame.TotalsExist = true
			parsedGame.OverPoint = market.Outcomes[0].Point
			parsedGame.OverPrice = market.Outcomes[0].Price
			parsedGame.UnderPoint = market.Outcomes[1].Point
			parsedGame.UnderPrice = market.Outcomes[1].Price
		}
	}

	return parsedGame, nil
}

func (s *GameServer) UpdateGameCacheAndLog(league string, games []Game) error {
	// only need max 3 bookmakers
	// only want draftkings odds bc others dont always have all markets
	parsedGames := []ParsedGame{}
	// maxBookmakers := 1
	for _, game := range games {
		parsedGame, err := s.ParseGame(game)
		if err != nil {
			fmt.Println(err)
			return err
		}

		parsedGames = append(parsedGames, parsedGame)
	}

	s.cache.gameCache[league] = parsedGames
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
	reqUrl := fmt.Sprintf("https://api.the-odds-api.com/v4/sports/%s/odds?api_key=%s&regions=us&markets=h2h,totals,spreads&oddsFormat=american", league, apiKey)
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
		// check if game exists: if so delete
		filter := bson.M{"game_id": game.Id}
		if _, err := coll.DeleteOne(context.TODO(), filter); err != nil {
			log.Println(err.Error())
		}

		res, err := coll.InsertOne(context.TODO(), game)
		if err != nil {
			fmt.Println(err.Error())
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
		fmt.Println(league)
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

		// for _, game := range games {
		// 	fmt.Println(game.SportTitle)
		// }

		s.UpdateGameCacheAndLog(league, games)

		totalGames += len(games)
	}

	fmt.Printf("Added %d games to in memory store\n", totalGames)
	return nil
}

// !!!! DANGER
func (s *GameServer) ClearDatabase() error {
	coll := s.client.Database("backup").Collection("games")

	filter := bson.D{{}}
	res, err := coll.DeleteMany(context.TODO(), filter)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("Deleted %d games from database\n", res.DeletedCount)
	return err
}
