// Package server provides the primary HTTP server for the Mega application.
//
// It configures timeouts, establishes routing rules, and serves frontend web resources.
package server

import (
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/tarcis-io/mega/internal/config"
	"github.com/tarcis-io/mega/web"
)

const (
	defaultHost = ""
	defaultPort = "8080"
)

const (
	publicDir    = "public"
	publicPrefix = "/public/"
)

const (
	defaultReadTimeout       = 15 * time.Second
	defaultReadHeaderTimeout = 5 * time.Second
	defaultWriteTimeout      = 15 * time.Second
	defaultIdleTimeout       = 60 * time.Second
)

// Run initializes and starts the primary HTTP server.
//
// It orchestrates configuration, routing, and lifecycle management.
func Run(_ *config.Config) error {
	addr := resolveAddr()

	router, err := newRouter()
	if err != nil {
		return err
	}

	srv := newServer(addr, router)

	return serve(srv)
}

// resolveAddr determines the server binding address from environment variables.
func resolveAddr() string {
	host := os.Getenv("HOST")
	if host == "" {
		host = defaultHost
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	return net.JoinHostPort(host, port)
}

// newRouter prepares the HTTP request multiplexer and registers all routes.
func newRouter() (*http.ServeMux, error) {
	publicFS, err := fs.Sub(web.PublicFS, publicDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize public file system: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle(publicPrefix, http.StripPrefix(publicPrefix, http.FileServer(http.FS(publicFS))))

	return mux, nil
}

// newServer configures the underlying HTTP server instance with strict timeouts.
func newServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadTimeout:       getEnvDuration("HTTP_READ_TIMEOUT", defaultReadTimeout),
		ReadHeaderTimeout: getEnvDuration("HTTP_READ_HEADER_TIMEOUT", defaultReadHeaderTimeout),
		WriteTimeout:      getEnvDuration("HTTP_WRITE_TIMEOUT", defaultWriteTimeout),
		IdleTimeout:       getEnvDuration("HTTP_IDLE_TIMEOUT", defaultIdleTimeout),
	}
}

// getEnvDuration parses an environment variable into a time.Duration, falling back if missing or invalid.
func getEnvDuration(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(val)
	if err != nil {
		log.Printf("WARNING: invalid duration for %s: %q", key, val)
		return fallback
	}

	return parsed
}

// serve starts the HTTP server and blocks until it stops.
func serve(srv *http.Server) error {
	displayAddr := srv.Addr
	if displayAddr != "" && displayAddr[0] == ':' {
		displayAddr = "localhost" + displayAddr
	}

	log.Printf("HTTP server is starting on http://%s", displayAddr)
	return srv.ListenAndServe()
}
