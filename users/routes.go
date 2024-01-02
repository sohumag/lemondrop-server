package users

import (
	"context"
	"fmt"
	"net/http"
	"os"
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

func (s *UserServer) HandleSignUpRoute(c *fiber.Ctx) error {
	// Parse JSON request body into the map
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	// Normalize email to lowercase
	data["email"] = strings.ToLower(data["email"])

	// Set up Stripe API key
	stripe.Key = os.Getenv("STRIPE_SECRET_TEST_KEY")

	// params := &stripe.AccountParams{Type: stripe.String(string(stripe.AccountTypeExpress))}
	// result, err := account.New(params)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }
	// accountId := result.ID
	// fmt.Println(accountId)

	// linkParams := &stripe.AccountLinkParams{
	// 	Account:    stripe.String(accountId),
	// 	RefreshURL: stripe.String("https://example.com/reauth"),
	// 	ReturnURL:  stripe.String("https://example.com/return"),
	// 	Type:       stripe.String("account_onboarding"),
	// }
	// _, err = accountlink.New(linkParams)

	// Create a customer on Stripe
	params := &stripe.CustomerParams{
		Email: stripe.String(data["email"]),
		Name:  stripe.String(data["first_name"] + data["last_name"]),
		Phone: stripe.String(data["phone_number"]),
	}
	cust, _ := customer.New(params)
	fmt.Println(cust.ID)

	// Create a User struct
	user := User{
		FirstName:           data["first_name"],
		LastName:            data["last_name"],
		PhoneNumber:         data["phone_number"],
		Email:               data["email"],
		Password:            data["password"],
		UserId:              primitive.NewObjectID(),
		DateJoined:          time.Now(),
		CurrentBalance:      0,
		CurrentAvailability: 0,
		CurrentFreePlay:     0,
		CurrentPending:      0,
		TotalProfit:         0,
		StripeCustomerId:    cust.ID,
	}

	// Encrypt password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 10)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)

	// Add user to the database
	coll := s.client.Database("users-db").Collection("users")
	res, err := coll.InsertOne(context.TODO(), &user)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Added user with id: %v\n", res.InsertedID)

	// Generate JWT token
	jwt, err := GenerateJWT(user.Email)
	if err != nil {
		return err
	}

	// Send JWT token in the response
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
