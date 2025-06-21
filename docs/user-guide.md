# FluxCLI User Guide

This comprehensive guide covers all aspects of using FluxCLI for managing FluxCD resources across Kubernetes clusters.

## Getting Started

### First Launch

When you first launch FluxCLI, it will:

1. Load your kubeconfig from the default location
2. Detect available Kubernetes contexts
3. Connect to the current context
4. Scan for FluxCD resources

```bash
fluxcli
```

### Initial Setup

Configure your clusters for optimal experience:

```bash
# Add a new cluster
:cluster add production --context prod-cluster --kubeconfig ~/.kube/prod

# Set default namespace
:config set namespace flux-system

# Configure refresh interval
:config set refresh_interval 10s
```

## Navigation

### Main Interface

The FluxCLI interface consists of several key areas:

```
┌─ FluxCLI - [PRODUCTION] ─ flux-system ────────────────────┐
│ Status: Connected ✓  │  Last Update: 2025-01-21 10:30:15  │
├─────────────────────────────────────────────────────────────┤
│ Resource Type: GitRepositories                    [1/5]     │
├─────────────────────────────────────────────────────────────┤
│ NAME              READY   STATUS        AGE     MESSAGE     │
│ ▶ flux-system     True    Stored        2d      Stored     │
│   podinfo         False   Failed        4h      Auth fail  │
│   microservices   True    Stored        12h     Stored     │
├─────────────────────────────────────────────────────────────┤
│ Events                                                      │
│ 10:29:45 GitRepository/podinfo reconciliation failed       │
│ 10:25:12 Kustomization/apps successfully reconciled        │
└─────────────────────────────────────────────────────────────┘
Press ? for help  │  : for commands  │  / to search  │  q to quit
```

### Keyboard Navigation

#### Global Navigation

| Key | Action |
|-----|--------|
| `?` | Show help screen |
| `q` | Quit application |
| `:` | Enter command mode |
| `/` | Search/filter resources |
| `Esc` | Cancel current operation |
| `Tab` | Switch between panes |
| `Shift+Tab` | Switch between panes (reverse) |

#### Resource Navigation

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `g` | Go to top |
| `G` | Go to bottom |
| `Enter` | View resource details |
| `Space` | Toggle resource selection |
| `r` | Refresh current view |

#### Cluster Navigation

| Key | Action |
|-----|--------|
| `Ctrl+K` | Next cluster |
| `Ctrl+J` | Previous cluster |
| `Ctrl+A` | Show all clusters view |
| `Ctrl+H` | Cluster health dashboard |

## Resource Management

### Resource Types

FluxCLI supports the following FluxCD resource types:

1. **GitRepository** (`:gr` or `:gitrepo`)
2. **HelmRepository** (`:hr` or `:helmrepo`)
3. **Kustomization** (`:ks` or `:kustomization`)
4. **HelmRelease** (`:helmrelease`)
5. **ResourceSet** (`:rs` or `:resourceset`)

### Viewing Resources

#### List Resources
```bash
# View GitRepositories
:gr

# View all resource types
:all

# View resources in specific namespace
:gr -n kube-system
```

#### Resource Details

Press `Enter` on any resource to view detailed information:

```
┌─ GitRepository: flux-system ───────────────────────────────┐
│ Name: flux-system                                          │
│ Namespace: flux-system                                     │
│ Ready: True                                                │
│ Status: Stored                                             │
│ URL: https://github.com/fluxcd/flux2-kustomize-helm       │
│ Branch: main                                               │
│ Revision: main/a1b2c3d                                     │
│ Age: 2d5h                                                  │
├─────────────────────────────────────────────────────────────┤
│ Conditions:                                                │
│ ✓ Ready: True (Stored)                                     │
│   Last Transition: 2025-01-19T08:25:14Z                   │
│   Message: Stored revision: main/a1b2c3d                   │
├─────────────────────────────────────────────────────────────┤
│ Recent Events:                                             │
│ 08:25:14 Successfully stored revision                      │
│ 08:20:10 Reconciliation started                           │
│ 08:15:05 Previous reconciliation completed                │
└─────────────────────────────────────────────────────────────┘
[d] Delete  [s] Suspend  [r] Resume  [f] Force Reconcile  [q] Back
```

### Resource Operations

#### Suspend/Resume Resources

```bash
# Suspend a resource
:suspend gitrepository/flux-system

# Resume a resource
:resume gitrepository/flux-system

# Suspend all resources of a type
:suspend gitrepository --all
```

#### Force Reconciliation

```bash
# Force reconcile a specific resource
:reconcile gitrepository/flux-system

# Reconcile with annotation
:reconcile gitrepository/flux-system --with-source
```

#### Delete Resources

```bash
# Delete a resource (with confirmation)
:delete gitrepository/old-repo

# Force delete without confirmation
:delete gitrepository/old-repo --force
```

## Command Mode

Enter command mode with `:` to access advanced features.

### Resource Commands

```bash
# Navigate to resource types
:gitrepository, :gr          # GitRepositories
:helmrepository, :hr         # HelmRepositories  
:kustomization, :ks          # Kustomizations
:helmrelease                 # HelmReleases
:resourceset, :rs            # ResourceSets
:all                         # All resources

# Resource operations
:describe <resource>         # Show detailed info
:edit <resource>            # Edit resource YAML
:logs <resource>            # Show controller logs
:events <resource>          # Show resource events
```

### Cluster Commands

```bash
# Cluster management
:cluster list               # List configured clusters
:cluster add <name>         # Add new cluster
:cluster remove <name>      # Remove cluster
:cluster test <name>        # Test cluster connectivity
:cluster <name>             # Switch to cluster

# Multi-cluster operations
:all-clusters              # View all clusters
:health                    # Cluster health dashboard
```

### Configuration Commands

```bash
# View configuration
:config show

# Set configuration values
:config set refresh_interval 30s
:config set namespace flux-system
:config set theme dark

# Reset configuration
:config reset
```

## Search and Filtering

### Quick Search

Press `/` to enter search mode:

```bash
# Search by name
/flux-system

# Search by status
/ready:false

# Search by namespace
/namespace:kube-system

# Combined search
/ready:false namespace:default
```

### Advanced Filtering

```bash
# Filter by conditions
:filter ready=true
:filter status=failed

# Filter by labels
:filter label=app=podinfo

# Filter by age
:filter age>1d
:filter age<1h

# Clear filters
:filter clear
```

## Monitoring and Events

### Event Streaming

FluxCLI provides real-time event monitoring:

```bash
# Show all events
:events

# Filter events by type
:events --type=Warning

# Follow events in real-time
:events --follow

# Events for specific resource
:events gitrepository/flux-system
```

### Resource Monitoring

```bash
# Watch resource changes
:watch gitrepository/flux-system

# Monitor reconciliation progress
:reconcile gitrepository/flux-system --watch

# Health monitoring
:health --watch
```

## Customization

### Themes

```bash
# Available themes
:theme list

# Set theme
:theme dark
:theme light
:theme solarized

# Custom colors
:config set color.ready green
:config set color.failed red
```

### Layout Options

```bash
# Adjust column widths
:config set columns.name 30
:config set columns.status 15

# Toggle columns
:config set show.age false
:config set show.message true

# Pane sizes
:config set pane.events.height 5
:config set pane.details.width 50
```

### Keyboard Shortcuts

Customize keyboard shortcuts in `~/.fluxcli/keybindings.yaml`:

```yaml
keybindings:
  global:
    quit: ["q", "Ctrl+C"]
    help: ["?", "F1"]
    command: [":"]
    search: ["/", "Ctrl+F"]
  
  navigation:
    up: ["k", "Up"]
    down: ["j", "Down"]
    top: ["g"]
    bottom: ["G"]
    
  operations:
    suspend: ["s"]
    resume: ["r"]
    reconcile: ["f"]
    delete: ["d"]
```

## Troubleshooting

### Common Issues

#### "No FluxCD resources found"
- Verify FluxCD is installed: `kubectl get crd | grep fluxcd`
- Check namespace: `kubectl get gitrepositories -A`
- Verify RBAC permissions

#### "Connection failed"
- Test kubeconfig: `kubectl cluster-info`
- Check cluster configuration: `:cluster test <name>`
- Verify network connectivity

#### "Permission denied"
- Check RBAC: `kubectl auth can-i get gitrepositories`
- Verify service account permissions
- Contact cluster administrator

### Debug Mode

Enable verbose logging:

```bash
# Launch with debug mode
fluxcli --debug

# Set log level
fluxcli --log-level trace

# Save logs to file
fluxcli --log-file fluxcli.log
```

### Performance Tuning

```bash
# Reduce refresh frequency
:config set refresh_interval 60s

# Limit concurrent connections
:config set max_connections 5

# Disable real-time events
:config set events.enabled false
```

## Tips and Best Practices

### Workflow Optimization

1. **Use Aliases**: Create short aliases for frequently used commands
2. **Filter Early**: Use filters to reduce visual clutter
3. **Multiple Clusters**: Keep cluster contexts organized
4. **Event Monitoring**: Enable event streaming for active troubleshooting

### Security Best Practices

1. **Least Privilege**: Use minimal required RBAC permissions
2. **Credential Management**: Rotate kubeconfig credentials regularly
3. **Audit Trail**: Monitor FluxCLI access logs
4. **Network Security**: Use secure cluster connections

### Performance Tips

1. **Selective Monitoring**: Monitor only critical resources in real-time
2. **Batch Operations**: Use multi-resource operations when possible
3. **Connection Pooling**: Configure appropriate connection limits
4. **Cache Management**: Understand caching behavior for large clusters