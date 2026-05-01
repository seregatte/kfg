package logger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialize(t *testing.T) {
	// Reset global state
	resetLogger()

	tests := []struct {
		name          string
		verbose       string
		logFile       string
		logColor      string
		expectVerbose int
		expectColor   bool
	}{
		{
			name:          "default settings",
			verbose:       "",
			logFile:       "",
			logColor:      "",
			expectVerbose: 1,
			expectColor:   false, // not TTY in test
		},
		{
			name:          "verbose=1",
			verbose:       "1",
			logFile:       "",
			logColor:      "",
			expectVerbose: 1,
			expectColor:   false,
		},
		{
			name:          "verbose=2",
			verbose:       "2",
			logFile:       "",
			logColor:      "",
			expectVerbose: 2,
			expectColor:   false,
		},
		{
			name:          "verbose=3",
			verbose:       "3",
			logFile:       "",
			logColor:      "",
			expectVerbose: 3,
			expectColor:   false,
		},
		{
			name:          "color=always",
			verbose:       "1",
			logFile:       "",
			logColor:      "always",
			expectVerbose: 1,
			expectColor:   true,
		},
		{
			name:          "color=never",
			verbose:       "1",
			logFile:       "",
			logColor:      "never",
			expectVerbose: 1,
			expectColor:   false,
		},
		{
			name:          "invalid verbose defaults to 0",
			verbose:       "invalid",
			logFile:       "",
			logColor:      "",
			expectVerbose: 0,
			expectColor:   false,
		},
		{
			name:          "verbose > 3 defaults to 3",
			verbose:       "10",
			logFile:       "",
			logColor:      "",
			expectVerbose: 3,
			expectColor:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global state
			resetLogger()

			// Set environment variables
			if tt.verbose != "" {
				os.Setenv("KFG_VERBOSE", tt.verbose)
			} else {
				os.Unsetenv("KFG_VERBOSE")
			}
			if tt.logColor != "" {
				os.Setenv("KFG_LOG_COLOR", tt.logColor)
			} else {
				os.Unsetenv("KFG_LOG_COLOR")
			}
			os.Unsetenv("KFG_LOG_FILE")
			os.Unsetenv("KFG_LOG_DIR")

			// Initialize logger
			err := Initialize()
			require.NoError(t, err)

			// Check verbose level
			assert.Equal(t, tt.expectVerbose, GetVerbose())

			// Check color mode
			assert.Equal(t, tt.expectColor, IsColorEnabled())

			// Cleanup
			os.Unsetenv("KFG_VERBOSE")
			os.Unsetenv("KFG_LOG_COLOR")
		})
	}
}

func TestJSONLFileCreation(t *testing.T) {
	resetLogger()

	// Create temp directory for log file
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	os.Setenv("KFG_LOG_FILE", logFile)
	defer os.Unsetenv("KFG_LOG_FILE")

	err := Initialize()
	require.NoError(t, err)

	// Log a message
	Info("test:component", "test message")

	// Sync and close
	Sync()
	Close()

	// Check file exists
	assert.FileExists(t, logFile)

	// Read file content
	content, err := os.ReadFile(logFile)
	require.NoError(t, err)

	// Parse JSONL
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	assert.GreaterOrEqual(t, len(lines), 1)

	// Validate JSON schema
	var event map[string]interface{}
	err = json.Unmarshal([]byte(lines[len(lines)-1]), &event)
	require.NoError(t, err)

	// Check required fields
	assert.Contains(t, event, "ts")
	assert.Contains(t, event, "level")
	assert.Contains(t, event, "component")
	assert.Contains(t, event, "msg")
	assert.Contains(t, event, "source")
	assert.Contains(t, event, "pid")

	// Check field values
	assert.Equal(t, "info", event["level"])
	assert.Equal(t, "core:test:component", event["component"])
	assert.Equal(t, "test message", event["msg"])
	assert.Equal(t, "go", event["source"])
}

func TestVerboseGating(t *testing.T) {
	tests := []struct {
		name          string
		verbose       int
		logLevel      string
		shouldDisplay bool
	}{
		// verbose=0: no human output (JSONL only)
		{"verbose=0, error", 0, "error", false},
		{"verbose=0, info", 0, "info", false},
		{"verbose=0, debug", 0, "debug", false},
		// verbose=1: error only
		{"verbose=1, error", 1, "error", true},
		{"verbose=1, warn", 1, "warn", false},
		{"verbose=1, info", 1, "info", false},
		{"verbose=1, detail", 1, "detail", false},
		{"verbose=1, debug", 1, "debug", false},
		// verbose=2: error + warn + info
		{"verbose=2, error", 2, "error", true},
		{"verbose=2, warn", 2, "warn", true},
		{"verbose=2, info", 2, "info", true},
		{"verbose=2, detail", 2, "detail", false},
		{"verbose=2, debug", 2, "debug", false},
		// verbose=3: all levels
		{"verbose=3, error", 3, "error", true},
		{"verbose=3, warn", 3, "warn", true},
		{"verbose=3, info", 3, "info", true},
		{"verbose=3, detail", 3, "detail", true},
		{"verbose=3, debug", 3, "debug", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check shouldShowLevel function
			level := parseLevelString(tt.logLevel)
			verbose = tt.verbose
			assert.Equal(t, tt.shouldDisplay, shouldShowLevel(level))
		})
	}
}

func TestColorOutputModes(t *testing.T) {
	tests := []struct {
		name        string
		colorMode   ColorMode
		isTTY       bool
		expectColor bool
	}{
		{"auto, not TTY", ColorAuto, false, false},
		{"auto, TTY", ColorAuto, true, true},
		{"always, not TTY", ColorAlways, false, true},
		{"always, TTY", ColorAlways, true, true},
		{"never, not TTY", ColorNever, false, false},
		{"never, TTY", ColorNever, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			colorMode = tt.colorMode
			isTTY = tt.isTTY
			assert.Equal(t, tt.expectColor, shouldUseColor())
		})
	}
}

func TestContextEnrichment(t *testing.T) {
	resetLogger()

	// Set context environment variables
	os.Setenv("KFG_WORKFLOW_NAME", "test-workflow")
	os.Setenv("KFG_KUSTOMIZATION_NAME", "test-kustomization")
	defer func() {
		os.Unsetenv("KFG_WORKFLOW_NAME")
		os.Unsetenv("KFG_KUSTOMIZATION_NAME")
	}()

	// Create temp log file
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")
	os.Setenv("KFG_LOG_FILE", logFile)
	defer os.Unsetenv("KFG_LOG_FILE")

	err := Initialize()
	require.NoError(t, err)

	// Log a message
	Info("test:component", "test message")

	Sync()
	Close()

	// Read and parse JSONL
	content, err := os.ReadFile(logFile)
	require.NoError(t, err)

	var event map[string]interface{}
	err = json.Unmarshal([]byte(strings.TrimSpace(string(content))), &event)
	require.NoError(t, err)

	// Check context fields
	assert.Equal(t, "test-workflow", event["workflow_name"])
	assert.Equal(t, "test-kustomization", event["kustomization_name"])
}

func TestSessionIDEnrichment(t *testing.T) {
	t.Run("session ID is enriched when set", func(t *testing.T) {
		resetLogger()

		// Set session ID
		os.Setenv("KFG_SESSION_ID", "1712938291-4821")
		defer os.Unsetenv("KFG_SESSION_ID")

		// Create temp log file
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "test.log")
		os.Setenv("KFG_LOG_FILE", logFile)
		defer os.Unsetenv("KFG_LOG_FILE")

		err := Initialize()
		require.NoError(t, err)

		// Log a message
		Info("test:component", "test message")

		Sync()
		Close()

		// Read and parse JSONL
		content, err := os.ReadFile(logFile)
		require.NoError(t, err)

		var event map[string]interface{}
		err = json.Unmarshal([]byte(strings.TrimSpace(string(content))), &event)
		require.NoError(t, err)

		// Check session_id field
		assert.Equal(t, "1712938291-4821", event["session_id"])
	})

	t.Run("session ID is absent when not set", func(t *testing.T) {
		resetLogger()

		// Ensure session ID is not set
		os.Unsetenv("KFG_SESSION_ID")

		// Create temp log file
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "test.log")
		os.Setenv("KFG_LOG_FILE", logFile)
		defer os.Unsetenv("KFG_LOG_FILE")

		err := Initialize()
		require.NoError(t, err)

		// Log a message
		Info("test:component", "test message")

		Sync()
		Close()

		// Read and parse JSONL
		content, err := os.ReadFile(logFile)
		require.NoError(t, err)

		var event map[string]interface{}
		err = json.Unmarshal([]byte(strings.TrimSpace(string(content))), &event)
		require.NoError(t, err)

		// Check session_id field is absent
		assert.NotContains(t, event, "session_id", "session_id should not be present when KFG_SESSION_ID is not set")
	})

	t.Run("session ID works with other context fields", func(t *testing.T) {
		resetLogger()

		// Set all context environment variables
		os.Setenv("KFG_WORKFLOW_NAME", "test-workflow")
		os.Setenv("KFG_KUSTOMIZATION_NAME", "test-kustomization")
		os.Setenv("KFG_SESSION_ID", "session-abc")
		defer func() {
			os.Unsetenv("KFG_WORKFLOW_NAME")
			os.Unsetenv("KFG_KUSTOMIZATION_NAME")
			os.Unsetenv("KFG_SESSION_ID")
		}()

		// Create temp log file
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "test.log")
		os.Setenv("KFG_LOG_FILE", logFile)
		defer os.Unsetenv("KFG_LOG_FILE")

		err := Initialize()
		require.NoError(t, err)

		// Log a message
		Info("test:component", "test message")

		Sync()
		Close()

		// Read and parse JSONL
		content, err := os.ReadFile(logFile)
		require.NoError(t, err)

		var event map[string]interface{}
		err = json.Unmarshal([]byte(strings.TrimSpace(string(content))), &event)
		require.NoError(t, err)

		// Check all context fields including session_id
		assert.Equal(t, "test-workflow", event["workflow_name"])
		assert.Equal(t, "test-kustomization", event["kustomization_name"])
		assert.Equal(t, "session-abc", event["session_id"])
	})
}

func TestLogMethods(t *testing.T) {
	resetLogger()

	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")
	os.Setenv("KFG_LOG_FILE", logFile)
	defer os.Unsetenv("KFG_LOG_FILE")

	err := Initialize()
	require.NoError(t, err)

	// Test all log methods from Go code
	Error("error:component", "error message")
	Warn("warn:component", "warn message")
	Info("info:component", "info message")
	Detail("detail:component", "detail message")
	Debug("debug:component", "debug message")

	Sync()
	Close()

	// Read and parse JSONL
	content, err := os.ReadFile(logFile)
	require.NoError(t, err)

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	assert.Equal(t, 5, len(lines))

	// Check each level and that Go code sets source="go"
	for _, line := range lines {
		var event map[string]interface{}
		err := json.Unmarshal([]byte(line), &event)
		require.NoError(t, err)

		// Verify source field is "go" for Go code
		assert.Equal(t, "go", event["source"], "Go code should have source=\"go\"")

		// Verify component field has "core:" prefix for Go logs
		component, ok := event["component"].(string)
		require.True(t, ok, "component should be a string")
		assert.True(t, strings.HasPrefix(component, "core:"), "Go logs should have 'core:' prefix in component field")
		assert.Contains(t, event, "msg")
		assert.Contains(t, event, "ts")
		assert.Contains(t, event, "pid")
	}
}

func TestGetWithoutInitialize(t *testing.T) {
	resetLogger()

	// Get should return a no-op logger when not initialized
	logger := Get()
	assert.NotNil(t, logger)

	// Test that logging methods don't panic with no-op logger
	logger.Info().Str("component", "test").Msg("test message")
}

func TestReinitialize(t *testing.T) {
	resetLogger()

	// Create temp log file
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")
	os.Setenv("KFG_LOG_FILE", logFile)
	defer os.Unsetenv("KFG_LOG_FILE")

	// Initialize with verbose=1
	os.Setenv("KFG_VERBOSE", "1")
	err := Initialize()
	require.NoError(t, err)
	assert.Equal(t, 1, GetVerbose())

	// Change env var to verbose=2
	os.Setenv("KFG_VERBOSE", "2")

	// Reinitialize
	err = Reinitialize()
	require.NoError(t, err)

	// Verify verbose is now 2
	assert.Equal(t, 2, GetVerbose())

	// Cleanup
	Close()
	os.Unsetenv("KFG_VERBOSE")
}

func TestShellLogsNoCorePrefix(t *testing.T) {
	resetLogger()

	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")
	os.Setenv("KFG_LOG_FILE", logFile)
	defer os.Unsetenv("KFG_LOG_FILE")

	err := Initialize()
	require.NoError(t, err)

	// Log from shell (via LogWithSource)
	LogWithSource("info", "shell", "feature:mcps", "shell message", nil)

	Sync()
	Close()

	// Read and parse JSONL
	content, err := os.ReadFile(logFile)
	require.NoError(t, err)

	var event map[string]interface{}
	err = json.Unmarshal([]byte(strings.TrimSpace(string(content))), &event)
	require.NoError(t, err)

	// Verify source field is "shell"
	assert.Equal(t, "shell", event["source"], "Shell logs should have source=\"shell\"")

	// Verify component field does NOT have "core:" prefix for shell logs
	component, ok := event["component"].(string)
	require.True(t, ok, "component should be a string")
	assert.False(t, strings.HasPrefix(component, "core:"), "Shell logs should NOT have 'core:' prefix in component field")
	assert.Equal(t, "feature:mcps", component, "Shell logs should preserve original component")
}

func TestLogFileExtension(t *testing.T) {
	resetLogger()

	// Unset any custom log file/dir settings
	os.Unsetenv("KFG_LOG_FILE")
	os.Unsetenv("KFG_LOG_DIR")

	err := Initialize()
	require.NoError(t, err)

	// Get the log file path
	logPath := GetJSONLPath()
	Close()

	// Verify file ends with .log, not .jsonl
	assert.True(t, strings.HasSuffix(logPath, ".log"), "Log file should have .log extension")
	assert.False(t, strings.HasSuffix(logPath, ".jsonl"), "Log file should NOT have .jsonl extension")
}

func TestGetJSONLPath(t *testing.T) {
	resetLogger()

	// Without initialization
	assert.Equal(t, "", GetJSONLPath())

	// With initialization
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")
	os.Setenv("KFG_LOG_FILE", logFile)
	defer os.Unsetenv("KFG_LOG_FILE")

	err := Initialize()
	require.NoError(t, err)

	assert.Equal(t, logFile, GetJSONLPath())
	Close()
}

// Helper function to reset global logger state
func resetLogger() {
	mu.Lock()
	defer mu.Unlock()

	globalLogger = nil
	jsonlFile = nil
	verbose = 0
	colorMode = ColorAuto
	isTTY = false
	initialized = false
}
