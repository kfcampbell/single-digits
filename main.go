package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
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

	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		return fmt.Errorf("empty token found! %v", token)
	}

	dChannelId := os.Getenv("DISCORD_CHANNEL_ID")
	if dChannelId == "" {
		return fmt.Errorf("empty channel ID found! %v", dChannelId)
	}

	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}
	msgs, err := bot.ChannelMessages(dChannelId, 100, "", "", "")
	if err != nil {
		return err
	}
	client := gosseract.NewClient()
	defer client.Close()

	for _, msg := range msgs {
		//fmt.Printf("message: %v\n", msg)
		if len(msg.Attachments) == 1 {
			img := msg.Attachments[0]

			resp, err := http.Get(img.URL)
			if err != nil {
				return err
			}

			if resp.StatusCode >= 300 {
				return fmt.Errorf("response was unsuccessful: HTTP %v", resp.StatusCode)
			}

			defer resp.Body.Close()

			imgBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			err = client.SetImageFromBytes(imgBytes)
			if err != nil {
				return err
			}
			text, err := client.Text()
			if err != nil {
				return err
			}
			score, err := parser.GetScoreFromText(text)
			if err != nil {
				return err
			}
			fmt.Printf("Author: %v, score: %v\n", msg.Author.Username, score)
		}
	}

	return nil
}
