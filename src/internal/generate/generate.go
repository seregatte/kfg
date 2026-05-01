package generate

import (
	"encoding/base64"
	"sort"
	"strings"

	"github.com/seregatte/kfg/src/internal/generate/templates"
	"github.com/seregatte/kfg/src/internal/manifest"
	"github.com/seregatte/kfg/src/internal/resolve"
	"github.com/seregatte/kfg/src/internal/resolver"
)

// Generator generates shell code from resolved resources.
type Generator struct {
	setName         string
	templateMgr     *templates.TemplateManager
	buildResultYAML string // Build result YAML to embed in generated shell
}

// NewGenerator creates a new Generator.
func NewGenerator(setName string) *Generator {
	templateMgr, err := templates.NewTemplateManager()
	if err != nil {
		panic(err)
	}

	return &Generator{
		setName:         setName,
		templateMgr:     templateMgr,
		buildResultYAML: "",
	}
}

// SetBuildResult sets the build result YAML to embed in generated shell.
// The build result is the full YAML output from the kustomize loader (all five kinds).
func (g *Generator) SetBuildResult(yaml string) {
	g.buildResultYAML = yaml
}

// GetBuildResult returns the current build result YAML.
func (g *Generator) GetBuildResult() string {
	return g.buildResultYAML
}

// GenerateKustomization generates shell code from a ResolvedKustomization.
func (g *Generator) GenerateKustomization(rk *resolve.ResolvedKustomization) (string, error) {
	var code strings.Builder

	// Header
	headerData := templates.HeaderData{
		KustomizationName: rk.Name,
		WorkflowName:      rk.Workflow.Workflow.Metadata.Name,
		Shell:             rk.Shell,
		SetName:           rk.Name,
	}
	header, err := g.templateMgr.ExecuteHeader(headerData)
	if err != nil {
		return "", err
	}
	code.WriteString(header)
	code.WriteString("\n")

	// Metadata environment variables
	g.generateMetadataEnv(&code, rk)

	// Artifacts array declaration
	g.generateArtifactsDeclaration(&code, rk)

	// Build result file setup (global scope)
	g.generateGlobalBuildResult(&code)

	// Helpers
	helpers, err := g.templateMgr.ExecuteHelper()
	if err != nil {
		return "", err
	}
	code.WriteString(helpers)
	code.WriteString("\n")

	// Step functions
	g.generateWorkflowStepFunctions(&code, rk.Workflow)

	// Cmd wrappers
	g.generateWorkflowCmdWrappers(&code, rk.Workflow)

	return code.String(), nil
}

// ResolvedMultiWorkflow represents the resolved output for multi-workflow shell generation.
type ResolvedMultiWorkflow struct {
	Name      string                      // Kustomization name (directory name)
	Shell     string                      // Shell type from first workflow
	Workflows []*resolve.ResolvedCmdWorkflow // All resolved workflows
	Steps     map[string]*manifest.Step   // All available steps
	Cmds      map[string]*manifest.Cmd   // All available cmds
}

// GenerateAllWorkflows generates shell code from multiple resolved workflows.
// Steps are deduplicated by name across all workflows.
// KFG_WORKFLOW_NAME is omitted as it's not meaningful with multiple workflows.
func (g *Generator) GenerateAllWorkflows(rm *ResolvedMultiWorkflow) (string, error) {
	var code strings.Builder

	// Header - use kustomization name, no workflow name
	headerData := templates.HeaderData{
		KustomizationName: rm.Name,
		WorkflowName:      "", // Empty for multi-workflow
		Shell:             rm.Shell,
		SetName:           rm.Name,
	}
	header, err := g.templateMgr.ExecuteHeader(headerData)
	if err != nil {
		return "", err
	}
	code.WriteString(header)
	code.WriteString("\n")

	// Metadata environment variables - no KFG_WORKFLOW_NAME for multi-workflow
	g.generateMultiWorkflowMetadataEnv(&code, rm)

	// Artifacts array declaration
	g.generateMultiWorkflowArtifactsDeclaration(&code, rm)

	// Build result file setup (global scope)
	g.generateGlobalBuildResult(&code)

	// Helpers
	helpers, err := g.templateMgr.ExecuteHelper()
	if err != nil {
		return "", err
	}
	code.WriteString(helpers)
	code.WriteString("\n")

	// Step functions (deduplicated across all workflows)
	g.generateMultiWorkflowStepFunctions(&code, rm.Workflows)

	// Cmd wrappers for all workflows
	g.generateMultiWorkflowCmdWrappers(&code, rm.Workflows)

	return code.String(), nil
}

// generateMultiWorkflowMetadataEnv generates environment variables for multi-workflow mode.
// KFG_WORKFLOW_NAME is omitted as it's not meaningful with multiple workflows.
func (g *Generator) generateMultiWorkflowMetadataEnv(code *strings.Builder, rm *ResolvedMultiWorkflow) {
	code.WriteString("# Metadata environment variables\n")
	code.WriteString("KFG_KUSTOMIZATION_NAME=" + rm.Name + "\n")
	// KFG_WORKFLOW_NAME omitted for multi-workflow mode
	code.WriteString("KFG_SHELL=" + rm.Shell + "\n")
	code.WriteString("\n")
}

// generateMultiWorkflowArtifactsDeclaration generates the artifacts array declaration.
func (g *Generator) generateMultiWorkflowArtifactsDeclaration(code *strings.Builder, rm *ResolvedMultiWorkflow) {
	code.WriteString("# Artifacts tracking\n")
	code.WriteString("declare -a KFG_ARTIFACTS=()\n")
	code.WriteString("\n")

	// Generate helper function to add artifacts
	code.WriteString("__kfg_add_artifact() {\n")
	code.WriteString("    local artifact=\"$1\"\n")
	code.WriteString("    KFG_ARTIFACTS+=(\"$artifact\")\n")
	code.WriteString("}\n")
	code.WriteString("\n")
}

// generateMultiWorkflowStepFunctions generates step functions from all workflows with deduplication.
// Steps with the same name from different workflows are treated as identical and generated once.
func (g *Generator) generateMultiWorkflowStepFunctions(code *strings.Builder, workflows []*resolve.ResolvedCmdWorkflow) {
	code.WriteString("# Step execution functions\n")
	code.WriteString("\n")

	// Collect all unique steps by name (deduplication)
	allSteps := make(map[string]*manifest.Step)

	for _, rw := range workflows {
		// Global before/after steps
		for _, step := range rw.BeforeSteps {
			allSteps[step.Step.Metadata.Name] = step.Step
		}
		for _, step := range rw.AfterSteps {
			allSteps[step.Step.Metadata.Name] = step.Step
		}

		// Per-cmd steps
		for _, entry := range rw.Cmds {
			for _, step := range entry.BeforeSteps {
				allSteps[step.Step.Metadata.Name] = step.Step
		}
			for _, step := range entry.AfterSteps {
				allSteps[step.Step.Metadata.Name] = step.Step
			}
		}
	}

	// Sort names for deterministic output
	var stepNames []string
	for name := range allSteps {
		stepNames = append(stepNames, name)
	}
	sort.Strings(stepNames)

	// Generate each step function
	for _, name := range stepNames {
		step := allSteps[name]
		stepData := g.convertStepToTemplateData(step)
		stepCode, err := g.templateMgr.ExecuteStep(stepData)
		if err != nil {
			panic(err)
		}
		code.WriteString(stepCode)
		code.WriteString("\n")
	}
}

// generateMultiWorkflowCmdWrappers generates cmd wrappers for all workflows.
func (g *Generator) generateMultiWorkflowCmdWrappers(code *strings.Builder, workflows []*resolve.ResolvedCmdWorkflow) {
	code.WriteString("# Cmd wrappers\n")
	code.WriteString("\n")

	// Collect all cmd names across workflows to detect conflicts
	allCmdNames := make(map[string]string) // cmdName -> workflowName

	for _, rw := range workflows {
		for cmdName := range rw.Cmds {
			if existingWorkflow, exists := allCmdNames[cmdName]; exists {
				// Log warning about duplicate cmd name (both workflows define same cmd)
				// In shell, later definitions override earlier ones, so we just proceed
				_ = existingWorkflow // suppress unused warning
			}
			allCmdNames[cmdName] = rw.Workflow.Metadata.Name
		}
	}

	// Sort workflow names for deterministic output
	workflowNames := make([]string, 0, len(workflows))
	workflowMap := make(map[string]*resolve.ResolvedCmdWorkflow)
	for _, rw := range workflows {
		name := rw.Workflow.Metadata.Name
		workflowNames = append(workflowNames, name)
		workflowMap[name] = rw
	}
	sort.Strings(workflowNames)

	// Generate cmd wrappers for each workflow
	for _, wfName := range workflowNames {
		rw := workflowMap[wfName]
		cmdNames := rw.GetAllCmdNames()

		for _, name := range cmdNames {
			entry := rw.Cmds[name]
			cmdData := g.convertWorkflowCmdToTemplateData(rw, entry)
			cmdCode := g.generateWorkflowCmdCode(cmdData)
			code.WriteString(cmdCode)
			code.WriteString("\n")
		}
	}
}

// generateMetadataEnv generates environment variables for metadata.
func (g *Generator) generateMetadataEnv(code *strings.Builder, rk *resolve.ResolvedKustomization) {
	code.WriteString("# Metadata environment variables\n")
	code.WriteString("KFG_KUSTOMIZATION_NAME=" + rk.Name + "\n")
	code.WriteString("KFG_WORKFLOW_NAME=" + rk.Workflow.Workflow.Metadata.Name + "\n")
	code.WriteString("KFG_SHELL=" + rk.Shell + "\n")
	code.WriteString("\n")
}

// generateArtifactsDeclaration generates the artifacts array declaration.
func (g *Generator) generateArtifactsDeclaration(code *strings.Builder, rk *resolve.ResolvedKustomization) {
	code.WriteString("# Artifacts tracking\n")
	code.WriteString("declare -a KFG_ARTIFACTS=()\n")
	code.WriteString("\n")

	// Generate helper function to add artifacts
	code.WriteString("__kfg_add_artifact() {\n")
	code.WriteString("    local artifact=\"$1\"\n")
	code.WriteString("    KFG_ARTIFACTS+=(\"$artifact\")\n")
	code.WriteString("}\n")
	code.WriteString("\n")
}

// generateGlobalBuildResult generates the build result file setup at global scope.
// This includes creating the temp file, decoding the base64-encoded YAML, exporting
// the environment variable, defining the helper function, and registering the EXIT trap.
func (g *Generator) generateGlobalBuildResult(code *strings.Builder) {
	if g.buildResultYAML == "" {
		return // No build result to emit
	}

	// Base64 encode the build result for safe embedding in shell code
	encodedBuildResult := base64.StdEncoding.EncodeToString([]byte(g.buildResultYAML))

	code.WriteString("# Build result file setup (global scope)\n")
	code.WriteString("__kfg_build_result_file=$(mktemp -t nixai-build-XXXXXX.yaml)\n")
	code.WriteString("echo \"" + encodedBuildResult + "\" | base64 -d > \"$__kfg_build_result_file\"\n")
	code.WriteString("export KFG_BUILD_RESULT_FILE=$__kfg_build_result_file\n")
	code.WriteString("\n")

	code.WriteString("# Build result helper function\n")
	code.WriteString("__kfg_build_result() {\n")
	code.WriteString("    cat \"$KFG_BUILD_RESULT_FILE\"\n")
	code.WriteString("}\n")
	code.WriteString("\n")

	code.WriteString("# Cleanup trap for build result file on shell exit\n")
	code.WriteString("trap 'rm -f \"$__kfg_build_result_file\"' EXIT\n")
	code.WriteString("\n")
}

// generateWorkflowStepFunctions generates step functions from workflow.
func (g *Generator) generateWorkflowStepFunctions(code *strings.Builder, rw *resolve.ResolvedCmdWorkflow) {
	code.WriteString("# Step execution functions\n")
	code.WriteString("\n")

	// Collect all unique steps
	allSteps := make(map[string]*manifest.Step)

	// Global before/after steps
	for _, step := range rw.BeforeSteps {
		allSteps[step.Step.Metadata.Name] = step.Step
	}
	for _, step := range rw.AfterSteps {
		allSteps[step.Step.Metadata.Name] = step.Step
	}

	// Per-cmd steps
	for _, entry := range rw.Cmds {
		for _, step := range entry.BeforeSteps {
			allSteps[step.Step.Metadata.Name] = step.Step
		}
		for _, step := range entry.AfterSteps {
			allSteps[step.Step.Metadata.Name] = step.Step
		}
	}

	// Sort names for deterministic output
	var stepNames []string
	for name := range allSteps {
		stepNames = append(stepNames, name)
	}
	sort.Strings(stepNames)

	// Generate each step function
	for _, name := range stepNames {
		step := allSteps[name]
		stepData := g.convertStepToTemplateData(step)
		stepCode, err := g.templateMgr.ExecuteStep(stepData)
		if err != nil {
			panic(err)
		}
		code.WriteString(stepCode)
		code.WriteString("\n")
	}
}

// generateWorkflowCmdWrappers generates cmd wrappers from workflow.
func (g *Generator) generateWorkflowCmdWrappers(code *strings.Builder, rw *resolve.ResolvedCmdWorkflow) {
	code.WriteString("# Cmd wrappers\n")
	code.WriteString("\n")

	// Sort cmd names for deterministic output
	cmdNames := rw.GetAllCmdNames()

	// Generate each cmd wrapper
	for _, name := range cmdNames {
		entry := rw.Cmds[name]
		cmdData := g.convertWorkflowCmdToTemplateData(rw, entry)
		cmdCode := g.generateWorkflowCmdCode(cmdData)
		code.WriteString(cmdCode)
		code.WriteString("\n")
	}
}

// convertWorkflowCmdToTemplateData converts a workflow cmd entry to template data.
func (g *Generator) convertWorkflowCmdToTemplateData(rw *resolve.ResolvedCmdWorkflow, entry *resolve.ResolvedCmdEntry) templates.WorkflowCmdData {
	// Collect all artifacts from Cmd and dependent steps
	allArtifacts := entry.Cmd.Spec.Artifacts
	for _, step := range entry.BeforeSteps {
		allArtifacts = append(allArtifacts, step.Step.Spec.Artifacts...)
	}
	for _, step := range entry.AfterSteps {
		allArtifacts = append(allArtifacts, step.Step.Spec.Artifacts...)
	}
	// Also include global workflow steps
	for _, step := range rw.BeforeSteps {
		allArtifacts = append(allArtifacts, step.Step.Spec.Artifacts...)
	}
	for _, step := range rw.AfterSteps {
		allArtifacts = append(allArtifacts, step.Step.Spec.Artifacts...)
	}
	// Sort for determinism
	sort.Strings(allArtifacts)

	data := templates.WorkflowCmdData{
		CmdName:   entry.Cmd.Metadata.CommandName,
		RunScript: strings.TrimSpace(entry.Cmd.Spec.Run),
		Artifacts: allArtifacts,
		Env:       formatCmdEnv(resolver.ResolveMap(entry.Cmd.Spec.Env)), // Cmd env: direct assignment
	}

	// Global before steps
	if len(rw.BeforeSteps) > 0 {
		data.HasGlobalBefore = true
		data.GlobalBeforeSteps = make([]templates.WorkflowStepData, len(rw.BeforeSteps))
		for i, step := range rw.BeforeSteps {
			whenCondition := ""
			if step.When != nil {
				whenCondition = g.generateWhenCondition(step.When)
			}
			data.GlobalBeforeSteps[i] = templates.WorkflowStepData{
				StepName:      step.Step.Metadata.Name,
				IgnoreFailure: step.FailurePolicy == "Ignore",
				WhenCondition: whenCondition,
				Env:           formatEnvDefaults(resolver.ResolveMap(step.Env)), // Resolve env placeholders
			}
		}
	}

	// Global after steps
	if len(rw.AfterSteps) > 0 {
		data.HasGlobalAfter = true
		data.GlobalAfterSteps = make([]templates.WorkflowStepData, len(rw.AfterSteps))
		for i, step := range rw.AfterSteps {
			whenCondition := ""
			if step.When != nil {
				whenCondition = g.generateWhenCondition(step.When)
			}
			data.GlobalAfterSteps[i] = templates.WorkflowStepData{
				StepName:      step.Step.Metadata.Name,
				IgnoreFailure: step.FailurePolicy == "Ignore",
				WhenCondition: whenCondition,
				Env:           formatEnvDefaults(resolver.ResolveMap(step.Env)), // Resolve env placeholders
			}
		}
	}

	// Per-cmd before steps
	if len(entry.BeforeSteps) > 0 {
		data.HasCmdBefore = true
		data.CmdBeforeSteps = make([]templates.WorkflowStepData, len(entry.BeforeSteps))
		for i, step := range entry.BeforeSteps {
			whenCondition := ""
			if step.When != nil {
				whenCondition = g.generateWhenCondition(step.When)
			}
			data.CmdBeforeSteps[i] = templates.WorkflowStepData{
				StepName:      step.Step.Metadata.Name,
				IgnoreFailure: step.FailurePolicy == "Ignore",
				WhenCondition: whenCondition,
				Env:           formatEnvDefaults(resolver.ResolveMap(step.Env)), // Resolve env placeholders
			}
		}
	}

	// Per-cmd after steps
	if len(entry.AfterSteps) > 0 {
		data.HasCmdAfter = true
		data.CmdAfterSteps = make([]templates.WorkflowStepData, len(entry.AfterSteps))
		for i, step := range entry.AfterSteps {
			whenCondition := ""
			if step.When != nil {
				whenCondition = g.generateWhenCondition(step.When)
			}
			data.CmdAfterSteps[i] = templates.WorkflowStepData{
				StepName:      step.Step.Metadata.Name,
				IgnoreFailure: step.FailurePolicy == "Ignore",
				WhenCondition: whenCondition,
				Env:           formatEnvDefaults(resolver.ResolveMap(step.Env)), // Resolve env placeholders
			}
		}
	}

	return data
}

// generateWorkflowCmdCode generates the shell code for a workflow cmd.
func (g *Generator) generateWorkflowCmdCode(data templates.WorkflowCmdData) string {
	var code strings.Builder

	code.WriteString(data.CmdName + "() {\n")
	code.WriteString("    local __kfg_status=0\n")
	code.WriteString("\n")

	// Generate session ID for this invocation (timestamp-random format)
	code.WriteString("    # Generate session ID for this invocation (timestamp-random format)\n")
	code.WriteString("    KFG_SESSION_ID=\"$(date +%s)-$RANDOM\"\n")
	code.WriteString("    export KFG_SESSION_ID\n")
	code.WriteString("\n")

	// Reset context
	code.WriteString("    # Reset context\n")
	code.WriteString("    __kfg_ctx_reset\n")
	code.WriteString("\n")

	// Cmd env exports (resolved placeholders)
	if len(data.Env) > 0 {
		code.WriteString("    # Cmd environment variables\n")
		// Sort env keys for deterministic output
		keys := make([]string, 0, len(data.Env))
		for k := range data.Env {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			code.WriteString("    export " + k + "=\"" + data.Env[k] + "\"\n")
		}
		code.WriteString("\n")
	}

	// Global before steps
	if data.HasGlobalBefore {
		code.WriteString("    # Global before steps\n")
		for _, step := range data.GlobalBeforeSteps {
			g.generateStepCall(&code, step, 4)
		}
		code.WriteString("\n")
	}

	// Cmd-specific before steps
	if data.HasCmdBefore {
		code.WriteString("    # Cmd before steps\n")
		for _, step := range data.CmdBeforeSteps {
			g.generateStepCall(&code, step, 4)
		}
		code.WriteString("\n")
	}

	// Main body
	code.WriteString("    # Main body\n")
	code.WriteString("    " + data.RunScript + "\n")
	code.WriteString("    __kfg_status=$?\n")
	code.WriteString("\n")

	// Cmd-specific after steps
	if data.HasCmdAfter {
		code.WriteString("    # Cmd after steps\n")
		for _, step := range data.CmdAfterSteps {
			g.generateAfterStepCall(&code, step, 4)
		}
		code.WriteString("\n")
	}

	// Global after steps
	if data.HasGlobalAfter {
		code.WriteString("    # Global after steps\n")
		for _, step := range data.GlobalAfterSteps {
			g.generateAfterStepCall(&code, step, 4)
		}
		code.WriteString("\n")
	}

	code.WriteString("    return $__kfg_status\n")
	code.WriteString("}\n")
	code.WriteString("export -f " + data.CmdName + "\n")

	return code.String()
}

// generateStepCall generates a step call for before steps.
func (g *Generator) generateStepCall(code *strings.Builder, step templates.WorkflowStepData, indent int) {
	spaces := strings.Repeat("    ", indent)

	// Check if we need subshell for env
	hasEnv := len(step.Env) > 0

	if step.WhenCondition != "" {
		code.WriteString(spaces + "if " + step.WhenCondition + "; then\n")
		if hasEnv {
			// Wrap in subshell with env exports
			code.WriteString(spaces + "    (\n")
			// Sort env keys for deterministic output
			keys := make([]string, 0, len(step.Env))
			for k := range step.Env {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				code.WriteString(spaces + "        export " + k + "=\"" + step.Env[k] + "\"\n")
			}
			code.WriteString(spaces + "        __kfg_run_step_" + step.StepName + "\n")
			code.WriteString(spaces + "    )")
			if step.IgnoreFailure {
				code.WriteString(" || true\n")
			} else {
				code.WriteString(" || return $?\n")
			}
		} else {
			if step.IgnoreFailure {
				code.WriteString(spaces + "    __kfg_run_step_" + step.StepName + " || true\n")
			} else {
				code.WriteString(spaces + "    __kfg_run_step_" + step.StepName + " || return $?\n")
			}
		}
		code.WriteString(spaces + "fi\n")
	} else {
		if hasEnv {
			// Wrap in subshell with env exports
			code.WriteString(spaces + "(\n")
			// Sort env keys for deterministic output
			keys := make([]string, 0, len(step.Env))
			for k := range step.Env {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				code.WriteString(spaces + "    export " + k + "=\"" + step.Env[k] + "\"\n")
			}
			code.WriteString(spaces + "    __kfg_run_step_" + step.StepName + "\n")
			code.WriteString(spaces + ")")
			if step.IgnoreFailure {
				code.WriteString(" || true\n")
			} else {
				code.WriteString(" || return $?\n")
			}
		} else {
			if step.IgnoreFailure {
				code.WriteString(spaces + "__kfg_run_step_" + step.StepName + " || true\n")
			} else {
				code.WriteString(spaces + "__kfg_run_step_" + step.StepName + " || return $?\n")
			}
		}
	}
}

// generateAfterStepCall generates a step call for after steps.
func (g *Generator) generateAfterStepCall(code *strings.Builder, step templates.WorkflowStepData, indent int) {
	spaces := strings.Repeat("    ", indent)

	// Check if we need subshell for env
	hasEnv := len(step.Env) > 0

	if step.WhenCondition != "" {
		code.WriteString(spaces + "if " + step.WhenCondition + "; then\n")
		if hasEnv {
			// Wrap in subshell with env exports
			code.WriteString(spaces + "    (\n")
			// Sort env keys for deterministic output
			keys := make([]string, 0, len(step.Env))
			for k := range step.Env {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				code.WriteString(spaces + "        export " + k + "=\"" + step.Env[k] + "\"\n")
			}
			code.WriteString(spaces + "        __kfg_run_step_" + step.StepName + "\n")
			code.WriteString(spaces + "    )")
			if step.IgnoreFailure {
				code.WriteString(" || true\n")
			} else {
				code.WriteString("\n")
			}
		} else {
			if step.IgnoreFailure {
				code.WriteString(spaces + "    __kfg_run_step_" + step.StepName + " || true\n")
			} else {
				code.WriteString(spaces + "    [ $__kfg_status -eq 0 ] && __kfg_run_step_" + step.StepName + "\n")
			}
		}
		code.WriteString(spaces + "fi\n")
	} else {
		if hasEnv {
			// Wrap in subshell with env exports
			code.WriteString(spaces + "(\n")
			// Sort env keys for deterministic output
			keys := make([]string, 0, len(step.Env))
			for k := range step.Env {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				code.WriteString(spaces + "    export " + k + "=\"" + step.Env[k] + "\"\n")
			}
			code.WriteString(spaces + "    __kfg_run_step_" + step.StepName + "\n")
			code.WriteString(spaces + ")")
			if step.IgnoreFailure {
				code.WriteString(" || true\n")
			} else {
				code.WriteString("\n")
			}
		} else {
			if step.IgnoreFailure {
				code.WriteString(spaces + "__kfg_run_step_" + step.StepName + " || true\n")
			} else {
				code.WriteString(spaces + "[ $__kfg_status -eq 0 ] && __kfg_run_step_" + step.StepName + "\n")
			}
		}
	}
}

// formatEnvDefaults formats environment variables as bash default expansions.
// Each value becomes "${KEY:-value}" so caller-set values take precedence.
// Used for step-level env where override capability is needed.
func formatEnvDefaults(env map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range env {
		result[k] = "${" + k + ":-" + v + "}"
	}
	return result
}

// formatCmdEnv formats environment variables for cmd-level export using direct
// assignment. Unlike formatEnvDefaults, this does NOT use ${KEY:-default}
// expansion, because cmd env must not leak between sequential cmd invocations
// in the same shell (the first cmd's exported value would persist).
func formatCmdEnv(env map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range env {
		result[k] = v
	}
	return result
}

// convertStepToTemplateData converts a Step to StepData.
func (g *Generator) convertStepToTemplateData(step *manifest.Step) templates.StepData {
	data := templates.StepData{
		StepName:      step.Metadata.Name,
		HasOutput:     step.Spec.Output != nil,
		Artifacts:     step.Spec.Artifacts,
		IgnoreFailure: false,
		Env:           formatEnvDefaults(resolver.ResolveMap(step.Spec.Env)), // Resolve env placeholders
	}

	if step.Spec.Output != nil {
		data.OutputName = step.Spec.Output.Name
		data.RunScript = strings.TrimSpace(step.Spec.Run)
	} else {
		runScript := strings.TrimSpace(step.Spec.Run)
		data.IsMultiLine = strings.Contains(runScript, "\n")

		if data.IsMultiLine {
			lines := strings.Split(runScript, "\n")
			data.RunLines = lines
		} else {
			data.RunScript = runScript
		}
	}

	return data
}

// generateWhenCondition generates a shell condition from a WhenClause.
func (g *Generator) generateWhenCondition(when *manifest.WhenClause) string {
	if when.Output != nil {
		return g.generateOutputCondition(when.Output)
	}
	if len(when.AllOf) > 0 {
		conditions := make([]string, len(when.AllOf))
		for i, c := range when.AllOf {
			conditions[i] = g.generateWhenCondition(&c)
		}
		return "__kfg_when_allof " + strings.Join(conditions, " ")
	}
	if len(when.AnyOf) > 0 {
		conditions := make([]string, len(when.AnyOf))
		for i, c := range when.AnyOf {
			conditions[i] = g.generateWhenCondition(&c)
		}
		return "__kfg_when_anyof " + strings.Join(conditions, " ")
	}
	if when.Not != nil {
		return "__kfg_when_not " + g.generateWhenCondition(when.Not)
	}
	return ""
}

// generateOutputCondition generates a shell condition from an OutputCondition.
func (g *Generator) generateOutputCondition(output *manifest.OutputCondition) string {
	step := output.Step
	name := output.Name

	if output.Equals != "" {
		return "__kfg_when_equals \"" + step + "\" \"" + name + "\" \"" + output.Equals + "\""
	}
	if len(output.In) > 0 {
		args := make([]string, len(output.In)+2)
		args[0] = "\"" + step + "\""
		args[1] = "\"" + name + "\""
		for i, v := range output.In {
			args[i+2] = "\"" + v + "\""
		}
		return "__kfg_when_in " + strings.Join(args, " ")
	}
	if output.Contains != "" {
		return "__kfg_when_contains \"" + step + "\" \"" + name + "\" \"" + output.Contains + "\""
	}
	if output.Matches != "" {
		return "__kfg_when_matches \"" + step + "\" \"" + name + "\" \"" + output.Matches + "\""
	}
	return ""
}
