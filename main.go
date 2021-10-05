package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/otiai10/gosseract/v2"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

// getScoreFromText extracts the time from the boilerplate NYT text
// and returns it in seconds.
func getScoreFromText(text string) (time.Duration, error) {
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

func run() error {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage("testdata/one-minute-forty-two-seconds.jpeg")
	text, err := client.Text()
	if err != nil {
		return err
	}
	score, err := getScoreFromText(text)
	if err != nil {
		return err
	}
	fmt.Printf("score: %v\n", score)

	client.SetImage("testdata/forty-two-seconds.jpeg")
	text, err = client.Text()
	if err != nil {
		return err
	}
	score, err = getScoreFromText(text)
	if err != nil {
		return err
	}
	fmt.Printf("score: %v\n", score)

	return nil
}
