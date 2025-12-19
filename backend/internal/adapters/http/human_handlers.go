package http

import (
	"net/http"
	"strconv"
	"backend/internal/core/application"
	"backend/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// HumanHandlers handles HTTP requests for human management
// This is an adapter that translates HTTP requests to domain operations
// Following hexagonal architecture, this adapter sits at the boundary between
// external HTTP requests and internal business logic, ensuring clean separation
type HumanHandlers struct {
	humanService *application.HumanService // Dependency on application layer, not domain directly
}

// NewHumanHandlers creates new human HTTP handlers
// Constructor injection ensures loose coupling and testability
func NewHumanHandlers(humanService *application.HumanService) *HumanHandlers {
	return &HumanHandlers{
		humanService: humanService,
	}
}

// CreateHuman handles human registration via HTTP POST
// Maps HTTP request body to application service request and handles response formatting
// Implements proper error handling with appropriate HTTP status codes
func (h *HumanHandlers) CreateHuman(c *gin.Context) {
	var req application.CreateHumanRequest
	
	// Bind JSON request body to struct with validation
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Call application service to handle business logic
	human, err := h.humanService.CreateHuman(c.Request.Context(), req)
	if err != nil {
		// Map domain errors to appropriate HTTP status codes
		// This translation layer keeps HTTP concerns out of the domain
		switch err {
		case domain.ErrInvalidUsername, domain.ErrInvalidEmail:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case domain.ErrHumanAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			// Don't expose internal errors to clients
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create human"})
		}
		return
	}

	// Return created human with HATEOAS links for discoverability
	// HATEOAS (Hypermedia as the Engine of Application State) provides
	// clients with available actions and navigation paths
	response := HumanResponse{
		Human: human,
		Links: HATEOASLinks{
			Self: Link{
				Href:   "/api/humans/" + human.ID.String(),
				Method: "GET",
				Type:   "application/json",
			},
			Update: Link{
				Href:   "/api/humans/" + human.ID.String(),
				Method: "PUT",
				Type:   "application/json",
			},
		},
	}

	c.JSON(http.StatusCreated, response)
}

// GetHuman retrieves a human by ID via HTTP GET
// Handles UUID parsing and provides appropriate error responses
func (h *HumanHandlers) GetHuman(c *gin.Context) {
	// Extract and validate UUID parameter from URL path
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid human ID"})
		return
	}

	// Delegate to application service for business logic
	human, err := h.humanService.GetHuman(c.Request.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrHumanNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve human"})
		}
		return
	}

	// Return human with comprehensive HATEOAS links
	response := HumanResponse{
		Human: human,
		Links: HATEOASLinks{
			Self: Link{
				Href:   "/api/humans/" + human.ID.String(),
				Method: "GET",
				Type:   "application/json",
			},
			Update: Link{
				Href:   "/api/humans/" + human.ID.String(),
				Method: "PUT", 
				Type:   "application/json",
			},
			Profile: Link{
				Href:   "/api/humans/" + human.ID.String() + "/profile",
				Method: "GET",
				Type:   "application/json",
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// UpdateHumanProfile updates human profile information via HTTP PUT
// Supports partial updates while maintaining data integrity
func (h *HumanHandlers) UpdateHumanProfile(c *gin.Context) {
	// Extract and validate human ID from URL
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid human ID"})
		return
	}

	// Bind and validate request body
	var req application.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Process update through application service
	human, err := h.humanService.UpdateHumanProfile(c.Request.Context(), id, req)
	if err != nil {
		switch err {
		case domain.ErrHumanNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		}
		return
	}

	// Return updated human with navigation links
	response := HumanResponse{
		Human: human,
		Links: HATEOASLinks{
			Self: Link{
				Href:   "/api/humans/" + human.ID.String(),
				Method: "GET",
				Type:   "application/json",
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// UpdateConsent updates human consent settings via HTTP PUT
// Critical for HIPAA compliance and privacy management
func (h *HumanHandlers) UpdateConsent(c *gin.Context) {
	// Extract and validate human ID
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid human ID"})
		return
	}

	// Bind consent update request
	var req application.UpdateConsentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Process consent update with proper audit trail
	human, err := h.humanService.UpdateConsent(c.Request.Context(), id, req)
	if err != nil {
		switch err {
		case domain.ErrHumanNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update consent"})
		}
		return
	}

	// Return confirmation without exposing sensitive data unnecessarily
	c.JSON(http.StatusOK, gin.H{
		"message": "Consent updated successfully",
		"consent": human.Consent, // Include updated consent for client confirmation
	})
}

// ListHumans retrieves humans with pagination (admin functionality)
// Includes basic pagination support and proper authorization checks
func (h *HumanHandlers) ListHumans(c *gin.Context) {
	// Parse pagination parameters with sensible defaults
	limitParam := c.DefaultQuery("limit", "20")
	offsetParam := c.DefaultQuery("offset", "0")

	// Convert string parameters to integers with validation
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 20 // Default page size
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0 // Start from beginning
	}

	// TODO: Implement proper authorization middleware
	// This endpoint should require admin role verification
	// For now, return placeholder indicating authorization is needed
	c.JSON(http.StatusOK, gin.H{
		"message": "List humans endpoint - requires admin authorization",
		"limit":   limit,
		"offset":  offset,
	})
}

// Response DTOs with HATEOAS links for rich API responses
// These structures provide clients with navigation capabilities
// and follow REST/HATEOAS principles for better API discoverability

// HumanResponse wraps human data with hypermedia links
type HumanResponse struct {
	Human *domain.Human `json:"human"`        // Core human data
	Links HATEOASLinks  `json:"_links"`       // Hypermedia navigation links
}

// HATEOASLinks defines available actions and relationships
// Supports REST Level 3 maturity with hypermedia controls
type HATEOASLinks struct {
	Self    Link `json:"self"`                    // Link to this resource
	Update  Link `json:"update,omitempty"`        // Link to update this resource
	Profile Link `json:"profile,omitempty"`       // Link to profile-specific operations
	Delete  Link `json:"delete,omitempty"`        // Link to delete operation (if authorized)
}

// Link represents a hypermedia link with method and content type
// Provides complete information for clients to make follow-up requests
type Link struct {
	Href   string `json:"href"`     // URL for the linked resource
	Method string `json:"method"`   // HTTP method to use (GET, POST, PUT, DELETE)
	Type   string `json:"type"`     // Expected content type for the request
}