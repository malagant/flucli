package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/malagant/fluxcli/internal/config"
	"github.com/malagant/fluxcli/pkg/ui"
)

var (
	cfgFile     string
	kubeconfig  string
	context     string
	namespace   string
	debug       bool
	logLevel    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fluxcli",
	Short: "Terminal UI for FluxCD Multi-Cluster Management",
	Long: `FluxCLI is a powerful terminal user interface (TUI) for managing FluxCD resources 
across multiple Kubernetes clusters. Inspired by tools like K9s, FluxCLI provides an 
intuitive, keyboard-driven interface specifically designed for GitOps workflows.

Features:
- Multi-Cluster Support - Seamlessly switch between and manage multiple Kubernetes clusters
- FluxCD Resource Management - View, monitor, and operate on GitRepository, HelmRepository, 
  Kustomization, HelmRelease, and ResourceSet resources  
- Real-time Monitoring - Live updates of resource status, events, and reconciliation progress
- Intuitive Navigation - K9s-inspired keyboard shortcuts and command patterns
- Advanced Filtering - Filter resources by namespace, status, cluster, and custom criteria
- Event Streaming - Monitor FluxCD events and reconciliation status in real-time`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration
		cfg, err := config.Load(cfgFile, kubeconfig, context, namespace)
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		// Initialize and run the TUI
		app := ui.NewApp(cfg)
		if err := app.Run(); err != nil {
			return fmt.Errorf("failed to run application: %w", err)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fluxcli/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig file (default is $HOME/.kube/config)")
	rootCmd.PersistentFlags().StringVar(&context, "context", "", "kubernetes context to use")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "kubernetes namespace to use")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug mode")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level (trace, debug, info, warn, error)")

	// Bind flags to viper
	viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))
	viper.BindPFlag("context", rootCmd.PersistentFlags().Lookup("context"))
	viper.BindPFlag("namespace", rootCmd.PersistentFlags().Lookup("namespace"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".fluxcli" (without extension).
		viper.AddConfigPath(home + "/.fluxcli")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if debug {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}
