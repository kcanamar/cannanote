package logging

import (
	"sync"

	"backend/internal/core/ports"
)

var (
	globalLogger ports.Logger
	once         sync.Once
)

// Initialize sets up the global logger instance
// This should be called once at application startup
func Initialize(logger ports.Logger) {
	once.Do(func() {
		globalLogger = logger
	})
}

// InitializeFromEnv sets up the global logger from environment variables
// This is a convenience function for typical startup
func InitializeFromEnv() {
	Initialize(NewConsoleLoggerFromEnv())
}

// L returns the global logger instance
// If not initialized, returns a default console logger
func L() ports.Logger {
	if globalLogger == nil {
		// Fallback to default logger if not initialized
		return NewConsoleLoggerFromEnv()
	}
	return globalLogger
}

// Debug logs a debug message using the global logger
func Debug(msg string, fields ...ports.Field) {
	L().Debug(msg, fields...)
}

// Info logs an info message using the global logger
func Info(msg string, fields ...ports.Field) {
	L().Info(msg, fields...)
}

// Warn logs a warning message using the global logger
func Warn(msg string, fields ...ports.Field) {
	L().Warn(msg, fields...)
}

// Error logs an error message using the global logger
func Error(msg string, fields ...ports.Field) {
	L().Error(msg, fields...)
}

// With returns a component-tagged logger from the global instance
func With(component string) ports.Logger {
	return L().With(component)
}
