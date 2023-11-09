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
	updateLog map[string]time.Time
	gameCache map[string][]ParsedGame
	// needs to take game id and prop name and return prop options
	// propsCache map[string]map[string][]string
}

type ParsedGame struct {
	Id           string    `json:"id" bson:"game_id"`
	SportKey     string    `json:"sport_key" bson:"sport_key"`
	SportTitle   string    `json:"sport_title" bson:"sport_title"`
	CommenceTime time.Time `json:"commence_time" bson:"commence_time"`
	HomeTeam     string    `json:"home_team" bson:"home_team"`
	AwayTeam     string    `json:"away_team" bson:"away_team"`

	LastUpdate string `json:"last_update" bson:"last_update"`

	MoneylinesExist bool `json:"moneylines_exist"`
	SpreadsExist    bool `json:"spreads_exist"`
	TotalsExist     bool `json:"totals_exist"`

	HomeMoneylinePrice  float64 `json:"home_moneyline_price"`
	AwayMoneylinePrice  float64 `json:"away_moneyline_price"`
	DrawMoneylineExists bool    `json:"draw_moneyline_exists"`
	DrawMoneylinePrice  float64 `json:"draw_moneyline_price"`

	HomeSpreadPrice float64 `json:"home_spread_price"`
	AwaySpreadPrice float64 `json:"away_spread_price"`
	HomeSpreadPoint float64 `json:"home_spread_point"`
	AwaySpreadPoint float64 `json:"away_spread_point"`

	OverPrice  float64 `json:"over_price"`
	UnderPrice float64 `json:"under_price"`
	OverPoint  float64 `json:"over_point"`
	UnderPoint float64 `json:"under_point"`
}

type Game struct {
	Id           string      `json:"id" bson:"game_id"`
	SportKey     string      `json:"sport_key" bson:"sport_key"`
	SportTitle   string      `json:"sport_title" bson:"sport_title"`
	CommenceTime time.Time   `json:"commence_time" bson:"commence_time"`
	HomeTeam     string      `json:"home_team" bson:"home_team"`
	AwayTeam     string      `json:"away_team" bson:"away_team"`
	Bookmakers   []Bookmaker `json:"bookmakers" bson:"bookmakers"`
	PropMarkets  []Market    `json:"prop_markets" bson:"prop_markets"`
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
	Name        string  `json:"name" bson:"name"`
	Description string  `json:"description" bson:"description"`
	Price       float64 `json:"price" bson:"price"`
	Point       float64 `json:"point" bson:"point"`
}
