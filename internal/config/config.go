// Package config provides environment-based configuration management for the application.
//
// It handles fetching configuration values from the system environment
// and enforces safe default values when variables are missing or invalid.
package config

import (
	"net"
	"os"
	"time"
)

const (
	// serverHostKey is the environment variable name for the server host.
	serverHostKey = "SERVER_HOST"

	// defaultServerHost is the default value for the server host.
	defaultServerHost = ""

	// serverPortKey is the environment variable name for the server port.
	serverPortKey = "SERVER_PORT"

	// defaultServerPort is the default value for the server port.
	defaultServerPort = "8080"

	// serverReadTimeoutKey is the environment variable name for the server read timeout.
	serverReadTimeoutKey = "SERVER_READ_TIMEOUT"

	// defaultServerReadTimeout is the default value for the server read timeout.
	defaultServerReadTimeout = 15 * time.Second

	// serverReadHeaderTimeoutKey is the environment variable name for the server read header timeout.
	serverReadHeaderTimeoutKey = "SERVER_READ_HEADER_TIMEOUT"

	// defaultServerReadHeaderTimeout is the default value for the server read header timeout.
	defaultServerReadHeaderTimeout = 5 * time.Second

	// serverWriteTimeoutKey is the environment variable name for the server write timeout.
	serverWriteTimeoutKey = "SERVER_WRITE_TIMEOUT"

	// defaultServerWriteTimeout is the default value for the server write timeout.
	defaultServerWriteTimeout = 15 * time.Second

	// serverIdleTimeoutKey is the environment variable name for the server idle timeout.
	serverIdleTimeoutKey = "SERVER_IDLE_TIMEOUT"

	// defaultServerIdleTimeout is the default value for the server idle timeout.
	defaultServerIdleTimeout = 60 * time.Second

	// serverShutdownTimeoutKey is the environment variable name for the server shutdown timeout.
	serverShutdownTimeoutKey = "SERVER_SHUTDOWN_TIMEOUT"

	// defaultServerShutdownTimeout is the default value for the server shutdown timeout.
	defaultServerShutdownTimeout = 30 * time.Second
)

// Config holds all runtime configuration parameters for the application.
type Config struct {
	// ServerHost is the network interface for the HTTP server.
	//
	// Default: ""
	ServerHost string

	// ServerPort is the port for the HTTP server.
	//
	// Default: "8080"
	ServerPort string

	// ServerReadTimeout is the maximum duration for reading the entire HTTP request.
	//
	// Default: 15 * time.Second
	ServerReadTimeout time.Duration

	// ServerReadHeaderTimeout is the maximum duration for reading the HTTP request headers.
	//
	// Default: 5 * time.Second
	ServerReadHeaderTimeout time.Duration

	// ServerWriteTimeout is the maximum duration for writing the HTTP response.
	//
	// Default: 15 * time.Second
	ServerWriteTimeout time.Duration

	// ServerIdleTimeout is the maximum duration for an idle HTTP connection.
	//
	// Default: 60 * time.Second
	ServerIdleTimeout time.Duration

	// ServerShutdownTimeout is the maximum duration for shutting down the HTTP server.
	//
	// Default: 30 * time.Second
	ServerShutdownTimeout time.Duration
}

// Load reads the configuration from environment variables and returns a new [Config] instance.
func Load() *Config {
	return &Config{
		ServerHost:              getEnv(serverHostKey, defaultServerHost),
		ServerPort:              getEnv(serverPortKey, defaultServerPort),
		ServerReadTimeout:       getEnvDuration(serverReadTimeoutKey, defaultServerReadTimeout),
		ServerReadHeaderTimeout: getEnvDuration(serverReadHeaderTimeoutKey, defaultServerReadHeaderTimeout),
		ServerWriteTimeout:      getEnvDuration(serverWriteTimeoutKey, defaultServerWriteTimeout),
		ServerIdleTimeout:       getEnvDuration(serverIdleTimeoutKey, defaultServerIdleTimeout),
		ServerShutdownTimeout:   getEnvDuration(serverShutdownTimeoutKey, defaultServerShutdownTimeout),
	}
}

// ServerAddress returns the server address in the format "host:port".
func (c *Config) ServerAddress() string {
	return net.JoinHostPort(c.ServerHost, c.ServerPort)
}

// getEnv retrieves the value of the environment variable named by the key.
//
// If the environment variable is not present, the fallback value is returned.
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return fallback
}

// getEnvDuration retrieves the value of the environment variable named by the key
// and attempts to parse it as a [time.Duration].
//
// If the environment variable is not present or cannot be parsed, the fallback value is returned.
func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if val, ok := os.LookupEnv(key); ok {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}

	return fallback
}
