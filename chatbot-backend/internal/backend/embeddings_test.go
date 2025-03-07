package backend

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestCreateEmbedding(t *testing.T) {
	// Load environment variables from .env file
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Logf("Error loading .env file: %v", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping test: OPENAI_API_KEY not set in environment")
	}

	client, err := NewEmbeddingClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test single embedding
	text := "This is a test sentence for embedding generation."
	userID := 123

	embedding, err := CreateEmbedding(client, text, userID)
	if err != nil {
		t.Fatalf("Failed to create embedding: %v", err)
	}

	// Check that we got a non-empty embedding
	if len(embedding) == 0 {
		t.Error("Expected non-empty embedding, got empty slice")
	}

	// The text-embedding-3-small model produces 1536-dimensional embeddings
	expectedDimensions := 1536
	if len(embedding) != expectedDimensions {
		t.Errorf("Expected embedding with %d dimensions, got %d", expectedDimensions, len(embedding))
	}
}

func TestCreateEmbeddings(t *testing.T) {
	// Load environment variables from .env file
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Logf("Error loading .env file: %v", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping test: OPENAI_API_KEY not set in environment")
	}

	client, err := NewEmbeddingClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test multiple embeddings
	texts := []string{
		"This is the first test sentence.",
		"This is the second test sentence.",
		"This is the third test sentence with different content.",
	}
	userID := 123

	embeddings, err := CreateEmbeddings(client, texts, userID)
	if err != nil {
		t.Fatalf("Failed to create embeddings: %v", err)
	}

	// Check that we got the right number of embeddings
	if len(embeddings) != len(texts) {
		t.Errorf("Expected %d embeddings, got %d", len(texts), len(embeddings))
	}

	// Check that each embedding has the expected dimensions
	expectedDimensions := 1536
	for i, embedding := range embeddings {
		if len(embedding) != expectedDimensions {
			t.Errorf("Embedding %d: expected %d dimensions, got %d", i, expectedDimensions, len(embedding))
		}
	}

	// Test empty input
	emptyEmbeddings, err := CreateEmbeddings(client, []string{}, userID)
	if err != nil {
		t.Fatalf("Failed on empty input: %v", err)
	}
	if len(emptyEmbeddings) != 0 {
		t.Errorf("Expected empty result for empty input, got %d embeddings", len(emptyEmbeddings))
	}
}
