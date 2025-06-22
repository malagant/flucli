# Nix Flake Quick Reference

This is a quick reference for developers using the FluxCLI Nix flake. For comprehensive documentation, see [Development Environment](development-environment.md).

## Basic Commands

```bash
# Enter development shell
nix develop

# Build FluxCLI package
nix build

# Run FluxCLI directly
nix run

# Update flake inputs (dependencies)
nix flake update

# Check flake structure
nix flake show
```

## Development Workflow

```bash
# 1. Enter development environment
nix develop

# 2. Build project
./dev.sh build

# 3. Run tests
./dev.sh test

# 4. Run linter
./dev.sh lint
```

## Available Tools

| Category | Tool | Purpose |
|----------|------|---------|
| **Go Development** | go, gopls, gotools | Language runtime & tooling |
| | delve, golangci-lint | Debugging & linting |
| **Kubernetes** | kubectl, kubecolor | Kubernetes CLI with colors |
| | kind, helm, fluxcd | Local clusters & GitOps |
| **Shell Enhancement** | starship | Modern shell prompt |
| | zsh-autosuggestions | Fish-like suggestions |
| | zsh-syntax-highlighting | Command highlighting |
| | fzf, atuin | Fuzzy search & history |
| **Development** | git, make, tmux | Version control & build |
| | eza, bat, ripgrep, fd | Enhanced CLI tools |
| | jq, yq-go | JSON/YAML processing |

## FluxCLI Aliases & Functions

### Quick Commands
```bash
fk              # ./fluxcli
fks             # ./fluxcli --all-namespaces  
fkr             # flux reconcile source git
fluxdev         # Build and run FluxCLI
fluxtest        # Run tests and linting
```

### Kubernetes & FluxCD
```bash
kgf             # Get all FluxCD resources
kdf             # kubectl describe
fluxcli-test    # Test against different clusters
flux-ctx        # Switch kubectl context
```

## Environment Variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `FLUXCLI_DEV` | `true` | Development mode flag |
| `FLUXCLI_LOG_LEVEL` | `debug` | Default log level |
| `GOPATH` | `$HOME/go` | Go workspace (if not set) |

## Direnv Integration

```bash
# Install direnv
nix-env -iA nixpkgs.direnv

# Enable for this project
echo "use flake" > .envrc
direnv allow

# Now the environment activates automatically when entering the directory
```

## Common Issues

### Slow Initial Setup
```bash
# Use binary cache for faster downloads
nix develop --option substituters "https://cache.nixos.org"
```

### Clear Cache
```bash
# Clean up disk space
nix-collect-garbage -d
```

### Shell Integration Issues
```bash
# Verify you're in nix shell
echo $IN_NIX_SHELL  # Should output "impure"

# Check available tools
which go kubectl flux
```

## Cross-Platform Notes

### macOS
- All tools work natively on Apple Silicon and Intel
- Shell integration respects existing `.zshrc`

### Linux
- Includes additional Linux-specific tools
- Works with any shell (bash, zsh, fish)

## See Also

- [Development Environment](development-environment.md) - Complete setup guide
- [Nix Manual](https://nixos.org/manual/nix/stable/) - Official Nix documentation
- [Flakes Wiki](https://nixos.wiki/wiki/Flakes) - Community flakes documentation
