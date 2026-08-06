package config

import (
	"os"
)

const (
	serverHostKey     = "SERVER_HOST"
	defaultServerHost = "0.0.0.0"
	serverPortKey     = "SERVER_PORT"
	defaultServerPort = "8080"
)

type Config struct {
	ServerHost string
	ServerPort string
}

func Load() *Config {
	return &Config{
		ServerHost: getEnv(serverHostKey, defaultServerHost),
		ServerPort: getEnv(serverPortKey, defaultServerPort),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return fallback
}
