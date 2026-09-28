// Package config provides configuration loading and validation for the application.
//
// It reads settings from environment variables, applies sensible defaults,
// and ensures all values are correctly typed and validated.
package config

import (
	"errors"
	"fmt"
	"net"
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

// hostnameRegexp is the regular expression for network hostname validation.
var hostnameRegexp = regexp.MustCompile(`^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])(\.[a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])*$`)

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

// Addr returns the network address in the format <host>:port.
//
// It is IPv6-safe.
func (s *Server) Addr() string {
	return net.JoinHostPort(s.Host, strconv.Itoa(s.Port))
}

// parser holds the environment lookup function and accumulates validation errors.
type parser struct {
	lookup func(string) (string, bool)
	err    error
}

// Host retrieves the network hostname associated with the provided key.
//
// It returns the fallback if the key is unset or if the value contains schemes, ports,
// or invalid characters.
func (p *parser) Host(key, fallback string) string {
	return p.parse(key, fallback, func(s string) (string, error) {
		val := s

		if val == "" {
			return val, nil
		}

		if strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]") {
			val = val[1 : len(val)-1]
		}

		if strings.Contains(val, "://") {
			return "", errors.New("must not contain a URL scheme (e.g., http://)")
		}

		if _, _, err := net.SplitHostPort(val); err == nil {
			return "", errors.New("must not include a port")
		}

		if net.ParseIP(val) == nil && !hostnameRegexp.MatchString(val) {
			return "", errors.New("must be a valid IP address or an RFC 1123 hostname")
		}

		return val, nil
	})
}

// Port retrieves the network port associated with the provided key.
//
// It returns the fallback if the key is unset, if the value fails to parse as an integer,
// or if it falls outside the valid TCP/UDP range.
func (p *parser) Port(key string, fallback int) int {
	return p.parse(key, fallback, func(s string) (int, error) {
		val, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}

		if val < minPort || val > maxPort {
			return 0, fmt.Errorf("must be between %d and %d", minPort, maxPort)
		}

		return val, nil
	})
}

// Timeout retrieves the time duration associated with the provided key.
//
// It returns the fallback if the key is unset, if the value fails to parse as a [time.Duration],
// or if it is negative.
func (p *parser) Timeout(key string, fallback time.Duration) time.Duration {
	return p.parse(key, fallback, func(s string) (time.Duration, error) {
		val, err := time.ParseDuration(s)
		if err != nil {
			return 0, err
		}

		if val < 0 {
			return 0, errors.New("must be non-negative")
		}

		return val, nil
	})
}

// Err returns all accumulated parsing errors bundled into a single error.
func (p *parser) Err() error {
	return p.err
}

// parse is a generic helper that fetches, sanitizes, and evaluates an environment variable.
//
// It returns the fallback if the key is unset. If the parsing or validation fails,
// it records the error internally and returns the fallback.
func (p *parser) parse[T any](key string, fallback T, parseFn func(string) (T, error)) T {
	raw, ok := p.lookup(key)
	if !ok {
		return fallback
	}

	val, err := parseFn(strings.TrimSpace(raw))
	if err != nil {
		p.err = errors.Join(p.err, fmt.Errorf("invalid configuration %s=%q: %v", key, raw, err))
		return fallback
	}

	return val
}
