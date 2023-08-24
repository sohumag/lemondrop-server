package games

import (
	"log"

	"github.com/gofiber/fiber/v2"
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
	app.Get("/api", func(c *fiber.Ctx) error {
		return c.SendString("API is running")
	})

	// SPORTS/CATEGORIES ---------------------

	app.Get("/api/categories/all", func(c *fiber.Ctx) error {
		return g.SendAllSportsCategories(c)
	})

	app.Get("/api/sports/:category", func(c *fiber.Ctx) error {
		return g.SendSportsFromCategories(c, c.Params("category"))
	})

	// GAMES --------------------------------

	// app.Get("/api/games/")

	app.Get("/api/games/all", func(c *fiber.Ctx) error {
		return g.SendGamesUpcoming(c)
	})

	app.Get("/api/games/:sport", func(c *fiber.Ctx) error {
		return g.SendGamesBySport(c, c.Params("sport"))
	})

	app.Listen(g.port)

	return nil
}

func (g *GameServer) SendGamesBySport(c *fiber.Ctx, sport string) error {
	games, err := g.GetUpcomingGamesBySport(sport)
	if err != nil {
		return err
	}

	c.JSON(games)
	return nil
}

func (g *GameServer) SendGamesUpcoming(c *fiber.Ctx) error {
	games, err := g.GetAllUpcomingGames(50)
	if err != nil {
		return err
	}

	c.JSON(games)

	return nil
}
