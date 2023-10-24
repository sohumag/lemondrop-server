package users

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	// User Added
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	JWT         string `json:"jwt"`

	// Collected
	UserId              int64     `json:"user_id"`
	DateJoined          time.Time `json:"date_joined"`
	CurrentBalance      float64   `json:"current_balance"`
	CurrentAvailability float64   `json:"current_availability"`
	CurrentFreePlay     float64   `json:"current_free_play"`
	CurrentPending      float64   `json:"current_pending"`
}

type UserServer struct {
	client *mongo.Client
}

type ClientUser struct {
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	PhoneNumber         string    `json:"phone_number"`
	Email               string    `json:"email"`
	JWT                 string    `json:"jwt"`
	UserId              int64     `json:"user_id" bson:"user_id"`
	DateJoined          time.Time `json:"date_joined"`
	CurrentBalance      float64   `json:"current_balance"`
	CurrentAvailability float64   `json:"current_availability"`
	CurrentFreePlay     float64   `json:"current_free_play"`
	CurrentPending      float64   `json:"current_pending"`
}
