# Epistemic Technology Contact Form Backend

This is a simple service that processes contact form submissions from the Epistemic Technology website and forwards them via email using either SendGrid or Mailgun.

## Architecture

The service is a Go application that provides a single REST endpoint for receiving contact form submissions. It supports two email providers: SendGrid and Mailgun, with automatic provider detection based on available API keys.

Endpoints:

- POST / -- Handles contact form submissions and sends the information to the configured recipient via the selected email provider.

## Environment Variables

### Required Environment Variables

- `CONTACT_EMAIL` - Email address that will receive contact form submissions
- `CONTACT_SENDER_EMAIL` - Email address used as the sender for contact form emails

### Email Provider Configuration

The service supports two email providers. You must configure at least one:

**Option 1: SendGrid**
- `SENDGRID_API_KEY` - API key for SendGrid service

**Option 2: Mailgun**
- `MAILGUN_API_KEY` - API key for Mailgun service
- `MAILGUN_DOMAIN` - Your Mailgun domain

### Optional Environment Variables

- `EMAIL_PROVIDER` - Explicitly set email provider ("sendgrid" or "mailgun"). If not set, the service will auto-detect based on available API keys, preferring SendGrid if both are available.
- `ALLOWED_ORIGIN` - Domain allowed for CORS (Cross-Origin Resource Sharing). Defaults to "http://localhost:1313" if not specified.
- `PORT` - Port on which the server will listen. Defaults to 8080 if not specified.

## Provider Selection Logic

1. If `EMAIL_PROVIDER` is explicitly set to "sendgrid" or "mailgun", that provider will be used.
2. If `EMAIL_PROVIDER` is not set, the service will:
   - Use SendGrid if `SENDGRID_API_KEY` is available
   - Use Mailgun if `MAILGUN_API_KEY` and `MAILGUN_DOMAIN` are available
   - Log an error if no provider is configured

## Running

The service is expected to be run through Docker.

### Using SendGrid:

```bash
docker build -t contact-backend .
docker run -p 8080:8080 \
  -e SENDGRID_API_KEY=your-sendgrid-api-key \
  -e CONTACT_EMAIL=your-recipient@example.com \
  -e CONTACT_SENDER_EMAIL=sender@example.com \
  -e ALLOWED_ORIGIN=https://your-website.com \
  contact-backend
```

### Using Mailgun:

```bash
docker build -t contact-backend .
docker run -p 8080:8080 \
  -e MAILGUN_API_KEY=your-mailgun-api-key \
  -e MAILGUN_DOMAIN=your-mailgun-domain.com \
  -e CONTACT_EMAIL=your-recipient@example.com \
  -e CONTACT_SENDER_EMAIL=sender@example.com \
  -e ALLOWED_ORIGIN=https://your-website.com \
  contact-backend
```

### Explicitly Setting Provider:

```bash
docker run -p 8080:8080 \
  -e EMAIL_PROVIDER=mailgun \
  -e MAILGUN_API_KEY=your-mailgun-api-key \
  -e MAILGUN_DOMAIN=your-mailgun-domain.com \
  -e CONTACT_EMAIL=your-recipient@example.com \
  -e CONTACT_SENDER_EMAIL=sender@example.com \
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

**Field Requirements:**
- `name` - Required
- `email` - Required  
- `subject` - Optional
- `message` - Required

Returns a JSON response on success:

```json
{
  "success": true,
  "message": "Contact submission received successfully"
}
```

Returns a JSON error response on failure:

```json
{
  "success": false,
  "message": "Error description"
}
```

### CORS Support

The service handles CORS preflight requests (OPTIONS) and sets appropriate headers based on the `ALLOWED_ORIGIN` environment variable.

## Email Features

- Emails are sent asynchronously to avoid blocking the API response
- The reply-to header is set to the original submitter's email address
- Comprehensive logging for debugging email delivery issues
- Fallback sender email generation for Mailgun (uses noreply@domain if CONTACT_SENDER_EMAIL not set)