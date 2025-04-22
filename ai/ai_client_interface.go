package ai

import (
	"context"
)

// InputMessage represents a single message in a conversation
type InputMessage struct {
	Role    string
	Content string
	Images  []Image // Optional images for multimodal models
}

// Image represents an image to be processed by AI models
type Image struct {
	Format string
	Data   []byte
	URL    string // Optional URL alternative to inline data
}

// ModelConfig represents configuration parameters for an AI model
type ModelConfig struct {
	Temperature   float32
	TopP          float32
	TopK          int32
	MaxTokens     int32
	SystemPrompt  string
	StopSequences []string
}

// Response represents a standardized response from any AI provider
type Response struct {
	Text       string
	TokenUsage TokenUsage
	Raw        interface{} // Raw provider-specific response
}

// TokenUsage stores token usage information
type TokenUsage struct {
	InputTokens  int
	OutputTokens int
	TotalTokens  int
}

// Client is the common interface for all AI providers
type Client interface {
	// Initialize initializes the client with the given options
	Initialize(ctx context.Context, opts ClientOptions) error

	// TextCompletion sends a text request to the AI provider
	TextCompletion(ctx context.Context, messages []InputMessage, config ModelConfig) (Response, error)

	// ImageRecognition sends images with optional text for processing
	ImageRecognition(ctx context.Context, messages []InputMessage, config ModelConfig) (Response, error)

	// Close releases any resources held by the client
	Close() error
}

// ClientOptions Current Approach (Functional Options Pattern)
// ClientOptions contains all configuration options
type ClientOptions struct {
	AccessKey   string
	SecretKey   string
	APIKey      string
	Region      string
	EndpointURL string
	ModelID     string
	Timeout     int
}

//
//// ClientOption is used for configuring the client
//type ClientOption func(*ClientOptions)
//
//// WithAPIKey sets the API key for authentication
//func WithAPIKey(apiKey string) ClientOption {
//	return func(o *ClientOptions) {
//		o.APIKey = apiKey
//	}
//}
//
//// WithRegion sets the region for services like AWS
//func WithRegion(region string) ClientOption {
//	return func(o *ClientOptions) {
//		o.Region = region
//	}
//}
//
//// WithEndpointURL sets a custom endpoint URL
//func WithEndpointURL(url string) ClientOption {
//	return func(o *ClientOptions) {
//		o.EndpointURL = url
//	}
//}
//
//// WithModelID sets the default model ID to use
//func WithModelID(modelID string) ClientOption {
//	return func(o *ClientOptions) {
//		o.ModelID = modelID
//	}
//}
//
//// WithAccessKey sets the access key for authentication
//func WithAccessKey(accessKey string) ClientOption {
//	return func(o *ClientOptions) {
//		o.AccessKey = accessKey
//	}
//}
//
//// WithSecretKey sets the secret key for authentication
//func WithSecretKey(secretKey string) ClientOption {
//	return func(o *ClientOptions) {
//		o.SecretKey = secretKey
//	}
//}
//
//// WithTimeout sets a timeout for requests in seconds
//func WithTimeout(seconds int) ClientOption {
//	return func(o *ClientOptions) {
//		o.Timeout = seconds
//	}
//}
