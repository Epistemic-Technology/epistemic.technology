package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Epistemic-Technology/epistemic.technology/chatbot-backend/internal/backend"
	"github.com/Epistemic-Technology/epistemic.technology/chatbot-backend/internal/chatbot"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	dbPathFlag := flag.String("db", "", "Path to the database file (overrides DATABASE_PATH env var)")
	apiKeyFlag := flag.String("api-key", "", "OpenAI API key (overrides OPENAI_API_KEY env var)")
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

	fmt.Println("Welcome to the Chatbot CLI!")
	fmt.Println("Type 'exit' or 'quit' to end the session.")
	fmt.Println("Type 'clear' to clear the conversation history.")
	fmt.Println("---------------------------------------------")

	history := ""
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\nYou: ")
		if !scanner.Scan() {
			break
		}

		userInput := scanner.Text()
		if userInput == "" {
			continue
		}

		if userInput == "exit" || userInput == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		if userInput == "clear" {
			history = ""
			fmt.Println("Conversation history cleared.")
			continue
		}

		if history != "" {
			history += "\n"
		}
		history += "User: " + userInput

		response, _, sources, err := chatbot.Chat(bot, 1, userInput, history)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		history += "\nBot: " + response

		fmt.Println("\nBot:", response)

		if len(sources) > 0 {
			fmt.Println("\nSources:")
			sourceMap := make(map[int]string)
			for _, source := range sources {
				if source.Title != "" {
					sourceMap[source.ID] = source.Title
				} else {
					sourceMap[source.ID] = source.FilePath
				}
			}

			for _, title := range sourceMap {
				fmt.Printf("- %s\n", title)
			}
		}
	}

	if scanner.Err() != nil {
		log.Fatalf("Error reading input: %v", scanner.Err())
	}
}
