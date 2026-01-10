# CannaNote Project Makefile
# Manage database, backend, and deployment tasks

.PHONY: help db-local db-prod db-reset db-seed backend-dev backend-build deploy-fly

# Default target
help:
	@echo "CannaNote Development Commands"
	@echo "=============================="
	@echo "Database:"
	@echo "  db-local     - Seed local Supabase database"
	@echo "  db-prod      - Seed production Supabase database" 
	@echo "  db-reset     - Reset local database with fresh schema"
	@echo "  db-seed      - Load reference data only"
	@echo ""
	@echo "Backend:"
	@echo "  backend-dev  - Start Go backend development server"
	@echo "  backend-build- Build backend for production"
	@echo ""
	@echo "Deployment:"
	@echo "  deploy       - Full deployment pipeline with checks"
	@echo "  deploy-fast  - Quick deployment (skip pre-checks)"
	@echo "  deploy-status- Check deployment status and health"
	@echo "  deploy-logs  - View deployment logs"
	@echo "  deploy-rollback - Rollback to previous version"

# Database commands
db-local:
	@echo "Setting up LOCAL database..."
	@cd supabase && supabase db reset --debug
	@echo "Local database seeded successfully"
	@echo "Studio URL: http://localhost:54323"

db-prod:
	@echo "Setting up PRODUCTION database..."
	@cd supabase && supabase db push
	@echo "Loading reference data to production..."
	@echo "Enter your Supabase database password when prompted..."
	@cd supabase && PGPASSWORD="$$DB_PASSWORD" psql \
		postgresql://postgres.citdskdmralncvjyybin:$$DB_PASSWORD@db.citdskdmralncvjyybin.supabase.co:5432/postgres \
		--file=reference-data/cannabinoids.sql
	@cd supabase && PGPASSWORD="$$DB_PASSWORD" psql \
		postgresql://postgres.citdskdmralncvjyybin:$$DB_PASSWORD@db.citdskdmralncvjyybin.supabase.co:5432/postgres \
		--file=reference-data/terpenes.sql
	@echo "Production database seeded successfully"

db-reset:
	@cd supabase && supabase db reset --debug

db-seed:
	@echo "Loading reference data to local database..."
	@cd supabase && psql \
		--host=127.0.0.1 \
		--port=54322 \
		--username=postgres \
		--dbname=postgres \
		--file=reference-data/cannabinoids.sql
	@cd supabase && psql \
		--host=127.0.0.1 \
		--port=54322 \
		--username=postgres \
		--dbname=postgres \
		--file=reference-data/terpenes.sql
	@echo "Reference data loaded"

# Backend commands  
backend-dev:
	@cd backend && make dev

backend-build:
	@cd backend && make build

# Deployment commands
deploy:
	@echo "Starting CannaNote deployment pipeline..."
	@cd backend && make deploy
	@echo "CannaNote deployment complete"

deploy-fast:
	@echo "Fast deployment (skipping checks)..."
	@cd backend && make deploy-fast

deploy-status:
	@cd backend && make fly-status

deploy-logs:
	@cd backend && make fly-logs

deploy-rollback:
	@cd backend && make rollback

# Quick setup for new developers
setup: db-local
	@echo "Setup complete! Run 'make backend-dev' to start coding."