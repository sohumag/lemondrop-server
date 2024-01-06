package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rlvgl/bookie-server/assetLookup"
	"github.com/rlvgl/bookie-server/bets"
	"github.com/rlvgl/bookie-server/games"
	"github.com/rlvgl/bookie-server/mailing"
	"github.com/rlvgl/bookie-server/markets"
	"github.com/rlvgl/bookie-server/messages"
	"github.com/rlvgl/bookie-server/users"
	"github.com/rlvgl/bookie-server/verification"
)

func StartAPI(port int) error {
	app := fiber.New()
	app.Use(cors.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://lemondrop.bet, https://lemondrop.ag, http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// API Group Router
	api := app.Group("/api")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API is running")
	})

	gs := games.NewGameServer()
	gs.Start(api)

	// USER SERVER ------------------------
	us := users.NewUserServer()
	us.Start(api)

	bs := bets.NewBetServer()
	bs.Start(api)

	ms := markets.NewMarketServer()
	ms.Start(api)

	mes := messages.NewMessageServer()
	mes.Start(api)

	as := assetLookup.NewAssetServer()
	as.Start(api)

	vs := verification.NewVerificationServer()
	vs.Start(api)

	mas := mailing.NewMailingServer()
	mas.Start(api)

	// log.Printf("Starting API on port %d\n", port)
	app.Listen(fmt.Sprintf(":%d", port))

	return nil
}
