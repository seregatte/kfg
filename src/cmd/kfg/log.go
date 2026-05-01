package main

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/seregatte/kfg/src/internal/logger"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log <level> <component> [message...]",
	Short: "Write structured log entry",
	Long: `Write a structured log entry to JSONL file and optionally stderr.

The log entry is persisted to JSONL file at the configured location.
Human-readable output appears in stderr based on KFG_VERBOSE level.

Levels: error, warn, info, detail, debug

Arguments:
  level      - Log level (required): error, warn, info, detail, debug
  component  - Component name (required): identifies the source of the log
  message    - Log message (optional): the message content

Flags:
  --session-id - Session ID for log correlation (overrides KFG_SESSION_ID env var)
                 Empty string omits session_id field

Session IDs enable per-invocation log correlation. Each generated command wrapper
auto-generates a session ID at invocation start (format: timestamp-random).

Examples:
  kfg sys log info "feature:mcps" "syncing for claude"
  kfg sys log error "cmd:build" "failed to parse manifest"
  kfg sys log debug "store:push" "artifact already exists"
  kfg sys log detail "resolve:workflow" ""

  # Set session ID via flag
  kfg sys log --session-id "custom-123" info "test:comp" "test message"

  # Omit session_id field
  kfg sys log --session-id "" info "test:comp" "test message"

JSONL output includes session_id field when set:
  {"level":"info","component":"test:comp","msg":"test message","session_id":"custom-123",...}`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Parse arguments
		level := strings.ToLower(args[0])
		component := args[1]
		message := ""
		if len(args) > 2 {
			// Join remaining args as message
			message = strings.Join(args[2:], " ")
		}

		// Get source flag (defaults to "shell" since kfg sys log is typically called from shell)
		source, _ := cmd.Flags().GetString("source")

		// Get session-id flag
		// If flag was not provided, pass nil to use env var
		// If flag was provided (even empty), pass pointer to override env var
		var sessionID *string
		if cmd.Flags().Changed("session-id") {
			sid, _ := cmd.Flags().GetString("session-id")
			sessionID = &sid
		}

		// Validate level
		validLevels := []string{"error", "warn", "info", "detail", "debug"}
		valid := false
		for _, validLevel := range validLevels {
			if level == validLevel {
				valid = true
				break
			}
		}

		if !valid {
			logger.Error("log", "Levels: error, warn, info, detail, debug")
			os.Exit(1)
		}

		// Log the entry with custom source and optional session ID
		logger.LogWithSource(level, source, component, message, sessionID)
	},
}

func init() {
	// Add logCmd to sysCmd (sys group for internal commands)
	sysCmd.AddCommand(logCmd)

	// Add --source flag (hidden, defaults to "shell")
	// This allows shell helpers to log with source="shell"
	logCmd.Flags().String("source", "shell", "Source identifier for the log entry")
	logCmd.Flags().MarkHidden("source")

	// Add --session-id flag for per-invocation log correlation
	logCmd.Flags().String("session-id", "", "Session ID for log correlation (overrides KFG_SESSION_ID env var)")
}