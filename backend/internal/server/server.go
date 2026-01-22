package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"backend/internal/database"
	"backend/internal/core/application"
	"backend/internal/adapters/repository"
	"backend/internal/adapters/external"
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

	// Initialize Supertokens FIRST (before any route registration)
	log.Println("Initializing Supertokens authentication...")
	if err := external.InitializeSupertokens(); err != nil {
		log.Printf("WARNING: Supertokens initialization failed: %v", err)
		log.Println("Continuing without Supertokens - auth features will be disabled")
	}

	// Initialize database service (existing infrastructure code)
	dbService := database.New()

	// Get the underlying *sql.DB for repositories
	// Repository adapters need direct database access for SQL operations
	rawDB := dbService.GetDB()

	// Wire the hexagonal architecture layers from inside out:

	// 1. Repositories (adapters implementing ports)
	// These translate between domain contracts and external systems
	humanRepo := repository.NewSupabaseHumanRepository(rawDB)

	// 2. Authentication service (Supertokens adapter)
	authService := external.NewSupertokensAuthService()

	// 3. Application services (use cases and business workflows)
	// These orchestrate domain entities and coordinate with repositories
	// Pass authService as the third parameter (was nil before)
	humanService := application.NewHumanService(humanRepo, nil, authService)

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