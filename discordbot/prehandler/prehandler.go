package prehandler

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	DISCORD_CHANNEL string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading ../.env file")
	}

	DISCORD_CHANNEL = os.Getenv("DISCORD_CHANNEL")
}

func PreHandler(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	logUserMessage(m)

	if isBotMessage(s, m) {
		return true
	}

	if !isBotChannel(m) {
		log.Printf("isBotChannel:false [%s]", m.ChannelID)
		return true
	}

	return false
}

func PreHandlerByCannelId(s *discordgo.Session, m *discordgo.MessageCreate, channelId string) bool {
	logUserMessage(m)

	if isBotMessage(s, m) {
		return true
	}

	if !isChannel(m, channelId) {
		log.Printf("isChannel:false [%s]", m.ChannelID)
		return true
	}

	return false
}

func logUserMessage(m *discordgo.MessageCreate) {
	log.Printf("user:[%s](%s)-ch:(%s)-msg:[%s]", m.Author.Username, m.Author.ID, m.ChannelID, m.Content)
}

func isBotMessage(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	return m.Author.ID == s.State.User.ID
}

func isBotChannel(m *discordgo.MessageCreate) bool {
	return m.ChannelID == DISCORD_CHANNEL
}

func isChannel(m *discordgo.MessageCreate, channelId string) bool {
	return m.ChannelID == channelId
}
