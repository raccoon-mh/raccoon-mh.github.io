package handler

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	DISCORD_CHANNEL = os.Getenv("DISCORD_CHANNEL")
)

func PingPong(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Printf("[%s] (%s) message: %s", m.Author.Username, m.Author.ID, m.Content)
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.ChannelID != DISCORD_CHANNEL {
		fmt.Println(m.ChannelID, DISCORD_CHANNEL)
		return
	}

	if m.Content == "ping" {
		s.ChannelMessageSendReply(m.ChannelID, "Pong!", m.MessageReference)
	}
	if m.Content == "pong" {
		s.ChannelMessageSendReply(m.ChannelID, "Ping!", m.MessageReference)
	}
}
