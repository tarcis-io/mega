//go:build js && wasm

// Package main is the entry point for the about page WebAssembly module.
//
// It renders the user interface and keeps the runtime alive.
package main

import (
	"github.com/tarcis-io/mega/internal/wasm/page/about"
)

// main renders the about page user interface.
//
// It blocks until the application is terminated to allow DOM event handling.
func main() {
	about.Render()
	select {}
}
