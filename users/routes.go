package users

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// @ WORKS
func (s *UserServer) HandleSignUpRoute(c *fiber.Ctx) error {
	pieces := ParseRequestBody(c.Body())

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

	// encrypting password
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(pieces["password"]), 15)
	user.Password = string(passwordHash)

	// adding user to database
	coll := s.client.Database("users-db").Collection("users")

	result, err := coll.InsertOne(context.TODO(), &user)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("added user with id: %v\n", result.InsertedID)

	jwt, err := GenerateJWT(user.Email)
	if err != nil {
		return err
	}

	c.Send([]byte(jwt))

	return nil
}

func (s *UserServer) HandleLoginRoute(c *fiber.Ctx) error {
	// login with jwt
	email, err := ParseRequestForJWT(c)
	if err != nil {
		// invalid token
		if err.Error() == "invalid token" {
			c.SendStatus(http.StatusBadRequest)
		}
		if err.Error() == "no token in header" {
			// fmt.Println("login without jwt")
			s.HandleLoginWithoutJWT(c)
		}
	} else {
		s.HandleLoginWithJWT(c, email)
	}
	return nil
}

func (s *UserServer) HandleLoginWithJWT(c *fiber.Ctx, email string) error {
	coll := s.client.Database("users-db").Collection("users")
	user := ClientUser{}
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "email", Value: email}}).Decode(&user); err != nil {
		c.SendStatus(http.StatusInternalServerError)
	}

	c.JSON(user)
	return nil
}

func (s *UserServer) HandleLoginWithoutJWT(c *fiber.Ctx) error {
	pieces := ParseRequestBody(c.Body())
	for _, val := range []string{"email", "password"} {
		if _, ok := pieces[val]; !ok {
			c.SendStatus(http.StatusBadRequest)
			return fiber.ErrBadRequest
		}
	}

	email := pieces["email"]
	password := pieces["password"]

	user := User{}
	coll := s.client.Database("users-db").Collection("users")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "email", Value: email}}).Decode(&user); err != nil {
		c.SendStatus(http.StatusInternalServerError)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err == nil {
		// user logged in correctly
		c.JSON(user)
	} else {
		c.SendStatus(http.StatusBadRequest)
	}

	return nil
}

func ParseRequestBody(body []byte) map[string]string {
	bodyRaw := strings.Trim(string(body), " \n{}")
	toks := strings.Split(bodyRaw, ",")
	pieces := map[string]string{}
	for _, tok := range toks {
		arr := strings.Split(tok, ":")
		if len(arr) >= 2 {
			key := strings.Trim(arr[0], " \"\n")
			val := strings.Trim(arr[1], " \"\n")
			pieces[key] = val
		}
	}

	return pieces
}
