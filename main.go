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

// GetText returns the text from the image.
func GetText(imgPath string) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(imgPath)
	text, err := client.Text()
	if err != nil {
		return "", err
	}
	return text, nil
}

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

func run() error {
	client := gosseract.NewClient()
	defer client.Close()

	text, err := GetText("testdata/one-minute-forty-two-seconds.jpeg")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", text)
	score, err := GetScoreFromText(text)
	if err != nil {
		return err
	}
	fmt.Printf("score: %v\n", score)

	text, err = GetText("testdata/forty-two-seconds.jpeg")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", text)
	score, err = GetScoreFromText(text)
	if err != nil {
		return err
	}
	fmt.Printf("score: %v\n", score)

	text, err = GetText("testdata/fifty-seven-seconds-plain.jpeg")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", text)
	/*score, err = GetScoreFromText(text)
	if err != nil {
		return err
	}
	fmt.Printf("score: %v\n", score)*/

	text, err = GetText("testdata/thirty-six-seconds.jpeg")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", text)
	/*score, err = GetScoreFromText(text)
	if err != nil {
		return err
	}
	fmt.Printf("score: %v\n", score)*/

	text, err = GetText("testdata/twenty-three-seconds.jpeg")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", text)
	/*score, err = GetScoreFromText(text)
	if err != nil {
		return err
	}
	fmt.Printf("score: %v\n", score)*/

	return nil
}
