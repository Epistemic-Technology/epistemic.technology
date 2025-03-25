# Epistemic Technology Chatbot Backend

This is the backend for the Epistemic Technology chatbot. It is a Go application that uses RAG to answer questions about the Epistemic Technology website.

## Architecture

For an explanation of the architecture, see [About our Chatbot](../site/content/about-our-chatbot.md). This service works in conjunction with the [chatbot-frontend](../chatbot-frontend/README.md).

## Environment Variables

The following environment variables are required for the application to function properly:

- `OPENAI_API_KEY` - API key for OpenAI services used for generating embeddings and LLM responses
- `DATABASE_PATH` - Path to the SQLite database file where document embeddings are stored
- `HUGO_CONTENT_PATH` - Path to the Hugo content directory containing the website content to be embedded
- `PORT` - Port on which the server will listen (e.g., "8181")

These environment variables can be set in a `.env` file in the project root directory, or they can be provided as command-line flags when starting the application:

- `--api-key` - Overrides the OPENAI_API_KEY environment variable
- `--db` - Overrides the DATABASE_PATH environment variable
- `--hugo-content-path` - Overrides the HUGO_CONTENT_PATH environment variable
- `--port` - Overrides the PORT environment variable

## CLI

The chatbot backend provides two command-line interfaces:

### Document Management CLI

The `cli.go` tool helps manage the document database:

```
Usage:
  cli embed-hugo-directory <directory> [--recursive] [--db=<path>]
  cli embed-hugo-file <file> [--db=<path>]
  cli list-documents [--db=<path>]
  cli get-document-details <document_id> [--db=<path>]
  cli get-db-stats [--db=<path>]
  cli list-chunks [--db=<path>]
  cli list-chunks-for-document <document_id> [--db=<path>]
```

This tool allows you to:

- Process and embed Hugo content files into the database
- Inspect documents and chunks stored in the database
- View database statistics

### Interactive Chat CLI

The `chat.go` tool provides a simple interactive chat interface:

```
Usage:
  chat [--db=<path>] [--api-key=<key>]
```

This tool:

- Provides an interactive prompt for chatting with the bot
- Maintains conversation history for context
- Displays sources used in responses
- Supports special commands (`exit`, `quit`, `clear`)

## Running the API

The service is expected to be run through Docker. The main entrypoint is [main.go](main.go). It starts the API and embeds documents into the database. If the database exists (such as through a volume mount), it will avoid duplicating documents by computing their hash.
