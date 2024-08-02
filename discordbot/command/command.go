package command

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "날씨",
			Description: "초단기 예보를 해준다구리! 1시간내 아주 정확한 정보야구리",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "지역",
					Description: "주소를 입력해구리",
					Required:    true,
				},
			},
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"날씨": weatherCommand,
	}
)

var (
	KISANG_TOKEN string
	KAKAO_TOKEN  string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading ../.env file")
	}

	KAKAO_TOKEN = os.Getenv("KAKAO_TOKEN")
	KISANG_TOKEN = os.Getenv("KISANG_TOKEN")
}
