package http

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/cmd/web"
	"backend/internal/adapters/external"
	"backend/internal/docs"
)

// DocsHandler manages documentation requests
type DocsHandler struct {
	cache       *docs.DocsCache
	contentDir  string
	authService *external.AuthService
}

// NewDocsHandler creates a new docs handler
func NewDocsHandler(contentDir string, db *sql.DB) (*DocsHandler, error) {
	cache := docs.NewDocsCache()

	// Load all content at startup
	if err := cache.LoadContent(contentDir); err != nil {
		return nil, err
	}

	return &DocsHandler{
		cache:       cache,
		contentDir:  contentDir,
		authService: external.NewAuthService(db),
	}, nil
}

// HandleDocsRequest handles documentation page requests
func (h *DocsHandler) HandleDocsRequest(c *gin.Context) {
	// Get the requested path, removing the /docs prefix
	requestPath := strings.TrimPrefix(c.Request.URL.Path, "/docs")

	// Clean the path - remove leading/trailing slashes
	requestPath = strings.Trim(requestPath, "/")

	// If path is empty, serve the index page
	if requestPath == "" {
		requestPath = "" // This will look for index.md in the root
	}

	// Debug: log what we're looking for
	c.Header("X-Debug-Request-Path", requestPath)

	// Try to get the document - first try exact path
	doc, exists := h.cache.GetDoc(requestPath)

	// If not found, try with /index appended (for directory requests like /getting-started -> /getting-started/index)
	if !exists && requestPath != "" {
		indexPath := requestPath + "/index"
		doc, exists = h.cache.GetDoc(indexPath)
		if exists {
			c.Header("X-Debug-Found-Path", indexPath)
			requestPath = indexPath
		}
	}

	// If still not found, return debug info
	if !exists {
		// Get sidebar to see what docs are loaded
		sidebar := h.cache.GetSidebar()
		c.JSON(http.StatusNotFound, gin.H{
			"error":          "Documentation page not found",
			"requested_path": requestPath,
			"available_docs": len(sidebar),
			"debug":          "Check server logs for content loading errors",
		})
		return
	}

	// Get user email from context if authenticated (set by OptionalSession middleware)
	userEmail := ""
	if userIDStr := GetUserIDFromGinContext(c); userIDStr != "" {
		if userID, err := uuid.Parse(userIDStr); err == nil {
			if email, err := h.authService.GetUserEmailByID(c.Request.Context(), userID); err == nil {
				userEmail = email
			}
		}
	}

	// Check if this is an HTMX request
	isHTMX := c.GetHeader("HX-Request") == "true"

	if isHTMX {
		// Render partial content only for HTMX swaps
		component := web.DocsContentPartial(doc)
		templ.Handler(component).ServeHTTP(c.Writer, c.Request)
	} else {
		// Render full layout for direct visits (SSR)
		sidebar := h.cache.GetSidebar()
		currentPath := "/docs/" + requestPath
		component := web.DocsContentWithAuth(doc, sidebar, currentPath, userEmail)
		templ.Handler(component).ServeHTTP(c.Writer, c.Request)
	}
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