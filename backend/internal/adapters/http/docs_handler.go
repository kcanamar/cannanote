package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/docs"
)

// DocsHandler manages documentation requests
type DocsHandler struct {
	cache      *docs.DocsCache
	contentDir string
}

// NewDocsHandler creates a new docs handler
func NewDocsHandler(contentDir string) (*DocsHandler, error) {
	cache := docs.NewDocsCache()
	
	// Load all content at startup
	if err := cache.LoadContent(contentDir); err != nil {
		return nil, err
	}

	return &DocsHandler{
		cache:      cache,
		contentDir: contentDir,
	}, nil
}

// HandleDocsRequest handles documentation page requests
func (h *DocsHandler) HandleDocsRequest(c *gin.Context) {
	// For now, redirect to the fallback docs page until content loading is working
	c.Redirect(http.StatusTemporaryRedirect, "/docs-fallback")
}

// HandleDocsSearch handles search requests
func (h *DocsHandler) HandleDocsSearch(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No search query provided"})
		return
	}

	// TODO: Implement search functionality
	// For now, return empty results
	c.JSON(http.StatusOK, gin.H{
		"results": []interface{}{},
		"query":   query,
	})
}