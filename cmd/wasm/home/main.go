//go:build js && wasm

// Package main is the entry point for the home page WebAssembly module.
//
// It renders the user interface and keeps the runtime alive.
package main

import (
	"github.com/tarcis-io/mega/internal/wasm/page/home"
)

// main renders the home page user interface.
//
// It blocks until the application is terminated to allow DOM event handling.
func main() {
	home.Render()
	select {}
}
