FROM golang:1.22.5-alpine

# Install build dependencies and SQLite
RUN apk add --no-cache gcc musl-dev sqlite sqlite-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o /chatbot-backend ./cmd/cli/main.go

# Create directory for SQLite database
RUN mkdir -p /data

CMD ["/chatbot-backend"] 