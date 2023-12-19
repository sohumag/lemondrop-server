package users

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	pieces["email"] = strings.ToLower(pieces["email"])

	stripe.Key = "sk_test_WQ4y1OC1xfTS8CCcu8nTKf29"

	params := &stripe.CustomerParams{
		Email: stripe.String(pieces["email"]),
		Name:  stripe.String(pieces["first_name"] + pieces["last_name"]),
		Phone: stripe.String(pieces["phone_number"]),
	}
	cust, _ := customer.New(params)
	fmt.Println(cust.ID)

	user := User{
		FirstName:           pieces["first_name"],
		LastName:            pieces["last_name"],
		PhoneNumber:         pieces["phone_number"],
		Email:               pieces["email"],
		Password:            pieces["password"],
		UserId:              primitive.NewObjectID(),
		DateJoined:          time.Now(),
		CurrentBalance:      0,
		CurrentAvailability: 0,
		CurrentFreePlay:     0,
		CurrentPending:      0,
		TotalProfit:         0,
		StripeCustomerId:    cust.ID,
	}

	fmt.Println(user)

	// encrypting password
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(pieces["password"]), 10)
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
			fmt.Println("invalid token")
			c.SendStatus(http.StatusBadRequest)
		}
		if err.Error() == "no token in header" {
			// fmt.Println("login without jwt")
			fmt.Println("no token in header")
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
	var data map[string]string
	// Unmarshal JSON into the struct
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	email := strings.ToLower(data["email"])
	password := data["password"]

	user := User{}

	coll := s.client.Database("users-db").Collection("users")
	err = coll.FindOne(context.TODO(), bson.D{{Key: "email", Value: email}}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		c.SendStatus(http.StatusNotFound)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err == nil {
		// user logged in correctly
		jwt, err := GenerateJWT(user.Email)
		if err != nil {
			return err
		}

		cuser := ClientUser{
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
			JWT:         jwt,

			UserId:              user.UserId,
			DateJoined:          user.DateJoined,
			CurrentBalance:      user.CurrentBalance,
			CurrentAvailability: user.CurrentAvailability,
			CurrentFreePlay:     user.CurrentFreePlay,
			CurrentPending:      user.CurrentPending,
		}

		user.JWT = jwt

		c.JSON(cuser)
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
