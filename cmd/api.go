package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rlvgl/bookie-server/games"
	"github.com/rlvgl/bookie-server/news"
	"github.com/rlvgl/bookie-server/users"
	"github.com/rlvgl/bookie-server/wheels"
)

func StartAPI(port int) error {
	app := fiber.New()
	app.Use(cors.New())

	// API Group Router
	api := app.Group("/api")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API is running")
	})

	// GAME SERVER ------------------------
	gs := games.NewGameServer()
	gs.Start(api)

	// gs.MigrateAllGames()

	// NEWS SERVER ------------------------
	ns := news.NewNewsServer()
	ns.Start(api)

	// USER SERVER ------------------------
	us := users.NewUserServer()
	us.Start(api)

	// WHEELS SERVER ----------------------
	ws := wheels.NewWheelServer()
	ws.Start(api)

	// log.Printf("Starting API on port %d\n", port)
	app.Listen(fmt.Sprintf(":%d", port))

	return nil
}
