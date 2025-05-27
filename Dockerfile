# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env.example ./.env
COPY --from=builder /app/migrations ./migrations

# Create a script to run migrations and start the app
RUN echo '#!/bin/sh\n\
echo "Running database migrations..."\n\
./main migrate\n\
echo "Starting application..."\n\
exec ./main' > /app/start.sh && chmod +x /app/start.sh

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/start.sh"] 