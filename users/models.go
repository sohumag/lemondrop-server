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
	DateJoined          time.Time          `json:"date_joined"`
	CurrentBalance      float64            `json:"current_balance"`
	CurrentAvailability float64            `json:"current_availability"`
	CurrentFreePlay     float64            `json:"current_free_play"`
	CurrentPending      float64            `json:"current_pending"`
}

type UserServer struct {
	client *mongo.Client
}

type ClientUser struct {
	FirstName           string             `json:"first_name"`
	LastName            string             `json:"last_name"`
	PhoneNumber         string             `json:"phone_number"`
	Email               string             `json:"email"`
	JWT                 string             `json:"jwt"`
	UserId              primitive.ObjectID `json:"user_id" bson:"_id"`
	DateJoined          time.Time          `json:"date_joined"`
	CurrentBalance      float64            `json:"current_balance"`
	CurrentAvailability float64            `json:"current_availability"`
	CurrentFreePlay     float64            `json:"current_free_play"`
	CurrentPending      float64            `json:"current_pending"`
}
