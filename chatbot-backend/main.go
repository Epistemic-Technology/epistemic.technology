package main

import (
	"flag"
	"log"
	"os"

	"github.com/Epistemic-Technology/epistemic.technology/chatbot-backend/internal/api"
	"github.com/Epistemic-Technology/epistemic.technology/chatbot-backend/internal/backend"
	"github.com/Epistemic-Technology/epistemic.technology/chatbot-backend/internal/chatbot"
	"github.com/joho/godotenv"
	// imported as openai
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	dbPathFlag := flag.String("db", "", "Path to the database file (overrides DATABASE_PATH env var)")
	apiKeyFlag := flag.String("api-key", "", "OpenAI API key (overrides OPENAI_API_KEY env var)")
	portFlag := flag.String("port", "", "Port to run the API on (overrides PORT env var)")
	flag.Parse()

	dbPath := *dbPathFlag
	if dbPath == "" {
		dbPath = os.Getenv("DATABASE_PATH")
		if dbPath == "" {
			log.Fatal("Error: No database path provided. Use --db flag or set DATABASE_PATH environment variable")
		}
	}

	apiKey := *apiKeyFlag
	if apiKey != "" {
		os.Setenv("OPENAI_API_KEY", apiKey)
	} else if os.Getenv("OPENAI_API_KEY") == "" {
		log.Fatal("Error: No OpenAI API key provided. Use --api-key flag or set OPENAI_API_KEY environment variable")
	}

	port := *portFlag
	if port != "" {
		os.Setenv("PORT", port)
	} else if os.Getenv("PORT") == "" {
		log.Fatal("Error: No port provided. Use --port flag or set PORT environment variable")
	}

	database, err := backend.GetDB(dbPath)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer backend.Close(database)

	embeddingClient, err := backend.NewEmbeddingClient()
	if err != nil {
		log.Fatalf("Error creating embeddings client: %v", err)
	}

	llmClient, err := backend.NewLLMClient()
	if err != nil {
		log.Fatalf("Error creating LLM client: %v", err)
	}

	bot := chatbot.NewChatBot(database, embeddingClient, llmClient)

	api.StartAPI(bot)
}
