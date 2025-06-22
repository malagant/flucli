package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"k8s.io/client-go/util/homedir"
)

// Config represents the application configuration
type Config struct {
	Clusters         []ClusterConfig `yaml:"clusters"`
	Defaults         DefaultConfig   `yaml:"defaults"`
	UI               UIConfig        `yaml:"ui"`
	Debug            bool            `yaml:"debug"`
	LogLevel         string          `yaml:"log_level"`
	CurrentKubeConfig string         `yaml:"-"` // Runtime only
	CurrentContext   string          `yaml:"-"` // Runtime only
	CurrentNamespace string          `yaml:"-"` // Runtime only
}

// ClusterConfig represents a single cluster configuration
type ClusterConfig struct {
	Name        string `yaml:"name"`
	Context     string `yaml:"context"`
	Kubeconfig  string `yaml:"kubeconfig"`
	Namespace   string `yaml:"namespace"`
	Color       string `yaml:"color"`
	Description string `yaml:"description"`
}

// DefaultConfig represents default settings
type DefaultConfig struct {
	Namespace            string        `yaml:"namespace"`
	RefreshInterval      time.Duration `yaml:"refresh_interval"`
	MaxConcurrentClusters int          `yaml:"max_concurrent_clusters"`
	EventsEnabled        bool          `yaml:"events_enabled"`
}

// UIConfig represents UI-specific settings
type UIConfig struct {
	Theme           string `yaml:"theme"`
	ShowAge         bool   `yaml:"show_age"`
	ShowMessage     bool   `yaml:"show_message"`
	ShowNamespace   bool   `yaml:"show_namespace"`
	PaneEventsHeight int   `yaml:"pane_events_height"`
	ColumnsName     int    `yaml:"columns_name"`
	ColumnsStatus   int    `yaml:"columns_status"`
}

// Load loads configuration from file and command line arguments
func Load(configFile, kubeconfig, context, namespace string) (*Config, error) {
	cfg := &Config{
		Defaults: DefaultConfig{
			Namespace:            "flux-system",
			RefreshInterval:      5 * time.Second,
			MaxConcurrentClusters: 10,
			EventsEnabled:        true,
		},
		UI: UIConfig{
			Theme:           "dark",
			ShowAge:         true,
			ShowMessage:     true,
			ShowNamespace:   true,
			PaneEventsHeight: 4,
			ColumnsName:     30,
			ColumnsStatus:   15,
		},
		Debug:    viper.GetBool("debug"),
		LogLevel: viper.GetString("log-level"),
	}

	// Load from config file if it exists
	if err := loadConfigFile(cfg, configFile); err != nil {
		return nil, err
	}

	// Override with command line arguments
	if kubeconfig != "" {
		cfg.CurrentKubeConfig = kubeconfig
	} else if cfg.CurrentKubeConfig == "" {
		// Check KUBECONFIG environment variable first
		if kubeconfigEnv := os.Getenv("KUBECONFIG"); kubeconfigEnv != "" {
			cfg.CurrentKubeConfig = kubeconfigEnv
		} else {
			// Use default kubeconfig location as final fallback
			if home := homedir.HomeDir(); home != "" {
				cfg.CurrentKubeConfig = filepath.Join(home, ".kube", "config")
			}
		}
	}

	if context != "" {
		cfg.CurrentContext = context
	}

	if namespace != "" {
		cfg.CurrentNamespace = namespace
	} else if cfg.CurrentNamespace == "" {
		cfg.CurrentNamespace = cfg.Defaults.Namespace
	}

	return cfg, nil
}

// loadConfigFile loads configuration from file
func loadConfigFile(cfg *Config, configFile string) error {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		configDir := filepath.Join(home, ".fluxcli")
		configPath := filepath.Join(configDir, "config.yaml")

		// Create config directory if it doesn't exist
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		// Create default config file if it doesn't exist
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			if err := createDefaultConfig(configPath); err != nil {
				return fmt.Errorf("failed to create default config: %w", err)
			}
		}

		viper.SetConfigFile(configPath)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	return viper.Unmarshal(cfg)
}

// createDefaultConfig creates a default configuration file
func createDefaultConfig(path string) error {
	defaultConfig := `# FluxCLI Configuration

clusters: []

defaults:
  namespace: flux-system
  refresh_interval: 5s
  max_concurrent_clusters: 10
  events_enabled: true

ui:
  theme: dark
  show_age: true
  show_message: true
  show_namespace: true
  pane_events_height: 4
  columns_name: 30
  columns_status: 15
`

	return os.WriteFile(path, []byte(defaultConfig), 0644)
}

// GetCluster returns cluster configuration by name
func (c *Config) GetCluster(name string) (*ClusterConfig, bool) {
	for _, cluster := range c.Clusters {
		if cluster.Name == name {
			return &cluster, true
		}
	}
	return nil, false
}

// AddCluster adds a new cluster configuration
func (c *Config) AddCluster(cluster ClusterConfig) {
	// Remove existing cluster with same name
	for i, existing := range c.Clusters {
		if existing.Name == cluster.Name {
			c.Clusters[i] = cluster
			return
		}
	}
	// Add new cluster
	c.Clusters = append(c.Clusters, cluster)
}

// RemoveCluster removes a cluster configuration by name
func (c *Config) RemoveCluster(name string) bool {
	for i, cluster := range c.Clusters {
		if cluster.Name == name {
			c.Clusters = append(c.Clusters[:i], c.Clusters[i+1:]...)
			return true
		}
	}
	return false
}

// Save saves the configuration to file
func (c *Config) Save() error {
	return viper.WriteConfig()
}

// SaveTo saves the configuration to the specified file path
func (c *Config) SaveTo(filepath string) error {
	return viper.WriteConfigAs(filepath)
}
