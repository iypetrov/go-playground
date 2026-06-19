NAME           := otelcol-sdnotify
REPO_ROOT      := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BIN_DIR        := $(REPO_ROOT)/bin
BUILD_DIR      := $(REPO_ROOT)/_build
TOOLS_MOD_FILE := $(REPO_ROOT)/internal/tools/go.mod
LD_FLAGS       ?= "-s -w"

.DEFAULT_GOAL := build

$(BIN_DIR):
	@mkdir -p $@

# Resolve the ocb (builder) tool from internal/tools/go.mod.
.PHONY: tools
tools:
	@cd $(REPO_ROOT)/internal/tools && go mod tidy

# Generate _build/ from manifest.yml. Pure codegen; no compilation.
# Runs from the repo root so manifest's `replaces:` paths resolve relative
# to it; the builder tool itself is pinned in internal/tools/go.mod.
.PHONY: generate-distribution
generate-distribution:
	@echo "Generating OpenTelemetry Collector distribution into $(BUILD_DIR)"
	@go tool -modfile $(TOOLS_MOD_FILE) builder \
		--skip-compilation \
		--config $(REPO_ROOT)/manifest.yml

# Build the collector binary into bin/.
.PHONY: build
build: $(BIN_DIR) generate-distribution
	@echo "Building $(NAME)"
	@cd $(BUILD_DIR) && \
		go mod tidy && \
		CGO_ENABLED=0 go build -ldflags $(LD_FLAGS) -o $(BIN_DIR)/$(NAME) .

# Run with the sample config. Watches ./input/ and prints to stdout.
.PHONY: run
run: build
	@mkdir -p $(REPO_ROOT)/input
	@$(BIN_DIR)/$(NAME) --config=$(REPO_ROOT)/config.yaml

.PHONY: test
test:
	@cd $(REPO_ROOT)/extension/sdnotify && go test -v ./...

.PHONY: clean
clean:
	@rm -rf $(BUILD_DIR) $(BIN_DIR)
