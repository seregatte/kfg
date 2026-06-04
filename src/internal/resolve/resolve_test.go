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

// ============================================================================
// Cache Merge Tests
// ============================================================================

func TestMergeCache_RefOverride(t *testing.T) {
	// StepReference cache takes precedence
	stepCache := &manifest.CacheConfig{Enabled: boolPtr(true)}
	refCache := &manifest.CacheConfig{Enabled: boolPtr(false)}

	result := MergeCache(stepCache, refCache)

	assert.NotNil(t, result)
	assert.NotNil(t, result.Enabled)
	assert.False(t, *result.Enabled) // ref takes precedence
}

func TestMergeCache_OnlyStepCache(t *testing.T) {
	// Use Step default when no ref override
	stepCache := &manifest.CacheConfig{Enabled: boolPtr(true)}
	var refCache *manifest.CacheConfig // nil

	result := MergeCache(stepCache, refCache)

	assert.NotNil(t, result)
	assert.NotNil(t, result.Enabled)
	assert.True(t, *result.Enabled)
}

func TestMergeCache_OnlyRefCache(t *testing.T) {
	// Use ref cache when step has no cache
	var stepCache *manifest.CacheConfig // nil
	refCache := &manifest.CacheConfig{Enabled: boolPtr(true)}

	result := MergeCache(stepCache, refCache)

	assert.NotNil(t, result)
	assert.NotNil(t, result.Enabled)
	assert.True(t, *result.Enabled)
}

func TestMergeCache_BothNil(t *testing.T) {
	result := MergeCache(nil, nil)

	assert.Nil(t, result)
}

func TestMergeCache_RefEnabledNil(t *testing.T) {
	// Ref cache with nil Enabled should still be used entirely
	stepCache := &manifest.CacheConfig{Enabled: boolPtr(true)}
	refCache := &manifest.CacheConfig{Enabled: nil}

	result := MergeCache(stepCache, refCache)

	assert.NotNil(t, result)
	assert.Nil(t, result.Enabled) // ref's nil Enabled is used
}

// Helper function to create bool pointer
func boolPtr(b bool) *bool {
	return &b
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

// ============================================================================
// StepRefName Tests (StepReference.Name runtime identity)
// ============================================================================

func TestResolveStepReference_StepRefName(t *testing.T) {
	step := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "detect-step"},
		Spec: manifest.StepSpec{
			Run:    "echo detect",
			Output: &manifest.Output{Name: "AGENT"},
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
					Name: "detect-agent", // StepReference.Name (runtime identity)
					Step: "detect-step",  // Step metadata.name
					Env: map[string]string{
						"MODE": "auto",
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

	// Verify StepRefName is populated from StepReference.Name
	resolvedStep := resolved.Workflow.BeforeSteps[0]
	assert.Equal(t, "detect-agent", resolvedStep.Name, "ResolvedStep.Name should be StepReference.Name")
	assert.Equal(t, "detect-step", resolvedStep.Step.Metadata.Name, "Step metadata name should be unchanged")
}

func TestResolveStepReference_StepRefNameWithEnvOverride(t *testing.T) {
	step := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "setup-step"},
		Spec: manifest.StepSpec{
			Run: "echo setup",
			Env: map[string]string{
				"SRC":  "default",
				"DEST": "default",
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
					Name: "setup-claude", // StepReference.Name
					Step: "setup-step",
					Env: map[string]string{
						"DEST": "CLAUDE.md",
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
	assert.Equal(t, "setup-claude", resolvedStep.Name)
	assert.Equal(t, "setup-step", resolvedStep.Step.Metadata.Name)

	// Verify env is merged properly
	assert.Equal(t, "default", resolvedStep.Env["SRC"])
	assert.Equal(t, "CLAUDE.md", resolvedStep.Env["DEST"])
}

func TestResolveStepReference_MultipleSameStepDifferentNames(t *testing.T) {
	// Test that the same Step can be used multiple times with different StepRefNames
	step := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "copy-step"},
		Spec: manifest.StepSpec{
			Run: "cp $SRC $DEST",
			Env: map[string]string{
				"SRC":  "default",
				"DEST": "default",
			},
			Output: &manifest.Output{Name: "RESULT"},
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
					Name: "copy-claude",
					Step: "copy-step",
					Env:  map[string]string{"DEST": "CLAUDE.md"},
				},
				{
					Name: "copy-gemini",
					Step: "copy-step",
					Env:  map[string]string{"DEST": "GEMINI.md"},
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
	assert.Len(t, resolved.Workflow.BeforeSteps, 2)

	// Verify both step references have different names but same step
	assert.Equal(t, "copy-claude", resolved.Workflow.BeforeSteps[0].Name)
	assert.Equal(t, "copy-step", resolved.Workflow.BeforeSteps[0].Step.Metadata.Name)
	assert.Equal(t, "CLAUDE.md", resolved.Workflow.BeforeSteps[0].Env["DEST"])

	assert.Equal(t, "copy-gemini", resolved.Workflow.BeforeSteps[1].Name)
	assert.Equal(t, "copy-step", resolved.Workflow.BeforeSteps[1].Step.Metadata.Name)
	assert.Equal(t, "GEMINI.md", resolved.Workflow.BeforeSteps[1].Env["DEST"])
}

func TestResolveStepReference_AfterStepsWithNames(t *testing.T) {
	step := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "cleanup-step"},
		Spec:       manifest.StepSpec{Run: "rm -rf temp"},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "test"},
		Spec: manifest.CmdWorkflowSpec{
			Cmds: []string{},
			After: []manifest.StepReference{
				{
					Name: "cleanup-final",
					Step: "cleanup-step",
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
	assert.Len(t, resolved.Workflow.AfterSteps, 1)

	// Verify StepRefName is populated in after steps too
	assert.Equal(t, "cleanup-final", resolved.Workflow.AfterSteps[0].Name)
	assert.Equal(t, "cleanup-step", resolved.Workflow.AfterSteps[0].Step.Metadata.Name)
}

// ============================================================================
// Cache Resolution Tests
// ============================================================================

func TestResolveStepReferences_CachePopulation(t *testing.T) {
	step := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "cached-step"},
		Spec: manifest.StepSpec{
			Run:   "echo cached",
			Cache: &manifest.CacheConfig{Enabled: boolPtr(true)},
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
					Name:  "cached-ref",
					Step:  "cached-step",
					Cache: &manifest.CacheConfig{Enabled: boolPtr(false)},
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

	// Verify cache is merged with StepReference precedence
	resolvedStep := resolved.Workflow.BeforeSteps[0]
	assert.NotNil(t, resolvedStep.Cache)
	assert.NotNil(t, resolvedStep.Cache.Enabled)
	assert.False(t, *resolvedStep.Cache.Enabled) // ref overrides step
}

func TestResolveStepReferences_CacheNoOverride(t *testing.T) {
	step := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "cached-step"},
		Spec: manifest.StepSpec{
			Run:   "echo cached",
			Cache: &manifest.CacheConfig{Enabled: boolPtr(true)},
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
					Name: "cached-ref",
					Step: "cached-step",
					// No cache override - should use step default
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

	// Verify step default cache is used
	resolvedStep := resolved.Workflow.BeforeSteps[0]
	assert.NotNil(t, resolvedStep.Cache)
	assert.NotNil(t, resolvedStep.Cache.Enabled)
	assert.True(t, *resolvedStep.Cache.Enabled)
}

func TestResolveStepReferences_NoCache(t *testing.T) {
	step := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "no-cache-step"},
		Spec:       manifest.StepSpec{Run: "echo nocache"},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "test"},
		Spec: manifest.CmdWorkflowSpec{
			Cmds: []string{},
			Before: []manifest.StepReference{
				{
					Name: "no-cache-ref",
					Step: "no-cache-step",
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

	// Verify no cache
	resolvedStep := resolved.Workflow.BeforeSteps[0]
	assert.Nil(t, resolvedStep.Cache)
}

func TestResolveStepReferences_RefCacheOnly(t *testing.T) {
	// Step has no cache, but ref adds cache
	step := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "nocache-step"},
		Spec:       manifest.StepSpec{Run: "echo nocache"},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "test"},
		Spec: manifest.CmdWorkflowSpec{
			Cmds: []string{},
			Before: []manifest.StepReference{
				{
					Name:  "ref-cached",
					Step:  "nocache-step",
					Cache: &manifest.CacheConfig{Enabled: boolPtr(true)},
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

	// Verify ref cache is applied
	resolvedStep := resolved.Workflow.BeforeSteps[0]
	assert.NotNil(t, resolvedStep.Cache)
	assert.NotNil(t, resolvedStep.Cache.Enabled)
	assert.True(t, *resolvedStep.Cache.Enabled)
}
