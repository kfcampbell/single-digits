package utils

import (
	"fmt"

	sc "github.com/kfcampbell/single-digits/pkg/score"
)

// GetWinnersMessage takes a list of Score structs, sorts them, and returns
// a nicely formatted podium string.
func GetWinnersMessage(scores []sc.Score, title string) string {
	scores = sc.SortScores(scores)
	// case where only one time is submitted for a day
	if len(scores) == 1 {
		return fmt.Sprintf(`
		Results for %v:
		🥇 - by default, %v is the winner with a time of %v
		`, title,
			scores[0].Author, scores[0].Score)
	}
	// case where only two times are submitted for a day
	if len(scores) == 2 {
		if scores[0].Score == scores[1].Score {
			return fmt.Sprintf(`
			Results for %v:
			🥇 - tie for first! %v and %v with times of %v
			`, title,
				scores[0].Author, scores[1].Author, scores[0].Score)
		}
		return fmt.Sprintf(`
		Results for %v:
		🥇 - %v with a time of %v
		🥈 - %v with a time of %v
		`, title,
			scores[0].Author, scores[0].Score,
			scores[1].Author, scores[1].Score)
	}
	// tie for first place
	if scores[0].Score == scores[1].Score {
		return fmt.Sprintf(`
		Results for %v:
		🥇 - tie for first! %v and %v with times of %v
		🥉 - %v with a time of %v
		🤏 - %v with a time of %v
		`, title,
			scores[0].Author, scores[1].Author, scores[0].Score,
			scores[2].Author, scores[2].Score,
			scores[3].Author, scores[3].Score)
	}

	// tie for second place
	if scores[1].Score == scores[2].Score {
		return fmt.Sprintf(`
		Results for %v:
		🥇 - %v with a time of %v
		🥈 - tie for second! %v and %v with times of %v
		🤏 - %v with a time of %v
		`, title,
			scores[0].Author, scores[0].Score,
			scores[1].Author, scores[2].Author, scores[1].Score,
			scores[3].Author, scores[3].Score)
	}

	// tie for third place
	if len(scores) > 3 && scores[2].Score == scores[3].Score {
		return fmt.Sprintf(`
		Results for %v:
		🥇 - %v with a time of %v
		🥈 - %v with a time of %v
		🥉 - tie for third! %v and %v with times of %v
		`, title,
			scores[0].Author, scores[0].Score,
			scores[1].Author, scores[1].Score,
			scores[2].Author, scores[3].Author, scores[2].Score)
	}
	return fmt.Sprintf(`
		Results for %v:
	🥇 - %v with a time of %v
	🥈 - %v with a time of %v
	🥉 - %v with a time of %v
	`, title,
		scores[0].Author, scores[0].Score,
		scores[1].Author, scores[1].Score,
		scores[2].Author, scores[2].Score)
}
