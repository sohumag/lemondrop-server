package main

import (
	"flag"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/rlvgl/bookie-server/games"
)

func main() {
	godotenv.Load("./.env")

	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	gs := games.NewGameServer(fmt.Sprintf(":%d", *port))
	gs.Start()

	// gs.MigrateAllGames()
	// gs.RunGameScoreUpdates()
	// err := gs.GetScoresForGames()
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
