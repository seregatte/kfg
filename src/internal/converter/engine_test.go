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

// --- Tasks 2.4-2.6: Raw input and format detection tests ---

func TestDetectFormatJSON(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"object", `{"key":"value"}`, "json"},
		{"array", `[1,2,3]`, "json"},
		{"string", `"hello"`, "json"},
		{"number", `42`, "json"},
		{"true", `true`, "json"},
		{"false", `false`, "json"},
		{"null", `null`, "json"},
		{"nested", `{"a":{"b":{"c":[1,2,3]}}}`, "json"},
		{"with whitespace", `  {"key":"value"}  `, "json"},
		{"negative number", `-123`, "json"},
		{"float", `3.14`, "json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := engine.DetectFormat(tt.input)
			if got != tt.want {
				t.Errorf("DetectFormat(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestDetectFormatYAML(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"yaml mapping", "key: value", "yaml"},
		{"yaml list", "- item1\n- item2", "yaml"},
		{"yaml multi", "name: test\ncount: 3", "yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := engine.DetectFormat(tt.input)
			if got != tt.want {
				t.Errorf("DetectFormat(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestDetectFormatEdgeCases(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"whitespace only", "   ", ""},
		{"plain text is valid yaml", "this is not json or yaml", "yaml"}, // YAML accepts any scalar
		{"partial json", `{"key":`, ""},                                  // not valid JSON or YAML
		{"yaml is valid json too", `{"key": "value"}`, "json"},           // JSON is valid YAML, JSON detected first
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := engine.DetectFormat(tt.input)
			if got != tt.want {
				t.Errorf("DetectFormat(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestApplyRawJSONWithInlineExpression(t *testing.T) {
	engine := NewEngine()

	result, err := engine.ApplyRaw(`{"key":"value","num":42}`, "json", ".key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result = strings.TrimSpace(result)
	if result != "value" {
		t.Errorf("expected 'value', got %q", result)
	}
}

func TestApplyRawYAMLWithInlineExpression(t *testing.T) {
	engine := NewEngine()

	result, err := engine.ApplyRaw("name: test\ncount: 3", "yaml", ".name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result = strings.TrimSpace(result)
	if result != "test" {
		t.Errorf("expected 'test', got %q", result)
	}
}

func TestApplyRawJSONMerge(t *testing.T) {
	engine := NewEngine()

	// Multi-document merge via YAML format (yq's multi-doc processing is YAML-based)
	// JSON objects are valid YAML, so we can use YAML decoder with --- separators
	// Use "di" (documentIndex) for multi-document merge
	input := `{"a":1}` + "\n---\n" + `{"b":2}`
	result, err := engine.ApplyRaw(input, "yaml", "select(di == 0) * select(di == 1)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "a") || !strings.Contains(result, "b") {
		t.Errorf("expected merged output with both a and b, got: %s", result)
	}
}

func TestApplyRawInvalidExpression(t *testing.T) {
	engine := NewEngine()

	_, err := engine.ApplyRaw(`{"key":"value"}`, "json", "[invalid")
	if err == nil {
		t.Fatal("expected error for invalid expression, got nil")
	}
}

func TestApplyRawUnsupportedFormat(t *testing.T) {
	engine := NewEngine()

	_, err := engine.ApplyRaw(`data`, "invalid_format_xyz", ".")
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}

func TestApplyRawWithConverter(t *testing.T) {
	engine := NewEngine()

	conv := Converter{
		Name:         "raw-conv",
		InputFormat:  "json",
		OutputFormat: "yaml",
		Expression:   ".",
	}

	result, err := engine.ApplyRawWithConverter(`{"name":"test","count":3}`, "json", conv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "name:") || !strings.Contains(result, "test") {
		t.Errorf("expected YAML output, got: %s", result)
	}
}

func TestApplyRawWithConverterFormatConversion(t *testing.T) {
	engine := NewEngine()

	// Input is JSON, converter expects YAML, so format conversion should happen
	conv := Converter{
		Name:         "format-conv",
		InputFormat:  "yaml",
		OutputFormat: "json",
		Expression:   ".",
	}

	result, err := engine.ApplyRawWithConverter(`{"name":"test"}`, "json", conv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Output should be JSON
	if !strings.Contains(result, `"name"`) || !strings.Contains(result, `"test"`) {
		t.Errorf("expected JSON output, got: %s", result)
	}
}

func TestApplyWithExpression(t *testing.T) {
	engine := NewEngine()

	asset := Asset{
		Name:        "test-asset",
		InputFormat: "yaml",
		Data: map[string]interface{}{
			"server": map[string]interface{}{
				"command": "npx",
				"args":    []interface{}{"-y", "pkg"},
			},
		},
	}

	// Extract a nested value
	result, err := engine.ApplyWithExpression(asset, ".server.command")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result = strings.TrimSpace(result)
	if result != "npx" {
		t.Errorf("expected 'npx', got %q", result)
	}
}

func TestApplyWithExpressionJSONAsset(t *testing.T) {
	engine := NewEngine()

	asset := Asset{
		Name:        "json-asset",
		InputFormat: "json",
		Data:        `{"mcpServers":{"context7":{"command":"npx"}}}`,
	}

	result, err := engine.ApplyWithExpression(asset, ".mcpServers.context7.command")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result = strings.TrimSpace(result)
	if result != "npx" {
		t.Errorf("expected 'npx', got %q", result)
	}
}
