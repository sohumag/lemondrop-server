package bets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *BetServer) AddBetToDB(c *fiber.Ctx) error {
	bet := &Bet{}
	if err := c.BodyParser(&bet); err != nil {
		fmt.Println(err)
	}
	bet.BetStatus = "Pending"
	// fmt.Println(bet)

	coll := s.client.Database("bets-db").Collection("bets")
	result, err := coll.InsertOne(context.TODO(), &bet)
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		fmt.Println(err)
	}

	fmt.Printf("Bet placed with id: %v\n", result.InsertedID)
	return nil
}

type GetBetParams struct {
	UserId string `json:"user_id"`
}

func (s *BetServer) GetAllBetsByUserId(c *fiber.Ctx, userId string) error {
	coll := s.client.Database("bets-db").Collection("bets")

	filter := bson.D{{Key: "user_id", Value: userId}}
	opts := options.Find().SetSort(bson.D{{Key: "bet_placed_time", Value: -1}})
	cursor, err := coll.Find(context.TODO(), filter, opts)

	if err != nil {
		return err
	}

	allBets := []Bet{}
	for cursor.Next(context.TODO()) {
		b := Bet{}
		cursor.Decode(&b)
		allBets = append(allBets, b)
	}

	c.JSON(allBets)

	return nil
}

// admin route: can protect later
func (s *BetServer) GetAllBets(c *fiber.Ctx) error {
	coll := s.client.Database("bets-db").Collection("bets")

	filter := bson.D{{}}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return err
	}

	allBets := []Bet{}
	for cursor.Next(context.TODO()) {
		b := Bet{}
		cursor.Decode(&b)
		allBets = append(allBets, b)
	}

	c.JSON(allBets)

	return nil
}

func ValidateBet(bet *Bet) error {
	if bet.UserId == "" {
		return fmt.Errorf("Invalid HTTP Request")
	}
	if bet.UserEmail == "" {
		return fmt.Errorf("Invalid HTTP Request")
	}

	if bet.GameId == "" {
		return fmt.Errorf("Invalid HTTP Request")
	}
	if bet.HomeTeam == "" {
		return fmt.Errorf("Invalid HTTP Request")
	}
	if bet.AwayTeam == "" {
		return fmt.Errorf("Invalid HTTP Request")
	}
	//if bet.GameStartTime.IsZero() {
	//	return fmt.Errorf("Invalid HTTP Request")
	//}

	if bet.BetType == "" {
		return fmt.Errorf("Invalid HTTP Request")
	}
	if bet.BetOnTeam == "" {
		return fmt.Errorf("Invalid HTTP Request")
	}
	// if bet.BetCategory == "" {
	// 	return fmt.Errorf("Invalid HTTP Request")
	// }

	// if bet.BetPoint == 0 {
	// 	return fmt.Errorf("Invalid HTTP Request")
	// }
	if bet.BetPrice == "" {
		return fmt.Errorf("Invalid HTTP Request")
	}
	if bet.BetAmount == "" {
		return fmt.Errorf("Invalid HTTP Request")
	}

	return nil
}
