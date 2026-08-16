package config

import (
	"errors"
	"fmt"
	"net"
	"os"
)

const (
	serverHostKey = "SERVER_HOST"
	serverPortKey = "SERVER_PORT"

	defaultServerHost = ""
	defaultServerPort = "8080"
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
	Host string
	Port string
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
