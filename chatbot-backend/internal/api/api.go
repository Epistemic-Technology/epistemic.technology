package api

import (
	"encoding/json"
	"net/http"

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

func HandleChat(w http.ResponseWriter, r *http.Request) {
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

	// Get the chatbot from the context
	bot, ok := r.Context().Value("chatbot").(*chatbot.ChatBot)
	if !ok {
		http.Error(w, "Chatbot not available", http.StatusInternalServerError)
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
