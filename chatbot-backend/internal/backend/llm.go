package backend

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type LLMClient struct {
	client *openai.Client
}

func NewLLMClient() (*LLMClient, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &LLMClient{client: client}, nil
}

//go:embed system_prompt.md
var systemPrompt string

func Chat(c *LLMClient, query string) (string, error) {
	ctx := context.Background()

	response, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F(
			[]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(query),
				openai.SystemMessage(systemPrompt),
			},
		),
		Model: openai.F(openai.ChatModelGPT4oMini),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	return response.Choices[0].Message.Content, nil
}
