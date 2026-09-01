// Package config provides configuration loading and validation for the application.
//
// It reads settings from environment variables, applies sensible defaults,
// and ensures all values are correctly typed and validated before the application starts.
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

// Environment variable keys and default values.
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

// Port constraints for valid TCP/UDP network ranges.
const (
	minPort = 1
	maxPort = 65535
)

// Config represents the top-level application configuration.
//
// It holds all domain-specific configuration groups required to run the application.
type Config struct {
	Server Server
}

// Load reads the application configuration from the system's environment variables.
//
// It validates all inputs and returns an error if any variables are malformed.
func Load() (*Config, error) {
	return load(os.LookupEnv)
}

// load parses the configuration using the provided lookup function.
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
		return nil, err
	}

	return cfg, nil
}

// Server represents the HTTP server configuration.
//
// It holds the network binding settings and connection timeout limits.
type Server struct {
	Host              string
	Port              int
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

// Address returns the formatted [host]:port string suitable for binding the HTTP server.
//
// It is IPv6-safe.
func (s *Server) Address() string {
	return net.JoinHostPort(s.Host, strconv.Itoa(s.Port))
}

// parser acts as a stateful accumulator for configuration parsing errors.
type parser struct {
	lookup func(key string) (string, bool)
	errs   []error
}

// String retrieves the value associated with the provided key.
//
// It returns the fallback if the key is unset.
func (p *parser) String(key, fallback string) string {
	if val, ok := p.lookup(key); ok {
		return val
	}

	return fallback
}

// Int retrieves the integer value associated with the provided key.
//
// It returns the fallback if the key is unset or fails to parse as an integer.
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

// Host retrieves the hostname string associated with the provided key.
//
// It returns the fallback if the key is unset or contains whitespace.
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

// Port retrieves the network port associated with the provided key.
//
// It returns the fallback if the key is unset, fails to parse, or falls outside the valid TCP/UDP range.
func (p *parser) Port(key string, fallback int) int {
	val := p.Int(key, fallback)
	if val < minPort || val > maxPort {
		p.addErrorf("invalid port %s=%q: must be between %d and %d", key, strconv.Itoa(val), minPort, maxPort)
		return fallback
	}

	return val
}

// Duration retrieves the time duration associated with the provided key.
//
// It returns the fallback if the key is unset or fails to parse as a duration.
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

// Timeout retrieves the time duration associated with the provided key.
//
// It returns the fallback if the key is unset, fails to parse, or is a negative value.
func (p *parser) Timeout(key string, fallback time.Duration) time.Duration {
	val := p.Duration(key, fallback)
	if val < 0 {
		p.addErrorf("invalid timeout %s=%q: must not be negative", key, val.String())
		return fallback
	}

	return val
}

// Err returns all accumulated parsing errors bundled into a single error.
func (p *parser) Err() error {
	return errors.Join(p.errs...)
}

// addErrorf formats and appends an error to the parser's internal error list.
func (p *parser) addErrorf(format string, args ...any) {
	p.errs = append(p.errs, fmt.Errorf(format, args...))
}
