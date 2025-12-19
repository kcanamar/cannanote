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
type SupabaseHumanRepository struct {
	db *sql.DB
}

// NewSupabaseHumanRepository creates a new Supabase human repository
func NewSupabaseHumanRepository(db *sql.DB) ports.HumanRepository {
	return &SupabaseHumanRepository{
		db: db,
	}
}

// Create stores a new human in Supabase
func (r *SupabaseHumanRepository) Create(ctx context.Context, human *domain.Human) error {
	// Serialize profile and consent to JSON for PostgreSQL
	profileJSON, err := json.Marshal(human.Profile)
	if err != nil {
		return err
	}

	consentJSON, err := json.Marshal(human.Consent)
	if err != nil {
		return err
	}

	// Supabase RLS will handle row-level security
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
func (r *SupabaseHumanRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Human, error) {
	query := `
		SELECT id, username, email, profile, consent, created_at, updated_at
		FROM humans 
		WHERE id = $1`

	var human domain.Human
	var profileJSON, consentJSON []byte

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

// GetByEmail retrieves a human by their email address
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
func (r *SupabaseHumanRepository) Update(ctx context.Context, human *domain.Human) error {
	// Serialize profile and consent to JSON
	profileJSON, err := json.Marshal(human.Profile)
	if err != nil {
		return err
	}

	consentJSON, err := json.Marshal(human.Consent)
	if err != nil {
		return err
	}

	// Update with current timestamp
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
func (r *SupabaseHumanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM humans WHERE id = $1`
	
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves humans with pagination
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

// Search finds humans by criteria
func (r *SupabaseHumanRepository) Search(ctx context.Context, criteria ports.HumanSearchCriteria) ([]*domain.Human, error) {
	// Build dynamic query based on criteria
	query := `
		SELECT id, username, email, profile, consent, created_at, updated_at
		FROM humans WHERE 1=1`
	
	args := []interface{}{}
	argCount := 0

	if criteria.Username != "" {
		argCount++
		query += " AND username ILIKE $" + string(rune(argCount))
		args = append(args, "%"+criteria.Username+"%")
	}

	if criteria.Email != "" {
		argCount++
		query += " AND email ILIKE $" + string(rune(argCount))
		args = append(args, "%"+criteria.Email+"%")
	}

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