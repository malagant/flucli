# FluxCLI - Terminal UI for FluxCD Multi-Cluster Management

FluxCLI is a powerful terminal user interface (TUI) for managing FluxCD resources across multiple Kubernetes clusters. Inspired by tools like K9s, FluxCLI provides an intuitive, keyboard-driven interface specifically designed for GitOps workflows.

## Features

- **Multi-Cluster Support** - Seamlessly switch between and manage multiple Kubernetes clusters
- **FluxCD Resource Management** - View, monitor, and operate on GitRepository, HelmRepository, Kustomization, HelmRelease, and ResourceSet resources
- **Real-time Monitoring** - Live updates of resource status, events, and reconciliation progress
- **Intuitive Navigation** - K9s-inspired keyboard shortcuts and command patterns
- **Advanced Filtering** - Filter resources by namespace, status, cluster, and custom criteria
- **Event Streaming** - Monitor FluxCD events and reconciliation status in real-time

## Quick Start

### Installation

```bash
# Install from releases
curl -sL https://github.com/your-org/fluxcli/releases/latest/download/fluxcli-linux-amd64.tar.gz | tar xz
sudo mv fluxcli /usr/local/bin/

# Or build from source
git clone https://github.com/your-org/fluxcli.git
cd fluxcli
make build
```

### Basic Usage

```bash
# Launch FluxCLI
fluxcli

# Launch with specific kubeconfig
fluxcli --kubeconfig ~/.kube/config

# Launch with specific cluster context
fluxcli --context my-cluster
```

### Kubernetes Configuration

FluxCLI respects the standard Kubernetes configuration conventions:

#### Environment Variable Support

```bash
# Use KUBECONFIG environment variable (recommended)
export KUBECONFIG=/path/to/your/kubeconfig
fluxcli

# Multiple kubeconfig files
export KUBECONFIG=/path/to/config1:/path/to/config2
fluxcli
```

#### Configuration Priority

1. `--kubeconfig` command line flag (highest priority)
2. `KUBECONFIG` environment variable
3. `$HOME/.kube/config` (default fallback)

This follows the same pattern as `kubectl` and other Kubernetes tools.

### Key Navigation

| Key | Action |
|-----|--------|
| `?` | Show help |
| `:` | Enter command mode |
| `/` | Search/filter |
| `Tab` | Switch between panes |
| `Enter` | View resource details |
| `Ctrl+C` | Exit |

## Configuration

FluxCLI uses a configuration file located at `~/.fluxcli/config.yaml`:

```yaml
clusters:
  - name: production
    context: prod-cluster
    kubeconfig: ~/.kube/prod-config
  - name: staging  
    context: staging-cluster
    kubeconfig: ~/.kube/staging-config

defaults:
  namespace: flux-system
  refresh_interval: 5s
```

## Documentation

- [Architecture](architecture.md) - System design and components
- [User Guide](user-guide.md) - Comprehensive usage documentation
- [Multi-Cluster Support](multi-cluster-support.md) - Cluster management details
- [Feature Specifications](specs/) - Detailed feature documentation

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and contribution guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.