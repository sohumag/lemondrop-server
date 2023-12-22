package bets

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *BetServer) BetChecker() error {
	allBets, err := s.GetPendingBets()
	if err != nil {
		return err
	}
	allScores, err := s.GetAllScores()
	if err != nil {
		fmt.Println(err)
	}

	for _, bet := range allBets {
		if bet.IsParlay {
			continue
		} else {
			_, err = s.CheckBetSelection(&bet.Selections[0], allScores)
		}
	}

	return nil
}

// can implement some sort of cache
func (s *BetServer) CheckBetSelection(selection *BetSelection, allScores []Score) (string, error) { // returns new status, error
	homeTeam := selection.HomeTeamName
	awayTeam := selection.HomeTeamName
	// leagueId := selection.LeagueID

	// get info from mongo using game hash

	// homeTeam := "Cleveland Cavaliers"
	// awayTeam := "Utah Jazz"
	// needs to have game commence time
	leagueId := "basketball_nba"

	allScoresAwayTeams := map[string]string{}
	allScoresHomeTeams := map[string]string{}
	for _, score := range allScores {
		if score.LeagueID != leagueId {
			continue
		}
		allScoresAwayTeams[score.AwayTeamName] = score.Hash
		allScoresHomeTeams[score.HomeTeamName] = score.Hash
	}

	closestGameHash, homeTeamDist := findClosestMatch(homeTeam, allScoresHomeTeams)
	closestGameHash, awayTeamDist := findClosestMatch(awayTeam, allScoresAwayTeams)

	fmt.Printf("%v @ %v: %v %v. Hash is %v\n", awayTeam, homeTeam, awayTeamDist, homeTeamDist, closestGameHash)
	return "", nil
}

func findClosestMatch(query string, nameHash map[string]string) (string, float64) {
	var maxSimilarity float64
	var closestMatch string

	for teamName := range nameHash {
		overlapCount := overlapCount(query, teamName)

		// Calculate Jaccard Similarity
		jaccardSimilarity := float64(overlapCount) / float64(len(query)+len(teamName)-overlapCount)

		if closestMatch == "" || jaccardSimilarity > maxSimilarity {
			maxSimilarity = jaccardSimilarity
			closestMatch = teamName
		}
	}

	return closestMatch, maxSimilarity
}

func (s *BetServer) GetPendingBets() ([]Bet, error) {
	// Access the specified database and collection
	coll := s.client.Database("bets-db").Collection("bets")

	// Define a filter to get all documents (if you have any specific filter criteria, you can modify it)
	filter := bson.D{{"bet_status", "Pending"}}

	// Perform the query to get all bets
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Decode the results into a slice of Bet objects
	var allBets []Bet
	for cursor.Next(context.TODO()) {
		var bet Bet
		if err := cursor.Decode(&bet); err != nil {
			log.Println("Error decoding bet:", err)
		}
		allBets = append(allBets, bet)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return allBets, nil
}

func (s *BetServer) GetAllScores() ([]Score, error) {
	coll := s.client.Database("games-db").Collection("scraped-scores")

	// Define a filter to get scores in the last 3 days
	threeDaysAgo := time.Now().AddDate(0, 0, -2)
	filter := bson.D{
		{"start_date", bson.D{{"$gte", threeDaysAgo}}},
	}

	// Define options to sort by start_date in descending order
	options := options.Find().SetSort(bson.D{{"start_date", -1}})

	// Perform the query to get scores
	cursor, err := coll.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Decode the results into a slice of Score objects
	var allScores []Score
	for cursor.Next(context.TODO()) {
		var score Score
		if err := cursor.Decode(&score); err != nil {
			log.Println("Error decoding score:", err)
		}
		allScores = append(allScores, score)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return allScores, nil
}
