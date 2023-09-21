package users

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
	usersApi.Post("/", func(c *fiber.Ctx) error {
		return s.AddUserToDB(c)
	})

	return nil
}

func (s *UserServer) AddUserToDB(c *fiber.Ctx) error {
	bodyRaw := strings.Trim(string(c.Body()), " \n{}")

	toks := strings.Split(bodyRaw, ",")

	pieces := map[string]string{}

	for _, tok := range toks {
		arr := strings.Split(tok, ":")
		key := strings.Trim(arr[0], " \"\n")
		val := strings.Trim(arr[1], " \"\n")
		pieces[key] = val
	}

	// data validation, regex validation on frontend
	for _, val := range []string{"last_name", "first_name", "phone_number", "email", "password"} {
		if _, ok := pieces[val]; !ok {
			c.SendStatus(http.StatusBadRequest)
			return fiber.ErrBadRequest
		}
	}

	user := User{
		FirstName:   pieces["first_name"],
		LastName:    pieces["last_name"],
		PhoneNumber: pieces["phone_number"],
		Email:       pieces["email"],
		// encrypt!!!
		Password: pieces["password"],

		DateJoined:          time.Now(),
		CurrentBalance:      0,
		CurrentAvailability: 0,
		CurrentFreePlay:     0,
		CurrentPending:      0,
	}

	// adding user to database
	coll := s.client.Database("users-db").Collection("users")

	if err := coll.FindOne(context.TODO(), bson.D{{Key: "email", Value: user.Email}}); err != nil {
		// user is new
		result, err := coll.InsertOne(context.TODO(), &user)
		if err != nil {
			return err
		}

		fmt.Printf("added user with id: %v\n", result.InsertedID)
	} else {
		// user exists already
		c.SendStatus(http.StatusBadRequest)
		return fiber.ErrBadRequest
	}

	return nil
}

func (s *UserServer) Start(api fiber.Router) error {
	s.StartUserServerAPI(api)
	return nil
}
