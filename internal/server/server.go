package server

import (
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/tarcis-io/mega/web"
)

const (
	publicDir    = "public"
	publicPrefix = "/public/"
	defaultHost  = "localhost"
	defaultPort  = "8080"
)

func Run() error {
	publicFS, err := fs.Sub(web.PublicFS, publicDir)
	if err != nil {
		return fmt.Errorf("failed to initialize public file system: %w", err)
	}

	serveMux := http.NewServeMux()
	serveMux.Handle(publicPrefix, http.StripPrefix(publicPrefix, http.FileServer(http.FS(publicFS))))

	addr := getServerAddress()
	log.Printf("HTTP server is starting on http://%s", addr)

	if err := http.ListenAndServe(addr, serveMux); err != nil {
		return fmt.Errorf("HTTP server stopped unexpectedly: %w", err)
	}

	return nil
}

func getServerAddress() string {
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
