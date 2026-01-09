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

	// Test critical production endpoints
	health.APIs["consumption_methods"] = testSupabaseAPIEndpoint("consumption_methods", "?select=name&limit=1")
	health.APIs["cannabinoids"] = testSupabaseAPIEndpoint("cannabinoids", "?select=name,full_name,description,psychoactive,reported_experiences,compound_notes&order=name&limit=1")
	health.APIs["profiles"] = testSupabaseAPIEndpoint("profiles", "?select=id&limit=1")
	
	// Test Supabase Auth service
	health.APIs["supabase_auth"] = testSupabaseAuthHealth()
	
	// Add environment configuration check
	health.APIs["environment"] = testEnvironmentConfig()

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

	// Return appropriate HTTP status - sanitize for production
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "production" {
		// Production mode - minimal information disclosure
		sanitizedHealth := sanitizeHealthForProduction(health)
		switch health.Status {
		case "up":
			c.JSON(http.StatusOK, sanitizedHealth)
		case "degraded":
			c.JSON(http.StatusPartialContent, sanitizedHealth)
		default:
			c.JSON(http.StatusServiceUnavailable, sanitizedHealth)
		}
	} else {
		// Development mode - full details
		switch health.Status {
		case "up":
			c.JSON(http.StatusOK, health)
		case "degraded":
			c.JSON(http.StatusPartialContent, health)
		default:
			c.JSON(http.StatusServiceUnavailable, health)
		}
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
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbDatabase := os.Getenv("DB_DATABASE")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	
	if dbPassword == "" || dbHost == "" {
		return DatabaseHealth{
			Status:     "down",
			Connection: "configuration_missing",
			Error:      "Missing required database environment variables",
		}
	}

	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUsername, dbPassword, dbHost, dbPort, dbDatabase)
	
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

// testSupabaseAuthHealth tests Supabase Auth service availability
func testSupabaseAuthHealth() APIHealth {
	supabaseURL := os.Getenv("SUPABASE_URL")
	apiKey := os.Getenv("SUPABASE_ANON_KEY")
	
	if supabaseURL == "" || apiKey == "" {
		return APIHealth{
			Status:         "down",
			Endpoint:       "configuration_missing",
			Error:          "Missing Supabase Auth configuration",
			ResponseTimeMs: 0,
		}
	}

	// Test auth health endpoint
	authURL := supabaseURL + "/auth/v1/health"
	
	start := time.Now()
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	req, err := http.NewRequest("GET", authURL, nil)
	if err != nil {
		return APIHealth{
			Status:         "down",
			Endpoint:       authURL,
			Error:          "Failed to create auth request: " + err.Error(),
			ResponseTimeMs: time.Since(start).Milliseconds(),
		}
	}

	// Health endpoints typically don't require authentication - try without auth first
	// If this fails, we can add auth headers back

	resp, err := client.Do(req)
	responseTime := time.Since(start).Milliseconds()
	
	if err != nil {
		return APIHealth{
			Status:         "down",
			Endpoint:       authURL,
			Error:          "Auth connection failed: " + err.Error(),
			ResponseTimeMs: responseTime,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return APIHealth{
			Status:         "up",
			Endpoint:       authURL,
			ResponseTimeMs: responseTime,
		}
	}

	return APIHealth{
		Status:         "down", 
		Endpoint:       authURL,
		Error:          fmt.Sprintf("Auth HTTP %d: %s", resp.StatusCode, resp.Status),
		ResponseTimeMs: responseTime,
	}
}

// testEnvironmentConfig validates production environment configuration
func testEnvironmentConfig() APIHealth {
	start := time.Now()
	
	appEnv := os.Getenv("APP_ENV")
	ginMode := os.Getenv("GIN_MODE")
	debug := os.Getenv("DEBUG")
	appURL := os.Getenv("APP_URL")
	
	var errors []string
	
	// Check production environment settings
	if appEnv != "production" {
		errors = append(errors, fmt.Sprintf("APP_ENV is '%s', expected 'production'", appEnv))
	}
	
	if ginMode != "release" {
		errors = append(errors, fmt.Sprintf("GIN_MODE is '%s', expected 'release'", ginMode))
	}
	
	if debug != "false" {
		errors = append(errors, fmt.Sprintf("DEBUG is '%s', expected 'false'", debug))
	}
	
	if appURL == "" {
		errors = append(errors, "APP_URL is not set")
	}
	
	responseTime := time.Since(start).Milliseconds()
	
	if len(errors) == 0 {
		return APIHealth{
			Status:         "up",
			Endpoint:       "environment_config",
			ResponseTimeMs: responseTime,
		}
	}
	
	return APIHealth{
		Status:         "down",
		Endpoint:       "environment_config", 
		Error:          fmt.Sprintf("Environment issues: %v", errors),
		ResponseTimeMs: responseTime,
	}
}

// sanitizeHealthForProduction removes sensitive information from health response
func sanitizeHealthForProduction(health HealthResponse) HealthResponse {
	sanitized := HealthResponse{
		Status:    health.Status,
		Timestamp: health.Timestamp,
		Summary:   health.Summary,
	}
	
	// Simple status message without details
	switch health.Status {
	case "up":
		sanitized.Message = "All systems operational"
	case "degraded":
		sanitized.Message = "Service partially available"
	default:
		sanitized.Message = "Service unavailable"
	}
	
	// Sanitized database health - no connection details
	sanitized.Database = DatabaseHealth{
		Status: health.Database.Status,
	}
	if health.Database.Status != "up" {
		sanitized.Database.Error = "Database connectivity issues"
	}
	
	// Sanitized API health - no endpoints or detailed errors
	sanitized.APIs = make(map[string]APIHealth)
	for name, api := range health.APIs {
		apiHealth := APIHealth{
			Status: api.Status,
		}
		if api.Status != "up" {
			apiHealth.Error = "Service connectivity issues"
		}
		sanitized.APIs[name] = apiHealth
	}
	
	return sanitized
}