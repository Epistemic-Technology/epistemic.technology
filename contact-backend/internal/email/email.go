package email

import (
	"log"
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

// SendEmail routes the email to the appropriate provider based on configuration
func SendEmail(submission ContactSubmission) {
	// First check if EMAIL_PROVIDER environment variable is set
	provider := os.Getenv("EMAIL_PROVIDER")

	if provider == "sendgrid" {
		sendSendGridEmail(submission)
		return
	}

	// If EMAIL_PROVIDER is not set, check for provider API keys
	if provider == "" {
		// Check for SendGrid API key
		if os.Getenv("SENDGRID_API_KEY") != "" {
			log.Printf("EMAIL_PROVIDER not set, but SENDGRID_API_KEY found - using SendGrid")
			sendSendGridEmail(submission)
			return
		}

		log.Printf("No email provider configured - set EMAIL_PROVIDER or provide API keys")
		return
	}

	log.Printf("Unsupported email provider: %s", provider)
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
