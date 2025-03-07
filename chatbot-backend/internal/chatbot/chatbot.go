package chatbot

import (
	"fmt"
	"strconv"

	"github.com/Epistemic-Technology/epistemic.technology/chatbot-backend/internal/backend"
)

type ChatBot struct {
	db              *backend.DB
	embeddingClient *backend.EmbeddingClient
	llmClient       *backend.LLMClient
}

func NewChatBot(db *backend.DB, embeddingClient *backend.EmbeddingClient, llmClient *backend.LLMClient) *ChatBot {
	return &ChatBot{
		db:              db,
		embeddingClient: embeddingClient,
		llmClient:       llmClient,
	}
}

func Chat(c *ChatBot, userID int, query string, history string) (response string, references []backend.Chunk, sources []backend.Document, err error) {
	queryEmbedding, err := backend.CreateEmbedding(c.embeddingClient, query, userID)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to create embedding: %w", err)
	}

	// Get similar chunks
	chunks, err := backend.SimilaritySearch(c.db, queryEmbedding, 5)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to search for similar chunks: %w", err)
	}

	// Get response from LLM
	finalQuery := buildUserQuery(query, history, chunks)
	response, err = backend.Chat(c.llmClient, finalQuery)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to get chat response: %w", err)
	}

	// Get source documents
	sources, err = backend.DocumentsFromChunks(chunks, c.db)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to get source documents: %w", err)
	}

	return response, chunks, sources, nil
}

func buildUserQuery(query string, history string, chunks []backend.Chunk) string {
	finalQuery := "This is our conversation history: " + history
	finalQuery += "\n\n"
	finalQuery += "This is a list of documents that are relevant to the conversation: "
	for _, chunk := range chunks {
		finalQuery += "Document ID: " + strconv.Itoa(chunk.DocumentID) + "\n"
		finalQuery += "Document Content: " + chunk.Content + "\n"
	}
	finalQuery += "This is the user's query: " + query
	return finalQuery
}
