package domain

import (
	"time"

	"github.com/google/uuid"
)

// Session represents a cannabis consumption session
// This is the core domain entity for session logging
type Session struct {
	ID              uuid.UUID      `json:"id"`
	HumanID         uuid.UUID      `json:"human_id"`
	ProductID       *uuid.UUID     `json:"product_id,omitempty"`
	Method          ConsumptionMethod `json:"method"`
	Amount          float64        `json:"amount"`
	Unit            AmountUnit     `json:"unit"`
	MoodBefore      *int           `json:"mood_before,omitempty"`
	MindBefore      *int           `json:"mind_before,omitempty"`
	BodyBefore      *int           `json:"body_before,omitempty"`
	CheckInMode     CheckInMode    `json:"check_in_mode"`
	CheckInInterval int            `json:"check_in_interval"`
	Status          SessionStatus  `json:"status"`
	SyncStatus      SyncStatus     `json:"sync_status"`
	Notes           string         `json:"notes,omitempty"`
	Timestamp       time.Time      `json:"timestamp"`
	CompletedAt     *time.Time     `json:"completed_at,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// ConsumptionMethod represents how cannabis is consumed
type ConsumptionMethod string

const (
	MethodVape    ConsumptionMethod = "vape"
	MethodSmoke   ConsumptionMethod = "smoke"
	MethodEdible  ConsumptionMethod = "edible"
	MethodTincture ConsumptionMethod = "tincture"
	MethodTopical ConsumptionMethod = "topical"
	MethodDab     ConsumptionMethod = "dab"
)

// AmountUnit represents the unit of measurement for consumption
type AmountUnit string

const (
	UnitPuffs        AmountUnit = "puffs"
	UnitSeconds      AmountUnit = "seconds"
	UnitHits         AmountUnit = "hits"
	UnitGrams        AmountUnit = "g"
	UnitJoints       AmountUnit = "joints"
	UnitMg           AmountUnit = "mg"
	UnitMl           AmountUnit = "ml"
	UnitDrops        AmountUnit = "drops"
	UnitPieces       AmountUnit = "pieces"
	UnitApplications AmountUnit = "applications"
)

// DefaultUnitsForMethod returns the default units for a consumption method
func DefaultUnitsForMethod(method ConsumptionMethod) []AmountUnit {
	switch method {
	case MethodVape:
		return []AmountUnit{UnitPuffs, UnitSeconds}
	case MethodSmoke:
		return []AmountUnit{UnitHits, UnitGrams, UnitJoints}
	case MethodEdible:
		return []AmountUnit{UnitMg, UnitPieces}
	case MethodTincture:
		return []AmountUnit{UnitMg, UnitMl, UnitDrops}
	case MethodTopical:
		return []AmountUnit{UnitApplications}
	case MethodDab:
		return []AmountUnit{UnitHits, UnitGrams}
	default:
		return []AmountUnit{UnitHits}
	}
}

// CheckInMode determines how check-ins are handled
type CheckInMode string

const (
	CheckInModeNone     CheckInMode = "none"
	CheckInModeSingle   CheckInMode = "single"
	CheckInModeMultiple CheckInMode = "multiple"
)

// SessionStatus represents the current state of a session
type SessionStatus string

const (
	SessionStatusActive    SessionStatus = "active"
	SessionStatusCompleted SessionStatus = "completed"
)

// SyncStatus represents sync state for offline-first
type SyncStatus string

const (
	SyncStatusLocal   SyncStatus = "local"
	SyncStatusPending SyncStatus = "pending"
	SyncStatusSynced  SyncStatus = "synced"
)

// CheckIn represents an effect check-in during or after a session
type CheckIn struct {
	ID               uuid.UUID  `json:"id"`
	SessionID        uuid.UUID  `json:"session_id"`
	Mood             int        `json:"mood"`
	Mind             int        `json:"mind"`
	Body             int        `json:"body"`
	Intensity        int        `json:"intensity"`
	MeetsExpectations *string   `json:"meets_expectations,omitempty"`
	Notes            string     `json:"notes,omitempty"`
	Timestamp        time.Time  `json:"timestamp"`
}

// FeelingLevel represents a 1-5 scale for mood/mind/body grids
type FeelingLevel int

const (
	FeelingVeryLow  FeelingLevel = 1
	FeelingLow      FeelingLevel = 2
	FeelingNeutral  FeelingLevel = 3
	FeelingGood     FeelingLevel = 4
	FeelingExcellent FeelingLevel = 5
)

// Validate performs domain validation on the session
func (s *Session) Validate() error {
	if s.Method == "" {
		return ErrInvalidSession
	}
	if s.Amount < 0 {
		return ErrInvalidSession
	}
	return nil
}

// Complete marks the session as completed
func (s *Session) Complete() {
	now := time.Now()
	s.Status = SessionStatusCompleted
	s.CompletedAt = &now
	s.UpdatedAt = now
}

// IsActive returns true if the session is still active
func (s *Session) IsActive() bool {
	return s.Status == SessionStatusActive
}

// CanAddCheckIn returns true if check-ins are allowed
func (s *Session) CanAddCheckIn() bool {
	return s.CheckInMode != CheckInModeNone && s.IsActive()
}

// DurationMinutes returns the session duration in minutes
func (s *Session) DurationMinutes() int {
	if s.CompletedAt != nil {
		return int(s.CompletedAt.Sub(s.Timestamp).Minutes())
	}
	return int(time.Since(s.Timestamp).Minutes())
}

// IsWithinCheckInWindow returns true if session is within 4-hour check-in window
func (s *Session) IsWithinCheckInWindow() bool {
	return s.DurationMinutes() <= 240
}

// Validate performs domain validation on the check-in
func (c *CheckIn) Validate() error {
	if c.SessionID == uuid.Nil {
		return ErrInvalidCheckIn
	}
	if c.Mood < 1 || c.Mood > 5 {
		return ErrInvalidCheckIn
	}
	if c.Mind < 1 || c.Mind > 5 {
		return ErrInvalidCheckIn
	}
	if c.Body < 1 || c.Body > 5 {
		return ErrInvalidCheckIn
	}
	if c.Intensity < 1 || c.Intensity > 10 {
		return ErrInvalidCheckIn
	}
	return nil
}
