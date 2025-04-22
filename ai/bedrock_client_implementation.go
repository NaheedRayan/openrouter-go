package ai

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

// BedrockClient implements the Client interface for AWS Bedrock
type BedrockClient struct {
	client  *bedrockruntime.Client
	options ClientOptions
	modelID string
}

// NewBedrockClient creates a new AWS Bedrock client
func NewBedrockClient() *BedrockClient {
	return &BedrockClient{}
}

// Initialize sets up the Bedrock client with AWS credentials
func (c *BedrockClient) Initialize(ctx context.Context, opts ClientOptions) error {

	// Setup AWS credentials
	var awsConfig aws.Config
	var err error

	cred := credentials.NewStaticCredentialsProvider(opts.AccessKey, opts.SecretKey, "")

	awsConfig, err = config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(cred),
		config.WithRegion(opts.Region),
	)

	if err != nil {
		return fmt.Errorf("unable to load AWS SDK config: %v", err)
	}

	// Apply options
	c.options = opts
	c.client = bedrockruntime.NewFromConfig(awsConfig)
	c.modelID = opts.ModelID

	return nil
}

// TextCompletion sends a text request to Bedrock
func (c *BedrockClient) TextCompletion(ctx context.Context, messages []InputMessage, config ModelConfig) (Response, error) {
	// Convert to Bedrock format messages
	bedrockMessages := []map[string]interface{}{}

	for _, msg := range messages {
		content := []map[string]interface{}{
			{"text": msg.Content},
		}

		bedrockMessages = append(bedrockMessages, map[string]interface{}{
			"role":    msg.Role,
			"content": content,
		})
	}

	// Create request payload
	requestPayload := map[string]interface{}{
		"schemaVersion": "messages-v1",
		"messages":      bedrockMessages,
		"system": []map[string]string{
			{"text": config.SystemPrompt},
		},
		"inferenceConfig": map[string]interface{}{
			"maxTokens":   config.MaxTokens,
			"topP":        config.TopP,
			"topK":        config.TopK,
			"temperature": config.Temperature,
		},
	}

	// Marshal request body to JSON
	jsonBytes, err := json.Marshal(requestPayload)
	if err != nil {
		return Response{}, fmt.Errorf("error marshaling request: %v", err)
	}

	// Prepare the Bedrock API request
	input := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(c.modelID),
		Body:        jsonBytes,
		ContentType: aws.String("application/json"),
	}

	// Call the Bedrock API
	response, err := c.client.InvokeModel(ctx, input)
	if err != nil {
		return Response{}, fmt.Errorf("error calling Bedrock API: %v", err)
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
		Usage struct {
			InputTokens  int `json:"inputTokens"`
			OutputTokens int `json:"outputTokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(response.Body, &responseBody); err != nil {
		return Response{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	// Format the standard response
	result := Response{
		Raw: responseBody,
		TokenUsage: TokenUsage{
			InputTokens:  responseBody.Usage.InputTokens,
			OutputTokens: responseBody.Usage.OutputTokens,
			TotalTokens:  responseBody.Usage.InputTokens + responseBody.Usage.OutputTokens,
		},
	}

	// Extract the text from the response
	if len(responseBody.Output.Message.Content) > 0 {
		result.Text = responseBody.Output.Message.Content[0].Text
	}

	return result, nil
}

// ImageRecognition sends images with optional text to Bedrock
func (c *BedrockClient) ImageRecognition(ctx context.Context, messages []InputMessage, config ModelConfig) (Response, error) {
	// We'll only use the last message for simplicity
	if len(messages) == 0 {
		return Response{}, fmt.Errorf("no messages provided")
	}

	message := messages[len(messages)-1]

	// Create content array for the request
	var contentArray []map[string]interface{}

	// Add images to content array
	for _, img := range message.Images {
		imgContent := map[string]interface{}{
			"image": map[string]interface{}{
				"format": img.Format,
				"source": map[string]string{
					"bytes": base64.StdEncoding.EncodeToString(img.Data),
				},
			},
		}
		contentArray = append(contentArray, imgContent)
	}

	// Add text to content array if present
	if message.Content != "" {
		contentArray = append(contentArray, map[string]interface{}{
			"text": message.Content,
		})
	}

	// Create request payload
	requestPayload := map[string]interface{}{
		"schemaVersion": "messages-v1",
		"messages": []map[string]interface{}{
			{
				"role":    message.Role,
				"content": contentArray,
			},
		},
		"system": []map[string]string{
			{"text": config.SystemPrompt},
		},
		"inferenceConfig": map[string]interface{}{
			"maxTokens":   config.MaxTokens,
			"topP":        config.TopP,
			"topK":        config.TopK,
			"temperature": config.Temperature,
		},
	}

	// Marshal request body to JSON
	jsonBytes, err := json.Marshal(requestPayload)
	if err != nil {
		return Response{}, fmt.Errorf("error marshaling request: %v", err)
	}

	// Prepare the Bedrock API request
	input := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(c.modelID),
		Body:        jsonBytes,
		ContentType: aws.String("application/json"),
	}

	// Call the Bedrock API
	response, err := c.client.InvokeModel(ctx, input)
	if err != nil {
		return Response{}, fmt.Errorf("error calling Bedrock API: %v", err)
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
		Usage struct {
			InputTokens  int `json:"inputTokens"`
			OutputTokens int `json:"outputTokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(response.Body, &responseBody); err != nil {
		return Response{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	// Format the standard response
	result := Response{
		Raw: responseBody,
		TokenUsage: TokenUsage{
			InputTokens:  responseBody.Usage.InputTokens,
			OutputTokens: responseBody.Usage.OutputTokens,
			TotalTokens:  responseBody.Usage.InputTokens + responseBody.Usage.OutputTokens,
		},
	}

	// Extract the text from the response
	if len(responseBody.Output.Message.Content) > 0 {
		result.Text = responseBody.Output.Message.Content[0].Text
	}

	return result, nil
}

// Close releases any resources
func (c *BedrockClient) Close() error {
	// AWS SDK doesn't require explicit cleanup
	return nil
}
