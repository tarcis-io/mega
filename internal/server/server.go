package server

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/tarcis-io/mega/internal/config"
	"github.com/tarcis-io/mega/web"
)

const (
	publicDir    = "public"
	publicPrefix = "/public/"
)

func Run(cfg *config.Config) error {
	router, err := newRouter()
	if err != nil {
		return err
	}

	srv := newServer(cfg, router)
	return srv.ListenAndServe()
}

func newRouter() (*http.ServeMux, error) {
	publicFS, err := fs.Sub(web.PublicFS, publicDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize public file system: %w", err)
	}

	router := http.NewServeMux()
	router.Handle(publicPrefix, http.StripPrefix(publicPrefix, http.FileServerFS(publicFS)))

	return router, nil
}

func newServer(cfg *config.Config, router http.Handler) *http.Server {
	return &http.Server{
		Addr:              cfg.Server.Address(),
		ReadTimeout:       cfg.Server.ReadTimeout,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
		Handler:           router,
	}
}
