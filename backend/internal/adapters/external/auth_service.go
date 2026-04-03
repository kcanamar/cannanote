package external

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"backend/internal/adapters/logging"
	"backend/internal/core/ports"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	// Session duration: 30 days
	SessionDuration = 30 * 24 * time.Hour
	// Email verification token duration: 24 hours
	EmailVerificationDuration = 24 * time.Hour
	// Password reset token duration: 1 hour
	PasswordResetDuration = 1 * time.Hour
	// bcrypt cost factor
	BcryptCost = 12
)

// AuthService handles authentication operations
type AuthService struct {
	db  *sql.DB
	log ports.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{
		db:  db,
		log: logging.With("auth"),
	}
}

// UserProfile represents a user profile from the database
type UserProfile struct {
	ID                uuid.UUID
	Email             string
	PasswordHash      sql.NullString
	EmailVerified     bool
	BetaStatus        string
	SubscriptionTier  string
	ActivatedAt       sql.NullTime
	CreatedAt         time.Time
}

// Session represents an active session
type Session struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	ExpiresAt      time.Time
	LastAccessedAt time.Time
}

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// CheckPassword verifies a password against its hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateSecureToken creates a cryptographically secure random token
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// HashToken creates a SHA-256 hash of a token for storage
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// CreateUser creates a new user with email (password set later after verification)
func (a *AuthService) CreateUser(ctx context.Context, email string) (*UserProfile, error) {
	// Check if user already exists
	var existingID uuid.UUID
	err := a.db.QueryRowContext(ctx, "SELECT id FROM profiles WHERE email = $1", email).Scan(&existingID)
	if err == nil {
		return nil, fmt.Errorf("email %s already registered", email)
	}
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Generate email verification token
	verificationToken, err := GenerateSecureToken(32)
	if err != nil {
		return nil, err
	}

	// Create user
	userID := uuid.New()
	expiresAt := time.Now().Add(EmailVerificationDuration)

	_, err = a.db.ExecContext(ctx, `
		INSERT INTO profiles (id, email, beta_status, beta_joined_at, subscription_tier,
		                      email_verified, email_verification_token, email_verification_expires_at, created_at)
		VALUES ($1, $2, 'waitlist', NOW(), 'beta_grandfathered',
		        FALSE, $3, $4, NOW())
	`, userID, email, HashToken(verificationToken), expiresAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	a.log.Info("Created new user", ports.F("email", email), ports.F("user_id", userID))

	return &UserProfile{
		ID:            userID,
		Email:         email,
		EmailVerified: false,
		BetaStatus:    "waitlist",
	}, nil
}

// GetEmailVerificationToken returns the raw token for sending in email
// This should only be called right after CreateUser
func (a *AuthService) CreateEmailVerificationToken(ctx context.Context, userID uuid.UUID) (string, error) {
	token, err := GenerateSecureToken(32)
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(EmailVerificationDuration)

	_, err = a.db.ExecContext(ctx, `
		UPDATE profiles
		SET email_verification_token = $1, email_verification_expires_at = $2
		WHERE id = $3
	`, HashToken(token), expiresAt, userID)

	if err != nil {
		return "", fmt.Errorf("failed to set verification token: %w", err)
	}

	return token, nil
}

// VerifyEmail verifies a user's email using the token
func (a *AuthService) VerifyEmail(ctx context.Context, token string) (*UserProfile, error) {
	tokenHash := HashToken(token)

	var user UserProfile
	err := a.db.QueryRowContext(ctx, `
		SELECT id, email, email_verified, beta_status, email_verification_expires_at
		FROM profiles
		WHERE email_verification_token = $1
	`, tokenHash).Scan(&user.ID, &user.Email, &user.EmailVerified, &user.BetaStatus, &user.ActivatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invalid or expired verification token")
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check if already verified
	if user.EmailVerified {
		return &user, nil
	}

	// Check expiration
	if user.ActivatedAt.Valid && time.Now().After(user.ActivatedAt.Time) {
		return nil, fmt.Errorf("verification token has expired")
	}

	// Mark email as verified and activate
	_, err = a.db.ExecContext(ctx, `
		UPDATE profiles
		SET email_verified = TRUE,
		    email_verification_token = NULL,
		    email_verification_expires_at = NULL,
		    beta_status = 'active',
		    activated_at = NOW()
		WHERE id = $1
	`, user.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to verify email: %w", err)
	}

	user.EmailVerified = true
	user.BetaStatus = "active"
	a.log.Info("Email verified", ports.F("email", user.Email), ports.F("user_id", user.ID))

	return &user, nil
}

// CreatePasswordResetToken generates a password reset token
func (a *AuthService) CreatePasswordResetToken(ctx context.Context, userID uuid.UUID) (string, error) {
	token, err := GenerateSecureToken(32)
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(PasswordResetDuration)

	_, err = a.db.ExecContext(ctx, `
		UPDATE profiles
		SET password_reset_token = $1, password_reset_expires_at = $2
		WHERE id = $3
	`, HashToken(token), expiresAt, userID)

	if err != nil {
		return "", fmt.Errorf("failed to create reset token: %w", err)
	}

	return token, nil
}

// SetPassword sets or updates a user's password using a reset token
func (a *AuthService) SetPassword(ctx context.Context, token, password string) (*UserProfile, error) {
	tokenHash := HashToken(token)

	// Find user with this reset token
	var user UserProfile
	var expiresAt time.Time
	err := a.db.QueryRowContext(ctx, `
		SELECT id, email, password_reset_expires_at
		FROM profiles
		WHERE password_reset_token = $1
	`, tokenHash).Scan(&user.ID, &user.Email, &expiresAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invalid or expired reset token")
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check expiration
	if time.Now().After(expiresAt) {
		return nil, fmt.Errorf("reset token has expired")
	}

	// Hash password
	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Set password and clear reset token
	_, err = a.db.ExecContext(ctx, `
		UPDATE profiles
		SET password_hash = $1,
		    password_reset_token = NULL,
		    password_reset_expires_at = NULL
		WHERE id = $2
	`, passwordHash, user.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to set password: %w", err)
	}

	a.log.Info("Password set", ports.F("email", user.Email), ports.F("user_id", user.ID))
	return &user, nil
}

// SignIn authenticates a user with email and password
func (a *AuthService) SignIn(ctx context.Context, email, password string) (*UserProfile, error) {
	var user UserProfile
	err := a.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, email_verified, beta_status, subscription_tier, activated_at, created_at
		FROM profiles
		WHERE email = $1
	`, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.EmailVerified,
		&user.BetaStatus, &user.SubscriptionTier, &user.ActivatedAt, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invalid email or password")
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check if password is set
	if !user.PasswordHash.Valid || user.PasswordHash.String == "" {
		return nil, fmt.Errorf("account not activated - please check your email")
	}

	// Verify password
	if !CheckPassword(password, user.PasswordHash.String) {
		return nil, fmt.Errorf("invalid email or password")
	}

	return &user, nil
}

// CreateSession creates a new session for a user
func (a *AuthService) CreateSession(ctx context.Context, userID uuid.UUID, userAgent, ipAddress string) (string, error) {
	token, err := GenerateSecureToken(32)
	if err != nil {
		return "", err
	}

	sessionID := uuid.New()
	expiresAt := time.Now().Add(SessionDuration)

	_, err = a.db.ExecContext(ctx, `
		INSERT INTO sessions (id, user_id, token_hash, expires_at, user_agent, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, sessionID, userID, HashToken(token), expiresAt, userAgent, ipAddress)

	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	a.log.Debug("Session created", ports.F("user_id", userID))
	return token, nil
}

// ValidateSession validates a session token and returns the user ID
func (a *AuthService) ValidateSession(ctx context.Context, token string) (*Session, error) {
	tokenHash := HashToken(token)

	var session Session
	err := a.db.QueryRowContext(ctx, `
		SELECT id, user_id, expires_at, last_accessed_at
		FROM sessions
		WHERE token_hash = $1 AND expires_at > NOW()
	`, tokenHash).Scan(&session.ID, &session.UserID, &session.ExpiresAt, &session.LastAccessedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invalid or expired session")
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Update last accessed time
	_, _ = a.db.ExecContext(ctx, `
		UPDATE sessions SET last_accessed_at = NOW() WHERE id = $1
	`, session.ID)

	return &session, nil
}

// RevokeSession revokes a session by token
func (a *AuthService) RevokeSession(ctx context.Context, token string) error {
	tokenHash := HashToken(token)

	result, err := a.db.ExecContext(ctx, `
		DELETE FROM sessions WHERE token_hash = $1
	`, tokenHash)

	if err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows > 0 {
		a.log.Debug("Session revoked")
	}

	return nil
}

// RevokeAllUserSessions revokes all sessions for a user
func (a *AuthService) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	_, err := a.db.ExecContext(ctx, `
		DELETE FROM sessions WHERE user_id = $1
	`, userID)

	if err != nil {
		return fmt.Errorf("failed to revoke sessions: %w", err)
	}

	a.log.Info("All sessions revoked", ports.F("user_id", userID))
	return nil
}

// GetUserByID retrieves a user by ID
func (a *AuthService) GetUserByID(ctx context.Context, userID uuid.UUID) (*UserProfile, error) {
	var user UserProfile
	err := a.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, email_verified, beta_status, subscription_tier, activated_at, created_at
		FROM profiles
		WHERE id = $1
	`, userID).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.EmailVerified,
		&user.BetaStatus, &user.SubscriptionTier, &user.ActivatedAt, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil
}

// CleanupExpiredSessions removes expired sessions from the database
func (a *AuthService) CleanupExpiredSessions(ctx context.Context) (int64, error) {
	result, err := a.db.ExecContext(ctx, `DELETE FROM sessions WHERE expires_at < NOW()`)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup sessions: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows > 0 {
		a.log.Info("Cleaned up expired sessions", ports.F("count", rows))
	}

	return rows, nil
}

// DeletionGracePeriod is the time between deletion request and actual deletion
const DeletionGracePeriod = 24 * time.Hour

// DeletionStatus represents the current deletion state for a user
type DeletionStatus struct {
	Scheduled   bool
	ScheduledAt *time.Time
	DeletesAt   *time.Time
}

// ScheduleDeletion schedules an account for deletion after the grace period
// User keeps full access during this time and can cancel
func (a *AuthService) ScheduleDeletion(ctx context.Context, userID uuid.UUID) (string, error) {
	// Check if deletion is already scheduled
	var existingScheduled sql.NullTime
	err := a.db.QueryRowContext(ctx, `
		SELECT deletion_scheduled_at FROM profiles WHERE id = $1
	`, userID).Scan(&existingScheduled)

	if err != nil {
		return "", fmt.Errorf("failed to check existing deletion: %w", err)
	}

	if existingScheduled.Valid {
		return "", fmt.Errorf("deletion already scheduled")
	}

	// Generate cancellation token
	cancelToken, err := GenerateSecureToken(32)
	if err != nil {
		return "", fmt.Errorf("failed to generate cancel token: %w", err)
	}

	scheduledAt := time.Now()
	expiresAt := scheduledAt.Add(DeletionGracePeriod)

	_, err = a.db.ExecContext(ctx, `
		UPDATE profiles
		SET deletion_scheduled_at = $1,
		    deletion_cancel_token = $2,
		    deletion_cancel_expires_at = $3
		WHERE id = $4
	`, scheduledAt, HashToken(cancelToken), expiresAt, userID)

	if err != nil {
		return "", fmt.Errorf("failed to schedule deletion: %w", err)
	}

	a.log.Info("Account deletion scheduled",
		ports.F("user_id", userID),
		ports.F("deletes_at", expiresAt),
	)

	return cancelToken, nil
}

// CancelDeletionByToken cancels a scheduled deletion using the email cancel token
func (a *AuthService) CancelDeletionByToken(ctx context.Context, token string) error {
	tokenHash := HashToken(token)

	// Find and verify the token
	var userID uuid.UUID
	var expiresAt time.Time
	err := a.db.QueryRowContext(ctx, `
		SELECT id, deletion_cancel_expires_at
		FROM profiles
		WHERE deletion_cancel_token = $1
		  AND deletion_scheduled_at IS NOT NULL
	`, tokenHash).Scan(&userID, &expiresAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("invalid or expired cancellation token")
	}
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	// Check if past the grace period (account already deleted or being deleted)
	if time.Now().After(expiresAt) {
		return fmt.Errorf("cancellation period has expired")
	}

	// Clear deletion scheduling
	_, err = a.db.ExecContext(ctx, `
		UPDATE profiles
		SET deletion_scheduled_at = NULL,
		    deletion_cancel_token = NULL,
		    deletion_cancel_expires_at = NULL
		WHERE id = $1
	`, userID)

	if err != nil {
		return fmt.Errorf("failed to cancel deletion: %w", err)
	}

	a.log.Info("Account deletion cancelled via token", ports.F("user_id", userID))
	return nil
}

// CancelDeletionByUserID cancels a scheduled deletion for an authenticated user
func (a *AuthService) CancelDeletionByUserID(ctx context.Context, userID uuid.UUID) error {
	result, err := a.db.ExecContext(ctx, `
		UPDATE profiles
		SET deletion_scheduled_at = NULL,
		    deletion_cancel_token = NULL,
		    deletion_cancel_expires_at = NULL
		WHERE id = $1
		  AND deletion_scheduled_at IS NOT NULL
	`, userID)

	if err != nil {
		return fmt.Errorf("failed to cancel deletion: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no deletion scheduled for this account")
	}

	a.log.Info("Account deletion cancelled via dashboard", ports.F("user_id", userID))
	return nil
}

// GetDeletionStatus returns the deletion status for a user
func (a *AuthService) GetDeletionStatus(ctx context.Context, userID uuid.UUID) (*DeletionStatus, error) {
	var scheduledAt sql.NullTime

	err := a.db.QueryRowContext(ctx, `
		SELECT deletion_scheduled_at
		FROM profiles
		WHERE id = $1
	`, userID).Scan(&scheduledAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	status := &DeletionStatus{
		Scheduled: scheduledAt.Valid,
	}

	if scheduledAt.Valid {
		status.ScheduledAt = &scheduledAt.Time
		deletesAt := scheduledAt.Time.Add(DeletionGracePeriod)
		status.DeletesAt = &deletesAt
	}

	return status, nil
}

// ExecutePendingDeletions finds and deletes accounts past their grace period
// Returns the count of deleted accounts
func (a *AuthService) ExecutePendingDeletions(ctx context.Context) (int, error) {
	// Find accounts scheduled for deletion more than 24 hours ago
	rows, err := a.db.QueryContext(ctx, `
		SELECT id, email FROM profiles
		WHERE deletion_scheduled_at IS NOT NULL
		  AND deletion_scheduled_at < NOW() - INTERVAL '24 hours'
	`)
	if err != nil {
		return 0, fmt.Errorf("failed to query pending deletions: %w", err)
	}
	defer rows.Close()

	var deleted int
	for rows.Next() {
		var id uuid.UUID
		var email string
		if err := rows.Scan(&id, &email); err != nil {
			a.log.Error("Failed to scan deletion row", ports.F("error", err))
			continue
		}

		// Execute deletion in transaction
		tx, err := a.db.BeginTx(ctx, nil)
		if err != nil {
			a.log.Error("Failed to begin deletion transaction", ports.F("error", err))
			continue
		}

		// Delete from humans table if exists (cannabis data)
		_, _ = tx.ExecContext(ctx, "DELETE FROM humans WHERE email = $1", email)

		// Delete from profiles table (cascades to sessions)
		_, err = tx.ExecContext(ctx, "DELETE FROM profiles WHERE id = $1", id)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				a.log.Error("Failed to rollback deletion transaction",
					ports.F("rollback_error", rbErr),
					ports.F("original_error", err),
					ports.F("user_id", id),
				)
			}
			a.log.Error("Failed to delete profile", ports.F("error", err), ports.F("user_id", id))
			continue
		}

		if err := tx.Commit(); err != nil {
			a.log.Error("Failed to commit deletion", ports.F("error", err))
			continue
		}

		deleted++
		// Log with hashed ID for privacy
		a.log.Info("Account deleted", ports.F("user_id_hash", HashToken(id.String())))
	}

	if deleted > 0 {
		a.log.Info("Executed pending account deletions", ports.F("count", deleted))
	}

	return deleted, nil
}

// GetUserEmailByID retrieves a user's email by their ID
func (a *AuthService) GetUserEmailByID(ctx context.Context, userID uuid.UUID) (string, error) {
	var email string
	err := a.db.QueryRowContext(ctx, `
		SELECT email FROM profiles WHERE id = $1
	`, userID).Scan(&email)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("user not found")
	}
	if err != nil {
		return "", fmt.Errorf("database error: %w", err)
	}

	return email, nil
}
