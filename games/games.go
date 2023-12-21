package games

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *GameServer) GetGamesBySport(c *fiber.Ctx) error {
	coll := s.client.Database("games-db").Collection("scraped-games")
	currentDate := time.Now()
	maxDate := currentDate.Add(time.Hour * 24 * 7)

	parsedSport := strings.Replace(c.Params("sport"), "%20", " ", -1)
	// parsedSport = strings.Replace(c.Params("sport"), "-", " ", -1)
	parsedSport = strings.Title(parsedSport)

	cursor, err := coll.Find(context.TODO(), bson.M{"start_date": bson.M{"$gt": currentDate, "$lt": maxDate}, "sport": parsedSport})
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	games, err := s.DecodeCursorIntoGames(cursor)
	if err != nil {
		return err
	}

	c.JSON(*games)

	return nil
}

func (s *GameServer) GetGamesByLeagueId(c *fiber.Ctx) error {
	fmt.Println("request handling..")
	coll := s.client.Database("games-db").Collection("scraped-games")
	currentDate := time.Now()
	// maxDate := currentDate.Add(time.Hour * 24 * 7)
	parsedLeague := strings.Replace(c.Params("league"), "%20", " ", -1)
	parsedLeague = strings.ToLower(parsedLeague)

	// fmt.Println(parsedLeague)

	// cursor, err := coll.Find(context.TODO(), bson.M{"league_id": parsedLeague})
	// cursor, err := coll.Find(context.TODO(), bson.M{"start_date": bson.M{"$gt": currentDate, "$lt": maxDate}, "league_id": parsedLeague})
	cursor, err := coll.Find(context.TODO(), bson.M{"start_date": bson.M{"$gt": currentDate}, "league_id": parsedLeague})
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	games, err := s.DecodeCursorIntoGames(cursor)
	if err != nil {
		return err
	}

	c.JSON(*games)

	return nil
}

func (s *GameServer) GetGameById(c *fiber.Ctx) error {
	coll := s.client.Database("games-db").Collection("scraped-games")
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	game := Game{}
	coll.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&game)

	c.JSON(game)

	return nil
}
