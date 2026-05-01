package templates

import (
	"embed"
	"fmt"
	"strings"
	"text/template"
)

// TemplateFS embeds the template files in the binary.
//
//go:embed *.tmpl
var TemplateFS embed.FS

// TemplateManager manages loaded templates.
type TemplateManager struct {
	templates map[string]*template.Template
}

// NewTemplateManager creates a new template manager and loads templates from embedded FS.
func NewTemplateManager() (*TemplateManager, error) {
	tm := &TemplateManager{
		templates: make(map[string]*template.Template),
	}

	// Load templates from embedded FS
	templateFiles := []string{
		"bash_header.tmpl",
		"bash_helper.tmpl",
		"bash_step.tmpl",
		"bash_command.tmpl",
	}

	for _, filename := range templateFiles {
		content, err := TemplateFS.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read template file %s: %w", filename, err)
		}

		tmpl, err := template.New(filename).Parse(string(content))
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", filename, err)
		}

		// Store template without .tmpl extension
		name := filename[:len(filename)-5] // Remove ".tmpl"
		tm.templates[name] = tmpl
	}

	return tm, nil
}

// Execute executes a named template with the given data.
func (tm *TemplateManager) Execute(name string, data interface{}) (string, error) {
	tmpl, ok := tm.templates[name]
	if !ok {
		return "", fmt.Errorf("template %s not found", name)
	}

	var buf strings.Builder
	err := tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	return buf.String(), nil
}

// ExecuteHeader executes the bash_header template.
func (tm *TemplateManager) ExecuteHeader(data HeaderData) (string, error) {
	return tm.Execute("bash_header", data)
}

// ExecuteHelper executes the bash_helper template.
func (tm *TemplateManager) ExecuteHelper() (string, error) {
	// bash_helper template doesn't need data, use nil
	return tm.Execute("bash_helper", nil)
}

// ExecuteStep executes the bash_step template.
func (tm *TemplateManager) ExecuteStep(data StepData) (string, error) {
	return tm.Execute("bash_step", data)
}

// ExecuteCommand executes the bash_command template.
func (tm *TemplateManager) ExecuteCommand(data CommandData) (string, error) {
	return tm.Execute("bash_command", data)
}

// GetTemplateNames returns all available template names.
func (tm *TemplateManager) GetTemplateNames() []string {
	names := make([]string, 0, len(tm.templates))
	for name := range tm.templates {
		names = append(names, name)
	}
	return names
}
