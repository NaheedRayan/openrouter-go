package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// Load environment variables
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found or couldn't load")
	}
}

func main() {
	LoadEnv()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: OPENAI_API_KEY not found in environment variables")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey))

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a helpful AI assistant."), // System message
			openai.UserMessage("write 1000 word essay"),
		}),
		Model:      openai.F(openai.ChatModelGPT4oMini),
		Temperature: openai.F(0.7), // Adjust creativity (0.0-2.0)
		TopP:        openai.F(0.9), // Controls randomness (0.0-1.0)
		// we dont have topK in the openai-go library
		MaxTokens:   openai.Int(100),  // Max tokens to generate
	})

	if err != nil {
		log.Fatalf("Error getting response: %v", err)
	}

	fmt.Println(chatCompletion.Choices[0].Message.Content)
}
