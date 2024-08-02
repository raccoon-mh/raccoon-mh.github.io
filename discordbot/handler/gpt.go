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

var (
	gpt_InitPrompt string
)

var (
	GPT_DISCORD_CHANNEL_ID string
	OPENAI_TOKEN           string
)

type GptCompletion struct {
	ID                string      `json:"id"`
	Object            string      `json:"object"`
	Created           int64       `json:"created"`
	Model             string      `json:"model"`
	Choices           []GptChoice `json:"choices"`
	Usage             GptUsage    `json:"usage"`
	SystemFingerprint string      `json:"system_fingerprint"`
}

type GptChoice struct {
	Index        int          `json:"index"`
	GptMessage   GptMessage   `json:"message"`
	Logprobs     *interface{} `json:"logprobs,omitempty"`
	FinishReason string       `json:"finish_reason"`
}

type GptMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GptUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type GptRequest struct {
	Model      string       `json:"model"`
	GptMessage []GptMessage `json:"messages"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading ../.env file")
	}

	GPT_DISCORD_CHANNEL_ID = os.Getenv("GPT_DISCORD_CHANNEL_ID")
	OPENAI_TOKEN = os.Getenv("OPENAI_TOKEN")
	gpt_InitPrompt = os.Getenv("GPT_InitPrompt")
}

func Openaigpt(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !(strings.HasPrefix(m.Content, "!구리 ") && m.ChannelID == GPT_DISCORD_CHANNEL_ID) {
		return
	}

	messages, err := getMessagesAndReplies(s, m.ChannelID, m.Author.ID, "!구리 ")
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, "이전 대화를 읽고있는 중에 문제가 생겼구리..", m.Reference())
		if err != nil {
			log.Println(err)
		}
		return
	}

	prompt := []GptMessage{{
		Role:    "system",
		Content: gpt_InitPrompt,
	}}
	for _, msg := range messages {
		if msg.UserID == m.Author.ID {
			prompt = append(prompt, GptMessage{Role: "user", Content: m.Author.GlobalName + ": " + msg.Content})
		} else if msg.UserID == DISCORD_BOT_ID {
			prompt = append(prompt, GptMessage{Role: "assistant", Content: msg.Content})
		}
	}
	prompt = append(prompt, GptMessage{Role: "user", Content: m.Author.GlobalName + ": " + strings.TrimPrefix(m.Content, "!구리 ")})
	fmt.Println("@@ gpt prompt :", prompt)
	answer, err := OpenaiGptRequest(prompt)
	if err != nil {
		log.Println(err)
		_, err := s.ChannelMessageSendReply(m.ChannelID, "답을 할 수 없어 구리...", m.Reference())
		if err != nil {
			log.Println(err)
		}
	}

	_, err = s.ChannelMessageSendReply(m.ChannelID, answer, m.Reference())
	if err != nil {
		log.Println(err)
	}
}

func OpenaiGptRequest(prompt []GptMessage) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"
	contentType := "application/json"
	authorization := OPENAI_TOKEN

	requestData := &GptRequest{
		Model:      "gpt-4o-mini",
		GptMessage: prompt,
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

	fmt.Println("openai Response Status:", resp.Status)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("통신문제발생")
	}

	fmt.Println("openai Response Body:", string(body))

	var data GptCompletion
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return data.Choices[0].GptMessage.Content, nil
}
