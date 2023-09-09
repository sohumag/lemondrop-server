package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rlvgl/bookie-server/games"
)

func StartAPI(port int) error {
	app := fiber.New()
	app.Use(cors.New())

	// API Group Router
	api := app.Group("/api")
	app.Get("/api", func(c *fiber.Ctx) error {
		return c.SendString("API is running")
	})

	// GAME SERVER ------------------------
	gs := games.NewGameServer()
	gs.Start(api)

	log.Printf("Starting API on port %d\n", port)
	app.Listen(fmt.Sprintf(":%d", port))

	return nil
}
