package games

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func (g *GameServer) AddAllSportsEndpointsToFile() error {
	reqUrl := "https://api.the-odds-api.com/v4/sports?apiKey=e0ae2e9cd2c145da9659ce53ddbc4442"

	res, err := http.Get(reqUrl)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	os.WriteFile("./allsports.txt", body, 0644)

	return nil
}

// need to add all games for all sports
func (g *GameServer) AddAllGamesEndpointsToFile() error {
	reqUrl := "https://api.the-odds-api.com/v4/sports/upcoming/odds/?apiKey=e0ae2e9cd2c145da9659ce53ddbc4442&regions=us&oddsFormat=american&bookmakers=draftkings,fanduel,unibet,sugarhouse,barstool,bovada"

	res, err := http.Get(reqUrl)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	os.WriteFile("./gamesData/allUpcoming.json", body, 0644)

	return nil
}

func (g *GameServer) GetGamesBySport(sportKey string) ([]Game, error) {
	fetchURL := fmt.Sprintf("https://api.the-odds-api.com/v4/sports/%s/odds/?apiKey=e0ae2e9cd2c145da9659ce53ddbc4442&regions=us&markets=h2h,spreads,totals&oddsFormat=american&bookmakers=draftkings,fanduel,unibet,sugarhouse,barstool,bovada", sportKey)
	res, err := http.Get(fetchURL)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	games := []Game{}
	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, &games)

	return games, nil
}
