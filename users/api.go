package users

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
	new user no jwt logs in/signs up-> return jwt and user info if valid
	new user no jwt fails log in -> return err
	new user with jwt logs in -> return user info
	new user wit invalid jwt logs in -> return failure. redirect to login

*/

/*
	Referrals:
	- each user has: unique 6 digit code and shareable link

*/

func NewUserServer() *UserServer {
	return &UserServer{
		client: ConnectDB(),
	}
}

func ConnectDB() *mongo.Client {
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client
}

func (s *UserServer) StartUserServerAPI(api fiber.Router) error {
	log.Println("Adding user server endpoints to API")

	// USERS
	usersApi := api.Group("/users")
	usersApi.Post("/signup", func(c *fiber.Ctx) error {
		return s.HandleSignUpRoute(c)
	})

	usersApi.Post("/login", func(c *fiber.Ctx) error {
		return s.HandleLoginRoute(c)
	})

	return nil
}

func (s *UserServer) Start(api fiber.Router) error {
	// err := s.RunPaymentMetricWeeklyReset()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	s.StartUserServerAPI(api)
	return nil
}

func (s *UserServer) StartUserPaymentMetricsRepeater() error {

	c := cron.New()

	// every monday at 12am
	_, err := c.AddFunc("0 0 * * 1", func() {
		s.RunPaymentMetricWeeklyReset()
	})

	if err != nil {
		fmt.Println("Error adding cron job:", err)
		return err
	}

	c.Start()

	select {}
}

func (s *UserServer) RunPaymentMetricWeeklyReset() error {
	// get all users from db
	// for each user,if balance != 0, add user id, email, name, phone number, balance, date to weekly end stats database
	// 		-going to be helpful for admin dashboard
	// reset user balance and availability. add free play if warranted. if pending != 0, add to next week.

	userColl := s.client.Database("users-db").Collection("users")
	filter := bson.D{{Key: "current_balance", Value: bson.D{{Key: "$ne", Value: 0}}}}

	// Query the collection
	cursor, err := userColl.Find(context.Background(), filter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error querying MongoDB: %v\n", err)
		return err
	}
	defer cursor.Close(context.Background())

	// Decode and print the results
	var users []User
	for cursor.Next(context.Background()) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding document: %v\n", err)
			return err
		}
		users = append(users, user)
	}

	balancesColl := s.client.Database("weekly-metrics").Collection("balances")
	for _, user := range users {

		balance := WeeklyBalanceMetric{
			UserId:      user.UserId,
			Email:       user.Email,
			Name:        user.FirstName + " " + user.LastName,
			PhoneNumber: user.PhoneNumber,
			Balance:     user.CurrentBalance,
			CurrentDate: time.Now(),
			Paid:        false,
		}

		insertResult, err := balancesColl.InsertOne(context.Background(), balance)
		if err != nil {
			fmt.Printf("Error inserting document: %v\n", err)
			return err
		}

		fmt.Printf("Inserted document with ID: %v\n", insertResult.InsertedID)

		newBalance := 0.0
		newAvailability := user.MaxAvailability
		newPending := user.CurrentPending
		newFreePlay := 0.0
		if user.CurrentBalance < 0 {
			newFreePlay = -.15 * user.CurrentBalance
		}

		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "current_balance", Value: newBalance},
				{Key: "current_availability", Value: newAvailability},
				{Key: "current_pending", Value: newPending},
				{Key: "current_free_play", Value: newFreePlay},
			}},
		}

		// Define a filter to find the document in "userColl" collection
		filter := bson.D{{Key: "_id", Value: user.UserId}}

		// Update document in "userColl" collection
		_, err = userColl.UpdateOne(context.Background(), filter, update)
		if err != nil {
			fmt.Printf("Error updating document in 'userColl' collection: %v\n", err)
			return err
		}
		fmt.Printf("Updated document in 'userColl' collection for user with ID: %v\n", user.UserId)
	}

	return nil
}
