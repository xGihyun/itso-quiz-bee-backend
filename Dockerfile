# Multi-stage build for Go application
# Stage 1: Build
FROM golang:1.23.1-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Production
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies and goose for migrations
RUN apk --no-cache add ca-certificates tzdata wget postgresql-client

# Install goose
RUN wget -O goose https://github.com/pressly/goose/releases/download/v3.17.0/goose_linux_x86_64 && \
    chmod +x goose

# Copy built binary from builder
COPY --from=builder /app/main .

# Copy any required files (migrations, etc.)
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations

# Copy scripts
COPY scripts ./scripts

# Copy entrypoint script
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh /app/scripts/*.sh

# Expose port
EXPOSE 3002

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3002/ || exit 1

# Run the application via entrypoint script
CMD ["/app/entrypoint.sh"]
