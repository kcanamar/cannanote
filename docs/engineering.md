# CannaNote Engineering Guide

## Quick Start for New Developers

This guide will get you from zero to contributing code to CannaNote. Read this first if you're joining the team.

### Prerequisites

**Required:**
- Go 1.24+
- Git
- Text editor/IDE

**Recommended:**
- Docker (for local development)
- Make (for build commands)

### Getting Started

1. **Clone and run locally:**
   ```bash
   git clone [repo-url]
   cd cannanote/backend
   make dev
   ```

2. **Access the application:**
   - Local app: http://localhost:3001
   - Health check: http://localhost:3001/health

3. **Make a change and see it live:**
   - Edit any `.templ` file
   - Edit any `.go` file  
   - Changes auto-reload via air

## Architecture Overview

### **Philosophy: Simple, Fast, Reliable**

We prioritize:
1. **Developer velocity** - Fast to develop and deploy
2. **User experience** - Quick loading, responsive interactions
3. **Maintenance** - Easy to understand and modify

### **Why Go + templ + HTMX**

**Go** - Fast compilation, excellent performance, simple deployment
**templ** - Type-safe HTML templates that compile to Go code
**HTMX** - Server-side rendering with dynamic interactions, no complex JS

This stack lets us build rich web applications with minimal complexity.

### Current Technology Stack

```go
// Backend
Go 1.24+                    // Main language
Gin router                  // HTTP routing
templ                      // Type-safe templates

// Frontend  
HTMX                       // Dynamic interactions
Tailwind CSS               // Utility-first styling
Alpine.js (minimal)        // Lightweight JS when needed

// Database & Auth
Supabase PostgreSQL        // Database with built-in auth
Supabase Auth              // OAuth + session management

// Infrastructure
Fly.io                     // Deployment platform
Docker                     // Containerization
```

## Project Structure

```
cannanote/
├── backend/                # Main Go application
│   ├── cmd/
│   │   └── api/
│   │       └── main.go    # Application entry point
│   │   └── web/           # templ templates
│   ├── internal/          # Private application code
│   │   ├── adapters/      # External integrations
│   │   │   ├── http/      # HTTP handlers
│   │   │   └── repository/# Data access layer
│   │   ├── core/          # Business logic
│   │   │   ├── domain/    # Entities and business rules  
│   │   │   ├── ports/     # Interfaces
│   │   │   └── application/# Use cases
│   │   ├── database/      # Database connection
│   │   └── server/        # Server setup and routing
│   ├── certs/             # SSL certificates
│   ├── Dockerfile         # Container definition
│   ├── fly.toml          # Deployment configuration
│   ├── Makefile          # Build commands
│   └── go.mod            # Go dependencies
├── supabase/              # Database schema and config
├── docs/                  # Documentation
└── README.md             # Project overview
```

### **Hexagonal Architecture**

We use hexagonal (ports and adapters) architecture:

- **Core Domain** (`internal/core/`) - Pure business logic, no external dependencies
- **Ports** (`internal/core/ports/`) - Interfaces that define what we need from the outside world
- **Adapters** (`internal/adapters/`) - Concrete implementations that connect to external services

This makes the code easy to test and modify.

## Development Workflow

### **Daily Development**

```bash
# Start development server with hot reload
make dev

# Run tests
make test

# Run linting
make lint

# Build for production  
make build

# Deploy to production
make deploy
```

### **Adding a New Feature**

1. **Define the domain model** in `internal/core/domain/`
2. **Create the port interface** in `internal/core/ports/`
3. **Implement the use case** in `internal/core/application/`
4. **Create the adapter** in `internal/adapters/`
5. **Add HTTP handlers** in `internal/adapters/http/`
6. **Create templ templates** in `cmd/web/`
7. **Add routes** in `internal/server/routes.go`
8. **Write tests**

### **Code Standards**

**Go Code:**
- Use `gofmt` for formatting
- Follow standard Go naming conventions
- Keep functions small and focused
- Write tests for business logic

**templ Templates:**
- Use semantic HTML
- Follow Tailwind CSS conventions
- Keep templates focused on presentation
- Use HTMX attributes for interactions

**Database:**
- Use meaningful table and column names
- Add appropriate indexes
- Use foreign keys for relationships
- Include created_at/updated_at timestamps

## Key Concepts

### **templ Templates**

templ compiles to Go code, giving us type safety:

```go
// hello.templ
package web

templ HelloPage(name string) {
    <html>
        <body>
            <h1>Hello, { name }!</h1>
        </body>
    </html>
}
```

Generates Go code you can call from handlers:

```go
func HelloHandler(c *gin.Context) {
    name := c.Query("name")
    component := HelloPage(name)
    c.Header("Content-Type", "text/html")
    component.Render(c.Request.Context(), c.Writer)
}
```

### **HTMX Interactions**

HTMX lets us add dynamic behavior without writing JavaScript:

```html
<button 
    hx-post="/api/entries"
    hx-target="#entries-list" 
    hx-swap="afterbegin">
    Add Entry
</button>
```

This makes an AJAX POST request and inserts the response at the top of the entries list.

### **Database Integration**

We use Supabase PostgreSQL with direct SQL queries:

```go
type Human struct {
    ID       string `json:"id" db:"id"`
    Username string `json:"username" db:"username"`  
    Email    string `json:"email" db:"email"`
}

func (r *supabaseHumanRepository) GetByID(ctx context.Context, id string) (*Human, error) {
    var human Human
    query := "SELECT id, username, email FROM humans WHERE id = $1"
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &human.ID, &human.Username, &human.Email,
    )
    return &human, err
}
```

### **Authentication**

We use Supabase Auth for user management:

```go
// Check if user is authenticated
func RequireAuth(next gin.HandlerFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.Redirect(http.StatusTemporaryRedirect, "/login")
            return
        }
        // Validate token with Supabase
        next(c)
    }
}
```

## Testing

### **Test Structure**

- Unit tests for business logic in `internal/core/`
- Integration tests for adapters
- HTTP tests for handlers

```go
func TestCreateHuman(t *testing.T) {
    service := application.NewHumanService(mockRepo)
    
    human := &domain.Human{
        Username: "testuser",
        Email:    "test@example.com",
    }
    
    result, err := service.CreateHuman(context.Background(), human)
    
    assert.NoError(t, err)
    assert.Equal(t, "testuser", result.Username)
}
```

### **Running Tests**

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Deployment

### **Production Deployment**

```bash
# Deploy to Fly.io
make deploy
```

This will:
1. Run tests and linting
2. Set environment variables from .env file
3. Build Docker image
4. Deploy to Fly.io
5. Run health checks

### **Environment Variables**

Required in `.env` file:

```bash
# Database
DB_HOST=db.citdskdmralncvjyybin.supabase.co
DB_PORT=5432  
DB_DATABASE=postgres
DB_USERNAME=postgres
DB_PASSWORD=your_password
DB_SCHEMA=public

# Supabase
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your_anon_key

# Application
PORT=3001
APP_ENV=local
GIN_MODE=debug
```

### **Health Monitoring**

The application includes comprehensive health checking:

- **Database connectivity** - Direct PostgreSQL connection test
- **API endpoints** - Supabase REST API tests
- **Response times** - Performance monitoring
- **Overall status** - Aggregated health percentage

Access health check at: `/health`

## Common Tasks

### **Adding a New Page**

1. Create templ template:
   ```go
   // cmd/web/new_page.templ
   templ NewPage() {
       @Base() {
           <h1>New Page</h1>
       }
   }
   ```

2. Generate Go code:
   ```bash
   make build  # Runs templ generate
   ```

3. Add handler:
   ```go
   func NewPageHandler(c *gin.Context) {
       component := NewPage()
       c.Header("Content-Type", "text/html")
       component.Render(c.Request.Context(), c.Writer)
   }
   ```

4. Add route:
   ```go
   // internal/server/routes.go
   r.GET("/new-page", NewPageHandler)
   ```

### **Adding Database Tables**

1. Add migration to Supabase:
   ```sql
   -- supabase/migrations/create_new_table.sql
   CREATE TABLE new_table (
       id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
       name VARCHAR(255) NOT NULL,
       created_at TIMESTAMP DEFAULT NOW()
   );
   ```

2. Update locally:
   ```bash
   cd supabase && supabase db reset
   ```

3. Deploy to production:
   ```bash
   cd supabase && supabase db push
   ```

### **Adding HTMX Interactions**

Use HTMX attributes in templates:

```go
templ EntryForm() {
    <form hx-post="/api/entries" hx-target="#entries" hx-swap="afterbegin">
        <input type="text" name="strain" />
        <button type="submit">Add Entry</button>
    </form>
}
```

Create corresponding handler:

```go
func CreateEntryHandler(c *gin.Context) {
    strain := c.PostForm("strain")
    // Create entry...
    
    // Return just the new entry HTML
    entry := EntryCard(newEntry)
    entry.Render(c.Request.Context(), c.Writer)
}
```

## Performance Guidelines

### **Database Queries**

- Use appropriate indexes
- Limit query results with LIMIT
- Use prepared statements
- Monitor query performance

### **Templates**

- Keep templates focused and small
- Minimize nested components
- Use HTMX to update only changed parts

### **HTMX Optimization**

- Use `hx-target` to update specific elements
- Use `hx-swap` strategies appropriately
- Implement proper loading states

## Security Considerations

### **Authentication**

- All sensitive routes require authentication
- Use Supabase RLS (Row Level Security) for data access
- Validate all user inputs

### **Data Handling**

- Sanitize user inputs
- Use parameterized queries to prevent SQL injection
- Validate data on both client and server sides

### **Environment**

- Keep secrets in environment variables
- Never commit sensitive data to git
- Use HTTPS in production

## Getting Help

### **Common Issues**

**Templates not updating:**
```bash
make build  # Regenerates templ templates
```

**Database connection issues:**
```bash
# Check health endpoint
curl http://localhost:3001/health
```

**HTMX not working:**
- Check browser network tab for failed requests
- Verify handler returns correct HTML
- Ensure HTMX attributes are correct

### **Useful Resources**

- [templ documentation](https://templ.guide/)
- [HTMX documentation](https://htmx.org/docs/)
- [Go documentation](https://golang.org/doc/)
- [Supabase documentation](https://supabase.com/docs)

### **Team Communication**

For questions about:
- **Architecture decisions** - Ask in team chat
- **Business logic** - Review domain models first
- **UI/UX** - Check style guide (docs/style-guide.md)
- **Deployment issues** - Check Fly.io logs (`make fly-logs`)

## Next Steps for New Developers

1. **Get the app running locally** - Follow Quick Start
2. **Make a small change** - Edit a template or add a route
3. **Read the codebase** - Start with `main.go` and follow the flow
4. **Look at existing tests** - Understand how we test features
5. **Pick up a small task** - Start with a simple feature or bug fix

This guide covers the essentials for contributing to CannaNote. The codebase is designed to be approachable - when in doubt, follow existing patterns and ask questions.