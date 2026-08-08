# Makefile for the Mega application.
#
# It automates the setup and compilation of the web assets and WebAssembly modules.

# Default target to execute when no arguments are provided.
.DEFAULT_GOAL := all

# The TinyGo executable used to compile Go code to WebAssembly.
#
# Fail fast if tinygo is missing.
TINYGO ?= tinygo
ifeq (, $(shell command -v $(TINYGO) 2> /dev/null))
	$(error Could not find tinygo. Is it installed correctly and in PATH?)
endif

# The Tailwind CSS CLI executable used to compile CSS styles.
#
# Fail fast if tailwindcss is missing.
TAILWINDCSS ?= tailwindcss
ifeq (, $(shell command -v $(TAILWINDCSS) 2> /dev/null))
	$(error Could not find tailwindcss. Is it installed correctly and in PATH?)
endif

# Build tool flags.
TINYGO_FLAGS   ?= -target=wasm -opt=s -no-debug
TAILWIND_FLAGS ?= --minify --content="./**/*.go" --content="./**/*.tmpl"

# Resolved path to the TinyGo installation directory.
#
# Fail fast if TINYGOROOT is not set.
TINYGOROOT := $(shell $(TINYGO) env TINYGOROOT 2>/dev/null)
ifeq ($(TINYGOROOT), )
	$(error Could not determine TINYGOROOT. Is tinygo configured correctly?)
endif

# Directories to exclude from source tracking.
IGNORE_DIRS := -type d -name .git -prune -o

# Tracks all Go files to trigger WASM rebuilds on internal package changes.
GO_SRCS := $(shell find . $(IGNORE_DIRS) -type f -name '*.go' -print)

# Tracks all UI files to trigger CSS rebuilds on Tailwind utility class changes.
UI_SRCS := $(shell find . $(IGNORE_DIRS) -type f \( -name '*.go' -o -name '*.tmpl' \) -print)

# Hierarchical directory structure for source code and compiled assets.
CMD                = cmd
CMD_WASM           = $(CMD)/wasm
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
WASM_EXEC_JS_OUTPUT = $(WEB_PUBLIC_JS_WASM)/wasm_exec.js

# List of WebAssembly modules to be compiled.
WASM_MODULES = \
	$(WEB_PUBLIC_WASM)/about.wasm \
	$(WEB_PUBLIC_WASM)/home.wasm

# Non-file action aliases (phony targets).
.PHONY: all setup build build-css build-wasm clean help

# all is the default target. It sets up the environment and compiles all application assets.
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

# Compile Tailwind CSS.
#
# Tracks UI files to catch utility class changes.
$(APP_CSS_OUTPUT): $(APP_CSS_INPUT) $(UI_SRCS)
	@echo "Compiling CSS..."
	@mkdir -p $(dir $@)
	@$(TAILWINDCSS) $(TAILWIND_FLAGS) -i $< -o $@

# build-wasm compiles the Go packages to WebAssembly modules.
build-wasm: $(WASM_MODULES)

# Compile WebAssembly modules.
#
# Tracks all Go files to catch internal package changes.
$(WASM_MODULES): $(WEB_PUBLIC_WASM)/%.wasm: $(GO_SRCS)
	@echo "Compiling $* WebAssembly module..."
	@mkdir -p $(dir $@)
	@$(TINYGO) build $(TINYGO_FLAGS) -o $@ ./$(CMD_WASM)/$*

# clean removes all generated build artifacts and output directories.
clean:
	@echo "Cleaning generated build artifacts..."
	@rm -rf $(WEB_PUBLIC_CSS) $(WEB_PUBLIC_JS_WASM) $(WEB_PUBLIC_WASM)

# help displays this help message.
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpDesc = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "  \033[36m%-15s\033[0m %s\n", helpCommand, helpDesc; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
