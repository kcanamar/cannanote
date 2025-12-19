package domain

import "errors"

// Domain errors - these represent business rule violations
var (
	// Human domain errors
	ErrInvalidUsername   = errors.New("invalid username")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrHumanNotFound     = errors.New("human not found")
	ErrHumanAlreadyExists = errors.New("human already exists")
	ErrUnauthorized      = errors.New("unauthorized access")
	
	// Entry domain errors  
	ErrInvalidStrain     = errors.New("invalid strain name")
	ErrInvalidAmount     = errors.New("invalid consumption amount")
	ErrEntryNotFound     = errors.New("entry not found")
	ErrEntryNotOwned     = errors.New("entry not owned by human")
	
	// General domain errors
	ErrInvalidInput      = errors.New("invalid input provided")
	ErrDuplicateResource = errors.New("resource already exists")
)

// BusinessError wraps errors with additional context
type BusinessError struct {
	Err     error
	Code    string
	Message string
	Context map[string]interface{}
}

func (e BusinessError) Error() string {
	return e.Err.Error()
}

// NewBusinessError creates a new business error with context
func NewBusinessError(err error, code, message string, context map[string]interface{}) BusinessError {
	return BusinessError{
		Err:     err,
		Code:    code,
		Message: message,
		Context: context,
	}
}