# CannaNote Backend Service

High-performance Go backend service for CannaNote personal cannabis journaling platform, implementing hexagonal architecture with clean separation of concerns, comprehensive privacy controls, and scalable domain-driven design.

## Architecture Overview

The backend implements **Hexagonal Architecture** (Ports & Adapters) to achieve vendor-agnostic design, high testability, and maintainable business logic isolation. Core domain remains pure and independent from external frameworks, databases, or delivery mechanisms.

### Architectural Principles

- **Core Domain Purity** — No I/O, frameworks, or external dependencies in domain layer
- **Dependency Inversion** — Core depends on abstractions, never concrete implementations
- **Manual Dependency Injection** — Explicit wiring for security and clarity
- **Port-Driven Development** — Contracts defined by business needs, not infrastructure constraints
- **Privacy by Design** — Data protection and user consent integrated at architectural level

### Technology Stack

- **Language**: Go 1.24+ with modern language features
- **HTTP Framework**: Gin router with middleware pipeline
- **Template Engine**: templ for type-safe HTML template compilation
- **Frontend Enhancement**: HTMX for server-driven dynamic interactions
- **CSS Framework**: Tailwind CSS for utility-first styling  
- **Database**: Supabase PostgreSQL with Row Level Security
- **Authentication**: Supabase Auth with JWT token validation
- **Deployment**: Docker containers on Fly.io platform
- **Architecture Pattern**: Hexagonal (Ports & Adapters)
- **Dependency Management**: Manual injection (no container frameworks)

## Project Structure

```
backend/
├── cmd/                              # Application entry points
│   ├── api/                         # API server entry point
│   │   └── main.go                  # HTTP server bootstrap and configuration
│   └── web/                         # Web application layer
│       ├── assets/                  # Static asset management
│       │   ├── css/                 # Compiled Tailwind CSS
│       │   ├── js/                  # JavaScript dependencies (HTMX)
│       │   └── images/              # Brand assets and static images
│       │       └── logos/           # Logo variations (logo-book.svg, logo-square.svg)
│       ├── *.templ                  # Type-safe HTML templates
│       ├── *_templ.go              # Generated Go template code
│       ├── base.templ              # Base layout template
│       ├── hello.templ             # Example component template
│       └── efs.go                  # Embedded file system configuration
├── internal/                        # Private application code
│   ├── core/                       # Core business logic (hexagon center)
│   │   ├── domain/                 # Pure business entities and domain rules
│   │   │   ├── human.go           # Human entity with cannabis preferences
│   │   │   ├── entry.go           # Cannabis experience entries (planned)
│   │   │   └── errors.go          # Domain-specific error definitions
│   │   ├── application/           # Use cases and business workflows
│   │   │   ├── human_service.go   # Human management business logic
│   │   │   └── entry_service.go   # Entry management workflows (planned)
│   │   └── ports/                 # Interface contracts and abstractions
│   │       ├── human_repository.go # Data persistence contracts
│   │       └── auth_service.go     # Authentication service contracts
│   ├── adapters/                   # External integrations (hexagon periphery)
│   │   ├── http/                   # Primary adapters (inbound requests)
│   │   │   ├── health_handler.go   # Health check endpoint implementation
│   │   │   ├── human_handlers.go   # Human management REST API
│   │   │   └── learn_handlers.go   # Educational content handlers
│   │   ├── repository/             # Secondary adapters (outbound data)
│   │   │   └── supabase_human_repository.go # PostgreSQL implementation
│   │   └── external/               # External service integrations (planned)
│   ├── database/                   # Database service abstractions
│   │   └── database.go            # Connection management and configuration
│   └── server/                     # HTTP server infrastructure
│       ├── server.go              # Dependency injection and server setup
│       └── routes.go              # Route registration and middleware
├── tests/                          # Test suites and utilities
│   └── handler_test.go            # HTTP handler integration tests
├── certs/                          # SSL certificate management
│   └── prod-ca-2021.crt           # Production certificate authority
├── Dockerfile                      # Container build definition
├── fly.toml                        # Fly.io deployment configuration
├── Makefile                        # Build automation and development commands
├── docker-compose.yml              # Local development environment
├── go.mod                          # Go module dependencies
└── go.sum                          # Dependency checksums
```

## Current Implementation Status

### Implemented Features

#### Core Infrastructure
- **HTTP Server** — Gin-based HTTP server with middleware pipeline
- **Template System** — templ integration for type-safe HTML generation
- **Static Asset Management** — Embedded file system for CSS, JS, and images
- **Health Monitoring** — Comprehensive health check endpoints with database connectivity
- **Database Integration** — Supabase PostgreSQL connection with connection pooling

#### Human Management Domain
- **Domain Entity** — Human entity with profile management and preference storage
- **Business Logic** — Registration, authentication, and profile update workflows
- **Repository Pattern** — Abstract data persistence with Supabase implementation
- **HTTP API** — RESTful endpoints for human management operations
- **Template Integration** — Server-rendered UI components with HTMX enhancement

#### Educational Content System
- **Content Handlers** — Cannabis education content delivery endpoints
- **Template Rendering** — Educational content with cannabinoid and terpene information
- **Static Content** — Integrated educational resources and reference materials

### Current API Endpoints

```
GET  /health           # System health check with database connectivity
GET  /                 # Application home page
GET  /hello            # Example template rendering
GET  /learn/cannabinoids # Educational content for cannabinoids
POST /humans           # Create new human profile
GET  /humans/:id       # Retrieve human profile by ID
PUT  /humans/:id       # Update existing human profile
```

### Database Schema

Current PostgreSQL schema implementation:

```sql
-- Human profiles with cannabis preferences
CREATE TABLE humans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    preferences JSONB DEFAULT '{}',
    consent_given BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Row Level Security policies
ALTER TABLE humans ENABLE ROW LEVEL SECURITY;
CREATE POLICY "Users can only see their own profile" ON humans
    FOR ALL USING (auth.uid() = id);
```

## Development Environment

### Prerequisites

- **Go 1.24+** — Latest Go runtime with module support
- **Make** — Build automation and command execution
- **Docker** — Container runtime for local development database
- **Supabase CLI** — Database management and migration tools

### Environment Configuration

Required environment variables for development and production:

```bash
# Database Configuration
DB_HOST=db.citdskdmralncvjyybin.supabase.co
DB_PORT=5432
DB_DATABASE=postgres
DB_USERNAME=postgres
DB_PASSWORD=your_secure_password
DB_SCHEMA=public

# Supabase Configuration
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your_anonymous_key
SUPABASE_SERVICE_ROLE_KEY=your_service_role_key

# Application Configuration
PORT=3001
APP_ENV=development
GIN_MODE=debug

# Security Configuration
JWT_SECRET=your_jwt_secret_key
```

### Development Commands

```bash
# Development workflow
make dev                # Start development server with hot reload
make build              # Build application binary
make test               # Run complete test suite
make lint               # Execute linting and code quality checks
make clean              # Clean build artifacts and cache

# Database management
make db-reset           # Reset local database to clean state
make db-migrate         # Apply pending database migrations
make db-seed            # Populate database with test data

# Deployment workflow
make docker-build       # Build Docker container image
make deploy             # Deploy to production environment
make logs               # View production application logs
```

### Local Development Setup

1. **Clone and initialize:**
   ```bash
   git clone <repository-url>
   cd cannanote/backend
   cp .env.example .env
   # Configure environment variables in .env
   ```

2. **Start local database:**
   ```bash
   make docker-run
   # Starts PostgreSQL container with development schema
   ```

3. **Initialize database schema:**
   ```bash
   cd ../supabase
   supabase db reset
   cd ../backend
   ```

4. **Start development server:**
   ```bash
   make dev
   # Starts server with hot reload at http://localhost:3001
   ```

### Testing Strategy

#### Unit Testing
- **Domain Logic** — Pure function testing for business rules
- **Application Services** — Use case testing with mocked dependencies
- **Repository Implementations** — Data persistence testing with test database

#### Integration Testing
- **HTTP Handlers** — End-to-end API testing with real database
- **Template Rendering** — UI component testing with sample data
- **Database Operations** — Schema validation and query performance

#### Testing Commands
```bash
make test               # Run all test suites
make test-unit          # Unit tests only
make test-integration   # Integration tests only
make test-coverage      # Generate coverage reports
```

## Planned Feature Development

### Phase 1: Core Cannabis Tracking

#### Entry Management Domain
- **Domain Entities** — Cannabis experience entry modeling with consumption methods
- **Business Logic** — Logging workflows, pattern recognition, and data validation
- **Storage Implementation** — Optimized time-series data storage with PostgreSQL
- **API Development** — RESTful entry CRUD operations with filtering and search
- **Template Integration** — Entry forms and list views with HTMX interactions

#### Data Privacy Enhancements
- **Local-First Architecture** — Client-side data storage with server synchronization
- **Encryption at Rest** — Field-level encryption for sensitive consumption data
- **Consent Management** — Granular privacy controls and data sharing preferences
- **Data Portability** — Export functionality for user data ownership

### Phase 2: Intelligence and Insights

#### Pattern Recognition Engine
- **Analytics Service** — Statistical analysis of consumption patterns and effects
- **Machine Learning Pipeline** — Recommendation engine for strain and dosage optimization
- **Correlation Analysis** — Environmental factors and consumption outcome analysis
- **Visualization Components** — Chart and graph rendering for pattern insights

#### Advanced Features
- **Calendar Integration** — Microsoft Graph API for lifestyle correlation analysis
- **Notification System** — Personalized reminders and harm reduction alerts
- **Social Features** — Optional community sharing with privacy controls
- **Content Management** — Educational content system with peer-reviewed research

### Phase 3: Platform Integration

#### External Service Integration
- **Third-Party APIs** — Integration with health tracking platforms
- **Webhook System** — Event-driven architecture for real-time data updates
- **GraphQL API** — Flexible query interface for mobile and web clients
- **Microservice Architecture** — Domain extraction for independent scaling

#### Enterprise Capabilities
- **Multi-Tenancy** — Platform support for white-label implementations
- **API Rate Limiting** — Request throttling and quota management
- **Audit Logging** — Comprehensive activity tracking for compliance
- **Monitoring and Observability** — Application performance monitoring and alerting

## Security Considerations

### Data Protection
- **Input Validation** — Comprehensive sanitization for all user inputs
- **SQL Injection Prevention** — Parameterized queries and prepared statements
- **XSS Protection** — Template output encoding and content security policies
- **CSRF Protection** — Token-based request validation for state changes

### Authentication and Authorization
- **JWT Token Validation** — Supabase Auth integration with token verification
- **Row Level Security** — Database-level access control for multi-user data
- **Permission System** — Role-based access control for administrative functions
- **Session Management** — Secure session handling with automatic expiration

### Infrastructure Security
- **HTTPS Enforcement** — TLS encryption for all production communications
- **Environment Isolation** — Secure configuration management and secret handling
- **Container Security** — Distroless base images and minimal attack surface
- **Dependency Management** — Regular security updates and vulnerability scanning

## Performance Optimization

### Application Performance
- **Connection Pooling** — Database connection optimization for concurrent requests
- **Template Caching** — Compiled template caching for reduced rendering overhead
- **Static Asset Optimization** — Compression and caching headers for media files
- **Query Optimization** — Index strategy and query performance monitoring

### Monitoring and Observability
- **Health Check Endpoints** — Comprehensive system health monitoring
- **Metrics Collection** — Application performance metrics and business KPIs
- **Error Tracking** — Centralized error logging and alerting system
- **Performance Profiling** — Go pprof integration for performance analysis

## Deployment Architecture

### Production Environment
- **Container Orchestration** — Docker containers deployed on Fly.io platform
- **Database Management** — Managed Supabase PostgreSQL with automatic backups
- **CDN Integration** — Global content delivery for static assets
- **Load Balancing** — Automatic request distribution and failover

### CI/CD Pipeline
- **Automated Testing** — Complete test suite execution on code changes
- **Code Quality Gates** — Linting, formatting, and security scanning
- **Deployment Automation** — Zero-downtime deployments with health checks
- **Rollback Capabilities** — Automated rollback on deployment failures

## Contributing Guidelines

### Code Standards
- **Go Conventions** — Follow standard Go formatting and naming conventions
- **Architecture Compliance** — Maintain hexagonal architecture boundaries
- **Test Coverage** — Minimum 80% test coverage for new features
- **Documentation** — Comprehensive documentation for public interfaces

### Development Workflow
- **Feature Branches** — Isolated development with pull request reviews
- **Code Reviews** — Mandatory review process for all code changes
- **Integration Testing** — Full test suite execution before merge
- **Documentation Updates** — Synchronize documentation with code changes

This backend service provides the foundation for CannaNote's personal cannabis journaling platform, emphasizing privacy protection, performance optimization, and maintainable architecture for long-term platform evolution.