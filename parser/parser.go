package parser

import (
	"fmt"
	"strings"
	"time"
)

// todo(kfcampbell): make parsing logic actually work with test cases

// GetScoreFromText extracts the time from the boilerplate NYT text
// and returns it in seconds.
func GetScoreFromText(text string) (time.Duration, error) {

	// strip off the colon in the header time if it exists
	if strings.Contains(text, "puzzle in") {
		text = text[strings.Index(text, "puzzle in"):]
	}

	fmt.Printf("trimmed text: %v", text)

	const divider = ":"
	s := strings.Split(text, divider)
	if len(s) == 1 {
		// case: formatted like "You solved a mini puzzle in 35 seconds."

	}

	return 0, fmt.Errorf("could not parse time correctly: not implemented yet")

	/*const dividerText = "You solved a Mini puzzle in "
	s := strings.Split(text, dividerText)
	scoreText := s[1]
	// case when it's formatted like "35 seconds."
	if strings.Contains(scoreText, "seconds") {
		timePieces := strings.Split(scoreText, " ")
		secondsS := timePieces[0]
		return time.ParseDuration(secondsS + "s")
	} else { // case when it's formatted like "1:42"
		timePieces := strings.Split(scoreText, ":")
		minutesS := timePieces[0]
		secondsPieces := timePieces[1]
		secondsS := strings.Split(secondsPieces, ".")

		minutes, err := time.ParseDuration(minutesS + "m")
		if err != nil {
			return 0, err
		}

		seconds, err := time.ParseDuration(secondsS[0] + "s")
		if err != nil {
			return 0, err
		}

		return minutes + seconds, nil
	}*/
}
