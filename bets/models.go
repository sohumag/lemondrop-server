package bets

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Bet represents the structure of a bet in the sportsbook.
type Bet struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UserID     string             `bson:"user_id" json:"user_id"`
	Amount     string             `bson:"bet_amount" json:"bet_amount"`
	IsParlay   bool               `bson:"is_parlay" json:"is_parlay"`
	Selections []BetSelection     `bson:"selections" json:"selections"`

	// server managed
	Status         string    `bson:"bet_status" json:"bet_status"`
	ParlayFinished bool      `bson:"parlay_finished" json:"parlay_finished"`
	PlacedAt       time.Time `bson:"placed_at" json:"placed_at"`
	// jwt??
}

// BetSelection represents a selection within a bet.
type BetSelection struct {
	BetType    string `bson:"bet_type" json:"bet_type"`
	PropType   string `bson:"prop_type,omitempty" json:"prop_type,omitempty"`
	PropName   string `bson:"prop_name,omitempty" json:"prop_name,omitempty"`
	PlayerName string `bson:"player_name,omitempty" json:"player_name,omitempty"`
	BetOn      string `bson:"bet_on" json:"bet_on"`
	BetPoint   string `bson:"bet_point,omitempty" json:"bet_point,omitempty"`
	Odds       string `bson:"odds,omitempty" json:"odds,omitempty"`

	GameId       string `bson:"game_id,omitempty" json:"game_id"`
	GameHash     string `bson:"game_hash" json:"game_hash"`
	HomeTeamName string `bson:"home_team_name" json:"home_team_name"`
	AwayTeamName string `bson:"away_team_name" json:"away_team_name"`

	// need to add
	LeagueID   string `bson:"league_id" json:"league_id"`
	LeagueName string `bson:"league_name" json:"league_name"`
	SportName  string `bson:"sport_name" json:"sport_name"`

	// server managed
	BetStatus string `bson:"bet_status" json:"bet_status"` // Add BetStatus field
}

type Score struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Completed      bool               `bson:"completed" json:"completed"`
	AwayTeamName   string             `bson:"away_team_name" json:"away_team_name"`
	AwayFinalScore string             `bson:"away_final_score" json:"away_final_score"`
	HomeTeamName   string             `bson:"home_team_name" json:"home_team_name"`
	HomeFinalScore string             `bson:"home_final_score" json:"home_final_score"`
	StartDate      time.Time          `bson:"start_date" json:"start_date"`
	Hash           string             `bson:"hash" json:"hash"`
	LeagueID       string             `bson:"league_id" json:"league_id"`

	AwayFirstQuarter  string         `bson:"away_first_quarter" json:"away_first_quarter"`
	AwaySecondQuarter string         `bson:"away_second_quarter" json:"away_second_quarter"`
	AwayThirdQuarter  string         `bson:"away_third_quarter" json:"away_third_quarter"`
	AwayFourthQuarter string         `bson:"away_fourth_quarter" json:"away_fourth_quarter"`
	AwayPlayerStats   []PlayerStats  `bson:"away_player_stats" json:"away_player_stats"`
	HomeFirstQuarter  string         `bson:"home_first_quarter" json:"home_first_quarter"`
	HomeSecondQuarter string         `bson:"home_second_quarter" json:"home_second_quarter"`
	HomeThirdQuarter  string         `bson:"home_third_quarter" json:"home_third_quarter"`
	HomeFourthQuarter string         `bson:"home_fourth_quarter" json:"home_fourth_quarter"`
	HomePlayerScores  []PlayerScores `bson:"home_player_scores" json:"home_player_scores"`
}

type PlayerStats struct {
	// Define the fields for player stats
}

type PlayerScores struct {
	// Define the fields for player scores
}
