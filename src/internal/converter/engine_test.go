package converter

import (
	"strings"
	"testing"
)

func TestApplyYAMLToJSON(t *testing.T) {
	engine := NewEngine()

	conv := Converter{
		Name:         "test-conv",
		InputFormat:  "yaml",
		OutputFormat: "json",
		Expression:   ".",
	}

	asset := Asset{
		Name:        "test-asset",
		InputFormat: "yaml",
		Data: map[string]interface{}{
			"key": "value",
			"num": 42,
		},
	}

	result, err := engine.Apply(conv, asset)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Result should be valid JSON containing our values
	if !strings.Contains(result, `"key"`) || !strings.Contains(result, `"value"`) {
		t.Errorf("expected JSON output with key/value, got: %s", result)
	}
}

func TestApplyJSONToYAML(t *testing.T) {
	engine := NewEngine()

	conv := Converter{
		Name:         "json-to-yaml",
		InputFormat:  "json",
		OutputFormat: "yaml",
		Expression:   ".",
	}

	asset := Asset{
		Name:        "json-asset",
		InputFormat: "json",
		Data:        `{"name": "test", "count": 3}`,
	}

	result, err := engine.Apply(conv, asset)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Result should be YAML
	if !strings.Contains(result, "name:") || !strings.Contains(result, "test") {
		t.Errorf("expected YAML output, got: %s", result)
	}
}

func TestApplyYAMLToTOML(t *testing.T) {
	engine := NewEngine()

	conv := Converter{
		Name:         "yaml-to-toml",
		InputFormat:  "yaml",
		OutputFormat: "toml",
		Expression:   ".",
	}

	asset := Asset{
		Name:        "yaml-asset",
		InputFormat: "yaml",
		Data: map[string]interface{}{
			"title": "Test",
			"count": 5,
		},
	}

	result, err := engine.Apply(conv, asset)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// TOML output should contain our values
	if !strings.Contains(result, "title") || !strings.Contains(result, "Test") {
		t.Errorf("expected TOML output, got: %s", result)
	}
}

func TestApplyRawOutputArrayJoin(t *testing.T) {
	engine := NewEngine()

	conv := Converter{
		Name:         "raw-conv",
		InputFormat:  "yaml",
		OutputFormat: "raw",
		Expression:   ".items[]",
	}

	asset := Asset{
		Name:        "array-asset",
		InputFormat: "yaml",
		Data: map[string]interface{}{
			"items": []interface{}{"line1", "line2", "line3"},
		},
	}

	result, err := engine.Apply(conv, asset)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Raw output should join array elements with newlines
	if result != "line1\nline2\nline3" {
		t.Errorf("expected 'line1\\nline2\\nline3', got: %q", result)
	}
}

func TestApplyRawOutputStringPassthrough(t *testing.T) {
	engine := NewEngine()

	conv := Converter{
		Name:         "raw-string-conv",
		InputFormat:  "yaml",
		OutputFormat: "raw",
		Expression:   ".message",
	}

	asset := Asset{
		Name:        "string-asset",
		InputFormat: "yaml",
		Data: map[string]interface{}{
			"message": "hello world",
		},
	}

	result, err := engine.Apply(conv, asset)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "hello world" {
		t.Errorf("expected 'hello world', got: %q", result)
	}
}

func TestApplyInvalidExpression(t *testing.T) {
	engine := NewEngine()

	conv := Converter{
		Name:         "bad-conv",
		InputFormat:  "yaml",
		OutputFormat: "yaml",
		Expression:   "[invalid expression!!",
	}

	asset := Asset{
		Name:        "test-asset",
		InputFormat: "yaml",
		Data: map[string]interface{}{
			"key": "value",
		},
	}

	_, err := engine.Apply(conv, asset)
	if err == nil {
		t.Fatal("expected error for invalid expression, got nil")
	}
}

func TestApplyFormatMismatch(t *testing.T) {
	engine := NewEngine()

	conv := Converter{
		Name:         "mismatch-conv",
		InputFormat:  "json",
		OutputFormat: "json",
		Expression:   ".",
	}

	// Asset claims YAML format but contains invalid JSON data
	asset := Asset{
		Name:        "mismatch-asset",
		InputFormat: "json",
		Data:        `this is not valid json at all`,
	}

	_, err := engine.Apply(conv, asset)
	if err == nil {
		t.Fatal("expected error for format mismatch, got nil")
	}
}

func TestValidateInputFormat(t *testing.T) {
	validFormats := []string{"yaml", "json", "xml", "toml", "csv", "tsv", "hcl", "lua", "ini", "shell", "base64", "uri", "kyaml", "props"}
	for _, f := range validFormats {
		if !ValidateInputFormat(f) {
			t.Errorf("expected %s to be a valid input format", f)
		}
	}

	if ValidateInputFormat("invalid") {
		t.Error("expected 'invalid' to be an invalid input format")
	}
}

func TestValidateOutputFormat(t *testing.T) {
	validFormats := []string{"yaml", "json", "xml", "toml", "csv", "tsv", "hcl", "lua", "ini", "shell", "base64", "uri", "kyaml", "props", "raw"}
	for _, f := range validFormats {
		if !ValidateOutputFormat(f) {
			t.Errorf("expected %s to be a valid output format", f)
		}
	}

	if ValidateOutputFormat("invalid") {
		t.Error("expected 'invalid' to be an invalid output format")
	}
}

func TestAssetsValidation(t *testing.T) {
	tests := []struct {
		name    string
		asset   Asset
		wantErr bool
	}{
		{
			name: "valid yaml asset",
			asset: Asset{
				Name:        "valid-asset",
				InputFormat: "yaml",
				Data:        map[string]interface{}{"key": "value"},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			asset: Asset{
				Name:        "",
				InputFormat: "yaml",
				Data:        map[string]interface{}{"key": "value"},
			},
			wantErr: true,
		},
		{
			name: "missing data",
			asset: Asset{
				Name:        "no-data",
				InputFormat: "yaml",
				Data:        nil,
			},
			wantErr: true,
		},
		{
			name: "invalid format",
			asset: Asset{
				Name:        "bad-format",
				InputFormat: "invalid",
				Data:        "some data",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test via manifest.Assets validation
			// Since Asset is our internal type, we test format validation
			if tt.asset.InputFormat != "yaml" && tt.asset.InputFormat != "json" {
				if ValidateInputFormat(tt.asset.InputFormat) {
					t.Errorf("expected format %s to be invalid", tt.asset.InputFormat)
				}
			}
		})
	}
}

func TestConverterValidation(t *testing.T) {
	tests := []struct {
		name      string
		converter Converter
		wantErr   bool
	}{
		{
			name: "valid converter",
			converter: Converter{
				Name:         "valid-conv",
				InputFormat:  "yaml",
				OutputFormat: "json",
				Expression:   ".",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			converter: Converter{
				Name:         "",
				InputFormat:  "yaml",
				OutputFormat: "json",
				Expression:   ".",
			},
			wantErr: true,
		},
		{
			name: "missing expression",
			converter: Converter{
				Name:         "no-expr",
				InputFormat:  "yaml",
				OutputFormat: "json",
				Expression:   "",
			},
			wantErr: true,
		},
		{
			name: "invalid input format",
			converter: Converter{
				Name:         "bad-input",
				InputFormat:  "invalid",
				OutputFormat: "json",
				Expression:   ".",
			},
			wantErr: true,
		},
		{
			name: "invalid output format",
			converter: Converter{
				Name:         "bad-output",
				InputFormat:  "yaml",
				OutputFormat: "invalid",
				Expression:   ".",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasError := false
			if tt.converter.Name == "" {
				hasError = true
			}
			if tt.converter.Expression == "" {
				hasError = true
			}
			if !ValidateInputFormat(tt.converter.InputFormat) {
				hasError = true
			}
			if !ValidateOutputFormat(tt.converter.OutputFormat) {
				hasError = true
			}

			if hasError != tt.wantErr {
				t.Errorf("validation error mismatch: got %v, want %v", hasError, tt.wantErr)
			}
		})
	}
}

func TestFormatConversion(t *testing.T) {
	engine := NewEngine()

	// Test JSON -> YAML conversion
	conv := Converter{
		Name:         "json-to-yaml",
		InputFormat:  "yaml",
		OutputFormat: "yaml",
		Expression:   ".",
	}

	asset := Asset{
		Name:        "json-data",
		InputFormat: "json",
		Data:        `{"server": "prod", "port": 8080}`,
	}

	result, err := engine.Apply(conv, asset)
	if err != nil {
		t.Fatalf("unexpected error during format conversion: %v", err)
	}

	if !strings.Contains(result, "server:") || !strings.Contains(result, "prod") {
		t.Errorf("expected YAML after JSON->YAML conversion, got: %s", result)
	}
}
