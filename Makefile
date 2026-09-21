# Makefile for the Mega application.
#
# Automates the setup and compilation of web assets and WebAssembly modules.

# Default target to execute when no target is specified on make command.
.DEFAULT_GOAL := all

# Remove partially-written targets if a recipe fails.
.DELETE_ON_ERROR:

# Verbosity control. Run `make V=1` to see the actual commands being executed.
V ?= 0
ifeq ($(V),1)
Q :=
else
Q := @
endif

# Runs the build unless every requested goal is metadata-only: clean, help, or rebuild.
NEED_BUILD := $(filter-out clean help rebuild,$(or $(MAKECMDGOALS),all))

# Runs tool checks unless every requested goal is metadata-only: clean or help.
NEED_TOOLS := $(filter-out clean help,$(or $(MAKECMDGOALS),all))

# --- Directory Structure ---

CMD                := cmd
CMD_WASM           := $(CMD)/wasm
VENDOR             := vendor
WEB                := web
WEB_PUBLIC         := $(WEB)/public
WEB_PUBLIC_CSS     := $(WEB_PUBLIC)/css
WEB_PUBLIC_JS      := $(WEB_PUBLIC)/js
WEB_PUBLIC_JS_WASM := $(WEB_PUBLIC_JS)/wasm
WEB_PUBLIC_WASM    := $(WEB_PUBLIC)/wasm
WEB_SRC            := $(WEB)/src
WEB_SRC_CSS        := $(WEB_SRC)/css

# --- Tooling Setup ---

ifneq ($(NEED_TOOLS),)

# Fail fast if tinygo is missing.
TINYGO ?= tinygo
ifeq (, $(shell command -v $(TINYGO) 2> /dev/null))
$(error Could not find tinygo. Is it installed correctly and in PATH?)
endif

# Fail fast if TINYGOROOT is not set.
TINYGOROOT := $(shell $(TINYGO) env TINYGOROOT 2>/dev/null)
ifeq ($(TINYGOROOT),)
$(error Could not determine TINYGOROOT. Is tinygo configured correctly?)
endif

# Fail fast if tailwindcss is missing.
TAILWINDCSS ?= tailwindcss
ifeq (, $(shell command -v $(TAILWINDCSS) 2> /dev/null))
$(error Could not find tailwindcss. Is it installed correctly and in PATH?)
endif

endif

# --- Configuration & Flags ---

# Size-optimized WASM build, debug info stripped.
TINYGO_FLAGS ?= -target=wasm -opt=s -panic=trap -no-debug

# Minifies output CSS.
TAILWIND_FLAGS ?= --minify

# --- Source Tracking ---

ifneq ($(NEED_BUILD),)

# Directories to exclude from source tracking.
IGNORE_DIRS := -type d \( -name .git -o -name $(VENDOR) -o -path ./$(WEB_PUBLIC) \) -prune -o

# Tracks all Go files to trigger WASM rebuilds on internal package changes.
GO_SRCS := $(shell find . $(IGNORE_DIRS) -type f -name '*.go' ! -name '*_test.go' -print)

# Tracks all UI files to trigger CSS rebuilds on Tailwind utility class changes.
UI_SRCS := $(shell find . $(IGNORE_DIRS) -type f \( -name '*.go' ! -name '*_test.go' -o -name '*.tmpl' \) -print)

endif

# --- Inputs and Outputs ---

APP_CSS_INPUT      := $(WEB_SRC_CSS)/app.css
WASM_EXEC_JS_INPUT := $(TINYGOROOT)/targets/wasm_exec.js

APP_CSS_OUTPUT      := $(WEB_PUBLIC_CSS)/app.css
WASM_EXEC_JS_OUTPUT := $(WEB_PUBLIC_JS_WASM)/wasm_exec.js

# WebAssembly modules to be compiled.
WASM_MODULES := $(patsubst $(CMD_WASM)/%/main.go,$(WEB_PUBLIC_WASM)/%.wasm,$(wildcard $(CMD_WASM)/*/main.go))

# Non-file action aliases.
.PHONY: all setup build build-css build-wasm clean rebuild help

# --- Targets ---

# Default target: sets up the environment and compiles all assets.
all: setup build

# Executes all compilation targets for the application.
build: build-css build-wasm

# Prepares dependencies for the application.
setup: $(WASM_EXEC_JS_OUTPUT)

# Copies the WebAssembly execution script from TinyGo if it doesn't exist or is updated.
$(WASM_EXEC_JS_OUTPUT): $(WASM_EXEC_JS_INPUT)
	@echo "Setting up wasm_exec.js..."
	$(Q)mkdir -p $(@D)
	$(Q)cp $< $@

# Compiles the styles to output a CSS file.
build-css: $(APP_CSS_OUTPUT)

# Compiles Tailwind CSS. Tracks UI files to catch utility class changes.
$(APP_CSS_OUTPUT): $(APP_CSS_INPUT) $(UI_SRCS)
	@echo "Compiling CSS..."
	$(Q)mkdir -p $(@D)
	$(Q)$(TAILWINDCSS) $(TAILWIND_FLAGS) -i $< -o $@

# Compiles the Go packages to WebAssembly modules.
build-wasm: $(WASM_MODULES)

# Compiles WebAssembly modules. Tracks all Go files to catch internal package changes.
$(WASM_MODULES): $(WEB_PUBLIC_WASM)/%.wasm: $(GO_SRCS) go.mod $(wildcard go.sum)
	@echo "Compiling $* WebAssembly module..."
	$(Q)mkdir -p $(@D)
	$(Q)$(TINYGO) build $(TINYGO_FLAGS) -o $@ ./$(CMD_WASM)/$*

# Removes all generated build artifacts and output directories.
clean:
	@echo "Cleaning generated build artifacts..."
	$(Q)rm -rf $(WEB_PUBLIC_CSS) $(WEB_PUBLIC_JS_WASM) $(WEB_PUBLIC_WASM)
	$(Q)find $(WEB_PUBLIC) -type d -empty -delete 2>/dev/null || true

# Cleans the project and builds it from scratch.
rebuild:
	@echo "Rebuilding project..."
	$(Q)$(MAKE) clean
	$(Q)$(MAKE) build

# Displays this help message.
help:
	@echo "Usage: make [target] [V=1 (for verbose output)]"
	@echo ""
	@echo "Targets:"
	@awk '/^[a-zA-Z0-9_-]+:/ { \
		helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 1, index($$1, ":")-1); \
			helpDesc = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "  \033[36m%-15s\033[0m %s\n", helpCommand, helpDesc; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
