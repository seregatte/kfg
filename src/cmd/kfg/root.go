package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/seregatte/kfg/src/internal/logger"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kfg",
	Short: "Declarative shell compiler",
	Long: `KFG is a declarative shell compiler that transforms YAML manifests into bash functions.

It allows you to define shell commands, their dependencies, and execution steps
in YAML manifests, then generates shell integration code that can be sourced
or used interactively.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Check if --verbose flag was explicitly provided
		if cmd.Flags().Changed("verbose") {
			verboseFlag, err := cmd.Flags().GetInt("verbose")
			if err == nil {
				// Set KFG_VERBOSE env var to flag value
				os.Setenv("KFG_VERBOSE", strconv.Itoa(verboseFlag))
				// Reinitialize logger with new verbose level
				logger.Reinitialize()
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, show help
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Silence cobra's default error handling so we can customize exit codes
	rootCmd.SilenceErrors = true

	if err := rootCmd.Execute(); err != nil {
		// Check if it's a usage error (flag/argument parsing error)
		if isUsageError(err) {
			logger.Error("cli", err.Error())
			os.Exit(2)
		} else {
			logger.Error("cli", err.Error())
			os.Exit(1)
		}
	}
}

// isUsageError checks if an error is a usage/flag/argument error.
func isUsageError(err error) bool {
	errMsg := err.Error()

	// Common usage error patterns from cobra
	usagePatterns := []string{
		"required flag",
		"unknown flag",
		"unknown command",
		"flag needs an argument",
		"accepts",
		"received",
		"arg(s)",
	}

	for _, pattern := range usagePatterns {
		if strings.Contains(errMsg, pattern) {
			return true
		}
	}

	return false
}

func init() {
	// Set version for Cobra's built-in --version flag
	rootCmd.Version = formatVersion()

	// Add persistent flags that are available to all subcommands
	rootCmd.PersistentFlags().IntP("verbose", "v", 0, "Verbosity level (0-3: 0=quiet, 1=info, 2=detail, 3=debug)")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}