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
	@echo "Backend (Go):"
	@echo "  make dev           - Start backend with hot reload"
	@echo "  make dev-setup     - Install all dependencies"
	@echo "  make test          - Run backend tests"
	@echo ""
	@echo "Desktop (Rust/egui):"
	@echo "  make desktop-dev   - Run desktop app"
	@echo "  make desktop-build - Build release binary"
	@echo "  make desktop-test  - Run desktop tests"
	@echo ""
	@echo "Mobile (Flutter):"
	@echo "  make mobile-deps   - Get Flutter dependencies"
	@echo "  make mobile-dev    - Run mobile app"
	@echo "  make mobile-build-android - Build Android APK"
	@echo "  make mobile-build-ios     - Build iOS app"
	@echo ""
	@echo "Database (NeonDB):"
	@echo "  make db-import     - Import schema to NeonDB"
	@echo "  make db-verify     - Verify NeonDB data"
	@echo ""
	@echo "Deployment (Fly.io):"
	@echo "  make deploy        - Full deployment"
	@echo "  make deploy-fast   - Quick deployment"
	@echo ""
	@echo "Cross-Platform:"
	@echo "  make test-all      - Run all platform tests"
	@echo "  make clean-all     - Clean all build artifacts"

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

# ==============================================================================
# DESKTOP COMMANDS (Rust/egui)
# ==============================================================================

# Run desktop app in development mode
desktop-dev:
	@cd desktop && cargo run

# Build desktop release binary
desktop-build:
	@cd desktop && cargo build --release

# Run desktop tests
desktop-test:
	@cd desktop && cargo test

# Format and lint desktop code
desktop-lint:
	@cd desktop && cargo fmt && cargo clippy

# ==============================================================================
# MOBILE COMMANDS (Flutter)
# ==============================================================================

# Get mobile dependencies
mobile-deps:
	@cd mobile && flutter pub get

# Run mobile app in development mode
mobile-dev:
	@cd mobile && flutter run

# Build Android APK
mobile-build-android:
	@cd mobile && flutter build apk --release

# Build iOS app
mobile-build-ios:
	@cd mobile && flutter build ios --release

# Run mobile tests
mobile-test:
	@cd mobile && flutter test

# Generate code (Drift tables)
mobile-gen:
	@cd mobile && flutter pub run build_runner build --delete-conflicting-outputs

# Lint mobile code
mobile-lint:
	@cd mobile && flutter analyze

# ==============================================================================
# CROSS-PLATFORM COMMANDS
# ==============================================================================

# Run all tests across all platforms
test-all: test desktop-test mobile-test
	@echo "✅ All platform tests passed"

# Clean all build artifacts
clean-all: clean
	@cd desktop && cargo clean
	@cd mobile && flutter clean
	@echo "✅ All platforms cleaned"