# CannaNote Project Makefile
# Tech Stack: Go + NeonDB (PostgreSQL) + Supertokens + Resend
# Development-focused commands for local workflow

.PHONY: help dev dev-setup test db-import db-verify clean deploy deploy-fast deploy-status deploy-logs deploy-rollback

# ==============================================================================
# QUICK START
# ==============================================================================

# Default target shows help
help:
	@echo "🌱 CannaNote Development Commands"
	@echo "=================================="
	@echo ""
	@echo "Quick Start:"
	@echo "  make dev         - Start development server with hot reload"
	@echo "  make dev-setup   - Install all development dependencies"
	@echo "  make test        - Run tests"
	@echo ""
	@echo "Database (NeonDB):"
	@echo "  make db-import   - Import schema and reference data to NeonDB"
	@echo "  make db-verify   - Verify NeonDB data integrity"
	@echo ""
	@echo "Deployment (Fly.io):"
	@echo "  make deploy      - Full deployment with pre-checks"
	@echo "  make deploy-fast - Quick deployment (skip checks)"
	@echo "  make deploy-status - Check deployment status"
	@echo "  make deploy-logs - View production logs"
	@echo "  make deploy-rollback - Rollback to previous version"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean       - Clean build artifacts"

# ==============================================================================
# DEVELOPMENT COMMANDS
# ==============================================================================

# Primary development command - start backend with hot reload
dev:
	@cd backend && make dev

# Install all development dependencies
dev-setup:
	@echo "🔧 Setting up CannaNote development environment..."
	@cd backend && make dev-setup
	@echo "✅ Development environment ready!"
	@echo "💡 Run 'make dev' to start coding"

# Run tests
test:
	@cd backend && make test

# Run tests with coverage
test-coverage:
	@cd backend && make test-coverage

# Clean build artifacts
clean:
	@cd backend && make clean

# ==============================================================================
# DATABASE COMMANDS (NeonDB)
# ==============================================================================

# Import schema and reference data to NeonDB
db-import:
	@echo "📊 Importing schema and reference data to NeonDB..."
	@test -f backend/.env || (echo "❌ backend/.env not found!" && exit 1)
	@cd backend && chmod +x ../scripts/import_to_neondb.sh
	@cd backend && ../scripts/import_to_neondb.sh
	@echo "✅ Database import complete"

# Verify NeonDB data integrity
db-verify:
	@echo "🔍 Verifying NeonDB data integrity..."
	@test -f backend/.env || (echo "❌ backend/.env not found!" && exit 1)
	@. backend/.env && psql "$$DB_CONNECTION_URI" -c "\
		SELECT 'profiles' as table_name, COUNT(*) as row_count FROM profiles \
		UNION ALL \
		SELECT 'cannabinoids', COUNT(*) FROM cannabinoids \
		UNION ALL \
		SELECT 'terpenes', COUNT(*) FROM terpenes \
		ORDER BY table_name;"
	@echo "✅ Verification complete"

# ==============================================================================
# DEPLOYMENT COMMANDS (Fly.io)
# ==============================================================================

# Full deployment pipeline with pre-checks
deploy:
	@echo "🚀 Starting CannaNote deployment pipeline..."
	@cd backend && make deploy
	@echo "✅ CannaNote deployment complete"

# Quick deployment (skip pre-checks)
deploy-fast:
	@echo "⚡ Fast deployment (skipping checks)..."
	@cd backend && make deploy-fast

# Check deployment status
deploy-status:
	@cd backend && make fly-status

# View production logs
deploy-logs:
	@cd backend && make fly-logs

# Rollback to previous version
deploy-rollback:
	@cd backend && make rollback