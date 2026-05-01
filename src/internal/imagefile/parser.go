// Package imagefile implements a parser for NixAI Imagefile manifests.
// Imagefile is a declarative format for composing configuration images
// inspired by Dockerfile but tailored to NixAI's needs.
package imagefile

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Instruction represents a parsed Imagefile instruction.
type Instruction interface {
	Line() int
}

// FromInstruction represents a FROM instruction.
type FromInstruction struct {
	lineNum    int
	ImageRef   string // Image reference (e.g., "claude-base:v2" or "scratch")
	StageName  string // Optional stage name from AS clause (e.g., "claude")
}

func (i *FromInstruction) Line() int { return i.lineNum }

// CopyInstruction represents a COPY instruction.
type CopyInstruction struct {
	lineNum    int
	FromStage  string   // Optional --from=<stage> flag
	Sources    []string // Source paths (files or directories)
	Dest       string   // Destination path
}

func (i *CopyInstruction) Line() int { return i.lineNum }

// EnvInstruction represents an ENV instruction.
type EnvInstruction struct {
	lineNum int
	Vars    map[string]string // Environment variables
}

func (i *EnvInstruction) Line() int { return i.lineNum }

// RunInstruction represents a RUN instruction.
type RunInstruction struct {
	lineNum int
	Command string // Shell command to execute
}

func (i *RunInstruction) Line() int { return i.lineNum }

// WorkdirInstruction represents a WORKDIR instruction.
type WorkdirInstruction struct {
	lineNum int
	Path    string // Working directory path
}

func (i *WorkdirInstruction) Line() int { return i.lineNum }

// TagInstruction represents a TAG instruction.
type TagInstruction struct {
	lineNum int
	Name    string // Image name
	Tag     string // Image tag
}

func (i *TagInstruction) Line() int { return i.lineNum }

// Stage represents a build stage with its instructions.
type Stage struct {
	Name        string       // Stage name (from AS clause or generated)
	From        *FromInstruction
	Instructions []Instruction // COPY, ENV, RUN instructions
	Tag         *TagInstruction // Optional TAG instruction (final stage only)
}

// AST represents the parsed Imagefile abstract syntax tree.
type AST struct {
	Stages    []*Stage
	LineCount int
}

// ParseError represents a parsing error with line number context.
type ParseError struct {
	Line    int
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("line %d: %s", e.Line, e.Message)
}

// Parser parses Imagefile manifests into an AST.
type Parser struct {
	scanner *bufio.Scanner
	lineNum int
}

// NewParser creates a new Imagefile parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{
		scanner: bufio.NewScanner(r),
		lineNum: 0,
	}
}

// Validate validates the AST according to Imagefile rules.
func (ast *AST) Validate() error {
	// Check for TAG in non-final stages
	for i, stage := range ast.Stages {
		if stage.Tag != nil && i < len(ast.Stages)-1 {
			return &ParseError{
				Line:    stage.Tag.Line(),
				Message: "TAG instruction can only appear in the final stage",
			}
		}
	}
	return nil
}

// Parse parses the Imagefile and returns the AST.
func (p *Parser) Parse() (*AST, error) {
	ast := &AST{Stages: []*Stage{}}
	var currentStage *Stage

	for p.scanner.Scan() {
		p.lineNum++
		line := p.scanner.Text()

		// Skip blank lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Skip comments
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		// Parse the instruction
		instruction, err := p.parseInstruction(line)
		if err != nil {
			return nil, err
		}

		// Handle stage creation and instruction assignment
		switch instr := instruction.(type) {
		case *FromInstruction:
			// Create a new stage
			stageName := instr.StageName
			if stageName == "" {
				// Generate a stage name if not provided
				stageName = fmt.Sprintf("stage%d", len(ast.Stages))
			}
			currentStage = &Stage{
				Name:         stageName,
				From:         instr,
				Instructions: []Instruction{},
			}
			ast.Stages = append(ast.Stages, currentStage)
		case *TagInstruction:
			if currentStage == nil {
				return nil, &ParseError{
					Line:    p.lineNum,
					Message: "TAG instruction requires a FROM instruction first",
				}
			}
			currentStage.Tag = instr
		default:
			if currentStage == nil {
				return nil, &ParseError{
					Line:    p.lineNum,
					Message: "instruction requires a FROM instruction first",
				}
			}
			currentStage.Instructions = append(currentStage.Instructions, instruction)
		}
	}

	if err := p.scanner.Err(); err != nil {
		return nil, err
	}

	// Validate that we have at least one stage
	if len(ast.Stages) == 0 {
		return nil, &ParseError{
			Line:    0,
			Message: "Imagefile is empty or contains no FROM instructions",
		}
	}

	ast.LineCount = p.lineNum

	// Run validation rules
	if err := ast.Validate(); err != nil {
		return nil, err
	}

	return ast, nil
}

// parseInstruction parses a single line into an instruction.
func (p *Parser) parseInstruction(line string) (Instruction, error) {
	// Handle line continuations
	fullLine, err := p.handleContinuations(line)
	if err != nil {
		return nil, err
	}

	// Split into keyword and arguments
	parts := strings.Fields(fullLine)
	if len(parts) == 0 {
		return nil, &ParseError{
			Line:    p.lineNum,
			Message: "empty instruction",
		}
	}

	keyword := strings.ToUpper(parts[0])
	args := parts[1:]

	switch keyword {
	case "FROM":
		return p.parseFrom(args)
	case "COPY":
		return p.parseCopy(args, fullLine)
	case "ENV":
		return p.parseEnv(args, fullLine)
	case "RUN":
		return p.parseRun(args, fullLine)
	case "WORKDIR":
		return p.parseWorkdir(args)
	case "TAG":
		return p.parseTag(args)
	default:
		return nil, &ParseError{
			Line:    p.lineNum,
			Message: fmt.Sprintf("unknown instruction: %s (valid instructions: FROM, COPY, ENV, RUN, WORKDIR, TAG)", parts[0]),
		}
	}
}

// handleContinuations handles backslash line continuations.
func (p *Parser) handleContinuations(line string) (string, error) {
	result := line
	for strings.HasSuffix(strings.TrimSpace(line), "\\") {
		if !p.scanner.Scan() {
			return "", &ParseError{
				Line:    p.lineNum,
				Message: "unexpected end of file in line continuation",
			}
		}
		p.lineNum++
		line = p.scanner.Text()
		// Remove trailing backslash and append
		result = strings.TrimSuffix(result, "\\") + line
	}
	return result, nil
}

// parseFrom parses a FROM instruction.
func (p *Parser) parseFrom(args []string) (*FromInstruction, error) {
	if len(args) < 1 {
		return nil, &ParseError{
			Line:    p.lineNum,
			Message: "FROM requires an image reference",
		}
	}

	instr := &FromInstruction{lineNum: p.lineNum}
	instr.ImageRef = args[0]

	// Check for AS clause
	for i := 1; i < len(args); i++ {
		if strings.ToUpper(args[i]) == "AS" {
			if i+1 >= len(args) {
				return nil, &ParseError{
					Line:    p.lineNum,
					Message: "AS requires a stage name",
				}
			}
			instr.StageName = args[i+1]
			break
		}
	}

	return instr, nil
}

// parseCopy parses a COPY instruction.
func (p *Parser) parseCopy(args []string, fullLine string) (*CopyInstruction, error) {
	instr := &CopyInstruction{lineNum: p.lineNum}

	// Parse flags
	i := 0
	for i < len(args) {
		if strings.HasPrefix(args[i], "--from=") {
			instr.FromStage = strings.TrimPrefix(args[i], "--from=")
			i++
		} else {
			break
		}
	}

	// Remaining args should be sources and destination
	remaining := args[i:]
	if len(remaining) < 2 {
		return nil, &ParseError{
			Line:    p.lineNum,
			Message: "COPY requires at least one source and a destination",
		}
	}

	instr.Sources = remaining[:len(remaining)-1]
	instr.Dest = remaining[len(remaining)-1]

	return instr, nil
}

// parseEnv parses an ENV instruction.
func (p *Parser) parseEnv(args []string, fullLine string) (*EnvInstruction, error) {
	instr := &EnvInstruction{
		lineNum: p.lineNum,
		Vars:    make(map[string]string),
	}

	// Find the ENV keyword and get everything after it
	envKeyword := "ENV "
	envIdx := strings.Index(strings.ToUpper(fullLine), envKeyword)
	if envIdx == -1 {
		return nil, &ParseError{
			Line:    p.lineNum,
			Message: "ENV instruction malformed",
		}
	}

	envPart := strings.TrimSpace(fullLine[envIdx+len(envKeyword):])

	// Parse key=value pairs
	// Support both: ENV KEY=value and ENV KEY="value with spaces"
	vars, err := p.parseEnvVars(envPart)
	if err != nil {
		return nil, err
	}

	instr.Vars = vars
	return instr, nil
}

// parseEnvVars parses environment variable assignments.
func (p *Parser) parseEnvVars(s string) (map[string]string, error) {
	vars := make(map[string]string)
	i := 0

	for i < len(s) {
		// Skip whitespace
		for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
			i++
		}
		if i >= len(s) {
			break
		}

		// Find key
		keyStart := i
		for i < len(s) && s[i] != '=' && s[i] != ' ' && s[i] != '\t' {
			i++
		}
		key := s[keyStart:i]

		if i >= len(s) || s[i] != '=' {
			return nil, &ParseError{
				Line:    p.lineNum,
				Message: fmt.Sprintf("ENV variable %s missing '='", key),
			}
		}
		i++ // Skip '='

		// Parse value
		var value string
		if i < len(s) && s[i] == '"' {
			// Quoted value
			i++ // Skip opening quote
			valueStart := i
			for i < len(s) && s[i] != '"' {
				i++
			}
			if i >= len(s) {
				return nil, &ParseError{
					Line:    p.lineNum,
					Message: "unterminated quoted value in ENV",
				}
			}
			value = s[valueStart:i]
			i++ // Skip closing quote
		} else {
			// Unquoted value
			valueStart := i
			for i < len(s) && s[i] != ' ' && s[i] != '\t' {
				i++
			}
			value = s[valueStart:i]
		}

		vars[key] = value
	}

	if len(vars) == 0 {
		return nil, &ParseError{
			Line:    p.lineNum,
			Message: "ENV requires at least one variable assignment",
		}
	}

	return vars, nil
}

// parseRun parses a RUN instruction.
func (p *Parser) parseRun(args []string, fullLine string) (*RunInstruction, error) {
	// Find the RUN keyword and get everything after it
	runKeyword := "RUN "
	runIdx := strings.Index(strings.ToUpper(fullLine), runKeyword)
	if runIdx == -1 {
		return nil, &ParseError{
			Line:    p.lineNum,
			Message: "RUN instruction malformed",
		}
	}

	command := strings.TrimSpace(fullLine[runIdx+len(runKeyword):])
	if command == "" {
		return nil, &ParseError{
			Line:    p.lineNum,
			Message: "RUN requires a command",
		}
	}

	return &RunInstruction{
		lineNum: p.lineNum,
		Command: command,
	}, nil
}

// parseTag parses a TAG instruction.
func (p *Parser) parseTag(args []string) (*TagInstruction, error) {
	if len(args) < 1 {
		return nil, &ParseError{
			Line:    p.lineNum,
			Message: "TAG requires a name:tag reference",
		}
	}

	// Parse name:tag format
	ref := args[0]
	parts := strings.SplitN(ref, ":", 2)
	if len(parts) != 2 {
		return nil, &ParseError{
			Line:    p.lineNum,
			Message: "TAG requires name:tag format (e.g., my-image:v1.0)",
		}
	}

	return &TagInstruction{
		lineNum: p.lineNum,
		Name:    parts[0],
		Tag:     parts[1],
	}, nil
}

// parseWorkdir parses a WORKDIR instruction.
func (p *Parser) parseWorkdir(args []string) (*WorkdirInstruction, error) {
	if len(args) < 1 {
		return nil, &ParseError{
			Line:    p.lineNum,
			Message: "WORKDIR requires a path",
		}
	}

	return &WorkdirInstruction{
		lineNum: p.lineNum,
		Path:    args[0],
	}, nil
}