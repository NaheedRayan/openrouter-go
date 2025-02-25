package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/joho/godotenv"
)

// Load environment variables from .env file
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found or couldn't load")
	}
}

// Get environment variable or default value
func GetEnvDefault(key, defVal string) string {
	val, ex := os.LookupEnv(key)
	if !ex {
		log.Printf("Environment variable %s not found, using default value\n", key)
		return defVal
	}
	return val
}

func main() {
	// Load .env file
	LoadEnv()

	AWS_ACCESS_KEY := GetEnvDefault("AWS_ACCESS_KEY", "")
	AWS_SECRET_KEY := GetEnvDefault("AWS_SECRET_KEY", "")
	AWS_REGION := GetEnvDefault("AWS_REGION", "us-east-1")

	// Print loaded env variables (for debugging)
	fmt.Println("AWS_ACCESS_KEY:", AWS_ACCESS_KEY)
	fmt.Println("AWS_SECRET_KEY:", AWS_SECRET_KEY)
	fmt.Println("AWS_REGION:", AWS_REGION)

	// Ensure credentials are available
	if AWS_ACCESS_KEY == "" || AWS_SECRET_KEY == "" {
		log.Fatal("Missing AWS credentials. Check your .env file or environment variables.")
	}

	ctx := context.Background()

	cred := credentials.NewStaticCredentialsProvider(AWS_ACCESS_KEY, AWS_SECRET_KEY, "")

	_cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(cred),
		config.WithRegion(AWS_REGION),
	)

	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	client := bedrockruntime.NewFromConfig(_cfg)

	// Prepare request payload
	requestPayload := map[string]interface{}{
		"schemaVersion": "messages-v1",
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{"text": "hi how are you?"},
				},
			},
		},
		"system": []map[string]string{
			{"text": "You are an expert artist. "},
		},
		"inferenceConfig": map[string]interface{}{
			"maxTokens":   300,
			"topP":        0.1,
			"topK":        20,
			"temperature": 0.3,
		},
	}

	// Marshal request body to JSON
	jsonBytes, err := json.Marshal(requestPayload)
	if err != nil {
		log.Fatalf("error marshaling request: %v", err)
	}

	// Prepare the Bedrock API request
	input := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String("amazon.nova-lite-v1:0"),
		Body:        jsonBytes,
		ContentType: aws.String("application/json"),
	}

	// Call the Bedrock API
	response, err := client.InvokeModel(context.TODO(), input)
	if err != nil {
		log.Fatalf("error calling Bedrock API: %v", err)
	}

	// Parse the response
	var responseBody struct {
		Output struct {
			Message struct {
				Content []struct {
					Text string `json:"text"`
				} `json:"content"`
			} `json:"message"`
		} `json:"output"`
	}
	if err := json.Unmarshal(response.Body, &responseBody); err != nil {
		log.Fatalf("error unmarshaling response: %v", err)
	}

	// Print the response
	fmt.Println(responseBody.Output.Message.Content[0].Text)
}
