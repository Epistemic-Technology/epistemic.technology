# Epistemic Technology Chatbot Frontend

This is a SolidJS application that simulates a terminal interface for interacting with a chatbot. It runs on top of the [Epistemic Technology static site](../site/README.md).

## Architecture

For a description of the architecture, see [About our chatbot](../site/content/about-our-chatbot.md). It uses Vite for building the application.

## Environment Variables

This application uses the following environment variables:

- `VITE_API_URL`: The URL for the chatbot API endpoint. Defaults to "http://localhost:8181/chat" if not specified.
- `VITE_FILEPATH_BASE_DIR`: Base directory path used for resolving file paths to URLs when referencing content files.

### How Vite Handles Environment Variables

Vite provides built-in support for environment variables through the `import.meta.env` object.

Only variables prefixed with `VITE_` are exposed to your Vite-processed code. This is a security measure to prevent accidentally exposing sensitive environment variables to client-side code.

To use environment variables in your local development:

1. Create a `.env` file in the root of the project
2. Add your environment variables using the format `VITE_VARIABLE_NAME=value`

Example `.env` file:

```
VITE_API_URL=http://localhost:8181/chat
VITE_FILEPATH_BASE_DIR=/path/to/content/files
```

## Running

For local development, the app can be launched on top of a blank page with `npm run start`. For production, the app is built as part of the [static site build process](../site/Dockerfile) and rendered on to the [#chatbot](../site/layouts/partials/chatbot.html) div.
