package parser

import (
	"strings"
	"time"
)

// GetScoreFromText extracts the time from the boilerplate NYT text
// and returns it in seconds.
func GetScoreFromText(text string) (time.Duration, error) {
	const dividerText = "You solved a Mini puzzle in "
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
	}
}
