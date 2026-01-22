package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// SupertokensCORS returns CORS middleware for Supertokens
// This middleware should be added early in the middleware chain
func SupertokensCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		supertokens.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		})).ServeHTTP(c.Writer, c.Request)
	}
}

// VerifySession returns middleware that requires a valid Supertokens session
// Use this to protect routes that require authentication
func VerifySession(options *sessmodels.VerifySessionOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		session.VerifySession(options, func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		})(c.Writer, c.Request)
	}
}

// OptionalSession returns middleware that extracts session if present but doesn't require it
// Useful for routes that want to show different content for logged-in users
func OptionalSession(options *sessmodels.VerifySessionOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set checkDatabase to false for optional sessions
		if options == nil {
			options = &sessmodels.VerifySessionOptions{}
		}
		sessionCheckResult := false
		options.SessionRequired = &sessionCheckResult

		session.VerifySession(options, func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		})(c.Writer, c.Request)
	}
}

// GetSessionFromContext retrieves the Supertokens session from the request context
// Returns nil if no session exists
func GetSessionFromContext(c *gin.Context) sessmodels.SessionContainer {
	sessionContainer := session.GetSessionFromRequestContext(c.Request.Context())
	return sessionContainer
}

// GetUserIDFromSession extracts the user ID from a Supertokens session
// Returns empty string if session is nil
func GetUserIDFromSession(sessionContainer sessmodels.SessionContainer) string {
	if sessionContainer == nil {
		return ""
	}
	return sessionContainer.GetUserID()
}
