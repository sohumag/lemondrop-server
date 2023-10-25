package games

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type GameServer struct {
	client *mongo.Client
	cache  Cache
}

type Cache struct {
	// map[leagueName] -> lastUpdated
	// map[leagueName] -> []Games
	updateLog map[string]time.Time
	gameCache map[string][]Game
}

type Game struct {
	Id           string      `json:"id" bson:"game_id"`
	SportKey     string      `json:"sport_key" bson:"sport_key"`
	SportTitle   string      `json:"sport_title" bson:"sport_title"`
	CommenceTime time.Time   `json:"commence_time" bson:"commence_time"`
	HomeTeam     string      `json:"home_team" bson:"home_team"`
	AwayTeam     string      `json:"away_team" bson:"away_team"`
	Bookmakers   []Bookmaker `json:"bookmakers" bson:"bookmakers"`
}

type Bookmaker struct {
	Key        string   `json:"key" bson:"key"`
	Title      string   `json:"title" bson:"title"`
	LastUpdate string   `json:"last_update" bson:"last_update"`
	Markets    []Market `json:"markets" bson:"markets"`
}

type Market struct {
	Key        string    `json:"key" bson:"key"`
	LastUpdate string    `json:"last_update" bson:"last_update"`
	Outcomes   []Outcome `json:"outcomes" bson:"outcomes"`
}

type Outcome struct {
	Name  string  `json:"name" bson:"name"`
	Price float64 `json:"price" bson:"price"`
	Point float64 `json:"point" bson:"point"`
}
