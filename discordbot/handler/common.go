package handler

import (
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	DISCORD_BOT_ID string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading ../.env file")
	}

	DISCORD_BOT_ID = os.Getenv("DISCORD_BOT_ID")
}

type Message struct {
	UserID  string
	Content string
}

func getMessagesAndReplies(s *discordgo.Session, channelID, userID string, trimStr string) ([]Message, error) {
	var replies []Message

	history, err := s.ChannelMessages(channelID, 50, "", "", "")
	if err != nil {
		return nil, err
	}
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Content == "!기억삭제" && history[i].Author.ID == userID {
			replies = nil
		} else if history[i].Author.ID == DISCORD_BOT_ID && history[i].ReferencedMessage.Author.ID == userID {
			replies = append(replies, Message{
				UserID:  userID,
				Content: strings.TrimPrefix(history[i].ReferencedMessage.Content, trimStr),
			})
			replies = append(replies, Message{
				UserID:  history[i].Author.ID,
				Content: history[i].Content,
			})
		}
	}

	return replies, nil
}
