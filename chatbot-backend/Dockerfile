FROM golang:1.24.1 AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y sqlite3 libsqlite3-dev

# Copy necessary files
WORKDIR /app
COPY chatbot-backend/ ./chatbot-backend/
COPY site/content/ ./site/content/

WORKDIR /app/chatbot-backend
RUN go mod download
RUN go build -o chatbot-backend .

# Final stage
FROM debian:bookworm-slim

# Install only the runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/chatbot-backend/chatbot-backend .
COPY --from=builder /app/site/content/ ./site/content/

RUN mkdir -p /db
ENV DATABASE_PATH="/db/site.db"
ENV PORT="8080"

EXPOSE 8080
CMD ["/app/chatbot-backend"]