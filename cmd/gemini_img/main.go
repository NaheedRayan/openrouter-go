package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
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


	// Download the image.
    imageResp, err := http.Get("https://upload.wikimedia.org/wikipedia/commons/thumb/8/87/Palace_of_Westminster_from_the_dome_on_Methodist_Central_Hall.jpg/2560px-Palace_of_Westminster_from_the_dome_on_Methodist_Central_Hall.jpg")
    if err != nil {
        panic(err)
    }
    defer imageResp.Body.Close()

    imageBytes, err := io.ReadAll(imageResp.Body)
    if err != nil {
        panic(err)
    }

    // Create the request.
    req := []genai.Part{
        genai.ImageData("jpeg", imageBytes),

        genai.Text("what is this image about?"),
    }


	resp, err := model.GenerateContent(ctx, req...)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(resp.UsageMetadata.TotalTokenCount)
	fmt.Println(resp.Candidates[0].Content.Parts[0])
	// printResponse(resp)
	// [END text_gen_text_only_prompt]
}


