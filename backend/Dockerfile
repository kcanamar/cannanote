# CannaNote Backend Dockerfile
# Multi-stage build for production deployment on fly.io

# Build stage
FROM golang:1.24-alpine AS builder

# Install templ CLI for template generation
RUN go install github.com/a-h/templ/cmd/templ@latest

# Set working directory
WORKDIR /app

# Copy go mod files for dependency caching
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source code
COPY backend/ .

# Copy docs from root directory (accessible with root build context)
COPY docs/ ./docs/

# Generate templ templates
RUN templ generate

# Build the application with static linking for Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/api/main.go

# Production stage
FROM alpine:latest

# Install certificates for HTTPS/TLS requests + timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy web assets if needed (templates, static files)
COPY --from=builder /app/cmd/web ./cmd/web

# Copy documentation content
COPY --from=builder /app/docs ./docs

# Copy the Supabase CA cert (critical for verify-full SSL mode)
COPY --from=builder /app/certs ./certs

# Change ownership to non-root user (including the new certs dir)
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port (update if your app uses a different one)
EXPOSE 3001

# Run the application
CMD ["./main"]