# Multi-Cluster Support

FluxCLI provides comprehensive support for managing FluxCD resources across multiple Kubernetes clusters simultaneously.

## Configuration

### Cluster Definition

Define clusters in `~/.fluxcli/config.yaml`:

```yaml
clusters:
  - name: production
    context: prod-cluster
    kubeconfig: ~/.kube/prod-config
    namespace: flux-system
    color: red
    description: "Production environment"
    
  - name: staging
    context: staging-cluster
    kubeconfig: ~/.kube/staging-config
    namespace: flux-system
    color: yellow
    description: "Staging environment"
    
  - name: development
    context: dev-cluster
    kubeconfig: ~/.kube/dev-config
    namespace: flux-system
    color: green
    description: "Development environment"

defaults:
  namespace: flux-system
  refresh_interval: 5s
  max_concurrent_clusters: 10
```

### Environment Variables

Override cluster settings with environment variables:

```bash
export FLUXCLI_KUBECONFIG=/path/to/kubeconfig
export FLUXCLI_CONTEXT=my-cluster
export FLUXCLI_NAMESPACE=custom-namespace
```

## Cluster Management

### Adding Clusters

#### Interactive Mode
```bash
fluxcli cluster add
# Prompts for cluster details
```

#### Command Line
```bash
fluxcli cluster add --name production --context prod --kubeconfig ~/.kube/prod
```

### Removing Clusters
```bash
fluxcli cluster remove production
```

### Listing Clusters
```bash
fluxcli cluster list
```

## Navigation and Context Switching

### Cluster Switching in TUI

| Key Combination | Action |
|----------------|--------|
| `Ctrl+K` | Switch to next cluster |
| `Ctrl+J` | Switch to previous cluster |
| `:cluster <name>` | Switch to specific cluster |
| `:clusters` | Show cluster selector |

### Cluster Indicator

The current cluster is always visible in the UI:

```
┌─ FluxCLI - [PRODUCTION] ──────────────────────────────────┐
│ Namespace: flux-system                    Connected: ✓     │
├────────────────────────────────────────────────────────────┤
│ GitRepositories                                            │
│ NAME              READY   STATUS      AGE                  │
│ flux-system       True    Stored      2d                   │
└────────────────────────────────────────────────────────────┘
```

## Multi-Cluster Views

### Unified Resource View

View resources across all clusters simultaneously:

```bash
:all-clusters
```

Display format:
```
┌─ All Clusters - GitRepositories ──────────────────────────┐
│ CLUSTER      NAME           READY   STATUS     AGE        │
│ production   flux-system    True    Stored     2d         │
│ staging      flux-system    True    Stored     1d         │
│ development  flux-system    False   Failed     3h         │
└────────────────────────────────────────────────────────────┘
```

### Cluster Health Dashboard

Overview of all cluster statuses:

```bash
:health
```

### Cross-Cluster Operations

Perform operations across multiple clusters:

```bash
:reconcile --all-clusters gitrepository/flux-system
:suspend --clusters prod,staging helmrelease/my-app
```

## Connection Management

### Connection Pooling

FluxCLI maintains persistent connections to configured clusters:

- **Connection Pool Size**: Configurable per cluster
- **Keep-Alive**: Automatic connection health monitoring
- **Retry Logic**: Exponential backoff on connection failures
- **Timeout Handling**: Configurable timeouts for operations

### Connection Status Indicators

| Status | Indicator | Description |
|--------|-----------|-------------|
| Connected | `✓` | Healthy connection |
| Connecting | `⟳` | Establishing connection |
| Disconnected | `✗` | Connection failed |
| Degraded | `⚠` | Partial connectivity |

### Offline Mode

When clusters become unavailable:
- Cached resource data remains visible
- Clear indication of stale data
- Automatic reconnection attempts
- Graceful degradation of functionality

## Security and Authentication

### Authentication Methods

Supports all kubectl-compatible authentication:
- **Certificates**: Client certificate authentication
- **Tokens**: Bearer token authentication
- **OIDC**: OpenID Connect integration
- **AWS IAM**: AWS IAM Authenticator
- **GCP**: Google Cloud authentication
- **Azure**: Azure Active Directory

### Per-Cluster Authentication

Each cluster can use different authentication methods:

```yaml
clusters:
  - name: aws-cluster
    context: aws-prod
    kubeconfig: ~/.kube/aws-config
    auth_provider: aws
    
  - name: gcp-cluster
    context: gcp-prod
    kubeconfig: ~/.kube/gcp-config
    auth_provider: gcp
```

### RBAC Considerations

FluxCLI operations require appropriate RBAC permissions:

```yaml
# Minimum required permissions
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: fluxcli-user
rules:
- apiGroups: ["source.toolkit.fluxcd.io"]
  resources: ["*"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["kustomize.toolkit.fluxcd.io"]
  resources: ["*"]
  verbs: ["get", "list", "watch", "patch"]
- apiGroups: ["helm.toolkit.fluxcd.io"]
  resources: ["*"]
  verbs: ["get", "list", "watch", "patch"]
- apiGroups: [""]
  resources: ["events"]
  verbs: ["get", "list", "watch"]
```

## Performance Optimization

### Parallel Operations

- Concurrent resource fetching across clusters
- Independent refresh cycles per cluster
- Asynchronous event processing

### Resource Caching

- Intelligent caching with TTL
- Delta updates for changed resources
- Memory-efficient storage

### Network Optimization

- Connection reuse across operations
- Batch API calls where possible
- Configurable request timeouts

## Troubleshooting

### Common Issues

#### Cluster Connection Failures
```bash
# Check cluster connectivity
fluxcli cluster test production

# Validate kubeconfig
kubectl --context production get nodes
```

#### Authentication Problems
```bash
# Verify authentication
kubectl --context production auth can-i get gitrepositories

# Check token expiration
kubectl --context production auth whoami
```

#### Performance Issues
```bash
# Reduce refresh interval
fluxcli config set refresh_interval 30s

# Limit concurrent clusters
fluxcli config set max_concurrent_clusters 5
```

### Debug Mode

Enable verbose logging:
```bash
fluxcli --debug --log-level trace
```

### Health Checks

Built-in health monitoring:
- Periodic connectivity tests
- Resource access validation
- Performance metrics collection