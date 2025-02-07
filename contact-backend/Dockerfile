FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN go build -o contact-backend .

FROM scratch
COPY --from=builder /app/contact-backend /app/contact-backend
EXPOSE 8080
CMD ["/app/contact-backend"]