# Epistemic Technology Contact Form Backend

This is a simple service that processes contact form submissions from the Epistemic Technology website and forwards them via email using SendGrid.

## Architecture

The service is a Go application that provides a single REST endpoint for receiving contact form submissions. It uses the SendGrid API to send emails to the website administrators when users submit the contact form.

Endpoints:

- POST / -- Handles contact form submissions and sends the information to the configured recipient via SendGrid.

## Environment Variables

The following environment variables are required for the application to function properly:

- `SENDGRID_API_KEY` - API key for SendGrid service used to send contact form emails
- `CONTACT_EMAIL` - Email address that will receive contact form submissions
- `CONTACT_SENDER_EMAIL` - Email address used as the sender for contact form emails
- `ALLOWED_ORIGIN` - Domain allowed for CORS (Cross-Origin Resource Sharing), typically your frontend domain

Optional environment variables:

- `PORT` - Port on which the server will listen (defaults to 8080 if not specified)

## Running

The service is expected to be run through Docker.

```bash
docker build -t contact-backend .
docker run -p 8080:8080 \
  -e SENDGRID_API_KEY=your-api-key \
  -e CONTACT_EMAIL=your-recipient@example.com \
  -e CONTACT_SENDER_EMAIL=sender@example.com \
  -e ALLOWED_ORIGIN=https://your-website.com \
  contact-backend
```

## API

### POST /

Accepts a JSON payload with the following structure:

```json
{
  "name": "User Name",
  "email": "user@example.com",
  "subject": "Subject Line",
  "message": "Message content from the user"
}
```

Returns a JSON response:

```json
{
  "success": true,
  "message": "Contact submission received successfully"
}
```

In case of errors, returns an appropriate HTTP status code and error message.
