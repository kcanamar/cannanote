package application

import (
	"context"
	"time"
	"backend/internal/core/domain"
	"backend/internal/core/ports"

	"github.com/google/uuid"
)

// HumanService implements business logic for human management
// This orchestrates domain entities with port contracts - no I/O dependencies
type HumanService struct {
	humanRepo ports.HumanRepository
	roleRepo  ports.RoleRepository
	authSvc   ports.AuthService
}

// NewHumanService creates a new human service with injected dependencies
// No formal DI framework - just constructor injection for security
func NewHumanService(
	humanRepo ports.HumanRepository,
	roleRepo ports.RoleRepository,
	authSvc ports.AuthService,
) *HumanService {
	return &HumanService{
		humanRepo: humanRepo,
		roleRepo:  roleRepo,
		authSvc:   authSvc,
	}
}

// CreateHuman registers a new human in the system
func (s *HumanService) CreateHuman(ctx context.Context, req CreateHumanRequest) (*domain.Human, error) {
	// Validate business rules
	if req.Username == "" {
		return nil, domain.ErrInvalidUsername
	}
	if req.Email == "" {
		return nil, domain.ErrInvalidEmail
	}

	// Check if human already exists
	existing, _ := s.humanRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, domain.ErrHumanAlreadyExists
	}

	// Create new human entity
	human := &domain.Human{
		ID:        uuid.New(),
		Username:  req.Username,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Profile: domain.HumanProfile{
			PreferredStrains:     req.PreferredStrains,
			PreferredConsumption: req.PreferredConsumption,
			ExperienceLevel:      req.ExperienceLevel,
			PublicProfile:        req.PublicProfile,
			ShareExperiences:     req.ShareExperiences,
		},
		Consent: domain.ConsentSettings{
			DataCollection:  req.ConsentToDataCollection,
			MarketingEmails: req.ConsentToMarketing,
			DataSharing:     req.ConsentToDataSharing,
			ConsentDate:     time.Now(),
			ConsentVersion:  "v1.0", // Track consent version changes
		},
	}

	// Validate the entity
	if err := human.Validate(); err != nil {
		return nil, err
	}

	// Persist the human
	if err := s.humanRepo.Create(ctx, human); err != nil {
		return nil, err
	}

	// Assign default role
	if err := s.roleRepo.AssignRole(ctx, human.ID, domain.RoleHuman, human.ID); err != nil {
		return nil, err
	}

	return human, nil
}

// GetHuman retrieves a human by ID
func (s *HumanService) GetHuman(ctx context.Context, id uuid.UUID) (*domain.Human, error) {
	human, err := s.humanRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return human, nil
}

// UpdateHumanProfile modifies human profile information
func (s *HumanService) UpdateHumanProfile(ctx context.Context, humanID uuid.UUID, req UpdateProfileRequest) (*domain.Human, error) {
	// Retrieve existing human
	human, err := s.humanRepo.GetByID(ctx, humanID)
	if err != nil {
		return nil, err
	}

	// Update profile through domain method
	profile := domain.HumanProfile{
		PreferredStrains:     req.PreferredStrains,
		PreferredConsumption: req.PreferredConsumption,
		ExperienceLevel:      req.ExperienceLevel,
		PublicProfile:        req.PublicProfile,
		ShareExperiences:     req.ShareExperiences,
		PreferredDispensary:  req.PreferredDispensary,
	}

	human.UpdateProfile(profile)

	// Persist changes
	if err := s.humanRepo.Update(ctx, human); err != nil {
		return nil, err
	}

	return human, nil
}

// UpdateConsent modifies human consent settings
func (s *HumanService) UpdateConsent(ctx context.Context, humanID uuid.UUID, req UpdateConsentRequest) (*domain.Human, error) {
	// Retrieve existing human
	human, err := s.humanRepo.GetByID(ctx, humanID)
	if err != nil {
		return nil, err
	}

	// Update consent through domain method
	consent := domain.ConsentSettings{
		DataCollection:      req.DataCollection,
		MarketingEmails:     req.MarketingEmails,
		DataSharing:         req.DataSharing,
		ConsentVersion:      req.ConsentVersion,
		MedicalDataSharing:  req.MedicalDataSharing, // Future HIPAA prep
	}

	human.GrantConsent(consent)

	// Persist changes
	if err := s.humanRepo.Update(ctx, human); err != nil {
		return nil, err
	}

	return human, nil
}

// AuthenticateHuman validates human credentials
func (s *HumanService) AuthenticateHuman(ctx context.Context, token string) (*domain.Human, error) {
	// Validate token through auth service port
	claims, err := s.authSvc.ValidateToken(ctx, token)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	// Retrieve human details
	human, err := s.humanRepo.GetByID(ctx, claims.HumanID)
	if err != nil {
		return nil, domain.ErrHumanNotFound
	}

	return human, nil
}

// CheckPermission verifies if human can perform an action
func (s *HumanService) CheckPermission(ctx context.Context, humanID uuid.UUID, action string) (bool, error) {
	// Get human
	human, err := s.humanRepo.GetByID(ctx, humanID)
	if err != nil {
		return false, err
	}

	// Check permission through domain logic
	return human.CanAccessFeature(action), nil
}

// Request/Response DTOs for the application layer
type CreateHumanRequest struct {
	Username                  string   `json:"username"`
	Email                     string   `json:"email"`
	PreferredStrains          []string `json:"preferred_strains"`
	PreferredConsumption      []string `json:"preferred_consumption"`
	ExperienceLevel           string   `json:"experience_level"`
	PublicProfile             bool     `json:"public_profile"`
	ShareExperiences          bool     `json:"share_experiences"`
	ConsentToDataCollection   bool     `json:"consent_to_data_collection"`
	ConsentToMarketing        bool     `json:"consent_to_marketing"`
	ConsentToDataSharing      bool     `json:"consent_to_data_sharing"`
}

type UpdateProfileRequest struct {
	PreferredStrains     []string `json:"preferred_strains"`
	PreferredConsumption []string `json:"preferred_consumption"`
	ExperienceLevel      string   `json:"experience_level"`
	PublicProfile        bool     `json:"public_profile"`
	ShareExperiences     bool     `json:"share_experiences"`
	PreferredDispensary  string   `json:"preferred_dispensary"`
}

type UpdateConsentRequest struct {
	DataCollection      bool   `json:"data_collection"`
	MarketingEmails     bool   `json:"marketing_emails"`
	DataSharing         bool   `json:"data_sharing"`
	ConsentVersion      string `json:"consent_version"`
	MedicalDataSharing  bool   `json:"medical_data_sharing"`
}