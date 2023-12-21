package bets

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Bet represents the structure of a bet in the sportsbook.

type Bet struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"userId"`
	Timestamp time.Time          `bson:"timestamp"`
	Status    string             `bson:"status"`
	Amount    float64            `bson:"amount"`
	Payout    float64            `bson:"payout"`
	BetType   string             `bson:"betType"`
	Details   BetDetails         `bson:"details"`
	IsParlay  bool               `bson:"isParlay"`
}

// BetDetails contains details specific to each bet type.
type BetDetails struct {
	BetType      string `bson:"betType"` // moneyline, spread, total, prop
	PropType     string `bson:"propType,omitempty"`
	BetOn        string `bson:"betOn,omitempty"`
	GameHash     string `bson:"gameHash,omitempty"`
	AwayTeamName string `bson:"awayTeamName,omitempty"`
	HomeTeamName string `bson:"homeTeamName,omitempty"`

	Team string  `bson:"team,omitempty"`
	Odds float64 `bson:"odds"`
	Legs []Bet   `bson:"legs,omitempty"`
}
