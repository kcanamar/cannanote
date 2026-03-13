package version

import (
	"backend/internal/adapters/logging"
	"backend/internal/core/ports"
)

// These variables are set at build time via ldflags
// Example: go build -ldflags "-X backend/internal/version.BuildHash=abc123"
var (
	// BuildHash is the git commit hash or build timestamp
	// Set via: -ldflags "-X backend/internal/version.BuildHash=VALUE"
	BuildHash = "dev"

	// BuildTime is when the binary was built
	// Set via: -ldflags "-X backend/internal/version.BuildTime=VALUE"
	BuildTime = "unknown"

	// Environment indicates local vs production
	// Set via: -ldflags "-X backend/internal/version.Environment=VALUE"
	Environment = "development"
)

var versionLog = logging.With("version")

// LogBuildInfo logs the current build information at startup
func LogBuildInfo() {
	versionLog.Info("Build information",
		ports.F("hash", BuildHash),
		ports.F("time", BuildTime),
		ports.F("env", Environment),
	)
}

// GetServiceWorkerVersion returns the version string for the service worker
// This ensures cache invalidation on new deployments
func GetServiceWorkerVersion() string {
	return BuildHash
}

// IsProduction returns true if running in production environment
func IsProduction() bool {
	return Environment == "production"
}

// IsDevelopment returns true if running in development environment
func IsDevelopment() bool {
	return Environment == "development" || Environment == "dev"
}
