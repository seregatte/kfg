package main

import (
	"github.com/spf13/cobra"
)

// sysCmd represents the sys command group for internal infrastructure commands
var sysCmd = &cobra.Command{
	Use:   "sys",
	Short: "Internal system commands",
	Long: `Internal system commands for infrastructure and debugging.

These commands are intended for internal use and are not typically
needed for normal operation.

Subcommands:
  log    Write structured log entries
  cache  Cache operations for Step results

Examples:
  kfg sys log info "component" "message"
  kfg sys log error "cmd:build" "failed to parse manifest"
  kfg sys cache ls
  kfg sys cache exists ctx7.steps.install`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	// Add sysCmd to rootCmd
	rootCmd.AddCommand(sysCmd)
}
