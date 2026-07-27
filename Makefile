# Makefile for the Mega application.
#
# It automates the setup and compilation of the web assets and WebAssembly modules.

# Build tool executables.
TINYGO      ?= tinygo
TAILWINDCSS ?= tailwindcss

# Resolved path to the TinyGo installation directory.
TINYGOROOT := $(shell $(TINYGO) env TINYGOROOT)

# Hierarchical directory structure for source code and compiled assets.
CMD                = cmd
CMD_WASM           = $(CMD)/wasm
CMD_WASM_ABOUT     = $(CMD_WASM)/about
CMD_WASM_HOME      = $(CMD_WASM)/home
WEB                = web
WEB_PUBLIC         = $(WEB)/public
WEB_PUBLIC_CSS     = $(WEB_PUBLIC)/css
WEB_PUBLIC_JS      = $(WEB_PUBLIC)/js
WEB_PUBLIC_JS_WASM = $(WEB_PUBLIC_JS)/wasm
WEB_PUBLIC_WASM    = $(WEB_PUBLIC)/wasm
WEB_SRC            = $(WEB)/src
WEB_SRC_CSS        = $(WEB_SRC)/css

# Source files used as inputs for builds and environment setup.
APP_CSS_INPUT      = $(WEB_SRC_CSS)/app.css
WASM_EXEC_JS_INPUT = $(TINYGOROOT)/targets/wasm_exec.js

# Generated artifacts and output files from the build process.
APP_CSS_OUTPUT      = $(WEB_PUBLIC_CSS)/app.css
ABOUT_WASM_OUTPUT   = $(WEB_PUBLIC_WASM)/about.wasm
HOME_WASM_OUTPUT    = $(WEB_PUBLIC_WASM)/home.wasm
WASM_EXEC_JS_OUTPUT = $(WEB_PUBLIC_JS_WASM)/wasm_exec.js

# Non-file action aliases (phony targets).
.PHONY: all setup build build-css build-wasm clean

# all is the default target.
#
# It sets up the environment and compiles all application assets.
all: setup build

# build executes all compilation targets for the application.
build: build-css build-wasm

# setup prepares dependencies for the application.
setup: $(WASM_EXEC_JS_OUTPUT)

# Copy the WebAssembly execution script from TinyGo if it doesn't exist or is updated.
$(WASM_EXEC_JS_OUTPUT): $(WASM_EXEC_JS_INPUT)
	@echo "Setting up wasm_exec.js..."
	@mkdir -p $(dir $@)
	@cp $< $@

# build-css compiles the styles to output a CSS file.
build-css: $(APP_CSS_OUTPUT)

# Compile Tailwind CSS only if the source file changes.
$(APP_CSS_OUTPUT): $(APP_CSS_INPUT)
	@echo "Compiling CSS..."
	@mkdir -p $(dir $@)
	@$(TAILWINDCSS) -i $< -o $@ --minify

# build-wasm compiles the Go packages to WebAssembly modules.
build-wasm: $(HOME_WASM_OUTPUT) $(ABOUT_WASM_OUTPUT)

# Compile home WebAssembly module. Tracks changes in the home directory.
$(HOME_WASM_OUTPUT): $(wildcard $(CMD_WASM_HOME)/*.go)
	@echo "Compiling home WebAssembly module..."
	@mkdir -p $(dir $@)
	@$(TINYGO) build -target wasm -o $@ $(CMD_WASM_HOME)

# Compile about WebAssembly module. Tracks changes in the about directory.
$(ABOUT_WASM_OUTPUT): $(wildcard $(CMD_WASM_ABOUT)/*.go)
	@echo "Compiling about WebAssembly module..."
	@mkdir -p $(dir $@)
	@$(TINYGO) build -target wasm -o $@ $(CMD_WASM_ABOUT)

# clean removes all generated build artifacts.
clean:
	@echo "Cleaning generated build artifacts..."
	@rm -f $(APP_CSS_OUTPUT) $(WASM_EXEC_JS_OUTPUT) $(HOME_WASM_OUTPUT) $(ABOUT_WASM_OUTPUT)
