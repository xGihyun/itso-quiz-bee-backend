# Build stage
FROM golang:1.23.1-alpine as builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o itso-quiz-bee .

# Production stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS and postgresql client (for migrations if needed)
RUN apk --no-cache add ca-certificates postgresql-client

# Copy built binary from builder
COPY --from=builder /app/itso-quiz-bee .

# Copy migrations
COPY internal/database/migrations ./internal/database/migrations

# Copy .env file if needed
COPY .env.example .env

EXPOSE 3002

CMD ["./itso-quiz-bee"]
