package domain

import (
	"time"
	"github.com/google/uuid"
)

// Human represents a human in the CannaNote system
// This is a pure domain entity with no I/O dependencies
type Human struct {
	ID        uuid.UUID       `json:"id"`
	Username  string          `json:"username"`
	Email     string          `json:"email"`
	Profile   HumanProfile    `json:"profile"`
	Consent   ConsentSettings `json:"consent"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// HumanProfile contains cannabis preferences and profile information
type HumanProfile struct {
	PreferredStrains     []string `json:"preferred_strains"`
	PreferredConsumption []string `json:"preferred_consumption"`
	ExperienceLevel      string   `json:"experience_level"`
	PublicProfile        bool     `json:"public_profile"`
	ShareExperiences     bool     `json:"share_experiences"`
	PreferredDispensary  string   `json:"preferred_dispensary"`
}

// ConsentSettings tracks HIPAA-ready consent management
type ConsentSettings struct {
	DataCollection     bool      `json:"data_collection"`
	MarketingEmails    bool      `json:"marketing_emails"`
	DataSharing        bool      `json:"data_sharing"`
	MedicalDataSharing bool      `json:"medical_data_sharing"`
	ConsentDate        time.Time `json:"consent_date"`
	ConsentVersion     string    `json:"consent_version"`
}

// Role represents system roles for humans
type Role string

const (
	RoleHuman     Role = "human"
	RoleModerator Role = "moderator"
	RoleAdmin     Role = "admin"
)

// Domain methods - pure business logic with no I/O

// Validate performs domain validation on the human entity
func (h *Human) Validate() error {
	if h.Username == "" {
		return ErrInvalidUsername
	}
	if h.Email == "" {
		return ErrInvalidEmail
	}
	return nil
}

// UpdateProfile updates the human's profile through domain rules
func (h *Human) UpdateProfile(profile HumanProfile) {
	h.Profile = profile
	h.UpdatedAt = time.Now()
}

// GrantConsent updates consent settings with timestamp
func (h *Human) GrantConsent(consent ConsentSettings) {
	consent.ConsentDate = time.Now()
	h.Consent = consent
	h.UpdatedAt = time.Now()
}

// CanAccessFeature checks if human can access a feature based on consent
func (h *Human) CanAccessFeature(feature string) bool {
	switch feature {
	case "data_analytics":
		return h.Consent.DataCollection
	case "social_features":
		return h.Profile.ShareExperiences
	case "public_profile":
		return h.Profile.PublicProfile
	default:
		return true
	}
}

// IsExperienced returns true if human has significant cannabis experience
func (h *Human) IsExperienced() bool {
	return h.Profile.ExperienceLevel == "experienced" || h.Profile.ExperienceLevel == "expert"
}