package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/malagant/fluxcli/internal/config"
	"github.com/malagant/fluxcli/pkg/core"
	"github.com/malagant/fluxcli/pkg/k8s"
)

// AppModel represents the main application model
type AppModel struct {
	config          *config.Config
	manager         *core.Manager
	state           AppState
	currentView     ViewType
	resourceView    *ResourceView
	eventView       *EventView
	commandMode     bool
	commandInput    string
	statusMessage   string
	errorMessage    string
	width           int
	height          int
	ready           bool
}

// AppState represents the application state
type AppState struct {
	Resources       map[string]map[k8s.ResourceType][]k8s.Resource
	Events          map[string][]Event
	CurrentCluster  string
	CurrentResource k8s.ResourceType
	Filter          string
	ShowHelp        bool
}

// ViewType represents different view types
type ViewType int

const (
	ViewResources ViewType = iota
	ViewEvents
	ViewDetails
)

// Event represents a Kubernetes event for display
type Event struct {
	Type      string
	Reason    string
	Object    string
	Message   string
	Timestamp string
	Count     int
}

// NewApp creates a new FluxCLI application
func NewApp(cfg *config.Config) *AppModel {
	manager := core.NewManager(cfg)
	
	app := &AppModel{
		config:      cfg,
		manager:     manager,
		currentView: ViewResources,
		state: AppState{
			Resources:       make(map[string]map[k8s.ResourceType][]k8s.Resource),
			Events:          make(map[string][]Event),
			CurrentCluster:  cfg.CurrentContext,
			CurrentResource: k8s.ResourceTypeGitRepository,
		},
	}

	app.resourceView = NewResourceView(cfg)
	app.eventView = NewEventView(cfg)

	return app
}

// Run starts the TUI application
func (m *AppModel) Run() error {
	if err := m.manager.Start(); err != nil {
		return fmt.Errorf("failed to start manager: %w", err)
	}
	defer m.manager.Stop()

	program := tea.NewProgram(m, tea.WithAltScreen())
	
	// Start background update handlers
	go m.handleUpdates(program)
	
	_, err := program.Run()
	return err
}

// Init initializes the model
func (m *AppModel) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.resourceView.Init(),
		m.eventView.Init(),
	)
}

// Update handles messages and updates the model
func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		
		// Update child views
		m.resourceView.SetSize(m.width, m.height-4) // Reserve space for header/footer
		m.eventView.SetSize(m.width, m.height/3)
		
	case tea.KeyMsg:
		if m.commandMode {
			return m.handleCommandMode(msg)
		}
		return m.handleNormalMode(msg)
		
	case ResourceUpdateMsg:
		m.handleResourceUpdate(msg)
		
	case EventUpdateMsg:
		m.handleEventUpdate(msg)
		
	case ErrorUpdateMsg:
		m.errorMessage = msg.Error
		
	case ClearStatusMsg:
		m.statusMessage = ""
		m.errorMessage = ""
	}

	// Update current view
	switch m.currentView {
	case ViewResources:
		m.resourceView, cmd = m.resourceView.Update(msg)
		cmds = append(cmds, cmd)
	case ViewEvents:
		m.eventView, cmd = m.eventView.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the application
func (m *AppModel) View() string {
	if !m.ready {
		return "Initializing FluxCLI..."
	}

	var view strings.Builder
	
	// Header
	view.WriteString(m.renderHeader())
	view.WriteString("\n")
	
	// Main content
	switch m.currentView {
	case ViewResources:
		view.WriteString(m.resourceView.View())
	case ViewEvents:
		view.WriteString(m.eventView.View())
	}
	
	// Footer
	view.WriteString("\n")
	view.WriteString(m.renderFooter())
	
	return view.String()
}

// handleNormalMode handles keyboard input in normal mode
func (m *AppModel) handleNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
		
	case ":":
		m.commandMode = true
		m.commandInput = ""
		return m, nil
		
	case "?":
		m.state.ShowHelp = !m.state.ShowHelp
		return m, nil
		
	case "/":
		// TODO: Implement search/filter mode
		m.statusMessage = "Search mode not yet implemented"
		return m, tea.Tick(2000, func(time.Time) tea.Msg { return ClearStatusMsg{} })
		
	case "tab":
		// Switch between views
		switch m.currentView {
		case ViewResources:
			m.currentView = ViewEvents
		case ViewEvents:
			m.currentView = ViewResources
		}
		return m, nil
		
	case "1":
		m.state.CurrentResource = k8s.ResourceTypeGitRepository
		m.resourceView.SetResourceType(k8s.ResourceTypeGitRepository)
		
	case "2":
		m.state.CurrentResource = k8s.ResourceTypeHelmRepository
		m.resourceView.SetResourceType(k8s.ResourceTypeHelmRepository)
		
	case "3":
		m.state.CurrentResource = k8s.ResourceTypeKustomization
		m.resourceView.SetResourceType(k8s.ResourceTypeKustomization)
		
	case "4":
		m.state.CurrentResource = k8s.ResourceTypeHelmRelease
		m.resourceView.SetResourceType(k8s.ResourceTypeHelmRelease)
		
	case "ctrl+k":
		// Previous cluster
		clusters := m.manager.GetClusters()
		if len(clusters) > 1 {
			current := m.state.CurrentCluster
			for i, cluster := range clusters {
				if cluster == current {
					prev := (i - 1 + len(clusters)) % len(clusters)
					m.state.CurrentCluster = clusters[prev]
					m.manager.SetCurrentCluster(clusters[prev])
					break
				}
			}
		}
		
	case "ctrl+j":
		// Next cluster
		clusters := m.manager.GetClusters()
		if len(clusters) > 1 {
			current := m.state.CurrentCluster
			for i, cluster := range clusters {
				if cluster == current {
					next := (i + 1) % len(clusters)
					m.state.CurrentCluster = clusters[next]
					m.manager.SetCurrentCluster(clusters[next])
					break
				}
			}
		}
		
	case "r":
		// Manual refresh
		m.statusMessage = "Refreshing resources..."
		cmds = append(cmds, tea.Tick(2000, func(time.Time) tea.Msg { return ClearStatusMsg{} }))
	}
	
	return m, tea.Batch(cmds...)
}

// handleCommandMode handles keyboard input in command mode
func (m *AppModel) handleCommandMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		cmd := m.executeCommand(m.commandInput)
		m.commandMode = false
		m.commandInput = ""
		return m, cmd
		
	case "esc":
		m.commandMode = false
		m.commandInput = ""
		return m, nil
		
	case "backspace":
		if len(m.commandInput) > 0 {
			m.commandInput = m.commandInput[:len(m.commandInput)-1]
		}
		
	default:
		if len(msg.String()) == 1 {
			m.commandInput += msg.String()
		}
	}
	
	return m, nil
}

// executeCommand executes a command entered in command mode
func (m *AppModel) executeCommand(command string) tea.Cmd {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil
	}
	
	cmd := parts[0]
	args := parts[1:]
	
	switch cmd {
	case "quit", "q":
		return tea.Quit
		
	case "suspend", "s":
		if len(args) > 0 {
			resourceName := args[0]
			if err := m.manager.SuspendResource(m.state.CurrentResource, resourceName); err != nil {
				m.errorMessage = fmt.Sprintf("Failed to suspend %s: %v", resourceName, err)
			} else {
				m.statusMessage = fmt.Sprintf("Suspended %s", resourceName)
			}
		}
		
	case "resume", "r":
		if len(args) > 0 {
			resourceName := args[0]
			if err := m.manager.ResumeResource(m.state.CurrentResource, resourceName); err != nil {
				m.errorMessage = fmt.Sprintf("Failed to resume %s: %v", resourceName, err)
			} else {
				m.statusMessage = fmt.Sprintf("Resumed %s", resourceName)
			}
		}
		
	case "reconcile", "rec":
		if len(args) > 0 {
			resourceName := args[0]
			if err := m.manager.ReconcileResource(m.state.CurrentResource, resourceName); err != nil {
				m.errorMessage = fmt.Sprintf("Failed to reconcile %s: %v", resourceName, err)
			} else {
				m.statusMessage = fmt.Sprintf("Triggered reconciliation for %s", resourceName)
			}
		}
		
	default:
		m.errorMessage = fmt.Sprintf("Unknown command: %s", cmd)
	}
	
	return tea.Tick(3000, func(time.Time) tea.Msg { return ClearStatusMsg{} })
}

// renderHeader renders the application header
func (m *AppModel) renderHeader() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Render("FluxCLI")
		
	cluster := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Render(fmt.Sprintf("Cluster: %s", m.state.CurrentCluster))
		
	resource := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("81")).
		Render(fmt.Sprintf("Resource: %s", m.state.CurrentResource))
		
	namespace := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("226")).
		Render(fmt.Sprintf("Namespace: %s", m.manager.GetCurrentNamespace()))
	
	if m.commandMode {
		commandPrompt := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("196")).
			Render(fmt.Sprintf(":%s", m.commandInput))
		return fmt.Sprintf("%s | %s | %s | %s | %s", title, cluster, resource, namespace, commandPrompt)
	}
	
	return fmt.Sprintf("%s | %s | %s | %s", title, cluster, resource, namespace)
}

// renderFooter renders the application footer
func (m *AppModel) renderFooter() string {
	var footer strings.Builder
	
	if m.errorMessage != "" {
		error := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Render(fmt.Sprintf("Error: %s", m.errorMessage))
		footer.WriteString(error)
	} else if m.statusMessage != "" {
		status := lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Render(m.statusMessage)
		footer.WriteString(status)
	} else if m.state.ShowHelp {
		help := m.renderHelp()
		footer.WriteString(help)
	} else {
		shortcuts := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("? help | 1-4 resource types | tab switch views | : command mode | ctrl+k/j clusters | q quit")
		footer.WriteString(shortcuts)
	}
	
	return footer.String()
}

// renderHelp renders the help text
func (m *AppModel) renderHelp() string {
	helpText := `
Navigation:
  j/k or ↓/↑    Navigate up/down
  g/G           Go to top/bottom
  enter         View details
  tab           Switch between views
  
Resource Types:
  1             GitRepositories
  2             HelmRepositories  
  3             Kustomizations
  4             HelmReleases
  
Clusters:
  ctrl+k/j      Previous/Next cluster
  
Commands (: to enter command mode):
  suspend <name>    Suspend resource
  resume <name>     Resume resource
  reconcile <name>  Trigger reconciliation
  
Other:
  /             Search/Filter (coming soon)
  r             Manual refresh
  ?             Toggle this help
  q             Quit
`
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		Render(strings.TrimSpace(helpText))
}

// Message types for updates
type ResourceUpdateMsg struct {
	Cluster   string
	Resources []k8s.Resource
	Type      k8s.ResourceType
}

type EventUpdateMsg struct {
	Cluster string
	Events  []Event
}

type ErrorUpdateMsg struct {
	Error string
}

type ClearStatusMsg struct{}

// handleUpdates handles background updates from the manager
func (m *AppModel) handleUpdates(program *tea.Program) {
	for {
		select {
		case update := <-m.manager.GetResourceUpdates():
			program.Send(ResourceUpdateMsg{
				Cluster:   update.Cluster,
				Resources: update.Resources,
				Type:      update.Type,
			})
			
		case update := <-m.manager.GetEventUpdates():
			events := make([]Event, len(update.Events))
			for i, event := range update.Events {
				events[i] = Event{
					Type:      event.Type,
					Reason:    event.Reason,
					Object:    fmt.Sprintf("%s/%s", event.InvolvedObject.Kind, event.InvolvedObject.Name),
					Message:   event.Message,
					Timestamp: event.FirstTimestamp.Format("15:04:05"),
					Count:     int(event.Count),
				}
			}
			program.Send(EventUpdateMsg{
				Cluster: update.Cluster,
				Events:  events,
			})
			
		case update := <-m.manager.GetErrorUpdates():
			program.Send(ErrorUpdateMsg{
				Error: fmt.Sprintf("[%s] %v", update.Cluster, update.Error),
			})
		}
	}
}

// handleResourceUpdate handles resource updates
func (m *AppModel) handleResourceUpdate(msg ResourceUpdateMsg) {
	if m.state.Resources[msg.Cluster] == nil {
		m.state.Resources[msg.Cluster] = make(map[k8s.ResourceType][]k8s.Resource)
	}
	m.state.Resources[msg.Cluster][msg.Type] = msg.Resources
	
	// Update resource view if it matches current view
	if msg.Cluster == m.state.CurrentCluster && msg.Type == m.state.CurrentResource {
		m.resourceView.SetResources(msg.Resources)
	}
}

// handleEventUpdate handles event updates  
func (m *AppModel) handleEventUpdate(msg EventUpdateMsg) {
	m.state.Events[msg.Cluster] = msg.Events
	
	// Update event view if it matches current cluster
	if msg.Cluster == m.state.CurrentCluster {
		m.eventView.SetEvents(msg.Events)
	}
}
