package games

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Game struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	GameType    string             `json:"game_type" bson:"game_type"`
	League      string             `json:"league" bson:"league" `
	LeagueId    string             `json:"league_id" bson:"league_id"`
	Sport       string             `json:"sport" bson:"sport"`
	StartDate   time.Time          `json:"start_date" bson:"start_date"`
	LastUpdated time.Time          `json:"last_updated" bson:"last_updated"`
	Hash        string             `json:"hash" bson:"hash"`

	AwayTeamName    string `json:"away_team_name" bson:"away_team_name"`
	HomeTeamName    string `json:"home_team_name" bson:"home_team_name"`
	AwayMoneyline   string `json:"away_moneyline" bson:"away_moneyline"`
	HomeMoneyline   string `json:"home_moneyline" bson:"home_moneyline"`
	DrawMoneyline   string `json:"draw_moneyline" bson:"draw_moneyline"`
	AwaySpreadPoint string `json:"away_spread_point" bson:"away_spread_point"`
	AwaySpreadPrice string `json:"away_spread_price" bson:"away_spread_price"`
	HomeSpreadPoint string `json:"home_spread_point" bson:"home_spread_point"`
	HomeSpreadPrice string `json:"home_spread_price" bson:"home_spread_price"`
	UnderPoint      string `json:"under_point" bson:"under_point"`
	UnderPrice      string `json:"under_price" bson:"under_price"`
	OverPoint       string `json:"over_point" bson:"over_point"`
	OverPrice       string `json:"over_price" bson:"over_price"`

	AwayLogoUrl string `json:"away_logo_url" bson:"away_logo_url"`
	HomeLogoUrl string `json:"home_logo_url" bson:"home_logo_url"`
	AwayRecord  string `json:"away_record" bson:"away_record"`
	HomeRecord  string `json:"home_record" bson:"home_record"`

	Props []Prop `json:"props" bson:"props"`
}

type Prop struct {
	Name  string `json:"name" bson:"name"`
	Stats []Stat `json:"stats" bson:"stats"`
}

type Stat struct {
	Name        string `json:"name" bson:"name"`
	OverPoint   string `json:"over_point" bson:"over_point"`
	OverPrice   string `json:"over_price" bson:"over_price"`
	UnderPoint  string `json:"under_point" bson:"under_point"`
	UnderPrice  string `json:"under_price" bson:"under_price"`
	IsOverUnder bool   `json:"is_over_under" bson:"is_over_under"`
	Price       string `json:"price" bson:"price"`
	Point       string `json:"point" bson:"point"`
}

type Pick struct {
	Id               primitive.ObjectID `json:"id" bson:"_id"`
	PlayerName       string             `json:"player_name" bson:"player_name"`
	Point            string             `json:"point" bson:"point"`
	PlayerPictureUrl string             `json:"player_picture_url" bson:"player_picture_url"`
	TeamPosition     string             `json:"team_position" bson:"team_position"`
	StartDate        time.Time          `json:"start_date" bson:"start_date"`
	Sport            string             `json:"sport" bson:"sport"`
	LeagueName       string             `json:"league_name" bson:"league_name"`
	LeagueId         string             `json:"league_id" bson:"league_id"`
	Opponent         string             `json:"opponent" bson:"opponent"`
	Market           string             `json:"market" bson:"market"`
	Hash             string             `json:"hash" bson:"hash"`
}

type League struct {
	Id                primitive.ObjectID `json:"id" bson:"_id"`
	LeagueName        string             `json:"league_name" bson:"league_name"`
	LeagueId          string             `json:"league_id" bson:"league_id"`
	Sport             string             `json:"sport" bson:"sport"`
	LeagueDescription string             `json:"league_description" bson:"league_description"`
	PhotoUrl          string             `json:"url" bson:"url"`
	Active            bool               `json:"active" bson:"active"`
	Popular           bool               `json:"popular" bson:"popular"`
}

type Sport struct {
	SportId   string `json:"sport_id" bson:"sport_id"`
	SportName string `json:"name" bson:"name"`
	PhotoUrl  string `json:"url" bson:"url"`
}
