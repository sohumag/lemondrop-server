package games

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// * GET ALL CURRENT CATEGORIES
func (g *GameServer) SendAllSportsCategories(c *fiber.Ctx) error {

	categoryStrs := []string{
		"American Football",
		"Baseball",
		"Basketball",
		"Boxing",
		"Cricket",
		"Golf",
		"Ice Hockey",
		"Mixed Martial Arts",
		"Rubgy League",
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
