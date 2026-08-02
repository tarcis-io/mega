package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	defaultPort = "8080"
)

func Run() error {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	addr := net.JoinHostPort(host, port)

	serveMux := http.NewServeMux()

	log.Printf("HTTP server is starting on http://%s", addr)
	if err := http.ListenAndServe(addr, serveMux); err != nil {
		return fmt.Errorf("HTTP server stopped unexpectedly: %w", err)
	}

	return nil
}
