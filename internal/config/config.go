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
	defaultServerReadHeaderTimeout = 5 * time.Second
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
			Host:              p.getString(serverHostKey, defaultServerHost),
			Port:              p.getString(serverPortKey, defaultServerPort),
			ReadTimeout:       p.getDuration(serverReadTimeoutKey, defaultServerReadTimeout),
			ReadHeaderTimeout: p.getDuration(serverReadHeaderTimeoutKey, defaultServerReadHeaderTimeout),
			WriteTimeout:      p.getDuration(serverWriteTimeoutKey, defaultServerWriteTimeout),
			IdleTimeout:       p.getDuration(serverIdleTimeoutKey, defaultServerIdleTimeout),
			ShutdownTimeout:   p.getDuration(serverShutdownTimeoutKey, defaultServerShutdownTimeout),
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

func (p *parser) getString(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return fallback
}

func (p *parser) getDuration(key string, fallback time.Duration) time.Duration {
	valStr := p.getString(key, "")
	if valStr == "" {
		return fallback
	}

	val, err := time.ParseDuration(valStr)
	if err != nil {
		p.errs = append(p.errs, fmt.Errorf("failed to parse time.Duration %s=%q: %w", key, valStr, err))
		return fallback
	}

	return val
}
