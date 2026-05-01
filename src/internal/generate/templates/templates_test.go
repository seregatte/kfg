package templates

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTemplateManager(t *testing.T) {
	tm, err := NewTemplateManager()
	assert.NoError(t, err)
	assert.NotNil(t, tm)
	assert.NotNil(t, tm.templates)
	assert.Len(t, tm.templates, 4, "should have 4 templates loaded")
}

func TestTemplateManagerGetTemplateNames(t *testing.T) {
	tm, err := NewTemplateManager()
	assert.NoError(t, err)

	names := tm.GetTemplateNames()
	assert.Len(t, names, 4)

	expectedNames := []string{"bash_header", "bash_helper", "bash_step", "bash_command"}
	for _, expected := range expectedNames {
		assert.Contains(t, names, expected)
	}
}

func TestExecuteHeader(t *testing.T) {
	tm, err := NewTemplateManager()
	assert.NoError(t, err)

	data := HeaderData{
		SetName: "test-set",
		Shell:   "bash",
	}

	output, err := tm.ExecuteHeader(data)
	assert.NoError(t, err)
	assert.Contains(t, output, "# kfg shell integration")
	assert.Contains(t, output, "test-set")
	assert.Contains(t, output, "bash")
}

func TestExecuteHelper(t *testing.T) {
	tm, err := NewTemplateManager()
	assert.NoError(t, err)

	output, err := tm.ExecuteHelper()
	assert.NoError(t, err)
	assert.Contains(t, output, "__kfg_ctx_reset")
	assert.Contains(t, output, "__kfg_output_set")
	assert.Contains(t, output, "__kfg_output_get")
	assert.Contains(t, output, "__kfg_when_equals")
	assert.Contains(t, output, "__kfg_when_in")
	assert.Contains(t, output, "__kfg_when_allof")
}

func TestExecuteStep(t *testing.T) {
	tm, err := NewTemplateManager()
	assert.NoError(t, err)

	tests := []struct {
		name   string
		data   StepData
		checks []string
	}{
		{
			name: "simple step",
			data: StepData{
				StepName:  "test-step",
				RunScript: "echo hello",
			},
			checks: []string{
				"__kfg_run_step_test-step",
				"echo hello",
			},
		},
		{
			name: "step with output",
			data: StepData{
				StepName:   "output-step",
				HasOutput:  true,
				RunScript:  "echo test",
				OutputName: "result",
			},
			checks: []string{
				"__kfg_run_step_output-step",
				"__output",
				"__kfg_output_set",
			},
		},
		{
			name: "step with when condition",
			data: StepData{
				StepName:      "conditional-step",
				RunScript:     "echo conditional",
				WhenCondition: "__kfg_when_equals \"step1\" \"result\" \"expected\"",
			},
			checks: []string{
				"__kfg_run_step_conditional-step",
				"if !",
				"return 0  # Skipped",
			},
		},
		{
			name: "multi-line step",
			data: StepData{
				StepName:    "multi-step",
				IsMultiLine: true,
				RunLines:    []string{"echo line1", "echo line2", "echo line3"},
			},
			checks: []string{
				"__kfg_run_step_multi-step",
				"echo line1",
				"echo line2",
				"echo line3",
			},
		},
		{
			name: "step with ignore failure",
			data: StepData{
				StepName:      "ignore-step",
				RunScript:     "echo ignore",
				IgnoreFailure: true,
			},
			checks: []string{
				"__kfg_run_step_ignore-step",
				"|| true",
			},
		},
		{
			name: "step with env",
			data: StepData{
				StepName:  "env-step",
				RunScript: "echo $VAR1 $VAR2",
				Env:       map[string]string{"VAR1": "value1", "VAR2": "value2"},
			},
			checks: []string{
				"__kfg_run_step_env-step",
				"export VAR1=\"value1\"",
				"export VAR2=\"value2\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := tm.ExecuteStep(tt.data)
			assert.NoError(t, err)
			for _, check := range tt.checks {
				assert.True(t, strings.Contains(output, check),
					"output should contain %s", check)
			}
		})
	}
}

func TestExecuteCommand(t *testing.T) {
	tm, err := NewTemplateManager()
	assert.NoError(t, err)

	tests := []struct {
		name   string
		data   CommandData
		checks []string
	}{
		{
			name: "simple command",
			data: CommandData{
				CommandName: "test-cmd",
				MainRun:     "echo main",
			},
			checks: []string{
				"test-cmd()",
				"__kfg_ctx_reset",
				"echo main",
			},
		},
		{
			name: "command with before steps",
			data: CommandData{
				CommandName:    "cmd-with-before",
				MainRun:        "echo main",
				HasBeforeSteps: true,
				BeforeSteps: []BeforeStepData{
					{StepName: "before1", IgnoreFailure: false},
					{StepName: "before2", IgnoreFailure: true},
				},
			},
			checks: []string{
				"cmd-with-before()",
				"__kfg_run_step_before1",
				"__kfg_run_step_before2",
				"|| true",
			},
		},
		{
			name: "command with after steps",
			data: CommandData{
				CommandName:   "cmd-with-after",
				MainRun:       "echo main",
				HasAfterSteps: true,
				AfterSteps: []AfterStepData{
					{StepName: "after1", IgnoreFailure: false},
					{StepName: "after2", IgnoreFailure: true},
				},
			},
			checks: []string{
				"cmd-with-after()",
				"__kfg_run_step_after1",
				"__kfg_run_step_after2",
			},
		},
		{
			name: "command with before and after steps",
			data: CommandData{
				CommandName:    "full-cmd",
				MainRun:        "echo full",
				HasBeforeSteps: true,
				BeforeSteps:    []BeforeStepData{{StepName: "before-step"}},
				HasAfterSteps:  true,
				AfterSteps:     []AfterStepData{{StepName: "after-step"}},
			},
			checks: []string{
				"full-cmd()",
				"__kfg_run_step_before-step",
				"__kfg_run_step_after-step",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := tm.ExecuteCommand(tt.data)
			assert.NoError(t, err)
			for _, check := range tt.checks {
				assert.True(t, strings.Contains(output, check),
					"output should contain %s", check)
			}
		})
	}
}

func TestExecuteInvalidTemplate(t *testing.T) {
	tm, err := NewTemplateManager()
	assert.NoError(t, err)

	// Try to execute a non-existent template
	_, err = tm.Execute("nonexistent", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestTemplateValidation(t *testing.T) {
	// Test that all templates are valid and can be parsed
	tm, err := NewTemplateManager()
	assert.NoError(t, err)
	assert.NotNil(t, tm)

	// All templates should be loaded
	assert.NotNil(t, tm.templates["bash_header"])
	assert.NotNil(t, tm.templates["bash_helper"])
	assert.NotNil(t, tm.templates["bash_step"])
	assert.NotNil(t, tm.templates["bash_command"])
}

func TestHeaderDataStructure(t *testing.T) {
	data := HeaderData{
		SetName: "dev",
		Shell:   "bash",
	}

	assert.Equal(t, "dev", data.SetName)
	assert.Equal(t, "bash", data.Shell)
}

func TestStepDataStructure(t *testing.T) {
	data := StepData{
		StepName:      "test-step",
		WhenCondition: "some condition",
		HasOutput:     true,
		RunScript:     "echo test",
		OutputName:    "result",
		IsMultiLine:   false,
		IgnoreFailure: true,
		Env:           map[string]string{"VAR1": "value1", "VAR2": "value2"},
	}

	assert.Equal(t, "test-step", data.StepName)
	assert.True(t, data.HasOutput)
	assert.True(t, data.IgnoreFailure)
	assert.Equal(t, "value1", data.Env["VAR1"])
	assert.Equal(t, "value2", data.Env["VAR2"])
}

func TestCommandDataStructure(t *testing.T) {
	data := CommandData{
		CommandName:    "my-cmd",
		HasBeforeSteps: true,
		BeforeSteps:    []BeforeStepData{{StepName: "before", IgnoreFailure: false}},
		MainRun:        "echo main",
		HasAfterSteps:  true,
		AfterSteps:     []AfterStepData{{StepName: "after", IgnoreFailure: true}},
	}

	assert.Equal(t, "my-cmd", data.CommandName)
	assert.True(t, data.HasBeforeSteps)
	assert.True(t, data.HasAfterSteps)
	assert.Len(t, data.BeforeSteps, 1)
	assert.Len(t, data.AfterSteps, 1)
}
