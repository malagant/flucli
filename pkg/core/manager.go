package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/malagant/fluxcli/internal/config"
	"github.com/malagant/fluxcli/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
)

// Manager manages FluxCD resources across multiple clusters
type Manager struct {
	config   *config.Config
	clusters map[string]*k8s.Client
	mu       sync.RWMutex
	
	// Event channels for UI updates
	resourceUpdates chan ResourceUpdate
	eventUpdates    chan EventUpdate
	errorUpdates    chan ErrorUpdate
	
	// Internal state
	currentCluster   string
	currentNamespace string
	ctx              context.Context
	cancel           context.CancelFunc
}

// ResourceUpdate represents a resource state update
type ResourceUpdate struct {
	Cluster   string
	Resources []k8s.Resource
	Type      k8s.ResourceType
}

// EventUpdate represents an event update
type EventUpdate struct {
	Cluster string
	Events  []corev1.Event
}

// ErrorUpdate represents an error update
type ErrorUpdate struct {
	Cluster string
	Error   error
}

// NewManager creates a new resource manager
func NewManager(cfg *config.Config) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Manager{
		config:          cfg,
		clusters:        make(map[string]*k8s.Client),
		resourceUpdates: make(chan ResourceUpdate, 100),
		eventUpdates:    make(chan EventUpdate, 100),
		errorUpdates:    make(chan ErrorUpdate, 100),
		currentCluster:  cfg.CurrentContext,
		currentNamespace: cfg.CurrentNamespace,
		ctx:             ctx,
		cancel:          cancel,
	}
}

// Start initializes the manager and starts background processes
func (m *Manager) Start() error {
	// Initialize default cluster connection
	if err := m.connectToCluster(m.currentCluster, m.config.CurrentKubeConfig, m.config.CurrentContext); err != nil {
		return fmt.Errorf("failed to connect to default cluster: %w", err)
	}

	// Initialize configured clusters
	for _, clusterCfg := range m.config.Clusters {
		if err := m.connectToCluster(clusterCfg.Name, clusterCfg.Kubeconfig, clusterCfg.Context); err != nil {
			m.errorUpdates <- ErrorUpdate{
				Cluster: clusterCfg.Name,
				Error:   fmt.Errorf("failed to connect to cluster %s: %w", clusterCfg.Name, err),
			}
		}
	}

	// Start background refresh
	go m.startResourceRefresh()
	go m.startEventRefresh()

	return nil
}

// Stop stops the manager and closes all connections
func (m *Manager) Stop() {
	m.cancel()
	close(m.resourceUpdates)
	close(m.eventUpdates)
	close(m.errorUpdates)
}

// connectToCluster establishes a connection to a Kubernetes cluster
func (m *Manager) connectToCluster(name, kubeconfig, context string) error {
	client, err := k8s.NewClient(kubeconfig, context, m.currentNamespace)
	if err != nil {
		return err
	}

	// Test connection
	if err := client.TestConnection(m.ctx); err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}

	m.mu.Lock()
	m.clusters[name] = client
	m.mu.Unlock()

	return nil
}

// GetResourceUpdates returns the channel for resource updates
func (m *Manager) GetResourceUpdates() <-chan ResourceUpdate {
	return m.resourceUpdates
}

// GetEventUpdates returns the channel for event updates
func (m *Manager) GetEventUpdates() <-chan EventUpdate {
	return m.eventUpdates
}

// GetErrorUpdates returns the channel for error updates
func (m *Manager) GetErrorUpdates() <-chan ErrorUpdate {
	return m.errorUpdates
}

// GetClusters returns the list of available clusters
func (m *Manager) GetClusters() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	clusters := make([]string, 0, len(m.clusters))
	for name := range m.clusters {
		clusters = append(clusters, name)
	}
	return clusters
}

// SetCurrentCluster sets the current active cluster
func (m *Manager) SetCurrentCluster(cluster string) error {
	m.mu.RLock()
	_, exists := m.clusters[cluster]
	m.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("cluster %s not found", cluster)
	}
	
	m.currentCluster = cluster
	return nil
}

// GetCurrentCluster returns the current active cluster
func (m *Manager) GetCurrentCluster() string {
	return m.currentCluster
}

// SetCurrentNamespace sets the current namespace
func (m *Manager) SetCurrentNamespace(namespace string) {
	m.currentNamespace = namespace
}

// GetCurrentNamespace returns the current namespace
func (m *Manager) GetCurrentNamespace() string {
	return m.currentNamespace
}

// ListResources lists all FluxCD resources of a specific type
func (m *Manager) ListResources(resourceType k8s.ResourceType) ([]k8s.Resource, error) {
	m.mu.RLock()
	client, exists := m.clusters[m.currentCluster]
	m.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("cluster %s not connected", m.currentCluster)
	}

	ctx, cancel := context.WithTimeout(m.ctx, 10*time.Second)
	defer cancel()

	switch resourceType {
	case k8s.ResourceTypeGitRepository:
		return client.ListGitRepositories(ctx, m.currentNamespace)
	case k8s.ResourceTypeHelmRepository:
		return client.ListHelmRepositories(ctx, m.currentNamespace)
	case k8s.ResourceTypeKustomization:
		return client.ListKustomizations(ctx, m.currentNamespace)
	case k8s.ResourceTypeHelmRelease:
		return client.ListHelmReleases(ctx, m.currentNamespace)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

// SuspendResource suspends a FluxCD resource
func (m *Manager) SuspendResource(resourceType k8s.ResourceType, name string) error {
	m.mu.RLock()
	client, exists := m.clusters[m.currentCluster]
	m.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("cluster %s not connected", m.currentCluster)
	}

	ctx, cancel := context.WithTimeout(m.ctx, 10*time.Second)
	defer cancel()

	return client.SuspendResource(ctx, resourceType, name, m.currentNamespace)
}

// ResumeResource resumes a FluxCD resource
func (m *Manager) ResumeResource(resourceType k8s.ResourceType, name string) error {
	m.mu.RLock()
	client, exists := m.clusters[m.currentCluster]
	m.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("cluster %s not connected", m.currentCluster)
	}

	ctx, cancel := context.WithTimeout(m.ctx, 10*time.Second)
	defer cancel()

	return client.ResumeResource(ctx, resourceType, name, m.currentNamespace)
}

// ReconcileResource triggers reconciliation of a FluxCD resource
func (m *Manager) ReconcileResource(resourceType k8s.ResourceType, name string) error {
	m.mu.RLock()
	client, exists := m.clusters[m.currentCluster]
	m.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("cluster %s not connected", m.currentCluster)
	}

	ctx, cancel := context.WithTimeout(m.ctx, 10*time.Second)
	defer cancel()

	return client.ReconcileResource(ctx, resourceType, name, m.currentNamespace)
}

// startResourceRefresh starts the background resource refresh process
func (m *Manager) startResourceRefresh() {
	ticker := time.NewTicker(m.config.Defaults.RefreshInterval)
	defer ticker.Stop()

	resourceTypes := []k8s.ResourceType{
		k8s.ResourceTypeGitRepository,
		k8s.ResourceTypeHelmRepository,
		k8s.ResourceTypeKustomization,
		k8s.ResourceTypeHelmRelease,
	}

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.refreshResources(resourceTypes)
		}
	}
}

// refreshResources refreshes all resources for all clusters
func (m *Manager) refreshResources(resourceTypes []k8s.ResourceType) {
	m.mu.RLock()
	clusters := make(map[string]*k8s.Client)
	for name, client := range m.clusters {
		clusters[name] = client
	}
	m.mu.RUnlock()

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, m.config.Defaults.MaxConcurrentClusters)

	for clusterName, client := range clusters {
		wg.Add(1)
		go func(name string, c *k8s.Client) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			for _, resourceType := range resourceTypes {
				resources, err := m.listResourcesForCluster(c, resourceType)
				if err != nil {
					m.errorUpdates <- ErrorUpdate{
						Cluster: name,
						Error:   fmt.Errorf("failed to list %s: %w", resourceType, err),
					}
					continue
				}

				select {
				case m.resourceUpdates <- ResourceUpdate{
					Cluster:   name,
					Resources: resources,
					Type:      resourceType,
				}:
				case <-m.ctx.Done():
					return
				}
			}
		}(clusterName, client)
	}

	wg.Wait()
}

// listResourcesForCluster lists resources for a specific cluster and type
func (m *Manager) listResourcesForCluster(client *k8s.Client, resourceType k8s.ResourceType) ([]k8s.Resource, error) {
	ctx, cancel := context.WithTimeout(m.ctx, 10*time.Second)
	defer cancel()

	switch resourceType {
	case k8s.ResourceTypeGitRepository:
		return client.ListGitRepositories(ctx, "")
	case k8s.ResourceTypeHelmRepository:
		return client.ListHelmRepositories(ctx, "")
	case k8s.ResourceTypeKustomization:
		return client.ListKustomizations(ctx, "")
	case k8s.ResourceTypeHelmRelease:
		return client.ListHelmReleases(ctx, "")
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

// startEventRefresh starts the background event refresh process
func (m *Manager) startEventRefresh() {
	if !m.config.Defaults.EventsEnabled {
		return
	}

	ticker := time.NewTicker(2 * time.Second) // Events refresh more frequently
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.refreshEvents()
		}
	}
}

// refreshEvents refreshes events for all clusters
func (m *Manager) refreshEvents() {
	m.mu.RLock()
	clusters := make(map[string]*k8s.Client)
	for name, client := range m.clusters {
		clusters[name] = client
	}
	m.mu.RUnlock()

	for clusterName, client := range clusters {
		go func(name string, c *k8s.Client) {
			ctx, cancel := context.WithTimeout(m.ctx, 5*time.Second)
			defer cancel()

			events, err := c.GetEvents(ctx, "")
			if err != nil {
				m.errorUpdates <- ErrorUpdate{
					Cluster: name,
					Error:   fmt.Errorf("failed to get events: %w", err),
				}
				return
			}

			select {
			case m.eventUpdates <- EventUpdate{
				Cluster: name,
				Events:  events,
			}:
			case <-m.ctx.Done():
				return
			}
		}(clusterName, client)
	}
}
