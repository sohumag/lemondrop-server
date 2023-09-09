package games

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (g *GameServer) SendAllSports(c *fiber.Ctx) error {

	type SportSend struct {
		Group string `json:"group"`
		Title string `json:"title"`
		Id    string `json:"id"`
	}

	games := []SportSend{
		{
			Group: "Football",
			Title: "NCAAF",
			Id:    "americanfootball_ncaaf",
		},
		{
			Group: "Football",
			Title: "NFL",
			Id:    "americanfootball_nfl",
		},
		{
			Group: "Football",
			Title: "NFL Preseason",
			Id:    "americanfootball_nfl_preseason",
		},
		{
			Group: "Baseball",
			Title: "MLB",
			Id:    "baseball_mlb",
		},
		{
			Group: "Basketball",
			Title: "NBA",
			Id:    "basketball_nba",
		},
	}

	c.JSON(games)

	return nil
}

// * GET ALL CURRENT CATEGORIES
func (g *GameServer) SendAllSportsCategories(c *fiber.Ctx) error {

	categoryStrs := []string{
		"All",
		"American Football",
		"Baseball",
		"Basketball",
		// "Boxing",
		// "Cricket",
		// "Golf",
		// "Ice Hockey",
		// "Mixed Martial Arts",
		// "Rubgy League",
	}

	categories := []Category{}

	for _, c := range categoryStrs {
		categories = append(categories, Category{Name: c})
	}

	c.JSON(categories)
	return nil
}

// * GET ALL SPORTS IN CATEGORY
func (g *GameServer) SendSportsFromCategories(c *fiber.Ctx, category string) error {
	category = strings.Replace(category, "%20", " ", -1)

	allSports, err := g.GetAllSportsInCategory(category)
	if err != nil {
		return err
	}

	sportsInCategory := []Sport{}
	for _, sport := range allSports {
		if sport.Group == category {
			sportsInCategory = append(sportsInCategory, sport)
		}
	}

	if len(sportsInCategory) == 0 {
		c.SendString("Invalid category. check /api/categories/all for all valid categories")
	}

	c.JSON(sportsInCategory)

	return nil
}
