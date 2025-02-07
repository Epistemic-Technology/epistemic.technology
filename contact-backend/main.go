package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type ContactSubmission struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func handleContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ALLOWED_ORIGIN"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var submission ContactSubmission
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&submission); err != nil {
		response := Response{
			Success: false,
			Message: "Error parsing JSON request body",
		}
		sendJSONResponse(w, http.StatusBadRequest, response)
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	if submission.Name == "" || submission.Email == "" || submission.Message == "" {
		response := Response{
			Success: false,
			Message: "Name, email, and message are required fields",
		}
		sendJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	sendSendGridEmail(submission)

	response := Response{
		Success: true,
		Message: "Contact submission received successfully",
	}
	sendJSONResponse(w, http.StatusOK, response)
}

func sendSendGridEmail(submission ContactSubmission) {
	go func() {
		from := mail.NewEmail(submission.Name, os.Getenv("CONTACT_SENDER_EMAIL"))
		to := mail.NewEmail("Epistemic Technology", os.Getenv("CONTACT_EMAIL"))
		message := mail.NewSingleEmail(
			from,
			submission.Subject,
			to,
			submission.Message,
			submission.Message,
		)
		message.ReplyTo = mail.NewEmail(submission.Name, submission.Email)
		log.Printf(
			"Sending contact email from %s <%s> to %s with subject %s",
			submission.Name,
			submission.Email,
			os.Getenv("CONTACT_EMAIL"),
			submission.Subject,
		)
		apiKey := os.Getenv("SENDGRID_API_KEY")
		if apiKey == "" {
			log.Printf("SENDGRID_API_KEY is not set")
			return
		}
		client := sendgrid.NewSendClient(apiKey)
		response, err := client.Send(message)
		if err != nil {
			log.Printf("Error sending email: %v", err)
			return
		}
		if response.StatusCode >= 400 {
			log.Printf("Error sending email: %v", response)
			return
		}
		log.Printf("Email sent successfully")
	}()
}

func sendJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func main() {
	if os.Getenv("ALLOWED_ORIGIN") == "" {
		os.Setenv("ALLOWED_ORIGIN", "http://localhost:1313")
	}
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8080")
	}
	if os.Getenv("CONTACT_EMAIL") == "" {
		os.Setenv("CONTACT_EMAIL", "nobody@example.invalid")
	}
	if os.Getenv("CONTACT_SENDER_EMAIL") == "" {
		os.Setenv("CONTACT_SENDER_EMAIL", "nobody@example.invalid")
	}

	// Set up the contact endpoint
	http.HandleFunc("/", handleContact)

	port := ":" + os.Getenv("PORT")
	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("Allowing CORS from: %s\n", os.Getenv("ALLOWED_ORIGIN"))
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
