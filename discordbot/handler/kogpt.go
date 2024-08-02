package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type KogptGeneration struct {
	Text   string `json:"text"`
	Tokens int    `json:"tokens"`
}

type KogptUsage struct {
	PromptTokens    int `json:"prompt_tokens"`
	GeneratedTokens int `json:"generated_tokens"`
	TotalTokens     int `json:"total_tokens"`
}

type Kogptresponse struct {
	ID              string            `json:"id"`
	KogptGeneration []KogptGeneration `json:"generations"`
	KogptUsage      KogptUsage        `json:"usage"`
}

var (
	botName  = "너구리"
	userName = "사용자"
	endFix   = "[END]"
)

var (
	InitPrompt  = "시스템:너구리는 친절하고 상냥한 비서이다. 사용자에게 친절하게 답변하며 말끝에 '구리'를 붙여 답한다. 예를 들어, '너구리는 어떤 일을 하니'라는 물음에 대해 너구리는 다음처럼 응답한다. [END] "
	exampleConv = "사용자:자기소개를 부탁해 [END] 너구리:내 임무는 사람들이 더 행복해지도록 돕는 거야구리! [END] 사용자:자기소개를 해봐 [END] 너구리:나는 벙커의 너구리야구리! 너를 행복하게 만들어줄거야구리! [END] "
)

var (
	KOGPT_DISCORD_CHANNEL_ID string
	KAKAO_TOKEN              string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading ../.env file")
	}

	KOGPT_DISCORD_CHANNEL_ID = os.Getenv("KOGPT_DISCORD_CHANNEL_ID")
	KAKAO_TOKEN = os.Getenv("KAKAO_TOKEN")
}

func Kogpt(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !(strings.HasPrefix(m.Content, "!gpt ") && m.ChannelID == KOGPT_DISCORD_CHANNEL_ID) {
		return
	}

	prompt := InitPrompt + exampleConv
	messages, err := getMessagesAndReplies(s, m.ChannelID, m.Author.ID, "!gpt ")
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, "문제가 생겼구리..", m.Reference())
		if err != nil {
			log.Println(err)
		}
		return
	}

	for _, msg := range messages {
		if msg.UserID == m.Author.ID {
			prompt += fmt.Sprintf(" %s: %s %s ", userName, msg.Content, endFix)
		} else if msg.UserID == DISCORD_BOT_ID {
			prompt += fmt.Sprintf(" %s: %s %s ", botName, msg.Content, endFix)
		}
	}

	prompt += "시스템: 대화를 계속해라. [END] 사용자: " + strings.TrimPrefix(m.Content, "!gpt ") + " [END] 너구리:"
	fmt.Println("### kogpt prompt : ", prompt)
	answer, err := KogptRequest(prompt, 0)
	if err != nil {
		log.Println(err)
		_, err := s.ChannelMessageSendReply(m.ChannelID, "답을 할 수 없어 구리...", m.Reference())
		if err != nil {
			log.Println(err)
		}
	}

	msgResponse := getDataBeforeUser(answer)
	_, err = s.ChannelMessageSendReply(m.ChannelID, msgResponse, m.Reference())
	if err != nil {
		log.Println(err)
	}
}

func getDataBeforeUser(data string) string {
	parts := strings.Split(data, endFix)
	if len(parts) > 1 {
		return parts[0]
	}
	return data
}

func KogptRequest(prompt string, level int) (string, error) {
	url := "https://api.kakaobrain.com/v1/inference/kogpt/generation"
	contentType := "application/json"
	authorization := KAKAO_TOKEN
	max_tokens := 50

	requestData := map[string]interface{}{
		"prompt":      prompt,
		"max_tokens":  max_tokens,
		"temperature": 0.05,
		"top_p":       0.3,
	}

	requestDataBytes, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error marshaling request data:", err)
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestDataBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", authorization)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}

	fmt.Println("kogpt Response Status:", resp.Status)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("통신문제발생")
	}

	fmt.Println("kogpt Response Body:", string(body))

	var data Kogptresponse
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if !strings.Contains(data.KogptGeneration[0].Text, "[END]") {
		if level+1 > 3 {
			return "", fmt.Errorf("depth error")
		}
		fmt.Printf("### level : %d / prompt : %s", level, prompt+data.KogptGeneration[0].Text)
		result, err := KogptRequest(prompt+data.KogptGeneration[0].Text, level+1)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		fmt.Println("@@@@@@@@@@@@@@@@@@@@@", data.KogptGeneration[0].Text+result)
		return data.KogptGeneration[0].Text + result, nil
	} else if data.KogptUsage.GeneratedTokens < max_tokens {
		return data.KogptGeneration[0].Text, nil
	}

	return data.KogptGeneration[0].Text, nil
}
