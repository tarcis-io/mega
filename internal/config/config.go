package config

import (
	"net"
)

type Config struct {
	Server Server
}

func New() *Config {
	return &Config{}
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
