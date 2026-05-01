package kustomize

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/seregatte/kfg/src/internal/manifest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/kustomize/api/filesys"
)

func TestLoaderLoad(t *testing.T) {
	// Create temporary kustomization directory
	tmpDir, err := os.MkdirTemp("", "kustomize-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create kustomization.yaml
	kustomization := `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - step.yaml
  - cmd.yaml
  - workflow.yaml
`
	err = os.WriteFile(filepath.Join(tmpDir, "kustomization.yaml"), []byte(kustomization), 0644)
	require.NoError(t, err)

	// Create step.yaml
	step := `apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: test-step
spec:
  run: echo test
`
	err = os.WriteFile(filepath.Join(tmpDir, "step.yaml"), []byte(step), 0644)
	require.NoError(t, err)

	// Create cmd.yaml
	cmd := `apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: test-cmd
  commandName: testcmd
spec:
  run: echo cmd
`
	err = os.WriteFile(filepath.Join(tmpDir, "cmd.yaml"), []byte(cmd), 0644)
	require.NoError(t, err)

	// Create workflow.yaml
	workflow := `apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test-workflow
spec:
  cmds:
    - test-cmd
`
	err = os.WriteFile(filepath.Join(tmpDir, "workflow.yaml"), []byte(workflow), 0644)
	require.NoError(t, err)

	// Load kustomization
	loader := NewLoader(nil)
	resMap, err := loader.Load(tmpDir)

	require.NoError(t, err)
	require.NotNil(t, resMap)

	// Convert to resources
	adapter := NewAdapter()
	resources, err := adapter.ResMapToResources(resMap)

	require.NoError(t, err)
	assert.Len(t, resources, 3)
}

func TestAdapterIndexByKind(t *testing.T) {
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

	adapter := NewAdapter()
	steps, cmds, workflows := adapter.IndexByKind(resources)

	assert.Len(t, steps, 1)
	assert.Len(t, cmds, 1)
	assert.Len(t, workflows, 1)

	assert.Equal(t, step, steps["test-step"])
	assert.Equal(t, cmd, cmds["test-cmd"])
	assert.Equal(t, workflow, workflows["test-workflow"])
}

func TestAdapterGetSteps(t *testing.T) {
	step1 := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "step1"},
		Spec:       manifest.StepSpec{Run: "echo 1"},
	}

	step2 := &manifest.Step{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "Step",
		Metadata:   manifest.Metadata{Name: "step2"},
		Spec:       manifest.StepSpec{Run: "echo 2"},
	}

	resources := []manifest.ParsedResource{
		{Step: step1},
		{Step: step2},
	}

	adapter := NewAdapter()
	steps := adapter.GetSteps(resources)

	assert.Len(t, steps, 2)
}

func TestAdapterGetCmds(t *testing.T) {
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

	resources := []manifest.ParsedResource{
		{Cmd: cmd1},
		{Cmd: cmd2},
	}

	adapter := NewAdapter()
	cmds := adapter.GetCmds(resources)

	assert.Len(t, cmds, 2)
}

func TestAdapterGetCmdWorkflows(t *testing.T) {
	workflow1 := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "workflow1"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{}},
	}

	workflow2 := &manifest.CmdWorkflow{
		APIVersion: "kfg.dev/v1alpha1",
		Kind:       "CmdWorkflow",
		Metadata:   manifest.Metadata{Name: "workflow2"},
		Spec:       manifest.CmdWorkflowSpec{Cmds: []string{}},
	}

	resources := []manifest.ParsedResource{
		{CmdWorkflow: workflow1},
		{CmdWorkflow: workflow2},
	}

	adapter := NewAdapter()
	workflows := adapter.GetCmdWorkflows(resources)

	assert.Len(t, workflows, 2)
}

func TestLoaderWithInMemoryFS(t *testing.T) {
	// Create in-memory filesystem
	fSys := filesys.MakeFsInMemory()

	// Create kustomization.yaml
	kustomization := `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - step.yaml
`
	err := fSys.WriteFile("/test/kustomization.yaml", []byte(kustomization))
	require.NoError(t, err)

	// Create step.yaml
	step := `apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: test-step
spec:
  run: echo test
`
	err = fSys.WriteFile("/test/step.yaml", []byte(step))
	require.NoError(t, err)

	// Load kustomization
	loader := NewLoader(nil)
	resMap, err := loader.LoadWithFS(fSys, "/test")

	require.NoError(t, err)
	require.NotNil(t, resMap)
}
