# UI Navigation Specification

This document defines the user interface navigation patterns, keyboard shortcuts, and interaction models for FluxCLI.

## Interface Layout

### Main Application Window

```
┌─ FluxCLI - [CLUSTER] ─ NAMESPACE ─────────────── STATUS ┐
│ Connected: ✓  │  Resources: 25  │  Last Update: 10:30:15 │
├─────────────────────────────────────────────────────────┤
│ Resource Type: GitRepositories                   [1/5]  │
├─────────────────────────────────────────────────────────┤
│ FILTERS: [namespace:flux-system] [status:all]   [Clear] │
├─────────────────────────────────────────────────────────┤
│ NAME              READY   STATUS        AGE     MESSAGE │
│ ▶ flux-system     True    Stored        2d      OK      │
│   podinfo         False   Failed        4h      Auth    │
│   microservices   True    Stored        12h     OK      │
│   ...                                                   │
├─────────────────────────────────────────────────────────┤
│ Events                                           [3/10] │
│ 10:29:45 GitRepository/podinfo reconciliation failed    │
│ 10:25:12 Kustomization/apps successfully reconciled     │
│ 10:20:33 HelmRelease/nginx upgrade completed            │
└─────────────────────────────────────────────────────────┘
Help: ?  │  Commands: :  │  Search: /  │  Quit: q
```

### Pane Structure

**Header Pane** (Fixed Height: 3 lines):
- Cluster name and connection status
- Current namespace and resource counts
- Last update timestamp

**Filter Bar** (Conditional Height: 0-1 lines):
- Active filters display
- Clear filters option
- Only shown when filters are active

**Resource List Pane** (Variable Height):
- Main content area
- Resource table with scrolling
- Selection indicator and cursor

**Events Pane** (Fixed Height: 4 lines):
- Recent FluxCD events
- Real-time event streaming
- Scrollable event history

**Status Bar** (Fixed Height: 1 line):
- Key binding hints
- Current mode indicator
- Help text

## Navigation Modes

### Normal Mode (Default)

**Purpose**: Primary navigation and resource browsing

**Available Actions**:
- Move cursor up/down through resources
- View resource details
- Execute quick operations
- Switch between panes
- Enter other modes

**Key Bindings**:
| Key | Action | Description |
|-----|--------|-------------|
| `j` / `↓` | Move Down | Next resource |
| `k` / `↑` | Move Up | Previous resource |
| `g` | Go to Top | First resource |
| `G` | Go to Bottom | Last resource |
| `Enter` | View Details | Open resource detail view |
| `Space` | Toggle Select | Multi-select resource |
| `Tab` | Next Pane | Switch to events pane |
| `Shift+Tab` | Previous Pane | Switch pane reverse |
| `r` | Refresh | Refresh current view |
| `?` | Help | Show help overlay |
| `:` | Command Mode | Enter command interface |
| `/` | Search Mode | Enter search/filter |
| `q` | Quit | Exit application |

### Command Mode

**Purpose**: Advanced operations and configuration

**Activation**: Press `:` in normal mode

**Interface**:
```
┌─ Command ──────────────────────────────────────────────┐
│ :gitrepository flux-system suspend                    │
└────────────────────────────────────────────────────────┘
```

**Features**:
- Command history (↑/↓ arrows)
- Auto-completion (Tab key)
- Command validation
- Error display

**Command Categories**:
- **Navigation**: Resource type switching
- **Operations**: Resource manipulation
- **Configuration**: Settings management
- **Cluster**: Multi-cluster operations

### Search/Filter Mode

**Purpose**: Find and filter resources

**Activation**: Press `/` in normal mode

**Interface**:
```
┌─ Search/Filter ────────────────────────────────────────┐
│ /ready:false namespace:kube-system                     │
└────────────────────────────────────────────────────────┘
```

**Search Syntax**:
- Simple text: `podinfo`
- Field filters: `ready:false`
- Namespace filter: `namespace:kube-system`
- Label filter: `label:app=nginx`
- Age filter: `age>1d`

### Detail Mode

**Purpose**: View detailed resource information

**Activation**: Press `Enter` on selected resource

**Layout**:
```
┌─ GitRepository: flux-system ───────────────────────────┐
│ Name: flux-system                                      │
│ Namespace: flux-system                                 │
│ Ready: True                                            │
│ Status: Stored                                         │
│ URL: https://github.com/fluxcd/flux2-kustomize-helm   │
│ Branch: main                                           │
│ Revision: main/a1b2c3d4e5f                            │
│ Age: 2d5h23m                                           │
├────────────────────────────────────────────────────────┤
│ Conditions:                                            │
│ ✓ Ready: True (Stored)                                 │
│   Last Transition: 2025-01-19T08:25:14Z               │
│   Message: stored revision: main/a1b2c3d4e5f           │
├────────────────────────────────────────────────────────┤
│ Recent Events:                                         │
│ 08:25:14 Successfully stored revision                  │
│ 08:20:10 Reconciliation started                       │
│ 08:15:05 Previous reconciliation completed            │
├────────────────────────────────────────────────────────┤
│ Spec (YAML):                                           │
│ apiVersion: source.toolkit.fluxcd.io/v1               │
│ kind: GitRepository                                    │
│ metadata:                                              │
│   name: flux-system                                    │
│   namespace: flux-system                               │
│ spec:                                                  │
│   interval: 1m0s                                       │
│   ref:                                                 │
│     branch: main                                       │
│   url: https://github.com/fluxcd/flux2-kustomize-helm │
└────────────────────────────────────────────────────────┘
[s] Suspend [r] Resume [f] Force Reconcile [e] Edit [q] Back
```

**Navigation in Detail Mode**:
| Key | Action |
|-----|--------|
| `↓`/`j` | Scroll down |
| `↑`/`k` | Scroll up |
| `g` | Go to top |
| `G` | Go to bottom |
| `s` | Suspend resource |
| `r` | Resume resource |
| `f` | Force reconcile |
| `e` | Edit resource |
| `d` | Delete resource |
| `q` | Back to list |

### Help Mode

**Purpose**: Display keyboard shortcuts and commands

**Activation**: Press `?` in any mode

**Layout**:
```
┌─ Help ─────────────────────────────────────────────────┐
│                    FluxCLI Help                        │
├────────────────────────────────────────────────────────┤
│ Navigation:            │ Operations:                   │
│   j/↓     Move down    │   s        Suspend resource   │
│   k/↑     Move up      │   r        Resume resource    │
│   g       Go to top    │   f        Force reconcile    │
│   G       Go to bottom │   d        Delete resource    │
│   Enter   View details │   Space    Toggle selection   │
│   Tab     Switch pane  │                               │
├────────────────────────────────────────────────────────┤
│ Modes:                 │ Quick Commands:               │
│   :       Command mode │   :gr      GitRepositories   │
│   /       Search mode  │   :hr      HelmRepositories  │
│   ?       This help    │   :ks      Kustomizations    │
│   q       Quit app     │   :helm    HelmReleases       │
├────────────────────────────────────────────────────────┤
│ Multi-Cluster:         │ Filtering:                    │
│   Ctrl+K  Next cluster │   /text    Search by name    │
│   Ctrl+J  Prev cluster │   /ready:false Filter status │
│   :cluster <name>      │   /ns:name Filter namespace  │
└────────────────────────────────────────────────────────┘
Press any key to close help
```

## Resource Type Views

### View Switching

**Resource Type Navigation**:
| Command | Target View |
|---------|-------------|
| `:gr` or `:gitrepository` | GitRepositories |
| `:hr` or `:helmrepository` | HelmRepositories |
| `:ks` or `:kustomization` | Kustomizations |
| `:helmrelease` | HelmReleases |
| `:rs` or `:resourceset` | ResourceSets |
| `:all` | All resource types |

**View Indicators**:
- Resource type displayed in header
- Resource count and pagination info
- Active filters shown in filter bar

### All Resources View

**Layout**: Mixed resource types in single list

**Format**:
```
TYPE            NAME              READY   STATUS     AGE
GitRepository   flux-system       True    Stored     2d
GitRepository   podinfo           False   Failed     4h
HelmRepository  bitnami           True    Indexed    1d
Kustomization   apps              True    Applied    1d
HelmRelease     nginx             True    Released   2d
```

**Sorting Options**:
- By type (default)
- By name
- By age
- By status
- By ready state

## Multi-Cluster Navigation

### Cluster Switching

**Visual Indicators**:
- Current cluster name in header
- Cluster status icon (connected/disconnected)
- Cluster-specific color coding (optional)

**Switching Methods**:
| Method | Usage |
|--------|-------|
| `Ctrl+K` | Next cluster in list |
| `Ctrl+J` | Previous cluster in list |
| `:cluster <name>` | Switch to specific cluster |
| `:clusters` | Show cluster selection menu |

### Cluster Selection Menu

**Layout**:
```
┌─ Select Cluster ───────────────────────────────────────┐
│ ▶ production      Connected ✓    Resources: 45        │
│   staging         Connected ✓    Resources: 23        │
│   development     Connecting ⟳   Resources: --        │
│   testing         Failed ✗        Last seen: 5m ago   │
├────────────────────────────────────────────────────────┤
│ [a] All Clusters View  [t] Test Connection  [q] Cancel │
└────────────────────────────────────────────────────────┘
```

### All Clusters View

**Purpose**: View resources across all clusters simultaneously

**Activation**: `Ctrl+A` or `:all-clusters`

**Format**:
```
CLUSTER      TYPE            NAME           READY   STATUS
production   GitRepository   flux-system    True    Stored
production   HelmRelease     nginx          True    Released
staging      GitRepository   flux-system    True    Stored
staging      HelmRelease     nginx          False   Failed
development  GitRepository   flux-system    True    Stored
```

## Filtering and Search

### Filter Types

**Status Filters**:
- `ready:true` - Only ready resources
- `ready:false` - Only failed/not ready resources
- `status:failed` - Resources with failed status
- `status:progressing` - Resources currently updating

**Namespace Filters**:
- `namespace:flux-system` - Specific namespace
- `namespace:kube-*` - Wildcard matching
- `namespace:!default` - Exclude namespace

**Age Filters**:
- `age>1d` - Older than 1 day
- `age<1h` - Newer than 1 hour
- `age:5m-1h` - Between 5 minutes and 1 hour

**Label Filters**:
- `label:app=nginx` - Specific label value
- `label:env` - Has label key
- `label:!debug` - Does not have label

### Search Interface

**Real-time Filtering**:
- Results update as user types
- Highlight matching text in results
- Clear indication of active filters

**Search History**:
- Last 10 searches saved
- Navigate with ↑/↓ in search mode
- Persistent across sessions

**Filter Combinations**:
- Multiple filters with AND logic
- Space-separated filter terms
- Parentheses for complex logic (future)

## Keyboard Shortcuts Summary

### Global Shortcuts (Available in all modes)

| Key | Action |
|-----|--------|
| `Ctrl+C` | Emergency quit |
| `F1` | Help (alternative to `?`) |
| `F5` | Force refresh |
| `Ctrl+L` | Clear screen/redraw |

### Context-Sensitive Shortcuts

**Resource List Context**:
| Key | Action |
|-----|--------|
| `s` | Suspend selected resource |
| `r` | Resume selected resource |
| `f` | Force reconcile resource |
| `d` | Delete resource |
| `e` | Edit resource |
| `i` | Inspect/describe resource |
| `l` | View logs (if applicable) |
| `y` | Copy resource name to clipboard |

**Multi-Selection Context**:
| Key | Action |
|-----|--------|
| `Space` | Toggle current item selection |
| `Ctrl+A` | Select all visible |
| `Escape` | Clear all selections |
| `Enter` | Execute operation on selected |

### Customizable Shortcuts

**Configuration File**: `~/.fluxcli/keybindings.yaml`

**Example Configuration**:
```yaml
keybindings:
  global:
    quit: ["q", "Ctrl+C"]
    help: ["?", "F1"]
    refresh: ["r", "F5"]

  navigation:
    up: ["k", "Up", "Ctrl+P"]
    down: ["j", "Down", "Ctrl+N"]
    top: ["g", "Home"]
    bottom: ["G", "End"]

  operations:
    suspend: ["s"]
    resume: ["r"]
    reconcile: ["f", "F6"]
    delete: ["d", "Del"]
    edit: ["e"]

  modes:
    command: [":"]
    search: ["/", "Ctrl+F"]
    help: ["?", "F1"]

  cluster:
    next: ["Ctrl+K"]
    previous: ["Ctrl+J"]
    select: ["Ctrl+S"]
    all: ["Ctrl+A"]
```

## Responsive Design

### Terminal Size Adaptation

**Minimum Size**: 80x24 characters
**Optimal Size**: 120x30 characters

**Responsive Behavior**:
- Hide less critical columns in narrow terminals
- Adjust pane heights based on terminal height
- Horizontal scrolling for wide content
- Graceful degradation of features

### Column Priority (for narrow displays)

**High Priority** (always shown):
- Resource name
- Ready status
- Age

**Medium Priority** (shown if space available):
- Status message
- Namespace (if not filtered)

**Low Priority** (shown in wide terminals):
- Detailed status
- Additional metadata

## Accessibility

### Screen Reader Support

**ARIA Labels**: Meaningful labels for all interface elements
**Navigation Announcements**: Clear indication of focus changes
**Status Updates**: Announce important status changes

### High Contrast Mode

**Color Schemes**:
- High contrast theme option
- Bold text for important information
- Clear visual separation of elements

### Keyboard-Only Operation

**Complete Functionality**: All features accessible via keyboard
**Focus Indicators**: Clear indication of focused elements
**Logical Tab Order**: Intuitive navigation flow