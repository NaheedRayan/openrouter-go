package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/joho/godotenv"
	"github.com/nfnt/resize"
)

// Converts image URL to Base64
func urlToBase64(imageURL string) (string, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return "", err
	}

	// Resize image
	m := resize.Resize(256, 256, img, resize.Lanczos3)

	// Encode image to JPEG
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, m, &jpeg.Options{Quality: 50})
	if err != nil {
		return "", err
	}

	// Convert to Base64
	photoBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	return photoBase64, nil
}

// Load environment variables
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found or couldn't load")
	}
}

// Get environment variable or default value
func GetEnvDefault(key, defVal string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("Environment variable %s not found, using default value\n", key)
		return defVal
	}
	return val
}

func main() {
	// Load environment variables
	LoadEnv()

	AWS_ACCESS_KEY := GetEnvDefault("AWS_ACCESS_KEY", "")
	AWS_SECRET_KEY := GetEnvDefault("AWS_SECRET_KEY", "")
	AWS_REGION := GetEnvDefault("AWS_REGION", "us-east-1")

	if AWS_ACCESS_KEY == "" || AWS_SECRET_KEY == "" {
		log.Fatal("Missing AWS credentials. Check your .env file or environment variables.")
	}

	ctx := context.Background()

	// AWS Credentials
	cred := credentials.NewStaticCredentialsProvider(AWS_ACCESS_KEY, AWS_SECRET_KEY, "")

	// Load AWS Config
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(cred),
		config.WithRegion(AWS_REGION),
	)
	if err != nil {
		log.Fatalf("Unable to load AWS SDK config: %v", err)
	}

	client := bedrockruntime.NewFromConfig(cfg)

	// Image URL
	imageURL := "https://fastly.picsum.photos/id/26/4209/2769.jpg?hmac=vcInmowFvPCyKGtV7Vfh7zWcA_Z0kStrPDW3ppP0iGI" // Example image
	base64Image, err := urlToBase64(imageURL)
	if err != nil {
		log.Fatalf("Error encoding image: %v", err)
	}

	imageURL_2 := "https://picsum.photos/id/237/200/300" // Example image
	base64Image2, err := urlToBase64(imageURL_2)
	if err != nil {
		log.Fatalf("Error encoding image: %v", err)
	}

	// Prepare request payload
	requestPayload := map[string]interface{}{
		"schemaVersion": "messages-v1",
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"image": map[string]interface{}{
							"format": "jpeg",
							"source": map[string]string{"bytes": base64Image},
						},
					},
					{
						"image": map[string]interface{}{
							"format": "jpeg",
							"source": map[string]string{"bytes": base64Image2},
						},
					},
					{
						"text": "What are the images",
					},
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

	// Convert request payload to JSON
	jsonBytes, err := json.Marshal(requestPayload)
	if err != nil {
		log.Fatalf("Error marshalling request: %v", err)
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

	// Parse the response
	// var responseBody bedrock.ConverseResponse
	if err := json.Unmarshal(response.Body, &responseBody); err != nil {
		log.Fatalf("error unmarshaling response: %v", err)
	}

	// Print the response
	fmt.Println(responseBody.Output.Message.Content[0].Text)
}
