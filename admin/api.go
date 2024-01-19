package admin

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rlvgl/bookie-server/bets"
	"github.com/rlvgl/bookie-server/messages"
	"github.com/rlvgl/bookie-server/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewAdminServer() *AdminServer {
	return &AdminServer{
		client: ConnectDB(),
	}
}

type AdminServer struct {
	client *mongo.Client
}

func ConnectDB() *mongo.Client {
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client
}

func (s *AdminServer) StartAdminServerAPI(api fiber.Router) error {
	log.Println("Adding Admin server endpoints to API")
	// Admin
	mapi := api.Group("/admin")

	mapi.Get("/messages/all", func(c *fiber.Ctx) error {
		return s.GetAllMessages(c)
	})

	mapi.Get("/metrics/balances/all", func(c *fiber.Ctx) error {
		return s.GetAllBalanceMetrics(c)
	})

	mapi.Post("/metrics/balances/complete/:id", func(c *fiber.Ctx) error {
		return s.CompleteBalanceMetric(c)
	})

	mapi.Get("/users/all", func(c *fiber.Ctx) error {
		return s.GetAllUsers(c)
	})

	mapi.Get("/bets/all", func(c *fiber.Ctx) error {
		return s.GetAllBets(c)
	})

	mapi.Post("/bets/status/:betId", func(c *fiber.Ctx) error {
		return s.ChangeBetsStatus(c)
	})

	return nil
}

func calculateOdds(bet bets.Bet) float64 {
	currentOdds := 1.0

	for _, selection := range bet.Selections {
		priceFloat := 0
		priceString := fmt.Sprint(selection.Odds)

		if priceString[0] == '+' {
			priceFloat = parseInt(priceString[1:])
		} else {
			priceFloat = parseInt(priceString[1:]) * -1
		}

		decimalOdds := 1.0
		if priceFloat > 0 {
			decimalOdds = 1 + float64(priceFloat)/100
		} else {
			decimalOdds = 1 - 100/float64(priceFloat)
		}

		currentOdds *= decimalOdds
	}

	return currentOdds
}

func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}

func parseFloat(s string) (float64, error) {
	result, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0, err
	}
	return result, nil
}

func (s *AdminServer) ChangeBetsStatus(c *fiber.Ctx) error {
	// get bet id => get bet
	// update bet status
	// update user status

	ValidateUserAdmin(c)

	betId := c.Params("betId")

	// Convert id string to ObjectId
	objectID, err := primitive.ObjectIDFromHex(betId)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ObjectID"})
	}
	collection := s.client.Database("bets-db").Collection("bets")
	filter := bson.D{{"_id", objectID}}

	var bet bets.Bet
	err = collection.FindOne(context.Background(), filter).Decode(&bet)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Document not found"})
	}

	floatAmt, err := parseFloat(bet.Amount)
	if err != nil {
		return err
	}
	toWinAmount := floatAmt * calculateOdds(bet)
	betAmt, err := parseFloat(bet.Amount)

	session, err := s.client.StartSession()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to start session")
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), func(sessionContext mongo.SessionContext) (interface{}, error) {
		user := users.User{}
		userID, err := primitive.ObjectIDFromHex(bet.UserID)
		err = s.client.Database("users-db").Collection("users").FindOne(sessionContext, bson.M{"_id": userID}).Decode(&user)
		if err != nil {
			return nil, err
		}

		betsCollection := s.client.Database("bets-db").Collection("bets")
		userCollection := s.client.Database("users-db").Collection("users")

		// Update user based on bet status
		dataMap := make(map[string]interface{})
		if err := c.BodyParser(&dataMap); err != nil {
			return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Access parsed data
		status, ok := dataMap["status"].(string)
		if !ok {
			return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid 'status' field in request body"})
		}

		// Check if the status is the same as the current bet status
		// You can replace this logic with your specific use case
		if status == bet.Status {
			return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No new status update"})
		}

		switch dataMap["status"] {
		case "Won":
			fmt.Println("marking bet as won.")
			user.CurrentAvailability += toWinAmount
			user.CurrentBalance += toWinAmount
			user.CurrentPending -= betAmt
		case "Lost":
			fmt.Println("marking bet as lost.")
			user.CurrentPending -= betAmt
			user.CurrentBalance -= betAmt
			// availability is removed on bet placed -> pending
		case "Pushed":
			fmt.Println("marking bet as pushed.")
			user.CurrentPending -= betAmt
			user.CurrentAvailability += betAmt
		default:
			return nil, fmt.Errorf("Invalid bet status")
		}

		bet.Status = dataMap["status"].(string)
		// if err != nil {
		// 	return err
		// }
		// Update user in the database
		_, err = userCollection.UpdateOne(sessionContext, bson.M{"_id": userID}, bson.M{"$set": user})
		if err != nil {
			return nil, err
		}

		// Update bet status
		_, err = betsCollection.UpdateOne(sessionContext, bson.M{"_id": bet.ID}, bson.M{"$set": bson.M{"bet_status": bet.Status}})
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Transaction failed"})
	}

	return nil
}

func (s *AdminServer) GetAllBets(c *fiber.Ctx) error {
	ValidateUserAdmin(c)

	coll := s.client.Database("bets-db").Collection("bets")
	docs := []bets.Bet{}

	opts := options.Find().SetSort(bson.D{{"placed_at", -1}})
	cursor, err := coll.Find(context.Background(), bson.D{{}}, opts)
	if err != nil {
		return err
	}

	for cursor.Next(context.Background()) {
		var doc bets.Bet
		if err := cursor.Decode(&doc); err != nil {
			fmt.Println("Error decoding document:", err)
			continue
		}
		docs = append(docs, doc)
	}
	c.JSON(docs)

	return nil

}

func (s *AdminServer) GetAllUsers(c *fiber.Ctx) error {
	ValidateUserAdmin(c)
	coll := s.client.Database("users-db").Collection("users")
	docs := []users.User{}

	// Define the options to sort by date_joined in descending order
	opts := options.Find().SetSort(bson.D{{"date_joined", -1}})

	// Perform the MongoDB query with the sorting option
	cursor, err := coll.Find(context.Background(), bson.D{{}}, opts)
	if err != nil {
		return err
	}

	for cursor.Next(context.Background()) {
		var doc users.User
		if err := cursor.Decode(&doc); err != nil {
			fmt.Println("Error decoding document:", err)
			continue
		}
		docs = append(docs, doc)
	}
	c.JSON(docs)

	return nil
}

func (s *AdminServer) CompleteBalanceMetric(c *fiber.Ctx) error {
	ValidateUserAdmin(c)

	// Get the ID from the request URL
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID parameter is required"})
	}

	// Convert the ID parameter to ObjectID
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID parameter"})
	}

	// Define the filter to find the document by ID
	filter := bson.M{"_id": objectID}

	// Define the update to set the "paid" field to true and create it if it doesn't exist
	update := bson.M{
		"$set": bson.M{"paid": true},
	}

	// Update the document in the MongoDB collection
	coll := s.client.Database("weekly-metrics").Collection("balances")
	result, err := coll.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		// Document not found
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Document updated successfully"})
}

func (s *AdminServer) GetAllBalanceMetrics(c *fiber.Ctx) error {
	ValidateUserAdmin(c)

	coll := s.client.Database("weekly-metrics").Collection("balances")

	// Calculate the date one week ago
	oneWeekAgo := time.Now().AddDate(0, 0, -7)

	// Define the filter to select documents where current_date is less than a week old
	filter := bson.D{{Key: "current_date", Value: bson.D{{Key: "$gte", Value: oneWeekAgo}}}}

	docs := []users.WeeklyBalanceMetric{}
	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		return err
	}

	for cursor.Next(context.Background()) {
		var doc users.WeeklyBalanceMetric
		if err := cursor.Decode(&doc); err != nil {
			fmt.Println("Error decoding document:", err)
			continue
		}
		docs = append(docs, doc)
	}
	c.JSON(docs)

	return nil
}

func (s *AdminServer) GetAllMessages(c *fiber.Ctx) error {
	ValidateUserAdmin(c)

	coll := s.client.Database("messages").Collection("contact")
	msgs := []messages.ContactMessage{}
	cursor, err := coll.Find(context.Background(), bson.D{{}})
	if err != nil {
		fmt.Println(err)
		return err
	}

	for cursor.Next(context.Background()) {
		var msg messages.ContactMessage
		if err := cursor.Decode(&msg); err != nil {
			fmt.Println("Error decoding document:", err)
			continue
		}
		msgs = append(msgs, msg)
	}
	c.JSON(msgs)
	return nil
}

func (s *AdminServer) Start(api fiber.Router) error {
	s.StartAdminServerAPI(api)
	return nil
}

func ValidateUserAdmin(c *fiber.Ctx) {
	secretKey := os.Getenv("DASHBOARD_SECRET_KEY")
	userSecretKey := c.Get("secret-key")
	if secretKey != userSecretKey {
		c.SendStatus(http.StatusForbidden)
	}
}
