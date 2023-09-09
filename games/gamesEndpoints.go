package games

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func (g *GameServer) SendGameById(c *fiber.Ctx, id string) error {
	coll := g.client.Database("games-db").Collection("games")
	var game Game
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "gameid", Value: id}}).Decode(&game); err != nil {
		fmt.Println(err.Error())
		return err
	}

	c.JSON(ParseGames([]Game{game})[0])

	return nil
}

func (g *GameServer) SendGamesBySport(c *fiber.Ctx, sport string) error {
	games, err := g.GetUpcomingGamesBySport(sport)
	if err != nil {
		return err
	}

	c.JSON(ParseGames(games))
	return nil
}

func (g *GameServer) SendGamesUpcoming(c *fiber.Ctx) error {
	games, err := g.GetAllUpcomingGames(50)
	if err != nil {
		return err
	}

	c.JSON(ParseGames(games))

	return nil
}

type ClientGame struct {
	Home         string
	Away         string
	CommenceTime time.Time
	GameId       string
	LastUpdate   time.Time
	Scores       []ScoreStr
	SportKey     string
	SportTitle   string
	Completed    bool

	MoneylineHome   int
	MoneylineAway   int
	MoneylineExists bool

	SpreadHomePoint float64
	SpreadHomePrice int
	SpreadAwayPoint float64
	SpreadAwayPrice int
	SpreadExists    bool

	OverPoint    float64
	OverPrice    int
	UnderPoint   float64
	UnderPrice   int
	TotalsExists bool
}

func ParseGames(rawGames []Game) []ClientGame {

	newGames := []ClientGame{}

	for _, game := range rawGames {
		homeTeam := game.HomeTeam
		awayTeam := game.AwayTeam

		marketPrices := map[string][]OutcomeDate{
			"h2h":     {},
			"spreads": {},
			"totals":  {},
		}

		for _, bookmaker := range game.Bookmakers {
			markets := bookmaker.Markets

			for _, market := range markets {
				switch market.Key {
				case "h2h":
					marketPrices["h2h"] = append(marketPrices["h2h"], OutcomeDate{Outcomes: market.Outcomes, Date: market.LastUpdate})
				case "spreads":
					marketPrices["spreads"] = append(marketPrices["spreads"], OutcomeDate{Outcomes: market.Outcomes, Date: market.LastUpdate})
				case "totals":
					marketPrices["totals"] = append(marketPrices["totals"], OutcomeDate{Outcomes: market.Outcomes, Date: market.LastUpdate})
				}

			}
		}

		// returns (outcomes, exists)
		getMostRecentlyUpdatedOutcomes := func(outcomeDates []OutcomeDate) (Outcome, Outcome, bool) {

			if len(outcomeDates) == 0 {
				return Outcome{}, Outcome{}, false
			}

			latest := time.Date(2020, time.April,
				11, 21, 34, 01, 0, time.UTC)
			outcomes := []Outcome{}

			for _, m := range outcomeDates {
				if m.Date.Second() > latest.Second() {
					latest = m.Date
					outcomes = m.Outcomes
				}
			}

			if outcomes[0].Name == homeTeam || outcomes[0].Name == "Over" {
				return outcomes[0], outcomes[1], true
			} else {
				return outcomes[1], outcomes[0], true
			}
		}

		h2hHome, h2hAway, h2hExists := getMostRecentlyUpdatedOutcomes(marketPrices["h2h"])
		spreadHome, spreadAway, spreadExists := getMostRecentlyUpdatedOutcomes(marketPrices["spreads"])
		over, under, totalsExists := getMostRecentlyUpdatedOutcomes(marketPrices["totals"])

		newGame := ClientGame{
			Home:         homeTeam,
			Away:         awayTeam,
			CommenceTime: game.CommenceTime,
			GameId:       game.GameId,
			LastUpdate:   game.LastUpdate,
			Scores:       game.Scores,
			SportKey:     game.SportKey,
			SportTitle:   game.SportTitle,
			Completed:    game.Completed,
		}

		if !h2hExists {
			newGame.MoneylineExists = false
		} else {
			newGame.MoneylineExists = true
			newGame.MoneylineHome = int(h2hHome.Price)
			newGame.MoneylineAway = int(h2hAway.Price)
		}

		if !spreadExists {
			newGame.SpreadExists = false
		} else {
			newGame.SpreadExists = true
			newGame.SpreadHomePoint = spreadHome.Point
			newGame.SpreadHomePrice = int(spreadHome.Price)
			newGame.SpreadAwayPoint = spreadAway.Point
			newGame.SpreadAwayPrice = int(spreadAway.Price)
		}

		if !totalsExists {
			newGame.TotalsExists = false
		} else {
			newGame.TotalsExists = true
			newGame.OverPoint = over.Point
			newGame.OverPrice = int(over.Price)
			newGame.UnderPoint = under.Point
			newGame.UnderPrice = int(under.Price)
		}

		newGames = append(newGames, newGame)
	}

	return newGames
}

type OutcomeDate struct {
	Outcomes []Outcome
	Date     time.Time
}
