package bets

import (
	"strconv"

	"github.com/rlvgl/bookie-server/users"
)

// hasEnoughBalance checks if the user has enough balance for the bet.
func hasEnoughBalance(amt float64, user users.User) bool {
	return amt <= user.CurrentBalance
}

func stringToFloat64(input string) (float64, error) {
	result, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func overlapCount(s1, s2 string) int {
	set := make(map[rune]struct{})

	for _, char := range s1 {
		set[char] = struct{}{}
	}

	count := 0
	for _, char := range s2 {
		if _, ok := set[char]; ok {
			count++
		}
	}

	return count
}

// import (
// 	"context"

// 	"github.com/gofiber/fiber/v2"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// type GetBetParams struct {
// 	UserId string `json:"user_id"`
// }

// func (s *BetServer) GetAllBetsByUserId(c *fiber.Ctx, userId string) error {
// 	coll := s.client.Database("bets-db").Collection("bets")

// 	filter := bson.D{{Key: "user_id", Value: userId}}
// 	opts := options.Find().SetSort(bson.D{{Key: "bet_placed_time", Value: -1}})
// 	cursor, err := coll.Find(context.TODO(), filter, opts)

// 	if err != nil {
// 		return err
// 	}

// 	allBets := []Bet{}
// 	for cursor.Next(context.TODO()) {
// 		b := Bet{}
// 		cursor.Decode(&b)
// 		allBets = append(allBets, b)
// 	}

// 	c.JSON(allBets)
// 	return nil
// }

// // admin route: can protect later
// func (s *BetServer) GetAllBets(c *fiber.Ctx) error {
// 	coll := s.client.Database("bets-db").Collection("bets")

// 	filter := bson.D{{}}
// 	cursor, err := coll.Find(context.TODO(), filter)
// 	if err != nil {
// 		return err
// 	}

// 	allBets := []Bet{}
// 	for cursor.Next(context.TODO()) {
// 		b := Bet{}
// 		cursor.Decode(&b)
// 		allBets = append(allBets, b)
// 	}

// 	c.JSON(allBets)

// 	return nil
// }
