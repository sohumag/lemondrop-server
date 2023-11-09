package games

import (
	"fmt"

	"slices"
)

func ValidateLeagueExists(league string) error {
	// if invalid league sent
	if !slices.Contains(validLeagues, league) {
		return fmt.Errorf("Invalid league")
	}

	return nil
}
