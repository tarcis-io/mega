package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
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
	defaultServerPort              = 8080
	defaultServerReadTimeout       = 15 * time.Second
	defaultServerReadHeaderTimeout = 5 * time.Second
	defaultServerWriteTimeout      = 15 * time.Second
	defaultServerIdleTimeout       = 60 * time.Second
	defaultServerShutdownTimeout   = 30 * time.Second
)

const (
	minPort = 1
	maxPort = 65535
)

type Config struct {
	Server Server
}

func Load() (*Config, error) {
	return load(os.LookupEnv)
}

func load(lookup func(key string) (string, bool)) (*Config, error) {
	p := &parser{
		lookup: lookup,
	}

	cfg := &Config{
		Server: Server{
			Host:              p.Host(serverHostKey, defaultServerHost),
			Port:              p.Port(serverPortKey, defaultServerPort),
			ReadTimeout:       p.Timeout(serverReadTimeoutKey, defaultServerReadTimeout),
			ReadHeaderTimeout: p.Timeout(serverReadHeaderTimeoutKey, defaultServerReadHeaderTimeout),
			WriteTimeout:      p.Timeout(serverWriteTimeoutKey, defaultServerWriteTimeout),
			IdleTimeout:       p.Timeout(serverIdleTimeoutKey, defaultServerIdleTimeout),
			ShutdownTimeout:   p.Timeout(serverShutdownTimeoutKey, defaultServerShutdownTimeout),
		},
	}

	if err := p.Err(); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return cfg, nil
}

type Server struct {
	Host              string
	Port              int
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

func (s *Server) Address() string {
	return net.JoinHostPort(s.Host, strconv.Itoa(s.Port))
}

type parser struct {
	lookup func(key string) (string, bool)
	errs   []error
}

func (p *parser) String(key, fallback string) string {
	if val, ok := p.lookup(key); ok {
		return val
	}

	return fallback
}

func (p *parser) Int(key string, fallback int) int {
	valStr, ok := p.lookup(key)
	if !ok {
		return fallback
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		p.addErrorf("invalid int %s=%q: must be a number", key, valStr)
		return fallback
	}

	return val
}

func (p *parser) Host(key, fallback string) string {
	val, ok := p.lookup(key)
	if !ok {
		return fallback
	}

	if strings.ContainsAny(val, " \t\r\n") {
		p.addErrorf("invalid host %s=%q: cannot contain whitespace", key, val)
		return fallback
	}

	return val
}

func (p *parser) Port(key string, fallback int) int {
	val := p.Int(key, fallback)
	if val < minPort || val > maxPort {
		p.addErrorf("invalid port %s=%q: must be between %d and %d", key, strconv.Itoa(val), minPort, maxPort)
		return fallback
	}

	return val
}

func (p *parser) Duration(key string, fallback time.Duration) time.Duration {
	valStr, ok := p.lookup(key)
	if !ok {
		return fallback
	}

	val, err := time.ParseDuration(valStr)
	if err != nil {
		p.addErrorf("invalid duration %s=%q: %w", key, valStr, err)
		return fallback
	}

	return val
}

func (p *parser) Timeout(key string, fallback time.Duration) time.Duration {
	val := p.Duration(key, fallback)
	if val < 0 {
		p.addErrorf("invalid timeout %s=%q: must be positive", key, val.String())
		return fallback
	}

	return val
}

func (p *parser) Err() error {
	return errors.Join(p.errs...)
}

func (p *parser) addErrorf(format string, args ...any) {
	p.errs = append(p.errs, fmt.Errorf(format, args...))
}
