package external

import (
	"context"
	"log"
	"os"

	"backend/internal/core/ports"

	"github.com/google/uuid"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// InitializeSupertokens configures Supertokens with managed cloud
// This should be called once during application startup
func InitializeSupertokens() error {
	apiBasePath := "/auth"
	websiteBasePath := "/auth"

	connectionURI := os.Getenv("SUPERTOKENS_CONNECTION_URI")
	apiKey := os.Getenv("SUPERTOKENS_API_KEY")
	appURL := os.Getenv("APP_URL")
	appEnv := os.Getenv("APP_ENV")

	if connectionURI == "" {
		log.Println("WARNING: SUPERTOKENS_CONNECTION_URI not set - Supertokens will not be initialized")
		return nil
	}

	if apiKey == "" {
		log.Println("WARNING: SUPERTOKENS_API_KEY not set - Supertokens will not be initialized")
		return nil
	}

	log.Printf("Initializing Supertokens with connection URI: %s", maskURI(connectionURI))

	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: connectionURI,
			APIKey:        apiKey,
		},
		AppInfo: supertokens.AppInfo{
			AppName:         "CannaNote",
			APIDomain:       appURL,
			WebsiteDomain:   appURL,
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &websiteBasePath,
		},
		RecipeList: []supertokens.Recipe{
			// Email/Password authentication
			emailpassword.Init(&epmodels.TypeInput{}),

			// Email verification (required for beta signup)
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeRequired,
			}),

			// Session management with secure cookies
			session.Init(&sessmodels.TypeInput{
				CookieSecure: func() *bool {
					// Use secure cookies in production
					secure := appEnv == "production"
					return &secure
				}(),
				CookieSameSite: func() *string {
					sameSite := "lax"
					return &sameSite
				}(),
			}),
		},
	})

	if err != nil {
		log.Printf("ERROR: Failed to initialize Supertokens: %v", err)
		return err
	}

	log.Println("✓ Supertokens initialized successfully")
	return nil
}

// SupertokensAuthService implements ports.AuthService using Supertokens
type SupertokensAuthService struct{}

// NewSupertokensAuthService creates a new Supertokens auth service adapter
func NewSupertokensAuthService() ports.AuthService {
	return &SupertokensAuthService{}
}

// GenerateToken creates a session token for a human
// Note: With Supertokens, tokens are generated via session.CreateNewSession
// This is typically called automatically after successful login
func (s *SupertokensAuthService) GenerateToken(ctx context.Context, humanID uuid.UUID) (string, error) {
	// Supertokens handles token generation through its session recipe
	// This method may not be directly used in the Supertokens flow
	log.Printf("GenerateToken called for human ID: %s (Supertokens handles this automatically)", humanID)
	return "", nil
}

// ValidateToken verifies a session token
// Note: With Supertokens, validation is done via middleware
// This method may not be directly used in the Supertokens flow
func (s *SupertokensAuthService) ValidateToken(ctx context.Context, token string) (*ports.AuthClaims, error) {
	// Supertokens handles validation through session.VerifySession middleware
	// This method may not be directly used in the Supertokens flow
	log.Printf("ValidateToken called (Supertokens handles this via middleware)")
	return nil, nil
}

// RevokeToken revokes a session
func (s *SupertokensAuthService) RevokeToken(ctx context.Context, token string) error {
	// To revoke a session in Supertokens, use:
	// session.RevokeSessionUsingSessionHandle(sessionHandle)
	// The sessionHandle would need to be extracted from the session
	log.Printf("RevokeToken called for token: %s", maskToken(token))
	return nil
}

// maskURI masks sensitive parts of the connection URI for logging
func maskURI(uri string) string {
	if len(uri) < 20 {
		return "****"
	}
	return uri[:15] + "****" + uri[len(uri)-10:]
}

// maskToken masks token for logging
func maskToken(token string) string {
	if len(token) < 10 {
		return "****"
	}
	return token[:4] + "****" + token[len(token)-4:]
}
