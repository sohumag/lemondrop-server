package bets

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rlvgl/bookie-server/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
	on bet placed: check if balance will allow -> move to pending.
	on bet finished: move from pending to either profit or delete.
*/

func (s *BetServer) HandleBetRequest(c *fiber.Ctx) error {
	bet := &Bet{}
	if err := c.BodyParser(&bet); err != nil {
		fmt.Println(err)
	}

	// bet.BetStatus = "Pending"
	// bet.BetId = primitive.NewObjectID()

	if bet.IsParlay {
		for i := 0; i < len(bet.Bets); i++ {
			bet.Bets[i].BetStatus = "Pending"
		}
	}

	// check if user has balance available
	// if is parlay. check once.
	// if singles. add all amounts in bets body and check
	betAmt := 0.0
	if bet.IsParlay {
		betAmt, _ = strconv.ParseFloat(bet.BetAmount, 64)
	} else {
		for _, subBet := range bet.Bets {
			// fmt.Println(subBet.BetOnTeam)
			subBetAmt, _ := strconv.ParseFloat(subBet.BetAmount, 64)
			betAmt += subBetAmt
		}
	}

	userColl := s.client.Database("users-db").Collection("users")
	user := users.User{}
	id, _ := primitive.ObjectIDFromHex(bet.UserId)
	filter := bson.M{"_id": id}
	userColl.FindOne(context.TODO(), filter).Decode(&user)

	// will fail all if not enough in balance
	if betAmt > user.CurrentBalance {
		return fmt.Errorf("broke boy")
	}

	update := bson.M{"$set": bson.M{"current_balance": user.CurrentBalance - betAmt, "current_pending": user.CurrentPending + betAmt}}
	_, err := userColl.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println(err)
	}

	if bet.IsParlay {
		s.AddBetToDB(*bet, c)
	} else {
		for _, sb := range bet.Bets {
			s.AddBetToDB(sb, c)
		}
	}

	// s.AddBetToDB(*bet, c)

	return nil

}

func (s *BetServer) AddBetToDB(bet Bet, c *fiber.Ctx) error {
	// if parlay: add directly
	// else add for each bet within
	coll := s.client.Database("bets-db").Collection("bets")

	bet.BetStatus = "Pending"
	bet.BetId = primitive.NewObjectID()
	result, err := coll.InsertOne(context.TODO(), bet)
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
