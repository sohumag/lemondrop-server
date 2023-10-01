package wheels

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Wheel struct {
	PurchaseCost int       `json:"purchase_cost" bson:"purchase_cost"`
	Values       []float64 `json:"values" bson:"values"`
	Name         string    `json:"name" bson:"name"`
}

type WheelServer struct {
	client *mongo.Client
}

func NewWheelServer() *WheelServer {
	return &WheelServer{client: ConnectDB()}
}

func ConnectDB() *mongo.Client {
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client
}

func (s *WheelServer) StartWheelServerAPI(api fiber.Router) error {
	log.Println("Adding wheel server endpoints to API")

	// USERS
	wheelsApi := api.Group("/wheels")
	wheelsApi.Get("/all", func(c *fiber.Ctx) error {
		return s.GetAllWheels(c)
	})

	return nil
}

func (s *WheelServer) GetAllWheels(c *fiber.Ctx) error {
	coll := s.client.Database("gambling").Collection("wheels")
	cursor, err := coll.Find(context.TODO(), bson.D{}, options.Find())
	if err != nil {
		fmt.Println(err)
		return err
	}
	// decode cursor all
	wheels := []Wheel{}
	for cursor.Next(context.TODO()) {
		wheel := Wheel{}
		cursor.Decode(&wheel)
		wheels = append(wheels, wheel)
	}
	fmt.Println(wheels)
	json.NewEncoder(c).Encode(wheels)
	return nil
}

func (s *WheelServer) Start(api fiber.Router) error {
	s.StartWheelServerAPI(api)
	return nil
}
