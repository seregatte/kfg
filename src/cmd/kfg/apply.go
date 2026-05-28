package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/seregatte/kfg/src/internal/config"
	"github.com/seregatte/kfg/src/internal/converter"
	"github.com/seregatte/kfg/src/internal/generate"
	"github.com/seregatte/kfg/src/internal/kustomize"
	"github.com/seregatte/kfg/src/internal/logger"
	"github.com/seregatte/kfg/src/internal/manifest"
	"github.com/seregatte/kfg/src/internal/resolve"
	"github.com/spf13/cobra"
)

var (
	// Flags for the apply command
	applyKustomizePath string
	applyFile          string
	applyOutput        string
	applyWorkflow      string
	applyCmds          string
	applyConvert       string
	applyUse           string
	applyWith          string
	applyRefresh       bool
	applyStoreDir      string
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

Modes:
  1. Shell generation (default): Generate shell functions from Cmd/CmdWorkflow resources
     using -w/--workflow and -c/--cmds flags.

  2. Conversion mode: Transform Asset data using a Converter resource with --convert
     and --use flags. Outputs data in the Converter's specified format.

  3. Inline conversion: Use --convert with raw string input and --with for an inline
     yq expression, bypassing Converter resource lookup.

  4. Stdin pipeline: Use -f - with --with to pass stdin directly to the yq engine
     for multi-document merge operations.

The source can be provided as:
  - A positional argument (path or GitHub URL)
  - The -k flag (kustomization path or GitHub URL)
  - The -f flag (manifest file path or stdin)
  - The KFG_KPATH environment variable (used as fallback)

GitHub URLs are supported and will be cloned automatically:
  - https://github.com/owner/repo//path
  - https://github.com/owner/repo//path?ref=v1.0.0

Environment variables:
  KFG_KPATH      Default kustomization path if -k or -f not specified
  KFG_REFRESH    Set to "1" to invalidate and rebuild cache entries for cacheable Steps
  KFG_STORE_DIR  Custom store directory for cache entries (defaults to ~/.kfg/store)

Examples:
  # Shell generation
  kfg apply packages/domains/ai-agents/overlays/dev
  kfg apply -k packages/domains/ai-agents/overlays/dev
  kfg apply -k https://github.com/owner/repo//manifests
  kfg apply packages/domains/ai-agents/overlays/dev --workflow dev,openspec
  kfg apply -k packages/domains/ai-agents/overlays/dev --workflow ai-agents
  kfg apply -k packages/domains/ai-agents/overlays/dev --cmds claude
  kfg apply -f manifest.yaml
  kfg apply -f - (read from stdin)
  kfg apply -k packages/domains/ai-agents/overlays/dev --refresh  (invalidate and rebuild cache entries)
  KFG_KPATH=./manifests kfg apply
  KFG_KPATH=https://github.com/owner/repo//manifests kfg apply

  # Conversion mode
  kfg apply -f manifest.yaml --convert my-asset --use my-converter
  kfg apply -f manifest.yaml --convert my-asset --use my-converter -o output.json

  # Inline conversion
  kfg apply -f manifest.yaml --convert my-asset --with '.data | {"key": .value}'
  kfg apply -f manifest.yaml --convert '{"key":"value"}' --with '.key'

  # Stdin pipeline
  echo '{"a":1}---{"b":2}' | kfg apply -f - --with 'select(fi == 0) * select(fi == 1)'`,
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

		// KFG_KPATH fallback: if kustomize path is empty, use env var
		if applyKustomizePath == "" && applyFile == "" {
			applyKustomizePath = config.GetKPath()
		}

		// Validate flags/args
		if applyKustomizePath == "" && applyFile == "" {
			logger.Error("apply", "kustomization source required. Provide a path, use -k flag, -f flag, or set KFG_KPATH.")
			cmd.Help()
			os.Exit(2)
		}

		if applyKustomizePath != "" && applyFile != "" {
			logger.Error("apply", "Cannot use both -k and -f flags")
			os.Exit(2)
		}

		// Validate conversion mode mutual exclusivity
		if applyConvert != "" || applyUse != "" {
			// --convert and --use must be used together (--with can substitute for --use)
			if applyConvert == "" {
				logger.Error("apply", "--use requires --convert to be specified")
				os.Exit(2)
			}
			if applyUse == "" && applyWith == "" {
				logger.Error("apply", "--convert requires --use or --with to be specified")
				os.Exit(2)
			}
			// --convert/--use cannot be used with -w/--workflow
			if applyWorkflow != "" {
				logger.Error("apply", "--convert/--use cannot be used with --workflow/-w (shell generation flag)")
				os.Exit(2)
			}
			// --convert/--use cannot be used with -c/--cmds
			if applyCmds != "" {
				logger.Error("apply", "--convert/--use cannot be used with --cmds/-c (shell generation flag)")
				os.Exit(2)
			}
		}

		// Validate --with flag
		if err := validateWithFlag(applyWith, applyConvert, applyUse, applyFile, applyWorkflow, applyCmds); err != nil {
			logger.Error("apply", err.Error())
			os.Exit(2)
		}

		// Stdin raw mode: -f - with --with and no --convert
		// Must run BEFORE runApplyPipeline since stdin can only be read once
		if applyFile == "-" && applyWith != "" && applyConvert == "" {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				logger.Error("apply", fmt.Sprintf("Failed to read stdin: %v", err))
				os.Exit(1)
			}
			if err := runStdinConversion(string(data), applyWith, applyOutput); err != nil {
				logger.Error("apply", err.Error())
				os.Exit(1)
			}
			return
		}

		// Run the apply pipeline (GitHub URLs are passed directly to kustomize loader)
		result, err := runApplyPipeline(applyKustomizePath, applyFile)
		if err != nil {
			logger.Error("apply", err.Error())
			printApplyError(err, applyKustomizePath)
			os.Exit(1)
		}

		// Conversion mode (--convert + --use or --convert + --with)
		if applyConvert != "" {
			if err := runConversion(result.Resources, applyConvert, applyUse, applyWith, applyOutput); err != nil {
				logger.Error("apply", err.Error())
				os.Exit(1)
			}
			return
		}

		// Shell generation mode (existing behavior)
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

			// Prepend refresh header if refresh flag is set
			if applyRefresh {
				shellCode = "export KFG_REFRESH=1\n\n" + shellCode
			}
			if applyStoreDir != "" {
				shellCode = fmt.Sprintf("export KFG_STORE_DIR=%s\n\n", applyStoreDir) + shellCode
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

			// Prepend refresh header if refresh flag is set
			if applyRefresh {
				shellCode = "export KFG_REFRESH=1\n\n" + shellCode
			}
			if applyStoreDir != "" {
				shellCode = fmt.Sprintf("export KFG_STORE_DIR=%s\n\n", applyStoreDir) + shellCode
			}
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
	applyCmd.Flags().StringVar(&applyConvert, "convert", "", "Asset name for conversion mode")

	applyCmd.Flags().BoolVarP(&applyRefresh, "refresh", "r", false, "Invalidate and rebuild cache entries for cacheable Steps")
	applyCmd.Flags().StringVar(&applyUse, "use", "", "Converter name for conversion mode")
	applyCmd.Flags().StringVar(&applyWith, "with", "", "Inline yq expression for conversion mode (bypasses Converter lookup)")
	applyCmd.Flags().StringVar(&applyStoreDir, "store", "", "Custom store directory for cache entries (overrides KFG_STORE_DIR)")
}

// validateWithFlag checks --with flag mutual exclusivity and requirements.
// Returns an error if validation fails.
func validateWithFlag(with, convert, use, file, workflow, cmds string) error {
	if with == "" {
		return nil
	}
	if use != "" {
		return fmt.Errorf("--with and --use are mutually exclusive (use one or the other, not both)")
	}
	if convert == "" && file != "-" {
		return fmt.Errorf("--with requires --convert or -f - (stdin) to be specified")
	}
	if workflow != "" {
		return fmt.Errorf("--with cannot be used with --workflow/-w (shell generation flag)")
	}
	if cmds != "" {
		return fmt.Errorf("--with cannot be used with --cmds/-c (shell generation flag)")
	}
	return nil
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

// runConversion executes the conversion pipeline: find Asset (or treat as raw string),
// find Converter (or use inline expression), run engine, output result.
func runConversion(resources []manifest.ParsedResource, assetName, converterName, inlineExpression, outputFile string) error {
	engine := converter.NewEngine()

	// Try to find Asset by metadata.name
	var foundAsset *manifest.Assets
	var availableAssets []string
	for _, res := range resources {
		if res.Assets != nil {
			availableAssets = append(availableAssets, res.Assets.Metadata.Name)
			if res.Assets.Metadata.Name == assetName {
				foundAsset = res.Assets
			}
		}
	}

	// Asset found: use it
	if foundAsset != nil {
		asset := converter.MapManifestAssets(foundAsset)

		// --with: inline expression mode (skip Converter lookup)
		if inlineExpression != "" {
			result, err := engine.ApplyWithExpression(asset, inlineExpression)
			if err != nil {
				return fmt.Errorf("conversion failed: %w", err)
			}
			return writeOutput(result, outputFile)
		}

		// --use: Converter lookup mode (existing behavior)
		var foundConverter *manifest.Converter
		var availableConverters []string
		for _, res := range resources {
			if res.Converter != nil {
				availableConverters = append(availableConverters, res.Converter.Metadata.Name)
				if res.Converter.Metadata.Name == converterName {
					foundConverter = res.Converter
				}
			}
		}
		if foundConverter == nil {
			msg := fmt.Sprintf("Converter not found: %s", converterName)
			if len(availableConverters) > 0 {
				msg += fmt.Sprintf(" (available: %s)", strings.Join(availableConverters, ", "))
			}
			return fmt.Errorf("%s", msg)
		}

		conv := converter.MapManifestConverter(foundConverter)
		result, err := engine.Apply(conv, asset)
		if err != nil {
			return fmt.Errorf("conversion failed: %w", err)
		}
		return writeOutput(result, outputFile)
	}

	// Asset not found: try raw string fallback
	// Need either a converter or an inline expression to proceed
	if inlineExpression == "" && converterName == "" {
		msg := fmt.Sprintf("Asset not found: %s", assetName)
		if len(availableAssets) > 0 {
			msg += fmt.Sprintf(" (available: %s)", strings.Join(availableAssets, ", "))
		}
		return fmt.Errorf("%s", msg)
	}

	// Detect format for raw string fallback
	inputFormat := engine.DetectFormat(assetName)

	// Without --with, only accept structured JSON (objects/arrays) as raw input.
	// Plain YAML scalars like "nonexistent" are valid YAML but likely meant
	// to be asset names — don't silently treat them as raw input.
	if inlineExpression == "" && inputFormat == "yaml" && !looksLikeStructuredJSON(assetName) {
		msg := fmt.Sprintf("Asset not found: %s", assetName)
		if len(availableAssets) > 0 {
			msg += fmt.Sprintf(" (available: %s)", strings.Join(availableAssets, ", "))
		}
		return fmt.Errorf("%s", msg)
	}

	// 2. Find the converter (needed for both --with and raw input fallback)
	var foundConverter *manifest.Converter
	var availableConverters []string
	for _, res := range resources {
		if res.Converter != nil {
			availableConverters = append(availableConverters, res.Converter.Metadata.Name)
			if res.Converter.Metadata.Name == converterName {
				foundConverter = res.Converter
			}
		}
	}

	// 3. Raw input with converter (no --with): apply converter's expression to raw string
	if inlineExpression == "" && inputFormat != "" && converterName != "" {
		if foundConverter == nil {
			msg := fmt.Sprintf("Converter not found: %s", converterName)
			if len(availableConverters) > 0 {
				msg += fmt.Sprintf(" (available: %s)", strings.Join(availableConverters, ", "))
			}
			return fmt.Errorf("%s", msg)
		}
		conv := converter.MapManifestConverter(foundConverter)
		result, err := engine.ApplyRawWithConverter(assetName, inputFormat, conv)
		if err != nil {
			return fmt.Errorf("conversion failed: %w", err)
		}
		return writeOutput(result, outputFile)
	}

	// 4. --with mode: detect format, apply raw with inline expression
	if inputFormat == "" {
		msg := fmt.Sprintf("Asset not found: %s", assetName)
		if len(availableAssets) > 0 {
			msg += fmt.Sprintf(" (available: %s)", strings.Join(availableAssets, ", "))
		}
		msg += fmt.Sprintf("\nNo matching Asset found and input is not valid JSON or YAML")
		return fmt.Errorf("%s", msg)
	}

	result, err := engine.ApplyRaw(assetName, inputFormat, inlineExpression)
	if err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}
	return writeOutput(result, outputFile)
}

// looksLikeStructuredJSON returns true if the input looks like a JSON
// object or array (not a plain scalar). Used to distinguish intentional
// raw JSON input from asset names that happen to be valid YAML scalars.
func looksLikeStructuredJSON(input string) bool {
	trimmed := strings.TrimSpace(input)
	return strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[")
}

// runStdinConversion reads stdin and applies an inline yq expression directly.
// No manifest parsing occurs.
func runStdinConversion(input, expression, outputFile string) error {
	engine := converter.NewEngine()
	result, err := engine.ApplyRaw(input, engine.DetectFormat(input), expression)
	if err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}
	return writeOutput(result, outputFile)
}

// writeOutput writes the conversion result to a file or stdout.
func writeOutput(result, outputFile string) error {
	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(result), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		logger.Info("apply", fmt.Sprintf("Wrote conversion output to %s", outputFile))
	} else {
		fmt.Print(result)
	}
	return nil
}
