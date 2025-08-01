# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -ldflags="-s -w" -o main cmd/api/main.go

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]