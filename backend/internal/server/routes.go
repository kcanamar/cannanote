package server

import (
	"context"
	"fmt"
	"html"
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

	"backend/cmd/web"
	"backend/internal/adapters/external"
	httpAdapters "backend/internal/adapters/http"
	"backend/internal/adapters/logging"
	"backend/internal/core/ports"
	"backend/internal/version"

	"github.com/a-h/templ"
)

var routesLog = logging.With("routes")

// conditionalLog provides backward compatibility while using the new logging system
// This maps the old level strings to proper logger methods
func conditionalLog(level, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	switch level {
	case "DEBUG":
		routesLog.Debug(msg)
	case "INFO":
		routesLog.Info(msg)
	case "WARN", "WARNING":
		routesLog.Warn(msg)
	case "ERROR", "SECURITY":
		routesLog.Warn(msg) // Security events as warnings for visibility
	default:
		routesLog.Info(msg)
	}
}

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
				"connect-src 'self' https://fonts.googleapis.com https://fonts.gstatic.com; " +
				"form-action 'self'; " +
				"frame-ancestors 'none'; " +
				"base-uri 'self'"
		} else {
			// Relaxed CSP for development
			csp = "default-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
				"img-src 'self' data: https: blob:; " +
				"font-src 'self' data: https://fonts.gstatic.com; " +
				"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; " +
				"connect-src 'self' ws: wss: https://fonts.googleapis.com https://fonts.gstatic.com; " +
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


// Rate limiting for beta signup
type rateLimiter struct {
	mu            sync.RWMutex
	attempts      map[string][]time.Time
	maxAttempts   int
	windowMinutes int
	cleanupTimer  *time.Timer
}

// newRateLimiter creates a rate limiter with configurable limits
func newRateLimiter(maxAttempts, windowMinutes int) *rateLimiter {
	rl := &rateLimiter{
		attempts:      make(map[string][]time.Time),
		maxAttempts:   maxAttempts,
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
	maxAttempts := 3   // default
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
	// Set GIN mode based on environment (explicit, not just env var)
	if version.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		routesLog.Info("GIN running in release mode")
	} else {
		gin.SetMode(gin.DebugMode)
		routesLog.Info("GIN running in debug mode")
	}

	// Use gin.New() for explicit middleware control
	r := gin.New()

	// Recovery middleware - always needed
	r.Use(gin.Recovery())

	// Request logging - only in development
	if version.IsDevelopment() {
		r.Use(gin.Logger())
		r.Use(func(c *gin.Context) {
			clientIP := getTrustedClientIP(c)
			routesLog.Debug("Incoming request",
				ports.F("method", c.Request.Method),
				ports.F("path", c.Request.URL.Path),
				ports.F("ip", clientIP),
			)
			c.Next()
			routesLog.Debug("Request completed",
				ports.F("method", c.Request.Method),
				ports.F("path", c.Request.URL.Path),
				ports.F("status", c.Writer.Status()),
			)
		})
	}

	// Apply security middleware (always)
	r.Use(SecurityHeaders())
	r.Use(RequestSizeLimit(8 << 20)) // 8MB request limit

	// Create auth middleware
	authMiddleware := httpAdapters.NewAuthMiddleware(s.db.GetDB())

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

	r.GET("/health", httpAdapters.HealthHandler)

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

	// Serve robots.txt at root level
	r.GET("/robots.txt", func(c *gin.Context) {
		c.File("./cmd/web/assets/robots.txt")
	})

	// Serve favicon.ico at root level (browsers request this by default)
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File("./cmd/web/assets/images/favicon/favicon.ico")
	})

	// Service worker - served dynamically with version injection
	r.GET("/sw.js", func(c *gin.Context) {
		swVersion := version.GetServiceWorkerVersion()
		swLog := routesLog.With("sw")

		swLog.Debug("Serving service worker",
			ports.F("version", swVersion),
			ports.F("env", version.Environment),
		)

		// Read the service worker template
		swContent, err := os.ReadFile("./cmd/web/assets/sw.js")
		if err != nil {
			swLog.Error("Failed to read service worker file", ports.F("error", err))
			c.String(http.StatusInternalServerError, "Failed to load service worker")
			return
		}

		// Replace version placeholder with actual version
		swJS := strings.Replace(string(swContent), "{{SW_VERSION}}", swVersion, 1)

		// Set headers for service worker
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Content-Type", "application/javascript; charset=utf-8")
		c.String(http.StatusOK, swJS)
	})

	// Offline fallback page
	r.GET("/offline", func(c *gin.Context) {
		templ.Handler(web.Offline()).ServeHTTP(c.Writer, c.Request)
	})

	// Documentation routes - check multiple possible locations
	var contentDir string
	possiblePaths := []string{
		filepath.Join("..", "docs", "content"),    // Local development
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
	docsHandler, err := httpAdapters.NewDocsHandler(contentDir)

	// Always add fallback route
	r.GET("/docs-fallback", func(c *gin.Context) {
		templ.Handler(web.Docs()).ServeHTTP(c.Writer, c.Request)
	})

	if err != nil {
		routesLog.Error("Failed to initialize docs handler", ports.F("error", err))
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

	// Email verification/activation handler
	r.GET("/activate", s.ActivateAccountHandler)

	// Onboarding (password setup after email verification)
	r.GET("/onboarding", s.OnboardingHandler)

	// Password reset handler (consumes reset token and sets new password)
	r.POST("/auth/reset-password", s.ResetPasswordHandler)

	// Protected routes (require authentication)
	protected := r.Group("/app")
	protected.Use(authMiddleware.RequireSession())
	{
		protected.GET("", s.DashboardHandler)
		protected.GET("/", s.DashboardHandler)

		// Session routes
		protected.GET("/sessions", s.SessionsHandler)
		protected.GET("/sessions/new", s.SessionLogHandler)
		protected.GET("/sessions/:id", s.SessionDetailHandler)
		protected.GET("/sessions/:id/checkin", s.SessionCheckInHandler)
	}

	// Login page for unauthenticated users
	r.GET("/login", s.LoginPageHandler)

	// Sign in handler (processes login form)
	r.POST("/auth/signin", s.SignInHandler)

	// Sign out handler
	r.POST("/auth/signout", s.SignOutHandler)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
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
	err = s.createBetaUser(c.Request.Context(), email)
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

// createBetaUser creates a user and sends welcome email with verification link
func (s *Server) createBetaUser(ctx context.Context, email string) error {
	conditionalLog("DEBUG", "Creating beta user: %s", email)

	// Create auth service
	authService := external.NewAuthService(s.db.GetDB())

	// Create user (this checks for duplicates)
	user, err := authService.CreateUser(ctx, email)
	if err != nil {
		return err
	}

	// Generate verification token for email
	verificationToken, err := authService.CreateEmailVerificationToken(ctx, user.ID)
	if err != nil {
		conditionalLog("ERROR", "Failed to create verification token: %v", err)
		// User created but token failed - they can request new verification later
		verificationToken = ""
	}

	// Build verification link
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:8080"
	}
	verificationLink := ""
	if verificationToken != "" {
		verificationLink = fmt.Sprintf("%s/activate?token=%s", appURL, verificationToken)
	}

	conditionalLog("DEBUG", "Verification link generated for %s", email)

	// Send welcome email via Resend (non-blocking)
	go func() {
		resendService := external.NewResendService()
		if err := resendService.SendWelcomeEmail(email, verificationLink); err != nil {
			conditionalLog("ERROR", "Failed to send welcome email to %s: %v", email, err)
		} else {
			conditionalLog("INFO", "Welcome email sent to: %s", email)
		}
	}()

	return nil
}

// ActivateAccountHandler handles email verification from the welcome email
func (s *Server) ActivateAccountHandler(c *gin.Context) {
	conditionalLog("DEBUG", "ActivateAccountHandler called")

	token := c.Query("token")
	if token == "" {
		conditionalLog("DEBUG", "No token provided")
		templ.Handler(web.ActivationError("Missing activation token. Please check your email for the correct link.")).ServeHTTP(c.Writer, c.Request)
		return
	}

	// Verify the email using our auth service
	authService := external.NewAuthService(s.db.GetDB())
	user, err := authService.VerifyEmail(c.Request.Context(), token)
	if err != nil {
		conditionalLog("ERROR", "Email verification failed: %v", err)
		templ.Handler(web.ActivationError("Failed to verify your email. The link may have expired.")).ServeHTTP(c.Writer, c.Request)
		return
	}

	conditionalLog("INFO", "Email verified for user: %s", user.Email)

	// Generate password reset token so user can set their password
	resetToken, err := authService.CreatePasswordResetToken(c.Request.Context(), user.ID)
	if err != nil {
		conditionalLog("ERROR", "Failed to create reset token: %v", err)
		// Still redirect to onboarding, they can request password reset manually
		c.Redirect(http.StatusFound, "/onboarding?verified=true&email="+user.Email)
		return
	}

	// Redirect to onboarding with the reset token
	redirectURL := fmt.Sprintf("/onboarding?verified=true&email=%s&reset_token=%s", user.Email, resetToken)
	conditionalLog("DEBUG", "Redirecting to onboarding")

	c.Redirect(http.StatusFound, redirectURL)
}

// OnboardingHandler displays the onboarding page for new users to set their password
func (s *Server) OnboardingHandler(c *gin.Context) {
	verified := c.Query("verified") == "true"
	email := c.Query("email")
	resetToken := c.Query("reset_token")

	conditionalLog("DEBUG", "OnboardingHandler: verified=%v, email=%s, token_present=%v", verified, email, resetToken != "")

	templ.Handler(web.Onboarding(verified, email, resetToken)).ServeHTTP(c.Writer, c.Request)
}

// ResetPasswordHandler handles the password reset form submission
func (s *Server) ResetPasswordHandler(c *gin.Context) {
	conditionalLog("DEBUG", "ResetPasswordHandler called")

	token := c.PostForm("token")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")

	// Validate inputs
	if token == "" {
		templ.Handler(web.PasswordResetError("Missing reset token. Please try the activation link again.")).ServeHTTP(c.Writer, c.Request)
		return
	}

	if password == "" || len(password) < 8 {
		templ.Handler(web.PasswordResetError("Password must be at least 8 characters.")).ServeHTTP(c.Writer, c.Request)
		return
	}

	if password != confirmPassword {
		templ.Handler(web.PasswordResetError("Passwords do not match.")).ServeHTTP(c.Writer, c.Request)
		return
	}

	// Reset password using our auth service
	authService := external.NewAuthService(s.db.GetDB())
	user, err := authService.SetPassword(c.Request.Context(), token, password)
	if err != nil {
		conditionalLog("ERROR", "Password reset failed: %v", err)
		templ.Handler(web.PasswordResetError("Failed to reset password. The link may have expired.")).ServeHTTP(c.Writer, c.Request)
		return
	}

	conditionalLog("INFO", "Password set for user: %s", user.Email)

	// Create session for the user (auto-login)
	sessionToken, err := authService.CreateSession(
		c.Request.Context(),
		user.ID,
		c.Request.UserAgent(),
		getTrustedClientIP(c),
	)
	if err != nil {
		conditionalLog("ERROR", "Failed to create session: %v", err)
		// Session creation failed, but password was reset - show success and let them log in manually
		templ.Handler(web.PasswordResetSuccess()).ServeHTTP(c.Writer, c.Request)
		return
	}

	// Set session cookie
	httpAdapters.SetSessionCookie(c, sessionToken)

	conditionalLog("DEBUG", "Session created, showing success page with redirect")

	// Show success page that will redirect to /app
	templ.Handler(web.PasswordResetSuccessWithRedirect()).ServeHTTP(c.Writer, c.Request)
}

// DashboardHandler displays the protected app dashboard for authenticated users
func (s *Server) DashboardHandler(c *gin.Context) {
	userID := httpAdapters.GetUserIDFromGinContext(c)
	conditionalLog("DEBUG", "DashboardHandler: User ID: %s", userID)

	// Get user profile from database
	var email string
	var betaStatus string
	db := s.db.GetDB()
	err := db.QueryRow("SELECT email, beta_status FROM profiles WHERE id = $1", userID).Scan(&email, &betaStatus)
	if err != nil {
		conditionalLog("ERROR", "Failed to get profile: %v", err)
		email = "Unknown"
		betaStatus = "unknown"
	}

	conditionalLog("DEBUG", "Rendering dashboard for %s (status: %s)", email, betaStatus)

	templ.Handler(web.Dashboard(userID, email, betaStatus)).ServeHTTP(c.Writer, c.Request)
}

// LoginPageHandler displays the login page for unauthenticated users
// If user is already authenticated, redirects to /app
func (s *Server) LoginPageHandler(c *gin.Context) {
	// Check if user already has a valid session
	if cookie, err := c.Cookie(httpAdapters.SessionCookieName); err == nil && cookie != "" {
		authService := external.NewAuthService(s.db.GetDB())
		if _, err := authService.ValidateSession(c.Request.Context(), cookie); err == nil {
			// User is already authenticated, redirect to app
			conditionalLog("DEBUG", "LoginPageHandler: User already authenticated, redirecting to /app")
			c.Redirect(http.StatusFound, "/app")
			return
		}
	}

	templ.Handler(web.LoginPage()).ServeHTTP(c.Writer, c.Request)
}

// SignInHandler processes the login form and creates a session
func (s *Server) SignInHandler(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	conditionalLog("DEBUG", "SignInHandler: Attempting sign in for %s", email)

	// Validate inputs
	if email == "" || password == "" {
		templ.Handler(web.LoginError("Please enter both email and password.")).ServeHTTP(c.Writer, c.Request)
		return
	}

	// Sign in with our auth service
	authService := external.NewAuthService(s.db.GetDB())
	user, err := authService.SignIn(c.Request.Context(), email, password)
	if err != nil {
		conditionalLog("DEBUG", "Sign in failed: %v", err)
		templ.Handler(web.LoginError("Invalid email or password.")).ServeHTTP(c.Writer, c.Request)
		return
	}

	conditionalLog("INFO", "Sign in successful for user: %s", user.Email)

	// Create session
	sessionToken, err := authService.CreateSession(
		c.Request.Context(),
		user.ID,
		c.Request.UserAgent(),
		getTrustedClientIP(c),
	)
	if err != nil {
		conditionalLog("ERROR", "Failed to create session: %v", err)
		templ.Handler(web.LoginError("Failed to create session. Please try again.")).ServeHTTP(c.Writer, c.Request)
		return
	}

	// Set session cookie
	httpAdapters.SetSessionCookie(c, sessionToken)

	conditionalLog("DEBUG", "Session created, redirecting to /app")

	// Redirect to dashboard
	c.Redirect(http.StatusFound, "/app")
}

// SignOutHandler destroys the session and redirects to home
func (s *Server) SignOutHandler(c *gin.Context) {
	conditionalLog("DEBUG", "SignOutHandler called")

	// Get session token from cookie
	if cookie, err := c.Cookie(httpAdapters.SessionCookieName); err == nil && cookie != "" {
		// Revoke session in database
		authService := external.NewAuthService(s.db.GetDB())
		if err := authService.RevokeSession(c.Request.Context(), cookie); err != nil {
			conditionalLog("ERROR", "Error revoking session: %v", err)
		}
	}

	// Clear session cookie
	httpAdapters.ClearSessionCookie(c)

	// Redirect to home
	c.Redirect(http.StatusFound, "/")
}

// SessionsHandler displays the session history page
func (s *Server) SessionsHandler(c *gin.Context) {
	// Get user email for display
	var email string
	userID := httpAdapters.GetUserIDFromGinContext(c)
	db := s.db.GetDB()
	err := db.QueryRow("SELECT email FROM profiles WHERE id = $1", userID).Scan(&email)
	if err != nil {
		email = "Unknown"
	}

	templ.Handler(web.Sessions(email, c.Request.URL.Path)).ServeHTTP(c.Writer, c.Request)
}

// SessionLogHandler displays the new session logging form
func (s *Server) SessionLogHandler(c *gin.Context) {
	var email string
	userID := httpAdapters.GetUserIDFromGinContext(c)
	db := s.db.GetDB()
	err := db.QueryRow("SELECT email FROM profiles WHERE id = $1", userID).Scan(&email)
	if err != nil {
		email = "Unknown"
	}

	templ.Handler(web.SessionLog(email, c.Request.URL.Path)).ServeHTTP(c.Writer, c.Request)
}

// SessionDetailHandler displays a single session's details
func (s *Server) SessionDetailHandler(c *gin.Context) {
	sessionID := c.Param("id")
	var email string
	userID := httpAdapters.GetUserIDFromGinContext(c)
	db := s.db.GetDB()
	err := db.QueryRow("SELECT email FROM profiles WHERE id = $1", userID).Scan(&email)
	if err != nil {
		email = "Unknown"
	}

	templ.Handler(web.SessionDetail(email, sessionID, c.Request.URL.Path)).ServeHTTP(c.Writer, c.Request)
}

// SessionCheckInHandler displays the check-in form for an active session
func (s *Server) SessionCheckInHandler(c *gin.Context) {
	sessionID := c.Param("id")
	var email string
	userID := httpAdapters.GetUserIDFromGinContext(c)
	db := s.db.GetDB()
	err := db.QueryRow("SELECT email FROM profiles WHERE id = $1", userID).Scan(&email)
	if err != nil {
		email = "Unknown"
	}

	templ.Handler(web.SessionCheckIn(email, sessionID, c.Request.URL.Path)).ServeHTTP(c.Writer, c.Request)
}
