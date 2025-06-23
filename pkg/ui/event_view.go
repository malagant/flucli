package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/malagant/fluxcli/internal/config"
)

// EventView displays Kubernetes events in a table
type EventView struct {
	config *config.Config
	table  table.Model
	events []Event
	width  int
	height int
}

// NewEventView creates a new event view
func NewEventView(cfg *config.Config) *EventView {
	columns := []table.Column{
		{Title: "Type", Width: 6},
		{Title: "Reason", Width: 12},
		{Title: "Object", Width: 22},
		{Title: "Message", Width: 60},
		{Title: "Time", Width: 8},
		{Title: "Count", Width: 5},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(false), // Events view is not focused by default
		table.WithHeight(cfg.UI.PaneEventsHeight),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return &EventView{
		config: cfg,
		table:  t,
	}
}

// Init initializes the event view
func (v *EventView) Init() tea.Cmd {
	return nil
}

// Update handles messages for the event view
func (v *EventView) Update(msg tea.Msg) (*EventView, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle arrow keys by checking Type directly
		switch msg.Type {
		case tea.KeyDown:
			v.table, cmd = v.table.Update(msg)
		case tea.KeyUp:
			v.table, cmd = v.table.Update(msg)
		case tea.KeyLeft:
			v.table, cmd = v.table.Update(msg)
		case tea.KeyRight:
			v.table, cmd = v.table.Update(msg)
		case tea.KeyPgDown:
			v.table, cmd = v.table.Update(msg)
		case tea.KeyPgUp:
			v.table, cmd = v.table.Update(msg)
		case tea.KeyHome:
			if len(v.events) > 0 {
				v.table.GotoTop()
			}
		case tea.KeyEnd:
			if len(v.events) > 0 {
				v.table.GotoBottom()
			}
		case tea.KeyEnter, tea.KeySpace:
			// TODO: Show event details
			return v, nil
		default:
			// Handle string-based keys
			switch msg.String() {
			// Vertical navigation - j/k for vim users
			case "j":
				v.table, cmd = v.table.Update(tea.KeyMsg{Type: tea.KeyDown})
			case "k":
				v.table, cmd = v.table.Update(tea.KeyMsg{Type: tea.KeyUp})
			
			// Page navigation with vim-style shortcuts
			case "ctrl+d":
				v.table, cmd = v.table.Update(tea.KeyMsg{Type: tea.KeyPgDown})
			case "ctrl+u":
				v.table, cmd = v.table.Update(tea.KeyMsg{Type: tea.KeyPgUp})
			
			// Vim-style navigation
			case "g":
				// Go to top
				if len(v.events) > 0 {
					v.table.GotoTop()
				}
			case "G":
				// Go to bottom
				if len(v.events) > 0 {
					v.table.GotoBottom()
				}
			case "H":
				// Go to top of visible area
				if len(v.events) > 0 {
					v.table.GotoTop()
				}
			case "M":
				// Go to middle of visible area
				if len(v.events) > 0 {
					middle := len(v.events) / 2
					for i := 0; i < middle; i++ {
						v.table, _ = v.table.Update(tea.KeyMsg{Type: tea.KeyDown})
					}
				}
			case "L":
				// Go to bottom of visible area
				if len(v.events) > 0 {
					v.table.GotoBottom()
				}
			}
		}
	}
	
	return v, cmd
}

// View renders the event view
func (v *EventView) View() string {
	if len(v.events) == 0 {
		emptyMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Render("No events found")
		
		// Create a bordered box for consistency
		box := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Height(v.config.UI.PaneEventsHeight).
			Padding(1, 2).
			Render(emptyMsg)
		
		return box
	}
	
	return v.table.View()
}

// SetEvents sets the events to display
func (v *EventView) SetEvents(events []Event) {
	v.events = events
	v.updateTable()
}

// SetSize sets the view dimensions
func (v *EventView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.table.SetHeight(height - 2) // Reserve space for borders
	v.updateTableColumns()
}

// SetFocused sets whether the event view is focused
func (v *EventView) SetFocused(focused bool) {
	v.table.Focus()
	if !focused {
		v.table.Blur()
	}
}

// updateTable updates the table with current events
func (v *EventView) updateTable() {
	rows := make([]table.Row, 0, len(v.events))
	
	// Sort events by timestamp (most recent first) and limit to recent events
	maxEvents := v.config.UI.PaneEventsHeight * 5 // Show more events than visible
	if maxEvents > len(v.events) {
		maxEvents = len(v.events)
	}
	
	for i := 0; i < maxEvents; i++ {
		event := v.events[i]
		row := v.createTableRow(event)
		rows = append(rows, row)
	}
	
	v.table.SetRows(rows)
}

// createTableRow creates a table row for an event
func (v *EventView) createTableRow(event Event) table.Row {
	// Use plain text without color styling to avoid display corruption
	eventType := event.Type
	if len(eventType) > 6 {
		eventType = eventType[:6]
	}
	
	// Format reason
	reason := event.Reason
	if len(reason) > 10 {
		reason = reason[:9] + "…"
	}
	
	// Format object
	object := event.Object
	if len(object) > 20 {
		object = object[:19] + "…"
	}
	
	// Format message (truncate if too long)
	message := event.Message
	maxMessageLength := 55
	if v.width > 0 {
		// Adjust message length based on available width
		usedWidth := 6 + 12 + 22 + 8 + 5 + 10 // Other columns + padding
		maxMessageLength = v.width - usedWidth
		if maxMessageLength < 20 {
			maxMessageLength = 20
		}
	}
	
	if len(message) > maxMessageLength {
		message = message[:maxMessageLength-3] + "…"
	}
	
	// Format timestamp
	timeFormatted := event.Timestamp
	
	// Format count
	countText := ""
	if event.Count > 1 {
		countText = fmt.Sprintf("%d", event.Count)
	}
	
	return table.Row{
		eventType,
		reason,
		object,
		message,
		timeFormatted,
		countText,
	}
}

// updateTableColumns updates table columns based on width
func (v *EventView) updateTableColumns() {
	baseColumns := []table.Column{
		{Title: "Type", Width: 6},
		{Title: "Reason", Width: 12},
		{Title: "Object", Width: 22},
		{Title: "Message", Width: 60},
		{Title: "Time", Width: 8},
		{Title: "Count", Width: 5},
	}

	// Adjust message column width based on available space
	if v.width > 0 {
		fixedWidth := 6 + 12 + 22 + 8 + 5 + 10 // Other columns + padding
		messageWidth := v.width - fixedWidth
		if messageWidth > 20 {
			baseColumns[3].Width = messageWidth
		}
	}

	v.table.SetColumns(baseColumns)
}
