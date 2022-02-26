package main

import (
	"log"
	"os"

	"github.com/otiai10/gosseract/v2"
)

func main() {
	// First argument is always the running program's filename, and the second argument will always be
	// `--` to signify the end of arguments to the `go` command.
	imagePath := os.Args[2]

	ocr := gosseract.NewClient()
	defer ocr.Close()

	err := ocr.SetImage(imagePath)
	if err != nil {
		log.Fatal(err)
	}

	text, err := ocr.Text()
	if err != nil {
		log.Fatal(err)
	}

	log.Print(text)
}
