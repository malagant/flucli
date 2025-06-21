package ui

import (
	"fmt"
	"strings"

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
		{Title: "Type", Width: 8},
		{Title: "Reason", Width: 15},
		{Title: "Object", Width: 25},
		{Title: "Message", Width: 50},
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
		switch msg.String() {
		case "j", "down":
			v.table, cmd = v.table.Update(msg)
		case "k", "up":
			v.table, cmd = v.table.Update(msg)
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
	// Format event type with color
	eventType := event.Type
	typeColor := lipgloss.Color("244") // Default gray
	switch strings.ToLower(event.Type) {
	case "normal":
		typeColor = lipgloss.Color("46") // Green
	case "warning":
		typeColor = lipgloss.Color("226") // Yellow
	case "error":
		typeColor = lipgloss.Color("196") // Red
	}
	typeFormatted := lipgloss.NewStyle().Foreground(typeColor).Render(eventType)
	
	// Format reason
	reason := event.Reason
	if len(reason) > 12 {
		reason = reason[:12] + "..."
	}
	reasonFormatted := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render(reason)
	
	// Format object
	object := event.Object
	if len(object) > 22 {
		object = object[:22] + "..."
	}
	objectFormatted := lipgloss.NewStyle().Foreground(lipgloss.Color("117")).Render(object)
	
	// Format message (truncate if too long)
	message := event.Message
	maxMessageLength := 45
	if v.width > 0 {
		// Adjust message length based on available width
		usedWidth := 8 + 15 + 25 + 8 + 5 + 10 // Other columns + padding
		maxMessageLength = v.width - usedWidth
		if maxMessageLength < 20 {
			maxMessageLength = 20
		}
	}
	
	if len(message) > maxMessageLength {
		message = message[:maxMessageLength-3] + "..."
	}
	messageColor := lipgloss.Color("244")
	if strings.ToLower(event.Type) == "warning" {
		messageColor = lipgloss.Color("226")
	} else if strings.ToLower(event.Type) == "error" {
		messageColor = lipgloss.Color("196")
	}
	messageFormatted := lipgloss.NewStyle().Foreground(messageColor).Render(message)
	
	// Format timestamp
	timeFormatted := lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Render(event.Timestamp)
	
	// Format count
	countText := ""
	if event.Count > 1 {
		countText = fmt.Sprintf("%d", event.Count)
	}
	countFormatted := lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Render(countText)
	
	return table.Row{
		typeFormatted,
		reasonFormatted,
		objectFormatted,
		messageFormatted,
		timeFormatted,
		countFormatted,
	}
}

// updateTableColumns updates table columns based on width
func (v *EventView) updateTableColumns() {
	baseColumns := []table.Column{
		{Title: "Type", Width: 8},
		{Title: "Reason", Width: 15},
		{Title: "Object", Width: 25},
		{Title: "Message", Width: 45},
		{Title: "Time", Width: 8},
		{Title: "Count", Width: 5},
	}

	// Adjust message column width based on available space
	if v.width > 0 {
		fixedWidth := 8 + 15 + 25 + 8 + 5 + 10 // Other columns + padding
		messageWidth := v.width - fixedWidth
		if messageWidth > 20 {
			baseColumns[3].Width = messageWidth
		}
	}

	v.table.SetColumns(baseColumns)
}
