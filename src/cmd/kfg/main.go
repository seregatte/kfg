package main

import (
	"os"

	"github.com/seregatte/kfg/src/internal/config"
	"github.com/seregatte/kfg/src/internal/logger"
)

// Build-time version information (injected via ldflags)
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// Initialize configuration first
	if err := config.Initialize(); err != nil {
		os.Exit(1)
	}

	// Initialize logger (reads directly from environment)
	if err := logger.Initialize(); err != nil {
		os.Exit(1)
	}
	defer logger.Close()

	// Execute the root command
	Execute()
}