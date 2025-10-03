// Package config loads, parses, and provides the application configuration.
package config

import (
	"time"
)

const (
	// Server configuration.
	serverAddressEnvKey     = "SERVER_ADDRESS"
	serverAddressEnvDefault = "localhost:8080"
)

type (
	// Config holds the application configuration.
	Config struct {
		// ServerAddress
		ServerAddress string

		// ServerReadTimeout
		ServerReadTimeout time.Duration

		// ServerReadHeaderTimeout
		ServerReadHeaderTimeout time.Duration

		// ServerWriteTimeout
		ServerWriteTimeout time.Duration

		// ServerIdleTimeout
		ServerIdleTimeout time.Duration

		// ServerShutdownTimeout
		ServerShutdownTimeout time.Duration
	}
)
