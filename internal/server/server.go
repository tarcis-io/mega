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
	return nil
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
