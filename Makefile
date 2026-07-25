# Makefile for the Mega application.
#
# It automates the setup and compilation of the web assets and WebAssembly modules.

TINYGO      ?= tinygo
TAILWINDCSS ?= tailwindcss

.PHONY: all setup build build-css build-wasm clean

# all is the default target.
#
# It sets up the environment and compiles all application assets.
all: setup build

# build executes all compilation targets for the application.
build: build-css build-wasm

# setup prepares dependencies for the application.
setup:
	@echo "Setting up application..."
	@mkdir -p web/public/js/wasm
	@cp $(shell $(TINYGO) env TINYGOROOT)/targets/wasm_exec.js web/public/js/wasm/wasm_exec.js

# build-css compiles the styles to output a CSS file.
build-css:
	@echo "Compiling CSS..."
	@mkdir -p web/public/css
	@$(TAILWINDCSS) -i web/src/css/app.css -o web/public/css/app.css --minify

# build-wasm compiles the Go packages to WebAssembly modules.
build-wasm:
	@echo "Compiling WebAssembly modules..."
	@mkdir -p web/public/wasm
	@$(TINYGO) build -target wasm -o web/public/wasm/home.wasm cmd/wasm/home

# clean removes all generated build artifacts.
clean:
	@echo "Cleaning generated build artifacts..."
	@rm -f web/public/css/app.css web/public/js/wasm/wasm_exec.js web/public/wasm/*.wasm
