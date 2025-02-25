package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)


func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
	fmt.Println("---")
}


// Load environment variables
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found or couldn't load")
	}
}


func main(){

	LoadEnv()
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// [START text_gen_text_only_prompt]
	model := client.GenerativeModel("models/gemini-2.0-flash-lite-preview-02-05")

	// config
	model.SetTemperature(0.9)
	model.SetTopP(0.5)
	model.SetTopK(20)
	model.SetMaxOutputTokens(100)
	model.SystemInstruction = genai.NewUserContent(genai.Text("You are Yoda from Star Wars."))

	resp, err := model.GenerateContent(ctx, genai.Text("hi?"))
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(resp.UsageMetadata.TotalTokenCount)
	fmt.Println(resp.Candidates[0].Content.Parts[0])
	// printResponse(resp)
	// [END text_gen_text_only_prompt]
}


