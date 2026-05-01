package resolve

import (
	"fmt"
	"sort"
	"strings"

	"github.com/seregatte/kfg/src/internal/manifest"
)

// ============================================================================
// Index
// ============================================================================

// Index indexes resources by name for efficient lookup.
type Index struct {
	steps        map[string]*manifest.Step
	cmds         map[string]*manifest.Cmd
	cmdWorkflows map[string]*manifest.CmdWorkflow
}

// NewIndex creates a new Index from a list of ParsedResources.
// Only execution kinds (Step, Cmd, CmdWorkflow) are indexed.
// Source kinds (Assets, Converter) are explicitly skipped as they are not part of the execution model.
func NewIndex(resources []manifest.ParsedResource) *Index {
	idx := &Index{
		steps:        make(map[string]*manifest.Step),
		cmds:         make(map[string]*manifest.Cmd),
		cmdWorkflows: make(map[string]*manifest.CmdWorkflow),
	}

	for _, res := range resources {
		switch {
		case res.Step != nil:
			idx.steps[res.Step.Metadata.Name] = res.Step
		case res.Cmd != nil:
			idx.cmds[res.Cmd.Metadata.Name] = res.Cmd
		case res.CmdWorkflow != nil:
			idx.cmdWorkflows[res.CmdWorkflow.Metadata.Name] = res.CmdWorkflow
		default:
			// Unknown or empty ParsedResource - skip
		}
	}

	return idx
}

// GetStep returns a specific Step by name.
func (idx *Index) GetStep(name string) (*manifest.Step, bool) {
	step, ok := idx.steps[name]
	return step, ok
}

// GetSteps returns all Step resources.
func (idx *Index) GetSteps() []*manifest.Step {
	steps := make([]*manifest.Step, 0, len(idx.steps))
	for _, step := range idx.steps {
		steps = append(steps, step)
	}
	return steps
}

// GetCmd returns a specific Cmd by name.
func (idx *Index) GetCmd(name string) (*manifest.Cmd, bool) {
	cmd, ok := idx.cmds[name]
	return cmd, ok
}

// GetCmds returns all Cmd resources.
func (idx *Index) GetCmds() []*manifest.Cmd {
	cmds := make([]*manifest.Cmd, 0, len(idx.cmds))
	for _, cmd := range idx.cmds {
		cmds = append(cmds, cmd)
	}
	return cmds
}

// GetCmdWorkflow returns a specific CmdWorkflow by name.
func (idx *Index) GetCmdWorkflow(name string) (*manifest.CmdWorkflow, bool) {
	workflow, ok := idx.cmdWorkflows[name]
	return workflow, ok
}

// GetCmdWorkflows returns all CmdWorkflow resources.
func (idx *Index) GetCmdWorkflows() []*manifest.CmdWorkflow {
	workflows := make([]*manifest.CmdWorkflow, 0, len(idx.cmdWorkflows))
	for _, workflow := range idx.cmdWorkflows {
		workflows = append(workflows, workflow)
	}
	return workflows
}

// ============================================================================
// Resolver
// ============================================================================

// Resolver resolves dependencies and prepares resources for shell generation.
type Resolver struct {
	index *Index
}

// NewResolver creates a new Resolver with an Index.
func NewResolver(index *Index) *Resolver {
	return &Resolver{
		index: index,
	}
}

// ============================================================================
// Resolved Types
// ============================================================================

// ResolvedStep represents a resolved step.
// Execution order is determined by YAML order.
type ResolvedStep struct {
	Step          *manifest.Step
	When          *manifest.WhenClause
	FailurePolicy string            // "Fail" (default) or "Ignore"
	Env           map[string]string // Merged env: StepSpec.Env + StepReference.Env
}

// ResolvedCmd represents a resolved Cmd (pure function, no before/after).
type ResolvedCmd struct {
	Cmd *manifest.Cmd
}

// ResolvedCmdEntry represents a Cmd with its per-command steps.
type ResolvedCmdEntry struct {
	Cmd         *manifest.Cmd
	BeforeSteps []ResolvedStep // Per-cmd steps (if any)
	AfterSteps  []ResolvedStep // Per-cmd steps (if any)
}

// ResolvedCmdWorkflow represents a fully resolved CmdWorkflow.
// Global before/after steps apply to ALL cmds in the workflow.
type ResolvedCmdWorkflow struct {
	Workflow    *manifest.CmdWorkflow
	Cmds        map[string]*ResolvedCmdEntry // cmd name -> resolved entry
	BeforeSteps []ResolvedStep               // Global before (all cmds)
	AfterSteps  []ResolvedStep               // Global after (all cmds)
	Shell       string                       // Shell from metadata.shell
}

// ResolvedKustomization represents the resolved output for shell generation.
type ResolvedKustomization struct {
	Name     string                    // Kustomization name (directory name)
	Shell    string                    // Shell type from workflow
	Workflow *ResolvedCmdWorkflow      // The resolved workflow
	Steps    map[string]*manifest.Step // All available steps
	Cmds     map[string]*manifest.Cmd  // All available cmds
}

// ============================================================================
// Resolution Methods
// ============================================================================

// ResolveKustomization resolves a kustomization for shell generation.
// It finds the CmdWorkflow (by name or the only workflow) and applies optional cmd filter.
func (r *Resolver) ResolveKustomization(workflowName string, cmdFilter []string) (*ResolvedKustomization, error) {
	// Find workflow
	var workflow *manifest.CmdWorkflow
	if workflowName != "" {
		w, ok := r.index.GetCmdWorkflow(workflowName)
		if !ok {
			return nil, ResolutionError{
				ResourceName: workflowName,
				ResourceKind: "CmdWorkflow",
				Message:      "CmdWorkflow not found",
			}
		}
		workflow = w
	} else {
		// Find the only workflow in the index
		workflows := r.index.GetCmdWorkflows()
		if len(workflows) == 0 {
			return nil, ResolutionError{
				ResourceName: "",
				ResourceKind: "CmdWorkflow",
				Message:      "no CmdWorkflow found in kustomization",
			}
		}
		if len(workflows) > 1 {
			names := make([]string, len(workflows))
			for i, w := range workflows {
				names[i] = w.Metadata.Name
			}
			return nil, ResolutionError{
				ResourceName: "",
				ResourceKind: "CmdWorkflow",
				Message:      "multiple CmdWorkflows found, specify one: " + strings.Join(names, ", "),
			}
		}
		workflow = workflows[0]
	}

	// Resolve the workflow with optional cmd filter
	resolvedWorkflow, err := r.ResolveCmdWorkflow(workflow.Metadata.Name, cmdFilter)
	if err != nil {
		return nil, err
	}

	return &ResolvedKustomization{
		Name:     workflow.Metadata.Name,
		Shell:    resolvedWorkflow.Shell,
		Workflow: resolvedWorkflow,
		Steps:    r.index.steps,
		Cmds:     r.index.cmds,
	}, nil
}

// ResolveCmdWorkflow resolves a CmdWorkflow with optional cmd filter.
func (r *Resolver) ResolveCmdWorkflow(name string, cmdFilter []string) (*ResolvedCmdWorkflow, error) {
	workflow, ok := r.index.GetCmdWorkflow(name)
	if !ok {
		return nil, ResolutionError{
			ResourceName: name,
			ResourceKind: "CmdWorkflow",
			Message:      "CmdWorkflow not found",
		}
	}

	// Determine shell (default to "bash")
	shell := workflow.Metadata.Shell
	if shell == "" {
		shell = "bash"
	}

	// Resolve global before/after steps (preserve YAML order)
	beforeSteps, err := r.ResolveStepReferences(workflow.Spec.Before)
	if err != nil {
		return nil, err
	}
	afterSteps, err := r.ResolveStepReferences(workflow.Spec.After)
	if err != nil {
		return nil, err
	}

	// Determine which cmds to include
	cmdsToInclude := workflow.Spec.Cmds
	if len(cmdFilter) > 0 {
		// Validate filter against workflow cmds
		for _, filteredCmd := range cmdFilter {
			found := false
			for _, workflowCmd := range workflow.Spec.Cmds {
				if workflowCmd == filteredCmd {
					found = true
					break
				}
			}
			if !found {
				return nil, ResolutionError{
					ResourceName: filteredCmd,
					ResourceKind: "Cmd",
					Message:      "Cmd not in workflow. Available cmds: " + strings.Join(workflow.Spec.Cmds, ", "),
				}
			}
		}
		cmdsToInclude = cmdFilter
	}

	// Resolve each cmd
	resolvedCmds := make(map[string]*ResolvedCmdEntry)
	for _, cmdName := range cmdsToInclude {
		cmd, ok := r.index.GetCmd(cmdName)
		if !ok {
			return nil, ResolutionError{
				ResourceName: cmdName,
				ResourceKind: "Cmd",
				Message:      "Cmd not found",
			}
		}

		// For now, each cmd has no per-cmd steps (those come from patches)
		resolvedCmds[cmdName] = &ResolvedCmdEntry{
			Cmd:         cmd,
			BeforeSteps: nil,
			AfterSteps:  nil,
		}
	}

	// Create resolved workflow
	result := &ResolvedCmdWorkflow{
		Workflow:    workflow,
		Cmds:        resolvedCmds,
		BeforeSteps: beforeSteps,
		AfterSteps:  afterSteps,
		Shell:       shell,
	}

	return result, nil
}

// ComputeAllHashes computes hashes for all Cmds in the workflow after resolution.
// Must be called after ResolveCmdWorkflow completes.

// ResolveStepReferences resolves a list of step references.
// Returns resolved steps in YAML order (no weight sorting).
func (r *Resolver) ResolveStepReferences(refs []manifest.StepReference) ([]ResolvedStep, error) {
	var resolved []ResolvedStep

	for _, ref := range refs {
		step, ok := r.index.GetStep(ref.Step)
		if !ok {
			return nil, ResolutionError{
				ResourceName: ref.Step,
				ResourceKind: "Step",
				Message:      "Step not found",
			}
		}

		failurePolicy := ref.FailurePolicy
		if failurePolicy == "" {
			failurePolicy = "Fail" // Default
		}

		// Merge step default env with reference override env
		mergedEnv := MergeEnv(step.Spec.Env, ref.Env)

		resolved = append(resolved, ResolvedStep{
			Step:          step,
			When:          ref.When,
			FailurePolicy: failurePolicy,
			Env:           mergedEnv,
		})
	}

	return resolved, nil
}

// ResolveCmd resolves a single Cmd by name.
func (r *Resolver) ResolveCmd(name string) (*ResolvedCmd, error) {
	cmd, ok := r.index.GetCmd(name)
	if !ok {
		return nil, ResolutionError{
			ResourceName: name,
			ResourceKind: "Cmd",
			Message:      "Cmd not found",
		}
	}

	result := &ResolvedCmd{
		Cmd: cmd,
	}
	return result, nil
}

// ResolveAllWorkflows resolves all CmdWorkflows in the index.
// Returns all workflows resolved, or an error if none exist.
func (r *Resolver) ResolveAllWorkflows() ([]*ResolvedCmdWorkflow, error) {
	workflows := r.index.GetCmdWorkflows()
	if len(workflows) == 0 {
		return nil, ResolutionError{
			ResourceName: "",
			ResourceKind: "CmdWorkflow",
			Message:      "no CmdWorkflow found in kustomization",
		}
	}

	results := make([]*ResolvedCmdWorkflow, 0, len(workflows))
	for _, workflow := range workflows {
		resolved, err := r.ResolveCmdWorkflow(workflow.Metadata.Name, nil)
		if err != nil {
			return nil, err
		}
		results = append(results, resolved)
	}

	return results, nil
}

// ResolveWorkflowsByName resolves specific CmdWorkflows by name.
// Returns workflows that match the provided names, or an error if any name is not found.
func (r *Resolver) ResolveWorkflowsByName(names []string) ([]*ResolvedCmdWorkflow, error) {
	if len(names) == 0 {
		return r.ResolveAllWorkflows()
	}

	results := make([]*ResolvedCmdWorkflow, 0, len(names))
	var availableNames []string

	for _, name := range names {
		resolved, err := r.ResolveCmdWorkflow(name, nil)
		if err != nil {
			// Collect available names for error message
			if len(availableNames) == 0 {
				workflows := r.index.GetCmdWorkflows()
				availableNames = make([]string, len(workflows))
				for i, w := range workflows {
					availableNames[i] = w.Metadata.Name
				}
			}
			return nil, ResolutionError{
				ResourceName: name,
				ResourceKind: "CmdWorkflow",
				Message:      "CmdWorkflow not found. Available: " + strings.Join(availableNames, ", "),
			}
		}
		results = append(results, resolved)
	}

	return results, nil
}

// ============================================================================
// Error Types
// ============================================================================

// ResolutionError represents an error during resolution.
type ResolutionError struct {
	ResourceName string
	ResourceKind string
	Message      string
}

// Error returns the error message.
func (e ResolutionError) Error() string {
	return fmt.Sprintf("%s/%s: %s", e.ResourceKind, e.ResourceName, e.Message)
}

// ============================================================================
// Helper Methods
// ============================================================================

// MergeEnv merges base env with override env (shallow merge, override wins).
// Returns nil if both are nil/empty.
func MergeEnv(base, override map[string]string) map[string]string {
	// If both are nil/empty, return nil
	if len(base) == 0 && len(override) == 0 {
		return nil
	}

	// Create result map with base values
	result := make(map[string]string)
	for k, v := range base {
		result[k] = v
	}

	// Override with values from override map
	for k, v := range override {
		result[k] = v
	}

	return result
}

// GetCmdFunctionName returns the bash function name for the Cmd.
func (rc *ResolvedCmd) GetCmdFunctionName() string {
	return rc.Cmd.Metadata.CommandName
}

// GetAllCmdNames returns sorted list of all cmd names in the workflow.
func (rw *ResolvedCmdWorkflow) GetAllCmdNames() []string {
	names := make([]string, 0, len(rw.Cmds))
	for name := range rw.Cmds {
		names = append(names, name)
	}
	sort.Strings(names) // Sort for deterministic output
	return names
}

// GetAllStepNames returns sorted list of all step names in the workflow.
func (rw *ResolvedCmdWorkflow) GetAllStepNames() []string {
	names := make([]string, 0)

	for _, step := range rw.BeforeSteps {
		names = append(names, step.Step.Metadata.Name)
	}

	for _, entry := range rw.Cmds {
		for _, step := range entry.BeforeSteps {
			names = append(names, step.Step.Metadata.Name)
		}
		for _, step := range entry.AfterSteps {
			names = append(names, step.Step.Metadata.Name)
		}
	}

	for _, step := range rw.AfterSteps {
		names = append(names, step.Step.Metadata.Name)
	}

	sort.Strings(names) // Sort for deterministic output
	return names
}
