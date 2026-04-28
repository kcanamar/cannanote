#!/bin/bash

# ============================================================================
# NeonDB Schema Import Script
# ============================================================================
# This script imports the CannaNote database schema and reference data
# into NeonDB in the correct order.
#
# Prerequisites:
# - PostgreSQL client (psql) installed
# - NeonDB connection details in .env file
# ============================================================================

set -e  # Exit on any error

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║         CannaNote NeonDB Schema Import Script                  ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""

# ============================================================================
# Load environment variables
# ============================================================================

if [ ! -f "backend/.env" ]; then
    echo -e "${RED}Error: backend/.env file not found!${NC}"
    echo "Please ensure you're running this from the project root directory."
    exit 1
fi

# Source the .env file
set -a
source backend/.env
set +a

echo -e "${BLUE}→${NC} Loaded environment variables from backend/.env"

# ============================================================================
# Construct connection string
# ============================================================================

# Check if DB_CONNECTION_URI is available (preferred)
if [ -n "$DB_CONNECTION_URI" ]; then
    CONN_STRING="$DB_CONNECTION_URI"
    echo -e "${GREEN}✓${NC} Using DB_CONNECTION_URI from .env"
else
    # Fallback to constructing from individual variables
    CONN_STRING="postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DATABASE}?sslmode=require"
    echo -e "${GREEN}✓${NC} Constructed connection string from DB_* variables"
fi

# ============================================================================
# Test connection
# ============================================================================

echo ""
echo -e "${BLUE}→${NC} Testing connection to NeonDB..."

if psql "$CONN_STRING" -c "SELECT version();" > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Successfully connected to NeonDB"
else
    echo -e "${RED}✗${NC} Failed to connect to NeonDB"
    echo "Please check your connection details in backend/.env"
    exit 1
fi

# ============================================================================
# Confirmation prompt
# ============================================================================

echo ""
echo -e "${YELLOW}⚠ WARNING:${NC} This will drop existing tables and recreate the schema."
echo -e "Database: ${BLUE}${DB_DATABASE}${NC}"
echo -e "Host: ${BLUE}${DB_HOST}${NC}"
echo ""
read -p "Do you want to continue? (yes/no): " -r
echo

if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo -e "${YELLOW}Import cancelled.${NC}"
    exit 0
fi

# ============================================================================
# Import schema
# ============================================================================

echo ""
echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  Step 1: Importing Main Schema                                 ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""

if [ ! -f "neondb_schema.sql" ]; then
    echo -e "${RED}✗${NC} File not found: neondb_schema.sql"
    exit 1
fi

echo -e "${BLUE}→${NC} Importing schema from neondb_schema.sql..."

if psql "$CONN_STRING" -f neondb_schema.sql > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Schema imported successfully"
else
    echo -e "${RED}✗${NC} Schema import failed"
    echo "Run manually to see errors: psql \"\$DB_CONNECTION_URI\" -f neondb_schema.sql"
    exit 1
fi

# ============================================================================
# Import cannabinoids reference data
# ============================================================================

echo ""
echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  Step 2: Importing Cannabinoids Reference Data                 ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""

if [ ! -f "neondb_reference_cannabinoids.sql" ]; then
    echo -e "${RED}✗${NC} File not found: neondb_reference_cannabinoids.sql"
    exit 1
fi

echo -e "${BLUE}→${NC} Importing cannabinoids data..."

if psql "$CONN_STRING" -f neondb_reference_cannabinoids.sql > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Cannabinoids reference data imported successfully"

    # Count records
    CANNABINOID_COUNT=$(psql "$CONN_STRING" -t -c "SELECT COUNT(*) FROM cannabinoids;")
    echo -e "${GREEN}✓${NC} Imported ${CANNABINOID_COUNT} cannabinoid records"
else
    echo -e "${RED}✗${NC} Cannabinoids import failed"
    exit 1
fi

# ============================================================================
# Import terpenes reference data
# ============================================================================

echo ""
echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  Step 3: Importing Terpenes Reference Data                     ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""

if [ ! -f "neondb_reference_terpenes.sql" ]; then
    echo -e "${RED}✗${NC} File not found: neondb_reference_terpenes.sql"
    exit 1
fi

echo -e "${BLUE}→${NC} Importing terpenes data..."

if psql "$CONN_STRING" -f neondb_reference_terpenes.sql > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Terpenes reference data imported successfully"

    # Count records
    TERPENE_COUNT=$(psql "$CONN_STRING" -t -c "SELECT COUNT(*) FROM terpenes;")
    echo -e "${GREEN}✓${NC} Imported ${TERPENE_COUNT} terpene records"
else
    echo -e "${RED}✗${NC} Terpenes import failed"
    exit 1
fi

# ============================================================================
# Verify import
# ============================================================================

echo ""
echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  Step 4: Verifying Database State                              ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""

echo -e "${BLUE}→${NC} Verifying all tables exist..."

TABLES=$(psql "$CONN_STRING" -t -c "
    SELECT table_name
    FROM information_schema.tables
    WHERE table_schema = 'public'
    ORDER BY table_name;
")

echo ""
echo -e "${GREEN}Tables created:${NC}"
echo "$TABLES" | while read -r table; do
    if [ -n "$table" ]; then
        ROW_COUNT=$(psql "$CONN_STRING" -t -c "SELECT COUNT(*) FROM $table;")
        echo -e "  ${GREEN}✓${NC} $table (${ROW_COUNT} rows)"
    fi
done

# ============================================================================
# Summary
# ============================================================================

echo ""
echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  Import Complete!                                              ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}✓${NC} NeonDB schema imported successfully"
echo -e "${GREEN}✓${NC} Reference data loaded"
echo -e "${GREEN}✓${NC} Database is ready for use"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "  1. Update backend/internal/database/database.go (remove cert path)"
echo "  2. Test backend connection: cd backend && go run cmd/api/main.go"
echo "  3. Verify health endpoint: curl http://localhost:3001/health"
echo ""