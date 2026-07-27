# Makefile for the Mega application.
#
# It automates the setup and compilation of the web assets and WebAssembly modules.

# Build tool executables.
TINYGO      ?= tinygo
TAILWINDCSS ?= tailwindcss

# Build tool flags.
TINYGO_FLAGS   ?= -target=wasm -opt=s -no-debug
TAILWIND_FLAGS ?= --minify

# Resolved path to the TinyGo installation directory.
TINYGOROOT := $(shell $(TINYGO) env TINYGOROOT)

# Directories to exclude from source tracking.
IGNORE_DIRS := -type d -name .git -prune -o

# Tracks all Go files to trigger WASM rebuilds on internal package changes.
GO_SRCS := $(shell find . $(IGNORE_DIRS) -type f -name '*.go' -print)

# Tracks all UI files to trigger CSS rebuilds on Tailwind utility class changes.
UI_SRCS := $(shell find . $(IGNORE_DIRS) -type f \( -name '*.go' -o -name '*.tmpl' \) -print)

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
.PHONY: all setup build build-css build-wasm clean help

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

# Compile Tailwind CSS.
#
# Tracks UI files to catch utility class changes.
$(APP_CSS_OUTPUT): $(APP_CSS_INPUT) $(UI_SRCS)
	@echo "Compiling CSS..."
	@mkdir -p $(dir $@)
	@$(TAILWINDCSS) $(TAILWIND_FLAGS) -i $< -o $@

# build-wasm compiles the Go packages to WebAssembly modules.
build-wasm: $(HOME_WASM_OUTPUT) $(ABOUT_WASM_OUTPUT)

# Compile home WebAssembly module.
#
# Tracks all Go files to catch internal package changes.
$(HOME_WASM_OUTPUT): $(GO_SRCS)
	@echo "Compiling home WebAssembly module..."
	@mkdir -p $(dir $@)
	@$(TINYGO) build $(TINYGO_FLAGS) -o $@ ./$(CMD_WASM_HOME)

# Compile about WebAssembly module.
#
# Tracks all Go files to catch internal package changes.
$(ABOUT_WASM_OUTPUT): $(GO_SRCS)
	@echo "Compiling about WebAssembly module..."
	@mkdir -p $(dir $@)
	@$(TINYGO) build $(TINYGO_FLAGS) -o $@ ./$(CMD_WASM_ABOUT)

# clean removes all generated build artifacts and empty directories.
clean:
	@echo "Cleaning generated build artifacts..."
	@rm -f $(APP_CSS_OUTPUT) $(WASM_EXEC_JS_OUTPUT) $(HOME_WASM_OUTPUT) $(ABOUT_WASM_OUTPUT)
	@[ -d $(WEB_PUBLIC) ] && find $(WEB_PUBLIC) -type d -empty -delete 2>/dev/null || true

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
