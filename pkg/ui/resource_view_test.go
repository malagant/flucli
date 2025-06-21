package ui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	
	"github.com/malagant/fluxcli/internal/config"
	"github.com/malagant/fluxcli/pkg/k8s"
)

func TestNewResourceView(t *testing.T) {
	cfg, err := config.Load("", "", "", "")
	require.NoError(t, err)
	
	rv := NewResourceView(cfg)
	
	assert.NotNil(t, rv)
	assert.NotNil(t, rv.table)
}

func TestResourceView_Update(t *testing.T) {
	cfg, err := config.Load("", "", "", "")
	require.NoError(t, err)
	
	rv := NewResourceView(cfg)
	
	// Test key press
	msg := tea.KeyMsg{Type: tea.KeyDown}
	model, cmd := rv.Update(msg)
	
	assert.NotNil(t, model)
	assert.Nil(t, cmd) // Navigation doesn't produce commands
}

func TestResourceView_SetResources(t *testing.T) {
	cfg, err := config.Load("", "", "", "")
	require.NoError(t, err)
	
	rv := NewResourceView(cfg)
	
	// Create test resources
	resources := []k8s.Resource{
		createTestResource("test-repo", "default", k8s.ResourceTypeGitRepository),
		createTestResource("another-repo", "flux-system", k8s.ResourceTypeGitRepository),
	}
	
	rv.SetResources(resources)
	
	// Check that table has been updated
	assert.Equal(t, 2, len(rv.table.Rows()))
}

func TestResourceView_SetResourceType(t *testing.T) {
	cfg, err := config.Load("", "", "", "")
	require.NoError(t, err)
	
	rv := NewResourceView(cfg)
	
	rv.SetResourceType(k8s.ResourceTypeHelmRepository)
	assert.Equal(t, k8s.ResourceTypeHelmRepository, rv.resourceType)
	
	rv.SetResourceType(k8s.ResourceTypeKustomization)
	assert.Equal(t, k8s.ResourceTypeKustomization, rv.resourceType)
	
	rv.SetResourceType(k8s.ResourceTypeHelmRelease)
	assert.Equal(t, k8s.ResourceTypeHelmRelease, rv.resourceType)
}

func TestResourceView_SetSize(t *testing.T) {
	cfg, err := config.Load("", "", "", "")
	require.NoError(t, err)
	
	rv := NewResourceView(cfg)
	
	// Test resize
	rv.SetSize(100, 30)
	
	// Table should be resized (we can't easily test exact dimensions 
	// without access to internal table state, but we can ensure 
	// the method doesn't panic)
	assert.NotNil(t, rv.table)
}

func TestResourceView_GetSelectedResource(t *testing.T) {
	cfg, err := config.Load("", "", "", "")
	require.NoError(t, err)
	
	rv := NewResourceView(cfg)
	
	// With no resources, should return nil
	selected := rv.GetSelectedResource()
	assert.Nil(t, selected)
	
	// Add some resources
	resources := []k8s.Resource{
		createTestResource("test-repo", "default", k8s.ResourceTypeGitRepository),
		createTestResource("another-repo", "flux-system", k8s.ResourceTypeGitRepository),
	}
	rv.SetResources(resources)
	
	// Should return first resource (index 0)
	selected = rv.GetSelectedResource()
	assert.NotNil(t, selected)
	assert.Equal(t, "test-repo", selected.Name)
}

// Helper function to create test Resource
func createTestResource(name, namespace string, resourceType k8s.ResourceType) k8s.Resource {
	return k8s.Resource{
		Type:       resourceType,
		Name:       name,
		Namespace:  namespace,
		Ready:      true,
		Status:     "Ready",
		Message:    "Reconciliation succeeded",
		Age:        5 * time.Minute,
		LastUpdate: time.Now(),
		Conditions: []k8s.Condition{
			{
				Type:               "Ready",
				Status:             "True",
				Reason:             "ReconciliationSucceeded", 
				Message:            "Reconciliation succeeded",
				LastTransitionTime: time.Now(),
			},
		},
		Suspended: false,
		URL:       "https://github.com/example/repo",
		Revision:  "main@sha256:abc123",
	}
}
