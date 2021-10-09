package parser

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// GetScoreFromText extracts the time from the boilerplate NYT text
// and returns it in seconds.
func GetScoreFromText(text string) (time.Duration, error) {

	// strip off the colon in the header time if it exists
	if strings.Contains(text, "puzzle in") {
		text = text[strings.Index(text, "puzzle in"):]
	}

	fmt.Printf("trimmed text: %v", text)

	// regex to get digits: ([0-9]+)
	r := regexp.MustCompile("[0-9]+")

	const divider = ":"
	s := strings.Split(text, divider)
	if len(s) == 1 {
		// case: formatted like "You solved a mini puzzle in 35 seconds."
		match := r.FindString(text)
		fmt.Printf("Time: %v", match)

		seconds, err := time.ParseDuration(match + "s")
		if err != nil {
			return 0, fmt.Errorf("could not parse time: %v", err)
		}
		return seconds, nil
	}

	// expect there to be two time pieces
	pieces := strings.Split(text, ":")
	minutesPiece := r.FindString(pieces[0])
	minutes, err := time.ParseDuration(minutesPiece + "m")
	if err != nil {
		return 0, fmt.Errorf("could not parse time: %v", err)
	}

	secondsPiece := r.FindString(pieces[1])
	seconds, err := time.ParseDuration(secondsPiece + "s")
	if err != nil {
		return 0, fmt.Errorf("could not parse time: %v", err)
	}

	return minutes + seconds, nil
}
