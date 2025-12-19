package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"backend/internal/database"
	"backend/internal/core/application"
	"backend/internal/adapters/repository"
	httpAdapters "backend/internal/adapters/http"
)

// Server holds the hexagonal architecture dependencies
type Server struct {
	port int
	db   database.Service

	// Hexagonal architecture components
	humanService  *application.HumanService
	humanHandlers *httpAdapters.HumanHandlers
}

// NewServer creates a new server with hexagonal architecture wiring
// Manual dependency injection - no formal DI framework for security
func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	// Initialize database service (existing code)
	dbService := database.New()

	// Get the underlying *sql.DB for repositories
	rawDB := dbService.GetDB()

	// Wire the hexagonal architecture
	// Repositories (adapters for ports)
	humanRepo := repository.NewSupabaseHumanRepository(rawDB)

	// TODO: Implement these when we add authentication
	// roleRepo := repository.NewSupabaseRoleRepository(rawDB)
	// authService := external.NewSupabaseAuthService()

	// Application services (use cases)
	humanService := application.NewHumanService(humanRepo, nil, nil)

	// HTTP handlers (adapters)
	humanHandlers := httpAdapters.NewHumanHandlers(humanService)

	// Create server instance
	newServer := &Server{
		port:          port,
		db:            dbService,
		humanService:  humanService,
		humanHandlers: humanHandlers,
	}

	// Configure HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

