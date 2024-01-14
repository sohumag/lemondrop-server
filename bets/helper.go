package bets

import (
	"strconv"

	"github.com/rlvgl/bookie-server/users"
)

// hasEnoughBalance checks if the user has enough balance for the bet.
func hasEnoughBalance(amt float64, user users.User) bool {
	return amt <= user.CurrentBalance
}

func stringToFloat64(input string) (float64, error) {
	result, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func overlapCount(s1, s2 string) int {
	set := make(map[rune]struct{})

	for _, char := range s1 {
		set[char] = struct{}{}
	}

	count := 0
	for _, char := range s2 {
		if _, ok := set[char]; ok {
			count++
		}
	}

	return count
}
