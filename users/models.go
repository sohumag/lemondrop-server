package users

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	// User Added
	FirstName   string `json:"first_name" bson:"first_name"`
	LastName    string `json:"last_name" bson:"last_name"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
	Email       string `json:"email" bson:"email"`
	Password    string `json:"password" bson:"password"`
	JWT         string `json:"jwt" bson:"jwt"`

	// Collected
	UserId              primitive.ObjectID `json:"user_id" bson:"_id"`
	DateJoined          time.Time          `json:"date_joined" bson:"date_joined"`
	CurrentBalance      float64            `json:"current_balance" bson:"current_balance"`
	CurrentAvailability float64            `json:"current_availability" bson:"current_availability"`
	CurrentFreePlay     float64            `json:"current_free_play" bson:"current_free_play"`
	CurrentPending      float64            `json:"current_pending" bson:"current_pending"`
	TotalProfit         float64            `json:"total_profit" bson:"total_profit"`
	StripeCustomerId    string             `json:"stripe_customer_id" bson:"stripe_customer_id"`

	ReferralCode     string `json:"referral_code" bson:"referral_code"`
	ReferredFromCode string `json:"referred_from_code" bson:"referred_from_code"`
}

type UserServer struct {
	client *mongo.Client
}

type ClientUser struct {
	FirstName    string `json:"first_name" bson:"first_name"`
	LastName     string `json:"last_name" bson:"last_name"`
	PhoneNumber  string `json:"phone_number" bson:"phone_number"`
	Email        string `json:"email" bson:"email"`
	JWT          string `json:"jwt" bson:"jwt"`
	ReferralCode string `json:"referral_code" bson:"referral_code"`

	UserId              primitive.ObjectID `json:"user_id" bson:"_id" bson:"user_id"`
	DateJoined          time.Time          `json:"date_joined" bson:"date_joined"`
	CurrentBalance      float64            `json:"current_balance" bson:"current_balance"`
	CurrentAvailability float64            `json:"current_availability" bson:"current_availability"`
	CurrentFreePlay     float64            `json:"current_free_play" bson:"current_free_play"`
	CurrentPending      float64            `json:"current_pending" bson:"current_pending"`
	TotalProfit         float64            `json:"total_profit" bson:"total_profit"`
}
