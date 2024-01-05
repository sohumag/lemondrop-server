package users

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/account"
	"github.com/stripe/stripe-go/v76/accountlink"
	"github.com/stripe/stripe-go/v76/customer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func GenerateReferralCode() (string, error) {
	// Define the length of the referral code
	codeLength := 6

	// Define the character set for the code
	charSet := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Calculate the required number of bytes
	byteLength := (codeLength * 6) / 8

	// Generate random bytes
	randomBytes := make([]byte, byteLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Convert random bytes to a big integer
	randomInt := new(big.Int).SetBytes(randomBytes)

	// Generate the referral code using the character set
	referralCode := make([]byte, codeLength)
	for i := range referralCode {
		index := new(big.Int).Mod(randomInt, big.NewInt(int64(len(charSet))))
		referralCode[i] = charSet[index.Int64()]
		randomInt.Div(randomInt, big.NewInt(int64(len(charSet))))
	}

	return string(referralCode), nil
}

func IsValidEmail(email string) bool {
	// Define the email validation regex pattern
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regex pattern
	regex := regexp.MustCompile(emailPattern)

	// Test if the email matches the pattern
	return regex.MatchString(email)
}

// returns Account Link URL, error
func (s *UserServer) HandleAccountLink(accountId string) (string, error) {
	accountLinkParams := &stripe.AccountLinkParams{
		Account:    stripe.String(accountId),
		RefreshURL: stripe.String("https://lemondrop.ag/bets"),
		ReturnURL:  stripe.String("https://lemondrop.ag/bets"),
		Type:       stripe.String("account_onboarding"),
		Collect:    stripe.String("eventually_due"),
	}

	accountLinkResult, err := accountlink.New(accountLinkParams)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return accountLinkResult.URL, nil
}

func (s *UserServer) HandleSignUpRoute(c *fiber.Ctx) error {
	referralCode, err := GenerateReferralCode()
	if err != nil {
		return err
	}

	// Parse JSON request body into the map
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	referredFromCode := data["referred_from_code"]

	// Normalize email to lowercase
	data["email"] = strings.ToLower(data["email"])
	if !IsValidEmail(data["email"]) {
		return fmt.Errorf("invalid email")
	}
	email := data["email"]

	// Set up Stripe API key
	stripe.Key = os.Getenv("STRIPE_SECRET_TEST_KEY")

	accountParams := &stripe.AccountParams{
		Country: stripe.String("US"),
		Type:    stripe.String(string(stripe.AccountTypeExpress)),
		Capabilities: &stripe.AccountCapabilitiesParams{
			CardPayments: &stripe.AccountCapabilitiesCardPaymentsParams{
				Requested: stripe.Bool(true),
			},
			Transfers: &stripe.AccountCapabilitiesTransfersParams{Requested: stripe.Bool(true)},
		},
		BusinessType: stripe.String(string(stripe.AccountBusinessTypeIndividual)),
		Email:        &email,
		// BusinessProfile: &stripe.AccountBusinessProfileParams{
		// 	URL: stripe.String("https://example.com"),
		// },
		BusinessProfile: &stripe.AccountBusinessProfileParams{MCC: stripe.String("7995")},
	}
	result, err := account.New(accountParams)
	if err != nil {
		fmt.Println(err)
	}
	accountId := result.ID
	fmt.Println(accountId)

	redirectUrl, err := s.HandleAccountLink(accountId)

	// Create a customer on Stripe
	params := &stripe.CustomerParams{
		Email: stripe.String(data["email"]),
		Name:  stripe.String(data["first_name"] + data["last_name"]),
		Phone: stripe.String(data["phone_number"]),
	}
	cust, _ := customer.New(params)

	// Create a User struct
	user := User{
		FirstName:   data["first_name"],
		LastName:    data["last_name"],
		PhoneNumber: data["phone_number"],
		Email:       data["email"],
		Password:    data["password"],
		UserId:      primitive.NewObjectID(),
		DateJoined:  time.Now(),

		CurrentBalance:      0,
		CurrentAvailability: 0,
		CurrentFreePlay:     0,
		CurrentPending:      0,
		TotalProfit:         0,

		StripeCustomerId: cust.ID,
		StripeExpressId:  accountId,
		DetailsSubmitted: false,

		ReferralCode:     referralCode,
		ReferredFromCode: referredFromCode,
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

	c.JSON(fiber.Map{
		"jwt":          jwt,
		"redirect_url": redirectUrl,
	})

	// dataMap := map[interface{}]interface{}{
	// 	"jwt":          jwt,
	// 	"redirect_url": redirectUrl,
	// 	// "stripe_account_id": accountId,
	// 	// "details_submitted": false,
	// }

	// fmt.Println(dataMap)

	// c.JSON(dataMap)
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
			FirstName:           user.FirstName,
			LastName:            user.LastName,
			PhoneNumber:         user.PhoneNumber,
			Email:               user.Email,
			JWT:                 jwt,
			ReferralCode:        user.ReferralCode,
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
