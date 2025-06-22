package main

import (
	"fmt"
	"os"

	"github.com/malagant/fluxcli/cmd"
)

var (
	version   = "dev"
	commit    = "none"
	buildTime = "unknown"
)

func main() {
	// Set version information in cmd package
	cmd.SetVersionInfo(version, commit, buildTime)
	
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
