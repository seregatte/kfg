package converter

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"gopkg.in/yaml.v3"
)

// Engine is the yq-go based transformation engine.
type Engine struct {
	evaluator yqlib.StringEvaluator
}

// NewEngine creates a new conversion engine.
func NewEngine() *Engine {
	return &Engine{
		evaluator: yqlib.NewStringEvaluator(),
	}
}

// Apply executes the conversion pipeline: Asset + Converter -> output string.
func (e *Engine) Apply(conv Converter, asset Asset) (string, error) {
	// 1. Serialize asset data to string in asset's format
	inputData, err := e.serializeAsset(asset)
	if err != nil {
		return "", fmt.Errorf("failed to serialize asset data: %w", err)
	}

	// 2. Convert asset data format to converter's expected input format if they differ
	if asset.InputFormat != conv.InputFormat {
		inputData, err = e.convertFormat(inputData, asset.InputFormat, conv.InputFormat)
		if err != nil {
			return "", fmt.Errorf("failed to convert format from %s to %s: %w", asset.InputFormat, conv.InputFormat, err)
		}
	}

	// 3. Get decoder for the converter's input format
	decoder, err := e.getDecoder(conv.InputFormat)
	if err != nil {
		return "", fmt.Errorf("failed to get decoder for format %s: %w", conv.InputFormat, err)
	}

	// 4. Handle raw output format separately
	if conv.OutputFormat == "raw" {
		return e.evaluateRaw(conv.Expression, inputData, decoder)
	}

	// 5. Get encoder for the output format
	encoder, err := e.getEncoder(conv.OutputFormat)
	if err != nil {
		return "", fmt.Errorf("failed to get encoder for format %s: %w", conv.OutputFormat, err)
	}

	// 6. Evaluate expression with yq-go
	result, err := e.evaluator.Evaluate(conv.Expression, inputData, encoder, decoder)
	if err != nil {
		return "", fmt.Errorf("expression evaluation failed: %w", err)
	}

	return result, nil
}

// serializeAsset converts asset data to a string representation.
func (e *Engine) serializeAsset(asset Asset) (string, error) {
	switch data := asset.Data.(type) {
	case string:
		// String data (non-YAML formats) - return as-is
		return data, nil
	case nil:
		return "", fmt.Errorf("asset data is nil")
	default:
		// YAML data (map) - serialize to YAML string
		yamlBytes, err := yaml.Marshal(data)
		if err != nil {
			return "", fmt.Errorf("failed to marshal YAML data: %w", err)
		}
		return string(yamlBytes), nil
	}
}

// convertFormat converts data from one format to another using yq-go.
func (e *Engine) convertFormat(data, fromFormat, toFormat string) (string, error) {
	decoder, err := e.getDecoder(fromFormat)
	if err != nil {
		return "", fmt.Errorf("failed to get decoder for %s: %w", fromFormat, err)
	}
	encoder, err := e.getEncoder(toFormat)
	if err != nil {
		return "", fmt.Errorf("failed to get encoder for %s: %w", toFormat, err)
	}

	// Use "." expression to pass through data with format conversion
	result, err := e.evaluator.Evaluate(".", data, encoder, decoder)
	if err != nil {
		return "", fmt.Errorf("format conversion failed: %w", err)
	}

	return result, nil
}

// evaluateRaw evaluates an expression and returns plain text output.
func (e *Engine) evaluateRaw(expression string, inputData string, decoder yqlib.Decoder) (string, error) {
	// Evaluate with YAML encoder - yq returns plain text for string results
	yamlFormat, err := yqlib.FormatFromString("yaml")
	if err != nil {
		return "", fmt.Errorf("failed to get yaml format: %w", err)
	}

	result, err := e.evaluator.Evaluate(expression, inputData, yamlFormat.GetConfiguredEncoder(), decoder)
	if err != nil {
		return "", fmt.Errorf("expression evaluation failed: %w", err)
	}

	// yq's YAML encoder already returns plain strings for scalar results.
	// For array expressions like .items[], it joins with newlines.
	// Just trim trailing whitespace and return.
	return strings.TrimRight(result, "\n"), nil
}

// extractRawValue extracts a raw string value from a YAML node.
func (e *Engine) extractRawValue(node *yaml.Node) string {
	if node == nil {
		return ""
	}

	// Handle document node
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		return e.extractRawValue(node.Content[0])
	}

	// Handle sequence (array) - join with newlines
	if node.Kind == yaml.SequenceNode {
		var parts []string
		for _, item := range node.Content {
			parts = append(parts, e.extractRawValue(item))
		}
		return strings.Join(parts, "\n")
	}

	// Handle scalar (string/number/etc.)
	if node.Kind == yaml.ScalarNode {
		return node.Value
	}

	// For mapping nodes, marshal back to YAML as fallback
	if node.Kind == yaml.MappingNode {
		bytes, err := yaml.Marshal(node)
		if err == nil {
			return strings.TrimRight(string(bytes), "\n")
		}
	}

	// Fallback: try to marshal
	bytes, err := yaml.Marshal(node)
	if err != nil {
		return ""
	}
	return strings.TrimRight(string(bytes), "\n")
}

// getDecoder returns a yqlib.Decoder for the given format.
func (e *Engine) getDecoder(format string) (yqlib.Decoder, error) {
	yqFormat, err := yqlib.FormatFromString(format)
	if err != nil {
		return nil, err
	}
	// yq formats don't have GetConfiguredDecoder, we need to use the decoder directly
	return yqFormat.DecoderFactory(), nil
}

// getEncoder returns a yqlib.Encoder for the given format.
func (e *Engine) getEncoder(format string) (yqlib.Encoder, error) {
	yqFormat, err := yqlib.FormatFromString(format)
	if err != nil {
		return nil, err
	}
	return yqFormat.GetConfiguredEncoder(), nil
}

// MapManifestAssets converts a manifest.Assets to converter.Asset.
func MapManifestAssets(assets any) Asset {
	// This function works with the manifest.Assets type at runtime
	// We use reflection to extract the fields
	v := reflect.ValueOf(assets)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return Asset{}
	}

	getField := func(name string) reflect.Value {
		f := v.FieldByName(name)
		if f.IsValid() {
			return f
		}
		return reflect.Value{}
	}

	metadata := getField("Metadata")
	spec := getField("Spec")

	name := ""
	inputFormat := "yaml"
	var data any

	if metadata.IsValid() {
		nameField := metadata.FieldByName("Name")
		if nameField.IsValid() {
			name = nameField.String()
		}
	}

	if spec.IsValid() {
		inputField := spec.FieldByName("Input")
		if inputField.IsValid() {
			formatField := inputField.FieldByName("Format")
			if formatField.IsValid() && formatField.String() != "" {
				inputFormat = formatField.String()
			}
		}
		dataField := spec.FieldByName("Data")
		if dataField.IsValid() {
			data = dataField.Interface()
		}
	}

	return Asset{
		Name:        name,
		InputFormat: inputFormat,
		Data:        data,
	}
}

// MapManifestConverter converts a manifest.Converter to converter.Converter.
func MapManifestConverter(converter any) Converter {
	v := reflect.ValueOf(converter)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return Converter{}
	}

	getField := func(name string) reflect.Value {
		f := v.FieldByName(name)
		if f.IsValid() {
			return f
		}
		return reflect.Value{}
	}

	metadata := getField("Metadata")
	spec := getField("Spec")

	name := ""
	inputFormat := "yaml"
	outputFormat := "yaml"
	expression := ""

	if metadata.IsValid() {
		nameField := metadata.FieldByName("Name")
		if nameField.IsValid() {
			name = nameField.String()
		}
	}

	if spec.IsValid() {
		inputField := spec.FieldByName("Input")
		if inputField.IsValid() {
			formatField := inputField.FieldByName("Format")
			if formatField.IsValid() && formatField.String() != "" {
				inputFormat = formatField.String()
			}
		}

		outputField := spec.FieldByName("Output")
		if outputField.IsValid() {
			formatField := outputField.FieldByName("Format")
			if formatField.IsValid() && formatField.String() != "" {
				outputFormat = formatField.String()
			}
		}

		engineField := spec.FieldByName("Engine")
		if engineField.IsValid() {
			exprField := engineField.FieldByName("Expression")
			if exprField.IsValid() {
				expression = exprField.String()
			}
		}
	}

	return Converter{
		Name:         name,
		InputFormat:  inputFormat,
		OutputFormat: outputFormat,
		Expression:   expression,
	}
}

// ValidateOutputFormat checks if the output format is supported.
func ValidateOutputFormat(format string) bool {
	for _, f := range supportedOutputFormats {
		if f == format {
			return true
		}
	}
	return false
}

// ValidateInputFormat checks if the input format is supported.
func ValidateInputFormat(format string) bool {
	for _, f := range supportedInputFormats {
		if f == format {
			return true
		}
	}
	return false
}

var supportedInputFormats = []string{
	"yaml", "json", "xml", "props", "csv", "tsv", "toml",
	"hcl", "lua", "ini", "shell", "base64", "uri", "kyaml",
}

var supportedOutputFormats = []string{
	"yaml", "json", "xml", "props", "csv", "tsv", "toml",
	"hcl", "lua", "ini", "shell", "base64", "uri", "kyaml",
	"raw",
}

// init validates that yqlib supports our expected formats at startup
func init() {
	// Verify yq-go can handle our formats by attempting to create a format
	_, err := yqlib.FormatFromString("yaml")
	if err != nil {
		panic(fmt.Sprintf("yq-go does not support yaml format: %v", err))
	}

	// Also check JSON
	_, err = yqlib.FormatFromString("json")
	if err != nil {
		panic(fmt.Sprintf("yq-go does not support json format: %v", err))
	}

	// Check that our internal JSON marshaling works
	_, err = json.Marshal(map[string]string{"test": "value"})
	if err != nil {
		panic(fmt.Sprintf("json marshaling failed: %v", err))
	}
}
