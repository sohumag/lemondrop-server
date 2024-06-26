package mailing

import (
	"context"
	"crypto/tls"
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
	"gopkg.in/gomail.v2"
)

type MailingServer struct {
	client *mongo.Client
}

func NewMailingServer() *MailingServer {
	return &MailingServer{
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

func (s *MailingServer) StartMailingServerAPI(api fiber.Router) error {
	log.Println("Adding mailing server endpoints to API")

	g := api.Group("/mailing")
	g.Post("/add", func(c *fiber.Ctx) error {
		return s.AddUserToMailingList(c)
	})

	return nil
}

type Email struct {
	Email      string    `json:"email" bson:"email"`
	JoinedDate time.Time `json:"joined_date" bson:"joined_date"`
}

func (s *MailingServer) AddUserToMailingList(c *fiber.Ctx) error {
	coll := s.client.Database("mailing").Collection("emails")
	email := Email{}
	email.JoinedDate = time.Now()
	c.BodyParser(&email)
	// fmt.Println(email)

	if email.Email == "" {
		c.SendStatus(http.StatusBadRequest)
		return fmt.Errorf("Missing required field email")
	}

	email.Email = strings.ToLower(email.Email)

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "email", Value: email.Email}, {Key: "joined_date", Value: email.JoinedDate}}}}
	opts := options.Update().SetUpsert(true)
	coll.UpdateOne(context.TODO(), bson.M{"email": email.Email}, update, opts)

	s.SendJoinedMailingListEmail(email.Email, c)

	return nil
}

func (s *MailingServer) SendJoinedMailingListEmail(email string, c *fiber.Ctx) {
	sendingEmail := os.Getenv("IMPROVMX_CLIENT_SERVICES_EMAIL")
	password := os.Getenv("IMPROVMX_CLIENT_SERVICES_PASSWORD")
	d := gomail.NewDialer("smtp.improvmx.com", 587, sendingEmail, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	subjectLine := "Lemondrop: Congrats! You're In."
	m.SetHeader("From", sendingEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subjectLine)
	body, _ := GetMailingListEmail()
	m.SetBody("text/html", body)

	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	}

	return

}

func (s *MailingServer) Start(api fiber.Router) error {
	s.StartMailingServerAPI(api)
	return nil
}
