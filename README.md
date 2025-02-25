# AI Router Library for Golang

This repository provides a Golang library for interacting with various AI models such as OpenAI (GPT), Google Gemini, and AWS Bedrock. The library supports both text completion and multimodal image recognition.

## Installation

Ensure you have Go installed (1.18+ recommended). To use this library, install the required dependencies:

```sh
go get github.com/NaheedRayan/ai-router/ai
```

Also, install `godotenv` for environment variable management:

```sh
go get github.com/joho/godotenv
```

## Environment Variables

Create a `.env` file in your project root and add your API keys:

```
OPENAI_API_KEY=your_openai_key
GEMINI_API_KEY=your_gemini_key
AWS_ACCESS_KEY=your_aws_access_key
AWS_SECRET_KEY=your_aws_secret_key
```

## Usage

### Import the Library

```go
import (
    "context"
    "github.com/NaheedRayan/ai-router/ai"
)
```

### Initialize the Client

You can initialize a client for different AI providers:

#### OpenAI Example

```go
ctx := context.Background()
client, err := ai.InitializeClient(ctx, ai.ProviderOpenAI,
    ai.WithAPIKey("your_openai_key"),
    ai.WithModelID("gpt-4o-mini"),
)
if err != nil {
    log.Fatalf("Failed to initialize OpenAI client: %v", err)
}
defer client.Close()
```

#### Gemini Example

```go
client, err := ai.InitializeClient(ctx, ai.ProviderGemini,
    ai.WithAPIKey("your_gemini_key"),
    ai.WithModelID("models/gemini-2.0-flash-lite-preview-02-05"),
)
```

#### AWS Bedrock Example

```go
client, err := ai.InitializeClient(ctx, ai.ProviderBedrock,
    ai.WithAPIKey("your_aws_access_key"),
    ai.WithEndpointURL("your_aws_secret_key"),
    ai.WithRegion("us-east-1"),
    ai.WithModelID("amazon.nova-lite-v1:0"),
)
```

### Sending Requests

#### Text Completion

```go
messages := []ai.InputMessage{
    { Role: "user", Content: "Write a haiku about programming." },
}

config := ai.ModelConfig{
    Temperature:  0.7,
    TopP:         0.9,
    MaxTokens:    100,
    SystemPrompt: "You are a helpful AI assistant.",
}

response, err := client.TextCompletion(ctx, messages, config)
if err != nil {
    log.Fatalf("Error calling AI model: %v", err)
}

fmt.Println("AI Response:", response.Text)
```

#### Multimodal Image Recognition

```go
image := ai.Image{
    Format: "jpeg",
    URL: "https://example.com/sample.jpg",
}

messages := []ai.InputMessage{
    { Role: "user", Content: "What is shown in this image?", Images: []ai.Image{image} },
}

response, err := client.ImageRecognition(ctx, messages, config)
if err != nil {
    log.Fatalf("Error calling AI model: %v", err)
}

fmt.Println("AI Response:", response.Text)
```

## Running the Example

To run the example provided in `main.go`, use:

```sh
go run main.go
```

## Contributing

Feel free to submit pull requests or report issues if you find any bugs or improvements.

## License

This project is licensed under the MIT License.



