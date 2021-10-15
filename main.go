package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

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
		if msg.Author.ID == botId {
			botMsgs = append(botMsgs, msg)
		}
	}
	botMsgs = sortMessages(botMsgs)
	return botMsgs[0], botMsgs[1]
}

// Score represents a single instance of the crossword puzzle
type Score struct {
	Author   string
	AuthorId string
	Score    time.Duration
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

	scores := make([]Score, 0)
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
			log.Printf("Author: %v, time: %v\n", msg.Author.Username, time)
			score := &Score{
				Author:   msg.Author.Username,
				AuthorId: msg.Author.ID,
				Score:    time,
			}
			scores = append(scores, *score)
		}
	}

	scores = sortScores(scores)

	// subtract a day so it's the Pacific day instead of the UTC one
	date := time.Now().Add(-24 * time.Hour).Format("Jan 2, 2006")
	announcement := fmt.Sprintf(`
		Results for %v:
	🥇 - %v with a time of %v
	🥈 - %v with a time of %v
	🥉 - %v with a time of %v
	`, date,
		scores[0].Author, scores[0].Score,
		scores[1].Author, scores[1].Score,
		scores[2].Author, scores[2].Score)

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

func sortScores(scores []Score) []Score {
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score < scores[j].Score
	})
	return scores
}
