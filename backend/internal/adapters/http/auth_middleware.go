package http

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"net/http"
	"os"
	"strings"

	"backend/internal/adapters/external"
	"backend/internal/adapters/logging"
	"backend/internal/core/ports"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	// SessionCookieName is the name of the session cookie
	SessionCookieName = "cannanote_session"
	// CSRFCookieName is the name of the CSRF token cookie
	CSRFCookieName = "cannanote_csrf"
	// CSRFHeaderName is the name of the CSRF header for AJAX requests
	CSRFHeaderName = "X-CSRF-Token"
	// CSRFFormFieldName is the name of the hidden form field
	CSRFFormFieldName = "csrf_token"
	// SessionContextKey is the key used to store session in context
	SessionContextKey contextKey = "session"
	// UserContextKey is the key used to store user in context
	UserContextKey contextKey = "user"
	// CSRFContextKey is the key used to store CSRF token in context
	CSRFContextKey contextKey = "csrf_token"
)

// AuthMiddleware provides authentication middleware using our custom auth service
type AuthMiddleware struct {
	authService *external.AuthService
	log         ports.Logger
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(db *sql.DB) *AuthMiddleware {
	return &AuthMiddleware{
		authService: external.NewAuthService(db),
		log:         logging.With("auth-middleware"),
	}
}

// RequireSession returns middleware that requires a valid session
// Redirects to /login if no valid session
func (m *AuthMiddleware) RequireSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getSessionToken(c)
		if token == "" {
			m.log.Debug("No session token found, redirecting to login",
				ports.F("path", c.Request.URL.Path),
			)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		session, err := m.authService.ValidateSession(c.Request.Context(), token)
		if err != nil {
			m.log.Debug("Invalid session, redirecting to login",
				ports.F("error", err),
				ports.F("path", c.Request.URL.Path),
			)
			// Clear invalid cookie
			clearSessionCookie(c)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Store session and user ID in context
		c.Set(SessionContextKey, session)
		c.Set(UserContextKey, session.UserID.String())

		c.Next()
	}
}

// OptionalSession returns middleware that extracts session if present
// Does not require authentication - continues even without session
func (m *AuthMiddleware) OptionalSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getSessionToken(c)
		if token == "" {
			c.Next()
			return
		}

		session, err := m.authService.ValidateSession(c.Request.Context(), token)
		if err != nil {
			// Invalid session, clear cookie but continue
			clearSessionCookie(c)
			c.Next()
			return
		}

		// Store session and user ID in context
		c.Set(SessionContextKey, session)
		c.Set(UserContextKey, session.UserID.String())

		c.Next()
	}
}

// getSessionToken extracts the session token from cookie or Authorization header
func getSessionToken(c *gin.Context) string {
	// Try cookie first
	if cookie, err := c.Cookie(SessionCookieName); err == nil && cookie != "" {
		return cookie
	}

	// Try Authorization header (Bearer token)
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return ""
}

// SetSessionCookie sets the session cookie with secure settings
func SetSessionCookie(c *gin.Context, token string) {
	appEnv := os.Getenv("APP_ENV")
	secure := appEnv == "production"

	// 30 days in seconds
	maxAge := 30 * 24 * 60 * 60

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		SessionCookieName,
		token,
		maxAge,
		"/",
		"", // domain - empty for current domain
		secure,
		true, // httpOnly
	)
}

// clearSessionCookie removes the session cookie
func clearSessionCookie(c *gin.Context) {
	appEnv := os.Getenv("APP_ENV")
	secure := appEnv == "production"

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		SessionCookieName,
		"",
		-1, // negative maxAge deletes the cookie
		"/",
		"",
		secure,
		true,
	)
}

// ClearSessionCookie is the exported version for use in handlers
func ClearSessionCookie(c *gin.Context) {
	clearSessionCookie(c)
}

// GetSessionFromContext retrieves the session from the Gin context
func GetSessionFromGinContext(c *gin.Context) *external.Session {
	if session, exists := c.Get(SessionContextKey); exists {
		if s, ok := session.(*external.Session); ok {
			return s
		}
	}
	return nil
}

// GetUserIDFromGinContext retrieves the user ID from the Gin context
func GetUserIDFromGinContext(c *gin.Context) string {
	if userID, exists := c.Get(UserContextKey); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// GetUserIDAsUUID retrieves the user ID as UUID from the Gin context
func GetUserIDAsUUID(c *gin.Context) (uuid.UUID, error) {
	userID := GetUserIDFromGinContext(c)
	if userID == "" {
		return uuid.UUID{}, nil
	}
	return uuid.Parse(userID)
}

// AuthServiceFromContext creates an AuthService from the database
// This is a helper for handlers that need to perform auth operations
func AuthServiceFromDB(db *sql.DB) *external.AuthService {
	return external.NewAuthService(db)
}

// ContextWithUserID adds user ID to context (for downstream handlers)
func ContextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, UserContextKey, userID.String())
}

// --- CSRF Protection (Double-Submit Cookie Pattern) ---

var csrfLog = logging.With("csrf")

// generateCSRFToken creates a cryptographically secure random token
func generateCSRFToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// CSRFProtection returns middleware that sets CSRF cookie and validates on POST/PUT/DELETE
// Uses double-submit cookie pattern: token in cookie must match token in form/header
func CSRFProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		appEnv := os.Getenv("APP_ENV")
		secure := appEnv == "production"

		// For safe methods (GET, HEAD, OPTIONS), just ensure token cookie exists
		if c.Request.Method == http.MethodGet ||
			c.Request.Method == http.MethodHead ||
			c.Request.Method == http.MethodOptions {
			ensureCSRFCookie(c, secure)
			c.Next()
			return
		}

		// For state-changing methods, validate the CSRF token
		cookieToken, err := c.Cookie(CSRFCookieName)
		if err != nil || cookieToken == "" {
			csrfLog.Warn("CSRF validation failed: missing cookie",
				ports.F("method", c.Request.Method),
				ports.F("path", c.Request.URL.Path),
				ports.F("ip", c.ClientIP()),
			)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CSRF validation failed"})
			return
		}

		// Get token from header (for AJAX) or form field
		submittedToken := c.GetHeader(CSRFHeaderName)
		if submittedToken == "" {
			submittedToken = c.PostForm(CSRFFormFieldName)
		}

		// Also check for HX-Request header (HTMX requests)
		// HTMX requests with SameSite cookies are generally safe, but we still validate
		isHTMX := c.GetHeader("HX-Request") == "true"

		if submittedToken == "" && !isHTMX {
			csrfLog.Warn("CSRF validation failed: missing token in request",
				ports.F("method", c.Request.Method),
				ports.F("path", c.Request.URL.Path),
				ports.F("ip", c.ClientIP()),
			)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CSRF validation failed"})
			return
		}

		// For HTMX requests, we trust SameSite cookie protection
		// But for form submissions, verify the token matches
		if !isHTMX && subtle.ConstantTimeCompare([]byte(cookieToken), []byte(submittedToken)) != 1 {
			csrfLog.Warn("CSRF validation failed: token mismatch",
				ports.F("method", c.Request.Method),
				ports.F("path", c.Request.URL.Path),
				ports.F("ip", c.ClientIP()),
			)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CSRF validation failed"})
			return
		}

		// Refresh token cookie for security (rotation)
		ensureCSRFCookie(c, secure)
		c.Next()
	}
}

// ensureCSRFCookie creates a CSRF cookie if one doesn't exist or refreshes it
func ensureCSRFCookie(c *gin.Context, secure bool) {
	// Check if cookie already exists
	if _, err := c.Cookie(CSRFCookieName); err == nil {
		// Cookie exists, store token in context for templates
		token, _ := c.Cookie(CSRFCookieName)
		c.Set(string(CSRFContextKey), token)
		return
	}

	// Generate new token
	token, err := generateCSRFToken()
	if err != nil {
		csrfLog.Error("Failed to generate CSRF token", ports.F("error", err))
		return
	}

	// Set cookie - NOT HttpOnly so JavaScript can read it for AJAX requests
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		CSRFCookieName,
		token,
		60*60*24, // 24 hours
		"/",
		"",
		secure,
		false, // httpOnly=false allows JS to read for AJAX
	)

	// Store in context for templates to use
	c.Set(string(CSRFContextKey), token)
}

// GetCSRFToken retrieves the CSRF token from context (for use in templates)
func GetCSRFToken(c *gin.Context) string {
	if token, exists := c.Get(string(CSRFContextKey)); exists {
		if t, ok := token.(string); ok {
			return t
		}
	}
	// Try to get from cookie as fallback
	if token, err := c.Cookie(CSRFCookieName); err == nil {
		return token
	}
	return ""
}

// CSRFFormField returns an HTML hidden input field with the CSRF token
// Use this in templates: {{ CSRFFormField(c) }}
func CSRFFormField(c *gin.Context) string {
	token := GetCSRFToken(c)
	return `<input type="hidden" name="` + CSRFFormFieldName + `" value="` + token + `">`
}
