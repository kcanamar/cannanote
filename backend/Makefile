# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	@templ generate
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Create DB container
docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./tests -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/air-verse/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

# Development with hot reload
dev: watch

# Lint the code
lint:
	@echo "Linting..."
	@golangci-lint run || echo "Install golangci-lint for better linting"

# Security scan
security:
	@echo "Security scanning..."
	@govulncheck ./... || echo "Install govulncheck: go install golang.org/x/vuln/cmd/govulncheck@latest"

# Full pre-deploy checks
pre-deploy: clean lint security test build
	@echo "All pre-deploy checks passed"

# Fly.io deployment commands
fly-auth:
	@echo "Checking Fly.io authentication..."
	@flyctl auth whoami || (echo "ERROR: Please run: flyctl auth login" && exit 1)

fly-secrets:
	@echo "Setting up Fly.io secrets..."
	@echo "Checking .env file exists..."
	@test -f .env || (echo "ERROR: .env file not found!" && exit 1)
	@echo "Setting database configuration..."
	@DB_HOST=$$(grep '^DB_HOST=' .env | cut -d'=' -f2-); \
	 DB_PORT=$$(grep '^DB_PORT=' .env | cut -d'=' -f2-); \
	 DB_DATABASE=$$(grep '^DB_DATABASE=' .env | cut -d'=' -f2-); \
	 DB_USERNAME=$$(grep '^DB_USERNAME=' .env | cut -d'=' -f2-); \
	 DB_PASSWORD=$$(grep '^DB_PASSWORD=' .env | cut -d'=' -f2-); \
	 DB_SCHEMA=$$(grep '^DB_SCHEMA=' .env | cut -d'=' -f2-); \
	 SUPABASE_URL=$$(grep '^SUPABASE_URL=' .env | cut -d'=' -f2-); \
	 SUPABASE_ANON_KEY=$$(grep '^SUPABASE_ANON_KEY=' .env | cut -d'=' -f2-); \
	 if [ -z "$$DB_PASSWORD" ] || [ -z "$$SUPABASE_URL" ] || [ -z "$$SUPABASE_ANON_KEY" ]; then \
	   echo "ERROR: Missing required environment variables in .env file!"; \
	   echo "Required: DB_PASSWORD, SUPABASE_URL, SUPABASE_ANON_KEY"; \
	   exit 1; \
	 fi; \
	 echo "Setting secrets..."; \
	 flyctl secrets set \
	   "DB_HOST=$$DB_HOST" \
	   "DB_PORT=$$DB_PORT" \
	   "DB_DATABASE=$$DB_DATABASE" \
	   "DB_USERNAME=$$DB_USERNAME" \
	   "DB_PASSWORD=$$DB_PASSWORD" \
	   "DB_SCHEMA=$$DB_SCHEMA" \
	   "SUPABASE_URL=$$SUPABASE_URL" \
	   "SUPABASE_ANON_KEY=$$SUPABASE_ANON_KEY"
	@echo "All secrets configured successfully"

fly-secrets-verify:
	@echo "Verifying secrets are set..."
	@flyctl secrets list | grep -q "DB_PASSWORD" || (echo "ERROR: DB_PASSWORD not set!" && exit 1)
	@flyctl secrets list | grep -q "SUPABASE_URL" || (echo "ERROR: SUPABASE_URL not set!" && exit 1)
	@flyctl secrets list | grep -q "SUPABASE_ANON_KEY" || (echo "ERROR: SUPABASE_ANON_KEY not set!" && exit 1)
	@echo "All required secrets are set"

fly-deploy:
	@echo "Deploying to Fly.io..."
	@echo "Using Fly.io remote builder (no local Docker needed)..."
	@flyctl deploy --remote-only
	@echo "Deployment complete"

fly-status:
	@echo "Checking deployment status..."
	@flyctl status
	@echo "\nApplication URLs:"
	@flyctl info
	@echo "\nHealth check:"
	@curl -s https://cannanote-backend.fly.dev/health | jq . || echo "Health check endpoint not ready yet"

fly-logs:
	@echo "Showing recent logs..."
	@flyctl logs

fly-ssh:
	@echo "Opening SSH session..."
	@flyctl ssh console

fly-machines:
	@echo "Listing machines..."
	@flyctl machines list

fly-stop:
	@echo "Stopping all machines..."
	@flyctl scale count 0

fly-destroy:
	@echo "Destroying machines (WARNING: This will delete data)..."
	@flyctl machines list
	@read -p "Are you sure? Type 'yes' to continue: " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		flyctl machines destroy --force; \
	else \
		echo "Cancelled."; \
	fi

# Complete deployment pipeline
deploy: fly-auth pre-deploy fly-secrets fly-secrets-verify fly-deploy fly-status
	@echo "Deployment pipeline complete"
	@echo "Your app is live at: https://cannanote-backend.fly.dev"
	@echo "Health check: https://cannanote-backend.fly.dev/health"

# Quick deploy (skip checks - use with caution)
deploy-fast: fly-auth fly-secrets fly-secrets-verify fly-deploy fly-status
	@echo "Fast deployment complete"

# Rollback deployment
rollback:
	@echo "Rolling back deployment..."
	@flyctl releases list
	@read -p "Enter release version to rollback to: " version; \
	flyctl releases rollback $$version

.PHONY: all build run test clean watch dev lint security pre-deploy fly-auth fly-secrets fly-secrets-verify fly-deploy fly-status fly-logs fly-ssh fly-machines fly-stop fly-destroy deploy deploy-fast rollback
