package http

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"backend/internal/adapters/external"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// SessionCookieName is the name of the session cookie
	SessionCookieName = "cannanote_session"
	// SessionContextKey is the key used to store session in context
	SessionContextKey = "session"
	// UserContextKey is the key used to store user in context
	UserContextKey = "user"
)

// AuthMiddleware provides authentication middleware using our custom auth service
type AuthMiddleware struct {
	authService *external.AuthService
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(db *sql.DB) *AuthMiddleware {
	return &AuthMiddleware{
		authService: external.NewAuthService(db),
	}
}

// RequireSession returns middleware that requires a valid session
// Redirects to /login if no valid session
func (m *AuthMiddleware) RequireSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getSessionToken(c)
		if token == "" {
			log.Printf("[AUTH] No session token found, redirecting to login")
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		session, err := m.authService.ValidateSession(c.Request.Context(), token)
		if err != nil {
			log.Printf("[AUTH] Invalid session: %v", err)
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
