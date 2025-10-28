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
	// Config
	Config struct {
		// LogLevel
		LogLevel slog.Level

		// LogFormat
		LogFormat string

		// LogOutput
		LogOutput io.Writer
	}
)

// New
func New() (*Config, error) {
	return nil, nil
}
