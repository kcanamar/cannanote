# Docker Containerization Strategy

## Overview
This document outlines the Docker containerization approach for CannaNote 2.0, designed for HIPAA compliance and enterprise deployment.

## Container Architecture

### Production Container (`Dockerfile.prod`)

```dockerfile
# Multi-stage build for security and size optimization
FROM golang:1.21-alpine AS builder

# Security: Run as non-root user
RUN adduser -D -s /bin/sh appuser

WORKDIR /app

# Copy dependency files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o cannanote ./cmd/server

# Production stage - use distroless for security
FROM gcr.io/distroless/static-debian11

# Copy user from builder
COPY --from=builder /etc/passwd /etc/passwd

# Copy binary
COPY --from=builder /app/cannanote /cannanote

# Copy static assets if needed
COPY --from=builder /app/static /static
COPY --from=builder /app/templates /templates

# Security: Run as non-root
USER appuser:appuser

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/cannanote", "-health-check"]

# Expose port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/cannanote"]
```

### Development Container (`Dockerfile.dev`)

```dockerfile
FROM golang:1.21-alpine

# Install development tools
RUN apk add --no-cache git curl

# Install Air for hot reloading
RUN go install github.com/cosmtrek/air@latest

# Create app user
RUN adduser -D -s /bin/sh appuser

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Change ownership to appuser
RUN chown -R appuser:appuser /app
USER appuser

# Expose port for development
EXPOSE 8080

# Default command runs Air for hot reloading
CMD ["air", "-c", ".air.toml"]
```

### Docker Compose Configuration

```yaml
# docker-compose.yml
version: '3.8'

services:
  app-dev:
    build:
      context: .
      dockerfile: Dockerfile.dev
    ports:
      - "8080:8080"
    volumes:
      - .:/app
      - go-mod-cache:/go/pkg/mod
    environment:
      - ENVIRONMENT=development
      - SUPABASE_URL=${SUPABASE_URL}
      - SUPABASE_ANON_KEY=${SUPABASE_ANON_KEY}
    depends_on:
      - prometheus
      - grafana
    networks:
      - cannanote-dev

  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./config/prometheus.yml:/etc/prometheus/prometheus.yml
    networks:
      - cannanote-dev

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-storage:/var/lib/grafana
    networks:
      - cannanote-dev

volumes:
  go-mod-cache:
  grafana-storage:

networks:
  cannanote-dev:
    driver: bridge
```

## Security Considerations

### Container Security Best Practices
- **Distroless base image**: Minimal attack surface
- **Non-root user**: Security principle of least privilege
- **Multi-stage builds**: Smaller production images
- **Health checks**: Monitoring and auto-recovery
- **Read-only filesystem**: Immutable container design

### HIPAA Compliance Features
- **Encrypted secrets**: OpenBao integration for runtime secrets
- **Audit logging**: All container events logged to Sentry
- **Network isolation**: Secure container networking
- **Access controls**: Container-level permission management

## Deployment Strategy

### Development Environment
```bash
# Start development environment
docker-compose up -d

# View logs
docker-compose logs -f app-dev

# Execute commands in container
docker-compose exec app-dev go test ./...
```

### Production Deployment (Fly.io)
```bash
# Build and deploy to Fly.io
fly deploy --dockerfile Dockerfile.prod

# Scale application
fly scale count 2

# View application logs
fly logs
```

## Environment Variables

### Required Environment Variables
```env
# Application
ENVIRONMENT=production
PORT=8080
SECRET_KEY=${SECRET_KEY}

# Supabase
SUPABASE_URL=${SUPABASE_URL}
SUPABASE_SERVICE_KEY=${SUPABASE_SERVICE_KEY}

# Observability
SENTRY_DSN=${SENTRY_DSN}
PROMETHEUS_ENDPOINT=${PROMETHEUS_ENDPOINT}

# OpenBao
OPENBAO_TOKEN=${OPENBAO_TOKEN}
OPENBAO_URL=${OPENBAO_URL}
```

## Monitoring Integration

### Health Check Endpoint
```go
func healthCheck(w http.ResponseWriter, r *http.Request) {
    status := map[string]string{
        "status": "healthy",
        "version": version,
        "timestamp": time.Now().UTC().Format(time.RFC3339),
    }
    
    // Check database connection
    if err := db.Ping(); err != nil {
        status["status"] = "unhealthy"
        status["database"] = "disconnected"
        w.WriteHeader(http.StatusServiceUnavailable)
    }
    
    json.NewEncoder(w).Encode(status)
}
```

### Metrics Collection
- **Prometheus metrics**: Application performance metrics
- **Sentry monitoring**: Error tracking and performance monitoring
- **Container metrics**: CPU, memory, network usage

## Compliance Documentation

### Audit Requirements
- Container build logs retained for 6 years
- Deployment artifacts signed and verified
- Security scanning results documented
- Access logs for all container operations

### Change Management
- All container changes require approval
- Automated testing before deployment
- Rollback procedures documented
- Security review for all updates

## Implementation Timeline

**Week 1:**
- [ ] Create Dockerfiles (dev/prod)
- [ ] Set up docker-compose for development
- [ ] Configure health checks

**Week 2:**
- [ ] Integrate with Fly.io deployment
- [ ] Set up monitoring and logging
- [ ] Security testing and optimization

**Week 3:**
- [ ] Documentation and training
- [ ] Compliance verification
- [ ] Production readiness review