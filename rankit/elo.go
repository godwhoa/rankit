package rankit

import (
	"math"
)

const (
	DEFAULT_ELO_RATING = 1000
	kFactor            = 32
)

// CalculateElo calculates the new elo ratings for the winner and loser
func CalculateElo(winnerRating, loserRating int) (int, int) {
	ea := 1 / (1 + math.Pow(10, float64(loserRating-winnerRating)/400))

	winnerRating += int(math.Ceil(kFactor * (1 - ea)))
	loserRating += int(math.Ceil(kFactor * (0 - (1 - ea))))

	return winnerRating, loserRating
}
