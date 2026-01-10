package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
	"backend/internal/core/domain"
	"backend/internal/core/ports"

	"github.com/google/uuid"
)

// SupabaseHumanRepository implements ports.HumanRepository for Supabase PostgreSQL
// This is an adapter - it translates between domain contracts and Supabase specifics
// Following hexagonal architecture, this adapter implements the port interface
// defined in the core layer, keeping infrastructure concerns separate from business logic
type SupabaseHumanRepository struct {
	db *sql.DB // Direct SQL connection for Supabase PostgreSQL
}

// NewSupabaseHumanRepository creates a new Supabase human repository
// This factory function follows dependency injection patterns by accepting
// the database connection rather than creating it internally
func NewSupabaseHumanRepository(db *sql.DB) ports.HumanRepository {
	return &SupabaseHumanRepository{
		db: db,
	}
}

// Create stores a new human in Supabase PostgreSQL
// Uses JSON columns for complex data types (profile, consent) to leverage PostgreSQL's
// native JSON support while maintaining flexibility for schema evolution
func (r *SupabaseHumanRepository) Create(ctx context.Context, human *domain.Human) error {
	// Serialize profile and consent to JSON for PostgreSQL storage
	// This allows us to store complex nested data while maintaining query capabilities
	profileJSON, err := json.Marshal(human.Profile)
	if err != nil {
		return err
	}

	consentJSON, err := json.Marshal(human.Consent)
	if err != nil {
		return err
	}

	// Supabase RLS (Row Level Security) will handle access control at the database level
	// This query assumes the 'humans' table exists with appropriate RLS policies
	query := `
		INSERT INTO humans (id, username, email, profile, consent, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = r.db.ExecContext(ctx, query,
		human.ID,
		human.Username,
		human.Email,
		profileJSON,
		consentJSON,
		human.CreatedAt,
		human.UpdatedAt,
	)

	return err
}

// GetByID retrieves a human by their UUID
// Returns domain.ErrHumanNotFound if no human exists with the given ID
func (r *SupabaseHumanRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Human, error) {
	query := `
		SELECT id, username, email, profile, consent, created_at, updated_at
		FROM humans 
		WHERE id = $1`

	var human domain.Human
	var profileJSON, consentJSON []byte

	// Query single row with context for cancellation support
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&human.ID,
		&human.Username,
		&human.Email,
		&profileJSON,
		&consentJSON,
		&human.CreatedAt,
		&human.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Convert SQL not found error to domain error
			return nil, domain.ErrHumanNotFound
		}
		return nil, err
	}

	// Deserialize JSON fields back into Go structs
	if err := json.Unmarshal(profileJSON, &human.Profile); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(consentJSON, &human.Consent); err != nil {
		return nil, err
	}

	return &human, nil
}

// GetByEmail retrieves a human by their email address
// Email is expected to be unique in the system for authentication purposes
func (r *SupabaseHumanRepository) GetByEmail(ctx context.Context, email string) (*domain.Human, error) {
	query := `
		SELECT id, username, email, profile, consent, created_at, updated_at
		FROM humans 
		WHERE email = $1`

	var human domain.Human
	var profileJSON, consentJSON []byte

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&human.ID,
		&human.Username,
		&human.Email,
		&profileJSON,
		&consentJSON,
		&human.CreatedAt,
		&human.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrHumanNotFound
		}
		return nil, err
	}

	// Deserialize JSON fields
	if err := json.Unmarshal(profileJSON, &human.Profile); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(consentJSON, &human.Consent); err != nil {
		return nil, err
	}

	return &human, nil
}

// GetByUsername retrieves a human by their username
// Username is also expected to be unique for user-friendly identification
func (r *SupabaseHumanRepository) GetByUsername(ctx context.Context, username string) (*domain.Human, error) {
	query := `
		SELECT id, username, email, profile, consent, created_at, updated_at
		FROM humans 
		WHERE username = $1`

	var human domain.Human
	var profileJSON, consentJSON []byte

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&human.ID,
		&human.Username,
		&human.Email,
		&profileJSON,
		&consentJSON,
		&human.CreatedAt,
		&human.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrHumanNotFound
		}
		return nil, err
	}

	// Deserialize JSON fields
	if err := json.Unmarshal(profileJSON, &human.Profile); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(consentJSON, &human.Consent); err != nil {
		return nil, err
	}

	return &human, nil
}

// Update modifies an existing human's information
// Automatically updates the updated_at timestamp to track changes
func (r *SupabaseHumanRepository) Update(ctx context.Context, human *domain.Human) error {
	// Serialize profile and consent to JSON for storage
	profileJSON, err := json.Marshal(human.Profile)
	if err != nil {
		return err
	}

	consentJSON, err := json.Marshal(human.Consent)
	if err != nil {
		return err
	}

	// Update with current timestamp for audit trail
	human.UpdatedAt = time.Now()

	query := `
		UPDATE humans 
		SET username = $2, email = $3, profile = $4, consent = $5, updated_at = $6
		WHERE id = $1`

	_, err = r.db.ExecContext(ctx, query,
		human.ID,
		human.Username,
		human.Email,
		profileJSON,
		consentJSON,
		human.UpdatedAt,
	)

	return err
}

// Delete removes a human from the system
// In production, consider implementing soft deletes for audit compliance
func (r *SupabaseHumanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM humans WHERE id = $1`
	
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves humans with pagination support
// Orders by creation date descending to show newest humans first
func (r *SupabaseHumanRepository) List(ctx context.Context, limit, offset int) ([]*domain.Human, error) {
	query := `
		SELECT id, username, email, profile, consent, created_at, updated_at
		FROM humans 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure proper cleanup of database resources

	var humans []*domain.Human

	// Iterate through all returned rows
	for rows.Next() {
		var human domain.Human
		var profileJSON, consentJSON []byte

		err := rows.Scan(
			&human.ID,
			&human.Username,
			&human.Email,
			&profileJSON,
			&consentJSON,
			&human.CreatedAt,
			&human.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Deserialize JSON fields for each human
		if err := json.Unmarshal(profileJSON, &human.Profile); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(consentJSON, &human.Consent); err != nil {
			return nil, err
		}

		humans = append(humans, &human)
	}

	// Check for errors that occurred during row iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return humans, nil
}

// Search finds humans by flexible criteria
// Uses ILIKE for case-insensitive partial matching on text fields
// Limited to 100 results to prevent performance issues
func (r *SupabaseHumanRepository) Search(ctx context.Context, criteria ports.HumanSearchCriteria) ([]*domain.Human, error) {
	// Build dynamic query based on provided search criteria
	// Starting with a base query that always evaluates to true
	query := `
		SELECT id, username, email, profile, consent, created_at, updated_at
		FROM humans WHERE 1=1`
	
	args := []interface{}{}
	argCount := 0

	// Add username filter if provided
	if criteria.Username != "" {
		argCount++
		query += " AND username ILIKE $" + string(rune('0'+argCount))
		args = append(args, "%"+criteria.Username+"%") // Partial match with wildcards
	}

	// Add email filter if provided
	if criteria.Email != "" {
		argCount++
		query += " AND email ILIKE $" + string(rune('0'+argCount))
		args = append(args, "%"+criteria.Email+"%")
	}

	// Order by creation date and limit results for performance
	query += " ORDER BY created_at DESC LIMIT 100"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var humans []*domain.Human

	for rows.Next() {
		var human domain.Human
		var profileJSON, consentJSON []byte

		err := rows.Scan(
			&human.ID,
			&human.Username,
			&human.Email,
			&profileJSON,
			&consentJSON,
			&human.CreatedAt,
			&human.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Deserialize JSON fields  
		if err := json.Unmarshal(profileJSON, &human.Profile); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(consentJSON, &human.Consent); err != nil {
			return nil, err
		}

		humans = append(humans, &human)
	}

	return humans, nil
}