package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	apiKey := viper.GetString("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY not found in config file or command-line arguments")
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

	history.addMessage(openai.ChatMessageRoleSystem, fmt.Sprintf("You are a %s expert", history.expert))

	for {
		fmt.Println()
		fmt.Printf("%s expert here ask your question('stop' to end): ", history.expert)
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

		resp, err := getResponse(client, context.Background(), history.messages)
		if err != nil {
			log.Printf("Error getting response: %s", err)
			continue
		}

		fmt.Printf("%s expert: %s\n", history.expert, resp)
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
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
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
