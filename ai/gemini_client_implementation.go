package ai

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeminiClient implements the Client interface for Google Gemini
type GeminiClient struct {
	client  *genai.Client
	model   *genai.GenerativeModel
	options ClientOptions
}

// Ensure GeminiClient implements Client interface
var _ Client = (*GeminiClient)(nil)

// NewGeminiClient creates a new Google Gemini client
func NewGeminiClient() *GeminiClient {
	return &GeminiClient{}
}

// Initialize sets up the Gemini client
func (c *GeminiClient) Initialize(ctx context.Context, opts ...ClientOption) error {
	// Apply options
	options := ClientOptions{
		ModelID: "models/gemini-2.0-flash-lite-preview-02-05", // Default model
	}

	for _, opt := range opts {
		opt(&options)
	}

	c.options = options

	// Create the Gemini client
	client, err := genai.NewClient(ctx, option.WithAPIKey(options.APIKey))
	if err != nil {
		return fmt.Errorf("failed to create Gemini client: %v", err)
	}

	c.client = client

	// Create the model
	model := client.GenerativeModel(options.ModelID)
	c.model = model

	return nil
}

// TextCompletion sends a text request to Gemini
func (c *GeminiClient) TextCompletion(ctx context.Context, messages []InputMessage, config ModelConfig) (Response, error) {
	// Apply configuration
	c.model.SetTemperature(float32(config.Temperature))
	c.model.SetTopP(float32(config.TopP))
	c.model.SetTopK(int32(config.TopK))
	c.model.SetMaxOutputTokens(int32(config.MaxTokens))

	// Set system prompt if provided
	if config.SystemPrompt != "" {
		c.model.SystemInstruction = genai.NewUserContent(genai.Text(config.SystemPrompt))
	}

	// Extract the last message for prompt
	if len(messages) == 0 {
		return Response{}, fmt.Errorf("no messages provided")
	}

	prompt := messages[len(messages)-1].Content

	// Generate content
	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return Response{}, fmt.Errorf("failed to generate content: %v", err)
	}

	// Format the standard response
	result := Response{
		Raw: resp,
	}

	// Extract usage information if available
	if resp.UsageMetadata != nil {
		result.TokenUsage = TokenUsage{
			InputTokens:  int(resp.UsageMetadata.PromptTokenCount),
			OutputTokens: int(resp.UsageMetadata.CandidatesTokenCount),
			TotalTokens:  int(resp.UsageMetadata.TotalTokenCount),
		}
	}

	// Extract text from response
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil && len(resp.Candidates[0].Content.Parts) > 0 {
		part := resp.Candidates[0].Content.Parts[0]
		if str, ok := part.(genai.Text); ok {
			result.Text = string(str)
		}
	}

	return result, nil
}

// ImageRecognition sends images with optional text to Gemini
func (c *GeminiClient) ImageRecognition(ctx context.Context, messages []InputMessage, config ModelConfig) (Response, error) {
	// Apply configuration
	c.model.SetTemperature(float32(config.Temperature))
	c.model.SetTopP(float32(config.TopP))
	c.model.SetTopK(int32(config.TopK))
	c.model.SetMaxOutputTokens(int32(config.MaxTokens))

	// Set system prompt if provided
	if config.SystemPrompt != "" {
		c.model.SystemInstruction = genai.NewUserContent(genai.Text(config.SystemPrompt))
	}

	// Extract the last message
	if len(messages) == 0 {
		return Response{}, fmt.Errorf("no messages provided")
	}

	message := messages[len(messages)-1]

	// Create parts for the request
	var parts []genai.Part

	// Add images to parts
	for _, img := range message.Images {
		if len(img.Data) > 0 {
			// Use inline image data
			parts = append(parts, genai.ImageData(img.Format, img.Data))
		} else if img.URL != "" {
			// Download from URL
			resp, err := http.Get(img.URL)
			if err != nil {
				return Response{}, fmt.Errorf("failed to download image: %v", err)
			}
			defer resp.Body.Close()

			imgBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return Response{}, fmt.Errorf("failed to read image data: %v", err)
			}

			parts = append(parts, genai.ImageData(img.Format, imgBytes))
		}
	}

	// Add text to parts if present
	if message.Content != "" {
		parts = append(parts, genai.Text(message.Content))
	}

	// Generate content
	resp, err := c.model.GenerateContent(ctx, parts...)
	if err != nil {
		return Response{}, fmt.Errorf("failed to generate content: %v", err)
	}

	// Format the standard response
	result := Response{
		Raw: resp,
	}

	// Extract usage information if available
	if resp.UsageMetadata != nil {
		result.TokenUsage = TokenUsage{
			InputTokens:  int(resp.UsageMetadata.PromptTokenCount),
			OutputTokens: int(resp.UsageMetadata.CandidatesTokenCount),
			TotalTokens:  int(resp.UsageMetadata.TotalTokenCount),
		}
	}

	// Extract text from response
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil && len(resp.Candidates[0].Content.Parts) > 0 {
		part := resp.Candidates[0].Content.Parts[0]
		if str, ok := part.(genai.Text); ok {
			result.Text = string(str)
		}
	}

	return result, nil
}

// Close releases resources
func (c *GeminiClient) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
