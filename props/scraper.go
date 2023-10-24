package props

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Scraping Bovada

func (s *PropServer) ScrapeAllBovadaSports() error {
	//! this is only for future events
	//! need to implement futures and live
	// allSports := []string{
	// 	"football/nfl",
	// 	"football/college-football",

	// 	"baseball/mlb",

	// 	"hockey/nhl",

	// 	"basketball/college-basketball",
	// 	"basketball/nba",

	// 	// soccer
	// 	// tennis
	// }

	sport := "basketball/nba"
	s.ScrapeBovadaSport(sport)

	return nil
}

func (s *PropServer) ScrapeBovadaSport(sportEndpoint string) error {
	endpointUrl := fmt.Sprintf("https://bovada.lv/sports/%s", sportEndpoint)

	res, err := http.Get(endpointUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)

	os.WriteFile("./nba.html", bytes, 0644)

	return err
}
