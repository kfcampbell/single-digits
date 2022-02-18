package utils_test

import (
	"fmt"
	"testing"
	"time"

	sc "github.com/kfcampbell/single-digits/pkg/score"
	"github.com/kfcampbell/single-digits/pkg/utils"
)

// getTestTimes is a convenience func that takes in integers
// and returns time.Duration. It panics if it encounters an error.
// This is gross but done for the sake of simplifying annoying
// test setup logic.
func getTestTimesFromSeconds(times []int) []time.Duration {
	durs := make([]time.Duration, 0, len(times))
	for _, t := range times {
		dur, err := time.ParseDuration(fmt.Sprintf("%vs", t))
		if err != nil {
			panic(fmt.Errorf("could not parse time duration in GetTestTimesFromSeconds: %v", err))
		}
		durs = append(durs, dur)
	}
	return durs
}

// getTestWinners is a convenience func that takes in variadic times and
// returns test scores with dummy names.
func getTestWinners(times []time.Duration) []sc.Score {
	scores := make([]sc.Score, 0, len(times))
	for i, v := range times {
		temp := sc.Score{
			Author:   fmt.Sprintf("TestAuthor%v", i),
			AuthorId: fmt.Sprintf("TestAuthorID%v", i),
			Score:    v,
		}
		scores = append(scores, temp)
	}
	return scores
}

func TestGetWinnersMessage(t *testing.T) {

	simpleThreeScores := make([]int, 0, 3)
	simpleThreeScores = append(simpleThreeScores, 45, 55, 35)

	simpleMoreScores := make([]int, 0, 8)
	simpleMoreScores = append(simpleMoreScores, 32, 68, 92, 13, 48, 36, 29, 42)

	tieForFirstManyScores := make([]int, 0, 4)
	tieForFirstManyScores = append(tieForFirstManyScores, 17, 17, 29, 34)

	cases := []struct {
		info     string
		times    []int
		expected string
	}{
		{"simpleThreeScores", simpleThreeScores, `
		Results for simpleThreeScores:
	🥇 - TestAuthor2 with a time of 35s
	🥈 - TestAuthor0 with a time of 45s
	🥉 - TestAuthor1 with a time of 55s
	`},
		{"simpleMoreScores", simpleMoreScores, `
		Results for simpleMoreScores:
	🥇 - TestAuthor3 with a time of 13s
	🥈 - TestAuthor6 with a time of 29s
	🥉 - TestAuthor0 with a time of 32s
	`},
		{"tieForFirstManyScores", tieForFirstManyScores, `
		Results for tieForFirstManyScores:
		🥇 - tie for first! TestAuthor0 and TestAuthor1 with times of 17s
		🥉 - TestAuthor2 with a time of 29s
		🤏 - TestAuthor3 with a time of 34s
		`},
	}

	for _, tc := range cases {
		actualWinners := getTestWinners(getTestTimesFromSeconds(tc.times))
		actualMessage := utils.GetWinnersMessage(actualWinners, tc.info)
		if tc.expected != actualMessage {
			t.Errorf("%s: expected %s and got %s", tc.info, tc.expected, actualMessage)
		}
	}
}
