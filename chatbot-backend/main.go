package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Epistemic-Technology/epistemic.technology/chatbot-backend/internal/backend"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go" // imported as openai
	"github.com/openai/openai-go/option"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	client := openai.NewClient(option.WithAPIKey(os.Getenv("OPENAI_API_KEY")))
	fmt.Println(client)

	// Example usage of the backend package
	_ = backend.Document{}
}
