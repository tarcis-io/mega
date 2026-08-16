package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	serverHostKey              = "SERVER_HOST"
	serverPortKey              = "SERVER_PORT"
	serverReadTimeoutKey       = "SERVER_READ_TIMEOUT"
	serverReadHeaderTimeoutKey = "SERVER_READ_HEADER_TIMEOUT"
	serverWriteTimeoutKey      = "SERVER_WRITE_TIMEOUT"
	serverIdleTimeoutKey       = "SERVER_IDLE_TIMEOUT"
	serverShutdownTimeoutKey   = "SERVER_SHUTDOWN_TIMEOUT"

	defaultServerHost              = ""
	defaultServerPort              = "8080"
	defaultServerReadTimeout       = 15 * time.Second
	defaultServerReadHeaderTimeout = 05 * time.Second
	defaultServerWriteTimeout      = 15 * time.Second
	defaultServerIdleTimeout       = 60 * time.Second
	defaultServerShutdownTimeout   = 30 * time.Second
)

type Config struct {
	Server Server
}

func Load() (*Config, error) {
	p := &parser{}

	cfg := &Config{
		Server: Server{
			Host: p.string(serverHostKey, defaultServerHost),
			Port: p.string(serverPortKey, defaultServerPort),
		},
	}

	if err := p.Err(); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return cfg, nil
}

type Server struct {
	Host              string
	Port              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

func (s *Server) Address() string {
	return net.JoinHostPort(s.Host, s.Port)
}

type parser struct {
	errs []error
}

func (p *parser) Err() error {
	return errors.Join(p.errs...)
}

func (p *parser) string(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return fallback
}
