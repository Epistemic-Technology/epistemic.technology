package chatbot

import (
	"fmt"
	"log"
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

func EmbedHugoDirectory(c *ChatBot, directory string) error {
	log.Println("Embedding Hugo directory: ", directory)
	docs, err := backend.HugoDirectoryToDocuments(directory, true)
	if err != nil {
		return fmt.Errorf("failed to convert directory to documents: %w", err)
	}
	log.Println("Converted directory to documents: ", len(docs))
	
	user := &backend.User{ID: 1}

	for i, doc := range docs {
		log.Printf("Processing document %d/%d: %s", i+1, len(docs), doc.Title)

		if err := backend.InsertDocument(c.db, &doc); err != nil {
			return fmt.Errorf("failed to insert document: %v", err)
		}

		chunks, err := backend.ChunkDocument(&doc, c.embeddingClient, user, c.db)
		if err != nil {
			return fmt.Errorf("error creating chunks: %v", err)
		}

		log.Printf("Created %d chunks", len(chunks))

		for j, chunk := range chunks {
			chunk.DocumentID = doc.ID
			if err := backend.InsertChunk(c.db, &chunk); err != nil {
				return fmt.Errorf("error inserting chunk %d: %v", j, err)
			}
		}
	}

	return nil
}

