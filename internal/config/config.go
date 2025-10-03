// Package config loads, parses, and provides the application configuration.
package config

import (
	"os"
	"time"
)

const (
	// Server configuration.
	serverAddressEnvKey               = "SERVER_ADDRESS"
	serverAddressEnvDefault           = "localhost:8080"
	serverReadTimeoutEnvKey           = "SERVER_READ_TIMEOUT"
	serverReadTimeoutEnvDefault       = "5s"
	serverReadHeaderTimeoutEnvKey     = "SERVER_READ_HEADER_TIMEOUT"
	serverReadHeaderTimeoutEnvDefault = "2s"
	serverWriteTimeoutEnvKey          = "SERVER_WRITE_TIMEOUT"
	serverWriteTimeoutEnvDefault      = "10s"
	serverIdleTimeoutEnvKey           = "SERVER_IDLE_TIMEOUT"
	serverIdleTimeoutEnvDefault       = "60s"
	serverShutdownTimeoutEnvKey       = "SERVER_SHUTDOWN_TIMEOUT"
	serverShutdownTimeoutEnvDefault   = "15s"
)

type (
	// Config holds the application configuration.
	Config struct {
		// ServerAddress specifies the TCP address for the server to listen on,
		// in the form "host:port".
		// The default is "localhost:8080".
		ServerAddress string

		// ServerReadTimeout is the maximum duration for reading the entire
		// request, including the body.
		// The default is "5s".
		ServerReadTimeout time.Duration

		// ServerReadHeaderTimeout is the amount of time allowed to read
		// request headers.
		// The default is "2s".
		ServerReadHeaderTimeout time.Duration

		// ServerWriteTimeout is the maximum duration before timing out writes
		// of the response.
		// The default is "10s".
		ServerWriteTimeout time.Duration

		// ServerIdleTimeout is the maximum amount of time to wait for the next
		// request when keep-alives are enabled.
		// The default is "60s".
		ServerIdleTimeout time.Duration

		// ServerShutdownTimeout is the maximum duration to wait for a graceful
		// shutdown.
		// The default is "15s".
		ServerShutdownTimeout time.Duration
	}
)

// getEnv retrieves the value of the environment variable named by the key.
// It returns the default value if the environment variable is not set.
func getEnv(key, envDefault string) string {
	if val, found := os.LookupEnv(key); found {
		return val
	}
	return envDefault
}
