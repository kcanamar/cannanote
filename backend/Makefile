# CannaNote Backend Makefile
# Development-focused commands for local workflow

# ==============================================================================
# DEVELOPMENT COMMANDS
# ==============================================================================

# Primary development command - build assets and start server with hot reload
dev:
	@echo "🌱 Starting CannaNote development environment..."
	@$(MAKE) dev-setup dev-build
	@echo "🚀 Starting server with hot reload..."
	@if command -v air > /dev/null; then \
	    air; \
	else \
	    read -p "Go's 'air' is not installed. Install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/air-verse/air@latest; \
	        air; \
	    else \
	        echo "❌ Cannot run dev mode without air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

# Install development dependencies
dev-setup:
	@echo "Installing development dependencies..."
	@command -v air > /dev/null || go install github.com/air-verse/air@latest
	@command -v templ > /dev/null || go install github.com/a-h/templ/cmd/templ@latest
	@command -v golangci-lint > /dev/null || echo "⚠️  Consider installing golangci-lint for linting"
	@command -v govulncheck > /dev/null || echo "⚠️  Consider installing govulncheck: go install golang.org/x/vuln/cmd/govulncheck@latest"
	@command -v npx > /dev/null || echo "⚠️  npm not found - JS minification will be skipped"
	@echo "Development setup complete"

# Build all development assets (CSS, JS, templates)
dev-build:
	@echo "Building CSS..."
	@./tailwindcss -i ./cmd/web/assets/css/input.css -o ./cmd/web/assets/css/output.css --content "./cmd/**/*.templ"
	@echo "Generating templates..."
	@templ generate
	@echo "Assets built successfully"

# Run server without hot reload
run:
	@echo "Starting server..."
	@go run cmd/api/main.go

# Watch CSS for changes (run in separate terminal)
css-watch:
	@echo "👀 Watching CSS for changes..."
	@./tailwindcss -i ./cmd/web/assets/css/input.css -o ./cmd/web/assets/css/output.css --content "./cmd/**/*.templ" --watch

# Run tests
test:
	@echo "Running tests..."
	@go test ./tests -v

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test ./tests -v -cover -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Lint the code
lint:
	@echo "🔍 Linting code..."
	@golangci-lint run || echo "⚠️  Install golangci-lint for better linting"

# Security scan
security:
	@echo "Running security scan..."
	@govulncheck ./... || echo "⚠️  Install govulncheck: go install golang.org/x/vuln/cmd/govulncheck@latest"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f main
	@rm -f coverage.out coverage.html
	@echo "Cleaned"

# ==============================================================================
# PRODUCTION BUILD
# ==============================================================================

# Build for production
build:
	@echo "Building for production..."
	@./tailwindcss -i ./cmd/web/assets/css/input.css -o ./cmd/web/assets/css/output.css --content "./cmd/**/*.templ" --minify
	@echo "Minifying JavaScript..."
	@if command -v npx > /dev/null && npx terser --version > /dev/null 2>&1; then \
		npx terser cmd/web/assets/js/app.js -o cmd/web/assets/js/app.min.js -c -m; \
	else \
		echo "terser not available, copying unminified JS"; \
		cp cmd/web/assets/js/app.js cmd/web/assets/js/app.min.js; \
	fi
	@templ generate
	@go build -o main cmd/api/main.go
	@echo "Production build complete"

# Full pre-deploy checks
pre-deploy: clean lint security test build
	@echo "All pre-deploy checks passed"

# ==============================================================================
# DOCKER COMMANDS (if needed for local PostgreSQL)
# ==============================================================================

docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# ==============================================================================
# FLY.IO DEPLOYMENT (Production - NeonDB + Supertokens + Resend)
# ==============================================================================

fly-auth:
	@echo "Checking Fly.io authentication..."
	@flyctl auth whoami || (echo "❌ Please run: flyctl auth login" && exit 1)

fly-secrets:
	@echo "Setting up Fly.io secrets..."
	@test -f .env || (echo "❌ .env file not found!" && exit 1)
	@echo "Reading environment variables..."
	@DB_HOST=$$(grep '^DB_HOST=' .env | cut -d'=' -f2-); \
	 DB_PORT=$$(grep '^DB_PORT=' .env | cut -d'=' -f2-); \
	 DB_DATABASE=$$(grep '^DB_DATABASE=' .env | cut -d'=' -f2-); \
	 DB_USERNAME=$$(grep '^DB_USERNAME=' .env | cut -d'=' -f2-); \
	 DB_PASSWORD=$$(grep '^DB_PASSWORD=' .env | cut -d'=' -f2-); \
	 DB_SCHEMA=$$(grep '^DB_SCHEMA=' .env | cut -d'=' -f2-); \
	 SUPERTOKENS_CONNECTION_URI=$$(grep '^SUPERTOKENS_CONNECTION_URI=' .env | cut -d'=' -f2-); \
	 SUPERTOKENS_API_KEY=$$(grep '^SUPERTOKENS_API_KEY=' .env | cut -d'=' -f2-); \
	 APP_URL=$$(grep '^APP_URL=' .env | cut -d'=' -f2-); \
	 RESEND_API_KEY=$$(grep '^RESEND_API_KEY=' .env | cut -d'=' -f2-); \
	 if [ -z "$$DB_PASSWORD" ] || [ -z "$$SUPERTOKENS_CONNECTION_URI" ] || [ -z "$$SUPERTOKENS_API_KEY" ] || [ -z "$$APP_URL" ] || [ -z "$$RESEND_API_KEY" ]; then \
	   echo "❌ Missing required environment variables!"; \
	   echo "Required: DB_PASSWORD, SUPERTOKENS_CONNECTION_URI, SUPERTOKENS_API_KEY, APP_URL, RESEND_API_KEY"; \
	   exit 1; \
	 fi; \
	 echo "Setting secrets on Fly.io..."; \
	 flyctl secrets set \
	   "DB_HOST=$$DB_HOST" \
	   "DB_PORT=$$DB_PORT" \
	   "DB_DATABASE=$$DB_DATABASE" \
	   "DB_USERNAME=$$DB_USERNAME" \
	   "DB_PASSWORD=$$DB_PASSWORD" \
	   "DB_SCHEMA=$$DB_SCHEMA" \
	   "SUPERTOKENS_CONNECTION_URI=$$SUPERTOKENS_CONNECTION_URI" \
	   "SUPERTOKENS_API_KEY=$$SUPERTOKENS_API_KEY" \
	   "APP_ENV=production" \
	   "GIN_MODE=release" \
	   "DEBUG=false" \
	   "APP_URL=$$APP_URL" \
	   "RESEND_API_KEY=$$RESEND_API_KEY"
	@echo "All secrets configured"

fly-secrets-verify:
	@echo "Verifying secrets..."
	@flyctl secrets list | grep -q "DB_PASSWORD" || (echo "❌ DB_PASSWORD not set!" && exit 1)
	@flyctl secrets list | grep -q "SUPERTOKENS_CONNECTION_URI" || (echo "❌ SUPERTOKENS_CONNECTION_URI not set!" && exit 1)
	@flyctl secrets list | grep -q "SUPERTOKENS_API_KEY" || (echo "❌ SUPERTOKENS_API_KEY not set!" && exit 1)
	@flyctl secrets list | grep -q "APP_ENV" || (echo "❌ APP_ENV not set!" && exit 1)
	@flyctl secrets list | grep -q "APP_URL" || (echo "❌ APP_URL not set!" && exit 1)
	@flyctl secrets list | grep -q "RESEND_API_KEY" || (echo "❌ RESEND_API_KEY not set!" && exit 1)
	@echo "All required secrets verified"

fly-deploy:
	@echo "Deploying to Fly.io..."
	@echo "Using Fly.io remote builder..."
	@cd .. && flyctl deploy . --config backend/fly.toml --dockerfile backend/Dockerfile --remote-only -a cannanote-backend
	@echo "Deployment complete"

fly-status:
	@echo "Checking deployment status..."
	@flyctl status
	@echo "\n  Application URLs:"
	@flyctl info
	@echo "\n  Health check:"
	@curl -s https://cannanote-backend.fly.dev/health | jq . || echo "  Health check not ready yet"

fly-logs:
	@echo "Recent logs:"
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
	@echo "WARNING: Destroying machines (data will be deleted)..."
	@flyctl machines list
	@read -p "Type 'yes' to continue: " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		flyctl machines destroy --force; \
	else \
		echo "Cancelled."; \
	fi

# Complete deployment pipeline
deploy: fly-auth pre-deploy fly-secrets fly-secrets-verify fly-deploy fly-status
	@echo "Deployment pipeline complete"
	@echo "Live at: https://cannanote-backend.fly.dev"
	@echo "Health: https://cannanote-backend.fly.dev/health"

# Quick deploy (skip pre-deploy checks)
deploy-fast: fly-auth fly-secrets fly-secrets-verify fly-deploy fly-status
	@echo "⚡ Fast deployment complete"

# Rollback to previous version
rollback:
	@echo "Rolling back deployment..."
	@flyctl releases list
	@read -p "Enter release version to rollback to: " version; \
	flyctl releases rollback $$version

.PHONY: dev dev-setup dev-build run css-watch test test-coverage lint security clean build pre-deploy docker-run docker-down fly-auth fly-secrets fly-secrets-verify fly-deploy fly-status fly-logs fly-ssh fly-machines fly-stop fly-destroy deploy deploy-fast rollback
