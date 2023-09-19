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

func (g *GameServer) StartGameServerAPI(api fiber.Router) error {
	log.Println("Adding game server endpoints to API")

	// CATEGORIES ---------------------
	categoriesApi := api.Group("/categories")
	categoriesApi.Get("/all", func(c *fiber.Ctx) error {
		return g.SendAllSportsCategories(c)
	})

	// SPORTS ----------------------------
	sportsApi := api.Group("/sports")
	sportsApi.Get("/all", func(c *fiber.Ctx) error {
		return g.SendAllSports(c)
	})

	sportsApi.Get("/:category", func(c *fiber.Ctx) error {
		return g.SendSportsFromCategories(c, c.Params("category"))
	})

	// GAMES --------------------------------
	gamesApi := api.Group("/games")
	gamesApi.Get("/all", func(c *fiber.Ctx) error {
		return g.SendGamesUpcoming(c)
	})

	gamesApi.Get("/All", func(c *fiber.Ctx) error {
		return g.SendGamesUpcoming(c)
	})

	gamesApi.Get("/:sport", func(c *fiber.Ctx) error {
		return g.SendGamesBySport(c, c.Params("sport"))
	})

	gamesApi.Get("/game/:id", func(c *fiber.Ctx) error {
		return g.SendGameById(c, c.Params("id"))
	})

	return nil
}
