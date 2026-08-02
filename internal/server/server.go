package server

import (
	"fmt"
	"log"
	"net/http"
)

func Run() error {
	addr := ":8080"

	serveMux := http.NewServeMux()

	log.Printf("HTTP server is starting on http://%s", addr)
	if err := http.ListenAndServe(addr, serveMux); err != nil {
		return fmt.Errorf("HTTP server stopped unexpectedly: %w", err)
	}

	return nil
}
