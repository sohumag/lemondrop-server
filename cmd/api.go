package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rlvgl/bookie-server/assetLookup"
	"github.com/rlvgl/bookie-server/bets"
	"github.com/rlvgl/bookie-server/games"
	"github.com/rlvgl/bookie-server/markets"
	"github.com/rlvgl/bookie-server/payments"
	"github.com/rlvgl/bookie-server/verification"

	// "github.com/rlvgl/bookie-server/scores"
	"github.com/rlvgl/bookie-server/users"
)

func StartAPI(port int) error {
	app := fiber.New()
	app.Use(cors.New())

	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: "http://localhost:5173, https://lemondrop.bet",
	// 	AllowHeaders: "Origin, Content-Type, Accept",
	// }))

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

	ps := payments.NewPaymentServer()
	ps.Start(api)

	bs := bets.NewBetServer()
	bs.Start(api)

	ms := markets.NewMarketServer()
	ms.Start(api)

	as := assetLookup.NewAssetServer()
	as.Start(api)

	vs := verification.NewVerificationServer()
	vs.Start(api)

	// log.Printf("Starting API on port %d\n", port)
	app.Listen(fmt.Sprintf(":%d", port))

	return nil
}
