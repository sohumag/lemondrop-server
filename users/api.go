package users

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
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

	usersApi.Post("/account-link/:accountId", func(c *fiber.Ctx) error {
		fmt.Println("linking account..")
		accountId := c.Params("accountId")
		url, err := s.HandleAccountLink(accountId)
		if err != nil {
			return err
		}
		c.JSON(fiber.Map{"url": url})
		return nil
	})

	return nil
}

func (s *UserServer) Start(api fiber.Router) error {
	s.StartUserServerAPI(api)
	return nil
}
