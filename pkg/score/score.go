package score

import (
	"sort"
	"time"
)

// Score represents a single instance of a completed crossword puzzle
type Score struct {
	Author   string
	AuthorId string
	Score    time.Duration
}

// SortScores does what it says on the tin, lowest times first
func SortScores(scores []Score) []Score {
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score < scores[j].Score
	})
	return scores
}
