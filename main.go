package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	tokenizer "github.com/samber/go-gpt-3-encoder"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
)

const (
	MAX_TOKEN_LIMIT    = 4096
	MAX_RESPONSE_TOKEN = 250
)

func main() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	apiKey := viper.GetString("API_KEY")
	if apiKey == "" {
		log.Fatalln("API_KEY not found in config file or command-line arguments")
	}

	if err := runChat(); err != nil {
		log.Fatalf("Error executing root command: %s", err)
	}
}

// runChat is the main loop for the chatbot
func runChat() error {
	client := openai.NewClient(viper.GetString("API_KEY"))
	scanner := bufio.NewScanner(os.Stdin)

	history := chatHistory{}
	fmt.Print("pick an expert: ")
	scanner.Scan()
	history.expert = strings.TrimSpace(scanner.Text())
	if history.expert == "" {
		log.Fatalln("Need to pick an expert")
	}
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	history.addMessage(openai.ChatMessageRoleSystem, fmt.Sprintf("You are a %s expert", history.expert))

	for {
		fmt.Println()
		green.Printf("%s expert here ask your question('stop' to end): ", history.expert)
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "stop" {
			break
		}

		history.addMessage(openai.ChatMessageRoleUser, input)
		history.maintainTokenLimit()

		resp, err := getResponse(client, context.Background(), history.messages)
		if err != nil {
			log.Printf("Error getting response: %s", err)
			continue
		}

		blue.Printf("%s expert: %s\n", history.expert, resp)
		//fmt.Printf("%s expert: %s\n", history.expert, resp)
		fmt.Println()

		history.addMessage(openai.ChatMessageRoleAssistant, resp)
	}

	return nil
}

// getResponse sends the conversation history to the OpenAI API to generate a response
func getResponse(client *openai.Client, ctx context.Context, messages []openai.ChatCompletionMessage) (string, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT3Dot5Turbo,
			Messages:  messages,
			MaxTokens: MAX_RESPONSE_TOKEN,
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

// chatHistory keeps track of the conversation history
type chatHistory struct {
	messages []openai.ChatCompletionMessage
	expert   string
}

// addMessage adds a message to the conversation history
func (h *chatHistory) addMessage(role string, content string) {
	h.messages = append(h.messages, openai.ChatCompletionMessage{
		Role:    role,
		Content: content,
	})
}

// maintainTokenLimit maintains token limit for conversation history
func (h *chatHistory) maintainTokenLimit() {
	out, err := json.Marshal(h.messages)
	if err != nil {
		log.Fatalln(err)
	}
	encoder, err := tokenizer.NewEncoder()
	if err != nil {
		log.Fatalln(err)
	}
	encoded, err := encoder.Encode(string(out))
	if err != nil {
		log.Fatalln(err)
	}

	if len(encoded)+MAX_RESPONSE_TOKEN >= MAX_TOKEN_LIMIT {
		h.removeOldMessages()
	}
}

// removeOldMessages removes older messages from chat history
func (h *chatHistory) removeOldMessages() {
	half := len(h.messages) / 2
	newMessages := make([]openai.ChatCompletionMessage, half+1)
	newMessages[0] = h.messages[0]
	copy(newMessages[1:], h.messages[half:])
	h.messages = newMessages
}
