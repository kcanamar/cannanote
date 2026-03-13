package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"backend/internal/core/ports"
)

// ConsoleLogger implements ports.Logger with human-readable console output
type ConsoleLogger struct {
	level     ports.Level
	component string
	fields    []ports.Field
	mu        sync.Mutex
	useColor  bool
}

// Config holds configuration for the console logger
type Config struct {
	Level    ports.Level
	UseColor bool
}

// NewConsoleLogger creates a new console logger with the given configuration
func NewConsoleLogger(cfg Config) *ConsoleLogger {
	return &ConsoleLogger{
		level:    cfg.Level,
		useColor: cfg.UseColor,
		fields:   make([]ports.Field, 0),
	}
}

// NewConsoleLoggerFromEnv creates a console logger configured from environment variables
// LOG_LEVEL: debug, info, warn, error (default: info)
// LOG_COLOR: true, false (default: true for TTY, false otherwise)
func NewConsoleLoggerFromEnv() *ConsoleLogger {
	levelStr := os.Getenv("LOG_LEVEL")
	if levelStr == "" {
		levelStr = "info"
	}

	useColor := true
	if colorStr := os.Getenv("LOG_COLOR"); colorStr != "" {
		useColor = colorStr == "true" || colorStr == "1"
	} else {
		// Auto-detect: use color if stdout is a terminal
		if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
			useColor = false
		}
	}

	return NewConsoleLogger(Config{
		Level:    ports.ParseLevel(levelStr),
		UseColor: useColor,
	})
}

// Debug logs a debug-level message
func (l *ConsoleLogger) Debug(msg string, fields ...ports.Field) {
	l.log(ports.LevelDebug, msg, fields...)
}

// Info logs an info-level message
func (l *ConsoleLogger) Info(msg string, fields ...ports.Field) {
	l.log(ports.LevelInfo, msg, fields...)
}

// Warn logs a warning-level message
func (l *ConsoleLogger) Warn(msg string, fields ...ports.Field) {
	l.log(ports.LevelWarn, msg, fields...)
}

// Error logs an error-level message
func (l *ConsoleLogger) Error(msg string, fields ...ports.Field) {
	l.log(ports.LevelError, msg, fields...)
}

// With returns a new Logger with the given component tag
func (l *ConsoleLogger) With(component string) ports.Logger {
	return &ConsoleLogger{
		level:     l.level,
		component: component,
		fields:    l.fields,
		useColor:  l.useColor,
	}
}

// WithFields returns a new Logger with preset fields
func (l *ConsoleLogger) WithFields(fields ...ports.Field) ports.Logger {
	newFields := make([]ports.Field, len(l.fields)+len(fields))
	copy(newFields, l.fields)
	copy(newFields[len(l.fields):], fields)

	return &ConsoleLogger{
		level:     l.level,
		component: l.component,
		fields:    newFields,
		useColor:  l.useColor,
	}
}

// log handles the actual logging with proper formatting
func (l *ConsoleLogger) log(level ports.Level, msg string, fields ...ports.Field) {
	// Skip if below configured level
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Get caller info (skip 2 frames: log -> Debug/Info/etc -> actual caller)
	_, file, line, ok := runtime.Caller(2)
	caller := "unknown:0"
	if ok {
		caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	// Build the log line
	var sb strings.Builder

	// Timestamp: 2024-01-15 10:30:45
	sb.WriteString(time.Now().Format("2006-01-02 15:04:05"))
	sb.WriteString(" ")

	// Level with optional color
	levelStr := l.formatLevel(level)
	sb.WriteString(levelStr)
	sb.WriteString(" ")

	// Component tag: [email]
	if l.component != "" {
		sb.WriteString("[")
		sb.WriteString(l.component)
		sb.WriteString("] ")
	}

	// Caller: resend_service.go:42
	sb.WriteString(caller)
	sb.WriteString(" ")

	// Message
	sb.WriteString(msg)

	// Combine preset fields and call-specific fields
	allFields := append(l.fields, fields...)
	if len(allFields) > 0 {
		for _, f := range allFields {
			sb.WriteString(" ")
			sb.WriteString(f.Key)
			sb.WriteString("=")
			sb.WriteString(formatValue(f.Value))
		}
	}

	sb.WriteString("\n")

	// Write to stdout (info, debug) or stderr (warn, error)
	if level >= ports.LevelWarn {
		os.Stderr.WriteString(sb.String())
	} else {
		os.Stdout.WriteString(sb.String())
	}
}

// formatLevel returns the level string with optional ANSI color codes
func (l *ConsoleLogger) formatLevel(level ports.Level) string {
	if !l.useColor {
		return padLevel(level.String())
	}

	// ANSI color codes
	var color string
	switch level {
	case ports.LevelDebug:
		color = "\033[36m" // Cyan
	case ports.LevelInfo:
		color = "\033[32m" // Green
	case ports.LevelWarn:
		color = "\033[33m" // Yellow
	case ports.LevelError:
		color = "\033[31m" // Red
	default:
		color = "\033[0m" // Reset
	}

	reset := "\033[0m"
	return color + padLevel(level.String()) + reset
}

// padLevel pads the level string to 5 characters for alignment
func padLevel(s string) string {
	if len(s) >= 5 {
		return s
	}
	return s + strings.Repeat(" ", 5-len(s))
}

// formatValue converts a value to a string representation
func formatValue(v any) string {
	switch val := v.(type) {
	case string:
		// Quote strings with spaces
		if strings.Contains(val, " ") {
			return fmt.Sprintf("%q", val)
		}
		return val
	case error:
		return fmt.Sprintf("%q", val.Error())
	default:
		return fmt.Sprintf("%v", val)
	}
}

// Ensure ConsoleLogger implements ports.Logger
var _ ports.Logger = (*ConsoleLogger)(nil)
