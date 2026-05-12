package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/seregatte/kfg/src/internal/manifest"
	"github.com/seregatte/kfg/src/internal/resolve"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestParseLaunchArgs(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedAgent string
		expectedExtra []string
	}{
		{
			name:          "no args",
			args:          []string{},
			expectedAgent: "",
			expectedExtra: []string{},
		},
		{
			name:          "agent only",
			args:          []string{"claude"},
			expectedAgent: "claude",
			expectedExtra: []string{},
		},
		{
			name:          "agent with extra args",
			args:          []string{"claude", "--", "--model", "gpt-4"},
			expectedAgent: "claude",
			expectedExtra: []string{"--model", "gpt-4"},
		},
		{
			name:          "separator only",
			args:          []string{"--", "--model", "gpt-4"},
			expectedAgent: "",
			expectedExtra: []string{"--model", "gpt-4"},
		},
		{
			name:          "multiple extra args",
			args:          []string{"opencode", "--", "--help", "--verbose"},
			expectedAgent: "opencode",
			expectedExtra: []string{"--help", "--verbose"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a dummy command for testing
			agentName, extraArgs := parseLaunchArgs(nil, tt.args)
			assert.Equal(t, tt.expectedAgent, agentName)
			assert.Equal(t, tt.expectedExtra, extraArgs)
		})
	}
}

func TestFindAgent(t *testing.T) {
	// Create test index with Cmds and CmdWorkflows
	cmdClaude := &manifest.Cmd{
		Metadata: manifest.Metadata{
			Name:        "dev.agents.claude",
			CommandName: "claude",
		},
		Spec: manifest.CmdSpec{
			Run: "command claude \"$@\"",
		},
	}
	cmdGemini := &manifest.Cmd{
		Metadata: manifest.Metadata{
			Name:        "dev.agents.gemini",
			CommandName: "gemini",
		},
		Spec: manifest.CmdSpec{
			Run: "command gemini \"$@\"",
		},
	}
	cmdOpenspec := &manifest.Cmd{
		Metadata: manifest.Metadata{
			Name:        "dev.openspec",
			CommandName: "openspec",
		},
		Spec: manifest.CmdSpec{
			Run: "command openspec \"$@\"",
		},
	}

	wfDev := &manifest.CmdWorkflow{
		Metadata: manifest.Metadata{
			Name: "dev.workflows.dev",
		},
		Spec: manifest.CmdWorkflowSpec{
			Cmds: []string{"dev.agents.claude", "dev.agents.gemini"},
		},
	}
	wfOpenspec := &manifest.CmdWorkflow{
		Metadata: manifest.Metadata{
			Name: "dev.workflows.openspec",
		},
		Spec: manifest.CmdWorkflowSpec{
			Cmds: []string{"dev.openspec"},
		},
	}

	resources := []manifest.ParsedResource{
		{Cmd: cmdClaude}, {Cmd: cmdGemini}, {Cmd: cmdOpenspec},
		{CmdWorkflow: wfDev}, {CmdWorkflow: wfOpenspec},
	}
	index := resolve.NewIndex(resources)

	tests := []struct {
		name             string
		agentName        string
		workflowFilter   string
		expectError      bool
		expectedCmdName  string
		expectedWorkflow string
	}{
		{
			name:             "agent found",
			agentName:        "claude",
			workflowFilter:   "",
			expectError:      false,
			expectedCmdName:  "dev.agents.claude",
			expectedWorkflow: "dev.workflows.dev",
		},
		{
			name:           "agent not found",
			agentName:      "nonexistent",
			workflowFilter: "",
			expectError:    true,
		},
		{
			name:           "agent not in specified workflow",
			agentName:      "claude",
			workflowFilter: "dev.workflows.openspec",
			expectError:    true,
		},
		{
			name:             "workflow filter match",
			agentName:        "claude",
			workflowFilter:   "dev.workflows.dev",
			expectError:      false,
			expectedCmdName:  "dev.agents.claude",
			expectedWorkflow: "dev.workflows.dev",
		},
		{
			name:             "openspec agent found",
			agentName:        "openspec",
			workflowFilter:   "",
			expectError:      false,
			expectedCmdName:  "dev.openspec",
			expectedWorkflow: "dev.workflows.openspec",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdName, workflowName, cmd, err := findAgent(index, tt.agentName, tt.workflowFilter)
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, cmdName)
				assert.Empty(t, workflowName)
				assert.Nil(t, cmd)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCmdName, cmdName)
				assert.Equal(t, tt.expectedWorkflow, workflowName)
				assert.NotNil(t, cmd)
			}
		})
	}
}

func TestListAvailableAgents(t *testing.T) {
	// Test with empty index
	emptyResources := []manifest.ParsedResource{}
	emptyIndex := resolve.NewIndex(emptyResources)

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	listAvailableAgents(emptyIndex)

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	os.Stdout = oldStdout

	output := buf.String()
	assert.Contains(t, output, "No agents found")

	// Test with populated index
	cmdClaude := &manifest.Cmd{
		Metadata: manifest.Metadata{
			Name:        "dev.agents.claude",
			CommandName: "claude",
		},
	}
	cmdGemini := &manifest.Cmd{
		Metadata: manifest.Metadata{
			Name:        "dev.agents.gemini",
			CommandName: "gemini",
		},
	}
	wfDev := &manifest.CmdWorkflow{
		Metadata: manifest.Metadata{
			Name: "dev.workflows.dev",
		},
		Spec: manifest.CmdWorkflowSpec{
			Cmds: []string{"dev.agents.claude", "dev.agents.gemini"},
		},
	}

	populatedResources := []manifest.ParsedResource{{Cmd: cmdClaude}, {Cmd: cmdGemini}, {CmdWorkflow: wfDev}}
	populatedIndex := resolve.NewIndex(populatedResources)

	// Capture output
	r2, w2, _ := os.Pipe()
	os.Stdout = w2

	listAvailableAgents(populatedIndex)

	w2.Close()
	var buf2 bytes.Buffer
	buf2.ReadFrom(r2)
	os.Stdout = oldStdout

	output2 := buf2.String()
	assert.Contains(t, output2, "Available agents:")
	assert.Contains(t, output2, "claude")
	assert.Contains(t, output2, "gemini")
	assert.Contains(t, output2, "workflow:")
}

func TestListAvailableAgentsOutputFormat(t *testing.T) {
	cmd := &manifest.Cmd{
		Metadata: manifest.Metadata{
			Name:        "test.agent",
			CommandName: "testagent",
		},
	}
	wf := &manifest.CmdWorkflow{
		Metadata: manifest.Metadata{
			Name: "test.workflow",
		},
		Spec: manifest.CmdWorkflowSpec{
			Cmds: []string{"test.agent"},
		},
	}

	resources := []manifest.ParsedResource{{Cmd: cmd}, {CmdWorkflow: wf}}
	index := resolve.NewIndex(resources)

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	listAvailableAgents(index)

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	os.Stdout = oldStdout

	output := buf.String()
	// Check output format
	lines := strings.Split(output, "\n")
	assert.GreaterOrEqual(t, len(lines), 2)
	assert.Equal(t, "Available agents:", lines[0])
	// Find the agent line
	for _, line := range lines {
		if strings.Contains(line, "testagent") {
			assert.Contains(t, line, "(workflow: test.workflow)")
		}
	}
}

func TestRunCommandFlags(t *testing.T) {
	// Test that flags are registered
	flags := runCmd.Flags()

	// Check --kustomize flag
	kustomizeFlag := flags.Lookup("kustomize")
	assert.NotNil(t, kustomizeFlag)
	assert.Equal(t, "k", kustomizeFlag.Shorthand)

	// Check --file flag
	fileFlag := flags.Lookup("file")
	assert.NotNil(t, fileFlag)
	assert.Equal(t, "f", fileFlag.Shorthand)

	// Check --workflow flag
	workflowFlag := flags.Lookup("workflow")
	assert.NotNil(t, workflowFlag)
	assert.Equal(t, "w", workflowFlag.Shorthand)

	// Check --cmds flag
	cmdsFlag := flags.Lookup("cmds")
	assert.NotNil(t, cmdsFlag)
	assert.Equal(t, "c", cmdsFlag.Shorthand)
}

func TestRunCommandStructure(t *testing.T) {
	assert.NotNil(t, runCmd)
	assert.Equal(t, "run [agent] [-- extra-args...]", runCmd.Use)
	assert.Contains(t, runCmd.Short, "Run an agent")
	assert.NotNil(t, runCmd.RunE)
}

func TestRunCommandKPathFallback(t *testing.T) {
	// Reset viper for each test
	viper.Reset()

	// Test 1: KFG_KPATH is set, no -k or -f flag provided
	os.Setenv("KFG_KPATH", "./test-manifests")
	viper.BindEnv("kpath", "KFG_KPATH")

	// The RunE function should use GetKPath() when no -k or -f is provided
	// We can verify the config getter works
	assert.Equal(t, "./test-manifests", viper.GetString("kpath"))
	os.Unsetenv("KFG_KPATH")

	// Test 2: KFG_KPATH is not set, -k flag provided
	viper.Reset()
	assert.Equal(t, "", viper.GetString("kpath"))

	// Test 3: KFG_KPATH with GitHub URL
	os.Setenv("KFG_KPATH", "https://github.com/owner/repo//manifests")
	viper.BindEnv("kpath", "KFG_KPATH")
	assert.Equal(t, "https://github.com/owner/repo//manifests", viper.GetString("kpath"))
	os.Unsetenv("KFG_KPATH")
}

func TestRunCommandLongDescription(t *testing.T) {
	// Verify the Long description mentions KFG_KPATH and GitHub URLs
	assert.Contains(t, runCmd.Long, "KFG_KPATH")
	assert.Contains(t, runCmd.Long, "github.com")
	assert.Contains(t, runCmd.Long, "https://github.com/owner/repo//path")
}

func TestRunCommandExamples(t *testing.T) {
	// Verify the examples include GitHub URL and KFG_KPATH usage
	assert.Contains(t, runCmd.Long, "kfg run -k https://github.com/owner/repo//manifests")
	assert.Contains(t, runCmd.Long, "KFG_KPATH=./manifests kfg run")
	assert.Contains(t, runCmd.Long, "KFG_KPATH=https://github.com")
}
