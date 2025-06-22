# Development Environment with Nix Flake

FluxCLI uses [Nix flakes](https://nixos.wiki/wiki/Flakes) to provide a reproducible, isolated development environment. This ensures all developers have access to the same versions of tools and dependencies, regardless of their host system.

## Overview

The Nix flake provides:
- **Go toolchain** with consistent versioning
- **Kubernetes tools** (kubectl, kind, helm, flux)
- **Development utilities** (linting, testing, documentation)
- **Cross-platform support** (Linux, macOS)
- **Shell integration** that respects your existing configuration

## Prerequisites

### Install Nix

**macOS/Linux:**
```bash
# Install Nix with flakes support
curl --proto '=https' --tlsv1.2 -sSf -L https://install.determinate.systems/nix | sh -s -- install
```

**Enable flakes** (if using older Nix):
```bash
mkdir -p ~/.config/nix
echo 'experimental-features = nix-command flakes' >> ~/.config/nix/nix.conf
```

## Quick Start

### 1. Enter Development Environment

```bash
cd /path/to/fluxcli
nix develop
```

This will:
- Download and cache all required tools
- Set up environment variables
- Source your existing `.zshrc` (if using zsh)
- Display available development commands

### 2. Build FluxCLI

```bash
# Using the dev script (recommended)
./dev.sh build

# Or directly with Go
go build -o fluxcli
```

### 3. Run Tests

```bash
./dev.sh test
# or
go test ./...
```

## Available Tools

When you enter the development environment, you'll have access to:

### Core Development
- **Go** - Language runtime and toolchain
- **gopls** - Go language server
- **gotools** - Additional Go utilities
- **go-tools** - Static analysis tools
- **delve** - Go debugger
- **golangci-lint** - Comprehensive Go linter

### Kubernetes Ecosystem
- **kubectl** - Kubernetes CLI
- **kubecolor** - Colorized kubectl output
- **kind** - Kubernetes in Docker
- **kubernetes-helm** - Helm package manager
- **fluxcd** - Flux CLI for GitOps

### Shell Enhancements
- **starship** - Modern, fast shell prompt
- **zsh-autosuggestions** - Fish-like autosuggestions
- **zsh-syntax-highlighting** - Command syntax highlighting
- **zsh-completions** - Additional tab completions
- **fzf** - Fuzzy finder for files and history
- **atuin** - Shell history management

### Development Utilities
- **git** - Version control
- **gnumake** - Build automation
- **tmux** - Terminal multiplexer
- **curl** - HTTP client
- **jq** - JSON processor
- **yq-go** - YAML processor
- **mdbook** - Documentation generator
- **eza** - Modern ls replacement with colors
- **bat** - Better cat with syntax highlighting
- **ripgrep** - Fast grep replacement
- **fd** - Better find alternative
- **tree** - Directory tree display
- **direnv** - Environment management

## Environment Variables

The flake sets up several environment variables for development:

```bash
# Core Go configuration
GOPATH=$HOME/go                    # Go workspace (if not already set)
PATH=$GOPATH/bin:$PATH            # Go binaries in PATH

# FluxCLI specific
FLUXCLI_DEV=true                  # Development mode flag
FLUXCLI_LOG_LEVEL=debug           # Default log level (customizable)

# Session management
FLUXCLI_WELCOME_SHOWN=1           # Prevents repeated welcome messages
FLUXCLI_ZSHRC_SOURCED=1          # Tracks .zshrc sourcing

# Shell enhancements
ZSH_AUTOSUGGEST_STRATEGY=(history completion)  # Autosuggestion strategy
ZSH_AUTOSUGGEST_BUFFER_MAX_SIZE=20             # Max suggestion buffer
```

## Shell Integration

### Respecting Your Configuration

The flake is designed to **enhance** rather than **replace** your existing shell setup:

- **Preserves** your existing `GOPATH` if already set
- **Maintains** your `PATH` order by only prepending if necessary
- **Sources** your `~/.zshrc` automatically in zsh shells
- **Honors** existing environment variables like `FLUXCLI_LOG_LEVEL`

### Welcome Message

On first entry, you'll see:
```
🚀 FluxCLI development environment loaded
📦 Go version: go version go1.24.0 darwin/arm64
🖥️  Platform: aarch64-darwin
☸️  kubectl version: Client Version: v1.31.0
🌊 Flux version: flux version 2.4.0

Available commands:
  ./dev.sh build    - Build FluxCLI
  ./dev.sh test     - Run tests
  ./dev.sh lint     - Run linter
  make help         - Show all available make targets
```

## Development Workflow

### 1. Daily Development

```bash
# Enter the environment
nix develop

# Build and test in one command
./dev.sh build && ./dev.sh test

# Run linting
./dev.sh lint

# Start development server/watching
make watch
```

### 2. IDE Integration

The environment works with any editor/IDE:

**VS Code:**
```bash
# Start VS Code from within the nix shell
nix develop
code .
```

**Vim/Neovim:**
```bash
# LSP and tools will be available automatically
nix develop
nvim .
```

### 3. Kubernetes Development

```bash
# Create a local test cluster
kind create cluster --name fluxcli-dev

# Install Flux
flux install

# Test FluxCLI against the cluster
./fluxcli
```

## Advanced Usage

### Custom Environment Variables

Override defaults by setting variables before entering:

```bash
export FLUXCLI_LOG_LEVEL=trace
export GOPATH=/custom/path
nix develop
```

### Direnv Integration

For automatic environment activation:

```bash
# Install direnv (if not using the flake)
nix-env -iA nixpkgs.direnv

# Create .envrc
echo "use flake" > .envrc
direnv allow
```

### Building the Package

Build FluxCLI as a Nix package:

```bash
# Build the package
nix build

# Run the built package
./result/bin/fluxcli
```

## FluxCLI-Specific Enhancements

The flake provides several FluxCLI-specific aliases and functions that are only added if not already defined by the user:

### Aliases

```bash
# FluxCLI shortcuts
fk                    # Quick alias for ./fluxcli
fks                   # FluxCLI with --all-namespaces
fkr                   # flux reconcile source git

# Development helpers
fluxdev               # Build and run FluxCLI in one command
fluxtest              # Run tests and linting

# Enhanced kubectl for FluxCD
kgf                   # Get all FluxCD resources across namespaces
kdf                   # kubectl describe (if not already defined)
```

### Functions

```bash
# Test FluxCLI against different clusters
fluxcli-test [context]    # Build and test against specified context

# Quick context switching
flux-ctx <context>        # Switch kubectl context and show cluster info
```

### Example Usage

```bash
# Quick development workflow
fluxdev                   # Build and run FluxCLI

# Test against different clusters
fluxcli-test staging     # Test against staging cluster
fluxcli-test production  # Test against production cluster

# View all FluxCD resources
kgf                      # Show GitRepositories, HelmRepositories, etc.
```

## Troubleshooting

### Common Issues

**Go not found:**
```bash
# Verify you're in the nix shell
echo $IN_NIX_SHELL  # Should output "impure"

# Re-enter the environment
exit
nix develop
```

**PATH issues:**
```bash
# Check if Go is in PATH
which go
echo $PATH | tr ':' '\n' | grep go
```

**Performance:**
```bash
# Clear Nix cache if builds are slow
nix-collect-garbage -d

# Use binary cache for faster downloads
nix develop --option substituters "https://cache.nixos.org"
```

### Getting Help

- **Nix Manual**: https://nixos.org/manual/nix/stable/
- **Flakes Wiki**: https://nixos.wiki/wiki/Flakes
- **FluxCLI Issues**: https://github.com/malagant/fluxcli/issues

## Contributing

### Adding New Tools

To add a new development tool:

1. **Edit** `flake.nix`:
```nix
buildInputs = with pkgs; [
  # ... existing tools
  your-new-tool
];
```

2. **Test** the change:
```bash
nix develop
which your-new-tool
```

3. **Update** this documentation

### Platform-Specific Tools

Add platform-specific tools using conditionals:

```nix
] ++ pkgs.lib.optionals isLinux [
  # Linux-specific packages
  your-linux-tool
] ++ pkgs.lib.optionals isDarwin [
  # macOS-specific packages
  your-macos-tool
];
```

## Alternative: Traditional Setup

If you prefer not to use Nix, you can set up the development environment manually:

```bash
# Install Go 1.24+
# Install kubectl, kind, helm, flux
# Install development tools
# See the buildInputs section in flake.nix for the complete list
```

However, using Nix ensures consistency across all development environments and eliminates "works on my machine" issues.
