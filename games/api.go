package games

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

/*
* GET all categories ✅
* GET all sports within categories ✅

* GET all games by category
* GET all games by sport
* GET all games general
 */

func (g *GameServer) StartAPI() error {
	log.Println("Starting API on port", g.port)

	app := fiber.New()

	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: "https://gofiber.io, https://gofiber.net",
	// 	AllowHeaders: "Origin, Content-Type, Accept",
	// }))

	app.Use(cors.New())

	app.Get("/api", func(c *fiber.Ctx) error {
		return c.SendString("API is running")
	})

	// SPORTS/CATEGORIES ---------------------

	app.Get("/api/categories/all", func(c *fiber.Ctx) error {
		return g.SendAllSportsCategories(c)
	})

	app.Get("/api/sports/all", func(c *fiber.Ctx) error {
		return g.SendAllSports(c)
	})

	app.Get("/api/sports/:category", func(c *fiber.Ctx) error {
		return g.SendSportsFromCategories(c, c.Params("category"))
	})

	// GAMES --------------------------------

	// app.Get("/api/games/")

	app.Get("/api/games/all", func(c *fiber.Ctx) error {
		return g.SendGamesUpcoming(c)
	})

	app.Get("/api/games/All", func(c *fiber.Ctx) error {
		return g.SendGamesUpcoming(c)
	})

	app.Get("/api/games/:sport", func(c *fiber.Ctx) error {
		return g.SendGamesBySport(c, c.Params("sport"))
	})

	app.Get("/api/games/game/:id", func(c *fiber.Ctx) error {
		return g.SendGameById(c, c.Params("id"))
	})

	app.Listen(g.port)

	return nil
}
