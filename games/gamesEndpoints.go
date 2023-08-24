package games

import "github.com/gofiber/fiber/v2"

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
