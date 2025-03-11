package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Epistemic-Technology/epistemic.technology/chatbot-backend/internal/backend"
	"github.com/Epistemic-Technology/epistemic.technology/chatbot-backend/internal/chatbot"
)

type ChatRequest struct {
	Query   string `json:"query"`
	History string `json:"history"`
}

type ChatResponse struct {
	Response   string             `json:"response"`
	References []backend.Chunk    `json:"references"`
	Sources    []backend.Document `json:"sources"`
}

func StartAPI(bot *chatbot.ChatBot) {
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		HandleChat(w, r, bot)
	})

	port := ":" + os.Getenv("PORT")
	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func HandleChat(w http.ResponseWriter, r *http.Request, bot *chatbot.ChatBot) {
	// Set CORS headers
	if setCORSHeaders(w, r) {
		return
	}

	// Parse the request
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Process the chat request
	response, references, sources, err := chatbot.Chat(bot, 1, req.Query, req.History)
	if err != nil {
		http.Error(w, "Error processing chat: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	resp := ChatResponse{
		Response:   response,
		References: references,
		Sources:    sources,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func setCORSHeaders(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return true
	}
	return false
}
