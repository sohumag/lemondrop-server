package main

import (
	"github.com/joho/godotenv"
	"github.com/rlvgl/bookie-server/games"
)

func main() {
	godotenv.Load("./.env")

	gs := games.NewGameServer(":8080")
	// gs.MigrateAllGames()
	gs.Start()

	// gs.RunGameScoreUpdates()
	// err := gs.GetScoresForGames()
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
