package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sendgrid/sendgrid-go"
	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/time/rate"
)

type EmailRequest struct {
	Email string `json:"email"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type RateLimiter struct {
	visitors map[string]*rate.Limiter
	mu       sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*rate.Limiter),
	}
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(1*time.Minute), 5)
		rl.visitors[ip] = limiter
	}

	return limiter
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func validateEmail(email string) bool {
	_, err := sgmail.ParseEmail(email)
	return err == nil
}

func initDB() (*sql.DB, error) {
	// Ensure the /db/ directory exists
	if _, err := os.Stat("/db"); os.IsNotExist(err) {
		log.Printf("Creating /db directory")
		if err := os.MkdirAll("/db", 0755); err != nil {
			return nil, fmt.Errorf("failed to create /db directory: %v", err)
		}
	}
	db, err := sql.Open("sqlite3", "/db/subscriptions.db")
	if err != nil {
		return nil, err
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS subscription_attempts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL,
		token TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		confirmed BOOLEAN DEFAULT FALSE
	);`

	_, err = db.Exec(createTable)
	return db, err
}

func setCORSHeaders(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ALLOWED_ORIGIN"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		log.Printf("Responding to OPTIONS request")
		w.WriteHeader(http.StatusOK)
		return true
	}
	return false
}

func handleSignup(w http.ResponseWriter, r *http.Request, db *sql.DB, rl *RateLimiter) {
	log.Printf("Received request to handleSignup: %s", r.Method)

	if setCORSHeaders(w, r) {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ip := r.RemoteAddr
	limiter := rl.getLimiter(ip)
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	var req EmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !validateEmail(req.Email) {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
		return
	}

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM subscription_attempts WHERE email = ? AND created_at > datetime('now', '-1 hour')", req.Email).Scan(&count)
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if count > 0 && os.Getenv("GO_DEBUG") != "true" {
		response := Response{
			Success: false,
			Message: "A subscription attempt was made recently. Please check your email or wait before trying again.",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	go func() {
		_, err = db.Exec("DELETE FROM subscription_attempts WHERE created_at < datetime('now', '-1 hour')")
		if err != nil {
			log.Printf("Error cleaning up old subscription attempts: %v", err)
		}

		token, err := generateToken()
		if err != nil {
			log.Printf("Error generating token: %v", err)
			return
		}

		_, err = db.Exec("INSERT INTO subscription_attempts (email, token) VALUES (?, ?)", req.Email, token)
		if err != nil {
			log.Printf("Error storing subscription attempt: %v", err)
			return
		}

		from := sgmail.NewEmail("Epistemic Technology", os.Getenv("SENDGRID_FROM_EMAIL"))
		to := sgmail.NewEmail("", req.Email)
		subject := "Confirm your subscription to Epistemic Technology"
		confirmURL := fmt.Sprintf("%s/confirm/?token=%s&email=%s", os.Getenv("SITE_URL"), token, req.Email)
		plainTextContent := fmt.Sprintf("Please confirm your subscription to Epistemic Technology by clicking this link: %s\n\nIf you need assistance, please contact us at %s\n\nYour email will not be shared with anyone outside of Epistemic Technology and will not be used for marketing purposes.", confirmURL, os.Getenv("SUPPORT_EMAIL"))
		htmlContent := fmt.Sprintf(`
			<h2>Confirm your subscription to Epistemic Technology</h2>
			<p>Please click the link below to confirm your subscription to our blog:</p>
			<p><a href="%s">Confirm Subscription</a></p>
			<p>If you need assistance, please contact us at %s</p>
			<p><em>Your email will not be shared with anyone outside of Epistemic Technology and will not be used for marketing purposes.</em></p>
		`, confirmURL, os.Getenv("SUPPORT_EMAIL"))

		message := sgmail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		message.ReplyTo = sgmail.NewEmail("Epistemic Technology", os.Getenv("SUPPORT_EMAIL"))
		client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
		confirmationResponse, err := client.Send(message)
		if err != nil {
			log.Printf("Error sending confirmation email: %v", err)
		}
		if confirmationResponse.StatusCode >= 400 {
			log.Printf("Confirmation email sent to %s with status code %d and body %s", req.Email, confirmationResponse.StatusCode, confirmationResponse.Body)
		} else {
			log.Printf("Confirmation email sent to %s", req.Email)
		}
	}()

	response := Response{
		Success: true,
		Message: "Please check your email to confirm your subscription.",
	}
	json.NewEncoder(w).Encode(response)
}

func handleConfirm(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Printf("Received request to handleConfirm: %s", r.Method)

	if setCORSHeaders(w, r) {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Token string `json:"token"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Token == "" || req.Email == "" {
		http.Error(w, "Missing token or email", http.StatusBadRequest)
		log.Printf("Missing token or email")
		return
	}

	var dbToken string
	var dbEmail string
	err := db.QueryRow("SELECT token, email FROM subscription_attempts WHERE token = ? AND email = ?", req.Token, req.Email).Scan(&dbToken, &dbEmail)
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if dbToken != req.Token || dbEmail != req.Email {
		http.Error(w, "Invalid token or email", http.StatusBadRequest)
		return
	}

	addSendGridSubscription(req.Email)
}

func addSendGridSubscription(email string) {
	request := sendgrid.GetRequest(
		os.Getenv("SENDGRID_API_KEY"),
		"/v3/marketing/contacts",
		"https://api.sendgrid.com",
	)
	request.Method = "PUT"

	contactData := map[string]interface{}{
		"list_ids": []string{os.Getenv("SENDGRID_LIST_ID")},
		"contacts": []map[string]string{
			{
				"email": email,
			},
		},
	}

	jsonData, err := json.Marshal(contactData)
	if err != nil {
		log.Printf("Error marshaling contact data: %v", err)
		return
	}
	request.Body = jsonData
	response, err := sendgrid.API(request)
	if err != nil {
		log.Printf("Error adding contact to SendGrid: %v", err)
		return
	}

	if response.StatusCode >= 400 {
		log.Printf("SendGrid API error: %d %s", response.StatusCode, response.Body)
		return
	}
}

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	rl := NewRateLimiter()

	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		handleSignup(w, r, db, rl)
	})

	http.HandleFunc("/confirm", func(w http.ResponseWriter, r *http.Request) {
		handleConfirm(w, r, db)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
