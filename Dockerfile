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

# Create start script
RUN echo '#!/bin/sh\n\
echo "Running database migrations..."\n\
./main migrate\n\
echo "Starting application..."\n\
exec ./main' > /app/start.sh && chmod +x /app/start.sh

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy the binary and start script from builder
COPY --from=builder /app/main .
COPY --from=builder /app/start.sh .
COPY --from=builder /app/.env.example ./.env
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 8080

# Run the application
CMD ["./start.sh"] 