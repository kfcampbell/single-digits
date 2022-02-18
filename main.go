package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kfcampbell/single-digits/pkg/parser"
	sc "github.com/kfcampbell/single-digits/pkg/score"
	"github.com/kfcampbell/single-digits/pkg/utils"
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

func sortMessages(msgs []*discordgo.Message) []*discordgo.Message {
	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].Timestamp > msgs[j].Timestamp
	})
	return msgs
}

func getMostRecentPuzzleAnnouncements(msgs []*discordgo.Message, botId string) (*discordgo.Message, *discordgo.Message) {
	// get the two most recent instances of the other bot's messages
	botMsgs := make([]*discordgo.Message, 0)
	for _, msg := range msgs {
		// get the latest messages from the other bot
		// make sure they're not GOMLs
		if msg.Author.ID == botId && strings.Contains(msg.Content, "https://www.nytimes.com/crosswords/game/mini") {
			botMsgs = append(botMsgs, msg)
		}
	}
	botMsgs = sortMessages(botMsgs)
	return botMsgs[0], botMsgs[1]
}

func run() error {

	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		return fmt.Errorf("empty token found! %v", token)
	}

	dServerId := os.Getenv("DISCORD_SERVER_ID")
	if dServerId == "" {
		return fmt.Errorf("empty server ID found! %v", dServerId)
	}

	dChannelId := os.Getenv("DISCORD_CHANNEL_ID")
	if dChannelId == "" {
		return fmt.Errorf("empty channel ID found! %v", dChannelId)
	}

	otherBotId := os.Getenv("DISCORD_OTHER_BOT_ID")
	if otherBotId == "" {
		return fmt.Errorf("empty ID for other bot: %v", otherBotId)
	}

	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	msgs, err := bot.ChannelMessages(dChannelId, 100, "", "", "")
	if err != nil {
		return err
	}

	recent, past := getMostRecentPuzzleAnnouncements(msgs, otherBotId)

	ocr := gosseract.NewClient()
	defer ocr.Close()

	scoreMsgs, err := bot.ChannelMessages(dChannelId, 100, recent.ID, past.ID, "")
	if err != nil {
		return err
	}

	scores := make([]sc.Score, 0)
	for _, msg := range scoreMsgs {
		// get rid of these garbage messages
		if msg.Timestamp >= recent.Timestamp || msg.Timestamp <= past.Timestamp {
			continue
		}
		// case where there's an image
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

			err = ocr.SetImageFromBytes(imgBytes)
			if err != nil {
				return err
			}
			text, err := ocr.Text()
			if err != nil {
				return err
			}

			if !parser.ContainsValidScore(text) {
				// if there's no valid score, assume it's a gif about how great
				// @kfcampbell is and move on
				continue
			}

			time, err := parser.GetScoreFromText(text)
			if err != nil {
				return err
			}

			// Discount instantaneous completions (the "Eddie Factor")
			if time == 0 {
				continue
			}

			log.Printf("Author: %v, time: %v\n", msg.Author.Username, time)
			score := &sc.Score{
				Author:   msg.Author.Username,
				AuthorId: msg.Author.ID,
				Score:    time,
			}
			scores = append(scores, *score)
		}
	}

	// subtract two days so it's the Pacific day yesterday instead of the UTC day today
	date := time.Now().Add(-48 * time.Hour).Format("Jan 2, 2006")
	announcement := utils.GetWinnersMessage(scores, date)
	fmt.Println(announcement)

	env := os.Getenv("ENVIRONMENT")
	if env == "PROD" {
		res, err := bot.ChannelMessageSend(dChannelId, announcement)
		if err != nil {
			return err
		}
		log.Printf("message sent: %v", res.Content)
	}

	return nil
}
