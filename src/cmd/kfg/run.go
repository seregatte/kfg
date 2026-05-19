package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/seregatte/kfg/src/internal/config"
	"github.com/seregatte/kfg/src/internal/generate"
	"github.com/seregatte/kfg/src/internal/logger"
	"github.com/seregatte/kfg/src/internal/manifest"
	"github.com/seregatte/kfg/src/internal/resolve"
	"github.com/spf13/cobra"
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
	Use:   "run [cmd] [-- extra-args...]",
	Short: "Run a command with one-shot execution",
	Long: `Run a command by generating shell code, sourcing it, and executing in one invocation.

This command provides a one-shot "generate → source → execute" experience
for running commands. It matches commands by their commandName (e.g., "my-cmd")
and auto-detects the workflow if not specified.

The source can be provided as:
  - The -k flag (kustomization path or GitHub URL)
  - The -f flag (manifest file path or stdin)
  - The KFG_KPATH environment variable (used as fallback)

GitHub URLs are supported and will be cloned automatically:
  - https://github.com/owner/repo//path
  - https://github.com/owner/repo//path?ref=v1.0.0

Arguments after '--' are passed directly to the command.

Examples:
  kfg run -k .manifests/overlay/dev my-cmd
  kfg run -k .manifests/overlay/dev my-cmd -- --flag value
  kfg run -k https://github.com/owner/repo//manifests my-cmd
  kfg run -k .manifests/overlay/dev -w dev my-cmd
  kfg run -f manifest.yaml my-cmd
  kfg run -k .manifests/overlay/dev (lists available commands)
  KFG_KPATH=./manifests kfg run my-cmd
  KFG_KPATH=https://github.com/owner/repo//manifests kfg run my-cmd`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// KFG_KPATH fallback: if kustomize path is empty, use env var
		if runKustomizePath == "" && runFile == "" {
			runKustomizePath = config.GetKPath()
		}

		// Validate flags
		if runKustomizePath == "" && runFile == "" {
			logger.Error("run", "kustomization source required. Provide a path, use -k flag, -f flag, or set KFG_KPATH.")
			cmd.Help()
			os.Exit(2)
		}

		if runKustomizePath != "" && runFile != "" {
			logger.Error("run", "Cannot use both -k and -f flags")
			os.Exit(2)
		}

		// Parse run args: command name and extra args
		cmdName, extraArgs := parseLaunchArgs(cmd, args)

		// No command name provided - list available commands
		if cmdName == "" {
			// Run the apply pipeline to get the index (GitHub URLs are passed directly to kustomize loader)
			result, err := runApplyPipeline(runKustomizePath, runFile)
			if err != nil {
				logger.Error("run", err.Error())
				os.Exit(1)
			}
			listAvailableCmds(result.Index)
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

		// Find the command in the index
		cmdMetadataName, workflowName, foundCmd, err := findCmd(result.Index, cmdName, workflowFilter)
		if err != nil {
			logger.Error("run", err.Error())
			listAvailableCmds(result.Index)
			os.Exit(1)
		}

		// Generate shell code for the specific command
		shellCode, shellType, err := generateForCmd(result, workflowName, cmdMetadataName, foundCmd)
		if err != nil {
			logger.Error("run", fmt.Sprintf("Failed to generate shell code: %v", err))
			os.Exit(1)
		}

		// Execute the command
		executeCmd(shellCode, shellType, cmdName, extraArgs)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Add flags for the run command
	runCmd.Flags().StringVarP(&runKustomizePath, "kustomize", "k", "", "Kustomization directory path")
	runCmd.Flags().StringVarP(&runFile, "file", "f", "", "Manifest file path (use '-' for stdin)")
	runCmd.Flags().StringVarP(&runWorkflow, "workflow", "w", "", "CmdWorkflow name (auto-detected if not specified)")
	runCmd.Flags().StringVarP(&runCmds, "cmds", "c", "", "Comma-separated list of cmds to run")
}

// parseLaunchArgs splits args using Cobra's dash boundary into command name and extra args.
// When the user passes `--`, Cobra strips the separator and records the split point
// via cmd.ArgsLenAtDash(). This function uses that metadata as the source of truth.
// Returns (cmdName, extraArgs).
func parseLaunchArgs(cmd *cobra.Command, args []string) (string, []string) {
	dashIndex := -1
	if cmd != nil {
		dashIndex = cmd.ArgsLenAtDash()
	}
	return splitArgsAtDash(dashIndex, args)
}

// splitArgsAtDash splits args at the dash boundary into command name and extra args.
// dashIndex: -1 if no `--` was present, 0 if `--` was before any positional args,
// or > 0 if there were positional args before `--`.
// Returns (cmdName, extraArgs).
func splitArgsAtDash(dashIndex int, args []string) (string, []string) {
	if dashIndex == -1 {
		// No separator - first arg is command name, no extra args
		if len(args) > 0 {
			return args[0], []string{}
		}
		return "", []string{}
	}

	// dashIndex == 0: separator was before any positional args, all args are extra
	if dashIndex == 0 {
		return "", args
	}

	// dashIndex > 0: command is args[0], extra args start at the dash boundary
	if len(args) > 0 {
		return args[0], args[dashIndex:]
	}
	return "", []string{}
}

// findCmd matches a command by commandName and finds its workflow.
// Returns (cmdMetadataName, workflowName, foundCmd, error).
func findCmd(index *resolve.Index, cmdName string, workflowFilter string) (string, string, *manifest.Cmd, error) {
	cmds := index.GetCmds()

	// Find cmd by commandName
	var foundCmd *manifest.Cmd
	var cmdMetadataName string
	for _, cmd := range cmds {
		if cmd.Metadata.CommandName == cmdName {
			foundCmd = cmd
			cmdMetadataName = cmd.Metadata.Name
			break
		}
	}

	if foundCmd == nil {
		return "", "", nil, fmt.Errorf("command '%s' not found", cmdName)
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
			return "", "", nil, fmt.Errorf("command '%s' not found in workflow '%s'", cmdName, workflowFilter)
		}
		return "", "", nil, fmt.Errorf("no workflow found containing command '%s'", cmdName)
	}

	return cmdMetadataName, workflowName, foundCmd, nil
}

// generateForCmd generates shell code for a specific command in a workflow.
// Returns (shellCode, shellType, error).
func generateForCmd(result *ApplyResult, workflowName string, cmdMetadataName string, cmd *manifest.Cmd) (string, string, error) {
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

// executeCmd writes a temp script, runs bash, and propagates exit code.
func executeCmd(shellCode string, shellType string, cmdName string, extraArgs []string) {
	// Generate unique temp file name
	hash := sha256.Sum256([]byte(shellCode))
	tempFileName := fmt.Sprintf("kfg-run-%s.%s", hex.EncodeToString(hash[:8]), shellType)
	tempFile := filepath.Join(os.TempDir(), tempFileName)

	// Build the script: generated shell code + trap + command call
	script := shellCode + "\n\n"
	script += fmt.Sprintf("trap 'rm -f %s' EXIT\n", tempFile)
	script += fmt.Sprintf("%s \"$@\"\n", cmdName)

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

// listAvailableCmds prints all Cmds with their workflow names.
func listAvailableCmds(index *resolve.Index) {
	cmds := index.GetCmds()
	workflows := index.GetCmdWorkflows()

	if len(cmds) == 0 {
		fmt.Println("No commands found in manifests")
		return
	}

	// Build a map of cmd -> workflow
	cmdToWorkflow := make(map[string]string)
	for _, wf := range workflows {
		for _, cmdRef := range wf.Spec.Cmds {
			cmdToWorkflow[cmdRef] = wf.Metadata.Name
		}
	}

	fmt.Println("Available commands:")
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
