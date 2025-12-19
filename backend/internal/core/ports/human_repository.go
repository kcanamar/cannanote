package ports

import (
	"context"
	"backend/internal/core/domain"

	"github.com/google/uuid"
)

// HumanRepository defines the contract for human data persistence
// This is a port interface - no implementation details, pure contract
type HumanRepository interface {
	// Create stores a new human in the system
	Create(ctx context.Context, human *domain.Human) error
	
	// GetByID retrieves a human by their unique identifier
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Human, error)
	
	// GetByEmail retrieves a human by their email address
	GetByEmail(ctx context.Context, email string) (*domain.Human, error)
	
	// GetByUsername retrieves a human by their username  
	GetByUsername(ctx context.Context, username string) (*domain.Human, error)
	
	// Update modifies an existing human's information
	Update(ctx context.Context, human *domain.Human) error
	
	// Delete removes a human from the system
	Delete(ctx context.Context, id uuid.UUID) error
	
	// List retrieves humans with pagination
	List(ctx context.Context, limit, offset int) ([]*domain.Human, error)
	
	// Search finds humans by criteria (for admin functionality)
	Search(ctx context.Context, criteria HumanSearchCriteria) ([]*domain.Human, error)
}

// RoleRepository manages human role assignments
type RoleRepository interface {
	// AssignRole grants a role to a human
	AssignRole(ctx context.Context, humanID uuid.UUID, role domain.Role, grantedBy uuid.UUID) error
	
	// RevokeRole removes a role from a human
	RevokeRole(ctx context.Context, humanID uuid.UUID, role domain.Role) error
	
	// GetRoles retrieves all roles for a human
	GetRoles(ctx context.Context, humanID uuid.UUID) ([]*domain.HumanRole, error)
	
	// HasRole checks if human has specific role
	HasRole(ctx context.Context, humanID uuid.UUID, role domain.Role) (bool, error)
}

// AuthService defines external authentication operations
type AuthService interface {
	// ValidateToken verifies an authentication token
	ValidateToken(ctx context.Context, token string) (*AuthClaims, error)
	
	// RefreshToken generates a new access token
	RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
	
	// RevokeToken invalidates a token
	RevokeToken(ctx context.Context, token string) error
}

// HumanSearchCriteria defines search parameters
type HumanSearchCriteria struct {
	Username    string
	Email       string  
	Role        domain.Role
	CreatedFrom *string
	CreatedTo   *string
}

// AuthClaims represents decoded authentication token information
type AuthClaims struct {
	HumanID  uuid.UUID `json:"human_id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Roles    []domain.Role `json:"roles"`
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}