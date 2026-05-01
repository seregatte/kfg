package generate

import (
	"strings"
	"testing"

	"github.com/seregatte/kfg/src/internal/manifest"
	"github.com/seregatte/kfg/src/internal/resolve"

	"github.com/stretchr/testify/assert"
)

func TestGenerateKustomization(t *testing.T) {
	// Create a minimal resolved kustomization
	step := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "test-step"},
		Spec:       manifest.StepSpec{Run: "echo 'hello'"},
	}

	cmd := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "test-cmd", CommandName: "testcmd"},
		Spec:       manifest.CmdSpec{Run: "echo 'cmd'"},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "test-workflow", Shell: "bash"},
		Spec: manifest.CmdWorkflowSpec{
			Cmds: []string{"test-cmd"},
			Before: []manifest.StepReference{
				{Step: "test-step"},
			},
		},
	}

	// Create resolved types
	resolvedWorkflow := &resolve.ResolvedCmdWorkflow{
		Workflow: workflow,
		Cmds: map[string]*resolve.ResolvedCmdEntry{
			"test-cmd": {
				Cmd: cmd,
			},
		},
		BeforeSteps: []resolve.ResolvedStep{
			{Step: step, FailurePolicy: "Fail"},
		},
		Shell: "bash",
	}

	resolved := &resolve.ResolvedKustomization{
		Name:     "test",
		Shell:    "bash",
		Workflow: resolvedWorkflow,
		Steps:    map[string]*manifest.Step{"test-step": step},
		Cmds:     map[string]*manifest.Cmd{"test-cmd": cmd},
	}

	// Generate
	gen := NewGenerator("test")
	code, err := gen.GenerateKustomization(resolved)

	assert.NoError(t, err)
	assert.NotEmpty(t, code)
	assert.Contains(t, code, "#!/bin/bash")
	assert.Contains(t, code, "KFG_KUSTOMIZATION_NAME=test")
	assert.Contains(t, code, "testcmd()")
	assert.Contains(t, code, "__kfg_run_step_test-step")
}

func TestGeneratorMetadataEnv(t *testing.T) {
	gen := NewGenerator("test")

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "dev", Shell: "bash"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"cmd1"}},
	}

	cmd := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "cmd1", CommandName: "cmd1"},
		Spec:       manifest.CmdSpec{Run: "echo test"},
	}

	resolved := &resolve.ResolvedKustomization{
		Name:  "my-kustomization",
		Shell: "bash",
		Workflow: &resolve.ResolvedCmdWorkflow{
			Workflow: workflow,
			Cmds:     map[string]*resolve.ResolvedCmdEntry{"cmd1": {Cmd: cmd}},
			Shell:    "bash",
		},
		Steps: map[string]*manifest.Step{},
		Cmds:  map[string]*manifest.Cmd{"cmd1": cmd},
	}

	code, err := gen.GenerateKustomization(resolved)

	assert.NoError(t, err)
	assert.Contains(t, code, "KFG_KUSTOMIZATION_NAME=my-kustomization")
	assert.Contains(t, code, "KFG_WORKFLOW_NAME=dev")
	assert.Contains(t, code, "KFG_SHELL=bash")
}

func TestGeneratorArtifacts(t *testing.T) {
	gen := NewGenerator("test")

	step := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "artifact-step"},
		Spec:       manifest.StepSpec{Run: "echo 'creating'", Artifacts: []string{"output.txt"}},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "test", Shell: "bash"},
		Spec: manifest.CmdWorkflowSpec{
			Cmds:   []string{"cmd1"},
			Before: []manifest.StepReference{{Step: "artifact-step"}},
		},
	}

	cmd := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "cmd1", CommandName: "cmd1"},
		Spec:       manifest.CmdSpec{Run: "echo test"},
	}

	resolved := &resolve.ResolvedKustomization{
		Name:  "test",
		Shell: "bash",
		Workflow: &resolve.ResolvedCmdWorkflow{
			Workflow: workflow,
			Cmds:     map[string]*resolve.ResolvedCmdEntry{"cmd1": {Cmd: cmd}},
			BeforeSteps: []resolve.ResolvedStep{
				{Step: step, FailurePolicy: "Fail"},
			},
			Shell: "bash",
		},
		Steps: map[string]*manifest.Step{"artifact-step": step},
		Cmds:  map[string]*manifest.Cmd{"cmd1": cmd},
	}

	code, err := gen.GenerateKustomization(resolved)

	assert.NoError(t, err)
	assert.Contains(t, code, "KFG_ARTIFACTS")
	assert.Contains(t, code, "__kfg_add_artifact")
	assert.Contains(t, code, "output.txt")
}

func TestGeneratorDeterministic(t *testing.T) {
	gen := NewGenerator("test")

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "test", Shell: "bash"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"cmd1", "cmd2"}},
	}

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

	resolved := &resolve.ResolvedKustomization{
		Name:  "test",
		Shell: "bash",
		Workflow: &resolve.ResolvedCmdWorkflow{
			Workflow: workflow,
			Cmds: map[string]*resolve.ResolvedCmdEntry{
				"cmd1": {Cmd: cmd1},
				"cmd2": {Cmd: cmd2},
			},
			Shell: "bash",
		},
		Steps: map[string]*manifest.Step{},
		Cmds:  map[string]*manifest.Cmd{"cmd1": cmd1, "cmd2": cmd2},
	}

	code1, err := gen.GenerateKustomization(resolved)
	assert.NoError(t, err)

	code2, err := gen.GenerateKustomization(resolved)
	assert.NoError(t, err)

	assert.Equal(t, code1, code2, "generated code should be deterministic")
}
func TestGeneratorSetGetBuildResult(t *testing.T) {
	gen := NewGenerator("test")

	// Initially empty
	assert.Equal(t, "", gen.GetBuildResult())

	// Set build result
	testYAML := "apiVersion: kfg.dev/v1alpha1\nkind: Step\nmetadata:\n  name: test\n"
	gen.SetBuildResult(testYAML)

	// Get should return the same value
	assert.Equal(t, testYAML, gen.GetBuildResult())
}

func TestGeneratorBuildResultSetup(t *testing.T) {
	gen := NewGenerator("test")

	// Set build result YAML
	buildResultYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test-assets
spec:
  schemaRef: schema://providers
  data:
    servers: []
---
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: test-cmd
  commandName: testcmd
spec:
  run: echo test
`
	gen.SetBuildResult(buildResultYAML)

	// Create minimal resolved kustomization
	cmd := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "test-cmd", CommandName: "testcmd"},
		Spec:       manifest.CmdSpec{Run: "echo test"},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "test-workflow", Shell: "bash"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"test-cmd"}},
	}

	resolved := &resolve.ResolvedKustomization{
		Name:  "test",
		Shell: "bash",
		Workflow: &resolve.ResolvedCmdWorkflow{
			Workflow: workflow,
			Cmds:     map[string]*resolve.ResolvedCmdEntry{"test-cmd": {Cmd: cmd}},
			Shell:    "bash",
		},
		Steps: map[string]*manifest.Step{},
		Cmds:  map[string]*manifest.Cmd{"test-cmd": cmd},
	}

	// Generate
	code, err := gen.GenerateKustomization(resolved)

	assert.NoError(t, err)
	assert.NotEmpty(t, code)

	// Verify build result setup is generated in global scope
	assert.Contains(t, code, "# Build result file setup (global scope)")
	assert.Contains(t, code, "mktemp")
	assert.Contains(t, code, "KFG_BUILD_RESULT_FILE")
	assert.Contains(t, code, "__kfg_build_result()")
	assert.Contains(t, code, "base64 -d")
	assert.Contains(t, code, "trap 'rm -f \"$__kfg_build_result_file\"' EXIT")
}

func TestGeneratorNoBuildResultSetup(t *testing.T) {
	gen := NewGenerator("test")

	// Do NOT set build result YAML (empty)

	// Create minimal resolved kustomization
	cmd := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "test-cmd", CommandName: "testcmd"},
		Spec:       manifest.CmdSpec{Run: "echo test"},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "test-workflow", Shell: "bash"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"test-cmd"}},
	}

	resolved := &resolve.ResolvedKustomization{
		Name:  "test",
		Shell: "bash",
		Workflow: &resolve.ResolvedCmdWorkflow{
			Workflow: workflow,
			Cmds:     map[string]*resolve.ResolvedCmdEntry{"test-cmd": {Cmd: cmd}},
			Shell:    "bash",
		},
		Steps: map[string]*manifest.Step{},
		Cmds:  map[string]*manifest.Cmd{"test-cmd": cmd},
	}

	// Generate
	code, err := gen.GenerateKustomization(resolved)

	assert.NoError(t, err)
	assert.NotEmpty(t, code)

	// Verify build result setup is NOT generated
	assert.NotContains(t, code, "# Build result setup")
	assert.NotContains(t, code, "KFG_BUILD_RESULT_FILE")
	assert.NotContains(t, code, "__kfg_build_result()")
}

func TestGeneratorBuildResultHelper(t *testing.T) {
	gen := NewGenerator("test")

	gen.SetBuildResult("test: data")

	// Create minimal resolved kustomization
	cmd := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "test-cmd", CommandName: "testcmd"},
		Spec:       manifest.CmdSpec{Run: "echo test"},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "test-workflow", Shell: "bash"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"test-cmd"}},
	}

	resolved := &resolve.ResolvedKustomization{
		Name:  "test",
		Shell: "bash",
		Workflow: &resolve.ResolvedCmdWorkflow{
			Workflow: workflow,
			Cmds:     map[string]*resolve.ResolvedCmdEntry{"test-cmd": {Cmd: cmd}},
			Shell:    "bash",
		},
		Steps: map[string]*manifest.Step{},
		Cmds:  map[string]*manifest.Cmd{"test-cmd": cmd},
	}

	code, err := gen.GenerateKustomization(resolved)

	assert.NoError(t, err)

	// Verify helper function structure (at global scope)
	assert.Contains(t, code, "__kfg_build_result() {")
	assert.Contains(t, code, "cat \"$KFG_BUILD_RESULT_FILE\"")
	assert.Contains(t, code, "}")

	// Verify cleanup trap uses EXIT (not RETURN) for global scope
	assert.Contains(t, code, "trap 'rm -f \"$__kfg_build_result_file\"' EXIT")
}

func TestGeneratorBuildResultWithMultipleCmds(t *testing.T) {
	gen := NewGenerator("test")

	gen.SetBuildResult("test: data")

	// Create resolved kustomization with multiple cmds
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
		Metadata:   manifest.Metadata{Name: "test-workflow", Shell: "bash"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"cmd1", "cmd2"}},
	}

	resolved := &resolve.ResolvedKustomization{
		Name:  "test",
		Shell: "bash",
		Workflow: &resolve.ResolvedCmdWorkflow{
			Workflow: workflow,
			Cmds: map[string]*resolve.ResolvedCmdEntry{
				"cmd1": {Cmd: cmd1},
				"cmd2": {Cmd: cmd2},
			},
			Shell: "bash",
		},
		Steps: map[string]*manifest.Step{},
		Cmds:  map[string]*manifest.Cmd{"cmd1": cmd1, "cmd2": cmd2},
	}

	code, err := gen.GenerateKustomization(resolved)

	assert.NoError(t, err)

	// Build result setup should appear exactly once in global scope
	// Count occurrences of global build result setup
	buildResultCount := strings.Count(code, "# Build result file setup (global scope)")
	assert.Equal(t, 1, buildResultCount, "Build result setup should appear exactly once in global scope")

	// Each cmd wrapper should NOT have build result setup
	assert.NotContains(t, code, "cmd1() {\n    local __kfg_build_result_file")
	assert.NotContains(t, code, "cmd2() {\n    local __kfg_build_result_file")
}

func TestGeneratorEnvPlaceholderResolution(t *testing.T) {
	// Set up test env var
	t.Setenv("EXP_API_KEY", "test_secret_key")

	gen := NewGenerator("test")

	// Create a step with env placeholder
	step := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "env-step"},
		Spec: manifest.StepSpec{
			Run: "echo 'hello'",
			Env: map[string]string{
				"API_KEY": "{env:EXP_API_KEY}",
				"MODEL":   "static-model", // No placeholder
			},
		},
	}

	// Create a cmd with env placeholder
	cmd := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "test-cmd", CommandName: "testcmd"},
		Spec: manifest.CmdSpec{
			Run: "echo 'cmd'",
			Env: map[string]string{
				"CMD_KEY": "{env:EXP_API_KEY}",
			},
		},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "test-workflow", Shell: "bash"},
		Spec: manifest.CmdWorkflowSpec{
			Cmds:   []string{"test-cmd"},
			Before: []manifest.StepReference{{Step: "env-step"}},
		},
	}

	// Create resolved types
	resolvedWorkflow := &resolve.ResolvedCmdWorkflow{
		Workflow: workflow,
		Cmds: map[string]*resolve.ResolvedCmdEntry{
			"test-cmd": {
				Cmd: cmd,
			},
		},
		BeforeSteps: []resolve.ResolvedStep{
			{Step: step, FailurePolicy: "Fail", Env: map[string]string{"API_KEY": "{env:EXP_API_KEY}"}},
		},
		Shell: "bash",
	}

	resolved := &resolve.ResolvedKustomization{
		Name:     "test",
		Shell:    "bash",
		Workflow: resolvedWorkflow,
		Steps:    map[string]*manifest.Step{"env-step": step},
		Cmds:     map[string]*manifest.Cmd{"test-cmd": cmd},
	}

	// Generate
	code, err := gen.GenerateKustomization(resolved)

	assert.NoError(t, err)
	assert.NotEmpty(t, code)

	// Verify placeholders are resolved to shell variable syntax
	// The env placeholder {env:EXP_API_KEY} should be resolved to $EXP_API_KEY
	assert.Contains(t, code, "$EXP_API_KEY")
	assert.NotContains(t, code, "{env:EXP_API_KEY}")

	// Verify static env value is preserved
	assert.Contains(t, code, "static-model")
}

func TestGeneratorCmdEnvPlaceholderResolution(t *testing.T) {
	// Set up test env var
	t.Setenv("KFG_MODEL", "gpt-4")

	gen := NewGenerator("test")

	// Create a cmd with env placeholder
	cmd := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "test-cmd", CommandName: "testcmd"},
		Spec: manifest.CmdSpec{
			Run: "echo 'cmd'",
			Env: map[string]string{
				"MODEL":  "{env:KFG_MODEL}",
				"STATIC": "value",
			},
		},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "test-workflow", Shell: "bash"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"test-cmd"}},
	}

	resolved := &resolve.ResolvedKustomization{
		Name:  "test",
		Shell: "bash",
		Workflow: &resolve.ResolvedCmdWorkflow{
			Workflow: workflow,
			Cmds: map[string]*resolve.ResolvedCmdEntry{
				"test-cmd": {Cmd: cmd},
			},
			Shell: "bash",
		},
		Steps: map[string]*manifest.Step{},
		Cmds:  map[string]*manifest.Cmd{"test-cmd": cmd},
	}

	// Generate
	code, err := gen.GenerateKustomization(resolved)

	assert.NoError(t, err)

	// Verify placeholder is resolved
	assert.Contains(t, code, "$KFG_MODEL")
	assert.NotContains(t, code, "{env:KFG_MODEL}")

	// Verify static env is preserved
	assert.Contains(t, code, "value")
}

func TestGeneratorStepEnvPlaceholderResolution(t *testing.T) {
	// Set up test env var
	t.Setenv("STEP_VAR", "step_value")

	gen := NewGenerator("test")

	// Create a step with env placeholder
	step := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "env-step"},
		Spec: manifest.StepSpec{
			Run: "echo 'step'",
			Env: map[string]string{
				"STEP_ENV": "{env:STEP_VAR}",
			},
		},
	}

	cmd := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "test-cmd", CommandName: "testcmd"},
		Spec:       manifest.CmdSpec{Run: "echo 'cmd'"},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "test-workflow", Shell: "bash"},
		Spec: manifest.CmdWorkflowSpec{
			Cmds:   []string{"test-cmd"},
			Before: []manifest.StepReference{{Step: "env-step"}},
		},
	}

	resolved := &resolve.ResolvedKustomization{
		Name:  "test",
		Shell: "bash",
		Workflow: &resolve.ResolvedCmdWorkflow{
			Workflow: workflow,
			Cmds: map[string]*resolve.ResolvedCmdEntry{
				"test-cmd": {Cmd: cmd},
			},
			BeforeSteps: []resolve.ResolvedStep{
				{Step: step, FailurePolicy: "Fail", Env: map[string]string{"STEP_ENV": "{env:STEP_VAR}"}},
			},
			Shell: "bash",
		},
		Steps: map[string]*manifest.Step{"env-step": step},
		Cmds:  map[string]*manifest.Cmd{"test-cmd": cmd},
	}

	// Generate
	code, err := gen.GenerateKustomization(resolved)

	assert.NoError(t, err)

	// Verify placeholder is resolved in step function
	assert.Contains(t, code, "$STEP_VAR")
	assert.NotContains(t, code, "{env:STEP_VAR}")

	// Verify export statement contains resolved placeholder with default expansion
	assert.Contains(t, code, "export STEP_ENV=\"${STEP_ENV:-$STEP_VAR}\"")
}

func TestGeneratorMultipleEnvPlaceholders(t *testing.T) {
	// Set up test env vars
	t.Setenv("HOST", "localhost")
	t.Setenv("PORT", "8080")

	gen := NewGenerator("test")

	// Create a cmd with multiple env placeholders
	cmd := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "test-cmd", CommandName: "testcmd"},
		Spec: manifest.CmdSpec{
			Run: "echo 'cmd'",
			Env: map[string]string{
				"URL":  "https://{env:HOST}:{env:PORT}/api",
				"PORT": "{env:PORT}",
			},
		},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "test-workflow", Shell: "bash"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"test-cmd"}},
	}

	resolved := &resolve.ResolvedKustomization{
		Name:  "test",
		Shell: "bash",
		Workflow: &resolve.ResolvedCmdWorkflow{
			Workflow: workflow,
			Cmds: map[string]*resolve.ResolvedCmdEntry{
				"test-cmd": {Cmd: cmd},
			},
			Shell: "bash",
		},
		Steps: map[string]*manifest.Step{},
		Cmds:  map[string]*manifest.Cmd{"test-cmd": cmd},
	}

	// Generate
	code, err := gen.GenerateKustomization(resolved)

	assert.NoError(t, err)

	// Verify all placeholders are resolved
	assert.Contains(t, code, "$HOST")
	assert.Contains(t, code, "$PORT")
	assert.Contains(t, code, "https://$HOST:$PORT/api")
	assert.NotContains(t, code, "{env:HOST}")
	assert.NotContains(t, code, "{env:PORT}")
}

func TestGeneratorMissingEnvVar(t *testing.T) {
	// Do NOT set MISSING_VAR

	gen := NewGenerator("test")

	// Create a cmd with env placeholder for missing var
	cmd := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "test-cmd", CommandName: "testcmd"},
		Spec: manifest.CmdSpec{
			Run: "echo 'cmd'",
			Env: map[string]string{
				"MISSING": "{env:MISSING_VAR}",
			},
		},
	}

	workflow := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "test-workflow", Shell: "bash"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"test-cmd"}},
	}

	resolved := &resolve.ResolvedKustomization{
		Name:  "test",
		Shell: "bash",
		Workflow: &resolve.ResolvedCmdWorkflow{
			Workflow: workflow,
			Cmds: map[string]*resolve.ResolvedCmdEntry{
				"test-cmd": {Cmd: cmd},
			},
			Shell: "bash",
		},
		Steps: map[string]*manifest.Step{},
		Cmds:  map[string]*manifest.Cmd{"test-cmd": cmd},
	}

	// Generate
	code, err := gen.GenerateKustomization(resolved)

	assert.NoError(t, err)

	// Missing env var should still be resolved to shell syntax
	// Shell envsubst will resolve to empty string at runtime
	assert.Contains(t, code, "$MISSING_VAR")
	assert.NotContains(t, code, "{env:MISSING_VAR}")
}

// ============================================================================
// Multi-Workflow Generation Tests
// ============================================================================

func TestGenerateAllWorkflows_Basic(t *testing.T) {
	gen := NewGenerator("test")

	// Create two workflows with separate cmds
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
		Metadata:   manifest.Metadata{Name: "dev", Shell: "bash"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"build"}},
	}

	workflow2 := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "prod", Shell: "bash"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"deploy"}},
	}

	// Create resolved workflows
	resolvedWorkflow1 := &resolve.ResolvedCmdWorkflow{
		Workflow: workflow1,
		Cmds:     map[string]*resolve.ResolvedCmdEntry{"build": {Cmd: cmd1}},
		Shell:    "bash",
	}

	resolvedWorkflow2 := &resolve.ResolvedCmdWorkflow{
		Workflow: workflow2,
		Cmds:     map[string]*resolve.ResolvedCmdEntry{"deploy": {Cmd: cmd2}},
		Shell:    "bash",
	}

	// Create multi-workflow resolved
	multi := &ResolvedMultiWorkflow{
		Name:      "test-kustomization",
		Shell:     "bash",
		Workflows: []*resolve.ResolvedCmdWorkflow{resolvedWorkflow1, resolvedWorkflow2},
		Steps:     map[string]*manifest.Step{},
		Cmds:      map[string]*manifest.Cmd{"build": cmd1, "deploy": cmd2},
	}

	// Generate
	code, err := gen.GenerateAllWorkflows(multi)

	assert.NoError(t, err)
	assert.NotEmpty(t, code)

	// Verify header contains kustomization name
	assert.Contains(t, code, "#!/bin/bash")
	assert.Contains(t, code, "KFG_KUSTOMIZATION_NAME=test-kustomization")

	// Verify KFG_WORKFLOW_NAME is NOT present (multi-workflow mode)
	assert.NotContains(t, code, "KFG_WORKFLOW_NAME")

	// Verify both cmd wrappers are present
	assert.Contains(t, code, "build()")
	assert.Contains(t, code, "deploy()")

	// Verify shell type
	assert.Contains(t, code, "KFG_SHELL=bash")
}

func TestGenerateAllWorkflows_StepDeduplication(t *testing.T) {
	gen := NewGenerator("test")

	// Create a shared step used by both workflows
	sharedStep := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "setup"},
		Spec:       manifest.StepSpec{Run: "echo setup"},
	}

	// Create separate cmds
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
		Metadata:   manifest.Metadata{Name: "dev", Shell: "bash"},
		Spec: manifest.CmdWorkflowSpec{
			Cmds:   []string{"build"},
			Before: []manifest.StepReference{{Step: "setup"}},
		},
	}

	workflow2 := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "prod", Shell: "bash"},
		Spec: manifest.CmdWorkflowSpec{
			Cmds:   []string{"deploy"},
			Before: []manifest.StepReference{{Step: "setup"}},
		},
	}

	// Create resolved workflows with shared step
	resolvedWorkflow1 := &resolve.ResolvedCmdWorkflow{
		Workflow: workflow1,
		Cmds:     map[string]*resolve.ResolvedCmdEntry{"build": {Cmd: cmd1}},
		BeforeSteps: []resolve.ResolvedStep{{Step: sharedStep, FailurePolicy: "Fail"}},
		Shell:    "bash",
	}

	resolvedWorkflow2 := &resolve.ResolvedCmdWorkflow{
		Workflow: workflow2,
		Cmds:     map[string]*resolve.ResolvedCmdEntry{"deploy": {Cmd: cmd2}},
		BeforeSteps: []resolve.ResolvedStep{{Step: sharedStep, FailurePolicy: "Fail"}},
		Shell:    "bash",
	}

	multi := &ResolvedMultiWorkflow{
		Name:      "test",
		Shell:     "bash",
		Workflows: []*resolve.ResolvedCmdWorkflow{resolvedWorkflow1, resolvedWorkflow2},
		Steps:     map[string]*manifest.Step{"setup": sharedStep},
		Cmds:      map[string]*manifest.Cmd{"build": cmd1, "deploy": cmd2},
	}

	// Generate
	code, err := gen.GenerateAllWorkflows(multi)

	assert.NoError(t, err)

	// Verify step function appears exactly once (deduplicated)
	stepCount := strings.Count(code, "__kfg_run_step_setup()")
	assert.Equal(t, 1, stepCount, "shared step should be deduplicated and appear exactly once")

	// Verify both cmds reference the same step
	assert.Contains(t, code, "build()")
	assert.Contains(t, code, "deploy()")
}

func TestGenerateAllWorkflows_MultipleCmdsPerWorkflow(t *testing.T) {
	gen := NewGenerator("test")

	// Create cmds for first workflow
	cmd1 := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "claude", CommandName: "claude"},
		Spec:       manifest.CmdSpec{Run: "echo claude"},
	}

	cmd2 := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "gemini", CommandName: "gemini"},
		Spec:       manifest.CmdSpec{Run: "echo gemini"},
	}

	// Create cmds for second workflow
	cmd3 := &manifest.Cmd{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Cmd",
		Metadata:   manifest.Metadata{Name: "openspec", CommandName: "openspec"},
		Spec:       manifest.CmdSpec{Run: "echo openspec"},
	}

	workflow1 := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "dev", Shell: "bash"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"claude", "gemini"}},
	}

	workflow2 := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "openspec-workflow", Shell: "bash"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"openspec"}},
	}

	resolvedWorkflow1 := &resolve.ResolvedCmdWorkflow{
		Workflow: workflow1,
		Cmds: map[string]*resolve.ResolvedCmdEntry{
			"claude": {Cmd: cmd1},
			"gemini": {Cmd: cmd2},
		},
		Shell: "bash",
	}

	resolvedWorkflow2 := &resolve.ResolvedCmdWorkflow{
		Workflow: workflow2,
		Cmds:     map[string]*resolve.ResolvedCmdEntry{"openspec": {Cmd: cmd3}},
		Shell:    "bash",
	}

	multi := &ResolvedMultiWorkflow{
		Name:      "test",
		Shell:     "bash",
		Workflows: []*resolve.ResolvedCmdWorkflow{resolvedWorkflow1, resolvedWorkflow2},
		Steps:     map[string]*manifest.Step{},
		Cmds:      map[string]*manifest.Cmd{"claude": cmd1, "gemini": cmd2, "openspec": cmd3},
	}

	// Generate
	code, err := gen.GenerateAllWorkflows(multi)

	assert.NoError(t, err)

	// Verify all cmds are present
	assert.Contains(t, code, "claude()")
	assert.Contains(t, code, "gemini()")
	assert.Contains(t, code, "openspec()")
}

func TestGenerateAllWorkflows_Deterministic(t *testing.T) {
	gen := NewGenerator("test")

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
		Metadata:   manifest.Metadata{Name: "dev", Shell: "bash"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"build"}},
	}

	workflow2 := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "prod", Shell: "bash"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"deploy"}},
	}

	resolvedWorkflow1 := &resolve.ResolvedCmdWorkflow{
		Workflow: workflow1,
		Cmds:     map[string]*resolve.ResolvedCmdEntry{"build": {Cmd: cmd1}},
		Shell:    "bash",
	}

	resolvedWorkflow2 := &resolve.ResolvedCmdWorkflow{
		Workflow: workflow2,
		Cmds:     map[string]*resolve.ResolvedCmdEntry{"deploy": {Cmd: cmd2}},
		Shell:    "bash",
	}

	multi := &ResolvedMultiWorkflow{
		Name:      "test",
		Shell:     "bash",
		Workflows: []*resolve.ResolvedCmdWorkflow{resolvedWorkflow1, resolvedWorkflow2},
		Steps:     map[string]*manifest.Step{},
		Cmds:      map[string]*manifest.Cmd{"build": cmd1, "deploy": cmd2},
	}

	// Generate twice
	code1, err := gen.GenerateAllWorkflows(multi)
	assert.NoError(t, err)

	code2, err := gen.GenerateAllWorkflows(multi)
	assert.NoError(t, err)

	// Should be identical
	assert.Equal(t, code1, code2, "multi-workflow generation should be deterministic")
}

func TestGenerateAllWorkflows_WithBuildResult(t *testing.T) {
	gen := NewGenerator("test")

	// Set build result
	gen.SetBuildResult("test: data")

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
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"build"}},
	}

	resolvedWorkflow := &resolve.ResolvedCmdWorkflow{
		Workflow: workflow,
		Cmds:     map[string]*resolve.ResolvedCmdEntry{"build": {Cmd: cmd}},
		Shell:    "bash",
	}

	multi := &ResolvedMultiWorkflow{
		Name:      "test",
		Shell:     "bash",
		Workflows: []*resolve.ResolvedCmdWorkflow{resolvedWorkflow},
		Steps:     map[string]*manifest.Step{},
		Cmds:      map[string]*manifest.Cmd{"build": cmd},
	}

	code, err := gen.GenerateAllWorkflows(multi)

	assert.NoError(t, err)

	// Verify build result setup
	assert.Contains(t, code, "# Build result file setup (global scope)")
	assert.Contains(t, code, "KFG_BUILD_RESULT_FILE")
}

func TestGenerateAllWorkflows_DifferentShells(t *testing.T) {
	gen := NewGenerator("test")

	// Note: design says to log warning if workflows use different shells
	// For now, we just use the first workflow's shell
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
		Metadata:   manifest.Metadata{Name: "dev", Shell: "bash"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"build"}},
	}

	workflow2 := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "prod", Shell: "zsh"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{"deploy"}},
	}

	resolvedWorkflow1 := &resolve.ResolvedCmdWorkflow{
		Workflow: workflow1,
		Cmds:     map[string]*resolve.ResolvedCmdEntry{"build": {Cmd: cmd1}},
		Shell:    "bash",
	}

	resolvedWorkflow2 := &resolve.ResolvedCmdWorkflow{
		Workflow: workflow2,
		Cmds:     map[string]*resolve.ResolvedCmdEntry{"deploy": {Cmd: cmd2}},
		Shell:    "zsh",
	}

	// Shell from ResolvedMultiWorkflow (should be first workflow's shell per design)
	multi := &ResolvedMultiWorkflow{
		Name:      "test",
		Shell:     "bash", // First workflow's shell
		Workflows: []*resolve.ResolvedCmdWorkflow{resolvedWorkflow1, resolvedWorkflow2},
		Steps:     map[string]*manifest.Step{},
		Cmds:      map[string]*manifest.Cmd{"build": cmd1, "deploy": cmd2},
	}

	code, err := gen.GenerateAllWorkflows(multi)

	assert.NoError(t, err)

	// Should use first workflow's shell
	assert.Contains(t, code, "#!/bin/bash")
	assert.Contains(t, code, "KFG_SHELL=bash")
}
