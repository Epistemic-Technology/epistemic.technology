package email

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
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

	if provider == "mailgun" {
		sendMailgunEmail(submission)
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

		// Check for Mailgun API key and domain
		if os.Getenv("MAILGUN_API_KEY") != "" && os.Getenv("MAILGUN_DOMAIN") != "" {
			log.Printf("EMAIL_PROVIDER not set, but MAILGUN_API_KEY and MAILGUN_DOMAIN found - using Mailgun")
			sendMailgunEmail(submission)
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

func sendMailgunEmail(submission ContactSubmission) {
	go func() {
		apiKey := os.Getenv("MAILGUN_API_KEY")
		domain := os.Getenv("MAILGUN_DOMAIN")

		if apiKey == "" {
			log.Printf("MAILGUN_API_KEY is not set")
			return
		}

		if domain == "" {
			log.Printf("MAILGUN_DOMAIN is not set")
			return
		}

		mg := mailgun.NewMailgun(domain, apiKey)

		sender := os.Getenv("CONTACT_SENDER_EMAIL")
		if sender == "" {
			sender = "noreply@" + domain
		}

		recipient := os.Getenv("CONTACT_EMAIL")
		if recipient == "" {
			log.Printf("CONTACT_EMAIL is not set")
			return
		}

		message := mailgun.NewMessage(
			submission.Name+" <"+sender+">",
			submission.Subject,
			submission.Message,
			recipient,
		)

		// Set reply-to to the original submitter
		message.SetReplyTo(submission.Email)

		log.Printf(
			"Sending contact email from %s <%s> to %s with subject %s via Mailgun",
			submission.Name,
			submission.Email,
			recipient,
			submission.Subject,
		)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		resp, id, err := mg.Send(ctx, message)
		if err != nil {
			log.Printf("Error sending email via Mailgun: %v", err)
			return
		}

		log.Printf("Email sent successfully via Mailgun - ID: %s, Response: %s", id, resp)
	}()
}
