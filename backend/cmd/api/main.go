package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Start server in a goroutine so we can listen for shutdown signals
	go func() {
		log.Info("Server listening", ports.F("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Server error", ports.F("error", err))
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal (Ctrl+C or SIGTERM from deploy systems)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	log.Info("Shutdown signal received, draining connections...",
		ports.F("signal", sig.String()),
	)

	// Give in-flight requests 10 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Shutdown error", ports.F("error", err))
		os.Exit(1)
	}

	log.Info("Server stopped gracefully")
}
