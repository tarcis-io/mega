// Package main is the entry point for the Mega application.
//
// It runs the primary HTTP server.
package main

import (
	"log"

	"github.com/tarcis-io/mega/internal/server"
)

// main runs the primary server responsible for handling incoming HTTP requests.
func main() {
	if err := server.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
