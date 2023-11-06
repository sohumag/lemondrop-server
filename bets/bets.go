package bets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *BetServer) AddBetToDB(c *fiber.Ctx) error {
	fmt.Println("received add bet http request")
	bet := &Bet{}
	if err := c.BodyParser(&bet); err != nil {
		fmt.Println(err)
	}

	if err := ValidateBet(bet); bet != nil {
		return err
	}

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

	// params := GetBetParams{}
	// if err := c.BodyParser(&params); err != nil {
	// 	return err
	// }

	filter := bson.D{{Key: "userid", Value: userId}}
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

	if bet.BetPoint == 0 {
		return fmt.Errorf("Invalid HTTP Request")
	}
	if bet.BetPrice == 0 {
		return fmt.Errorf("Invalid HTTP Request")
	}
	if bet.BetAmount == 0 {
		return fmt.Errorf("Invalid HTTP Request")
	}

	return nil
}
