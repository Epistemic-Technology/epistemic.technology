# Epistemic Technology Website & Services

This repository contains the complete codebase for the Epistemic Technology website and its supporting microservices. The system is designed as a modern, containerized web application with several backend services providing dynamic functionality to an otherwise static website.

## Project Overview

The Epistemic Technology website serves as the company's online presence, combining a static Hugo-generated site with dynamic interactive features:

- **Interactive Chatbot**: A RAG-powered AI chatbot that answers questions about Epistemic Technology
- **Contact Form**: A service that processes contact form submissions and sends them via email
- **Blog Subscription**: A service that manages blog signup and confirmation workflows

## Architecture

The project is built with a microservices architecture:

```
┌────────────────┐     ┌─────────────────┐     ┌────────────────┐
│                │     │                 │     │                │
│  Static Site   │────▶│  Contact Form   │────▶│    SendGrid    │
│    (Hugo)      │     │    Backend      │     │    Service     │
│                │     │                 │     │                │
└────────────────┘     └─────────────────┘     └────────────────┘
        │
        │               ┌─────────────────┐     ┌────────────────┐
        │               │                 │     │                │
        └──────────────▶│  Blog Signup    │────▶│    SendGrid    │
        │               │    Service      │     │    Service     │
        │               │                 │     │                │
        │               └─────────────────┘     └────────────────┘
        │
        │               ┌─────────────────┐     ┌────────────────┐     ┌────────────────┐
        │               │                 │     │                │     │                │
        └──────────────▶│    Chatbot      │────▶│    Chatbot     │────▶│     OpenAI     │
                        │   Frontend      │     │    Backend     │     │      API       │
                        │                 │     │                │     │                │
                        └─────────────────┘     └────────────────┘     └────────────────┘

```

## Components

The repository contains the following main components:

1. **[site](site/)**: Hugo-generated static website with the ETBASIC theme
2. **[contact-backend](contact-backend/)**: Go service for processing contact form submissions
3. **[blog-signup](blog-signup/)**: Go service for managing blog subscriptions
4. **[chatbot-frontend](chatbot-frontend/)**: SolidJS application providing a terminal-like chatbot interface
5. **[chatbot-backend](chatbot-backend/)**: Go service with RAG implementation for answering questions

## Getting Started

### Prerequisites

- Docker and Docker Compose
- SendGrid API key (for contact form and blog signup services)
- OpenAI API key (for chatbot functionality)

### Environment Setup

Create a `.env` file in the root directory with the following variables:

```
# Contact Form Backend
SENDGRID_API_KEY=your-sendgrid-api-key
CONTACT_EMAIL=recipient@example.com
CONTACT_SENDER_EMAIL=sender@example.com

# Blog Signup Service
SENDGRID_FROM_EMAIL=blog@example.com
SENDGRID_LIST_ID=your-sendgrid-list-id
SITE_URL=https://your-site.com
SUPPORT_EMAIL=support@example.com

# Chatbot Backend
OPENAI_API_KEY=your-openai-api-key
DATABASE_PATH=/app/data/chatbot.db
HUGO_CONTENT_PATH=/app/site/content
```

### Running the Application

The entire application can be run using Docker Compose:

```bash
docker-compose up
```

This will build and start all services according to the configuration in `compose.yaml`. Once running, the website will be available at http://localhost:8080.

### Development

For local development, each component can be run separately. Refer to the README in each component directory for specific development instructions.

## License

Code is [licensed](LICENSE.md) under an MIT license. Site content (blog posts, etc) is [licensed](CONTENT_LICENSE.md) under [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/).

## Contributing

This repository is made public for interest's sake and as a source of possible examples. It is not intended as a distributable application beyond the Epistemic Technology site. We are not currently accepting outside contributions to the site, though please contact us if you are interested in collaborating.

## Contact

Please [contact us](https://epistemic.technology/contact/) through the [Epistemic Technology](https://epistemic.technology/) website.
─
