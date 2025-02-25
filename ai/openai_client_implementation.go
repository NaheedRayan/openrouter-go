package ai

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// OpenAIClient implements the Client interface for OpenAI
type OpenAIClient struct {
	client  *openai.Client
	options ClientOptions
	modelID string
}

// Ensure OpenAIClient implements Client interface
var _ Client = (*OpenAIClient)(nil)

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient() *OpenAIClient {
	return &OpenAIClient{}
}

// Initialize sets up the OpenAI client
func (c *OpenAIClient) Initialize(ctx context.Context, opts ...ClientOption) error {
	// Apply options
	options := ClientOptions{
		ModelID: string(openai.ChatModelGPT4oMini), // Default model
	}

	for _, opt := range opts {
		opt(&options)
	}

	c.options = options
	c.modelID = options.ModelID

	// Create the OpenAI client
	c.client = openai.NewClient(option.WithAPIKey(options.APIKey))

	return nil
}

// TextCompletion sends a text request to OpenAI
func (c *OpenAIClient) TextCompletion(ctx context.Context, messages []InputMessage, config ModelConfig) (Response, error) {
	// Convert to OpenAI format messages
	var openAIMessages []openai.ChatCompletionMessageParamUnion

	// Add system message if provided
	if config.SystemPrompt != "" {
		openAIMessages = append(openAIMessages, openai.SystemMessage(config.SystemPrompt))
	}

	// Add input messages
	for _, msg := range messages {
		switch msg.Role {
		case "system":
			openAIMessages = append(openAIMessages, openai.SystemMessage(msg.Content))
		case "assistant":
			openAIMessages = append(openAIMessages, openai.AssistantMessage(msg.Content))
		default: // Default to user message
			openAIMessages = append(openAIMessages, openai.UserMessage(msg.Content))
		}
	}

	// Create request parameters
	params := openai.ChatCompletionNewParams{
		Messages:    openai.F(openAIMessages),
		Model:       openai.F(c.modelID),
		Temperature: openai.F(float64(config.Temperature)),
		TopP:        openai.F(float64(config.TopP)),
		MaxTokens:   openai.Int(int64(config.MaxTokens)),
	}

	// // Add stop sequences if provided
	// if len(config.StopSequences) > 0 {
	// 	params.Stop = openai.F(config.StopSequences)
	// }

	// Send request
	response, err := c.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return Response{}, fmt.Errorf("error getting response: %v", err)
	}

	// Format the standard response
	result := Response{
		Raw: response,
		TokenUsage: TokenUsage{
			InputTokens:  int(response.Usage.PromptTokens),
			OutputTokens: int(response.Usage.CompletionTokens),
			TotalTokens:  int(response.Usage.TotalTokens),
		},
	}

	// Extract text from response
	if len(response.Choices) > 0 {
		result.Text = response.Choices[0].Message.Content
	}

	return result, nil
}

// ImageRecognition sends images with optional text to OpenAI
func (c *OpenAIClient) ImageRecognition(ctx context.Context, messages []InputMessage, config ModelConfig) (Response, error) {
	// OpenAI API requires a different endpoint for image analysis
	// We're using GPT-4 Vision for this example
	return Response{}, fmt.Errorf("method not implemented for OpenAI client yet")
}

// Close releases resources
func (c *OpenAIClient) Close() error {
	// OpenAI Go SDK doesn't require explicit cleanup
	return nil
}
