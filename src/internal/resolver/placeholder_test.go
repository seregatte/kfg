// Package resolver provides placeholder resolution for environment variables.
// This file implements unit tests for placeholder.go.
package resolver

import (
	"os"
	"testing"
)

func TestResolveString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		envVars  map[string]string
		expected string
	}{
		{
			name:     "simple placeholder",
			input:    "{env:API_KEY}",
			envVars:  map[string]string{"API_KEY": "secret123"},
			expected: "$API_KEY",
		},
		{
			name:     "placeholder with underscores",
			input:    "{env:MY_VAR_123}",
			envVars:  map[string]string{"MY_VAR_123": "value"},
			expected: "$MY_VAR_123",
		},
		{
			name:     "placeholder in middle of string",
			input:    "https://{env:HOST}:8080",
			envVars:  map[string]string{"HOST": "localhost"},
			expected: "https://$HOST:8080",
		},
		{
			name:     "multiple placeholders",
			input:    "{env:USER}/{env:PROJECT}",
			envVars:  map[string]string{"USER": "john", "PROJECT": "myproject"},
			expected: "$USER/$PROJECT",
		},
		{
			name:     "missing env var",
			input:    "{env:MISSING_VAR}",
			envVars:  map[string]string{},
			expected: "$MISSING_VAR",
		},
		{
			name:     "no placeholder",
			input:    "plain string",
			envVars:  map[string]string{},
			expected: "plain string",
		},
		{
			name:     "invalid placeholder name (starts with number)",
			input:    "{env:123_INVALID}",
			envVars:  map[string]string{},
			expected: "{env:123_INVALID}", // Should not transform
		},
		{
			name:     "wrong prefix",
			input:    "{var:API_KEY}",
			envVars:  map[string]string{},
			expected: "{var:API_KEY}", // Should not transform
		},
		{
			name:     "empty env var",
			input:    "{env:EMPTY_VAR}",
			envVars:  map[string]string{"EMPTY_VAR": ""},
			expected: "$EMPTY_VAR",
		},
		{
			name:     "placeholder at start and end",
			input:    "{env:PREFIX}_middle_{env:SUFFIX}",
			envVars:  map[string]string{"PREFIX": "pre", "SUFFIX": "suf"},
			expected: "$PREFIX_middle_$SUFFIX",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up env vars
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			// Clear any env vars not in test
			for _, k := range []string{"API_KEY", "MY_VAR_123", "HOST", "USER", "PROJECT", "MISSING_VAR", "EMPTY_VAR", "PREFIX", "SUFFIX"} {
				if _, ok := tt.envVars[k]; !ok {
					os.Unsetenv(k)
				}
			}

			result := ResolveString(tt.input)
			if result != tt.expected {
				t.Errorf("ResolveString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestResolveEnvPlaceholders(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		envVars  map[string]string
		expected map[string]interface{}
	}{
		{
			name: "simple string value",
			input: map[string]interface{}{
				"apiKey": "{env:API_KEY}",
			},
			envVars: map[string]string{"API_KEY": "secret"},
			expected: map[string]interface{}{
				"apiKey": "$API_KEY",
			},
		},
		{
			name: "nested map",
			input: map[string]interface{}{
				"servers": map[string]interface{}{
					"main": map[string]interface{}{
						"apiKey": "{env:SERVER_API_KEY}",
					},
				},
			},
			envVars: map[string]string{"SERVER_API_KEY": "server_secret"},
			expected: map[string]interface{}{
				"servers": map[string]interface{}{
					"main": map[string]interface{}{
						"apiKey": "$SERVER_API_KEY",
					},
				},
			},
		},
		{
			name: "array with placeholders",
			input: map[string]interface{}{
				"args": []interface{}{"--api-key", "{env:API_KEY}"},
			},
			envVars: map[string]string{"API_KEY": "key123"},
			expected: map[string]interface{}{
				"args": []interface{}{"--api-key", "$API_KEY"},
			},
		},
		{
			name: "mixed content",
			input: map[string]interface{}{
				"url":    "https://{env:HOST}:8080",
				"apiKey": "{env:API_KEY}",
				"name":   "server-name",
			},
			envVars: map[string]string{"HOST": "localhost", "API_KEY": "secret"},
			expected: map[string]interface{}{
				"url":    "https://$HOST:8080",
				"apiKey": "$API_KEY",
				"name":   "server-name",
			},
		},
		{
			name: "deeply nested structure",
			input: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"level3": map[string]interface{}{
							"value": "{env:DEEP_VAR}",
						},
					},
				},
			},
			envVars: map[string]string{"DEEP_VAR": "deep"},
			expected: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"level3": map[string]interface{}{
							"value": "$DEEP_VAR",
						},
					},
				},
			},
		},
		{
			name: "non-string values preserved",
			input: map[string]interface{}{
				"stringVal": "{env:VAR}",
				"intVal":    42,
				"floatVal":  3.14,
				"boolVal":   true,
				"nilVal":    nil,
			},
			envVars: map[string]string{"VAR": "value"},
			expected: map[string]interface{}{
				"stringVal": "$VAR",
				"intVal":    42,
				"floatVal":  3.14,
				"boolVal":   true,
				"nilVal":    nil,
			},
		},
		{
			name:     "nil input",
			input:    nil,
			envVars:  map[string]string{},
			expected: nil,
		},
		{
			name:     "empty map",
			input:    map[string]interface{}{},
			envVars:  map[string]string{},
			expected: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up env vars
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			// Clear any env vars not in test
			for _, k := range []string{"API_KEY", "SERVER_API_KEY", "HOST", "DEEP_VAR", "VAR"} {
				if _, ok := tt.envVars[k]; !ok {
					os.Unsetenv(k)
				}
			}

			result := ResolveEnvPlaceholders(tt.input)
			if !mapsEqual(result, tt.expected) {
				t.Errorf("ResolveEnvPlaceholders() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestResolveMap(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]string
		envVars  map[string]string
		expected map[string]string
	}{
		{
			name: "simple env map",
			input: map[string]string{
				"API_KEY": "{env:EXP_API_KEY}",
				"MODEL":   "{env:KFG_MODEL}",
			},
			envVars: map[string]string{
				"EXP_API_KEY": "key123",
				"KFG_MODEL": "gpt-4",
			},
			expected: map[string]string{
				"API_KEY": "$EXP_API_KEY",
				"MODEL":   "$KFG_MODEL",
			},
		},
		{
			name: "mixed values in env map",
			input: map[string]string{
				"RESOLVED":   "{env:VAR}",
				"NOT_RESOLVED": "static-value",
			},
			envVars: map[string]string{"VAR": "value"},
			expected: map[string]string{
				"RESOLVED":   "$VAR",
				"NOT_RESOLVED": "static-value",
			},
		},
		{
			name:     "nil input",
			input:    nil,
			envVars:  map[string]string{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up env vars
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			// Clear any env vars not in test
			for _, k := range []string{"EXP_API_KEY", "KFG_MODEL", "VAR"} {
				if _, ok := tt.envVars[k]; !ok {
					os.Unsetenv(k)
				}
			}

			result := ResolveMap(tt.input)
			if !stringMapsEqual(result, tt.expected) {
				t.Errorf("ResolveMap() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Helper functions for comparing maps

func mapsEqual(a, b map[string]interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		bv, ok := b[k]
		if !ok {
			return false
		}
		if !valuesEqual(v, bv) {
			return false
		}
	}
	return true
}

func valuesEqual(a, b interface{}) bool {
	switch aVal := a.(type) {
	case string:
		bVal, ok := b.(string)
		return ok && aVal == bVal
	case int:
		bVal, ok := b.(int)
		return ok && aVal == bVal
	case float64:
		bVal, ok := b.(float64)
		return ok && aVal == bVal
	case bool:
		bVal, ok := b.(bool)
		return ok && aVal == bVal
	case nil:
		return b == nil
	case map[string]interface{}:
		bVal, ok := b.(map[string]interface{})
		return ok && mapsEqual(aVal, bVal)
	case []interface{}:
		bVal, ok := b.([]interface{})
		return ok && arraysEqual(aVal, bVal)
	default:
		return false
	}
}

func arraysEqual(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !valuesEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func stringMapsEqual(a, b map[string]string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		bv, ok := b[k]
		if !ok || v != bv {
			return false
		}
	}
	return true
}