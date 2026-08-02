// Package web provides embedded file systems for the application's frontend resources.
//
// It embeds web-related assets directly into the compiled binary.
package web

import (
	"embed"
)

var (
	// PublicFS contains public-facing static files such as images.
	//go:embed public
	PublicFS embed.FS
)
