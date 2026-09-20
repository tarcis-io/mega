// Package config provides configuration loading and validation for the application.
//
// It reads settings from environment variables, applies sensible defaults,
// and ensures all values are correctly typed and validated.
package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Environment variable keys.
const (
	envServerHost              = "SERVER_HOST"
	envServerPort              = "SERVER_PORT"
	envServerReadTimeout       = "SERVER_READ_TIMEOUT"
	envServerReadHeaderTimeout = "SERVER_READ_HEADER_TIMEOUT"
	envServerWriteTimeout      = "SERVER_WRITE_TIMEOUT"
	envServerIdleTimeout       = "SERVER_IDLE_TIMEOUT"
	envServerShutdownTimeout   = "SERVER_SHUTDOWN_TIMEOUT"
)

// Default configuration values.
const (
	defaultServerHost              = ""
	defaultServerPort              = 8080
	defaultServerReadTimeout       = 15 * time.Second
	defaultServerReadHeaderTimeout = 5 * time.Second
	defaultServerWriteTimeout      = 15 * time.Second
	defaultServerIdleTimeout       = 60 * time.Second
	defaultServerShutdownTimeout   = 30 * time.Second
)

// Constraints for valid TCP/UDP network port ranges.
const (
	minPort = 0
	maxPort = 65535
)

// hostnameRegexp is the expression for network hostname validation.
var hostnameRegexp = regexp.MustCompile(`^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])(\.[a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])*$`)

// Config represents the top-level application configuration.
//
// It holds all domain-specific configuration groups required to run the application.
type Config struct {
	Server Server
}

// Load reads the application configuration from the system's environment variables.
//
// It validates all inputs and returns a joined error containing all validation failures
// if multiple variables are malformed.
func Load() (*Config, error) {
	return load(os.LookupEnv)
}

func load(lookup func(key string) (string, bool)) (*Config, error) {
	p := &parser{
		lookup: lookup,
	}

	cfg := &Config{
		Server: loadServer(p),
	}

	if err := p.Err(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadServer(p *parser) Server {
	return Server{
		Host:              p.Host(envServerHost, defaultServerHost),
		Port:              p.Port(envServerPort, defaultServerPort),
		ReadTimeout:       p.NonNegativeDuration(envServerReadTimeout, defaultServerReadTimeout),
		ReadHeaderTimeout: p.NonNegativeDuration(envServerReadHeaderTimeout, defaultServerReadHeaderTimeout),
		WriteTimeout:      p.NonNegativeDuration(envServerWriteTimeout, defaultServerWriteTimeout),
		IdleTimeout:       p.NonNegativeDuration(envServerIdleTimeout, defaultServerIdleTimeout),
		ShutdownTimeout:   p.NonNegativeDuration(envServerShutdownTimeout, defaultServerShutdownTimeout),
	}
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

// Addr returns the formatted <host>:port string suitable for binding the HTTP server.
//
// It is IPv6-safe.
func (s *Server) Addr() string {
	return net.JoinHostPort(s.Host, strconv.Itoa(s.Port))
}

// parser acts as a stateful accumulator for configuration parsing errors.
type parser struct {
	lookup func(key string) (string, bool)
	err    error
}

// String retrieves the string value associated with the provided key.
//
// It returns the fallback if the key is unset.
func (p *parser) String(key, fallback string) string {
	if val, ok := p.get(key); ok {
		return val
	}

	return fallback
}

// Int retrieves the integer value associated with the provided key.
//
// It returns the fallback if the key is unset or if the value fails to parse as an integer.
func (p *parser) Int(key string, fallback int) int {
	valStr, ok := p.get(key)
	if !ok {
		return fallback
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		p.addErrorf("invalid int %s=%q: must be a valid integer", key, valStr)
		return fallback
	}

	return val
}

// Host retrieves the network hostname associated with the provided key.
//
// It returns the fallback if the key is unset or if the value contains schemes, ports,
// or invalid characters.
func (p *parser) Host(key, fallback string) string {
	valStr, ok := p.get(key)
	if !ok {
		return fallback
	}

	val := valStr
	if strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]") {
		val = val[1 : len(val)-1]
	}

	if strings.Contains(val, "://") {
		p.addErrorf("invalid host %s=%q: must not contain a URL scheme (e.g., http://)", key, valStr)
		return fallback
	}

	if _, _, err := net.SplitHostPort(val); err == nil {
		p.addErrorf("invalid host %s=%q: must not include a port", key, valStr)
		return fallback
	}

	if net.ParseIP(val) == nil {
		if !hostnameRegexp.MatchString(val) {
			p.addErrorf("invalid host %s=%q: must be a valid IP address or a hostname", key, valStr)
			return fallback
		}
	}

	return val
}

// Port retrieves the network port associated with the provided key.
//
// It returns the fallback if the key is unset or if the value fails to parse as an integer,
// or if it falls outside the valid TCP/UDP range.
func (p *parser) Port(key string, fallback int) int {
	val := p.Int(key, fallback)
	if val < minPort || val > maxPort {
		raw, _ := p.get(key)
		p.addErrorf("invalid port %s=%q: must be between %d and %d", key, raw, minPort, maxPort)
		return fallback
	}

	return val
}

// Duration retrieves the time duration associated with the provided key.
//
// If the parsed value is a unitless number, it is implicitly treated as seconds.
// It returns the fallback if the key is unset or if the value fails to parse as a [time.Duration].
func (p *parser) Duration(key string, fallback time.Duration) time.Duration {
	rawStr, ok := p.get(key)
	if !ok {
		return fallback
	}

	valStr := rawStr
	if _, err := strconv.ParseFloat(valStr, 64); err == nil {
		valStr += "s"
	}

	val, err := time.ParseDuration(valStr)
	if err != nil {
		p.addErrorf("invalid duration %s=%q: %w", key, rawStr, err)
		return fallback
	}

	return val
}

// NonNegativeDuration retrieves the time duration associated with the provided key.
//
// It returns the fallback if the key is unset, if the value fails to parse as a [time.Duration],
// or if it is negative.
func (p *parser) NonNegativeDuration(key string, fallback time.Duration) time.Duration {
	val := p.Duration(key, fallback)
	if val < 0 {
		raw, _ := p.get(key)
		p.addErrorf("invalid duration %s=%q: must be non-negative", key, raw)
		return fallback
	}

	return val
}

// Err returns all accumulated parsing errors bundled into a single error.
func (p *parser) Err() error {
	return p.err
}

// get retrieves and sanitizes the environment variable value associated with the provided key.
func (p *parser) get(key string) (string, bool) {
	val, ok := p.lookup(key)
	if !ok {
		return "", false
	}

	return strings.TrimSpace(val), true
}

// addErrorf formats and merges an error into the parser's internal error state.
func (p *parser) addErrorf(format string, args ...any) {
	p.err = errors.Join(p.err, fmt.Errorf(format, args...))
}
