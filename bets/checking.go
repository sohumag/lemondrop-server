package bets

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// type GameScore struct {
// 	ID             primitive.ObjectID `bson:"_id" json:"_id"`
// 	Completed      bool               `bson:"completed" json:"completed"`
// 	AwayTeamName   string             `json:"away_team_name" bson:"away_team_name"`
// 	HomeTeamName   string             `json:"home_team_name" bson:"home_team_name"`
// 	AwayFinalScore string             `json:"away_final_score" bson:"away_final_score"`
// 	HomeFinalScore string             `json:"home_final_score" bson:"home_final_score"`
// 	StartDate      time.Time          `json:"start_date" bson:"start_date"`
// 	EspnHash       string             `json:"espn_hash" bson:"espn_hash"`
// 	DkHash         string             `json:"dk_hash" bson:"dk_hash"`
// 	League         string             `json:"league" bson:"league"`
// }

// // RunBetCheckingRepeater fetches pending bets, checks their status, and updates the database.
// func (s *BetServer) RunBetCheckingRepeater() {
// 	coll := s.client.Database("bets-db").Collection("bets")
// 	filter := bson.D{{Key: "bet_status", Value: "Pending"}}
// 	cursor, err := coll.Find(context.TODO(), filter)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var allBets []Bet
// 	for cursor.Next(context.TODO()) {
// 		var bet Bet
// 		if err := cursor.Decode(&bet); err != nil {
// 			log.Fatal(err)
// 		}
// 		allBets = append(allBets, bet)
// 	}

// 	for _, bet := range allBets {
// 		if bet.IsParlay {
// 			s.CheckParlay(&bet)
// 		} else {
// 			s.CheckBet(&bet, false)
// 		}

// 		filter := bson.D{{Key: "_id", Value: bet.BetId}}
// 		update := bson.D{{Key: "$set", Value: bson.D{{Key: "bet_status", Value: bet.BetStatus}, {Key: "parlay_finished", Value: bet.ParlayFinished}}}}
// 		_, err := coll.UpdateOne(context.TODO(), filter, update)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 	}
// }

// // CheckParlay checks the status of a parlay bet.
// func (s *BetServer) CheckParlay(parlay *Bet) {
// 	s.UpdateParlayStatus(parlay)
// 	if parlay.ParlayFinished {
// 		return
// 	}

// 	for _, bet := range parlay.Bets {
// 		s.CheckBet(&bet, true)
// 	}

// 	s.UpdateParlayStatus(parlay)
// }

// // UpdateParlayStatus updates the status of a parlay bet.
// func (s *BetServer) UpdateParlayStatus(parlay *Bet) {
// 	inProgress := false
// 	lost := false

// 	for _, bet := range parlay.Bets {
// 		if bet.BetStatus == "Pending" {
// 			inProgress = true
// 			break
// 		} else if bet.BetStatus == "Lost" {
// 			lost = true
// 		}
// 	}

// 	if inProgress {
// 		return
// 	} else if lost {
// 		s.setParlayResult(parlay, "Lost")
// 		s.ChangeUserFundsWin(parlay.UserId, parlay.BetAmount)
// 	} else {
// 		s.setParlayResult(parlay, "Won")
// 		s.ChangeUserFundsLose(parlay.UserId, parlay.BetAmount)
// 	}
// }

// // setParlayResult sets the result for a parlay bet.
// func (s *BetServer) setParlayResult(parlay *Bet, result string) {
// 	parlay.ParlayFinished = true
// 	parlay.BetStatus = result
// }

// // CheckBet checks the status of a single bet.
// func (s *BetServer) CheckBet(bet *Bet, partOfParlay bool) {
// 	coll := s.client.Database("games-db").Collection("scraped-scores")
// 	gameScore := GameScore{}
// 	filter := bson.D{{Key: "dk_hash", Value: bet.GameHash}}

// 	if err := coll.FindOne(context.TODO(), filter).Decode(&gameScore); err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			// Game hasn't happened yet
// 			return
// 		}
// 		log.Fatal(err)
// 	}

// 	// Process the bet based on its type
// 	switch strings.ToLower(bet.BetType) {
// 	case "moneyline":
// 		s.processMoneyLineBet(bet, gameScore)
// 	case "spread":
// 		s.processSpreadBet(bet, gameScore)
// 	case "over", "under":
// 		s.processOverUnderBet(bet, gameScore)
// 	}
// }

// // processMoneyLineBet processes a moneyline bet.
// func (s *BetServer) processMoneyLineBet(bet *Bet, gameScore GameScore) {
// 	awayScore, _ := strconv.Atoi(gameScore.AwayFinalScore)
// 	homeScore, _ := strconv.Atoi(gameScore.HomeFinalScore)

// 	if bet.BetOnTeam == bet.AwayTeam && awayScore > homeScore {
// 		s.MarkBetAsValid(bet, false)
// 	} else if bet.BetOnTeam == bet.HomeTeam && homeScore > awayScore {
// 		s.MarkBetAsValid(bet, false)
// 	} else {
// 		s.MarkBetAsInvalid(bet, false)
// 	}
// }

// // processSpreadBet processes a spread bet.
// func (s *BetServer) processSpreadBet(bet *Bet, gameScore GameScore) {
// 	numPoint, _ := strconv.ParseFloat(bet.BetPoint[1:], 64)
// 	homeScore, _ := strconv.ParseFloat(gameScore.HomeFinalScore, 64)
// 	awayScore, _ := strconv.ParseFloat(gameScore.AwayFinalScore, 64)

// 	if strings.Contains(bet.BetOnTeam, bet.AwayTeam) && awayScore+numPoint > awayScore {
// 		s.MarkBetAsValid(bet, false)
// 	} else if strings.Contains(bet.BetOnTeam, bet.HomeTeam) && homeScore+numPoint > awayScore {
// 		s.MarkBetAsValid(bet, false)
// 	} else {
// 		s.MarkBetAsInvalid(bet, false)
// 	}
// }

// // processOverUnderBet processes an over/under bet.
// func (s *BetServer) processOverUnderBet(bet *Bet, gameScore GameScore) {
// 	point, _ := strconv.ParseFloat(bet.BetPoint, 64)
// 	homeScore, _ := strconv.ParseFloat(gameScore.HomeFinalScore, 64)
// 	awayScore, _ := strconv.ParseFloat(gameScore.AwayFinalScore, 64)
// 	teamTotal := homeScore + awayScore

// 	if strings.ToLower(bet.BetType) == "over" && teamTotal > point {
// 		s.MarkBetAsValid(bet, false)
// 	} else if strings.ToLower(bet.BetType) == "under" && teamTotal < point {
// 		s.MarkBetAsValid(bet, false)
// 	} else {
// 		s.MarkBetAsInvalid(bet, false)
// 	}
// }

// // MarkBetAsValid marks a bet as valid and updates user funds.
// func (s *BetServer) MarkBetAsValid(bet *Bet, partOfParlay bool) {
// 	bet.BetStatus = "Won"
// 	fmt.Println("Bet is valid. Marking as Won.")
// 	if !partOfParlay {
// 		s.ChangeUserFundsWin(bet.UserId, bet.BetAmount)
// 	}
// }

// // MarkBetAsInvalid marks a bet as invalid.
// func (s *BetServer) MarkBetAsInvalid(bet *Bet, partOfParlay bool) {
// 	bet.BetStatus = "Lost"
// 	fmt.Println("Bet is invalid. Marking as Lost.")
// }

// // ChangeUserFundsWin updates user funds for a winning bet.
// func (s *BetServer) ChangeUserFundsWin(userID string, amount string) {
// 	s.changeUserFunds(userID, amount, true)
// }

// // ChangeUserFundsLose updates user funds for a losing bet.
// func (s *BetServer) ChangeUserFundsLose(userID string, amount string) {
// 	s.changeUserFunds(userID, amount, false)
// }

// // changeUserFunds updates user funds based on the bet result.
// func (s *BetServer) changeUserFunds(userID string, amount string, isWin bool) {
// 	coll := s.client.Database("users-db").Collection("users")
// 	id, _ := primitive.ObjectIDFromHex(userID)

// 	var update bson.M
// 	if isWin {
// 		amt, _ := strconv.ParseFloat(amount, 64)
// 		update = bson.M{
// 			"$set": bson.M{
// 				"current_pending": 0,
// 				"current_balance": "$current_balance + " + amount,
// 				"total_profit":    "$total_profit + " + amount,
// 			},
// 		}
// 	} else {
// 		amtToSubtract, _ := strconv.ParseFloat(amount, 64)
// 		update = bson.M{
// 			"$set": bson.M{
// 				"current_pending": "$current_pending - " + amount,
// 			},
// 		}
// 	}

// 	filter := bson.M{"_id": id}
// 	coll.UpdateOne(context.TODO(), filter, update)
// }
