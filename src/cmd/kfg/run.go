package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/seregatte/kfg/src/internal/generate"
	"github.com/seregatte/kfg/src/internal/logger"
	"github.com/seregatte/kfg/src/internal/manifest"
	"github.com/seregatte/kfg/src/internal/resolve"
)

var (
	// Flags for the run command
	runKustomizePath string
	runFile          string
	runWorkflow      string
	runCmds          string
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [agent] [-- extra-args...]",
	Short: "Run an agent with one-shot execution",
	Long: `Run an agent by generating shell code, sourcing it, and executing in one invocation.

This command provides a one-shot "generate → source → execute" experience
for running agents. It matches agents by their commandName (e.g., "claude")
and auto-detects the workflow if not specified.

Arguments after '--' are passed directly to the agent.

Examples:
  kfg run -k .nixai/overlay/dev claude
  kfg run -k .nixai/overlay/dev claude -- --model gpt-4
  kfg run -k .nixai/overlay/dev -w dev claude
  kfg run -f manifest.yaml claude
  kfg run -k .nixai/overlay/dev (lists available agents)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate flags
		if runKustomizePath == "" && runFile == "" {
			logger.Error("run", "Either -k (kustomize path) or -f (file) is required")
			cmd.Help()
			os.Exit(2)
		}

		if runKustomizePath != "" && runFile != "" {
			logger.Error("run", "Cannot use both -k and -f flags")
			os.Exit(2)
		}

		// Parse run args: agent name and extra args
		agentName, extraArgs := parseLaunchArgs(cmd, args)

		// No agent name provided - list available agents
		if agentName == "" {
			// Run the apply pipeline to get the index
			result, err := runApplyPipeline(runKustomizePath, runFile)
			if err != nil {
				logger.Error("run", err.Error())
				os.Exit(1)
			}
			listAvailableAgents(result.Index)
			os.Exit(1)
		}

		// Run the apply pipeline
		result, err := runApplyPipeline(runKustomizePath, runFile)
		if err != nil {
			logger.Error("run", err.Error())
			os.Exit(1)
		}

		// Parse workflow filter
		var workflowFilter string
		if runWorkflow != "" {
			workflowFilter = runWorkflow
		}

		// Find the agent in the index
		cmdMetadataName, workflowName, foundCmd, err := findAgent(result.Index, agentName, workflowFilter)
		if err != nil {
			logger.Error("run", err.Error())
			listAvailableAgents(result.Index)
			os.Exit(1)
		}

		// Generate shell code for the specific agent
		shellCode, shellType, err := generateForAgent(result, workflowName, cmdMetadataName, foundCmd)
		if err != nil {
			logger.Error("run", fmt.Sprintf("Failed to generate shell code: %v", err))
			os.Exit(1)
		}

		// Execute the agent
		executeAgent(shellCode, shellType, agentName, extraArgs)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Add flags for the run command
	runCmd.Flags().StringVarP(&runKustomizePath, "kustomize", "k", "", "Kustomization directory path")
	runCmd.Flags().StringVarP(&runFile, "file", "f", "", "Manifest file path (use '-' for stdin)")
	runCmd.Flags().StringVarP(&runWorkflow, "workflow", "w", "", "CmdWorkflow name (auto-detected if not specified)")
	runCmd.Flags().StringVarP(&runCmds, "cmds", "c", "", "Comma-separated list of cmds (overrides agent matching)")
}

// parseLaunchArgs splits args at '--' into agent name and extra args.
// Returns (agentName, extraArgs).
func parseLaunchArgs(cmd *cobra.Command, args []string) (string, []string) {
	// Find the separator '--'
	separatorIndex := -1
	for i, arg := range args {
		if arg == "--" {
			separatorIndex = i
			break
		}
	}

	if separatorIndex == -1 {
		// No separator found - first arg is agent name, no extra args
		if len(args) > 0 {
			return args[0], []string{}
		}
		return "", []string{}
	}

	// Separator found - first arg before separator is agent name, rest are extra args
	if separatorIndex > 0 {
		return args[0], args[separatorIndex+1:]
	}
	return "", args[separatorIndex+1:]
}

// findAgent matches an agent by commandName and finds its workflow.
// Returns (cmdMetadataName, workflowName, foundCmd, error).
func findAgent(index *resolve.Index, agentName string, workflowFilter string) (string, string, *manifest.Cmd, error) {
	cmds := index.GetCmds()

	// Find cmd by commandName
	var foundCmd *manifest.Cmd
	var cmdMetadataName string
	for _, cmd := range cmds {
		if cmd.Metadata.CommandName == agentName {
			foundCmd = cmd
			cmdMetadataName = cmd.Metadata.Name
			break
		}
	}

	if foundCmd == nil {
		return "", "", nil, fmt.Errorf("agent '%s' not found", agentName)
	}

	// Find workflow containing this cmd
	workflows := index.GetCmdWorkflows()
	var workflowName string

	for _, wf := range workflows {
		// Check if workflow filter is specified
		if workflowFilter != "" && wf.Metadata.Name != workflowFilter {
			continue
		}

		// Check if this workflow contains the cmd
		for _, cmdRef := range wf.Spec.Cmds {
			if cmdRef == cmdMetadataName {
				workflowName = wf.Metadata.Name
				break
			}
		}

		if workflowName != "" {
			break
		}
	}

	if workflowName == "" {
		if workflowFilter != "" {
			return "", "", nil, fmt.Errorf("agent '%s' not found in workflow '%s'", agentName, workflowFilter)
		}
		return "", "", nil, fmt.Errorf("no workflow found containing agent '%s'", agentName)
	}

	return cmdMetadataName, workflowName, foundCmd, nil
}

// generateForAgent generates shell code for a specific agent in a workflow.
// Returns (shellCode, shellType, error).
func generateForAgent(result *ApplyResult, workflowName string, cmdMetadataName string, cmd *manifest.Cmd) (string, string, error) {
	// Resolve the workflow with the specific cmd filter
	resolver := result.Resolver
	cmdFilter := []string{cmdMetadataName}

	resolved, err := resolver.ResolveKustomization(workflowName, cmdFilter)
	if err != nil {
		return "", "", fmt.Errorf("failed to resolve workflow '%s': %w", workflowName, err)
	}

	// Generate shell code
	generator := generate.NewGenerator(resolved.Name)
	if result.BuildResultYAML != "" {
		generator.SetBuildResult(result.BuildResultYAML)
	}

	shellCode, err := generator.GenerateKustomization(resolved)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate shell code: %w", err)
	}

	return shellCode, resolved.Shell, nil
}

// executeAgent writes a temp script, runs bash, and propagates exit code.
func executeAgent(shellCode string, shellType string, agentName string, extraArgs []string) {
	// Generate unique temp file name
	hash := sha256.Sum256([]byte(shellCode))
	tempFileName := fmt.Sprintf("kfg-run-%s.%s", hex.EncodeToString(hash[:8]), shellType)
	tempFile := filepath.Join(os.TempDir(), tempFileName)

	// Build the script: generated shell code + trap + agent call
	script := shellCode + "\n\n"
	script += fmt.Sprintf("trap 'rm -f %s' EXIT\n", tempFile)
	script += fmt.Sprintf("%s \"$@\"\n", agentName)

	// Write temp file
	err := os.WriteFile(tempFile, []byte(script), 0644)
	if err != nil {
		logger.Error("run", fmt.Sprintf("Failed to write temp file: %v", err))
		os.Exit(1)
	}

	// Build command args: script file + extra args
	cmdArgs := []string{tempFile}
	cmdArgs = append(cmdArgs, extraArgs...)

	// Execute the shell
	cmd := exec.Command(shellType, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	// Propagate exit code
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}

	os.Exit(0)
}

// listAvailableAgents prints all Cmds with their workflow names.
func listAvailableAgents(index *resolve.Index) {
	cmds := index.GetCmds()
	workflows := index.GetCmdWorkflows()

	if len(cmds) == 0 {
		fmt.Println("No agents found in manifests")
		return
	}

	// Build a map of cmd -> workflow
	cmdToWorkflow := make(map[string]string)
	for _, wf := range workflows {
		for _, cmdRef := range wf.Spec.Cmds {
			cmdToWorkflow[cmdRef] = wf.Metadata.Name
		}
	}

	fmt.Println("Available agents:")
	for _, cmd := range cmds {
		commandName := cmd.Metadata.CommandName
		if commandName == "" {
			commandName = cmd.Metadata.Name
		}
		workflowName := cmdToWorkflow[cmd.Metadata.Name]
		if workflowName == "" {
			workflowName = "unknown"
		}
		fmt.Printf("  %s (workflow: %s)\n", commandName, workflowName)
	}
}