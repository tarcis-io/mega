// Package config loads and provides the application configuration.
package config

type (
	// LogLevel represents the severity level of log messages.
	LogLevel string
)

const (
	// LogLevelDebug defines the debug log level.
	// Debug messages are typically used for detailed information during
	// development.
	LogLevelDebug LogLevel = "debug"

	// LogLevelInfo defines the info log level.
	// Info messages are used to provide general information about the
	// application's execution.
	LogLevelInfo LogLevel = "info"

	// LogLevelWarn defines the warning log level.
	// Warning messages are used to indicate potential issues or unexpected
	// situations.
	LogLevelWarn LogLevel = "warn"

	// LogLevelError defines the error log level.
	// Error messages are used to report errors or exceptional conditions.
	LogLevelError LogLevel = "error"
)
