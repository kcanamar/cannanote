package domain

import "errors"

// Domain errors - these represent business rule violations
var (
	// Human validation errors
	ErrInvalidUsername     = errors.New("username cannot be empty")
	ErrInvalidEmail        = errors.New("email cannot be empty")
	ErrHumanAlreadyExists  = errors.New("human already exists")
	ErrHumanNotFound       = errors.New("human not found")

	// Authorization errors
	ErrUnauthorized        = errors.New("unauthorized access")
	ErrInsufficientPermissions = errors.New("insufficient permissions")

	// Consent errors
	ErrConsentRequired     = errors.New("consent required for this operation")
	ErrConsentExpired      = errors.New("consent has expired")

	// Role errors
	ErrRoleNotFound        = errors.New("role not found")
	ErrInvalidRole         = errors.New("invalid role")

	// Session errors
	ErrInvalidSession      = errors.New("invalid session data")
	ErrSessionNotFound     = errors.New("session not found")
	ErrSessionCompleted    = errors.New("session already completed")

	// Check-in errors
	ErrInvalidCheckIn      = errors.New("invalid check-in data")
	ErrCheckInNotAllowed   = errors.New("check-ins not allowed for this session")

	// Product errors
	ErrInvalidProduct      = errors.New("invalid product data")
	ErrProductNotFound     = errors.New("product not found")
)