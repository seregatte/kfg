package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/seregatte/kfg/src/internal/generate"
	"github.com/seregatte/kfg/src/internal/kustomize"
	"github.com/seregatte/kfg/src/internal/manifest"
	"github.com/seregatte/kfg/src/internal/resolve"
	"github.com/seregatte/kfg/src/internal/logger"
)

var (
	// Flags for the apply command
	applyKustomizePath string
	applyFile          string
	applyOutput        string
	applyWorkflow      string
	applyCmds          string
)

// ApplyResult holds the result of the apply pipeline (load → validate → index → resolve).
// It is used by both `apply` and `launch` commands.
type ApplyResult struct {
	Resources       []manifest.ParsedResource
	Shell           string
	BuildResultYAML string
	Index           *resolve.Index
	Resolver        *resolve.Resolver
}

// runApplyPipeline executes the core apply pipeline: load → validate → index → resolve.
// It returns an ApplyResult that can be used by both `apply` and `launch` commands.
// This function does not produce side effects (no output, no file writes).
func runApplyPipeline(kustomizePath, file string) (*ApplyResult, error) {
	var resources []manifest.ParsedResource
	var shell string
	var buildResultYAML string

	// Load resources based on flag
	if kustomizePath != "" {
		// Load via kustomize
		loader := kustomize.NewLoader(nil)
		resMap, err := loader.Load(kustomizePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load kustomization: %w", err)
		}

		// Convert ResMap to ParsedResource
		adapter := kustomize.NewAdapter()
		resources, err = adapter.ResMapToResources(resMap)
		if err != nil {
			return nil, fmt.Errorf("failed to convert resources: %w", err)
		}

		// Get build result YAML (all five kinds from ResMap)
		buildResultBytes, err := resMap.AsYaml()
		if err != nil {
			return nil, fmt.Errorf("failed to serialize build result: %w", err)
		}
		buildResultYAML = string(buildResultBytes)

		// Shell will be determined by workflow metadata
		shell = "bash" // Default, will be overridden by workflow

	} else if file != "" {
		// Load from file or stdin
		parser := manifest.NewParser()
		var data []byte
		var err error

		if file == "-" {
			// Read from stdin
			data, err = io.ReadAll(os.Stdin)
			if err != nil {
				return nil, fmt.Errorf("failed to read stdin: %w", err)
			}
		} else {
			// Read from file
			data, err = os.ReadFile(file)
			if err != nil {
				return nil, fmt.Errorf("failed to read file: %w", err)
			}
		}

		resources, err = parser.ParseData(file, data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse manifest: %w", err)
		}

		shell = "bash" // Default
	}

	// Validate resources
	for _, res := range resources {
		err := res.Validate()
		if err != nil {
			return nil, fmt.Errorf("validation failed for %s: %w", res.Name(), err)
		}
	}

	// Create index and resolver
	index := resolve.NewIndex(resources)
	resolver := resolve.NewResolver(index)

	return &ApplyResult{
		Resources:       resources,
		Shell:           shell,
		BuildResultYAML: buildResultYAML,
		Index:           index,
		Resolver:        resolver,
	}, nil
}

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply [path]",
	Short: "Apply kustomization and generate shell code",
	Long: `Apply a kustomization or manifest file and generate shell code.

This command processes a kustomization or manifest file, resolves the workflow,
and generates shell functions that can be sourced or used interactively.

The [path] argument can be a kustomization directory path or a manifest file.
If not provided, the -k or -f flag must be used.

Examples:
  kfg apply .nixai/overlay/dev
  kfg apply -k .nixai/overlay/dev
  kfg apply .nixai/overlay/dev --workflow dev,openspec
  kfg apply -k .nixai/overlay/dev --workflow ai-agents
  kfg apply -k .nixai/overlay/dev --cmds claude
  kfg apply -f manifest.yaml
  kfg apply -f - (read from stdin)`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Handle positional argument - if provided, use it as kustomization path
		if len(args) > 0 {
			// Positional argument provided - use as kustomization path
			if applyKustomizePath == "" && applyFile == "" {
				applyKustomizePath = args[0]
			} else if applyKustomizePath != "" {
				logger.Error("apply", "Cannot use positional argument with -k flag")
				os.Exit(2)
			} else if applyFile != "" {
				logger.Error("apply", "Cannot use positional argument with -f flag")
				os.Exit(2)
			}
		}

		// Validate flags/args
		if applyKustomizePath == "" && applyFile == "" {
			logger.Error("apply", "Either a positional path, -k (kustomize path), or -f (file) is required")
			cmd.Help()
			os.Exit(2)
		}

		if applyKustomizePath != "" && applyFile != "" {
			logger.Error("apply", "Cannot use both -k and -f flags")
			os.Exit(2)
		}

		// Run the apply pipeline
		result, err := runApplyPipeline(applyKustomizePath, applyFile)
		if err != nil {
			logger.Error("apply", err.Error())
			printApplyError(err, applyKustomizePath)
			os.Exit(1)
		}

		// Resolve cmd filter
		var cmdFilter []string
		if applyCmds != "" {
			cmds := strings.Split(applyCmds, ",")
			for i, c := range cmds {
				cmds[i] = strings.TrimSpace(c)
			}
			cmdFilter = cmds
		}

		// Parse workflow flag for comma-separated values
		workflowNames := parseWorkflowFlag(applyWorkflow)

		// Determine if we're in single-workflow or multi-workflow mode
		isMultiWorkflow := len(workflowNames) > 1 || (len(workflowNames) == 0 && len(result.Index.GetCmdWorkflows()) > 1)

		// Resolve and generate based on mode
		if isMultiWorkflow {
			// Multi-workflow mode
			shellCode, _, shellType, err := resolveAndGenerateMultiWorkflow(result.Resolver, result.Index, workflowNames, result.BuildResultYAML)
			if err != nil {
				logger.Error("apply", fmt.Sprintf("Failed to resolve workflows: %v", err))
				printResolutionErrorNew(err, result.Index)
				os.Exit(1)
			}
			result.Shell = shellType

			// Output to file or stdout
			if applyOutput != "" {
				err = os.WriteFile(applyOutput, []byte(shellCode), 0644)
				if err != nil {
					logger.Error("apply", fmt.Sprintf("Failed to write output file: %v", err))
					os.Exit(1)
				}
				logger.Info("apply", fmt.Sprintf("Wrote output to %s", applyOutput))
			} else {
				fmt.Print(shellCode)
			}
		} else {
			// Single-workflow mode (existing behavior)
			workflowName := ""
			if len(workflowNames) == 1 {
				workflowName = workflowNames[0]
			}

			resolved, err := result.Resolver.ResolveKustomization(workflowName, cmdFilter)
			if err != nil {
				logger.Error("apply", fmt.Sprintf("Failed to resolve kustomization: %v", err))
				printResolutionErrorNew(err, result.Index)
				os.Exit(1)
			}

			result.Shell = resolved.Shell

			// Generate shell code
			generator := generate.NewGenerator(resolved.Name)
			if result.BuildResultYAML != "" {
				generator.SetBuildResult(result.BuildResultYAML)
			}
			shellCode, err := generator.GenerateKustomization(resolved)
			if err != nil {
				logger.Error("apply", fmt.Sprintf("Failed to generate shell code: %v", err))
				os.Exit(1)
			}

			// Output to file or stdout
			if applyOutput != "" {
				err = os.WriteFile(applyOutput, []byte(shellCode), 0644)
				if err != nil {
					logger.Error("apply", fmt.Sprintf("Failed to write output file: %v", err))
					os.Exit(1)
				}
				logger.Info("apply", fmt.Sprintf("Wrote output to %s", applyOutput))
			} else {
				fmt.Print(shellCode)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Add flags for the apply command
	applyCmd.Flags().StringVarP(&applyKustomizePath, "kustomize", "k", "", "Kustomization directory path")
	applyCmd.Flags().StringVarP(&applyFile, "file", "f", "", "Manifest file path (use '-' for stdin)")
	applyCmd.Flags().StringVarP(&applyOutput, "output", "o", "", "Output file path (default: stdout)")
	applyCmd.Flags().StringVarP(&applyWorkflow, "workflow", "w", "", "CmdWorkflow name(s), comma-separated (default: all workflows)")
	applyCmd.Flags().StringVarP(&applyCmds, "cmds", "c", "", "Comma-separated list of cmds to generate")
}

func printApplyError(err error, path string) {
	logger.Error("apply", fmt.Sprintf("Failed to apply kustomization at %s: %v", path, err))

	if os.IsNotExist(err) {
		logger.Error("apply", "The specified path does not exist. Make sure the kustomization directory or file exists.")
	}
}

func printResolutionErrorNew(err error, index *resolve.Index) {
	// Check if it's a workflow not found error
	if strings.Contains(err.Error(), "CmdWorkflow not found") {
		workflows := index.GetCmdWorkflows()
		if len(workflows) > 0 {
			logger.Error("apply", fmt.Sprintf("CmdWorkflow not found: %v", err))
			logger.Info("apply", "Available CmdWorkflows:")
			for _, wf := range workflows {
				logger.Info("apply", fmt.Sprintf("  - %s", wf.Metadata.Name))
			}
		} else {
			logger.Error("apply", "No CmdWorkflows found in manifests")
			logger.Info("apply", "Create a CmdWorkflow resource with kind: CmdWorkflow")
		}
		return
	}

	// Check if it's a cmd not found error
	if strings.Contains(err.Error(), "Cmd not found") {
		cmds := index.GetCmds()
		if len(cmds) > 0 {
			logger.Error("apply", fmt.Sprintf("Cmd not found: %v", err))
			logger.Info("apply", "Available Cmds:")
			for _, cmd := range cmds {
				logger.Info("apply", fmt.Sprintf("  - %s", cmd.Metadata.Name))
			}
		} else {
			logger.Error("apply", "No Cmds found in manifests")
			logger.Info("apply", "Create Cmd resources with kind: Cmd")
		}
		return
	}

	// Check if it's a step not found error
	if strings.Contains(err.Error(), "step not found") {
		steps := index.GetSteps()
		if len(steps) > 0 {
			logger.Error("apply", fmt.Sprintf("Step not found: %v", err))
			logger.Info("apply", "Available Steps:")
			for _, step := range steps {
				logger.Info("apply", fmt.Sprintf("  - %s", step.Metadata.Name))
			}
		} else {
			logger.Error("apply", "No Steps found in manifests")
			logger.Info("apply", "Create Step resources with kind: Step")
		}
		return
	}

	// Generic resolution error
	logger.Error("apply", fmt.Sprintf("Resolution failed: %v", err))
}

// parseWorkflowFlag parses the workflow flag for comma-separated values.
// Returns a slice of workflow names, trimming whitespace around each.
// Returns empty slice if no workflow specified.
func parseWorkflowFlag(workflowFlag string) []string {
	if workflowFlag == "" {
		return []string{}
	}

	names := strings.Split(workflowFlag, ",")
	result := make([]string, 0, len(names))
	for _, name := range names {
		trimmed := strings.TrimSpace(name)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// resolveAndGenerateMultiWorkflow resolves multiple workflows and generates shell code.
// Returns the shell code, kustomization name, shell type, and any error.
func resolveAndGenerateMultiWorkflow(resolver *resolve.Resolver, index *resolve.Index, workflowNames []string, buildResultYAML string) (string, string, string, error) {
	// Resolve workflows
	var workflows []*resolve.ResolvedCmdWorkflow
	var err error

	if len(workflowNames) == 0 {
		// No workflow specified - resolve all workflows
		workflows, err = resolver.ResolveAllWorkflows()
		if err != nil {
			return "", "", "", err
		}
	} else {
		// Specific workflows specified
		workflows, err = resolver.ResolveWorkflowsByName(workflowNames)
		if err != nil {
			return "", "", "", err
		}
	}

	// Determine kustomization name and shell
	// Use first workflow's shell (assumes consistency per design)
	kustomizationName := "multi-workflow"
	shell := "bash"
	if len(workflows) > 0 {
		shell = workflows[0].Shell
		// Use kustomization name from workflow if available
		// Note: in multi-workflow mode, we use a generic name or derive from the overlay
	}

	// Collect all steps and cmds for the multi-workflow resolved
	allSteps := make(map[string]*manifest.Step)
	allCmds := make(map[string]*manifest.Cmd)
	
	// Collect from index (available steps and cmds)
	for _, step := range index.GetSteps() {
		allSteps[step.Metadata.Name] = step
	}
	for _, cmd := range index.GetCmds() {
		allCmds[cmd.Metadata.Name] = cmd
	}
	
	// Also add steps from resolved workflows (referenced steps)
	for _, wf := range workflows {
		for _, step := range wf.BeforeSteps {
			allSteps[step.Step.Metadata.Name] = step.Step
		}
		for _, step := range wf.AfterSteps {
			allSteps[step.Step.Metadata.Name] = step.Step
		}
		for _, entry := range wf.Cmds {
			allCmds[entry.Cmd.Metadata.Name] = entry.Cmd
			for _, step := range entry.BeforeSteps {
				allSteps[step.Step.Metadata.Name] = step.Step
			}
			for _, step := range entry.AfterSteps {
				allSteps[step.Step.Metadata.Name] = step.Step
			}
		}
	}

	// Create multi-workflow resolved
	multi := &generate.ResolvedMultiWorkflow{
		Name:      kustomizationName,
		Shell:     shell,
		Workflows: workflows,
		Steps:     allSteps,
		Cmds:      allCmds,
	}

	// Generate shell code
	generator := generate.NewGenerator(kustomizationName)
	if buildResultYAML != "" {
		generator.SetBuildResult(buildResultYAML)
	}
	shellCode, err := generator.GenerateAllWorkflows(multi)
	if err != nil {
		return "", "", "", err
	}

	return shellCode, kustomizationName, shell, nil
}