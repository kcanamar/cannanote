package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"backend/internal/adapters/repository"
	"backend/internal/core/application"
	"backend/internal/database"
	httpAdapters "backend/internal/adapters/http"
)

// Server holds the hexagonal architecture dependencies
// This struct orchestrates the wiring of all layers: domain, application, and adapters
// Following dependency inversion principle, higher-level modules don't depend on lower-level modules
type Server struct {
	port int
	db   database.Service

	// Hexagonal architecture components
	// Application layer services contain business logic
	humanService  *application.HumanService
	// HTTP adapters translate web requests to domain operations
	humanHandlers *httpAdapters.HumanHandlers
}

// NewServer creates a new server with hexagonal architecture wiring
// Manual dependency injection - no formal DI framework for security
// This approach keeps dependencies explicit and prevents hidden coupling
func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	// Initialize database service
	dbService := database.New()

	// Get the underlying *sql.DB for repositories
	rawDB := dbService.GetDB()

	// Wire the hexagonal architecture layers:

	// 1. Repositories (adapters implementing ports)
	humanRepo := repository.NewSupabaseHumanRepository(rawDB)

	// 2. Application services (use cases and business workflows)
	// Auth is now handled directly in route handlers via AuthService
	humanService := application.NewHumanService(humanRepo, nil, nil)

	// 4. HTTP handlers (adapters for web interface)
	// These translate HTTP requests/responses to application service calls
	humanHandlers := httpAdapters.NewHumanHandlers(humanService)

	// Create server instance with all dependencies wired
	newServer := &Server{
		port:          port,
		db:            dbService,
		humanService:  humanService,
		humanHandlers: humanHandlers,
	}

	// Configure HTTP server with production-ready timeouts
	// These timeouts prevent resource exhaustion from slow clients
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(), // Gin router with all endpoints
		IdleTimeout:  time.Minute,                // How long to keep connections open
		ReadTimeout:  10 * time.Second,          // Max time to read request
		WriteTimeout: 30 * time.Second,          // Max time to write response
	}

	return server
}