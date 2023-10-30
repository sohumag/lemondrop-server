package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rlvgl/bookie-server/assetLookup"
	"github.com/rlvgl/bookie-server/bets"
	"github.com/rlvgl/bookie-server/games"
	"github.com/rlvgl/bookie-server/markets"
	"github.com/rlvgl/bookie-server/news"
	"github.com/rlvgl/bookie-server/users"
	"github.com/rlvgl/bookie-server/wheels"
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

	// ARCHIVED GAME SERVER ------------------------
	// gs := games.NewGameServer()
	// gs.Start(api)

	// go func() {
	// 	s := gocron.NewScheduler(time.UTC)
	// 	s.Every(4).Hours().Do(func() {
	// 		gs.MigrateAllGames()
	// 	})
	// }()

	// NEWS SERVER ------------------------
	ns := news.NewNewsServer()
	ns.Start(api)

	// USER SERVER ------------------------
	us := users.NewUserServer()
	us.Start(api)

	// WHEELS SERVER ----------------------
	ws := wheels.NewWheelServer()
	ws.Start(api)

	bs := bets.NewBetServer()
	bs.Start(api)

	ms := markets.NewMarketServer()
	ms.Start(api)

	as := assetLookup.NewAssetServer()
	as.Start(api)

	// ps := props.NewPropServer()
	// ps.Start(api)
	// ps.GetAllPropsGames()
	// ps.ScrapeAllBovadaSports()

	// log.Printf("Starting API on port %d\n", port)
	app.Listen(fmt.Sprintf(":%d", port))

	return nil
}
