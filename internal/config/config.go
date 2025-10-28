package config

import (
	"io"
	"log/slog"
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
