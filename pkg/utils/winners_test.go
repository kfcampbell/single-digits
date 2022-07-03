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

	tieForSecondManyScores := make([]int, 0, 4)
	tieForSecondManyScores = append(tieForSecondManyScores, 15, 22, 22, 49)

	simpleTwoScores := make([]int, 0, 2)
	simpleTwoScores = append(simpleTwoScores, 45,55)

	twoScoresTie := make([]int, 0, 2)
	twoScoresTie = append(twoScoresTie, 30,30)

	simpleOneScore := make([]int, 0, 1)
	simpleOneScore = append(simpleOneScore, 19)

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
		{"tieForSecondManyScores", tieForSecondManyScores, `
		Results for tieForSecondManyScores:
		🥇 - TestAuthor0 with a time of 15s
		🥈 - tie for second! TestAuthor1 and TestAuthor2 with times of 22s
		🤏 - TestAuthor3 with a time of 49s
		`},
		{"simpleTwoScores",simpleTwoScores, `
		Results for simpleTwoScores:
		🥇 - TestAuthor0 with a time of 45s
		🥈 - TestAuthor1 with a time of 55s
		`},
		//whitespace formatting on this case is a lil jank bc the source code has
		//more indentation(should probably be running a trim on these strings but oh well)
		{"twoScoresTie",twoScoresTie, `
			Results for twoScoresTie:
			🥇 - tie for first! TestAuthor0 and TestAuthor1 with times of 30s
			`},
		{"simpleOneScore",simpleOneScore, `
		Results for simpleOneScore:
		🥇 - by default, TestAuthor0 is the winner with a time of 19s
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
