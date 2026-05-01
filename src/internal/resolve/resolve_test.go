package resolve

import (
	"testing"

	"github.com/seregatte/kfg/src/internal/manifest"

	"github.com/stretchr/testify/assert"
)

func TestNewIndex(t *testing.T) {
	step := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "test-step"},
		Spec:       manifest.StepSpec{Run: "echo test"},
	}

	cmd := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "test-cmd", CommandName: "testcmd"},
		Spec:       manifest.CmdSpec{Run: "echo test"},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "test-workflow"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"test-cmd"}},
	}

	resources := []manifest.ParsedResource{
		{Step: step},
		{Cmd: cmd},
		{CmdWorkflow: workflow},
	}

	index := NewIndex(resources)

	assert.NotNil(t, index)
	assert.Len(t, index.steps, 1)
	assert.Len(t, index.cmds, 1)
	assert.Len(t, index.cmdWorkflows, 1)
}

func TestIndexGetStep(t *testing.T) {
	step := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "my-step"},
		Spec:       manifest.StepSpec{Run: "echo test"},
	}

	index := NewIndex([]manifest.ParsedResource{{Step: step}})

	found, ok := index.GetStep("my-step")
	assert.True(t, ok)
	assert.Equal(t, step, found)

	_, ok = index.GetStep("nonexistent")
	assert.False(t, ok)
}

func TestIndexGetCmd(t *testing.T) {
	cmd := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "my-cmd", CommandName: "myCmd"},
		Spec:       manifest.CmdSpec{Run: "echo test"},
	}

	index := NewIndex([]manifest.ParsedResource{{Cmd: cmd}})

	found, ok := index.GetCmd("my-cmd")
	assert.True(t, ok)
	assert.Equal(t, cmd, found)

	_, ok = index.GetCmd("nonexistent")
	assert.False(t, ok)
}

func TestIndexGetCmdWorkflow(t *testing.T) {
	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "my-workflow"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"cmd1"}},
	}

	index := NewIndex([]manifest.ParsedResource{{CmdWorkflow: workflow}})

	found, ok := index.GetCmdWorkflow("my-workflow")
	assert.True(t, ok)
	assert.Equal(t, workflow, found)

	_, ok = index.GetCmdWorkflow("nonexistent")
	assert.False(t, ok)
}

func TestResolveKustomization(t *testing.T) {
	step := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "setup"},
		Spec:       manifest.StepSpec{Run: "echo setup"},
	}

	cmd := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "build", CommandName: "build"},
		Spec:       manifest.CmdSpec{Run: "echo build"},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "dev", Shell: "bash"},
		Spec: manifest.CmdWorkflowSpec{
			Cmds:   []string{"build"},
			Before: []manifest.StepReference{{Step: "setup"}},
		},
	}

	index := NewIndex([]manifest.ParsedResource{
		{Step: step},
		{Cmd: cmd},
		{CmdWorkflow: workflow},
	})

	resolver := NewResolver(index)
	resolved, err := resolver.ResolveKustomization("", nil)

	assert.NoError(t, err)
	assert.NotNil(t, resolved)
	assert.Equal(t, "dev", resolved.Name)
	assert.Equal(t, "bash", resolved.Shell)
	assert.NotNil(t, resolved.Workflow)
	assert.Len(t, resolved.Workflow.Cmds, 1)
	assert.Len(t, resolved.Workflow.BeforeSteps, 1)
}

func TestResolveKustomizationByName(t *testing.T) {
	workflow1 := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "dev"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{}},
	}

	workflow2 := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "prod"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{}},
	}

	index := NewIndex([]manifest.ParsedResource{
		{CmdWorkflow: workflow1},
		{CmdWorkflow: workflow2},
	})

	resolver := NewResolver(index)

	// Resolve by name
	resolved, err := resolver.ResolveKustomization("prod", nil)
	assert.NoError(t, err)
	assert.Equal(t, "prod", resolved.Name)

	// Should fail without workflow name (multiple workflows)
	_, err = resolver.ResolveKustomization("", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "multiple CmdWorkflows found")
}

func TestResolveCmdFilter(t *testing.T) {
	cmd1 := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "cmd1", CommandName: "cmd1"},
		Spec:       manifest.CmdSpec{Run: "echo 1"},
	}

	cmd2 := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "cmd2", CommandName: "cmd2"},
		Spec:       manifest.CmdSpec{Run: "echo 2"},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "test"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"cmd1", "cmd2"}},
	}

	index := NewIndex([]manifest.ParsedResource{
		{Cmd: cmd1},
		{Cmd: cmd2},
		{CmdWorkflow: workflow},
	})

	resolver := NewResolver(index)

	// Filter to only cmd1
	resolved, err := resolver.ResolveKustomization("test", []string{"cmd1"})
	assert.NoError(t, err)
	assert.Len(t, resolved.Workflow.Cmds, 1)
	assert.NotNil(t, resolved.Workflow.Cmds["cmd1"])

	// Invalid filter
	_, err = resolver.ResolveKustomization("test", []string{"nonexistent"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Cmd not in workflow")
}

func TestResolvedCmdWorkflowGetAllCmdNames(t *testing.T) {
	cmd1 := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "cmd1", CommandName: "cmd1"},
		Spec:       manifest.CmdSpec{Run: "echo 1"},
	}

	cmd2 := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "cmd2", CommandName: "cmd2"},
		Spec:       manifest.CmdSpec{Run: "echo 2"},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "test"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"cmd1", "cmd2"}},
	}

	index := NewIndex([]manifest.ParsedResource{
		{Cmd: cmd1},
		{Cmd: cmd2},
		{CmdWorkflow: workflow},
	})

	resolver := NewResolver(index)
	resolved, err := resolver.ResolveKustomization("test", nil)

	assert.NoError(t, err)
	names := resolved.Workflow.GetAllCmdNames()
	assert.Len(t, names, 2)
	assert.Contains(t, names, "cmd1")
	assert.Contains(t, names, "cmd2")
}

func TestNewIndex_MixedResources(t *testing.T) {
	// Create all three execution kinds
	step := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "setup"},
		Spec:       manifest.StepSpec{Run: "echo setup"},
	}

	cmd := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "build", CommandName: "build"},
		Spec:       manifest.CmdSpec{Run: "echo build"},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "dev"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"build"}},
	}

	// Create ParsedResource with all three execution kinds
	resources := []manifest.ParsedResource{
		{Step: step},
		{Cmd: cmd},
		{CmdWorkflow: workflow},
	}

	index := NewIndex(resources)

	// Verify all execution kinds are indexed
	assert.Len(t, index.steps, 1)
	assert.Len(t, index.cmds, 1)
	assert.Len(t, index.cmdWorkflows, 1)

	// Verify Step, Cmd, CmdWorkflow are indexed
	_, ok := index.GetStep("setup")
	assert.True(t, ok)
	_, ok = index.GetCmd("build")
	assert.True(t, ok)
	_, ok = index.GetCmdWorkflow("dev")
	assert.True(t, ok)

	t.Logf("Indexed %d execution kinds", 3)
}

func TestResolver_ExecutionKindsOnly(t *testing.T) {
	// Create execution kinds
	step := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "setup"},
		Spec:       manifest.StepSpec{Run: "echo setup"},
	}

	cmd := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "build", CommandName: "build"},
		Spec:       manifest.CmdSpec{Run: "echo build"},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "dev"},
		Spec: manifest.CmdWorkflowSpec{
			Cmds:   []string{"build"},
			Before: []manifest.StepReference{{Step: "setup"}},
		},
	}

	// Create index with execution kinds only
	index := NewIndex([]manifest.ParsedResource{
		{Step: step},
		{Cmd: cmd},
		{CmdWorkflow: workflow},
	})

	resolver := NewResolver(index)
	resolved, err := resolver.ResolveKustomization("", nil)

	// Verify resolution succeeds
	assert.NoError(t, err)
	assert.NotNil(t, resolved)
	assert.Equal(t, "dev", resolved.Name)

	// Verify resolved workflow has correct structure
	assert.Len(t, resolved.Workflow.Cmds, 1)
	assert.Len(t, resolved.Workflow.BeforeSteps, 1)

	t.Log("Resolver correctly processes execution kinds only")
}

// ============================================================================
// Env Merge Tests
// ============================================================================

func TestMergeEnv_BaseAndOverride(t *testing.T) {
	base := map[string]string{
		"SRC":  "default-src",
		"DEST": "default-dest",
		"MODE": "normal",
	}
	override := map[string]string{
		"DEST":  "override-dest",
		"DEBUG": "true",
	}

	result := MergeEnv(base, override)

	assert.NotNil(t, result)
	assert.Equal(t, "default-src", result["SRC"])
	assert.Equal(t, "override-dest", result["DEST"])
	assert.Equal(t, "normal", result["MODE"])
	assert.Equal(t, "true", result["DEBUG"])
	assert.Len(t, result, 4)
}

func TestMergeEnv_OnlyBase(t *testing.T) {
	base := map[string]string{
		"SRC":  "default-src",
		"DEST": "default-dest",
	}
	var override map[string]string // nil

	result := MergeEnv(base, override)

	assert.NotNil(t, result)
	assert.Equal(t, "default-src", result["SRC"])
	assert.Equal(t, "default-dest", result["DEST"])
	assert.Len(t, result, 2)
}

func TestMergeEnv_OnlyOverride(t *testing.T) {
	var base map[string]string // nil
	override := map[string]string{
		"SRC":   "override-src",
		"DEBUG": "true",
	}

	result := MergeEnv(base, override)

	assert.NotNil(t, result)
	assert.Equal(t, "override-src", result["SRC"])
	assert.Equal(t, "true", result["DEBUG"])
	assert.Len(t, result, 2)
}

func TestMergeEnv_BothNil(t *testing.T) {
	result := MergeEnv(nil, nil)

	assert.Nil(t, result)
}

func TestMergeEnv_BothEmpty(t *testing.T) {
	base := map[string]string{}
	override := map[string]string{}

	result := MergeEnv(base, override)

	assert.Nil(t, result)
}

func TestMergeEnv_OverrideWins(t *testing.T) {
	base := map[string]string{
		"KEY": "base-value",
	}
	override := map[string]string{
		"KEY": "override-value",
	}

	result := MergeEnv(base, override)

	assert.NotNil(t, result)
	assert.Equal(t, "override-value", result["KEY"])
	assert.Len(t, result, 1)
}

func TestResolveStepReferences_EnvPopulation(t *testing.T) {
	step := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "copy-step"},
		Spec: manifest.StepSpec{
			Run: "cp $SRC $DEST",
			Env: map[string]string{
				"SRC":  "docs/AGENTS.md",
				"DEST": "output.md",
			},
		},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "test"},
		Spec: manifest.CmdWorkflowSpec{
			Cmds: []string{},
			Before: []manifest.StepReference{
				{
					Step: "copy-step",
					Env: map[string]string{
						"DEST":  "CLAUDE.md",
						"DEBUG": "true",
					},
				},
			},
		},
	}

	index := NewIndex([]manifest.ParsedResource{
		{Step: step},
		{CmdWorkflow: workflow},
	})

	resolver := NewResolver(index)
	resolved, err := resolver.ResolveKustomization("test", nil)

	assert.NoError(t, err)
	assert.Len(t, resolved.Workflow.BeforeSteps, 1)

	resolvedStep := resolved.Workflow.BeforeSteps[0]
	assert.NotNil(t, resolvedStep.Env)
	assert.Equal(t, "docs/AGENTS.md", resolvedStep.Env["SRC"])
	assert.Equal(t, "CLAUDE.md", resolvedStep.Env["DEST"])
	assert.Equal(t, "true", resolvedStep.Env["DEBUG"])
}

func TestResolveStepReferences_EnvNoOverride(t *testing.T) {
	step := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "simple-step"},
		Spec: manifest.StepSpec{
			Run: "echo $VAR",
			Env: map[string]string{
				"VAR": "default-value",
			},
		},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "test"},
		Spec: manifest.CmdWorkflowSpec{
			Cmds: []string{},
			Before: []manifest.StepReference{
				{Step: "simple-step"}, // No env override
			},
		},
	}

	index := NewIndex([]manifest.ParsedResource{
		{Step: step},
		{CmdWorkflow: workflow},
	})

	resolver := NewResolver(index)
	resolved, err := resolver.ResolveKustomization("test", nil)

	assert.NoError(t, err)
	assert.Len(t, resolved.Workflow.BeforeSteps, 1)

	resolvedStep := resolved.Workflow.BeforeSteps[0]
	assert.NotNil(t, resolvedStep.Env)
	assert.Equal(t, "default-value", resolvedStep.Env["VAR"])
}

func TestResolveStepReferences_NoEnv(t *testing.T) {
	step := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "no-env-step"},
		Spec:       manifest.StepSpec{Run: "echo hello"},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "test"},
		Spec: manifest.CmdWorkflowSpec{
			Cmds: []string{},
			Before: []manifest.StepReference{
				{Step: "no-env-step"},
			},
		},
	}

	index := NewIndex([]manifest.ParsedResource{
		{Step: step},
		{CmdWorkflow: workflow},
	})

	resolver := NewResolver(index)
	resolved, err := resolver.ResolveKustomization("test", nil)

	assert.NoError(t, err)
	assert.Len(t, resolved.Workflow.BeforeSteps, 1)

	resolvedStep := resolved.Workflow.BeforeSteps[0]
	assert.Nil(t, resolvedStep.Env)
}

// ============================================================================
// ResolveAllWorkflows Tests
// ============================================================================

func TestResolveAllWorkflows_Single(t *testing.T) {
	cmd := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "build", CommandName: "build"},
		Spec:       manifest.CmdSpec{Run: "echo build"},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "dev"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"build"}},
	}

	index := NewIndex([]manifest.ParsedResource{
		{Cmd: cmd},
		{CmdWorkflow: workflow},
	})

	resolver := NewResolver(index)
	workflows, err := resolver.ResolveAllWorkflows()

	assert.NoError(t, err)
	assert.Len(t, workflows, 1)
	assert.Equal(t, "dev", workflows[0].Workflow.Metadata.Name)
	assert.Len(t, workflows[0].Cmds, 1)
}

func TestResolveAllWorkflows_Multiple(t *testing.T) {
	cmd1 := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "build", CommandName: "build"},
		Spec:       manifest.CmdSpec{Run: "echo build"},
	}

	cmd2 := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "deploy", CommandName: "deploy"},
		Spec:       manifest.CmdSpec{Run: "echo deploy"},
	}

	workflow1 := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "dev"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"build"}},
	}

	workflow2 := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "prod"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"deploy"}},
	}

	index := NewIndex([]manifest.ParsedResource{
		{Cmd: cmd1},
		{Cmd: cmd2},
		{CmdWorkflow: workflow1},
		{CmdWorkflow: workflow2},
	})

	resolver := NewResolver(index)
	workflows, err := resolver.ResolveAllWorkflows()

	assert.NoError(t, err)
	assert.Len(t, workflows, 2)

	// Verify both workflows are resolved
	workflowNames := make(map[string]bool)
	for _, w := range workflows {
		workflowNames[w.Workflow.Metadata.Name] = true
	}
	assert.True(t, workflowNames["dev"])
	assert.True(t, workflowNames["prod"])
}

func TestResolveAllWorkflows_None(t *testing.T) {
	index := NewIndex([]manifest.ParsedResource{})
	resolver := NewResolver(index)

	_, err := resolver.ResolveAllWorkflows()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no CmdWorkflow found")
}

// ============================================================================
// ResolveWorkflowsByName Tests
// ============================================================================

func TestResolveWorkflowsByName_Single(t *testing.T) {
	cmd := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "build", CommandName: "build"},
		Spec:       manifest.CmdSpec{Run: "echo build"},
	}

	workflow1 := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "dev"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"build"}},
	}

	workflow2 := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "prod"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{}},
	}

	index := NewIndex([]manifest.ParsedResource{
		{Cmd: cmd},
		{CmdWorkflow: workflow1},
		{CmdWorkflow: workflow2},
	})

	resolver := NewResolver(index)
	workflows, err := resolver.ResolveWorkflowsByName([]string{"dev"})

	assert.NoError(t, err)
	assert.Len(t, workflows, 1)
	assert.Equal(t, "dev", workflows[0].Workflow.Metadata.Name)
}

func TestResolveWorkflowsByName_Multiple(t *testing.T) {
	cmd1 := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "build", CommandName: "build"},
		Spec:       manifest.CmdSpec{Run: "echo build"},
	}

	cmd2 := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "deploy", CommandName: "deploy"},
		Spec:       manifest.CmdSpec{Run: "echo deploy"},
	}

	workflow1 := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "dev"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"build"}},
	}

	workflow2 := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "prod"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"deploy"}},
	}

	index := NewIndex([]manifest.ParsedResource{
		{Cmd: cmd1},
		{Cmd: cmd2},
		{CmdWorkflow: workflow1},
		{CmdWorkflow: workflow2},
	})

	resolver := NewResolver(index)
	workflows, err := resolver.ResolveWorkflowsByName([]string{"dev", "prod"})

	assert.NoError(t, err)
	assert.Len(t, workflows, 2)

	// Verify both workflows are resolved
	workflowNames := make(map[string]bool)
	for _, w := range workflows {
		workflowNames[w.Workflow.Metadata.Name] = true
	}
	assert.True(t, workflowNames["dev"])
	assert.True(t, workflowNames["prod"])
}

func TestResolveWorkflowsByName_EmptyList(t *testing.T) {
	cmd := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "build", CommandName: "build"},
		Spec:       manifest.CmdSpec{Run: "echo build"},
	}

	workflow1 := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "dev"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"build"}},
	}

	workflow2 := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "prod"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{}},
	}

	index := NewIndex([]manifest.ParsedResource{
		{Cmd: cmd},
		{CmdWorkflow: workflow1},
		{CmdWorkflow: workflow2},
	})

	resolver := NewResolver(index)
	// Empty list should return all workflows (same as ResolveAllWorkflows)
	workflows, err := resolver.ResolveWorkflowsByName([]string{})

	assert.NoError(t, err)
	assert.Len(t, workflows, 2)
}

func TestResolveWorkflowsByName_NotFound(t *testing.T) {
	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "dev"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{}},
	}

	index := NewIndex([]manifest.ParsedResource{{CmdWorkflow: workflow}})
	resolver := NewResolver(index)

	_, err := resolver.ResolveWorkflowsByName([]string{"nonexistent"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CmdWorkflow not found")
	assert.Contains(t, err.Error(), "Available:")
}

func TestResolveWorkflowsByName_PartialNotFound(t *testing.T) {
	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "dev"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{}},
	}

	index := NewIndex([]manifest.ParsedResource{{CmdWorkflow: workflow}})
	resolver := NewResolver(index)

	_, err := resolver.ResolveWorkflowsByName([]string{"dev", "nonexistent"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CmdWorkflow not found")
}
