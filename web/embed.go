package web

import (
	"embed"
)

var (
	//go:embed public
	PublicFS embed.FS
)
