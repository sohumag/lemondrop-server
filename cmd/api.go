package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/rlvgl/bookie-server/admin"
	"github.com/rlvgl/bookie-server/assetLookup"
	"github.com/rlvgl/bookie-server/bets"
	"github.com/rlvgl/bookie-server/games"
	"github.com/rlvgl/bookie-server/mailing"
	"github.com/rlvgl/bookie-server/markets"
	"github.com/rlvgl/bookie-server/messages"
	"github.com/rlvgl/bookie-server/users"
	"github.com/rlvgl/bookie-server/verification"
)

func validateApiKey(c *fiber.Ctx, key string) (bool, error) {
	hashedApiKey := sha256.Sum256([]byte(os.Getenv("LEMONDROP_API_KEY")))
	hashedKey := sha256.Sum256([]byte(key))

	if subtle.ConstantTimeCompare(hashedApiKey[:], hashedKey[:]) == 1 {
		return true, nil
	}

	return false, keyauth.ErrMissingOrMalformedAPIKey
}

func StartAPI(port int) error {
	app := fiber.New()
	app.Use(cors.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://lemondrop.bet, https://lemondrop.ag, http://localhost:3000, http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// app.Use(keyauth.New(keyauth.Config{
	// 	KeyLookup: "header:lemondrop_api_key",
	// 	Validator: validateApiKey,
	// }))

	// API Group Router
	api := app.Group("/api")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API is running")
	})

	adm := admin.NewAdminServer()
	adm.Start(api)

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
