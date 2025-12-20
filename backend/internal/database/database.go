package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	// GetDB returns the underlying database connection for repository use
	GetDB() *sql.DB
}

type service struct {
	db *sql.DB
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		log.Println("DEBUG: Reusing existing database connection")
		return dbInstance
	}
	
	// Debug environment variables (mask password)
	log.Printf("DEBUG: DB_HOST=%s", host)
	log.Printf("DEBUG: DB_PORT=%s", port)
	log.Printf("DEBUG: DB_DATABASE=%s", database)
	log.Printf("DEBUG: DB_USERNAME=%s", username)
	log.Printf("DEBUG: DB_PASSWORD=%s", maskPassword(password))
	log.Printf("DEBUG: DB_SCHEMA=%s", schema)
	
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=verify-full&sslrootcert=certs/prod-ca-2021.crt&search_path=%s", username, password, host, port, database, schema)
	
	// Debug connection string (mask password)
	log.Printf("DEBUG: Connection string: %s", maskConnectionString(connStr))
	
	log.Println("DEBUG: Attempting to open database connection...")
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Printf("ERROR: Failed to open database connection: %v", err)
		log.Fatal(err)
	}
	log.Println("DEBUG: Database connection opened successfully")
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	log.Println("DEBUG: Starting database health check...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	log.Println("DEBUG: Attempting to ping database...")
	err := s.db.PingContext(ctx)
	if err != nil {
		log.Printf("ERROR: Database ping failed: %v", err)
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		// Don't use log.Fatalf - just return the error for debugging
		return stats
	}
	
	log.Println("DEBUG: Database ping successful!")

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.db.Close()
}

// GetDB returns the underlying database connection
// This exposes the *sql.DB for use by repository adapters
func (s *service) GetDB() *sql.DB {
	return s.db
}

// maskPassword masks password for logging (shows first 2 and last 2 chars)
func maskPassword(password string) string {
	if len(password) <= 4 {
		return "****"
	}
	return password[:2] + "****" + password[len(password)-2:]
}

// maskConnectionString masks the password in connection string for logging
func maskConnectionString(connStr string) string {
	// Find password section in postgresql://user:password@host:port/db
	if strings.Contains(connStr, "://") && strings.Contains(connStr, "@") {
		parts := strings.Split(connStr, "@")
		if len(parts) == 2 {
			userPart := parts[0]
			if strings.Contains(userPart, ":") {
				userPassParts := strings.Split(userPart, ":")
				if len(userPassParts) >= 3 {
					// postgresql://user:password -> mask password
					userPassParts[2] = maskPassword(userPassParts[2])
					return strings.Join(userPassParts, ":") + "@" + parts[1]
				}
			}
		}
	}
	return "****MASKED****"
}
