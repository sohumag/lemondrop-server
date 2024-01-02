package verification

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewVerificationServer() *VerificationServer {
	return &VerificationServer{
		client: ConnectDB(),
	}
}

type VerificationServer struct {
	client *mongo.Client
}

func ConnectDB() *mongo.Client {
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client
}

func (s *VerificationServer) StartVerificationServerAPI(api fiber.Router) error {
	log.Println("Adding verification server endpoints to API")
	// verification
	vapi := api.Group("/verifications")

	vapi.Get("/", func(c *fiber.Ctx) error {
		c.Send([]byte("verifications work"))
		return nil
	})

	vapi.Post("/base", func(c *fiber.Ctx) error {
		return s.SendEmail()
	})

	return nil
}

func (s *VerificationServer) Start(api fiber.Router) error {
	s.StartVerificationServerAPI(api)
	return nil
}

func (s *VerificationServer) SendEmail() error {
	from := mail.NewEmail("Help Email", "help@lemondrop.ag")
	subject := "Please Verify Your Email"
	to := mail.NewEmail("Sohum", "agrawalsohum@gmail.com")
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)

	if err != nil {
		log.Println(err)
		return err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	return nil
}
