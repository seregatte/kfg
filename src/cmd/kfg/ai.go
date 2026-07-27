package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/seregatte/kfg/src/internal/logger"
	"github.com/spf13/cobra"
)

const (
	// localAIOverlay is the fixed working-directory-relative path for the AI overlay.
	localAIOverlay = "packages/domains/ai-agents/overlays/ai"
	// remoteAIOverlay is the default online source used outside kfg checkouts.
	remoteAIOverlay = "https://github.com/seregatte/kfg.git//packages/domains/ai-agents/overlays/ai?ref=main"
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

Overlay source resolution:
  - When run from a kfg checkout root, the local AI overlay is used directly.
  - Otherwise, the canonical online overlay is fetched from GitHub.
  - The KFG_KPATH environment variable does not affect kfg ai.

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
		agent := resolveAgent()
		kfgBinary, err := os.Executable()
		if err != nil {
			logger.Error("ai", fmt.Sprintf("Failed to get executable path: %v", err))
			return err
		}

		overlay, err := resolveAIOverlay()
		if err != nil {
			logger.Error("ai", fmt.Sprintf("Failed to resolve AI overlay: %v", err))
			return err
		}

		runArgs := buildRunArgs(overlay, agent, args)
		logger.Debug("ai", fmt.Sprintf("Executing: %s %s", kfgBinary, strings.Join(runArgs, " ")))

		subprocess := exec.Command(kfgBinary, runArgs...)
		subprocess.Stdin = os.Stdin
		subprocess.Stdout = os.Stdout
		subprocess.Stderr = os.Stderr

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

// resolveAgent returns the wizard agent from KFG_AI_AGENT, defaulting to pi.
func resolveAgent() string {
	agent := os.Getenv("KFG_AI_AGENT")
	if agent == "" {
		agent = "pi"
	}
	if !supportedAgents[agent] {
		logger.Warn("ai", fmt.Sprintf("Unsupported agent '%s', falling back to 'pi'", agent))
		agent = "pi"
	}
	return agent
}

// resolveAIOverlay selects the local AI overlay when it exists relative to the
// current working directory, or the canonical remote source when absent.
// Filesystem errors other than absence are returned so development defects are
// not silently hidden by a remote fallback.
func resolveAIOverlay() (string, error) {
	_, err := os.Stat(localAIOverlay)
	if err == nil {
		return localAIOverlay, nil
	}
	if os.IsNotExist(err) {
		return remoteAIOverlay, nil
	}
	return "", fmt.Errorf("cannot inspect %s: %w", localAIOverlay, err)
}

// buildRunArgs constructs the kfg run command arguments.
// The overlay is passed explicitly via -k so KFG_KPATH cannot redirect the wizard.
func buildRunArgs(overlay, agent string, args []string) []string {
	runArgs := []string{
		"run",
		"-k", overlay,
		agent,
	}
	if len(args) > 0 {
		runArgs = append(runArgs, "--")
		runArgs = append(runArgs, args...)
	}
	return runArgs
}

func init() {
	rootCmd.AddCommand(aiCmd)
}
