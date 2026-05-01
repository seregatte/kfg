package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStep(t *testing.T) {
	yaml := `apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: test-step
spec:
  run: echo 'hello'
`

	parser := NewParser()
	resources, err := parser.ParseData("test.yaml", []byte(yaml))

	assert.NoError(t, err)
	assert.Len(t, resources, 1)
	assert.NotNil(t, resources[0].Step)
	assert.Equal(t, "test-step", resources[0].Step.Metadata.Name)
	assert.Equal(t, "echo 'hello'", resources[0].Step.Spec.Run)
}

func TestParseCmd(t *testing.T) {
	yaml := `apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: test-cmd
  commandName: testcmd
spec:
  run: echo 'hello'
`

	parser := NewParser()
	resources, err := parser.ParseData("test.yaml", []byte(yaml))

	assert.NoError(t, err)
	assert.Len(t, resources, 1)
	assert.NotNil(t, resources[0].Cmd)
	assert.Equal(t, "test-cmd", resources[0].Cmd.Metadata.Name)
	assert.Equal(t, "testcmd", resources[0].Cmd.Metadata.CommandName)
}

func TestParseCmdWorkflow(t *testing.T) {
	yaml := `apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test-workflow
  shell: bash
spec:
  cmds:
    - cmd1
    - cmd2
  before:
    - step: setup-step
`

	parser := NewParser()
	resources, err := parser.ParseData("test.yaml", []byte(yaml))

	assert.NoError(t, err)
	assert.Len(t, resources, 1)
	assert.NotNil(t, resources[0].CmdWorkflow)
	assert.Equal(t, "test-workflow", resources[0].CmdWorkflow.Metadata.Name)
	assert.Equal(t, []string{"cmd1", "cmd2"}, resources[0].CmdWorkflow.Spec.Cmds)
	assert.Len(t, resources[0].CmdWorkflow.Spec.Before, 1)
}

func TestParseMultiDocument(t *testing.T) {
	yaml := `---
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: step1
spec:
  run: echo 'step1'
---
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: step2
spec:
  run: echo 'step2'
`

	parser := NewParser()
	resources, err := parser.ParseData("test.yaml", []byte(yaml))

	assert.NoError(t, err)
	assert.Len(t, resources, 2)
	assert.Equal(t, "step1", resources[0].Step.Metadata.Name)
	assert.Equal(t, "step2", resources[1].Step.Metadata.Name)
}

func TestParseInvalidKind(t *testing.T) {
	yaml := `apiVersion: kfg.dev/v1alpha1
kind: InvalidKind
metadata:
  name: test
`

	parser := NewParser()
	_, err := parser.ParseData("test.yaml", []byte(yaml))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported kind")
}

func TestParseInvalidYAML(t *testing.T) {
	yaml := `invalid yaml content
  broken indentation
    not valid
`

	parser := NewParser()
	_, err := parser.ParseData("test.yaml", []byte(yaml))

	assert.Error(t, err)
}

func TestParsedResourceKind(t *testing.T) {
	step := &Step{Metadata: Metadata{Name: "test"}}
	cmd := &Cmd{Metadata: Metadata{Name: "test"}}
	workflow := &CmdWorkflow{Metadata: Metadata{Name: "test"}}

	stepRes := ParsedResource{Step: step}
	assert.Equal(t, "Step", stepRes.Kind())

	cmdRes := ParsedResource{Cmd: cmd}
	assert.Equal(t, "Cmd", cmdRes.Kind())

	workflowRes := ParsedResource{CmdWorkflow: workflow}
	assert.Equal(t, "CmdWorkflow", workflowRes.Kind())
}

func TestParsedResourceName(t *testing.T) {
	step := &Step{Metadata: Metadata{Name: "my-step"}}
	res := ParsedResource{Step: step}
	assert.Equal(t, "my-step", res.Name())
}