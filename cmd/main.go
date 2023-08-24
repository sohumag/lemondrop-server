package main

import (
	"github.com/joho/godotenv"
	"github.com/rlvgl/bookie-server/games"
)

func main() {
	godotenv.Load("./.env")

	gs := games.NewGameServer(":8080")
	gs.Start()

}
