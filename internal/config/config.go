// Package config loads, parses, and provides the application configuration.
package config

import (
	"fmt"
	"net"
	"os"
	"time"
)

// Server configuration.
const (
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
// It returns the value (which may be empty) if the environment variable is
// set, otherwise it returns the default value.
func getEnv(key, defaultValue string) string {
	if value, isSet := os.LookupEnv(key); isSet {
		return value
	}
	return defaultValue
}

// getParsedEnv retrieves the value of the environment variable named by the
// key and parses it into a specified generic type T.
// If the environment variable is not set, the default value is used.
// On success, it returns the parsed value of type T and a nil error.
// On failure, it returns the zero value of type T and the corresponding
// parsing error.
func getParsedEnv[T any](key, defaultValue string, parse func(string) (T, error)) (T, error) {
	env := getEnv(key, defaultValue)
	parsedEnv, err := parse(env)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to parse environment variable (%s) got=%q: %w", key, env, err)
	}
	return parsedEnv, nil
}

// validateHostPort validates if the given string is a valid network address of
// the form "host:port".
// On success, it returns the original string and a nil error.
// On failure, it returns an empty string and the corresponding parsing error.
func validateHostPort(hostPort string) (string, error) {
	if _, _, err := net.SplitHostPort(hostPort); err != nil {
		return "", fmt.Errorf("invalid \"host:port\": %w", err)
	}
	return hostPort, nil
}
