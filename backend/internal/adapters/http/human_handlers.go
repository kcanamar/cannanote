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
// This is an adapter that translates HTTP to domain operations
type HumanHandlers struct {
	humanService *application.HumanService
}

// NewHumanHandlers creates new human HTTP handlers
func NewHumanHandlers(humanService *application.HumanService) *HumanHandlers {
	return &HumanHandlers{
		humanService: humanService,
	}
}

// CreateHuman handles human registration
func (h *HumanHandlers) CreateHuman(c *gin.Context) {
	var req application.CreateHumanRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	human, err := h.humanService.CreateHuman(c.Request.Context(), req)
	if err != nil {
		// Map domain errors to HTTP status codes
		switch err {
		case domain.ErrInvalidUsername, domain.ErrInvalidEmail:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case domain.ErrHumanAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create human"})
		}
		return
	}

	// Return created human with HATEOAS links
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

// GetHuman retrieves a human by ID
func (h *HumanHandlers) GetHuman(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid human ID"})
		return
	}

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

	// Return human with HATEOAS links
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

// UpdateHumanProfile updates human profile information
func (h *HumanHandlers) UpdateHumanProfile(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid human ID"})
		return
	}

	var req application.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

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

	// Return updated human with HATEOAS links
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

// UpdateConsent updates human consent settings
func (h *HumanHandlers) UpdateConsent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid human ID"})
		return
	}

	var req application.UpdateConsentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{
		"message": "Consent updated successfully",
		"consent": human.Consent,
	})
}

// ListHumans retrieves humans with pagination (admin only)
func (h *HumanHandlers) ListHumans(c *gin.Context) {
	// Parse pagination parameters
	limitParam := c.DefaultQuery("limit", "20")
	offsetParam := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	// TODO: Add authorization check for admin role
	// For now, this is a placeholder for the admin functionality

	c.JSON(http.StatusOK, gin.H{
		"message": "List humans endpoint - requires admin authorization",
		"limit":   limit,
		"offset":  offset,
	})
}

// Response DTOs with HATEOAS links
type HumanResponse struct {
	Human *domain.Human `json:"human"`
	Links HATEOASLinks  `json:"_links"`
}

type HATEOASLinks struct {
	Self    Link `json:"self"`
	Update  Link `json:"update,omitempty"`
	Profile Link `json:"profile,omitempty"`
	Delete  Link `json:"delete,omitempty"`
}

type Link struct {
	Href   string `json:"href"`
	Method string `json:"method"`
	Type   string `json:"type"`
}