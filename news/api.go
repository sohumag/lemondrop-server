package news

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func (s *NewsServer) StartNewsServerAPI(api fiber.Router) error {
	log.Println("Adding news server endpoints to API")

	newsApi := api.Group("/news")
	newsApi.Get("/", func(c *fiber.Ctx) error {
		c.JSON(map[string]string{"message": "news api is working"})
		return nil
	})

	return nil
}
