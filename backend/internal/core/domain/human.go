package domain

import (
	"time"

	"github.com/google/uuid"
)

// Human represents a person using the CannaNote platform
// This is the core entity for human management domain
type Human struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Profile information
	Profile     HumanProfile `json:"profile"`
	
	// Privacy and consent settings
	Consent     ConsentSettings `json:"consent"`
}

// HumanProfile contains non-PHI profile information safe to collect in Phase 1
type HumanProfile struct {
	// Cannabis preferences (recreational tracking)
	PreferredStrains      []string `json:"preferred_strains"`
	PreferredConsumption  []string `json:"preferred_consumption"` // flower, vape, edible, etc.
	ExperienceLevel      string   `json:"experience_level"`      // beginner, intermediate, experienced
	
	// Social settings
	PublicProfile        bool     `json:"public_profile"`
	ShareExperiences     bool     `json:"share_experiences"`
	
	// Dispensary affiliations (for B2B partnerships)
	PreferredDispensary  string   `json:"preferred_dispensary,omitempty"`
}

// ConsentSettings manages privacy and data sharing permissions
type ConsentSettings struct {
	DataCollection      bool      `json:"data_collection"`       // Consent to data collection
	MarketingEmails     bool      `json:"marketing_emails"`      // Marketing communications
	DataSharing         bool      `json:"data_sharing"`          // Anonymous analytics sharing
	ConsentDate         time.Time `json:"consent_date"`
	ConsentVersion      string    `json:"consent_version"`       // Track consent version changes
	
	// Future HIPAA preparation
	MedicalDataSharing  bool      `json:"medical_data_sharing"`  // For Phase 3
}

// Role represents user authorization levels
type Role string

const (
	RoleHuman            Role = "human"              // Standard platform user
	RoleAdmin            Role = "admin"              // Platform administrator  
	RoleDispensaryPartner Role = "dispensary_partner" // B2B dispensary partner
	RoleProvider         Role = "provider"           // Healthcare provider (Phase 3)
)

// HumanRole links humans to their platform roles
type HumanRole struct {
	HumanID     uuid.UUID `json:"human_id"`
	Role        Role      `json:"role"`
	GrantedAt   time.Time `json:"granted_at"`
	GrantedBy   uuid.UUID `json:"granted_by"`
	Active      bool      `json:"active"`
}

// Validate ensures the human entity is valid
func (h *Human) Validate() error {
	if h.Username == "" {
		return ErrInvalidUsername
	}
	if h.Email == "" {
		return ErrInvalidEmail
	}
	return nil
}

// HasRole checks if human has a specific role
func (h *Human) HasRole(role Role) bool {
	// This would typically query the roles through a port interface
	// For now, returning basic logic
	return true // Implementation will use repository
}

// CanAccessFeature checks if human can access premium features
func (h *Human) CanAccessFeature(feature string) bool {
	// Freemium logic - expand as needed
	switch feature {
	case "unlimited_entries":
		return h.HasRole(RoleAdmin) // Premium feature
	case "analytics_export":
		return h.HasRole(RoleAdmin) // Premium feature
	case "basic_tracking":
		return true // Free feature
	default:
		return false
	}
}

// UpdateProfile updates human profile information
func (h *Human) UpdateProfile(profile HumanProfile) {
	h.Profile = profile
	h.UpdatedAt = time.Now()
}

// GrantConsent updates consent settings
func (h *Human) GrantConsent(consent ConsentSettings) {
	consent.ConsentDate = time.Now()
	h.Consent = consent
	h.UpdatedAt = time.Now()
}