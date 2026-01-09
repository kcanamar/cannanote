package server

import (
	"fmt"
	"html"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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

// SecurityHeaders middleware adds security headers to all responses
func SecurityHeaders() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		
		// Prevent embedding in frames (clickjacking protection)
		c.Header("X-Frame-Options", "DENY")
		
		// XSS protection (legacy, but still useful)
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// Force HTTPS (HSTS)
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		
		// Content Security Policy - development friendly but still secure
		appEnv := os.Getenv("APP_ENV")
		var csp string
		if appEnv == "production" {
			// Production CSP with external scripts only and Google Fonts support
			csp = "default-src 'self'; " +
				"script-src 'self'; " +
				"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; " +
				"img-src 'self' data: https:; " +
				"font-src 'self' https://fonts.gstatic.com; " +
				"connect-src 'self'; " +
				"form-action 'self'; " +
				"frame-ancestors 'none'; " +
				"base-uri 'self'"
		} else {
			// Relaxed CSP for development
			csp = "default-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
				"img-src 'self' data: https: blob:; " +
				"font-src 'self' data: https://fonts.gstatic.com; " +
				"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; " +
				"connect-src 'self' ws: wss:; " +
				"form-action 'self'"
		}
		c.Header("Content-Security-Policy", csp)
		
		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Remove server information
		c.Header("Server", "")
		
		c.Next()
	})
}

// getTrustedClientIP extracts the real client IP from trusted sources
func getTrustedClientIP(c *gin.Context) string {
	// Check Cloudflare header first
	if ip := c.GetHeader("CF-Connecting-IP"); ip != "" && isValidIP(ip) {
		return ip
	}
	
	// Check X-Real-IP (common proxy header)
	if ip := c.GetHeader("X-Real-IP"); ip != "" && isValidIP(ip) {
		return ip
	}
	
	// Check X-Forwarded-For (may contain multiple IPs)
	if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
		// Take the first IP (client IP)
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if isValidIP(ip) {
				return ip
			}
		}
	}
	
	// Fall back to Gin's default (direct connection)
	return c.ClientIP()
}

// isValidIP validates an IP address
func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// RequestSizeLimit middleware limits request body size
func RequestSizeLimit(maxSize int64) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	})
}

// ConditionalLogging only logs debug info in development
func conditionalLog(level, format string, args ...interface{}) {
	debug := os.Getenv("DEBUG")
	appEnv := os.Getenv("APP_ENV")
	
	// Only log debug messages in development or when DEBUG=true
	if level == "DEBUG" && debug != "true" && appEnv == "production" {
		return
	}
	
	// Always log non-debug messages
	log.Printf("["+level+"] "+format, args...)
}

// Rate limiting for beta signup
type rateLimiter struct {
	mu           sync.RWMutex
	attempts     map[string][]time.Time
	maxAttempts  int
	windowMinutes int
	cleanupTimer *time.Timer
}

// newRateLimiter creates a rate limiter with configurable limits
func newRateLimiter(maxAttempts, windowMinutes int) *rateLimiter {
	rl := &rateLimiter{
		attempts:     make(map[string][]time.Time),
		maxAttempts:  maxAttempts,
		windowMinutes: windowMinutes,
	}
	
	// Auto-cleanup every hour to prevent memory bloat
	rl.cleanupTimer = time.AfterFunc(time.Hour, func() {
		rl.cleanup()
		rl.cleanupTimer.Reset(time.Hour)
	})
	
	return rl
}

// Get rate limit settings from environment (with defaults)
func getRateLimitSettings() (int, int) {
	maxAttempts := 3 // default
	windowMinutes := 15 // default
	
	if envMax := os.Getenv("RATE_LIMIT_MAX_ATTEMPTS"); envMax != "" {
		if parsed, err := strconv.Atoi(envMax); err == nil && parsed > 0 {
			maxAttempts = parsed
		}
	}
	
	if envWindow := os.Getenv("RATE_LIMIT_WINDOW_MINUTES"); envWindow != "" {
		if parsed, err := strconv.Atoi(envWindow); err == nil && parsed > 0 {
			windowMinutes = parsed
		}
	}
	
	return maxAttempts, windowMinutes
}

var signupRateLimiter = func() *rateLimiter {
	maxAttempts, windowMinutes := getRateLimitSettings()
	return newRateLimiter(maxAttempts, windowMinutes)
}()

// checkRateLimit checks if IP has exceeded rate limit (max 3 attempts per 15 minutes)
func (rl *rateLimiter) checkRateLimit(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	now := time.Now()
	cutoff := now.Add(-time.Duration(rl.windowMinutes) * time.Minute)
	
	// Clean old attempts for this IP
	if attempts, exists := rl.attempts[ip]; exists {
		validAttempts := []time.Time{}
		for _, attempt := range attempts {
			if attempt.After(cutoff) {
				validAttempts = append(validAttempts, attempt)
			}
		}
		rl.attempts[ip] = validAttempts
		
		// Remove empty entries to save memory
		if len(validAttempts) == 0 {
			delete(rl.attempts, ip)
		}
	}
	
	// Check if too many attempts
	if len(rl.attempts[ip]) >= rl.maxAttempts {
		conditionalLog("SECURITY", "Rate limit exceeded for IP %s: %d/%d attempts in %d minutes", 
			ip, len(rl.attempts[ip]), rl.maxAttempts, rl.windowMinutes)
		return false // Rate limit exceeded
	}
	
	// Add current attempt
	rl.attempts[ip] = append(rl.attempts[ip], now)
	return true // Allow request
}

// cleanup removes expired entries to prevent memory bloat
func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	now := time.Now()
	cutoff := now.Add(-time.Duration(rl.windowMinutes) * time.Minute)
	
	for ip, attempts := range rl.attempts {
		validAttempts := []time.Time{}
		for _, attempt := range attempts {
			if attempt.After(cutoff) {
				validAttempts = append(validAttempts, attempt)
			}
		}
		
		if len(validAttempts) == 0 {
			delete(rl.attempts, ip)
		} else {
			rl.attempts[ip] = validAttempts
		}
	}
	
	conditionalLog("DEBUG", "Rate limiter cleanup completed: %d active IPs", len(rl.attempts))
}

// getRateLimitStatus returns current status for monitoring
func (rl *rateLimiter) getRateLimitStatus(ip string) (int, int) {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	
	attempts := len(rl.attempts[ip])
	return attempts, rl.maxAttempts
}

// validateAndSanitizeEmail validates and sanitizes email input
func validateAndSanitizeEmail(email string) (string, error) {
	conditionalLog("DEBUG", "validateAndSanitizeEmail input: '%s' (length: %d)", email, len(email))
	
	// Trim whitespace only - no HTML escaping for email addresses
	email = strings.TrimSpace(email)
	conditionalLog("DEBUG", "After trimming: '%s' (length: %d)", email, len(email))
	
	// Check length limits (RFC 5321: 320 chars max, but 254 is practical)
	if len(email) > 254 {
		return "", fmt.Errorf("email address too long (max 254 characters)")
	}
	
	// Check minimum length (a@b.co = 6 chars minimum)
	if len(email) < 6 {
		conditionalLog("DEBUG", "Email too short - length is %d, minimum is 6", len(email))
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
	
	// Apply security middleware (production-aware)
	r.Use(SecurityHeaders())
	r.Use(RequestSizeLimit(8 << 20)) // 8MB request limit
	
	// Add custom logging middleware (conditional debug)
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		clientIP := getTrustedClientIP(c)
		conditionalLog("DEBUG", "Incoming request: %s %s from %s", c.Request.Method, c.Request.URL.Path, clientIP)
		c.Next()
		conditionalLog("DEBUG", "Request completed: %s %s -> %d", c.Request.Method, c.Request.URL.Path, c.Writer.Status())
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

	// Documentation routes - check multiple possible locations
	var contentDir string
	possiblePaths := []string{
		filepath.Join("..", "docs", "content"),     // Local development
		filepath.Join("docs", "content"),          // Fly.io deployment
		filepath.Join("/app", "docs", "content"),  // Docker deployment
	}
	
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			contentDir = path
			break
		}
	}
	
	if contentDir == "" {
		conditionalLog("ERROR", "Could not find docs content directory in any of these locations: %v", possiblePaths)
	}
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
	// Use secure IP detection for rate limiting
	clientIP := getTrustedClientIP(c)
	
	// Check rate limit first
	if !signupRateLimiter.checkRateLimit(clientIP) {
		attempts, maxAttempts := signupRateLimiter.getRateLimitStatus(clientIP)
		windowMinutes := signupRateLimiter.windowMinutes
		
		// More informative error message
		errorMsg := fmt.Sprintf("Too many signup attempts (%d/%d). Please try again in %d minutes.", 
			attempts, maxAttempts, windowMinutes)
		templ.Handler(web.BetaSignupError(errorMsg)).ServeHTTP(c.Writer, c.Request)
		return
	}
	
	// Get raw input
	emailRaw := c.PostForm("email")
	consentRaw := c.PostForm("consent")
	
	// Validate input sizes to prevent DoS
	if len(emailRaw) > 254 { // RFC 5321 email length limit
		conditionalLog("SECURITY", "Oversized email input from %s: %d characters", clientIP, len(emailRaw))
		templ.Handler(web.BetaSignupError("Email address is too long")).ServeHTTP(c.Writer, c.Request)
		return
	}
	
	// Debug logging (only in development)
	conditionalLog("DEBUG", "Raw email input: '%s' (length: %d)", emailRaw, len(emailRaw))
	conditionalLog("DEBUG", "Raw consent input: '%s'", consentRaw)
	
	// Validate and sanitize email
	email, err := validateAndSanitizeEmail(emailRaw)
	if err != nil {
		conditionalLog("SECURITY", "Invalid email input from %s: %v", clientIP, err)
		templ.Handler(web.BetaSignupError("Please enter a valid email address")).ServeHTTP(c.Writer, c.Request)
		return
	}
	
	// Validate consent
	err = validateConsent(consentRaw)
	if err != nil {
		conditionalLog("SECURITY", "Invalid consent input from %s: %v", clientIP, err)
		templ.Handler(web.BetaSignupError("Privacy consent is required to join the beta")).ServeHTTP(c.Writer, c.Request)
		return
	}
	
	// Create beta user profile
	err = s.createBetaUser(email)
	if err != nil {
		// Check if it's a duplicate email error
		if strings.Contains(err.Error(), "already registered") {
			conditionalLog("INFO", "Duplicate signup attempt: %s from %s", email, clientIP)
			templ.Handler(web.BetaSignupError("This email is already registered for beta access")).ServeHTTP(c.Writer, c.Request)
			return
		}
		
		// Generic error for other failures
		conditionalLog("ERROR", "Failed to create beta user %s: %v", email, err)
		templ.Handler(web.BetaSignupError("Unable to process signup. Please try again.")).ServeHTTP(c.Writer, c.Request)
		return
	}
	
	conditionalLog("INFO", "Beta user created successfully: %s from %s", email, clientIP)
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
	conditionalLog("INFO", "Creating Supabase auth user for: %s", email)
	
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
		conditionalLog("ERROR", "Failed to create profile for auth user %s: %v", authUser.ID, err)
		// TODO: Consider cleaning up auth user if profile creation fails
		return fmt.Errorf("failed to create user profile: %w", err)
	}
	
	conditionalLog("INFO", "Successfully created Supabase auth user and profile: %s (ID: %s)", email, authUser.ID)
	return nil
}
