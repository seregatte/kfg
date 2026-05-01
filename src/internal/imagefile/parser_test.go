package imagefile

import (
	"strings"
	"testing"
)

func TestParseFromInstruction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		wantRef  string
		wantStage string
	}{
		{
			name:     "simple FROM",
			input:    "FROM claude-base:v2",
			wantErr:  false,
			wantRef:  "claude-base:v2",
			wantStage: "stage0",
		},
		{
			name:     "FROM scratch",
			input:    "FROM scratch",
			wantErr:  false,
			wantRef:  "scratch",
			wantStage: "stage0",
		},
		{
			name:     "FROM with AS clause",
			input:    "FROM claude-base:v2 AS claude",
			wantErr:  false,
			wantRef:  "claude-base:v2",
			wantStage: "claude",
		},
		{
			name:     "FROM case insensitive",
			input:    "from claude-base:v2",
			wantErr:  false,
			wantRef:  "claude-base:v2",
			wantStage: "stage0",
		},
		{
			name:     "FROM with no image",
			input:    "FROM",
			wantErr:  true,
		},
		{
			name:     "FROM with AS but no stage name",
			input:    "FROM claude-base:v2 AS",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(strings.NewReader(tt.input))
			ast, err := parser.Parse()

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(ast.Stages) != 1 {
				t.Errorf("expected 1 stage, got %d", len(ast.Stages))
				return
			}

			stage := ast.Stages[0]
			if stage.From.ImageRef != tt.wantRef {
				t.Errorf("expected image ref %s, got %s", tt.wantRef, stage.From.ImageRef)
			}

			if stage.Name != tt.wantStage {
				t.Errorf("expected stage name %s, got %s", tt.wantStage, stage.Name)
			}
		})
	}
}

func TestParseCopyInstruction(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   bool
		wantFrom  string
		wantSrc   []string
		wantDest  string
	}{
		{
			name:     "COPY from workspace",
			input:    "FROM scratch\nCOPY docs/AGENTS.md AGENTS.md",
			wantErr:  false,
			wantFrom: "",
			wantSrc:  []string{"docs/AGENTS.md"},
			wantDest: "AGENTS.md",
		},
		{
			name:     "COPY from stage",
			input:    "FROM claude-base:v2 AS base\nFROM scratch\nCOPY --from=base .claude/ .claude/",
			wantErr:  false,
			wantFrom: "base",
			wantSrc:  []string{".claude/"},
			wantDest: ".claude/",
		},
		{
			name:     "COPY multiple sources",
			input:    "FROM scratch\nCOPY file1.txt file2.txt ./",
			wantErr:  false,
			wantSrc:  []string{"file1.txt", "file2.txt"},
			wantDest: "./",
		},
		{
			name:     "COPY missing destination",
			input:    "FROM scratch\nCOPY file1.txt",
			wantErr:  true,
		},
		{
			name:     "COPY case insensitive",
			input:    "FROM scratch\ncopy file1.txt ./",
			wantErr:  false,
			wantSrc:  []string{"file1.txt"},
			wantDest: "./",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(strings.NewReader(tt.input))
			ast, err := parser.Parse()

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Get the last stage (where COPY should be)
			stage := ast.Stages[len(ast.Stages)-1]
			if len(stage.Instructions) == 0 {
				t.Error("expected at least one instruction")
				return
			}

			copyInstr, ok := stage.Instructions[0].(*CopyInstruction)
			if !ok {
				t.Error("expected CopyInstruction")
				return
			}

			if copyInstr.FromStage != tt.wantFrom {
				t.Errorf("expected FromStage %s, got %s", tt.wantFrom, copyInstr.FromStage)
			}

			if len(copyInstr.Sources) != len(tt.wantSrc) {
				t.Errorf("expected %d sources, got %d", len(tt.wantSrc), len(copyInstr.Sources))
				return
			}

			for i, src := range tt.wantSrc {
				if copyInstr.Sources[i] != src {
					t.Errorf("expected source %s at index %d, got %s", src, i, copyInstr.Sources[i])
				}
			}

			if copyInstr.Dest != tt.wantDest {
				t.Errorf("expected destination %s, got %s", tt.wantDest, copyInstr.Dest)
			}
		})
	}
}

func TestParseEnvInstruction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		wantVars map[string]string
	}{
		{
			name:     "simple ENV",
			input:    "FROM scratch\nENV DEBUG=true",
			wantErr:  false,
			wantVars: map[string]string{"DEBUG": "true"},
		},
		{
			name:     "ENV with quoted value",
			input:    "FROM scratch\nENV MESSAGE=\"Hello World\"",
			wantErr:  false,
			wantVars: map[string]string{"MESSAGE": "Hello World"},
		},
		{
			name:     "ENV multiple variables",
			input:    "FROM scratch\nENV DEBUG=true TARGET=AGENTS.md",
			wantErr:  false,
			wantVars: map[string]string{"DEBUG": "true", "TARGET": "AGENTS.md"},
		},
		{
			name:     "ENV case insensitive",
			input:    "FROM scratch\nenv DEBUG=true",
			wantErr:  false,
			wantVars: map[string]string{"DEBUG": "true"},
		},
		{
			name:     "ENV missing value",
			input:    "FROM scratch\nENV DEBUG",
			wantErr:  true,
		},
		{
			name:     "ENV unterminated quote",
			input:    "FROM scratch\nENV MESSAGE=\"Hello",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(strings.NewReader(tt.input))
			ast, err := parser.Parse()

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			stage := ast.Stages[len(ast.Stages)-1]
			if len(stage.Instructions) == 0 {
				t.Error("expected at least one instruction")
				return
			}

			envInstr, ok := stage.Instructions[0].(*EnvInstruction)
			if !ok {
				t.Error("expected EnvInstruction")
				return
			}

			if len(envInstr.Vars) != len(tt.wantVars) {
				t.Errorf("expected %d vars, got %d", len(tt.wantVars), len(envInstr.Vars))
				return
			}

			for key, val := range tt.wantVars {
				if envInstr.Vars[key] != val {
					t.Errorf("expected var %s=%s, got %s=%s", key, val, key, envInstr.Vars[key])
				}
			}
		})
	}
}

func TestParseRunInstruction(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   bool
		wantCmd   string
	}{
		{
			name:     "simple RUN",
			input:    "FROM scratch\nRUN echo hello",
			wantErr:  false,
			wantCmd:  "echo hello",
		},
		{
			name:     "RUN with shell",
			input:    "FROM scratch\nRUN sh -c 'cat base.md override.md > AGENTS.md'",
			wantErr:  false,
			wantCmd:  "sh -c 'cat base.md override.md > AGENTS.md'",
		},
		{
			name:     "RUN case insensitive",
			input:    "FROM scratch\nrun echo hello",
			wantErr:  false,
			wantCmd:  "echo hello",
		},
		{
			name:     "RUN empty command",
			input:    "FROM scratch\nRUN",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(strings.NewReader(tt.input))
			ast, err := parser.Parse()

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			stage := ast.Stages[len(ast.Stages)-1]
			if len(stage.Instructions) == 0 {
				t.Error("expected at least one instruction")
				return
			}

			runInstr, ok := stage.Instructions[0].(*RunInstruction)
			if !ok {
				t.Error("expected RunInstruction")
				return
			}

			if runInstr.Command != tt.wantCmd {
				t.Errorf("expected command %s, got %s", tt.wantCmd, runInstr.Command)
			}
		})
	}
}

func TestParseTagInstruction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		wantName string
		wantTag  string
	}{
		{
			name:     "simple TAG",
			input:    "FROM scratch\nTAG my-image:v1.0",
			wantErr:  false,
			wantName: "my-image",
			wantTag:  "v1.0",
		},
		{
			name:     "TAG case insensitive",
			input:    "FROM scratch\ntag my-image:v1.0",
			wantErr:  false,
			wantName: "my-image",
			wantTag:  "v1.0",
		},
		{
			name:     "TAG missing tag",
			input:    "FROM scratch\nTAG my-image",
			wantErr:  true,
		},
		{
			name:     "TAG with no arguments",
			input:    "FROM scratch\nTAG",
			wantErr:  true,
		},
		{
			name:     "TAG in non-final stage",
			input:    "FROM scratch AS base\nTAG base:v1.0\nFROM claude-base:v2",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(strings.NewReader(tt.input))
			ast, err := parser.Parse()

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Get the stage with TAG
			var tagInstr *TagInstruction
			for _, stage := range ast.Stages {
				if stage.Tag != nil {
					tagInstr = stage.Tag
					break
				}
			}

			if tagInstr == nil {
				t.Error("expected TAG instruction")
				return
			}

			if tagInstr.Name != tt.wantName {
				t.Errorf("expected name %s, got %s", tt.wantName, tagInstr.Name)
			}

			if tagInstr.Tag != tt.wantTag {
				t.Errorf("expected tag %s, got %s", tt.wantTag, tagInstr.Tag)
			}
		})
	}
}

func TestLineContinuation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		wantCmd  string
	}{
		{
			name:     "RUN with continuation",
			input:    "FROM scratch\nRUN echo hello \\\nworld",
			wantErr:  false,
			wantCmd:  "echo hello world",
		},
		{
			name:     "RUN with multiple continuations",
			input:    "FROM scratch\nRUN echo \\\nhello \\\nworld",
			wantErr:  false,
			wantCmd:  "echo hello world",
		},
		{
			name:     "unterminated continuation",
			input:    "FROM scratch\nRUN echo \\",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(strings.NewReader(tt.input))
			ast, err := parser.Parse()

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			stage := ast.Stages[len(ast.Stages)-1]
			if len(stage.Instructions) == 0 {
				t.Error("expected at least one instruction")
				return
			}

			runInstr, ok := stage.Instructions[0].(*RunInstruction)
			if !ok {
				t.Error("expected RunInstruction")
				return
			}

			if runInstr.Command != tt.wantCmd {
				t.Errorf("expected command %s, got %s", tt.wantCmd, runInstr.Command)
			}
		})
	}
}

func TestCommentsAndBlankLines(t *testing.T) {
	input := `FROM scratch
# This is a comment

# Another comment
COPY file.txt ./


ENV DEBUG=true`

	parser := NewParser(strings.NewReader(input))
	ast, err := parser.Parse()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if len(ast.Stages) != 1 {
		t.Errorf("expected 1 stage, got %d", len(ast.Stages))
		return
	}

	stage := ast.Stages[0]
	if len(stage.Instructions) != 2 {
		t.Errorf("expected 2 instructions, got %d", len(stage.Instructions))
	}
}

func TestEmptyImagefile(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "empty file",
			input:   "",
			wantErr: true,
		},
		{
			name:    "only comments",
			input:   "# Comment\n# Another comment",
			wantErr: true,
		},
		{
			name:    "only blank lines",
			input:   "\n\n\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(strings.NewReader(tt.input))
			_, err := parser.Parse()

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestInvalidInstruction(t *testing.T) {
	input := "FROM scratch\nINVALID command"

	parser := NewParser(strings.NewReader(input))
	_, err := parser.Parse()

	if err == nil {
		t.Error("expected error for invalid instruction")
		return
	}

	// Check that error mentions the invalid instruction
	if !strings.Contains(err.Error(), "unknown instruction") {
		t.Errorf("expected error to mention 'unknown instruction', got: %v", err)
	}
}

func TestInstructionBeforeFrom(t *testing.T) {
	input := "COPY file.txt ./"

	parser := NewParser(strings.NewReader(input))
	_, err := parser.Parse()

	if err == nil {
		t.Error("expected error for instruction before FROM")
		return
	}

	if !strings.Contains(err.Error(), "FROM") {
		t.Errorf("expected error to mention FROM requirement, got: %v", err)
	}
}

func TestMultiStageBuild(t *testing.T) {
	input := `FROM claude-base:v2 AS base
COPY .claude/ .claude/
FROM scratch
COPY --from=base .claude/ .claude/
TAG my-image:v1.0`

	parser := NewParser(strings.NewReader(input))
	ast, err := parser.Parse()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if len(ast.Stages) != 2 {
		t.Errorf("expected 2 stages, got %d", len(ast.Stages))
		return
	}

	// Check first stage
	stage1 := ast.Stages[0]
	if stage1.Name != "base" {
		t.Errorf("expected stage name 'base', got %s", stage1.Name)
	}
	if stage1.From.ImageRef != "claude-base:v2" {
		t.Errorf("expected image ref 'claude-base:v2', got %s", stage1.From.ImageRef)
	}
	if len(stage1.Instructions) != 1 {
		t.Errorf("expected 1 instruction in stage1, got %d", len(stage1.Instructions))
	}

	// Check second stage
	stage2 := ast.Stages[1]
	if stage2.Name != "stage1" {
		t.Errorf("expected stage name 'stage1', got %s", stage2.Name)
	}
	if stage2.From.ImageRef != "scratch" {
		t.Errorf("expected image ref 'scratch', got %s", stage2.From.ImageRef)
	}
	if len(stage2.Instructions) != 1 {
		t.Errorf("expected 1 instruction in stage2, got %d", len(stage2.Instructions))
	}

	// Check COPY --from
	copyInstr, ok := stage2.Instructions[0].(*CopyInstruction)
	if !ok {
		t.Error("expected CopyInstruction in stage2")
		return
	}
	if copyInstr.FromStage != "base" {
		t.Errorf("expected FromStage 'base', got %s", copyInstr.FromStage)
	}

	// Check TAG in final stage
	if stage2.Tag == nil {
		t.Error("expected TAG instruction in final stage")
		return
	}
	if stage2.Tag.Name != "my-image" {
		t.Errorf("expected tag name 'my-image', got %s", stage2.Tag.Name)
	}
}

func TestCompleteImagefile(t *testing.T) {
	input := `# Multi-stage build example
FROM claude-base:v2 AS base
COPY docs/AGENTS.md AGENTS.md
ENV DEBUG=true

FROM scratch AS final
COPY --from=base AGENTS.md AGENTS.md
COPY .claude/ .claude/
ENV MODE=production
RUN echo "Building image"
TAG my-config:v1.0`

	parser := NewParser(strings.NewReader(input))
	ast, err := parser.Parse()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if len(ast.Stages) != 2 {
		t.Errorf("expected 2 stages, got %d", len(ast.Stages))
	}

	// Verify all instruction types are present
	stage1 := ast.Stages[0]
	hasCopy := false
	hasEnv := false
	for _, instr := range stage1.Instructions {
		switch instr.(type) {
		case *CopyInstruction:
			hasCopy = true
		case *EnvInstruction:
			hasEnv = true
		}
	}
	if !hasCopy {
		t.Error("stage1 missing COPY instruction")
	}
	if !hasEnv {
		t.Error("stage1 missing ENV instruction")
	}

	stage2 := ast.Stages[len(ast.Stages)-1]
	stage2HasCopyFrom := false
	stage2HasEnv := false
	stage2HasRun := false
	stage2HasTag := false
	for _, instr := range stage2.Instructions {
		switch i := instr.(type) {
		case *CopyInstruction:
			if i.FromStage != "" {
				stage2HasCopyFrom = true
			}
		case *EnvInstruction:
			stage2HasEnv = true
		case *RunInstruction:
			stage2HasRun = true
		}
	}
	if stage2.Tag != nil {
		stage2HasTag = true
	}

	if !stage2HasCopyFrom {
		t.Error("final stage missing COPY --from instruction")
	}
	if !stage2HasEnv {
		t.Error("final stage missing ENV instruction")
	}
	if !stage2HasRun {
		t.Error("final stage missing RUN instruction")
	}
	if !stage2HasTag {
		t.Error("final stage missing TAG instruction")
	}
}
// 6.7: WORKDIR parsing tests
func TestParseWorkdirInstruction(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantErr    bool
		wantPath   string
	}{
		{
			name:     "simple workdir absolute",
			input:    "FROM scratch\nWORKDIR /app",
			wantErr:  false,
			wantPath: "/app",
		},
		{
			name:     "workdir relative",
			input:    "FROM scratch\nWORKDIR subdir",
			wantErr:  false,
			wantPath: "subdir",
		},
		{
			name:     "workdir nested absolute",
			input:    "FROM scratch\nWORKDIR /usr/local/bin",
			wantErr:  false,
			wantPath: "/usr/local/bin",
		},
		{
			name:     "workdir with trailing slash",
			input:    "FROM scratch\nWORKDIR /app/",
			wantErr:  false,
			wantPath: "/app/",
		},
		{
			name:     "workdir case insensitive",
			input:    "FROM scratch\nworkdir /app",
			wantErr:  false,
			wantPath: "/app",
		},
		{
			name:     "workdir empty",
			input:    "FROM scratch\nWORKDIR",
			wantErr:  true,
		},
		{
			name:     "workdir dot",
			input:    "FROM scratch\nWORKDIR .",
			wantErr:  false,
			wantPath: ".",
		},
		{
			name:     "workdir double slash",
			input:    "FROM scratch\nWORKDIR //app",
			wantErr:  false,
			wantPath: "//app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(strings.NewReader(tt.input))
			ast, err := parser.Parse()

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			stage := ast.Stages[len(ast.Stages)-1]
			if len(stage.Instructions) == 0 {
				t.Error("expected at least one instruction")
				return
			}

			workdirInstr, ok := stage.Instructions[0].(*WorkdirInstruction)
			if !ok {
				t.Error("expected WorkdirInstruction")
				return
			}

			if workdirInstr.Path != tt.wantPath {
				t.Errorf("expected path %s, got %s", tt.wantPath, workdirInstr.Path)
			}
		})
	}
}

// Test WORKDIR in multi-stage build
func TestWorkdirInMultiStage(t *testing.T) {
	input := `FROM scratch AS base
WORKDIR /base
FROM scratch AS final
WORKDIR /final
TAG test:v1`

	parser := NewParser(strings.NewReader(input))
	ast, err := parser.Parse()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if len(ast.Stages) != 2 {
		t.Errorf("expected 2 stages, got %d", len(ast.Stages))
		return
	}

	// Check base stage WORKDIR
	baseStage := ast.Stages[0]
	if len(baseStage.Instructions) == 0 {
		t.Error("expected WORKDIR instruction in base stage")
		return
	}
	workdirBase, ok := baseStage.Instructions[0].(*WorkdirInstruction)
	if !ok {
		t.Error("expected WorkdirInstruction in base stage")
		return
	}
	if workdirBase.Path != "/base" {
		t.Errorf("expected base workdir /base, got %s", workdirBase.Path)
	}

	// Check final stage WORKDIR
	finalStage := ast.Stages[1]
	if len(finalStage.Instructions) == 0 {
		t.Error("expected WORKDIR instruction in final stage")
		return
	}
	workdirFinal, ok := finalStage.Instructions[0].(*WorkdirInstruction)
	if !ok {
		t.Error("expected WorkdirInstruction in final stage")
		return
	}
	if workdirFinal.Path != "/final" {
		t.Errorf("expected final workdir /final, got %s", workdirFinal.Path)
	}
}

// Test WORKDIR with other instructions
func TestWorkdirWithOtherInstructions(t *testing.T) {
	input := `FROM scratch
WORKDIR /app
COPY file.txt .
ENV DEBUG=true
RUN echo hello
TAG test:v1`

	parser := NewParser(strings.NewReader(input))
	ast, err := parser.Parse()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	stage := ast.Stages[0]
	if len(stage.Instructions) != 4 {
		t.Errorf("expected 4 instructions (WORKDIR, COPY, ENV, RUN), got %d", len(stage.Instructions))
		return
	}

	// Verify instruction types
	_, ok1 := stage.Instructions[0].(*WorkdirInstruction)
	_, ok2 := stage.Instructions[1].(*CopyInstruction)
	_, ok3 := stage.Instructions[2].(*EnvInstruction)
	_, ok4 := stage.Instructions[3].(*RunInstruction)

	if !ok1 {
		t.Error("expected first instruction to be WorkdirInstruction")
	}
	if !ok2 {
		t.Error("expected second instruction to be CopyInstruction")
	}
	if !ok3 {
		t.Error("expected third instruction to be EnvInstruction")
	}
	if !ok4 {
		t.Error("expected fourth instruction to be RunInstruction")
	}
}
