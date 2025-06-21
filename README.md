# FluxCLI

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey.svg)]()

FluxCLI is a powerful **Terminal User Interface (TUI)** for managing FluxCD resources across multiple Kubernetes clusters. Inspired by tools like K9s, FluxCLI provides an intuitive, keyboard-driven interface specifically designed for GitOps workflows.

## ✨ Features

- 🌐 **Multi-Cluster Support** - Seamlessly switch between and manage multiple Kubernetes clusters
- 🔄 **FluxCD Resource Management** - View, monitor, and operate on GitRepository, HelmRepository, Kustomization, HelmRelease, and ResourceSet resources
- ⚡ **Real-time Monitoring** - Live updates of resource status, events, and reconciliation progress
- ⌨️ **Intuitive Navigation** - K9s-inspired keyboard shortcuts and command patterns
- 🔍 **Advanced Filtering** - Filter resources by namespace, status, cluster, and custom criteria
- 📡 **Event Streaming** - Monitor FluxCD events and reconciliation status in real-time
- 🎨 **Beautiful Interface** - Clean, colorized terminal interface with responsive layouts

## 🚀 Quick Start

### Prerequisites

- Kubernetes cluster(s) with FluxCD v2.0+ installed
- kubectl configured with access to your clusters
- Go 1.21+ (for building from source)

### Installation

#### Using Nix (Recommended for Development)

```bash
# Clone the repository
git clone https://github.com/malagant/fluxcli.git
cd fluxcli

# Build using our development tools
./dev.sh build

# Run FluxCLI
./fluxcli
```

#### Using Make

```bash
# Clone and build
git clone https://github.com/malagant/fluxcli.git
cd fluxcli
make build

# Run FluxCLI
./fluxcli
```

#### Using Go

```bash
# Install directly from source
go install github.com/malagant/fluxcli@latest

# Or build locally
git clone https://github.com/malagant/fluxcli.git
cd fluxcli
go build -o fluxcli .
```

### First Run

1. **Configure FluxCLI**: Create a configuration file at `~/.fluxcli/config.yaml` or let FluxCLI create one for you
2. **Launch**: Run `fluxcli` to start the terminal interface
3. **Navigate**: Use keyboard shortcuts to browse your FluxCD resources

## 🎮 Usage

### Basic Navigation

| Key | Action |
|-----|--------|
| `j/k` | Move up/down in lists |
| `g/G` | Go to top/bottom |
| `Enter` | View resource details |
| `Tab` | Switch between views |
| `Ctrl+K/J` | Switch clusters |
| `1-4` | Switch resource types |
| `:` | Enter command mode |
| `?` | Toggle help |
| `q` | Quit |

### Command Mode

Press `:` to enter command mode for advanced operations:

- `:suspend <resource>` - Suspend a FluxCD resource
- `:resume <resource>` - Resume a FluxCD resource  
- `:reconcile <resource>` - Trigger reconciliation
- `:quit` - Exit FluxCLI

### Configuration

FluxCLI uses a YAML configuration file located at `~/.fluxcli/config.yaml`:

```yaml
# Multi-cluster configuration
clusters:
  - name: "production"
    kubeconfig: "~/.kube/config"
    context: "prod-cluster"
  - name: "staging"
    kubeconfig: "~/.kube/staging-config"
    context: "staging-cluster"

# Default settings
defaults:
  namespace: "flux-system"
  refresh_interval: "5s"
  max_concurrent_clusters: 10

# UI preferences
ui:
  theme: "dark"
  show_events: true
  columns:
    - "Name"
    - "Namespace" 
    - "Age"
    - "Status"
    - "Message"
```

## 🛠️ Development

FluxCLI provides several development tools and scripts:

### Development Environment

We use Nix for reproducible development environments:

```bash
# Enter development shell (provides Go, kubectl, helm, etc.)
./dev.sh shell

# Or use individual commands
./dev.sh build    # Build the binary
./dev.sh test     # Run tests
./dev.sh lint     # Run linter
./dev.sh tidy     # Tidy Go modules
```

### Using Make

```bash
make help         # Show all available commands
make build        # Build FluxCLI
make test         # Run tests
make lint         # Run linter
make run          # Build and run
make clean        # Clean artifacts
```

### Project Structure

```
fluxcli/
├── cmd/                 # CLI commands and entry points
├── pkg/
│   ├── core/           # Core application logic
│   ├── k8s/            # Kubernetes client and resource management
│   └── ui/             # Terminal UI components (Bubble Tea)
├── internal/
│   └── config/         # Configuration management
├── docs/               # Documentation
├── dev.sh              # Development helper script
├── Makefile            # Build automation
├── flake.nix           # Nix development environment
└── go.mod              # Go module definition
```

## 📚 Documentation

- [User Guide](docs/user-guide.md) - Comprehensive usage guide
- [Architecture](docs/architecture.md) - Technical architecture overview
- [Multi-Cluster Support](docs/multi-cluster-support.md) - Multi-cluster configuration
- [MVP Features](docs/specs/mvp-features.md) - MVP feature specifications

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Run linter (`make lint`)
6. Commit your changes (`git commit -m 'Add some amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## 📋 Requirements

- **Kubernetes**: 1.20+
- **FluxCD**: v2.0+
- **Go**: 1.21+ (for building)
- **Terminal**: Modern terminal with color support

## 🐛 Troubleshooting

### Common Issues

**FluxCLI won't connect to cluster:**
- Verify kubectl can connect: `kubectl cluster-info`
- Check kubeconfig path and context in configuration
- Ensure FluxCD is installed: `flux check`

**Resources not displaying:**
- Verify FluxCD resources exist: `kubectl get gitrepositories -A`
- Check namespace permissions
- Ensure correct FluxCD CRDs are installed

**Performance issues:**
- Reduce refresh interval in configuration
- Limit concurrent clusters in config
- Filter resources by namespace

For more troubleshooting help, see our [User Guide](docs/user-guide.md).

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [FluxCD](https://fluxcd.io/) - GitOps toolkit
- [K9s](https://k9scli.io/) - Inspiration for keyboard-driven K8s UI
- [Kubernetes](https://kubernetes.io/) - Container orchestration platform

## 🌟 Star History

[![Star History Chart](https://api.star-history.com/svg?repos=malagant/fluxcli&type=Date)](https://star-history.com/#malagant/fluxcli&Date)

---

**FluxCLI** - *Making GitOps management delightful* ⭐
