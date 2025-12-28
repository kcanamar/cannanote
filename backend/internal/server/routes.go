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

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", httpHandlers.HealthHandler)

	r.Static("/assets", "./cmd/web/assets")

	r.GET("/web", func(c *gin.Context) {
		templ.Handler(web.HelloForm()).ServeHTTP(c.Writer, c.Request)
	})

	r.POST("/hello", func(c *gin.Context) {
		web.HelloWebHandler(c.Writer, c.Request)
	})

	// Educational content routes
	r.GET("/learn/cannabinoids", func(c *gin.Context) {
		web.CannabinoidsHandler(c)
	})

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
