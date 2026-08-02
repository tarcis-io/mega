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
	addr := buildAddr()

	router, err := buildRouter()
	if err != nil {
		return err
	}

	server, err := buildServer(addr, router)
	if err != nil {
		return err
	}

	log.Printf("HTTP server is starting on http://%s", addr)

	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("HTTP server stopped unexpectedly: %w", err)
	}

	return nil
}

func buildAddr() string {
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

func buildRouter() (*http.ServeMux, error) {
	publicFS, err := fs.Sub(web.PublicFS, publicDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize public file system: %w", err)
	}

	serveMux := http.NewServeMux()
	serveMux.Handle(publicPrefix, http.StripPrefix(publicPrefix, http.FileServer(http.FS(publicFS))))

	return serveMux, nil
}

func buildServer(addr string, handler http.Handler) (*http.Server, error) {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return server, nil
}
