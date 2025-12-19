package ports

import (
	"context"
	"backend/internal/core/domain"
	"github.com/google/uuid"
)

// HumanRepository defines the contract for human data persistence
// This is a port - an interface that domain services depend on
type HumanRepository interface {
	Create(ctx context.Context, human *domain.Human) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Human, error)
	GetByEmail(ctx context.Context, email string) (*domain.Human, error)
	GetByUsername(ctx context.Context, username string) (*domain.Human, error)
	Update(ctx context.Context, human *domain.Human) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*domain.Human, error)
	Search(ctx context.Context, criteria HumanSearchCriteria) ([]*domain.Human, error)
}

// RoleRepository defines the contract for role management
type RoleRepository interface {
	AssignRole(ctx context.Context, humanID uuid.UUID, role domain.Role, assignedBy uuid.UUID) error
	RemoveRole(ctx context.Context, humanID uuid.UUID, role domain.Role) error
	GetRoles(ctx context.Context, humanID uuid.UUID) ([]domain.Role, error)
	HasRole(ctx context.Context, humanID uuid.UUID, role domain.Role) (bool, error)
}

// AuthService defines the contract for authentication operations
type AuthService interface {
	GenerateToken(ctx context.Context, humanID uuid.UUID) (string, error)
	ValidateToken(ctx context.Context, token string) (*AuthClaims, error)
	RevokeToken(ctx context.Context, token string) error
}

// DTOs and criteria for port contracts

type HumanSearchCriteria struct {
	Username string
	Email    string
	Role     domain.Role
}

type AuthClaims struct {
	HumanID uuid.UUID
	Roles   []domain.Role
}