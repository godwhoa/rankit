package rankit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEloRating(t *testing.T) {
	testCases := []struct {
		winningRating        int
		losingRating         int
		expectedWinnerRating int
		expectedLoserRating  int
	}{
		{1200, 1000, 1208, 993},
		{1500, 1500, 1516, 1484},
		{2000, 1800, 2008, 1793},
		{1200, 1200, 1216, 1184},
	}

	for _, tc := range testCases {
		newWinningRating, newLosingRating := CalculateElo(tc.winningRating, tc.losingRating)
		assert.Equal(t, tc.expectedWinnerRating, newWinningRating, "Winning ratings should match")
		assert.Equal(t, tc.expectedLoserRating, newLosingRating, "Losing ratings should match")
	}
}
