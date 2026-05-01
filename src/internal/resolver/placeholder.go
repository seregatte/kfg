// Package resolver provides placeholder resolution for environment variables.
// This file implements {env:VAR} → $VAR transformation for manifest data.
package resolver

import (
	"os"
	"regexp"
	"strings"

	"github.com/seregatte/kfg/src/internal/logger"
)

// envPlaceholderRegex matches {env:VAR_NAME} patterns where VAR_NAME is a valid
// shell variable name: starts with letter or underscore, followed by letters,
// digits, or underscores.
var envPlaceholderRegex = regexp.MustCompile(`\{env:([A-Za-z_][A-Za-z0-9_]*)\}`)

// ResolveEnvPlaceholders resolves all {env:VAR} placeholders in the data.
// It recursively processes nested maps and arrays.
// Placeholders are transformed to shell variable syntax: {env:VAR} → $VAR
// If an env var is not set, it resolves to empty string and logs a warning.
func ResolveEnvPlaceholders(data map[string]interface{}) map[string]interface{} {
	if data == nil {
		return nil
	}

	result := make(map[string]interface{}, len(data))
	for k, v := range data {
		result[k] = resolveValue(v)
	}
	return result
}

// resolveValue resolves placeholders in any value type.
func resolveValue(v interface{}) interface{} {
	switch val := v.(type) {
	case string:
		return resolveString(val)
	case map[string]interface{}:
		return resolveMap(val)
	case []interface{}:
		return resolveArray(val)
	default:
		// Return other types unchanged (int, float, bool, nil, etc.)
		return v
	}
}

// resolveString resolves {env:VAR} placeholders in a string.
// Pattern: {env:[A-Za-z_][A-Za-z0-9_]*}
// Transformation: {env:VAR} → $VAR (or empty string if VAR is not set)
func resolveString(s string) string {
	return envPlaceholderRegex.ReplaceAllStringFunc(s, func(match string) string {
		// Extract the variable name from the placeholder
		// match is "{env:VAR_NAME}", we need "VAR_NAME"
		varName := strings.TrimPrefix(match, "{env:")
		varName = strings.TrimSuffix(varName, "}")

		// Get the environment variable value
		envValue := os.Getenv(varName)

		if envValue == "" {
			// Log warning for missing/empty env var
			logger.Warn("resolver:placeholder", "env var "+varName+" not set")
		}

		// Return shell variable syntax: $VAR
		// This preserves the variable reference for shell envsubst resolution
		return "$" + varName
	})
}

// resolveMap recursively resolves placeholders in a nested map.
func resolveMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		return nil
	}

	result := make(map[string]interface{}, len(m))
	for k, v := range m {
		result[k] = resolveValue(v)
	}
	return result
}

// resolveArray recursively resolves placeholders in an array.
func resolveArray(a []interface{}) []interface{} {
	if a == nil {
		return nil
	}

	result := make([]interface{}, len(a))
	for i, v := range a {
		result[i] = resolveValue(v)
	}
	return result
}

// ResolveString is a convenience function to resolve placeholders in a single string.
// This is useful for resolving env values in Step/Cmd env fields.
func ResolveString(s string) string {
	return resolveString(s)
}

// ResolveMap is a convenience function to resolve placeholders in a map[string]string.
// This is useful for resolving env maps in Step/Cmd specifications.
func ResolveMap(env map[string]string) map[string]string {
	if env == nil {
		return nil
	}

	result := make(map[string]string, len(env))
	for k, v := range env {
		result[k] = resolveString(v)
	}
	return result
}