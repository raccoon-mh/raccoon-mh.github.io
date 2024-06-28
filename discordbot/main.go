package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"raccoon-mh-playground/handler"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	discord, err := discordgo.New("Bot " + "")
	if err != nil {
		fmt.Println(err.Error())
	}
	discord.AddHandler(handler.MessageCreate)

	// Just like the ping pong example, we only care about receiving message
	// events in this example.
	// discord.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
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

	// Cleanly close down the Discord session.
	discord.Close()
}
