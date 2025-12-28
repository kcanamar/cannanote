package http

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type HealthResponse struct {
	Status     string                    `json:"status"`
	Timestamp  string                    `json:"timestamp"`
	Database   DatabaseHealth            `json:"database"`
	APIs       map[string]APIHealth      `json:"apis"`
	Message    string                    `json:"message,omitempty"`
	Summary    ConnectionSummary         `json:"summary"`
}

type DatabaseHealth struct {
	Status         string `json:"status"`
	Connection     string `json:"connection"`
	Error          string `json:"error,omitempty"`
	ResponseTimeMs int64  `json:"response_time_ms,omitempty"`
}

type APIHealth struct {
	Status         string `json:"status"`
	Endpoint       string `json:"endpoint"`
	Error          string `json:"error,omitempty"`
	ResponseTimeMs int64  `json:"response_time_ms,omitempty"`
}

type ConnectionSummary struct {
	TotalConnections int `json:"total_connections"`
	HealthyCount     int `json:"healthy_count"`
	UnhealthyCount   int `json:"unhealthy_count"`
	OverallHealth    int `json:"overall_health_percentage"`
}

// Simple rate limiting for health checks
var (
	healthCheckCache    *HealthResponse
	healthCheckTime     time.Time
	healthCheckMutex    sync.RWMutex
	healthCheckInterval = 10 * time.Second // Cache health results for 10 seconds
)

// HealthHandler provides comprehensive health check for all external connections
func HealthHandler(c *gin.Context) {
	health := HealthResponse{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		APIs:      make(map[string]APIHealth),
	}

	// Test direct PostgreSQL connection
	dbHealth := testDirectPostgresConnection()
	health.Database = dbHealth

	// Test all Supabase REST API endpoints
	health.APIs["consumption_methods"] = testSupabaseAPIEndpoint("consumption_methods", "?select=name&limit=1")
	health.APIs["cannabinoids"] = testSupabaseAPIEndpoint("cannabinoids", "?select=name,full_name,description,psychoactive,reported_experiences,compound_notes&order=name&limit=1")
	
	// Add any future API endpoints here
	// health.APIs["strains"] = testSupabaseAPIEndpoint("strains", "?select=name&limit=1")
	// health.APIs["entries"] = testSupabaseAPIEndpoint("entries", "?select=id&limit=1")

	// Calculate summary statistics
	totalConnections := 1 + len(health.APIs) // 1 for database + APIs
	healthyCount := 0
	
	if dbHealth.Status == "up" {
		healthyCount++
	}
	
	for _, api := range health.APIs {
		if api.Status == "up" {
			healthyCount++
		}
	}
	
	health.Summary = ConnectionSummary{
		TotalConnections: totalConnections,
		HealthyCount:     healthyCount,
		UnhealthyCount:   totalConnections - healthyCount,
		OverallHealth:    (healthyCount * 100) / totalConnections,
	}

	// Determine overall status based on health percentage
	if health.Summary.OverallHealth == 100 {
		health.Status = "up"
		health.Message = "All systems operational"
	} else if health.Summary.OverallHealth >= 75 {
		health.Status = "degraded"
		health.Message = fmt.Sprintf("Partial service - %d%% systems healthy", health.Summary.OverallHealth)
	} else if health.Summary.OverallHealth >= 25 {
		health.Status = "degraded"
		health.Message = fmt.Sprintf("Major issues - %d%% systems healthy", health.Summary.OverallHealth)
	} else {
		health.Status = "down"
		health.Message = "Critical failure - multiple systems down"
	}

	// Return appropriate HTTP status
	switch health.Status {
	case "up":
		c.JSON(http.StatusOK, health)
	case "degraded":
		c.JSON(http.StatusPartialContent, health)
	default:
		c.JSON(http.StatusServiceUnavailable, health)
	}
}

// testSupabaseAPIEndpoint tests connectivity to a specific Supabase REST API endpoint
func testSupabaseAPIEndpoint(tableName, query string) APIHealth {
	supabaseURL := os.Getenv("SUPABASE_URL")
	apiKey := os.Getenv("SUPABASE_ANON_KEY")
	
	if supabaseURL == "" || apiKey == "" {
		return APIHealth{
			Status:         "down",
			Endpoint:       "configuration_missing",
			Error:          "Missing Supabase configuration",
			ResponseTimeMs: 0,
		}
	}

	testURL := supabaseURL + "/rest/v1/" + tableName + query
	
	// Record start time for response measurement
	start := time.Now()
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	req, err := http.NewRequest("GET", testURL, nil)
	if err != nil {
		return APIHealth{
			Status:         "down",
			Endpoint:       testURL,
			Error:          "Failed to create request: " + err.Error(),
			ResponseTimeMs: time.Since(start).Milliseconds(),
		}
	}

	req.Header.Add("apikey", apiKey)
	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	responseTime := time.Since(start).Milliseconds()
	
	if err != nil {
		return APIHealth{
			Status:         "down",
			Endpoint:       testURL,
			Error:          "Connection failed: " + err.Error(),
			ResponseTimeMs: responseTime,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return APIHealth{
			Status:         "up",
			Endpoint:       testURL,
			ResponseTimeMs: responseTime,
		}
	}

	return APIHealth{
		Status:         "down", 
		Endpoint:       testURL,
		Error:          fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status),
		ResponseTimeMs: responseTime,
	}
}

// testDirectPostgresConnection tests direct PostgreSQL connectivity
func testDirectPostgresConnection() DatabaseHealth {
	// Build connection string from environment variables
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return DatabaseHealth{
			Status:     "down",
			Connection: "configuration_missing",
			Error:      "Missing DB_PASSWORD environment variable",
		}
	}

	connectionString := fmt.Sprintf("postgresql://postgres:%s@db.citdskdmralncvjyybin.supabase.co:5432/postgres", dbPassword)
	
	// Record start time for response measurement
	start := time.Now()
	
	// Attempt connection
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return DatabaseHealth{
			Status:         "down",
			Connection:     "direct_postgresql",
			Error:          "Failed to open connection: " + err.Error(),
			ResponseTimeMs: time.Since(start).Milliseconds(),
		}
	}
	defer db.Close()

	// Test the connection with a simple query
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	var result int
	err = db.QueryRowContext(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return DatabaseHealth{
			Status:         "down",
			Connection:     "direct_postgresql",
			Error:          "Query failed: " + err.Error(),
			ResponseTimeMs: time.Since(start).Milliseconds(),
		}
	}

	return DatabaseHealth{
		Status:         "up",
		Connection:     "direct_postgresql",
		ResponseTimeMs: time.Since(start).Milliseconds(),
	}
}