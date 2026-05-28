package cache

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// RestoreResult represents the output of a restore operation.
type RestoreResult struct {
	Lines []string // Shell eval-safe lines to emit
}

// Restore restores cached Step results and emits shell eval-safe output.
// The output lines can be eval'd in bash to restore artifacts and outputs.
func Restore(stepRefName string, workdir string) (*RestoreResult, error) {
	entryPath := GetEntryPath(stepRefName)

	// Check if entry exists
	if !Exists(stepRefName) {
		return nil, fmt.Errorf("cache entry not found: %s", stepRefName)
	}

	// Read metadata
	entry, err := ReadCacheEntry(entryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache entry: %w", err)
	}

	var lines []string

	// Restore artifacts
	if len(entry.Artifacts) > 0 {
		restored, err := RestoreArtifacts(entryPath, workdir, entry.Artifacts)
		if err != nil {
			return nil, fmt.Errorf("failed to restore artifacts: %w", err)
		}

		// Emit __kfg_add_artifact calls for restored artifacts
		for _, artifact := range restored {
			lines = append(lines, fmt.Sprintf(`__kfg_add_artifact %q`, artifact))
		}
	}

	// Restore output if present
	if entry.HasOutput && entry.OutputName != "" {
		// Decode the output value
		outputValue := entry.OutputValue
		lines = append(lines, fmt.Sprintf(`__kfg_output_set %q %q %q`, stepRefName, entry.OutputName, outputValue))
	}

	return &RestoreResult{
		Lines: lines,
	}, nil
}

// RestoreToStdout restores and returns the shell eval-safe output as a string.
// This is the main entry point for the restore subcommand.
func RestoreToStdout(stepRefName string, workdir string) (string, error) {
	result, err := Restore(stepRefName, workdir)
	if err != nil {
		return "", err
	}

	return strings.Join(result.Lines, "\n"), nil
}

// buildShellQuoted builds a shell-safe quoted string.
func buildShellQuoted(s string) string {
	// Use single quotes and escape any single quotes within
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

// encodeOutputValue encodes an output value for storage.
func encodeOutputValue(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

// decodeOutputValue decodes a stored output value.
func decodeOutputValue(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode output value: %w", err)
	}
	return string(decoded), nil
}
