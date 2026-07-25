# Makefile for the Mega application.
#
# It automates the setup and compilation of the web assets and WebAssembly modules.

.PHONY: all setup css wasm clean

# all is the default target.
#
# It sets up the application and compiles the web assets and WebAssembly modules.
all: setup css wasm

# setup prepares dependencies for the application.
setup:
	@echo "Setting up application..."
	@mkdir -p web/public/js/wasm
	@cp $(shell tinygo env TINYGOROOT)/targets/wasm_exec.js web/public/js/wasm/wasm_exec.js

# css compiles the styles to output a CSS file.
css:
	@echo "Compiling css..."
	@mkdir -p web/public/css
	@tailwindcss -i web/src/css/app.css -o web/public/css/app.css --minify

# wasm compiles the Go packages to WebAssembly modules.
wasm:
	@echo "Compiling WebAssembly modules..."
	@mkdir -p web/public/wasm
	@tinygo build -target wasm -o web/public/wasm/home.wasm ./cmd/wasm/home

# clean removes all generated build artifacts.
clean:
	@echo "Cleaning generated build artifacts..."
	@rm -rf web/public/css/app.css
	@rm -rf web/public/js/wasm/wasm_exec.js
	@rm -rf web/public/wasm/*.wasm
