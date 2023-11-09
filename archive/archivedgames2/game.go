package games

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (s *GameServer) ReturnGameById(c *fiber.Ctx, id string) error {
	// fmt.Println(id)
	for _, league := range validLeagues {
		allGames := s.cache.gameCache[league]
		for _, game := range allGames {
			if game.Id == id {
				c.JSON(game)
				return nil
			}
		}
	}

	c.SendStatus(http.StatusBadRequest)
	return fmt.Errorf("Invalid game id")
}
