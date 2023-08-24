package games

func (g *GameServer) MigrateSports() error {
	// here to ensure not duplicating sports in db
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
	return nil // protection against expensive API call

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
