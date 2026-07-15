package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/seregatte/kfg/src/internal/logger"
	"github.com/spf13/cobra"
)

var (
	// Supported agents for the wizard
	supportedAgents = map[string]bool{
		"pi":       true,
		"opencode": true,
	}
)

// aiCmd represents the ai command
var aiCmd = &cobra.Command{
	Use:   "ai [-- extra-args...]",
	Short: "Start an AI wizard for interactive configuration generation",
	Long: `Start an AI agent pre-configured to generate kfg configurations interactively.

The wizard guides you through creating project-specific kfg configurations by asking
questions about your project and composing manifests from existing domain building blocks.

The agent is determined by the KFG_AI_AGENT environment variable (default: pi).
All arguments after -- are forwarded to the agent as its initial prompt.

Environment variables:
  KFG_AI_AGENT  Agent to use for the wizard (default: pi)
                Supported values: pi, opencode

Examples:
  kfg ai                                          # Start wizard with pi (default)
  kfg ai -- "create a new project"                # Start wizard with initial prompt
  KFG_AI_AGENT=opencode kfg ai                    # Start wizard with opencode
  kfg ai -- --model sonnet "create project"       # Forward args to agent

The wizard will:
  1. Ask about your project type, language, and requirements
  2. Select appropriate building blocks from the kfg domain catalog
  3. Present a summary of what will be generated
  4. Ask for confirmation before writing any files
  5. Generate complete, project-specific kfg manifests`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Determine which agent to use
		agent := os.Getenv("KFG_AI_AGENT")
		if agent == "" {
			agent = "pi"
		}

		// Validate agent, fall back to pi if unsupported
		if !supportedAgents[agent] {
			logger.Warn("ai", fmt.Sprintf("Unsupported agent '%s', falling back to 'pi'", agent))
			agent = "pi"
		}

		// Get the current kfg binary path
		kfgBinary, err := os.Executable()
		if err != nil {
			logger.Error("ai", fmt.Sprintf("Failed to get executable path: %v", err))
			return err
		}

		// Construct the kfg run command
		// kfg run -k packages/domains/ai-agents/overlays/ai <agent> -- <args>
		runArgs := []string{
			"run",
			"-k", "packages/domains/ai-agents/overlays/ai",
			agent,
		}

		// Forward any extra arguments after --
		if len(args) > 0 {
			runArgs = append(runArgs, "--")
			runArgs = append(runArgs, args...)
		}

		logger.Debug("ai", fmt.Sprintf("Executing: %s %s", kfgBinary, strings.Join(runArgs, " ")))

		// Execute kfg run as a subprocess
		subprocess := exec.Command(kfgBinary, runArgs...)
		subprocess.Stdin = os.Stdin
		subprocess.Stdout = os.Stdout
		subprocess.Stderr = os.Stderr

		// Run the subprocess and propagate exit code
		if err := subprocess.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			logger.Error("ai", fmt.Sprintf("Failed to execute wizard: %v", err))
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(aiCmd)
}
