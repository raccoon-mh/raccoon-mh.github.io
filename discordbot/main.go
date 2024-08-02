package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"raccoon-mh-playground/command"
	"raccoon-mh-playground/handler"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	DISCORD_BOT_TOKEN string
	DISCORD_SERVER    string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	DISCORD_BOT_TOKEN = os.Getenv("DISCORD_BOT_TOKEN")
	DISCORD_SERVER = os.Getenv("DISCORD_SERVER")
}

func main() {
	fmt.Println("Start main stream")

	DISCORD_BOT_TOKEN := os.Getenv("DISCORD_BOT_TOKEN")
	discord, err := discordgo.New("Bot " + DISCORD_BOT_TOKEN)
	if err != nil {
		fmt.Println(err.Error())
	}

	discord.AddHandler(handler.PingPong)
	discord.AddHandler(handler.Kogpt)
	discord.AddHandler(handler.Openaigpt)
	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := command.CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Start Register Commands")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(command.Commands))
	for i, v := range command.Commands {
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, DISCORD_SERVER, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	fmt.Println("Done Register Commands")

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()

	log.Println("Removing commands...")
	for _, v := range registeredCommands {
		err := discord.ApplicationCommandDelete(discord.State.User.ID, DISCORD_SERVER, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
	log.Println("Done!")
}
