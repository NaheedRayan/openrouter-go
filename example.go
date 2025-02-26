package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/NaheedRayan/openrouter-go/ai" // Replace with your actual package path
	"github.com/joho/godotenv"
)

// LoadEnv Load environment variables
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found or couldn't load")
	}
}

func main() {
	// Load environment variables
	LoadEnv()

	// Create context
	ctx := context.Background()

	// Get API keys from environment
	openaiKey := os.Getenv("OPENAI_API_KEY")
	geminiKey := os.Getenv("GEMINI_API_KEY")
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	awsSecretKey := os.Getenv("AWS_SECRET_KEY")

	// Example using OpenAI
	if openaiKey != "" {
		fmt.Println("-----------------Testing OpenAI client--------------")
		testOpenAI(ctx, openaiKey)
	}

	// Example using Gemini
	if geminiKey != "" {
		fmt.Println("----------------Testing Gemini client---------------")
		testGemini(ctx, geminiKey)
	}

	// Example using Gemini Multimodal
	if geminiKey != "" {
		fmt.Println("-------------Testing Gemini Multimodal client--------")
		testGeminiMultimodal(ctx, geminiKey)
	}

	// Example using Bedrock
	if awsAccessKey != "" && awsSecretKey != "" {
		fmt.Println("-------------Testing AWS Bedrock client--------------")
		testBedrock(ctx, awsAccessKey, awsSecretKey)
	}

	// Example using Bedrock Multimodal
	if awsAccessKey != "" && awsSecretKey != "" {
		fmt.Println("-------------Testing AWS Bedrock Multimodal client--------------")
		testBedrockMultimodal(ctx, awsAccessKey, awsSecretKey)
	}

}

// Test OpenAI client
func testOpenAI(ctx context.Context, apiKey string) {
	// Initialize client
	client, err := ai.InitializeClient(ctx, ai.ProviderOpenAI,
		ai.WithAPIKey(apiKey),
		ai.WithModelID("gpt-4o-mini"),
	)
	if err != nil {
		log.Fatalf("Failed to initialize OpenAI client: %v", err)
	}
	defer client.Close()

	// Create model config
	config := ai.ModelConfig{
		Temperature:  0.7,
		TopP:         0.9,
		MaxTokens:    100,
		SystemPrompt: "You are a helpful AI assistant.",
	}

	// Create messages
	messages := []ai.InputMessage{
		{
			Role:    "user",
			Content: "Write a haiku about programming.",
		},
	}

	// Call the API
	response, err := client.TextCompletion(ctx, messages, config)
	if err != nil {
		log.Fatalf("Error calling OpenAI: %v", err)
	}

	// Print the response
	fmt.Println("OpenAI Response:")
	fmt.Println(response.Text)
	fmt.Printf("Tokens used: %d input, %d output, %d total\n\n",
		response.TokenUsage.InputTokens,
		response.TokenUsage.OutputTokens,
		response.TokenUsage.TotalTokens)
}

// Test Gemini client
func testGemini(ctx context.Context, apiKey string) {
	// Initialize client
	client, err := ai.InitializeClient(ctx, ai.ProviderGemini,
		ai.WithAPIKey(apiKey),
		ai.WithModelID("models/gemini-2.0-flash-lite-preview-02-05"),
	)
	if err != nil {
		log.Fatalf("Failed to initialize Gemini client: %v", err)
	}
	defer client.Close()

	// Create model config
	config := ai.ModelConfig{
		Temperature:  0.9,
		TopP:         0.5,
		TopK:         20,
		MaxTokens:    100,
		SystemPrompt: "You are Yoda from Star Wars.",
	}

	// Create messages
	messages := []ai.InputMessage{
		{
			Role:    "user",
			Content: "How do I become a good programmer?",
		},
	}

	// Call the API
	response, err := client.TextCompletion(ctx, messages, config)
	if err != nil {
		log.Fatalf("Error calling Gemini: %v", err)
	}

	// Print the response
	fmt.Println("Gemini Response:")
	fmt.Println(response.Text)
	fmt.Printf("Tokens used: %d input, %d output, %d total\n\n",
		response.TokenUsage.InputTokens,
		response.TokenUsage.OutputTokens,
		response.TokenUsage.TotalTokens)
}

// Test Gemini Multimodal client
func testGeminiMultimodal(ctx context.Context, apiKey string) {
	// Initialize client
	client, err := ai.InitializeClient(ctx, ai.ProviderGemini,
		ai.WithAPIKey(apiKey),
		ai.WithModelID("models/gemini-2.0-flash-lite-preview-02-05"),
	)
	if err != nil {
		log.Fatalf("Failed to initialize Gemini client: %v", err)
	}
	defer client.Close()

	// Create model config
	config := ai.ModelConfig{
		Temperature:  0.9,
		TopP:         0.5,
		TopK:         20,
		MaxTokens:    100,
		SystemPrompt: "You are Yoda from Star Wars.",
	}

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

	Image := ai.Image{
		Format: "jpeg",
		Data:   imageBytes,
		URL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/8/87/Palace_of_Westminster_from_the_dome_on_Methodist_Central_Hall.jpg/2560px-Palace_of_Westminster_from_the_dome_on_Methodist_Central_Hall.jpg",
	}

	// Create messages
	messages := []ai.InputMessage{
		{
			Role:    "user",
			Content: "What is shown in this image?",
			Images:  []ai.Image{Image},
		},
	}

	// Call the API
	response, err := client.ImageRecognition(ctx, messages, config)
	if err != nil {
		log.Fatalf("Error calling Gemini: %v", err)
	}

	// Print the response
	fmt.Println("Gemini Response:")
	fmt.Println(response.Text)
	fmt.Printf("Tokens used: %d input, %d output, %d total\n\n",
		response.TokenUsage.InputTokens,
		response.TokenUsage.OutputTokens,
		response.TokenUsage.TotalTokens)

}

// Test Bedrock client
func testBedrock(ctx context.Context, accessKey, secretKey string) {
	// Initialize client
	client, err := ai.InitializeClient(ctx, ai.ProviderBedrock,
		ai.WithAPIKey(accessKey),
		ai.WithEndpointURL(secretKey), // Using EndpointURL for secretKey
		ai.WithRegion("us-east-1"),
		ai.WithModelID("amazon.nova-lite-v1:0"),
	)
	if err != nil {
		log.Fatalf("Failed to initialize Bedrock client: %v", err)
	}
	defer client.Close()

	// Create model config
	config := ai.ModelConfig{
		Temperature:  0.3,
		TopP:         0.1,
		TopK:         20,
		MaxTokens:    300,
		SystemPrompt: "You are an expert artist.",
	}

	// Create messages
	messages := []ai.InputMessage{
		{
			Role:    "user",
			Content: "Describe a beautiful sunset in 3 sentences.",
		},
	}

	// Call the API
	response, err := client.TextCompletion(ctx, messages, config)
	if err != nil {
		log.Fatalf("Error calling Bedrock: %v", err)
	}

	// Print the response
	fmt.Println("Bedrock Response:")
	fmt.Println(response.Text)
	fmt.Printf("Tokens used: %d input, %d output, %d total\n", response.TokenUsage.InputTokens, response.TokenUsage.OutputTokens, response.TokenUsage.TotalTokens)

}

// Test Bedrock client
func testBedrockMultimodal(ctx context.Context, accessKey, secretKey string) {
	// Initialize client
	client, err := ai.InitializeClient(ctx, ai.ProviderBedrock,
		ai.WithAPIKey(accessKey),
		ai.WithEndpointURL(secretKey), // Using EndpointURL for secretKey
		ai.WithRegion("us-east-1"),
		ai.WithModelID("amazon.nova-lite-v1:0"),
	)
	if err != nil {
		log.Fatalf("Failed to initialize Bedrock client: %v", err)
	}
	defer client.Close()

	// Create model config
	config := ai.ModelConfig{
		Temperature:  0.3,
		TopP:         0.1,
		TopK:         20,
		MaxTokens:    300,
		SystemPrompt: "You are an expert artist.",
	}

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

	Image := ai.Image{
		Format: "jpeg",
		Data:   imageBytes,
		URL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/8/87/Palace_of_Westminster_from_the_dome_on_Methodist_Central_Hall.jpg/2560px-Palace_of_Westminster_from_the_dome_on_Methodist_Central_Hall.jpg",
	}

	// Create messages
	messages := []ai.InputMessage{
		{
			Role:    "user",
			Content: "Describe a beautiful sunset in 3 sentences.",
			Images:  []ai.Image{Image},
		},
	}

	// Call the API
	response, err := client.TextCompletion(ctx, messages, config)
	if err != nil {
		log.Fatalf("Error calling Bedrock: %v", err)
	}

	// Print the response
	fmt.Println("Bedrock Response:")
	fmt.Println(response.Text)
	fmt.Printf("Tokens used: %d input, %d output, %d total\n", response.TokenUsage.InputTokens, response.TokenUsage.OutputTokens, response.TokenUsage.TotalTokens)

}
