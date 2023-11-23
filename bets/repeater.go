package bets

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
	BET STATUSES:
		- (Pending Verification), Pending, Success, Failure
*/

func (s *BetServer) RunBetCheckingRepeater() {
	// add cron job
	// get all bets where status is pending and game start time has passed
	//? get hash of score doc where hash = dk_hash
	// check score from doc and calculate whether bet hit or not
	// update bet status to success or failure
	// update user account balance/pending/earned

	coll := s.client.Database("bets-db").Collection("bets")
	// get all where: commence time is after. state is pending.
	filter := bson.D{{Key: "bet_status", Value: "Pending"}}
	// filter := bson.D{}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	allBets := []Bet{}
	placeholder := &Bet{}
	for cursor.Next(context.TODO()) {
		cursor.Decode(placeholder)
		allBets = append(allBets, *placeholder)
	}

	// fmt.Println(allBets)
	for _, bet := range allBets {
		s.CheckBet(&bet)
		filter := bson.D{{Key: "_id", Value: bet.BetId}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "bet_status", Value: bet.BetStatus}}}}
		// fmt.Println(bet.BetStatus, bet.GameHash)
		_, err := coll.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			fmt.Println(err)
			return
		}
		// fmt.Println(res.)
	}
}

type GameScore struct {
	Id             string    `bson:"_id" json:"_id"`
	Completed      bool      `bson:"completed" json:"completed"`
	AwayTeamName   string    `json:"away_team_name" bson:"away_team_name"`
	HomeTeamName   string    `json:"home_team_name" bson:"home_team_name"`
	AwayFinalScore string    `json:"away_final_score" bson:"away_final_score"`
	HomeFinalScore string    `json:"home_final_score" bson:"home_final_score"`
	StartDate      time.Time `json:"start_date" bson:"start_date"`
	EspnHash       string    `json:"espn_hash" bson:"espn_hash"`
	DkHash         string    `json:"dk_hash" bson:"dk_hash"`
	League         string    `json:"league" bson:"league"`
}

func (s *BetServer) CheckBet(bet *Bet) {
	// check scores db for where dk hash matches hash
	coll := s.client.Database("games-db").Collection("scraped-scores")

	gameScore := GameScore{}
	filter := bson.D{{Key: "dk_hash", Value: bet.GameHash}}
	err := coll.FindOne(context.TODO(), filter).Decode(&gameScore)
	if err != nil {
		// when game hasn't happened yet
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}
	// fmt.Println("checking bet:", bet.BetOnTeam, bet.GameHash)
	// fmt.Println(gameScore)

	// fmt.Printf("game score: %v\n", gameScore)
	awayScore, err := strconv.Atoi(gameScore.AwayFinalScore)
	homeScore, err := strconv.Atoi(gameScore.HomeFinalScore)
	// fmt.Printf("%v(%v) @ %v(%v)\n", bet.AwayTeam, awayScore, bet.HomeTeam, homeScore)

	// switch on bet type
	// fmt.Println(bet.BetOnTeam, bet.BetType, bet.AwayTeam)
	// fmt.Println("team names: ", bet.BetOnTeam, gameScore.AwayTeamName)

	// fmt.Println(bet.BetType)

	switch strings.ToLower(bet.BetType) {
	case "moneyline":
		// fmt.Println("validating ml")
		fmt.Println(bet.AwayTeam, bet.HomeTeam, bet.BetOnTeam)
		if bet.BetOnTeam == bet.AwayTeam {
			if awayScore > homeScore {
				fmt.Println("away team wins")
				s.MarkBetAsValid(bet)
			} else if awayScore < homeScore {
				fmt.Println("home team wins")
				s.MarkBetAsInvalid(bet)
			}
		} else if bet.BetOnTeam == bet.HomeTeam {
			if homeScore > awayScore {
				fmt.Println("home team wins")
				s.MarkBetAsValid(bet)
			} else if homeScore < awayScore {
				fmt.Println("away team wins")
				s.MarkBetAsInvalid(bet)
			}
		}
	case "spread":
		if strings.Contains(bet.BetOnTeam, bet.AwayTeam) {
			// parse away spread

		} else if strings.Contains(bet.BetOnTeam, bet.HomeTeam) {
			numPoint := 0.0
			if bet.BetPoint[0] == '-' {
				numPoint, err = strconv.ParseFloat(bet.BetPoint[1:], 64)
				numPoint = numPoint * -1
			} else if bet.BetPoint[0] == '+' {
				numPoint, err = strconv.ParseFloat(bet.BetPoint[1:], 64)
			}
			homeScore, _ := strconv.ParseFloat(gameScore.HomeFinalScore, 64)
			awayScore, _ := strconv.ParseFloat(gameScore.AwayFinalScore, 64)

			if homeScore+numPoint > awayScore {
				s.MarkBetAsValid(bet)
			} else if homeScore+numPoint < awayScore {
				s.MarkBetAsInvalid(bet)
			}
		}

	}

}

func (s *BetServer) MarkBetAsValid(bet *Bet) {
	bet.BetStatus = "Won"
	fmt.Println("bet is valid")
}

func (s *BetServer) MarkBetAsInvalid(bet *Bet) {
	bet.BetStatus = "Lost"
	fmt.Println("bet is invalid")
}

func (s *BetServer) MarkBetAsPush(bet *Bet) {
	bet.BetStatus = "Pushed"
	fmt.Println("bet is pushed")
}
