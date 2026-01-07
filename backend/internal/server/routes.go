package server

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
	"github.com/supabase-community/gotrue-go/types"

	"backend/cmd/web"
	httpHandlers "backend/internal/adapters/http"
	"github.com/a-h/templ"
)

// Email validation regex - RFC 5322 compliant but practical
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Rate limiting for beta signup
type rateLimiter struct {
	mu       sync.RWMutex
	attempts map[string][]time.Time
}

var signupRateLimiter = &rateLimiter{
	attempts: make(map[string][]time.Time),
}

// checkRateLimit checks if IP has exceeded rate limit (max 3 attempts per 15 minutes)
func (rl *rateLimiter) checkRateLimit(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	now := time.Now()
	cutoff := now.Add(-15 * time.Minute) // 15 minute window
	
	// Clean old attempts
	if attempts, exists := rl.attempts[ip]; exists {
		validAttempts := []time.Time{}
		for _, attempt := range attempts {
			if attempt.After(cutoff) {
				validAttempts = append(validAttempts, attempt)
			}
		}
		rl.attempts[ip] = validAttempts
	}
	
	// Check if too many attempts
	if len(rl.attempts[ip]) >= 3 {
		return false // Rate limit exceeded
	}
	
	// Add current attempt
	rl.attempts[ip] = append(rl.attempts[ip], now)
	return true // Rate limit OK
}

// validateAndSanitizeEmail validates and sanitizes email input
func validateAndSanitizeEmail(email string) (string, error) {
	log.Printf("DEBUG: validateAndSanitizeEmail input: '%s' (length: %d)", email, len(email))
	
	// Trim whitespace only - no HTML escaping for email addresses
	email = strings.TrimSpace(email)
	log.Printf("DEBUG: After trimming: '%s' (length: %d)", email, len(email))
	
	// Check length limits (RFC 5321: 320 chars max, but 254 is practical)
	if len(email) > 254 {
		return "", fmt.Errorf("email address too long (max 254 characters)")
	}
	
	// Check minimum length (a@b.co = 6 chars minimum)
	if len(email) < 6 {
		log.Printf("DEBUG: Email too short - length is %d, minimum is 6", len(email))
		return "", fmt.Errorf("email address too short")
	}
	
	// Validate format first
	if !emailRegex.MatchString(email) {
		return "", fmt.Errorf("invalid email format")
	}
	
	// Convert to lowercase for consistency
	email = strings.ToLower(email)
	
	// Additional security: check for suspicious patterns that could indicate injection attempts
	suspicious := []string{"<script", "javascript:", "data:", "vbscript:", "onload=", "';", "--", "/*"}
	emailLower := strings.ToLower(email)
	for _, pattern := range suspicious {
		if strings.Contains(emailLower, pattern) {
			return "", fmt.Errorf("invalid email format")
		}
	}
	
	return email, nil
}

// validateConsent validates the consent checkbox input
func validateConsent(consent string) error {
	// Sanitize input
	consent = strings.TrimSpace(html.EscapeString(consent))
	
	// HTML checkbox only sends "on" when checked, or empty string when unchecked
	if consent != "on" {
		return fmt.Errorf("privacy consent is required")
	}
	
	return nil
}

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	
	// Add custom debug middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		log.Printf("DEBUG: Incoming request: %s %s from %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
		c.Next()
		log.Printf("DEBUG: Request completed: %s %s -> %d", c.Request.Method, c.Request.URL.Path, c.Writer.Status())
	})

	// Landing page
	r.GET("/", func(c *gin.Context) {
		templ.Handler(web.Landing()).ServeHTTP(c.Writer, c.Request)
	})

	// Privacy page
	r.GET("/privacy", func(c *gin.Context) {
		templ.Handler(web.Privacy()).ServeHTTP(c.Writer, c.Request)
	})

	// Terms of Service page
	r.GET("/terms", func(c *gin.Context) {
		templ.Handler(web.Terms()).ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/health", httpHandlers.HealthHandler)

	// Cache headers middleware for static assets
	r.Use(func(c *gin.Context) {
		// Add cache headers for static assets (CSS, JS, fonts, images)
		if len(c.Request.URL.Path) > 8 && c.Request.URL.Path[:8] == "/assets/" {
			// Cache CSS and JS for 1 year (they have content hashes)
			if filepath.Ext(c.Request.URL.Path) == ".css" || filepath.Ext(c.Request.URL.Path) == ".js" {
				c.Header("Cache-Control", "public, max-age=31536000, immutable")
			} else {
				// Cache other assets for 30 days
				c.Header("Cache-Control", "public, max-age=2592000")
			}
		}
		c.Next()
	})

	// Static assets served with cache headers
	r.Static("/assets", "./cmd/web/assets")

	r.GET("/web", func(c *gin.Context) {
		templ.Handler(web.HelloForm()).ServeHTTP(c.Writer, c.Request)
	})

	r.POST("/hello", func(c *gin.Context) {
		web.HelloWebHandler(c.Writer, c.Request)
	})

	// Documentation routes
	contentDir := filepath.Join("..", "docs", "content")
	docsHandler, err := httpHandlers.NewDocsHandler(contentDir)
	
	// Always add fallback route
	r.GET("/docs-fallback", func(c *gin.Context) {
		templ.Handler(web.Docs()).ServeHTTP(c.Writer, c.Request)
	})
	
	if err != nil {
		log.Printf("Failed to initialize docs handler: %v", err)
		// Fallback to old docs page
		r.GET("/docs", func(c *gin.Context) {
			templ.Handler(web.Docs()).ServeHTTP(c.Writer, c.Request)
		})
	} else {
		// Dynamic docs routing - try the new handler
		r.GET("/docs/*path", docsHandler.HandleDocsRequest)
		r.GET("/docs", docsHandler.HandleDocsRequest)
		
		// Search endpoint
		r.GET("/api/docs/search", docsHandler.HandleDocsSearch)
	}

	// Legacy route for existing cannabinoids handler (will be migrated)
	r.GET("/learn/cannabinoids", func(c *gin.Context) {
		web.CannabinoidsHandler(c)
	})

	// Pricing page
	r.GET("/pricing", func(c *gin.Context) {
		templ.Handler(web.Pricing()).ServeHTTP(c.Writer, c.Request)
	})

	// Beta signup API
	r.POST("/api/beta-signup", s.BetaSignupHandler)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) BetaSignupHandler(c *gin.Context) {
	// Check rate limit first
	clientIP := c.ClientIP()
	if !signupRateLimiter.checkRateLimit(clientIP) {
		log.Printf("SECURITY: Rate limit exceeded for IP %s", clientIP)
		templ.Handler(web.BetaSignupError("Too many signup attempts. Please try again in 15 minutes.")).ServeHTTP(c.Writer, c.Request)
		return
	}
	
	// Get raw input
	emailRaw := c.PostForm("email")
	consentRaw := c.PostForm("consent")
	
	// Debug logging to identify the issue
	log.Printf("DEBUG: Raw email input: '%s' (length: %d)", emailRaw, len(emailRaw))
	log.Printf("DEBUG: Raw consent input: '%s'", consentRaw)
	
	// Validate and sanitize email
	email, err := validateAndSanitizeEmail(emailRaw)
	if err != nil {
		log.Printf("SECURITY: Invalid email input from %s: %v", c.ClientIP(), err)
		templ.Handler(web.BetaSignupError("Please enter a valid email address")).ServeHTTP(c.Writer, c.Request)
		return
	}
	
	// Validate consent
	err = validateConsent(consentRaw)
	if err != nil {
		log.Printf("SECURITY: Invalid consent input from %s: %v", c.ClientIP(), err)
		templ.Handler(web.BetaSignupError("Privacy consent is required to join the beta")).ServeHTTP(c.Writer, c.Request)
		return
	}
	
	// Create beta user profile
	err = s.createBetaUser(email)
	if err != nil {
		// Check if it's a duplicate email error
		if strings.Contains(err.Error(), "already registered") {
			log.Printf("INFO: Duplicate signup attempt: %s from %s", email, c.ClientIP())
			templ.Handler(web.BetaSignupError("This email is already registered for beta access")).ServeHTTP(c.Writer, c.Request)
			return
		}
		
		// Generic error for other failures
		log.Printf("ERROR: Failed to create beta user %s: %v", email, err)
		templ.Handler(web.BetaSignupError("Unable to process signup. Please try again.")).ServeHTTP(c.Writer, c.Request)
		return
	}
	
	log.Printf("INFO: Beta user created successfully: %s from %s", email, c.ClientIP())
	templ.Handler(web.BetaSignupSuccess(email)).ServeHTTP(c.Writer, c.Request)
}

// createBetaUser creates a Supabase auth user and corresponding profile
func (s *Server) createBetaUser(email string) error {
	// Initialize Supabase client with service role key for admin operations
	supabaseURL := os.Getenv("SUPABASE_URL")
	serviceKey := os.Getenv("SUPABASE_SECRET_KEY")
	
	if supabaseURL == "" || serviceKey == "" {
		return fmt.Errorf("missing Supabase configuration")
	}
	
	client, err := supabase.NewClient(supabaseURL, serviceKey, &supabase.ClientOptions{})
	if err != nil {
		return fmt.Errorf("failed to create Supabase client: %w", err)
	}
	
	// Check if user already exists by email in auth.users
	db := s.db.GetDB()
	var existingUserID string
	err = db.QueryRow("SELECT id FROM auth.users WHERE email = $1", email).Scan(&existingUserID)
	if err == nil {
		// User already exists
		return fmt.Errorf("email %s already registered for beta", email)
	}
	
	// Create Supabase auth user with email verification required
	log.Printf("Creating Supabase auth user for: %s", email)
	
	// Use the Auth API to create user - this will send verification email
	authUser, err := client.Auth.AdminCreateUser(types.AdminCreateUserRequest{
		Email:        email,
		EmailConfirm: false, // Require email verification for security
	})
	if err != nil {
		return fmt.Errorf("failed to create Supabase auth user: %w", err)
	}
	
	// Create corresponding profile using the real auth user ID
	_, err = db.Exec(`
		INSERT INTO profiles (id, email, beta_status, beta_joined_at, created_at) 
		VALUES ($1, $2, 'waitlist', NOW(), NOW())
	`, authUser.ID, email)
	
	if err != nil {
		log.Printf("ERROR: Failed to create profile for auth user %s: %v", authUser.ID, err)
		// TODO: Consider cleaning up auth user if profile creation fails
		return fmt.Errorf("failed to create user profile: %w", err)
	}
	
	log.Printf("Successfully created Supabase auth user and profile: %s (ID: %s)", email, authUser.ID)
	return nil
}
