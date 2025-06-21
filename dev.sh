#!/usr/bin/env bash

# FluxCLI Development Helper Script
# Provides an easy way to run commands in the Nix development environment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Print usage
usage() {
    echo "FluxCLI Development Helper"
    echo ""
    echo "Usage: $0 [command] [args...]"
    echo ""
    echo "Commands:"
    echo "  build          Build the FluxCLI binary"
    echo "  test           Run tests"
    echo "  lint           Run linter"
    echo "  tidy           Tidy Go modules"
    echo "  clean          Clean built artifacts"
    echo "  run [args...]  Run FluxCLI with arguments"
    echo "  shell          Enter Nix development shell"
    echo ""
    echo "Examples:"
    echo "  $0 build                     # Build the binary"
    echo "  $0 run --help               # Run FluxCLI with --help"
    echo "  $0 test                     # Run tests"
    echo "  $0 shell                    # Enter development shell"
}

# Run command with Nix
run_with_nix() {
    echo -e "${GREEN}Running:${NC} $*"
    nix shell 'nixpkgs#go' -c "$@"
}

# Main command processing
case "${1:-}" in
    "build")
        echo -e "${YELLOW}Building FluxCLI...${NC}"
        run_with_nix go build -o fluxcli
        echo -e "${GREEN}✅ Build completed successfully!${NC}"
        echo -e "Binary created: ${GREEN}./fluxcli${NC}"
        ;;
    
    "test")
        echo -e "${YELLOW}Running tests...${NC}"
        run_with_nix go test ./...
        ;;
    
    "lint")
        echo -e "${YELLOW}Running linter...${NC}"
        run_with_nix go vet ./...
        run_with_nix go fmt ./...
        echo -e "${GREEN}✅ Linting completed!${NC}"
        ;;
    
    "tidy")
        echo -e "${YELLOW}Tidying Go modules...${NC}"
        run_with_nix go mod tidy
        echo -e "${GREEN}✅ Modules tidied!${NC}"
        ;;
    
    "clean")
        echo -e "${YELLOW}Cleaning artifacts...${NC}"
        rm -f fluxcli
        echo -e "${GREEN}✅ Cleaned!${NC}"
        ;;
    
    "run")
        shift
        if [ ! -f "./fluxcli" ]; then
            echo -e "${YELLOW}Binary not found, building first...${NC}"
            run_with_nix go build -o fluxcli
        fi
        echo -e "${GREEN}Running FluxCLI:${NC} ./fluxcli $*"
        ./fluxcli "$@"
        ;;
    
    "shell")
        echo -e "${YELLOW}Entering Nix development shell...${NC}"
        echo -e "${GREEN}💡 You'll have access to Go and other development tools${NC}"
        nix shell 'nixpkgs#go' 'nixpkgs#kubectl' 'nixpkgs#kubernetes-helm'
        ;;
    
    "help"|"--help"|"-h"|"")
        usage
        ;;
    
    *)
        echo -e "${RED}Unknown command: $1${NC}"
        echo ""
        usage
        exit 1
        ;;
esac
