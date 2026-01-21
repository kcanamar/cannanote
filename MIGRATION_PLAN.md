# Migration Plan: Supabase → NeonDB + Supertokens + Resend

## Executive Summary

**Migration Goal**: Replace Supabase (auth + database) with NeonDB (database), Supertokens Managed Cloud (auth), and Resend (email - already implemented).

**Timeline**: 11-16 hours (1-2 focused development days)

**Risk Level**: LOW-MEDIUM
- Minimal Supabase usage (only beta signup)
- Hexagonal architecture designed for this swap
- No existing user sessions to migrate
- Database already uses vendor-neutral pgx driver

**Key Decision**: Using Supertokens Managed Cloud (not self-hosted) for faster time-to-market and zero operational burden during beta phase.

---

## Current State Analysis

### What We're Actually Using Supabase For

1. **Auth**: Only `AdminCreateUser` for beta signup (`backend/internal/server/routes.go:508-569`)
2. **Database**: Direct PostgreSQL via pgx/v5 (vendor-neutral, NOT using Supabase SDK)
3. **REST API**: Single endpoint fetching cannabinoid data (`backend/cmd/web/learn_handlers.go`)

### What We're NOT Using

- No user login/session management
- No protected routes or JWT middleware
- No Realtime, Storage, Edge Functions
- No complex auth flows

**This makes migration exceptionally clean.**

---

## Reusing Existing SQL Schema Files

### Available SQL Files

You have well-structured SQL files that can be reused:

- `supabase/seeds/seed.sql` - Complete application schema (humans, cannabinoids, terpenes, entries, etc.)
- `supabase/migrations/20260107210113_beta_profiles_table.sql` - Beta profiles table
- `supabase/reference-data/cannabinoids.sql` - Cannabinoid reference data
- `supabase/reference-data/terpenes.sql` - Terpene reference data

### Required Modifications for NeonDB

**The Issue**: These files contain Supabase-specific features:
1. Foreign keys to `auth.users(id)` (doesn't exist in plain PostgreSQL)
2. RLS policies using `auth.uid()` (Supabase function not available in NeonDB)

**The Solution**:

1. **Copy `supabase/seeds/seed.sql` to `neondb_schema.sql`**
2. **Remove all RLS policies** (lines 43-58, 272-284):
   - Delete all `CREATE POLICY` statements
   - Delete all `ALTER TABLE ... ENABLE ROW LEVEL SECURITY` statements
   - Supertokens uses session-based auth, not database-level RLS

3. **Keep all table definitions unchanged** - Everything else stays the same

4. **Append profiles table** (modified from migration file):

```sql
-- ============================================================================
-- PROFILES TABLE - Beta user profiles linked to Supertokens users
-- ============================================================================

CREATE TABLE IF NOT EXISTS profiles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT UNIQUE NOT NULL,
  beta_status TEXT DEFAULT 'waitlist',
  beta_joined_at TIMESTAMPTZ DEFAULT NOW(),
  subscription_tier TEXT DEFAULT 'beta_grandfathered',
  referral_code TEXT UNIQUE,
  invited_at TIMESTAMPTZ,
  activated_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index for fast email lookups
CREATE INDEX idx_profiles_email ON profiles(email);
CREATE INDEX idx_profiles_beta_status ON profiles(beta_status);
```

**Note**: Removed `REFERENCES auth.users(id)` since NeonDB doesn't have Supabase auth tables.

5. **Reference data files are ready as-is**:
   - `supabase/reference-data/cannabinoids.sql` - No changes needed
   - `supabase/reference-data/terpenes.sql` - No changes needed

---

## Environment Variables Required

### New .env Configuration

```bash
# ============================================
# APPLICATION CONFIGURATION
# ============================================
APP_ENV=production
APP_URL=https://cannanote.org
PORT=8080
GIN_MODE=release
DEBUG=false

# ============================================
# NEONDB POSTGRESQL CONFIGURATION
# ============================================
# Get from: https://console.neon.tech/app/projects
DB_HOST=ep-your-project-id.us-east-2.aws.neon.tech
DB_PORT=5432
DB_DATABASE=cannanote
DB_USERNAME=neondb_owner
DB_PASSWORD=<from-neon-dashboard>
DB_SCHEMA=public

# Note: NeonDB requires sslmode=verify-full (uses system CA certs)

# ============================================
# SUPERTOKENS MANAGED CLOUD CONFIGURATION
# ============================================
# Get from: https://supertokens.com/dashboard
SUPERTOKENS_CONNECTION_URI=https://dev-xxxxx.us.supertokens.io
SUPERTOKENS_API_KEY=<from-supertokens-dashboard>

# ============================================
# RESEND EMAIL CONFIGURATION
# ============================================
RESEND_API_KEY=re_J4cvQmF1_8HhtDxBLtZqVCET9SeJaqJQ1
# ^ Already configured, no changes needed

# ============================================
# RATE LIMITING
# ============================================
RATE_LIMIT_MAX_ATTEMPTS=3
RATE_LIMIT_WINDOW_MINUTES=15
```

### Variables to DELETE

```bash
# Remove from .env:
SUPABASE_URL
SUPABASE_ANON_KEY
SUPABASE_PUBLISH_KEY
SUPABASE_SECRET_KEY
```

---

## Required Documentation

### Supertokens Go SDK
- **Main Setup**: https://supertokens.com/docs/session/quick-setup/backend
- **Email/Password Recipe**: https://supertokens.com/docs/emailpassword/quickstart/introduction
- **API Reference**: https://supertokens.com/docs/go/api-reference/middleware/supertokens-middleware
- **Session Verification**: https://supertokens.com/docs/passwordless/common-customizations/sessions/session-verification-in-api/verify-session
- **GitHub**: https://github.com/supertokens/supertokens-golang
- **Go Docs**: https://pkg.go.dev/github.com/supertokens/supertokens-golang

### NeonDB
- **Secure Connection**: https://neon.tech/docs/connect/connect-securely
- **Connection Format**: https://neon.tech/docs/connect/connect-from-any-app
- **Troubleshooting**: https://neon.tech/docs/connect/connection-errors

---

## Migration Steps (Ordered)

### Phase 1: Database Migration (1-2 hours)

**Goal**: Move PostgreSQL data from Supabase to NeonDB

#### Step 1.1: Create NeonDB Project
```bash
# 1. Go to https://console.neon.tech
# 2. Create project: "CannaNote Production"
# 3. Select region: us-east-2 (or closest to backend)
# 4. Copy connection string from dashboard
```

#### Step 1.2: Prepare NeonDB Schema

```bash
# Copy seed schema
cp supabase/seeds/seed.sql neondb_schema.sql

# Edit neondb_schema.sql:
# 1. Remove all lines containing "CREATE POLICY"
# 2. Remove all lines containing "ENABLE ROW LEVEL SECURITY"
# 3. Append profiles table definition (see above)
```

#### Step 1.3: Import to NeonDB

```bash
# Import main schema (modified seed file)
psql "postgresql://<user>:<pass>@ep-xxxxx.us-east-2.aws.neon.tech:5432/cannanote?sslmode=verify-full" \
  -f neondb_schema.sql

# Import reference data (no modifications needed)
psql "postgresql://..." -f supabase/reference-data/cannabinoids.sql
psql "postgresql://..." -f supabase/reference-data/terpenes.sql
```

#### Step 1.4: Export/Import User Data (if beta users exist)

```bash
# Check if you have existing beta users in Supabase
psql "postgresql://postgres:VUNotEXZ1EDAMt11@db.citdskdmralncvjyybin.supabase.co:5432/postgres" \
  -c "SELECT COUNT(*) FROM profiles;"

# If count > 0, export the data
pg_dump -h db.citdskdmralncvjyybin.supabase.co \
  -U postgres \
  -a \
  -t profiles \
  -f profiles_data.sql \
  postgres

# Import to NeonDB
psql "postgresql://..." -f profiles_data.sql
```

#### Step 1.5: Verify Data Integrity

```sql
-- Compare row counts
SELECT 'profiles' as table_name, COUNT(*) FROM profiles;
SELECT 'cannabinoids' as table_name, COUNT(*) FROM cannabinoids;
SELECT 'terpenes' as table_name, COUNT(*) FROM terpenes;
SELECT 'humans' as table_name, COUNT(*) FROM humans;
SELECT 'entries' as table_name, COUNT(*) FROM entries;
```

---

### Phase 2: Update Database Connection (30 min)

**Goal**: Point backend to NeonDB instead of Supabase

#### Step 2.1: Update database.go

**File**: `/backend/internal/database/database.go` (lines 75-82)

**Change**:
```go
// BEFORE (Supabase with cert):
connStr := fmt.Sprintf(
    "postgresql://%s:%s@%s:%s/%s?sslmode=verify-full&sslrootcert=%s&search_path=%s",
    username, password, host, port, database, certPath, schema,
)

// AFTER (NeonDB uses system CA):
connStr := fmt.Sprintf(
    "postgresql://%s:%s@%s:%s/%s?sslmode=verify-full&search_path=%s",
    username, password, host, port, database, schema,
)
```

#### Step 2.2: Update .env

```bash
# Update these values in /backend/.env
DB_HOST=ep-your-neondb-project.us-east-2.aws.neon.tech
DB_USERNAME=neondb_owner
DB_PASSWORD=<from-neon-dashboard>
DB_DATABASE=cannanote
```

#### Step 2.3: Test Connection

```bash
cd backend
go run cmd/api/main.go
# Should see: "Database connection opened successfully"

# Test health endpoint
curl http://localhost:8080/health
```

---

### Phase 3: Supertokens Integration (3-4 hours)

**Goal**: Replace Supabase Auth with Supertokens Managed Cloud

#### Step 3.1: Create Supertokens Account

```bash
# 1. Go to https://supertokens.com/dashboard
# 2. Sign up and create app: "CannaNote Beta"
# 3. Select region: US
# 4. Copy connectionURI and apiKey
# 5. Add to .env:
SUPERTOKENS_CONNECTION_URI=https://dev-xxxxx.us.supertokens.io
SUPERTOKENS_API_KEY=<from-dashboard>
```

#### Step 3.2: Add Dependencies

```bash
cd backend
go get github.com/supertokens/supertokens-golang@latest
go mod tidy
```

#### Step 3.3: Create Supertokens Auth Service

**Create new file**: `/backend/internal/adapters/external/supertokens_auth_service.go`

```go
package external

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"os"
)

// InitializeSupertokens configures Supertokens with managed cloud
func InitializeSupertokens() error {
	apiBasePath := "/auth"
	websiteBasePath := "/auth"

	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: os.Getenv("SUPERTOKENS_CONNECTION_URI"),
			APIKey:        os.Getenv("SUPERTOKENS_API_KEY"),
		},
		AppInfo: supertokens.AppInfo{
			AppName:         "CannaNote",
			APIDomain:       os.Getenv("APP_URL"),
			WebsiteDomain:   os.Getenv("APP_URL"),
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &websiteBasePath,
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(&epmodels.TypeInput{}),
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeRequired,
			}),
			session.Init(&sessmodels.TypeInput{
				CookieSecure: func() *bool {
					secure := os.Getenv("APP_ENV") == "production"
					return &secure
				}(),
			}),
		},
	})

	return err
}

// NewSupertokensAuthService creates auth service adapter
func NewSupertokensAuthService() *SupertokensAuthService {
	return &SupertokensAuthService{}
}

type SupertokensAuthService struct{}

// Implement ports.AuthService interface methods
// (GenerateToken, ValidateToken, RevokeToken)
// Note: Supertokens handles most of this via middleware
```

#### Step 3.4: Create Middleware Helpers

**Create new file**: `/backend/internal/adapters/http/supertokens_middleware.go`

```go
package http

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// SupertokensCORS returns CORS middleware for Supertokens
func SupertokensCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		supertokens.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		})).ServeHTTP(c.Writer, c.Request)
	}
}

// VerifySession returns middleware requiring valid session
func VerifySession(options *sessmodels.VerifySessionOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		session.VerifySession(options, func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		})(c.Writer, c.Request)
	}
}
```

#### Step 3.5: Wire Dependencies in server.go

**File**: `/backend/internal/server/server.go` (lines 35-60)

**Changes**:
```go
func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	// INITIALIZE SUPERTOKENS FIRST
	if err := external.InitializeSupertokens(); err != nil {
		log.Fatalf("Failed to initialize Supertokens: %v", err)
	}

	dbService := database.New()
	rawDB := dbService.GetDB()

	// Wire dependencies
	humanRepo := repository.NewSupabaseHumanRepository(rawDB)

	// ADD AUTH SERVICE
	authService := external.NewSupertokensAuthService()

	// Pass authService to HumanService (3rd parameter)
	humanService := application.NewHumanService(humanRepo, nil, authService)

	humanHandlers := httpAdapters.NewHumanHandlers(humanService)

	newServer := &Server{
		port:          port,
		db:            dbService,
		humanService:  humanService,
		humanHandlers: humanHandlers,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
```

---

### Phase 4: Rewrite Beta Signup (2-3 hours)

**Goal**: Replace Supabase AdminCreateUser with Supertokens signup

#### Step 4.1: Update routes.go Imports

**File**: `/backend/internal/server/routes.go`

**Add imports**:
```go
import (
	"crypto/rand"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	httpAdapters "backend/internal/adapters/http"
)
```

#### Step 4.2: Add Supertokens Middleware to Router

**File**: `/backend/internal/server/routes.go` (in `RegisterRoutes` function)

**Add after SecurityHeaders**:
```go
func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(SecurityHeaders())
	r.Use(RequestSizeLimit(8 << 20))

	// ADD SUPERTOKENS MIDDLEWARE
	r.Use(httpAdapters.SupertokensCORS())

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// ... rest of routes
}
```

#### Step 4.3: Replace createBetaUser Function

**File**: `/backend/internal/server/routes.go` (lines 508-569)

**Replace entire function**:
```go
func (s *Server) createBetaUserSupertokens(email string) error {
	// Generate secure temporary password
	tempPassword := generateSecurePassword()

	// Create user in Supertokens
	signUpResponse, err := emailpassword.SignUp("public", email, tempPassword)
	if err != nil {
		return fmt.Errorf("failed to create Supertokens user: %w", err)
	}

	if signUpResponse.EmailAlreadyExistsError != nil {
		return fmt.Errorf("email %s already registered", email)
	}

	userID := signUpResponse.OK.User.ID

	// Create profile in NeonDB
	db := s.db.GetDB()
	_, err = db.Exec(`
		INSERT INTO profiles (id, email, beta_status, beta_joined_at, created_at)
		VALUES ($1, $2, 'waitlist', NOW(), NOW())
	`, userID, email)

	if err != nil {
		conditionalLog("ERROR", "Failed to create profile: %v", err)
		return fmt.Errorf("failed to create user profile: %w", err)
	}

	// Send welcome email (non-blocking)
	go func() {
		resendService := NewResendService()
		if err := resendService.SendWelcomeEmail(email, ""); err != nil {
			conditionalLog("ERROR", "Failed to send welcome email: %v", err)
		}
	}()

	conditionalLog("INFO", "Beta user created successfully: %s", email)
	return nil
}

func generateSecurePassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	b := make([]byte, 32)
	rand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}
```

#### Step 4.4: Update BetaSignupHandler Call

**File**: `/backend/internal/server/routes.go` (line ~489)

**Change**:
```go
// BEFORE
if err := s.createBetaUser(email); err != nil {

// AFTER
if err := s.createBetaUserSupertokens(email); err != nil {
```

#### Step 4.5: Update Resend Welcome Email

**File**: `/backend/internal/adapters/external/resend_service.go` (lines 76-77)

**Change message**:
```go
// OLD:
<li><strong>Verify your email</strong> - Check your inbox for the Supabase verification email</li>

// NEW:
<li><strong>Verify your email</strong> - Check your inbox for the verification email</li>
<li><strong>Set your password</strong> - Use the password reset link to secure your account</li>
```

---

### Phase 5: Testing & Validation (2-3 hours)

#### Step 5.1: Local Testing Checklist

```bash
# 1. Test database connection
cd backend
go run cmd/api/main.go
# Should see: "Database connection opened successfully"

# 2. Test health endpoint
curl http://localhost:8080/health
# Expected: 200 OK

# 3. Test beta signup
curl -X POST http://localhost:8080/api/beta-signup \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "email=test@example.com&consent=on"
# Expected: 200 OK with success message

# 4. Verify user in Supertokens dashboard
# Go to: https://supertokens.com/dashboard
# Should see: test@example.com in users list

# 5. Verify profile in NeonDB
psql "postgresql://..." -c "SELECT * FROM profiles WHERE email='test@example.com';"
# Expected: 1 row with beta_status='waitlist'

# 6. Check email received
# Check test@example.com inbox for welcome email
```

#### Step 5.2: Data Integrity Verification

**Create**: `/backend/scripts/verify_migration.sql`

```sql
-- 1. Verify no orphaned profiles
SELECT COUNT(*) as count FROM profiles WHERE id IS NULL;
-- Expected: 0

-- 2. Verify email uniqueness
SELECT email, COUNT(*) FROM profiles GROUP BY email HAVING COUNT(*) > 1;
-- Expected: 0 rows

-- 3. Verify all profiles complete
SELECT COUNT(*) FROM profiles
WHERE email IS NULL OR beta_status IS NULL OR created_at IS NULL;
-- Expected: 0

-- 4. Count check
SELECT COUNT(*) as total_profiles FROM profiles;
-- Compare with Supabase count

-- 5. Beta status distribution
SELECT beta_status, COUNT(*) FROM profiles GROUP BY beta_status;
```

---

### Phase 6: Cleanup & Deploy (1 hour)

#### Step 6.1: Remove Supabase Dependencies

**File**: `/backend/go.mod`

**Remove**:
```go
github.com/supabase-community/functions-go
github.com/supabase-community/gotrue-go
github.com/supabase-community/postgrest-go
github.com/supabase-community/storage-go
github.com/supabase-community/supabase-go
```

Run:
```bash
go mod tidy
```

#### Step 6.2: Remove Supabase Env Vars

**File**: `/backend/.env`

**Delete**:
```bash
SUPABASE_URL=
SUPABASE_ANON_KEY=
SUPABASE_PUBLISH_KEY=
SUPABASE_SECRET_KEY=
```

#### Step 6.3: Update Health Check

**File**: `/backend/internal/adapters/http/health_handler.go`

**Remove**:
- `testSupabaseAuthHealth` function
- Calls to `testSupabaseAPIEndpoint`

Keep only database health check.

#### Step 6.4: Build & Deploy

```bash
# Build production binary
cd backend
CGO_ENABLED=0 GOOS=linux go build -o backend cmd/api/main.go

# Deploy (update environment variables in production)
# Set: SUPERTOKENS_CONNECTION_URI, SUPERTOKENS_API_KEY
# Set: DB_HOST, DB_USERNAME, DB_PASSWORD (NeonDB values)
# Remove: SUPABASE_* variables
```

---

## Critical Files to Modify

### Must Modify (Core Migration)

1. **`/backend/go.mod`**
   - Add: `github.com/supertokens/supertokens-golang`
   - Remove: All `github.com/supabase-community/*` packages

2. **`/backend/internal/server/server.go`** (lines 35-60)
   - Add Supertokens initialization
   - Wire `SupertokensAuthService` into dependency graph

3. **`/backend/internal/server/routes.go`** (lines 508-569)
   - Replace `createBetaUser` with `createBetaUserSupertokens`
   - Add Supertokens middleware to router

4. **`/backend/internal/database/database.go`** (lines 75-82)
   - Update connection string for NeonDB (remove cert path)

5. **`/backend/.env`**
   - Add: `SUPERTOKENS_CONNECTION_URI`, `SUPERTOKENS_API_KEY`
   - Update: All `DB_*` variables to NeonDB values
   - Remove: All `SUPABASE_*` variables

### Must Create (New Files)

6. **`/backend/internal/adapters/external/supertokens_auth_service.go`**
   - `InitializeSupertokens()` function
   - `SupertokensAuthService` struct

7. **`/backend/internal/adapters/http/supertokens_middleware.go`**
   - `SupertokensCORS()` middleware
   - `VerifySession()` middleware

8. **`neondb_schema.sql`** (root directory)
   - Modified copy of `supabase/seeds/seed.sql` without RLS policies
   - Includes profiles table definition

---

## Post-Migration Checklist

- [ ] Database connected to NeonDB successfully
- [ ] Beta signup creates user in Supertokens
- [ ] Profile created in NeonDB profiles table
- [ ] Welcome email sent via Resend
- [ ] Health check endpoint returns 200 OK
- [ ] No Supabase dependencies in `go.mod`
- [ ] No Supabase env vars in production config
- [ ] Logs show no errors for 24 hours
- [ ] Rate limiting still works on beta signup

---

## Timeline Estimate

| Phase | Duration |
|-------|----------|
| 1. Database Migration | 1-2 hours |
| 2. Update DB Connection | 30 min |
| 3. Supertokens Integration | 3-4 hours |
| 4. Beta Signup Rewrite | 2-3 hours |
| 5. Testing & Validation | 2-3 hours |
| 6. Cleanup & Deploy | 1 hour |
| **TOTAL** | **11-16 hours** |

**Recommended Schedule**:
- **Day 1**: Phases 1-4 (infrastructure + implementation)
- **Day 2**: Phases 5-6 (testing + deployment)
- **Day 3**: Monitor production for 24 hours

---

## Success Criteria

Migration is successful when:

1. ✅ Beta users can sign up via `/api/beta-signup`
2. ✅ Users receive welcome + verification emails
3. ✅ User data stored in NeonDB `profiles` table
4. ✅ Authentication managed by Supertokens
5. ✅ Health check returns 200 OK
6. ✅ No Supabase dependencies remain in code
7. ✅ No errors in logs for 24 hours post-deployment

---

## Notes

- **Privacy Alignment**: This migration supports your "radical data transparency" goal by separating auth (Supertokens) from cannabis data (NeonDB), with all consumption data staying client-side encrypted.

- **Future Flutter App**: Same Supertokens setup will work for Flutter mobile app with Supertokens Flutter SDK.

- **Beta Grandfathering**: Users created during beta (in `profiles` table with `subscription_tier='beta_grandfathered'`) maintain lifetime sync access.

- **Cannabis Data Stays Local**: This migration only affects auth layer. Cannabis session data remains in IndexedDB (web) and Drift (mobile), encrypted client-side as designed.
