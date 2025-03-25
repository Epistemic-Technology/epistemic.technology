# Epistemic Technology Blog Signup Service

This is a simple service that allows users to sign up for the Epistemic Technology blog.

## Architecture

The service is a Go application that uses a SQLite database to temporarily store email addresses and identification tokens. It uses the SendGrid API to send emails to users.

Endpoints:

- POST /signup -- Handles submission of subscription form submissions and sends a confirmation email to the user.
- POST /confirm -- Handles confirmation of the user's email address when user clicks the link in the confirmation email, and adds the user to the SendGrid mailing list.

## Environment Variables

The following environment variables are required for the application to function properly:

- `SENDGRID_API_KEY` - API key for SendGrid service used to send confirmation emails and add users to mailing lists
- `SENDGRID_FROM_EMAIL` - Email address used as the sender for confirmation emails
- `SENDGRID_LIST_ID` - ID of the SendGrid mailing list to add confirmed subscribers to
- `SITE_URL` - Base URL of the website, used to construct confirmation links (e.g., "https://epistemic.technology")
- `SUPPORT_EMAIL` - Email address for user support inquiries, included in confirmation emails
- `ALLOWED_ORIGIN` - Domain allowed for CORS (Cross-Origin Resource Sharing), typically your frontend domain

Optional environment variables:

- `PORT` - Port on which the server will listen (defaults to 8080 if not specified)
- `GO_DEBUG` - When set to "true", allows bypass of the rate limiting for repeated subscription attempts (for development/testing)

## Running

The expected way to run the service is through Docker.
