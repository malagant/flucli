# FluxCLI Makefile
# Provides build, test, and development commands

.PHONY: help build test lint clean install run dev tidy format check deps

# Default target
.DEFAULT_GOAL := help

# Binary name
BINARY_NAME := fluxcli
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION := $(shell go version | cut -d ' ' -f 3)

# Build flags
LDFLAGS := -X main.Version=$(VERSION) \
          -X main.BuildTime=$(BUILD_TIME) \
          -X main.GoVersion=$(GO_VERSION)

# Go commands (with Nix shell wrapper)
GO_CMD := nix shell nixpkgs\#go -c go
GO_BUILD := $(GO_CMD) build
GO_TEST := $(GO_CMD) test
GO_VET := $(GO_CMD) vet
GO_FMT := $(GO_CMD) fmt
GO_MOD := $(GO_CMD) mod

help: ## Show this help message
	@echo 'FluxCLI Development Commands:'
	@echo ''
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
dev: ## Start development environment
	@echo "🚀 Starting FluxCLI development environment..."
	@echo "💡 Use 'make build' to build, 'make test' to test"
	@echo "📝 Use './dev.sh shell' to enter Nix development shell"

deps: ## Download Go module dependencies
	@echo "📦 Downloading dependencies..."
	$(GO_MOD) download
	$(GO_MOD) tidy

##@ Build
build: deps ## Build the FluxCLI binary
	@echo "🔨 Building $(BINARY_NAME)..."
	$(GO_BUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .
	@echo "✅ Build completed: ./$(BINARY_NAME)"

build-all: deps ## Build for all platforms
	@echo "🔨 Building for all platforms..."
	GOOS=linux GOARCH=amd64 $(GO_BUILD) -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 $(GO_BUILD) -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 $(GO_BUILD) -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GO_BUILD) -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 $(GO_BUILD) -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)-windows-amd64.exe .
	@echo "✅ Cross-platform builds completed in ./dist/"

install: build ## Install the binary to $GOPATH/bin
	@echo "📦 Installing $(BINARY_NAME)..."
	$(GO_CMD) install -ldflags "$(LDFLAGS)" .
	@echo "✅ Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

##@ Testing
test: deps ## Run tests
	@echo "🧪 Running tests..."
	$(GO_TEST) -v ./...

test-coverage: deps ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	$(GO_TEST) -v -coverprofile=coverage.out ./...
	$(GO_CMD) tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

benchmark: deps ## Run benchmarks
	@echo "⚡ Running benchmarks..."
	$(GO_TEST) -bench=. -benchmem ./...

##@ Code Quality
lint: deps ## Run linter
	@echo "🔍 Running linter..."
	$(GO_VET) ./...
	@echo "✅ Linting completed"

format: ## Format Go code
	@echo "🎨 Formatting code..."
	$(GO_FMT) ./...
	@echo "✅ Code formatted"

check: lint test ## Run all checks (lint + test)
	@echo "✅ All checks passed"

tidy: ## Tidy Go modules
	@echo "🧹 Tidying modules..."
	$(GO_MOD) tidy
	@echo "✅ Modules tidied"

##@ Running
run: build ## Build and run FluxCLI
	@echo "🚀 Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

run-help: build ## Build and run FluxCLI with --help
	@echo "📖 Running $(BINARY_NAME) --help..."
	./$(BINARY_NAME) --help

run-debug: build ## Build and run FluxCLI in debug mode
	@echo "🔍 Running $(BINARY_NAME) in debug mode..."
	./$(BINARY_NAME) --debug --log-level=debug

##@ Cleanup
clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	rm -f $(BINARY_NAME)
	rm -rf dist/
	rm -f coverage.out coverage.html
	@echo "✅ Cleaned"

##@ Info
version: ## Show version information
	@echo "FluxCLI Version Info:"
	@echo "  Version: $(VERSION)"
	@echo "  Build Time: $(BUILD_TIME)"
	@echo "  Go Version: $(GO_VERSION)"

info: ## Show project information
	@echo "FluxCLI Project Information:"
	@echo "  Binary: $(BINARY_NAME)"
	@echo "  Version: $(VERSION)"
	@echo "  Go Version: $(GO_VERSION)"
	@echo ""
	@echo "Available commands:"
	@echo "  make build    - Build the binary"
	@echo "  make test     - Run tests"
	@echo "  make lint     - Run linter"
	@echo "  make run      - Build and run"
	@echo "  make clean    - Clean artifacts"
	@echo ""
	@echo "Development tools:"
	@echo "  ./dev.sh      - Development helper script"
	@echo "  make help     - Show all available commands"
