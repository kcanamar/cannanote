package main

import (
	"fmt"

	"backend/internal/adapters/logging"
	"backend/internal/core/ports"
	"backend/internal/server"
	"backend/internal/version"
)

func main() {
	// Initialize logging first - before any other initialization
	// Configured via LOG_LEVEL and LOG_COLOR environment variables
	logging.InitializeFromEnv()

	log := logging.With("main")
	log.Info("Starting CannaNote server")

	// Log build information for debugging and verification
	version.LogBuildInfo()
	log.Info("Service worker cache version",
		ports.F("sw_version", version.GetServiceWorkerVersion()),
	)

	srv := server.NewServer()

	err := srv.ListenAndServe()
	if err != nil {
		log.Error("Server failed to start", ports.F("error", err))
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
