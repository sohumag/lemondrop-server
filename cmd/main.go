package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("./.env")

	port := flag.Int("p", 8080, "port to listen on")
	flag.Parse()

	// gs.MigrateAllGames()
	// gs.RunGameScoreUpdates()
	// err := gs.GetScoresForGames()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	log.Fatal(StartAPI(*port))

}
