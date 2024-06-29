package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"raccoon-mh-playground/handler"

	"github.com/bwmarrin/discordgo"
)

func main() {
	DISCORD_BOT_TOKEN := os.Getenv("DISCORD_BOT_TOKEN")
	discord, err := discordgo.New("Bot " + DISCORD_BOT_TOKEN)
	if err != nil {
		fmt.Println(err.Error())
	}
	discord.AddHandler(handler.PingPong)

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()
}
