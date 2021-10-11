package main

import (
	"fmt"
	"log"

	"github.com/kfcampbell/single-digits/parser"
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

func run() error {
	client := gosseract.NewClient()
	defer client.Close()

	text, err := GetText("testdata/one-minute-forty-two-seconds.jpeg")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", text)
	score, err := parser.GetScoreFromText(text)
	if err != nil {
		return err
	}
	fmt.Printf("score: %v\n", score)

	text, err = GetText("testdata/forty-two-seconds.jpeg")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", text)
	score, err = parser.GetScoreFromText(text)
	if err != nil {
		return err
	}
	fmt.Printf("score: %v\n", score)

	text, err = GetText("testdata/fifty-seven-seconds-plain.jpeg")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", text)
	score, err = parser.GetScoreFromText(text)
	if err != nil {
		return err
	}
	fmt.Printf("score: %v\n", score)

	text, err = GetText("testdata/thirty-six-seconds.jpeg")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", text)
	score, err = parser.GetScoreFromText(text)
	if err != nil {
		return err
	}
	fmt.Printf("score: %v\n", score)

	text, err = GetText("testdata/twenty-three-seconds.jpeg")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", text)
	score, err = parser.GetScoreFromText(text)
	if err != nil {
		return err
	}
	fmt.Printf("score: %v\n", score)

	return nil
}
