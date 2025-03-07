package chatbot

import (
	"os"
	"strings"
	"testing"

	"github.com/Epistemic-Technology/epistemic.technology/chatbot-backend/internal/backend"
	"github.com/joho/godotenv"
)

func setupTestEnvironment(t *testing.T) (*ChatBot, func()) {
	// Create a temporary database file
	tempFile, err := os.CreateTemp("", "test-db-*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	dbPath := tempFile.Name()

	// Clean up after the test
	cleanup := func() {
		os.Remove(dbPath)
	}

	// Initialize the database
	database, err := backend.GetDB(dbPath)
	if err != nil {
		cleanup()
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Create embedding client
	embeddingClient, err := backend.NewEmbeddingClient()
	if err != nil {
		cleanup()
		t.Fatalf("Failed to create embedding client: %v", err)
	}

	// Create LLM client
	llmClient, err := backend.NewLLMClient()
	if err != nil {
		cleanup()
		t.Fatalf("Failed to create LLM client: %v", err)
	}

	// Create and insert a test document
	doc := &backend.Document{
		Title:           "Test Document",
		Content:         "This is a test document about artificial intelligence and machine learning.",
		Author:          "Test Author",
		PublicationDate: "2023-01-01",
		URL:             "https://example.com/test",
	}

	err = backend.InsertDocument(database, doc)
	if err != nil {
		t.Fatalf("Failed to insert document: %v", err)
	}

	// Insert a test chunk
	chunk := &backend.Chunk{
		DocumentID: doc.ID,
		Content:    "This is a test chunk about artificial intelligence.",
		Hash:       []byte{1, 2, 3, 4},
		Embedding:  backend.Embedding{0.1, 0.2, 0.3, 0.4},
	}

	err = backend.InsertChunk(database, chunk)
	if err != nil {
		t.Fatalf("Failed to insert chunk: %v", err)
	}

	// Return cleanup function
	return NewChatBot(database, embeddingClient, llmClient), cleanup
}

func TestChat(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	// Skip test if OPENAI_API_KEY is not set
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("OPENAI_API_KEY not set, skipping test")
	}

	chatbot, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Test with a simple query
	query := "What is artificial intelligence?"
	history := ""
	userID := 1

	// Call the Chat function
	response, references, sources, err := Chat(chatbot, userID, query, history)
	if err != nil {
		t.Fatalf("Chat failed: %v", err)
	}

	// Verify response is not empty
	if response == "" {
		t.Error("Expected non-empty response, got empty string")
	}

	// Verify references were returned
	if len(references) == 0 {
		t.Error("Expected references, got none")
	}

	// Verify sources were returned
	if len(sources) == 0 {
		t.Error("Expected sources, got none")
	}

	// Verify the content of the references
	foundTestContent := false
	for _, ref := range references {
		if ref.Content == "This is a test document with some content for testing the chatbot." {
			foundTestContent = true
			break
		}
	}
	if !foundTestContent {
		t.Error("Expected to find test content in references")
	}
}

func TestChatWithNoRelevantDocuments(t *testing.T) {
	// Skip test if OPENAI_API_KEY is not set
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("OPENAI_API_KEY not set, skipping test")
	}

	chatbot, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Test with a query that should not match any documents
	query := "What is the meaning of life?"
	history := ""
	userID := 1

	// Call the Chat function
	response, references, sources, err := Chat(chatbot, userID, query, history)
	if err != nil {
		t.Fatalf("Chat failed: %v", err)
	}

	// Verify response is not empty
	if response == "" {
		t.Error("Expected non-empty response, got empty string")
	}

	// Even with an unrelated query, we should still get some references
	// since we're doing a similarity search with a limit
	if len(references) == 0 {
		t.Error("Expected references, got none")
	}

	// Verify sources were returned
	if len(sources) == 0 {
		t.Error("Expected sources, got none")
	}
}

func TestBuildUserQuery(t *testing.T) {
	// Create test chunks
	chunks := []backend.Chunk{
		{
			DocumentID: 1,
			Content:    "This is the first test chunk.",
			Hash:       []byte{1, 2, 3, 4},
		},
		{
			DocumentID: 2,
			Content:    "This is the second test chunk.",
			Hash:       []byte{5, 6, 7, 8},
		},
	}

	// Test with empty history
	query := "What is artificial intelligence?"
	history := ""
	result := buildUserQuery(query, history, chunks)

	// Verify the result contains the expected components
	if !contains(result, "This is our conversation history: ") {
		t.Error("Expected result to contain conversation history header")
	}
	if !contains(result, "This is a list of documents that are relevant to the conversation: ") {
		t.Error("Expected result to contain documents header")
	}
	if !contains(result, "Document ID: 1") {
		t.Error("Expected result to contain first document ID")
	}
	if !contains(result, "Document ID: 2") {
		t.Error("Expected result to contain second document ID")
	}
	if !contains(result, "This is the first test chunk.") {
		t.Error("Expected result to contain first chunk content")
	}
	if !contains(result, "This is the second test chunk.") {
		t.Error("Expected result to contain second chunk content")
	}
	if !contains(result, "This is the user's query: What is artificial intelligence?") {
		t.Error("Expected result to contain user query")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestChatWithHistory(t *testing.T) {
	// Skip test if OPENAI_API_KEY is not set
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("OPENAI_API_KEY not set, skipping test")
	}

	chatbot, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Test with history
	query := "Tell me more about it"
	userID := 1
	history := "User: What is in the test document?\nBot: The test document contains information about testing the chatbot."

	// Call the Chat function
	response, references, sources, err := Chat(chatbot, userID, query, history)
	if err != nil {
		t.Fatalf("Chat failed: %v", err)
	}

	// Verify response is not empty
	if response == "" {
		t.Error("Expected non-empty response, got empty string")
	}

	// Verify references were returned
	if len(references) == 0 {
		t.Error("Expected references, got none")
	}

	// Verify sources were returned
	if len(sources) == 0 {
		t.Error("Expected sources, got none")
	}
}

func TestChatWithEmbeddingError(t *testing.T) {
	// Save the original environment variable
	originalAPIKey := os.Getenv("OPENAI_API_KEY")
	defer os.Setenv("OPENAI_API_KEY", originalAPIKey)

	// Unset the API key to force an error
	os.Unsetenv("OPENAI_API_KEY")

	// Create a temporary database file
	tempFile, err := os.CreateTemp("", "test-db-*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	dbPath := tempFile.Name()
	defer os.Remove(dbPath)

	// Initialize the database
	database, err := backend.GetDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer backend.Close(database)

	// Create embedding client (should fail due to missing API key)
	_, err = backend.NewEmbeddingClient()
	if err == nil {
		t.Fatal("Expected error when creating embedding client without API key")
	}
}
