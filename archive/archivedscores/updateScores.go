package scores

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var validLeagues = games.ReturnAllLeagues()
var validLeagues = []string{}

// var validLeagues = []string{
// 	"americanfootball_nfl",
// }

func (s *ScoreServer) StartScoresUpdates() error {

	// scheduler := gocron.NewScheduler(time.UTC)
	// _, err := scheduler.Every(1).Day().At("23:59").At("15:00").Do(func() {
	// 	s.UpdateScores()
	// })

	// if err != nil {
	// 	return err
	// }

	return nil
}

func (s *ScoreServer) UpdateScores() error {
	for _, league := range validLeagues {
		// get scores
		// add to db
		games, err := s.GetScoresByLeague(league)
		if err != nil {
			// fmt.Println(err)
			return err
		}

		for _, game := range games {
			fmt.Println(game)
		}

		s.AddScoredGamesToDB(games)
	}

	return nil
}

func (s *ScoreServer) GetScoresByLeague(league string) ([]ScoredGame, error) {
	apiKey := os.Getenv("ODDS_API_KEY")

	url := fmt.Sprintf(
		"https://api.the-odds-api.com/v4/sports/%s/scores/?daysFrom=3&apiKey=%s",
		league,
		apiKey,
	)

	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	games := []ScoredGame{}
	json.Unmarshal(bytes, &games)

	return games, nil
}

func (s *ScoreServer) AddScoredGamesToDB(games []ScoredGame) error {
	coll := s.client.Database("scores-db").Collection("scores")

	for _, game := range games {

		// check if game exists: replace if does
		// using upsert

		filter := bson.D{{Key: "id", Value: game.Id}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "sport_key", Value: game.SportKey}, {Key: "sport_title", Value: game.SportTitle}, {Key: "commence_time", Value: game.CommenceTime}, {Key: "home_team", Value: game.HomeTeam}, {Key: "away_team", Value: game.AwayTeam}, {Key: "completed", Value: game.Completed}, {Key: "last_update", Value: game.LastUpdate}, {Key: "scores", Value: game.Scores}}}}
		opts := options.Update().SetUpsert(true)

		res, err := coll.UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			return err
		}

		fmt.Printf("updated %v docs\n", res.ModifiedCount)

		// if _, err := coll.InsertOne(context.TODO(), game); err != nil {
		// 	fmt.Println("failed to insert game: ", err)
		// 	continue
		// }
		// fmt.Println("inserted game into database")
	}

	return nil
}

type ScoredGame struct {
	Id           string      `json:"id"`
	SportKey     string      `json:"sport_key"`
	SportTitle   string      `json:"sport_title"`
	CommenceTime string      `json:"commence_time"`
	Completed    bool        `json:"completed"`
	HomeTeam     string      `json:"home_team"`
	AwayTeam     string      `json:"away_team"`
	Scores       []TeamScore `json:"scores"`
	LastUpdate   string      `json:"last_update"`
}

type TeamScore struct {
	Name  string `json:"name"`
	Score string `json:"score"`
}
