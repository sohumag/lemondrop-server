package games

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func (g *GameServer) MigrateSports() error {
	// protection against expensive api call
	return nil

	// if neeed to remigrate, nmeed to fetch all sports from file again
	// sports, err := g.GetAllSports()
	// if err != nil {
	// 	return err
	// }

	// for _, s := range sports {
	// 	g.AddSportToDB(&s)
	// }

	// return nil
}

func (g *GameServer) MigrateAllGames() error {

	sports := []string{
		"americanfootball_ncaaf",
		"americanfootball_cfl",
		"americanfootball_nfl",
		"americanfootball_nfl_preseason",
		"americanfootball_nfl_super_bowl_winner",

		"baseball_mlb",
		"baseball_mlb_world_series_winner",

		"basketball_nba",
		"basketball_nba_championship_winner",
		"basketball_wnba",

		"boxing_boxing",

		"cricket_asia_cup",
		"cricket_caribbean_premier_league",
		"cricket_international_t20",
		"cricket_odi",
		"cricket_the_hundred",

		"golf_masters_tournament_winner",
		"golf_pga_championship_winner",
		"golf_us_open_winnner",

		"icehockey_nhl_championship_winner",
		"icehockey_sweden_allsvenskan",
		"icehockey_sweden_hockey_league",

		"mma_mixed_martial_arts",

		"rugbyleague_nrl",
	}

	for _, sport := range sports {
		if err := g.MigrateGamesBySport(sport); err != nil {
			return err
		}
	}

	return nil
}

func (g *GameServer) MigrateGamesBySport(sport string) error {

	games, err := g.GetGamesBySport(sport)
	if err != nil {
		return err

	}

	for _, game := range games {

		err := g.AddGameToDB(&game)
		if err != nil {
			return err
		}
	}

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

// # HISTORIC NOT USED
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
