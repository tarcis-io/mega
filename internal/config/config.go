// Package config loads, parses, and provides the application configuration.
package config

import (
	"io"
	"log/slog"
)

// Supported log formats.
const (
	LogFormatJSON = "json"
	LogFormatText = "text"
)

// Supported log outputs.
const (
	logOutputStdout = "stdout"
	logOutputStderr = "stderr"
)

// Log output file mode.
const (
	logOutputFileMode = 0644
)

type (
	// Config holds the application configuration.
	Config struct {
		// LogLevel specifies the minimum level of log messages to output.
		// It will be slog.LevelDebug, slog.LevelInfo, slog.LevelWarn,
		// slog.LevelError, or a numerical level.
		// Default: slog.LevelInfo.
		LogLevel slog.Level

		// LogFormat specifies the output format of log messages.
		// It can be either "json" or "text".
		// Default: "json".
		LogFormat string

		// LogOutput specifies the destination of log messages.
		// It will be os.Stdout, os.Stderr, or an *os.File.
		// Default: os.Stdout.
		LogOutput io.Writer
	}
)

// New creates and returns a new Config instance by loading, parsing, and
// validating the application configuration.
// It returns a fully initialized Config on success.
// If any configuration value fails to parse or is invalid, it returns a nil
// Config and a single error that aggregates all errors found.
func New() (*Config, error) {
	return nil, nil
}

type (
	// parser is a helper struct for parsing the application configuration.
	parser struct {
		// errs holds all errors found during parsing.
		errs []error
	}
)

// newParser creates and returns a new initialized parser instance.
func newParser() *parser {
	return &parser{}
}
