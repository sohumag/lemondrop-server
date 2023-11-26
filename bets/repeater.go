package bets

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/rlvgl/bookie-server/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	for _, bet := range allBets {
		if bet.IsParlay {
			s.CheckParlay(&bet)
		} else { // singles
			s.CheckBet(&bet, false)
		}

		filter := bson.D{{Key: "_id", Value: bet.BetId}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "bet_status", Value: bet.BetStatus}, {Key: "parlay_finished", Value: bet.ParlayFinished}}}}
		_, err := coll.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (s *BetServer) CheckParlay(parlay *Bet) {
	// check if parlay is finisehd. not necessary but extrea check is nice
	s.UpdateParlayStatus(parlay)
	if parlay.ParlayFinished {
		return
	}

	// parlay is not completed:
	// for each bet in parlay
	// check status of game. if doesnt exist: next game
	// update status of parlay at end
	// update database w parlay: done at top level caller**

	for _, bet := range parlay.Bets {
		s.CheckBet(&bet, true)
	}

	s.UpdateParlayStatus(parlay)

	return
}

func (s *BetServer) UpdateParlayStatus(parlay *Bet) {
	inProgress := false
	lost := false
	for _, bet := range parlay.Bets {
		if bet.BetStatus == "Pending" {
			inProgress = true
			break
		} else if bet.BetStatus == "Lost" {
			lost = true
		}
	}
	if inProgress {
		return
	} else if lost {
		parlay.ParlayFinished = true
		parlay.BetStatus = "Lost"
		s.ChangeUserFundsWin(parlay.UserId, parlay.BetAmount)
	} else {
		parlay.ParlayFinished = true
		parlay.BetStatus = "Won"
		s.ChangeUserFundsLose(parlay.UserId, parlay.BetAmount)
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

func (s *BetServer) CheckBet(bet *Bet, partOfParlay bool) {
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

	// parse game and update bet type
	awayScore, err := strconv.Atoi(gameScore.AwayFinalScore)
	homeScore, err := strconv.Atoi(gameScore.HomeFinalScore)
	switch strings.ToLower(bet.BetType) {
	case "moneyline":
		// fmt.Println("validating ml")
		fmt.Println(bet.AwayTeam, bet.HomeTeam, bet.BetOnTeam)
		if bet.BetOnTeam == bet.AwayTeam {
			if awayScore > homeScore {
				fmt.Println("away team wins")
				s.MarkBetAsValid(bet, partOfParlay)
			} else if awayScore < homeScore {
				fmt.Println("home team wins")
				s.MarkBetAsInvalid(bet, partOfParlay)
			}
		} else if bet.BetOnTeam == bet.HomeTeam {
			if homeScore > awayScore {
				fmt.Println("home team wins")
				s.MarkBetAsValid(bet, partOfParlay)
			} else if homeScore < awayScore {
				fmt.Println("away team wins")
				s.MarkBetAsInvalid(bet, partOfParlay)
			}
		}
	case "spread":
		// fmt.Println(bet.AwayTeam, bet.HomeTeam, bet.BetOnTeam)
		numPoint := 0.0
		if bet.BetPoint[0] == '-' {
			numPoint, err = strconv.ParseFloat(bet.BetPoint[1:], 64)
			numPoint = numPoint * -1
		} else if bet.BetPoint[0] == '+' {
			numPoint, err = strconv.ParseFloat(bet.BetPoint[1:], 64)
		}
		homeScore, _ := strconv.ParseFloat(gameScore.HomeFinalScore, 64)
		awayScore, _ := strconv.ParseFloat(gameScore.AwayFinalScore, 64)

		if strings.Contains(bet.BetOnTeam, bet.AwayTeam) {
			if awayScore+numPoint > awayScore {
				s.MarkBetAsValid(bet, partOfParlay)
			} else if awayScore+numPoint < awayScore {
				s.MarkBetAsInvalid(bet, partOfParlay)
			}

		} else if strings.Contains(bet.BetOnTeam, bet.HomeTeam) {
			if homeScore+numPoint > awayScore {
				s.MarkBetAsValid(bet, partOfParlay)
			} else if homeScore+numPoint < awayScore {
				s.MarkBetAsInvalid(bet, partOfParlay)
			}
		}

	case "over", "under":
		point, _ := strconv.ParseFloat(bet.BetPoint, 64)
		homeScore, _ := strconv.ParseFloat(gameScore.HomeFinalScore, 64)
		awayScore, _ := strconv.ParseFloat(gameScore.AwayFinalScore, 64)
		teamTotal := homeScore + awayScore
		if strings.ToLower(bet.BetType) == "over" {
			if teamTotal > point {
				s.MarkBetAsValid(bet, partOfParlay)
			} else if teamTotal < point {
				s.MarkBetAsInvalid(bet, partOfParlay)
			} else if teamTotal == point {
				s.MarkBetAsPush(bet, partOfParlay)
			}
		} else if strings.ToLower(bet.BetType) == "under" {
			if teamTotal < point {
				s.MarkBetAsValid(bet, partOfParlay)
			} else if teamTotal > point {
				s.MarkBetAsInvalid(bet, partOfParlay)
			} else if teamTotal == point {
				s.MarkBetAsPush(bet, partOfParlay)
			}
		}

	}
}

func (s *BetServer) MarkBetAsValid(bet *Bet, partOfParlay bool) {
	bet.BetStatus = "Won"
	fmt.Println("bet is valid. Marking as Won.")
	if !partOfParlay {
		s.ChangeUserFundsWin(bet.UserId, bet.BetAmount)
	}
}

func (s *BetServer) MarkBetAsInvalid(bet *Bet, partOfParlay bool) {
	bet.BetStatus = "Lost"
	fmt.Println("bet is invalid. Marking as Lost.")

}

func (s *BetServer) MarkBetAsPush(bet *Bet, partOfParlay bool) {
	bet.BetStatus = "Pushed"
	fmt.Println("bet is pushed")
}

// takes user id
func (s *BetServer) ChangeUserFundsWin(userId string, amount string) {
	// remove from pending. Add to profit. Add to balance.
	coll := s.client.Database("users-db").Collection("users")
	id, _ := primitive.ObjectIDFromHex(userId)
	filter := bson.M{"_id": id}
	user := users.User{}
	coll.FindOne(context.TODO(), filter).Decode(&user)
	amt, _ := strconv.ParseFloat(amount, 64)
	newPendingAmt := user.CurrentPending - amt
	newBalanceAmt := user.CurrentBalance + amt
	newProfitAmt := user.TotalProfit + amt
	update := bson.M{"$set": bson.M{"current_pending": newPendingAmt, "current_balance": newBalanceAmt, "total_profit": newProfitAmt}}
	coll.UpdateOne(context.TODO(), filter, update)
}

func (s *BetServer) ChangeUserFundsLose(userId string, amount string) {
	// remove from pending
	coll := s.client.Database("users-db").Collection("users")
	id, _ := primitive.ObjectIDFromHex(userId)
	filter := bson.M{"_id": id}
	user := users.User{}
	coll.FindOne(context.TODO(), filter).Decode(&user)
	amtToSubtract, _ := strconv.ParseFloat(amount, 64)
	newAmt := user.CurrentPending - amtToSubtract
	update := bson.M{"$set": bson.M{"current_pending": newAmt}}
	coll.UpdateOne(context.TODO(), filter, update)
}
