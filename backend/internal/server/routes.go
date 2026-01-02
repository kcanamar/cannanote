package server

import (
	"log"
	"net/http"

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

	r.Static("/assets", "./cmd/web/assets")

	r.GET("/web", func(c *gin.Context) {
		templ.Handler(web.HelloForm()).ServeHTTP(c.Writer, c.Request)
	})

	r.POST("/hello", func(c *gin.Context) {
		web.HelloWebHandler(c.Writer, c.Request)
	})

	// Documentation routes
	r.GET("/docs", func(c *gin.Context) {
		templ.Handler(web.Docs()).ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/docs/guides/cannabinoids", func(c *gin.Context) {
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
