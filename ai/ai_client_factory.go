package ai

import (
	"context"
	"fmt"
)

const (
	ProviderOpenAI  = "openai"
	ProviderGemini  = "gemini"
	ProviderBedrock = "bedrock"
)

// NewClient creates a new AI client for the specified provider
func NewClient(provider string) (Client, error) {
	switch provider {
	case ProviderOpenAI:
		return NewOpenAIClient(), nil
	case ProviderGemini:
		return NewGeminiClient(), nil
	case ProviderBedrock:
		return NewBedrockClient(), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// InitializeClient is a helper to create and initialize a client in one step
func InitializeClient(ctx context.Context, provider string, opts ClientOptions) (Client, error) {
	client, err := NewClient(provider)
	if err != nil {
		return nil, err
	}

	if err := client.Initialize(ctx, opts); err != nil {
		return nil, err
	}

	return client, nil
}
