package server

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"backend/cmd/web"
	httpHandlers "backend/internal/adapters/http"
	"github.com/a-h/templ"
)

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
	email := c.PostForm("email")
	consent := c.PostForm("consent")
	
	// Basic validation
	if email == "" {
		templ.Handler(web.BetaSignupError("Email is required")).ServeHTTP(c.Writer, c.Request)
		return
	}
	
	if consent != "on" {
		templ.Handler(web.BetaSignupError("Privacy consent is required to join the beta")).ServeHTTP(c.Writer, c.Request)
		return
	}
	
	// TODO: Store in database with proper validation
	// For now, just return success response
	log.Printf("INFO: Beta signup request: %s (consent: %s)", email, consent)
	
	templ.Handler(web.BetaSignupSuccess(email)).ServeHTTP(c.Writer, c.Request)
}
