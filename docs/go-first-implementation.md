# Go-First Implementation Guide

## Document Control
- **Version**: 2.0 (Revised Strategy)
- **Date**: December 18, 2024
- **Classification**: Internal Use Only

## Strategic Overview

**Key Change**: We're launching with Go MVP from day 1, not migrating from Express. This eliminates technical debt while maintaining rapid time-to-market. HIPAA compliance is delayed until provider partnerships justify infrastructure costs (~$20K MRR threshold).

## Updated Technology Stack

### Core Stack (2025 Best Practices)
```go
// Backend
Go 1.22+                    // Latest stable, improved performance
Chi router                  // Lightweight, fast routing
templ                      // Type-safe templates (github.com/a-h/templ)

// Frontend  
HTMX                       // Hypermedia-driven interactions
Tailwind CSS               // Utility-first styling
Alpine.js (minimal)        // Lightweight JS for complex interactions

// Database & Auth
Supabase PostgreSQL        // HIPAA-capable when needed
Supabase Auth              // OAuth + session management
Casbin                     // RBAC for future compliance

// Infrastructure
Fly.io                     // SOC2 certified, edge deployment
Docker                     // Containerization
Fly.io secrets            // Built-in secret management (not OpenBao)

// Monitoring
Sentry                     // Error tracking (free tier)
Prometheus                 // Metrics collection  
Zap                        // Structured logging
Grafana                    // Dashboards (free tier)
```

## Week 1-6: MVP Development

### Directory Structure
```
cannanote/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── entries.go
│   │   └── analytics.go        # B2B dashboard
│   ├── models/
│   │   ├── human.go            # User model (note: "human" terminology)
│   │   ├── entry.go
│   │   └── dispensary.go       # B2B partnerships
│   ├── services/
│   │   ├── auth.go
│   │   ├── entry.go
│   │   └── analytics.go        # Dispensary insights
│   └── middleware/
│       ├── auth.go
│       ├── logging.go
│       └── rbac.go
├── templates/
│   ├── components/             # templ components
│   ├── pages/                 # Full page templates
│   └── layouts/               # Base layouts
├── static/
│   ├── css/
│   ├── js/
│   └── images/
├── migrations/
├── docker-compose.yml
├── Dockerfile
├── fly.toml
├── Makefile
└── go.mod
```

### Week 1-2: Foundation

#### Makefile Setup
```make
.PHONY: dev build test lint deploy

dev:
	@echo "Starting development server..."
	@templ generate
	@air -c .air.toml

build:
	@echo "Building production binary..."
	@templ generate
	@CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/cannanote ./cmd/server

test:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

lint:
	@echo "Running linters..."
	@gosec ./...
	@staticcheck ./...
	@golangci-lint run

deploy:
	@echo "Deploying to Fly.io..."
	@fly deploy --remote-only

docker-dev:
	@echo "Starting Docker development environment..."
	@docker-compose up -d
```

#### Core Models
```go
// internal/models/human.go
package models

import (
    "time"
    "github.com/google/uuid"
)

type Human struct {
    ID          uuid.UUID `json:"id" db:"id"`
    Username    string    `json:"username" db:"username"`
    Email       string    `json:"email" db:"email"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
    
    // Subscription info
    PlanType    string    `json:"plan_type" db:"plan_type"` // free, premium, enterprise
    SubscribedAt *time.Time `json:"subscribed_at,omitempty" db:"subscribed_at"`
    
    // Privacy settings
    PublicProfile bool    `json:"public_profile" db:"public_profile"`
    ShareData     bool    `json:"share_data" db:"share_data"`
}

// internal/models/entry.go
type CannabisEntry struct {
    ID               uuid.UUID  `json:"id" db:"id"`
    HumanID          uuid.UUID  `json:"human_id" db:"human_id"`
    
    // Cannabis details
    StrainName       string     `json:"strain_name" db:"strain_name"`
    StrainType       string     `json:"strain_type" db:"strain_type"` // indica, sativa, hybrid
    THCPercentage    *float64   `json:"thc_percentage,omitempty" db:"thc_percentage"`
    CBDPercentage    *float64   `json:"cbd_percentage,omitempty" db:"cbd_percentage"`
    ConsumptionMethod string    `json:"consumption_method" db:"consumption_method"`
    AmountConsumed   *float64   `json:"amount_consumed,omitempty" db:"amount_consumed"`
    
    // Experience tracking
    Effects          []string   `json:"effects" db:"effects"`
    MoodBefore       *int       `json:"mood_before,omitempty" db:"mood_before"` // 1-10 scale
    MoodAfter        *int       `json:"mood_after,omitempty" db:"mood_after"`
    SleepQuality     *int       `json:"sleep_quality,omitempty" db:"sleep_quality"`
    Notes           string     `json:"notes" db:"notes"`
    
    // Social features
    IsPublic        bool       `json:"is_public" db:"is_public"`
    Likes           int        `json:"likes" db:"likes"`
    
    // Timestamps
    ConsumedAt      time.Time  `json:"consumed_at" db:"consumed_at"`
    CreatedAt       time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}
```

### Week 3-4: HTMX Frontend with templ

#### templ Components
```go
// templates/components/entry_form.templ
package components

import "github.com/cannanote/internal/models"

templ EntryForm(strains []string) {
    <form 
        hx-post="/api/entries" 
        hx-target="#entries-list" 
        hx-swap="afterbegin"
        class="bg-white p-6 rounded-lg shadow-md">
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label class="block text-sm font-medium text-gray-700">Strain Name</label>
                <input 
                    type="text" 
                    name="strain_name"
                    list="strain-options"
                    class="mt-1 block w-full rounded-md border-gray-300 shadow-sm"
                    required>
                <datalist id="strain-options">
                    for _, strain := range strains {
                        <option value={ strain }></option>
                    }
                </datalist>
            </div>
            
            <div>
                <label class="block text-sm font-medium text-gray-700">Type</label>
                <select name="strain_type" class="mt-1 block w-full rounded-md border-gray-300 shadow-sm">
                    <option value="indica">Indica</option>
                    <option value="sativa">Sativa</option>
                    <option value="hybrid">Hybrid</option>
                </select>
            </div>
            
            <div>
                <label class="block text-sm font-medium text-gray-700">Consumption Method</label>
                <select name="consumption_method" class="mt-1 block w-full rounded-md border-gray-300 shadow-sm">
                    <option value="flower">Flower</option>
                    <option value="vape">Vape</option>
                    <option value="edible">Edible</option>
                    <option value="concentrate">Concentrate</option>
                </select>
            </div>
            
            <div>
                <label class="block text-sm font-medium text-gray-700">Mood Before (1-10)</label>
                <input 
                    type="range" 
                    name="mood_before" 
                    min="1" 
                    max="10"
                    class="w-full">
            </div>
        </div>
        
        <div class="mt-4">
            <label class="block text-sm font-medium text-gray-700">Notes</label>
            <textarea 
                name="notes"
                rows="3"
                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm"
                placeholder="How are you feeling? What effects are you hoping for?"></textarea>
        </div>
        
        <div class="mt-6">
            <button 
                type="submit"
                class="w-full bg-green-600 text-white py-2 px-4 rounded-md hover:bg-green-700 focus:ring-2 focus:ring-green-500">
                Track Experience
            </button>
        </div>
    </form>
}

// templates/components/entry_card.templ
templ EntryCard(entry models.CannabisEntry) {
    <div class="bg-white p-4 rounded-lg shadow-sm border border-gray-200 mb-4">
        <div class="flex justify-between items-start mb-2">
            <h3 class="text-lg font-semibold text-gray-900">{ entry.StrainName }</h3>
            <span class="text-sm text-gray-500">{ entry.ConsumedAt.Format("Jan 2, 3:04 PM") }</span>
        </div>
        
        <div class="flex gap-2 mb-2">
            <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                { entry.StrainType }
            </span>
            <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                { entry.ConsumptionMethod }
            </span>
        </div>
        
        if len(entry.Effects) > 0 {
            <div class="mb-2">
                <span class="text-sm text-gray-600">Effects: </span>
                for i, effect := range entry.Effects {
                    <span class="text-sm text-gray-800">{ effect }</span>
                    if i < len(entry.Effects)-1 {
                        <span class="text-gray-400">, </span>
                    }
                }
            </div>
        }
        
        if entry.MoodBefore != nil && entry.MoodAfter != nil {
            <div class="flex gap-4 text-sm text-gray-600 mb-2">
                <span>Mood: { fmt.Sprintf("%d → %d", *entry.MoodBefore, *entry.MoodAfter) }</span>
            </div>
        }
        
        if entry.Notes != "" {
            <p class="text-sm text-gray-700 mt-2">{ entry.Notes }</p>
        }
        
        <div class="flex justify-between items-center mt-3 pt-3 border-t border-gray-100">
            <button 
                hx-post={ fmt.Sprintf("/api/entries/%s/like", entry.ID) }
                hx-target="#like-count-" + entry.ID.String()
                hx-swap="innerHTML"
                class="flex items-center text-sm text-gray-500 hover:text-red-600">
                <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M3.172 5.172a4 4 0 015.656 0L10 6.343l1.172-1.171a4 4 0 115.656 5.656L10 17.657l-6.828-6.829a4 4 0 010-5.656z"/>
                </svg>
                <span id={ "like-count-" + entry.ID.String() }>{ fmt.Sprintf("%d", entry.Likes) }</span>
            </button>
            
            <div class="flex gap-2">
                <button 
                    hx-get={ fmt.Sprintf("/entries/%s/edit", entry.ID) }
                    hx-target="#modal"
                    class="text-sm text-gray-500 hover:text-blue-600">
                    Edit
                </button>
                <button 
                    hx-delete={ fmt.Sprintf("/api/entries/%s", entry.ID) }
                    hx-target="closest div"
                    hx-swap="outerHTML"
                    hx-confirm="Are you sure you want to delete this entry?"
                    class="text-sm text-gray-500 hover:text-red-600">
                    Delete
                </button>
            </div>
        </div>
    </div>
}
```

### Week 5-6: B2B Analytics Dashboard

#### Dispensary Analytics
```go
// internal/handlers/analytics.go
package handlers

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/cannanote/internal/services"
    "github.com/cannanote/templates/pages"
)

type AnalyticsHandler struct {
    analytics *services.AnalyticsService
}

func (h *AnalyticsHandler) DispensaryDashboard(w http.ResponseWriter, r *http.Request) {
    dispensaryID := chi.URLParam(r, "dispensaryID")
    
    // Get anonymized analytics for dispensary
    insights, err := h.analytics.GetDispensaryInsights(r.Context(), dispensaryID)
    if err != nil {
        http.Error(w, "Failed to load analytics", http.StatusInternalServerError)
        return
    }
    
    component := pages.DispensaryAnalytics(insights)
    component.Render(r.Context(), w)
}

// internal/services/analytics.go
type DispensaryInsights struct {
    DispensaryID     string                 `json:"dispensary_id"`
    TotalHumans      int                    `json:"total_humans"`
    ActiveHumans     int                    `json:"active_humans"`
    PopularStrains   []StrainPopularity     `json:"popular_strains"`
    EffectsReported  map[string]int         `json:"effects_reported"`
    ConsumptionTrends []ConsumptionTrend    `json:"consumption_trends"`
    SatisfactionScore float64               `json:"satisfaction_score"`
    RecommendationRate float64              `json:"recommendation_rate"`
}

type StrainPopularity struct {
    StrainName      string  `json:"strain_name"`
    StrainType      string  `json:"strain_type"`
    PurchaseCount   int     `json:"purchase_count"`
    AvgRating      float64  `json:"avg_rating"`
    RecommendRate   float64  `json:"recommend_rate"`
}

func (s *AnalyticsService) GetDispensaryInsights(ctx context.Context, dispensaryID string) (*DispensaryInsights, error) {
    // Query anonymized data only - no PHI
    query := `
        SELECT 
            COUNT(DISTINCT e.human_id) as total_humans,
            COUNT(DISTINCT CASE WHEN e.created_at > NOW() - INTERVAL '30 days' THEN e.human_id END) as active_humans,
            AVG(CASE WHEN e.mood_after IS NOT NULL AND e.mood_before IS NOT NULL 
                THEN e.mood_after - e.mood_before END) as satisfaction_score
        FROM entries e 
        JOIN human_dispensary_connections hdc ON e.human_id = hdc.human_id 
        WHERE hdc.dispensary_id = $1
        AND e.created_at > NOW() - INTERVAL '90 days'
    `
    
    // Execute query and build insights...
    return insights, nil
}
```

## Revenue Generation Strategy

### Freemium Model Implementation
```go
// internal/middleware/subscription.go
package middleware

func RequirePremium(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        human := GetHumanFromContext(r.Context())
        
        if human.PlanType == "free" {
            // Show upgrade modal instead of blocking
            component := components.UpgradeModal("This feature requires a Premium subscription")
            component.Render(r.Context(), w)
            return
        }
        
        next(w, r)
    }
}

// Feature limitations for free tier
func (s *EntryService) CanCreateEntry(humanID uuid.UUID) (bool, error) {
    if human.PlanType == "free" {
        count, err := s.GetMonthlyEntryCount(humanID)
        if err != nil {
            return false, err
        }
        return count < 10, nil // 10 entries per month for free
    }
    
    return true, nil
}
```

### B2B Partnership Integration
```go
// internal/models/dispensary.go
type Dispensary struct {
    ID              uuid.UUID `json:"id" db:"id"`
    Name            string    `json:"name" db:"name"`
    Location        string    `json:"location" db:"location"`
    ContactEmail    string    `json:"contact_email" db:"contact_email"`
    
    // Partnership details
    PartnershipType string    `json:"partnership_type" db:"partnership_type"` // analytics, affiliate, full
    MonthlyFee      *int      `json:"monthly_fee,omitempty" db:"monthly_fee"`
    CommissionRate  *float64  `json:"commission_rate,omitempty" db:"commission_rate"`
    
    // Settings
    IsActive        bool      `json:"is_active" db:"is_active"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type HumanDispensaryConnection struct {
    HumanID      uuid.UUID `db:"human_id"`
    DispensaryID uuid.UUID `db:"dispensary_id"`
    CreatedAt    time.Time `db:"created_at"`
    
    // Purchase tracking (for affiliate revenue)
    LastPurchase *time.Time `db:"last_purchase"`
    TotalSpent   *float64   `db:"total_spent"`
}
```

## Deployment Configuration

### Fly.io Setup
```toml
# fly.toml
app = "cannanote"
primary_region = "sea"

[env]
  ENVIRONMENT = "production"
  PORT = "8080"

[http_service]
  internal_port = 8080
  force_https = true
  auto_stop_machines = false
  auto_start_machines = true

  [[http_service.checks]]
    grace_period = "10s"
    interval = "30s"
    method = "GET"
    timeout = "5s"
    path = "/health"

[vm]
  cpu_kind = "shared"
  cpus = 1
  memory_mb = 512

[[services]]
  internal_port = 8080
  protocol = "tcp"

  [services.concurrency]
    hard_limit = 100
    soft_limit = 50

  [[services.http_checks]]
    grace_period = "10s"
    interval = "30s"
    method = "GET"
    path = "/health"
    timeout = "5s"

  [[services.ports]]
    handlers = ["http"]
    port = 80
    force_https = true

  [[services.ports]]
    handlers = ["tls", "http"]
    port = 443
```

### Docker Configuration
```dockerfile
# Dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

COPY . .
RUN templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o cannanote ./cmd/server

FROM gcr.io/distroless/static-debian11

COPY --from=builder /app/cannanote /cannanote
COPY --from=builder /app/static /static

EXPOSE 8080
USER 65534:65534

ENTRYPOINT ["/cannanote"]
```

### Development Environment
```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - .:/app
    environment:
      - ENVIRONMENT=development
      - SUPABASE_URL=${SUPABASE_URL}
      - SUPABASE_ANON_KEY=${SUPABASE_ANON_KEY}
      - SUPABASE_SERVICE_KEY=${SUPABASE_SERVICE_KEY}
    depends_on:
      - prometheus
    networks:
      - cannanote

  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
    networks:
      - cannanote

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    networks:
      - cannanote

networks:
  cannanote:
    driver: bridge
```

## Success Metrics & KPIs

### Week 1-6 KPIs
```go
// Track these metrics in Prometheus
var (
    humanSignups = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "cannanote_human_signups_total",
            Help: "Total number of human signups",
        },
        []string{"plan_type", "referrer"},
    )

    entriesCreated = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "cannanote_entries_created_total", 
            Help: "Total number of cannabis entries created",
        },
        []string{"plan_type", "strain_type"},
    )

    dispensaryPartnershipRequests = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "cannanote_dispensary_requests_total",
            Help: "Total number of dispensary partnership requests",
        },
        []string{"partnership_type", "status"},
    )
)
```

### Target Metrics
- **Human Acquisition Cost**: <$20 (through organic + referrals)
- **Weekly Active Humans**: 50% of total signups
- **Premium Conversion Rate**: 5-10% (freemium to paid)
- **Dispensary Partnership Rate**: 1 partnership per 50 outreach attempts
- **B2B Revenue per Partnership**: $100-500/month

This implementation guide provides a clear path to launching CannaNote 2.0 with Go from day 1, focusing on rapid revenue generation through both individual humans and B2B dispensary partnerships while keeping infrastructure costs minimal until HIPAA compliance is justified by revenue thresholds.