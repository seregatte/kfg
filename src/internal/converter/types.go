// Package converter implements the yq-go based transformation engine
// that evaluates expressions against Assets and produces output in the
// specified format.
package converter

// Converter represents a transformation definition.
type Converter struct {
	Name         string
	InputFormat  string
	OutputFormat string
	Expression   string
}

// Asset represents a data payload for conversion.
type Asset struct {
	Name        string
	InputFormat string
	Data        any // map[string]any for YAML, string for other formats
}
