package games

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *GameServer) GetAllSports(c *fiber.Ctx) error {
	sports := []Sport{
		{SportName: "Basketball"},
		{SportName: "Football"},
		{SportName: "Ice Hockey"},
		{SportName: "Soccer"},
		{SportName: "Combat Sports"},
	}

	c.JSON(sports)

	return nil
}

func (s *GameServer) GetAllSportsAndLeagues(c *fiber.Ctx) error {
	coll := s.client.Database("games-db").Collection("leagues")
	cursor, err := coll.Find(context.TODO(), bson.M{"active": true})
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	leagues := []League{}
	league := League{}
	for cursor.Next(context.TODO()) {
		cursor.Decode(&league)
		leagues = append(leagues, league)
	}

	// schema: array of [popular: []league, ]
	// get all leagues
	// get all leagues where popular is true
	// get all unique sports names

	// return: popular: [], sports: [], allLeagues: []

	var popularLeagues []League
	uniqueSports := []Sport{
		{SportName: "Basketball", SportId: "basketball", PhotoUrl: "https://a.espncdn.com/combiner/i?img=/redesign/assets/img/icons/ESPN-icon-basketball.png&h=80&w=80&scale=crop&cquality=40"},
		{SportName: "Football", SportId: "football", PhotoUrl: "https://a3.espncdn.com/combiner/i?img=%2Fredesign%2Fassets%2Fimg%2Ficons%2FESPN%2Dicon%2Dfootball%2Dcollege.png&w=80&h=80&scale=crop&cquality=40&location=origin"},
		{SportName: "Ice Hockey", SportId: "ice_hockey", PhotoUrl: "https://a4.espncdn.com/combiner/i?img=%2Fredesign%2Fassets%2Fimg%2Ficons%2FESPN%2Dicon%2Dhockey.png&w=80&h=80&scale=crop&cquality=40&location=origin"},
		{SportName: "Soccer", SportId: "soccer", PhotoUrl: "https://a.espncdn.com/combiner/i?img=/redesign/assets/img/icons/ESPN-icon-soccer.png&w=64&h=64&scale=crop&cquality=40&location=origin"},
		{SportName: "Combat Sports", SportId: "combat_sports", PhotoUrl: "https://a.espncdn.com/combiner/i?img=/redesign/assets/img/icons/ESPN-icon-mma.png&w=64&h=64&scale=crop&cquality=40&location=origin"},
	}

	// Loop through the original array
	for _, league := range leagues {
		// Check if the league is popular
		if league.Popular {
			popularLeagues = append(popularLeagues, league)
		}
	}

	allSportsLeagues := map[string]interface{}{
		"Popular":     popularLeagues,
		"Sports":      uniqueSports,
		"All Leagues": leagues,
	}

	c.JSON(allSportsLeagues)

	return nil
}

func (s *GameServer) GetAllLeagues(c *fiber.Ctx) error {
	coll := s.client.Database("games-db").Collection("leagues")
	cursor, err := coll.Find(context.TODO(), bson.M{"active": true})
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	leagues := []League{}
	league := League{}
	for cursor.Next(context.TODO()) {
		cursor.Decode(&league)
		leagues = append(leagues, league)
	}
	c.JSON(leagues)
	return nil
}

func (s *GameServer) GetAllLeaguesBySport(c *fiber.Ctx) error {
	sport := c.Params("sport")
	fmt.Println(sport)

	coll := s.client.Database("games-db").Collection("leagues")
	cursor, err := coll.Find(context.TODO(), bson.M{"active": true})

	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	leagues := []League{}
	league := League{}
	for cursor.Next(context.TODO()) {
		// add only if league is part of sport
		cursor.Decode(&league)
		if strings.Replace(strings.ToLower(league.Sport), " ", "_", -1) == sport {
			leagues = append(leagues, league)
		}
	}

	c.JSON(leagues)
	return nil
}
