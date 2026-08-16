package config

import (
	"errors"
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
	Server *Server
}

func New() *Config {
	return &Config{
		Server: &Server{
			Host: defaultServerHost,
			Port: defaultServerPort,
		},
	}
}

func Load() (*Config, error) {
	return New(), nil
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
