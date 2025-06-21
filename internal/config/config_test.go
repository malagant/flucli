package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadDefaultConfig(t *testing.T) {
	// Test loading without any config file
	config, err := Load("", "", "", "")
	require.NoError(t, err)

	// Should have default values
	assert.Equal(t, "flux-system", config.Defaults.Namespace)
	assert.Equal(t, 5*time.Second, config.Defaults.RefreshInterval)
	assert.Equal(t, 10, config.Defaults.MaxConcurrentClusters)
	assert.True(t, config.Defaults.EventsEnabled)
	assert.Equal(t, "dark", config.UI.Theme)
	assert.True(t, config.UI.ShowAge)
	assert.True(t, config.UI.ShowMessage)
	assert.True(t, config.UI.ShowNamespace)
	assert.Equal(t, 4, config.UI.PaneEventsHeight)
	assert.Equal(t, 30, config.UI.ColumnsName)
	assert.Equal(t, 15, config.UI.ColumnsStatus)
}

func TestLoadWithCommandLineOverrides(t *testing.T) {
	// Load config with command line overrides
	config, err := Load("", "/custom/kubeconfig", "custom-context", "custom-namespace")
	require.NoError(t, err)

	assert.Equal(t, "/custom/kubeconfig", config.CurrentKubeConfig)
	assert.Equal(t, "custom-context", config.CurrentContext)
	assert.Equal(t, "custom-namespace", config.CurrentNamespace)
}

func TestAddCluster(t *testing.T) {
	config, err := Load("", "", "", "")
	require.NoError(t, err)
	
	cluster := ClusterConfig{
		Name:       "new-cluster",
		Kubeconfig: "/path/to/kubeconfig",
		Context:    "new-context",
		Namespace:  "new-namespace",
	}

	config.AddCluster(cluster)
	
	assert.Len(t, config.Clusters, 1)
	assert.Equal(t, "new-cluster", config.Clusters[0].Name)
	assert.Equal(t, "/path/to/kubeconfig", config.Clusters[0].Kubeconfig)
	assert.Equal(t, "new-context", config.Clusters[0].Context)
	assert.Equal(t, "new-namespace", config.Clusters[0].Namespace)
}

func TestAddClusterReplaceExisting(t *testing.T) {
	config, err := Load("", "", "", "")
	require.NoError(t, err)
	
	// Add initial cluster
	cluster1 := ClusterConfig{
		Name:       "test-cluster",
		Kubeconfig: "/path1",
		Context:    "ctx1",
	}
	config.AddCluster(cluster1)
	assert.Len(t, config.Clusters, 1)
	
	// Add cluster with same name (should replace)
	cluster2 := ClusterConfig{
		Name:       "test-cluster",
		Kubeconfig: "/path2",
		Context:    "ctx2",
	}
	config.AddCluster(cluster2)
	
	assert.Len(t, config.Clusters, 1)
	assert.Equal(t, "/path2", config.Clusters[0].Kubeconfig)
	assert.Equal(t, "ctx2", config.Clusters[0].Context)
}

func TestRemoveCluster(t *testing.T) {
	config, err := Load("", "", "", "")
	require.NoError(t, err)
	
	// Add two clusters
	config.AddCluster(ClusterConfig{Name: "cluster1", Kubeconfig: "/path1", Context: "ctx1"})
	config.AddCluster(ClusterConfig{Name: "cluster2", Kubeconfig: "/path2", Context: "ctx2"})
	
	assert.Len(t, config.Clusters, 2)
	
	// Remove one cluster
	removed := config.RemoveCluster("cluster1")
	assert.True(t, removed)
	assert.Len(t, config.Clusters, 1)
	assert.Equal(t, "cluster2", config.Clusters[0].Name)
	
	// Try to remove non-existent cluster
	removed = config.RemoveCluster("non-existent")
	assert.False(t, removed)
	assert.Len(t, config.Clusters, 1)
}
