package backend

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Embedding []float64

// Client wraps the OpenAI client for embedding operations
type EmbeddingClient struct {
	client *openai.Client
}

// NewEmbeddingClient creates a new embedding client using the OpenAI API key from environment
func NewEmbeddingClient() (*EmbeddingClient, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &EmbeddingClient{client: client}, nil
}

// CreateEmbedding generates an embedding vector for a single string
func CreateEmbedding(c *EmbeddingClient, text string, userID int) (Embedding, error) {
	embeddings, err := CreateEmbeddings(c, []string{text}, userID)
	if err != nil {
		return nil, err
	}

	if len(embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return embeddings[0], nil
}

// CreateEmbeddings generates embedding vectors for multiple strings
func CreateEmbeddings(c *EmbeddingClient, texts []string, userID int) ([]Embedding, error) {
	if len(texts) == 0 {
		return []Embedding{}, nil
	}

	ctx := context.Background()
	userIDStr := strconv.Itoa(userID)

	embedding, err := c.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.F(
			openai.EmbeddingNewParamsInputUnion(
				openai.EmbeddingNewParamsInputArrayOfStrings(texts),
			),
		),
		Model:          openai.F(openai.EmbeddingModelTextEmbedding3Small),
		EncodingFormat: openai.F(openai.EmbeddingNewParamsEncodingFormatFloat),
		User:           openai.F(userIDStr),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create embeddings: %w", err)
	}

	result := make([]Embedding, len(embedding.Data))
	for i, embeddingData := range embedding.Data {
		result[i] = embeddingData.Embedding
	}

	return result, nil
}
