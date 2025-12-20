# CannaNote Backend Dockerfile
# Multi-stage build for production deployment on fly.io

# Build stage
FROM golang:1.24-alpine AS builder

# Install templ CLI for template generation
RUN go install github.com/a-h/templ/cmd/templ@latest

# Set working directory
WORKDIR /app

# Copy go mod files for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate templ templates
RUN templ generate

# Build the application with static linking for Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/api/main.go

# Production stage
FROM alpine:latest

# Install certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy web assets if needed (templates, static files)
COPY --from=builder /app/cmd/web ./cmd/web

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 3001

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3001/health || exit 1

# Run the application
CMD ["./main"]