package ports

// Level represents the severity of a log message
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// String returns the string representation of a log level
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel converts a string to a Level
func ParseLevel(s string) Level {
	switch s {
	case "debug", "DEBUG":
		return LevelDebug
	case "info", "INFO":
		return LevelInfo
	case "warn", "WARN", "warning", "WARNING":
		return LevelWarn
	case "error", "ERROR":
		return LevelError
	default:
		return LevelInfo
	}
}

// Field represents a structured key-value pair for logging
type Field struct {
	Key   string
	Value any
}

// F is a convenience function to create a Field
func F(key string, value any) Field {
	return Field{Key: key, Value: value}
}

// Logger defines the contract for application logging
// This is a port - adapters implement this interface
// Domain layer should NOT use this (returns errors instead)
// Application and Adapter layers use this freely
type Logger interface {
	// Debug logs a debug-level message (development/troubleshooting)
	Debug(msg string, fields ...Field)

	// Info logs an info-level message (normal operations)
	Info(msg string, fields ...Field)

	// Warn logs a warning-level message (unexpected but handled)
	Warn(msg string, fields ...Field)

	// Error logs an error-level message (failures requiring attention)
	Error(msg string, fields ...Field)

	// With returns a new Logger with the given component tag
	// Example: logger.With("email") logs as [email]
	With(component string) Logger

	// WithFields returns a new Logger with preset fields
	// These fields are included in every subsequent log call
	WithFields(fields ...Field) Logger
}
