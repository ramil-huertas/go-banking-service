# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies and tidy
RUN go mod download && go mod tidy

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata bash

# Copy the binary and other files from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env.example ./.env
COPY --from=builder /app/migrations ./migrations

# Create start script with explicit commands
RUN echo '#!/bin/bash' > start.sh && \
    echo 'set -e' >> start.sh && \
    echo 'echo "Running database migrations..."' >> start.sh && \
    echo './main migrate' >> start.sh && \
    echo 'echo "Starting application..."' >> start.sh && \
    echo 'exec ./main' >> start.sh && \
    chmod +x start.sh

# Expose port
EXPOSE 8080

# Run the application
CMD ["/bin/bash", "./start.sh"] 