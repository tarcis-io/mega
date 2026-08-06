package config

import (
	"net"
	"os"
	"time"
)

const (
	serverHostKey                  = "SERVER_HOST"
	defaultServerHost              = "0.0.0.0"
	serverPortKey                  = "SERVER_PORT"
	defaultServerPort              = "8080"
	serverReadTimeoutKey           = "SERVER_READ_TIMEOUT"
	defaultServerReadTimeout       = 15 * time.Second
	serverReadHeaderTimeoutKey     = "SERVER_READ_HEADER_TIMEOUT"
	defaultServerReadHeaderTimeout = 5 * time.Second
	serverWriteTimeoutKey          = "SERVER_WRITE_TIMEOUT"
	defaultServerWriteTimeout      = 15 * time.Second
	serverIdleTimeoutKey           = "SERVER_IDLE_TIMEOUT"
	defaultServerIdleTimeout       = 60 * time.Second
	serverShutdownTimeoutKey       = "SERVER_SHUTDOWN_TIMEOUT"
	defaultServerShutdownTimeout   = 30 * time.Second
)

type Config struct {
	ServerHost              string
	ServerPort              string
	ServerReadTimeout       time.Duration
	ServerReadHeaderTimeout time.Duration
	ServerWriteTimeout      time.Duration
	ServerIdleTimeout       time.Duration
	ServerShutdownTimeout   time.Duration
}

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

func (c *Config) ServerAddress() string {
	return net.JoinHostPort(c.ServerHost, c.ServerPort)
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if val, ok := os.LookupEnv(key); ok {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}

	return fallback
}
