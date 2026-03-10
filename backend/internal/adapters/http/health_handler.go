package http

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Services  map[string]ServiceHealth `json:"services"`
	Summary   HealthSummary          `json:"summary"`
}

type ServiceHealth struct {
	Status         string `json:"status"`
	ResponseTimeMs int64  `json:"response_time_ms,omitempty"`
	Details        string `json:"details,omitempty"`
	Error          string `json:"error,omitempty"`
}

type HealthSummary struct {
	Healthy   int    `json:"healthy"`
	Unhealthy int    `json:"unhealthy"`
	Total     int    `json:"total"`
	Percent   int    `json:"percent"`
	Message   string `json:"message"`
}

// HealthHandler provides health check for all services
func HealthHandler(c *gin.Context) {
	health := HealthResponse{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Services:  make(map[string]ServiceHealth),
	}

	// Check database connection
	health.Services["database"] = checkDatabase()

	// Check sessions table exists and is accessible
	health.Services["sessions"] = checkSessionsTable()

	// Check profiles table exists and is accessible
	health.Services["profiles"] = checkProfilesTable()

	// Check Resend API key is configured (optional - emails still work without it in dev)
	health.Services["email"] = checkEmailConfig()

	// Calculate summary
	healthy := 0
	for _, svc := range health.Services {
		if svc.Status == "up" {
			healthy++
		}
	}

	total := len(health.Services)
	percent := (healthy * 100) / total

	health.Summary = HealthSummary{
		Healthy:   healthy,
		Unhealthy: total - healthy,
		Total:     total,
		Percent:   percent,
	}

	// Determine overall status
	switch {
	case percent == 100:
		health.Status = "up"
		health.Summary.Message = "All systems operational"
	case percent >= 75:
		health.Status = "degraded"
		health.Summary.Message = "Most systems operational"
	case percent >= 50:
		health.Status = "degraded"
		health.Summary.Message = "Some systems degraded"
	default:
		health.Status = "down"
		health.Summary.Message = "Critical systems down"
	}

	// Return appropriate HTTP status
	// Database is critical - if it's down, return 503
	if health.Services["database"].Status != "up" {
		c.JSON(http.StatusServiceUnavailable, health)
		return
	}

	// Otherwise base on overall health
	switch health.Status {
	case "up":
		c.JSON(http.StatusOK, health)
	case "degraded":
		c.JSON(http.StatusOK, health) // Degraded is still functional
	default:
		c.JSON(http.StatusServiceUnavailable, health)
	}
}

// checkDatabase tests PostgreSQL connectivity
func checkDatabase() ServiceHealth {
	start := time.Now()

	connStr := os.Getenv("DB_CONNECTION_URI")
	if connStr == "" {
		// Build from individual vars
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		database := os.Getenv("DB_DATABASE")
		username := os.Getenv("DB_USERNAME")
		password := os.Getenv("DB_PASSWORD")

		if host == "" || password == "" {
			return ServiceHealth{
				Status:         "down",
				ResponseTimeMs: time.Since(start).Milliseconds(),
				Error:          "Database not configured",
			}
		}
		connStr = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username, password, host, port, database)
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return ServiceHealth{
			Status:         "down",
			ResponseTimeMs: time.Since(start).Milliseconds(),
			Error:          "Failed to open connection",
		}
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result int
	err = db.QueryRowContext(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return ServiceHealth{
			Status:         "down",
			ResponseTimeMs: time.Since(start).Milliseconds(),
			Error:          "Query failed",
		}
	}

	return ServiceHealth{
		Status:         "up",
		ResponseTimeMs: time.Since(start).Milliseconds(),
		Details:        "PostgreSQL connected",
	}
}

// checkSessionsTable verifies the sessions table exists
func checkSessionsTable() ServiceHealth {
	start := time.Now()

	connStr := os.Getenv("DB_CONNECTION_URI")
	if connStr == "" {
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		database := os.Getenv("DB_DATABASE")
		username := os.Getenv("DB_USERNAME")
		password := os.Getenv("DB_PASSWORD")
		connStr = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username, password, host, port, database)
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return ServiceHealth{
			Status:         "down",
			ResponseTimeMs: time.Since(start).Milliseconds(),
			Error:          "Database connection failed",
		}
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if sessions table exists
	var exists bool
	err = db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = 'sessions'
		)
	`).Scan(&exists)

	if err != nil {
		return ServiceHealth{
			Status:         "down",
			ResponseTimeMs: time.Since(start).Milliseconds(),
			Error:          "Table check failed",
		}
	}

	if !exists {
		return ServiceHealth{
			Status:         "down",
			ResponseTimeMs: time.Since(start).Milliseconds(),
			Error:          "Sessions table not found - run migration",
		}
	}

	return ServiceHealth{
		Status:         "up",
		ResponseTimeMs: time.Since(start).Milliseconds(),
		Details:        "Sessions table ready",
	}
}

// checkProfilesTable verifies the profiles table exists and has auth columns
func checkProfilesTable() ServiceHealth {
	start := time.Now()

	connStr := os.Getenv("DB_CONNECTION_URI")
	if connStr == "" {
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		database := os.Getenv("DB_DATABASE")
		username := os.Getenv("DB_USERNAME")
		password := os.Getenv("DB_PASSWORD")
		connStr = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username, password, host, port, database)
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return ServiceHealth{
			Status:         "down",
			ResponseTimeMs: time.Since(start).Milliseconds(),
			Error:          "Database connection failed",
		}
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if profiles table has password_hash column (migration applied)
	var hasPasswordHash bool
	err = db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.columns
			WHERE table_schema = 'public'
			AND table_name = 'profiles'
			AND column_name = 'password_hash'
		)
	`).Scan(&hasPasswordHash)

	if err != nil {
		return ServiceHealth{
			Status:         "down",
			ResponseTimeMs: time.Since(start).Milliseconds(),
			Error:          "Column check failed",
		}
	}

	if !hasPasswordHash {
		return ServiceHealth{
			Status:         "down",
			ResponseTimeMs: time.Since(start).Milliseconds(),
			Error:          "Auth columns missing - run migration",
		}
	}

	return ServiceHealth{
		Status:         "up",
		ResponseTimeMs: time.Since(start).Milliseconds(),
		Details:        "Profiles table ready",
	}
}

// checkEmailConfig verifies Resend is configured
func checkEmailConfig() ServiceHealth {
	start := time.Now()

	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return ServiceHealth{
			Status:         "down",
			ResponseTimeMs: time.Since(start).Milliseconds(),
			Error:          "RESEND_API_KEY not configured",
		}
	}

	// Just check it's configured, don't make API calls
	return ServiceHealth{
		Status:         "up",
		ResponseTimeMs: time.Since(start).Milliseconds(),
		Details:        "Resend configured",
	}
}
