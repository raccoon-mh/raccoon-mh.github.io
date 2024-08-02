package handler

import (
	"log"
	"raccoon-mh-playground/prehandler"

	"github.com/bwmarrin/discordgo"
)

func PingPong(s *discordgo.Session, m *discordgo.MessageCreate) {
	if prehandler.PreHandler(s, m) {
		return
	}

	if m.Content == "ping" {
		_, err := s.ChannelMessageSendReply(m.ChannelID, "Pong!", m.Reference())
		if err != nil {
			log.Println(err)
		}
	}
	if m.Content == "pong" {
		_, err := s.ChannelMessageSendReply(m.ChannelID, "Ping!", m.Reference())
		if err != nil {
			log.Println(err)
		}
	}
}
