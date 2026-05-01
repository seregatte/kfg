package manifest

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Parser parses YAML manifests into resources.
type Parser struct{}

// NewParser creates a new Parser.
func NewParser() *Parser {
	return &Parser{}
}

// ParsedResource is a container for parsed resources of any kind.
type ParsedResource struct {
	Step        *Step
	Cmd         *Cmd
	CmdWorkflow *CmdWorkflow
}

// Kind returns the resource kind.
func (r ParsedResource) Kind() string {
	if r.Step != nil {
		return "Step"
	}
	if r.Cmd != nil {
		return "Cmd"
	}
	if r.CmdWorkflow != nil {
		return "CmdWorkflow"
	}
	return ""
}

// Name returns the resource name.
func (r ParsedResource) Name() string {
	if r.Step != nil {
		return r.Step.Metadata.Name
	}
	if r.Cmd != nil {
		return r.Cmd.Metadata.Name
	}
	if r.CmdWorkflow != nil {
		return r.CmdWorkflow.Metadata.Name
	}
	return ""
}

// Identity returns the resource identity.
func (r ParsedResource) Identity() ResourceIdentity {
	if r.Step != nil {
		return r.Step.Identity()
	}
	if r.Cmd != nil {
		return r.Cmd.Identity()
	}
	if r.CmdWorkflow != nil {
		return r.CmdWorkflow.Identity()
	}
	return ResourceIdentity{}
}

// Validate validates the resource.
func (r ParsedResource) Validate() error {
	if r.Step != nil {
		return r.Step.Validate()
	}
	if r.Cmd != nil {
		return r.Cmd.Validate()
	}
	if r.CmdWorkflow != nil {
		return r.CmdWorkflow.Validate()
	}
	return fmt.Errorf("empty ParsedResource")
}

// GetStep returns the Step if this is a Step resource.
func (r ParsedResource) GetStep() *Step {
	return r.Step
}

// GetCmd returns the Cmd if this is a Cmd resource.
func (r ParsedResource) GetCmd() *Cmd {
	return r.Cmd
}

// GetCmdWorkflow returns the CmdWorkflow if this is a CmdWorkflow resource.
func (r ParsedResource) GetCmdWorkflow() *CmdWorkflow {
	return r.CmdWorkflow
}

// ParseFile parses a single YAML file into resources.
// Multi-document files (separated by ---) are supported.
func (p *Parser) ParseFile(path string) ([]ParsedResource, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, ParseError{
			File:    path,
			Message: fmt.Sprintf("failed to read file: %v", err),
		}
	}

	return p.ParseData(path, data)
}

// ParseData parses YAML data into resources.
// The path parameter is used for error reporting.
func (p *Parser) ParseData(path string, data []byte) ([]ParsedResource, error) {
	// Handle multi-document YAML files
	decoder := yaml.NewDecoder(bytes.NewReader(data))

	var resources []ParsedResource
	docIndex := 0

	for {
		var raw yaml.Node
		err := decoder.Decode(&raw)
		if err == io.EOF {
			break
		}
		if err != nil {
			line := getErrorLine(data, err)
			return nil, ParseError{
				File:    path,
				Line:    line,
				Message: fmt.Sprintf("failed to parse YAML: %v", err),
			}
		}

		// Skip empty documents
		if raw.Kind == 0 {
			docIndex++
			continue
		}

		resource, err := p.parseNode(path, &raw, docIndex)
		if err != nil {
			return nil, err
		}

		resources = append(resources, resource)
		docIndex++
	}

	return resources, nil
}

// parseNode parses a YAML node into a ParsedResource based on kind.
func (p *Parser) parseNode(path string, node *yaml.Node, docIndex int) (ParsedResource, error) {
	// First pass: decode just the kind field
	var kindOnly struct {
		Kind string `yaml:"kind"`
	}
	if err := node.Decode(&kindOnly); err != nil {
		line := node.Line
		if line == 0 {
			line = getErrorLineFromNode(node, err)
		}
		return ParsedResource{}, ParseError{
			File:    path,
			Line:    line,
			Message: fmt.Sprintf("failed to decode kind: %v", err),
		}
	}

	// Second pass: decode into the correct type based on kind
	switch kindOnly.Kind {
	case "Step":
		var step Step
		if err := node.Decode(&step); err != nil {
			line := node.Line
			if line == 0 {
				line = getErrorLineFromNode(node, err)
			}
			return ParsedResource{}, ParseError{
				File:    path,
				Line:    line,
				Message: fmt.Sprintf("failed to decode Step: %v", err),
			}
		}
		return ParsedResource{Step: &step}, nil

	case "Cmd":
		var cmd Cmd
		if err := node.Decode(&cmd); err != nil {
			line := node.Line
			if line == 0 {
				line = getErrorLineFromNode(node, err)
			}
			return ParsedResource{}, ParseError{
				File:    path,
				Line:    line,
				Message: fmt.Sprintf("failed to decode Cmd: %v", err),
			}
		}
		return ParsedResource{Cmd: &cmd}, nil

	case "CmdWorkflow":
		var workflow CmdWorkflow
		if err := node.Decode(&workflow); err != nil {
			line := node.Line
			if line == 0 {
				line = getErrorLineFromNode(node, err)
			}
			return ParsedResource{}, ParseError{
				File:    path,
				Line:    line,
				Message: fmt.Sprintf("failed to decode CmdWorkflow: %v", err),
			}
		}
		return ParsedResource{CmdWorkflow: &workflow}, nil

	default:
		line := node.Line
		return ParsedResource{}, ParseError{
			File:    path,
			Line:    line,
			Message: fmt.Sprintf("unsupported kind: %s (supported: %s)", kindOnly.Kind, strings.Join(SupportedKinds, ", ")),
		}
	}
}

// ParseDirectory parses all YAML files in a directory recursively.
// Files are processed in lexicographic order.
func (p *Parser) ParseDirectory(dir string) ([]ParsedResource, error) {
	var allResources []ParsedResource

	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Silently ignore missing directories
		return nil, nil
	}

	// Find all YAML files
	files, err := findYAMLElements(dir)
	if err != nil {
		return nil, ParseError{
			File:    dir,
			Message: fmt.Sprintf("failed to find YAML files: %v", err),
		}
	}

	// Sort files lexicographically
	sort.Strings(files)

	// Parse each file
	for _, file := range files {
		resources, err := p.ParseFile(file)
		if err != nil {
			return nil, err
		}
		allResources = append(allResources, resources...)
	}

	return allResources, nil
}

// ParsePath parses all YAML files from a path specification.
// Paths are separated by ':' (colon) and directories are loaded recursively.
func (p *Parser) ParsePath(pathSpec string) ([][]ParsedResource, error) {
	paths := splitPathSpec(pathSpec)
	var layers [][]ParsedResource

	for _, path := range paths {
		// Expand path if it contains environment variables or home directory
		path = expandPath(path)

		resources, err := p.ParseDirectory(path)
		if err != nil {
			return nil, err
		}

		layers = append(layers, resources)
	}

	return layers, nil
}

// findYAMLElements recursively finds all YAML files in a directory.
func findYAMLElements(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".yaml" || ext == ".yml" {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// expandPath expands environment variables and home directory in a path.
func expandPath(path string) string {
	// Handle ~ expansion
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[1:])
		}
	}

	// Handle ${XDG_CONFIG_HOME:-$HOME/.config} syntax
	if strings.Contains(path, "${") {
		path = expandEnvVars(path)
	}

	// Handle $HOME expansion
	path = os.ExpandEnv(path)

	return path
}

// expandEnvVars expands shell-style environment variable syntax with defaults.
func expandEnvVars(s string) string {
	// Handle ${VAR:-default} syntax
	result := s

	// Pattern: ${VAR:-default}
	start := 0
	for {
		idx := strings.Index(result[start:], "${")
		if idx == -1 {
			break
		}
		idx += start

		end := strings.Index(result[idx:], "}")
		if end == -1 {
			break
		}
		end += idx

		// Extract the variable specification
		varSpec := result[idx+2 : end]

		// Check for :- syntax (default value)
		if strings.Contains(varSpec, ":-") {
			parts := strings.SplitN(varSpec, ":-", 2)
			varName := parts[0]
			defaultVal := parts[1]

			// Get the value from environment or use default
			val := os.Getenv(varName)
			if val == "" {
				val = defaultVal
			}

			// Replace the variable specification with the value
			result = result[:idx] + val + result[end+1:]
		} else {
			// Simple variable reference
			val := os.Getenv(varSpec)
			result = result[:idx] + val + result[end+1:]
		}

		start = idx + len(varSpec)
	}

	return result
}

// splitPathSpec splits a path specification into individual paths.
// Paths are separated by ':' (colon), except on Windows where ';' is used.
func splitPathSpec(pathSpec string) []string {
	// Use ':' as separator (POSIX-style)
	return strings.Split(pathSpec, ":")
}

// getErrorLine attempts to extract the line number from a YAML parse error.
func getErrorLine(data []byte, err error) int {
	// Try to extract line number from yaml.TypeError
	if _, ok := err.(*yaml.TypeError); ok {
		// yaml.TypeError contains line information in its message
		// We'll try to parse it
		return 0 // Fallback to 0 if we can't parse
	}

	return 0
}

// getErrorLineFromNode attempts to get the error line from a YAML node.
func getErrorLineFromNode(node *yaml.Node, err error) int {
	if node != nil && node.Line > 0 {
		return node.Line
	}
	return getErrorLine(nil, err)
}
