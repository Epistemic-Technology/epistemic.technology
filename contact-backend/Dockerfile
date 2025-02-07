FROM golang:latest AS builder
WORKDIR /app
COPY . .
# Create a fully static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o contact-backend .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/contact-backend /contact-backend
EXPOSE ${PORT:-8080}
CMD ["/contact-backend"]