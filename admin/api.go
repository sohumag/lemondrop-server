package admin

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
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
