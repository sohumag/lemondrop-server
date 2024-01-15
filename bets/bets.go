package bets

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rlvgl/bookie-server/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// HandleBetRequest handles the incoming bet request.
func (s *BetServer) HandleBetRequest(c *fiber.Ctx) error {
	bets := make([]Bet, 0)
	if err := c.BodyParser(&bets); err != nil {
		fmt.Println(err)
		return c.Status(http.StatusBadRequest).SendString("Invalid request body")
	}

	for _, bet := range bets {
		bet.Status = "Pending"

		bet.PlacedAt = time.Now()
		// get game by hash and get time start. check if already passed. will work for parlays partially
		// coll := s.client.Database("games-db").Collection("scraped-games")
		// filter := bson.D{{Key: "hash", Value: bet.Selections[0].GameHash}}

		// game := games.Game{}
		// coll.FindOne(context.Background(), filter).Decode(&game)

		// if time.Now().After(game.StartDate) {
		// 	return fmt.Errorf("") // leave empty bc most likely hack attempt
		// }

		if bet.IsParlay {
			bet.ParlayFinished = false
			for _, b := range bet.Selections {
				b.BetStatus = "Pending"
			}
		}

		// fmt.Printf("Bet added to transaction: %v\n", bet)
	}

	session, err := s.client.StartSession()
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to start session")
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), func(sessionContext mongo.SessionContext) (interface{}, error) {
		for _, bet := range bets {
			if err := s.HandleBet(&bet, sessionContext); err != nil {
				return nil, err
			}
		}
		return nil, nil
	})

	if err != nil {
		fmt.Println("transaction failed")
		return c.Status(http.StatusInternalServerError).SendString("Transaction failed")
	}

	return c.SendString("Bets placed successfully")
}
func (s *BetServer) HandleBet(bet *Bet, sessionContext mongo.SessionContext) error {
	userColl := s.client.Database("users-db").Collection("users")
	user := users.User{}
	id, err := primitive.ObjectIDFromHex(bet.UserID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": id}

	betAmt, err := stringToFloat64(bet.Amount)
	if err != nil {
		return err
	}

	err = userColl.FindOne(sessionContext, filter).Decode(&user)
	if err != nil {
		return err
	}

	// Calculate the total available balance (free play + current availability + current balance)
	totalAvailable := user.CurrentFreePlay + user.CurrentAvailability + user.CurrentBalance

	// Use free play first if available
	if betAmt <= user.CurrentFreePlay {
		// If the bet amount is less than or equal to free play, use free play only
		update := bson.M{"$set": bson.M{
			"current_free_play": user.CurrentFreePlay - betAmt,
			"current_pending":   user.CurrentPending + betAmt,
		}}
		_, err = userColl.UpdateOne(sessionContext, filter, update)
		if err != nil {
			return err
		}
	} else if betAmt <= (user.CurrentFreePlay + user.CurrentAvailability) {
		// If the bet amount is less than or equal to the total of free play and current availability, use free play and availability
		update := bson.M{"$set": bson.M{
			"current_free_play":    0,
			"current_availability": user.CurrentAvailability - (betAmt - user.CurrentFreePlay),
			"current_pending":      user.CurrentPending + betAmt,
		}}
		_, err = userColl.UpdateOne(sessionContext, filter, update)
		if err != nil {
			return err
		}
	} else if betAmt <= totalAvailable {
		// If the bet amount is less than or equal to the total available balance, use free play, availability, and deduct the remaining from current balance
		update := bson.M{"$set": bson.M{
			"current_free_play":    0,
			"current_availability": 0,
			"current_balance":      user.CurrentBalance - (betAmt - user.CurrentFreePlay - user.CurrentAvailability),
			"current_pending":      user.CurrentPending + betAmt,
		}}
		_, err = userColl.UpdateOne(sessionContext, filter, update)
		if err != nil {
			return err
		}
	} else {
		// Insufficient funds
		return fmt.Errorf("insufficient availability")
	}

	return s.AddBetToDB(bet, sessionContext)
}

func (s *BetServer) AddBetToDB(bet *Bet, sessionContext mongo.SessionContext) error {
	coll := s.client.Database("bets-db").Collection("bets")

	bet.Status = "Pending"
	bet.ID = primitive.NewObjectID()
	bet.PlacedAt = time.Now()

	_, err := coll.InsertOne(sessionContext, bet)
	if err != nil {
		return err
	}

	fmt.Printf("Bet placed with ID: %v\n", bet.ID)
	return nil
}
