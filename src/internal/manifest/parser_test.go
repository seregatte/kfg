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
    - name: setup
      step: setup-step
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

// Tests for StepReference name validation

func TestValidateStepReferenceNames_Required(t *testing.T) {
	tests := []struct {
		name           string
		yaml           string
		expectError    bool
		errorContains  string
	}{
		{
			name: "workflow step reference without name should error",
			yaml: `apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test-workflow
  shell: bash
spec:
  cmds:
    - cmd1
  before:
    - step: setup-step
`,
			expectError:   true,
			errorContains: "step reference missing required 'name' field",
		},
		{
			name: "workflow step reference with name should pass",
			yaml: `apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test-workflow
  shell: bash
spec:
  cmds:
    - cmd1
  before:
    - name: setup
      step: setup-step
`,
			expectError: false,
		},
		{
			name: "multiple step references without name should error",
			yaml: `apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test-workflow
  shell: bash
spec:
  cmds:
    - cmd1
  before:
    - step: setup-step
    - step: another-step
`,
			expectError:   true,
			errorContains: "step reference missing required 'name' field",
		},
		{
			name: "step reference in after without name should error",
			yaml: `apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test-workflow
  shell: bash
spec:
  cmds:
    - cmd1
  before:
    - name: setup
      step: setup-step
  after:
    - step: cleanup-step
`,
			expectError:   true,
			errorContains: "step reference missing required 'name' field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			resources, err := parser.ParseData("test.yaml", []byte(tt.yaml))
			assert.NoError(t, err)
			assert.Len(t, resources, 1)

			err = resources[0].CmdWorkflow.Validate()
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateStepReferenceNames_Uniqueness(t *testing.T) {
	tests := []struct {
		name           string
		yaml           string
		expectError    bool
		errorContains  string
	}{
		{
			name: "duplicate step reference names should error",
			yaml: `apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test-workflow
  shell: bash
spec:
  cmds:
    - cmd1
  before:
    - name: setup
      step: setup-step
    - name: setup
      step: another-step
`,
			expectError:   true,
			errorContains: "duplicate step reference name",
		},
		{
			name: "unique step reference names should pass",
			yaml: `apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test-workflow
  shell: bash
spec:
  cmds:
    - cmd1
  before:
    - name: setup
      step: setup-step
    - name: verify
      step: another-step
`,
			expectError: false,
		},
		{
			name: "duplicate names in before and after should error",
			yaml: `apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test-workflow
  shell: bash
spec:
  cmds:
    - cmd1
  before:
    - name: cleanup
      step: cleanup-step
  after:
    - name: cleanup
      step: another-cleanup
`,
			expectError:   true,
			errorContains: "duplicate step reference name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			resources, err := parser.ParseData("test.yaml", []byte(tt.yaml))
			assert.NoError(t, err)
			assert.Len(t, resources, 1)

			err = resources[0].CmdWorkflow.Validate()
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateWhenOutputStepReferences(t *testing.T) {
	tests := []struct {
		name           string
		yaml           string
		expectError    bool
		errorContains  string
	}{
		{
			name: "valid when.output.step reference should pass",
			yaml: `apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test-workflow
  shell: bash
spec:
  cmds:
    - cmd1
  before:
    - name: detect
      step: detect-step
    - name: setup
      step: setup-step
      when:
        output:
          step: detect
          name: RESULT
          equals: "yes"
`,
			expectError: false,
		},
		{
			name: "invalid when.output.step reference should error",
			yaml: `apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test-workflow
  shell: bash
spec:
  cmds:
    - cmd1
  before:
    - name: detect
      step: detect-step
    - name: setup
      step: setup-step
      when:
        output:
          step: nonexistent
          name: RESULT
          equals: "yes"
`,
			expectError:   true,
			errorContains: "references non-existent step reference name",
		},
		{
			name: "when.output.step referencing step metadata name should error",
			yaml: `apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test-workflow
  shell: bash
spec:
  cmds:
    - cmd1
  before:
    - name: detect
      step: detect-step
    - name: setup
      step: setup-step
      when:
        output:
          step: detect-step
          name: RESULT
          equals: "yes"
`,
			expectError:   true,
			errorContains: "references non-existent step reference name",
		},
		{
			name: "when.output.step in after section should work",
			yaml: `apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test-workflow
  shell: bash
spec:
  cmds:
    - cmd1
  before:
    - name: detect
      step: detect-step
  after:
    - name: cleanup
      step: cleanup-step
      when:
        output:
          step: detect
          name: RESULT
          equals: "yes"
`,
			expectError: false,
		},
		{
			name: "when.output.step without operators should error",
			yaml: `apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test-workflow
  shell: bash
spec:
  cmds:
    - cmd1
  before:
    - name: detect
      step: detect-step
    - name: setup
      step: setup-step
      when:
        output:
          step: detect
          name: RESULT
`,
			expectError:   true,
			errorContains: "without any operator",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			resources, err := parser.ParseData("test.yaml", []byte(tt.yaml))
			assert.NoError(t, err)
			assert.Len(t, resources, 1)

			err = resources[0].CmdWorkflow.Validate()
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateEnvKfgOutputReferences(t *testing.T) {
	tests := []struct {
		name           string
		yaml           string
		expectError    bool
		errorContains  string
	}{
		{
			name: "valid $kfg.output reference should pass",
			yaml: `apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test-workflow
  shell: bash
spec:
  cmds:
    - cmd1
  before:
    - name: detect
      step: detect-step
    - name: setup
      step: setup-step
      env:
        RESULT: "$kfg.output(detect)"
`,
			expectError: false,
		},
		{
			name: "invalid $kfg.output reference should error",
			yaml: `apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test-workflow
  shell: bash
spec:
  cmds:
    - cmd1
  before:
    - name: detect
      step: detect-step
    - name: setup
      step: setup-step
      env:
        RESULT: "$kfg.output(nonexistent)"
`,
			expectError:   true,
			errorContains: "references non-existent step reference name",
		},
		{
			name: "$kfg.output with output name should pass",
			yaml: `apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test-workflow
  shell: bash
spec:
  cmds:
    - cmd1
  before:
    - name: detect
      step: detect-step
    - name: setup
      step: setup-step
      env:
        RESULT: "$kfg.output(detect.AGENT)"
`,
			expectError: false,
		},
		{
			name: "env without $kfg.output should pass",
			yaml: `apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test-workflow
  shell: bash
spec:
  cmds:
    - cmd1
  before:
    - name: detect
      step: detect-step
    - name: setup
      step: setup-step
      env:
        RESULT: "static-value"
`,
			expectError: false,
		},
		{
			name: "multiple $kfg.output references in same env value should all be validated",
			yaml: `apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test-workflow
  shell: bash
spec:
  cmds:
    - cmd1
  before:
    - name: detect
      step: detect-step
    - name: analyze
      step: analyze-step
    - name: setup
      step: setup-step
      env:
        RESULT: "$kfg.output(detect) and $kfg.output(analyze)"
`,
			expectError: false,
		},
		{
			name: "one invalid $kfg.output in multiple references should error",
			yaml: `apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test-workflow
  shell: bash
spec:
  cmds:
    - cmd1
  before:
    - name: detect
      step: detect-step
    - name: setup
      step: setup-step
      env:
        RESULT: "$kfg.output(detect) and $kfg.output(bad)"
`,
			expectError:   true,
			errorContains: "references non-existent step reference name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			resources, err := parser.ParseData("test.yaml", []byte(tt.yaml))
			assert.NoError(t, err)
			assert.Len(t, resources, 1)

			err = resources[0].CmdWorkflow.Validate()
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}