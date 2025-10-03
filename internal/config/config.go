// Package config loads, parses, and provides the application configuration.
package config

const (
	// Server configuration.
	serverAddressEnvKey     = "SERVER_ADDRESS"
	serverAddressEnvDefault = "localhost:8080"
)

type (
	// Config holds the application configuration.
	Config struct {
		ServerAddress string
	}
)
