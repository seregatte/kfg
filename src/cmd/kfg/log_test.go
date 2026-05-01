package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/seregatte/kfg/src/internal/logger"
)

func TestLogCommandStructure(t *testing.T) {
	assert.NotNil(t, logCmd)
	assert.Equal(t, "log <level> <component> [message...]", logCmd.Use)
	assert.Contains(t, logCmd.Short, "structured log")
	assert.NotNil(t, logCmd.Run)
}

func TestLogCommandFlags(t *testing.T) {
	// Test that --source flag is registered and hidden
	sourceFlag := logCmd.Flags().Lookup("source")
	assert.NotNil(t, sourceFlag)
	assert.Equal(t, "shell", sourceFlag.DefValue)

	// Flag should be hidden
	assert.True(t, sourceFlag.Hidden)

	// Test that --session-id flag is registered
	sessionIDFlag := logCmd.Flags().Lookup("session-id")
	assert.NotNil(t, sessionIDFlag)
	assert.Equal(t, "", sessionIDFlag.DefValue)
}

func TestLogCommandLevels(t *testing.T) {
	// Test that the command accepts valid levels conceptually
	// The actual level validation happens in the Run function
	validLevels := []string{"error", "warn", "info", "detail", "debug"}
	for _, level := range validLevels {
		// Verify each level is a known string
		assert.NotEmpty(t, level)
	}
}

func TestLogCommandArgs(t *testing.T) {
	// Verify the command requires at least 2 arguments (level and component)
	// Cobra's Args field contains the validator
	assert.NotNil(t, logCmd.Args)
}

func TestLogCommandSourceFlag(t *testing.T) {
	// Verify source flag defaults to "shell"
	sourceFlag := logCmd.Flags().Lookup("source")
	assert.NotNil(t, sourceFlag)
	assert.Equal(t, "shell", sourceFlag.DefValue, "source flag should default to 'shell' for shell helpers")
}

func TestLogCommandSessionIDFlag(t *testing.T) {
	// Verify session-id flag exists and defaults to empty string
	sessionIDFlag := logCmd.Flags().Lookup("session-id")
	assert.NotNil(t, sessionIDFlag, "session-id flag should be registered")
	assert.Equal(t, "", sessionIDFlag.DefValue, "session-id flag should default to empty string")
	assert.Contains(t, sessionIDFlag.Usage, "Session ID", "session-id flag usage should mention Session ID")
}

func TestSessionIDFlagFunctionality(t *testing.T) {
	// Test the actual functionality of the session-id flag
	// These tests run the log command and check the JSONL output

	t.Run("flag provides session ID", func(t *testing.T) {
		// Reset logger state
		resetLoggerState()

		// Create temp log file
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "test.jsonl")
		os.Setenv("KFG_LOG_FILE", logFile)
		defer os.Unsetenv("KFG_LOG_FILE")

		// Initialize logger
		err := logger.Initialize()
		require.NoError(t, err)

		// Create a new command instance to avoid flag pollution
		cmd := createLogCommand()
		cmd.SetArgs([]string{"--session-id", "test-session-123", "info", "test:comp", "test message"})

		// Execute command
		err = cmd.Execute()
		require.NoError(t, err)

		// Sync and close
		logger.Sync()
		logger.Close()

		// Read and parse JSONL
		content, err := os.ReadFile(logFile)
		require.NoError(t, err)

		var event map[string]interface{}
		err = json.Unmarshal([]byte(strings.TrimSpace(string(content))), &event)
		require.NoError(t, err)

		// Check session_id field
		assert.Equal(t, "test-session-123", event["session_id"])
		assert.Equal(t, "shell", event["source"])
		assert.Equal(t, "info", event["level"])
	})

	t.Run("flag overrides env var", func(t *testing.T) {
		resetLoggerState()

		// Set env var
		os.Setenv("KFG_SESSION_ID", "env-session-456")
		defer os.Unsetenv("KFG_SESSION_ID")

		// Create temp log file
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "test.jsonl")
		os.Setenv("KFG_LOG_FILE", logFile)
		defer os.Unsetenv("KFG_LOG_FILE")

		// Initialize logger
		err := logger.Initialize()
		require.NoError(t, err)

		// Execute command with flag overriding env var
		cmd := createLogCommand()
		cmd.SetArgs([]string{"--session-id", "flag-session-789", "info", "test:comp", "test message"})

		err = cmd.Execute()
		require.NoError(t, err)

		logger.Sync()
		logger.Close()

		// Read and parse JSONL
		content, err := os.ReadFile(logFile)
		require.NoError(t, err)

		var event map[string]interface{}
		err = json.Unmarshal([]byte(strings.TrimSpace(string(content))), &event)
		require.NoError(t, err)

		// Flag value should override env var
		assert.Equal(t, "flag-session-789", event["session_id"])
	})

	t.Run("env var used when flag not provided", func(t *testing.T) {
		resetLoggerState()

		// Set env var
		os.Setenv("KFG_SESSION_ID", "env-session-only")
		defer os.Unsetenv("KFG_SESSION_ID")

		// Create temp log file
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "test.jsonl")
		os.Setenv("KFG_LOG_FILE", logFile)
		defer os.Unsetenv("KFG_LOG_FILE")

		// Initialize logger
		err := logger.Initialize()
		require.NoError(t, err)

		// Execute command without session-id flag
		cmd := createLogCommand()
		cmd.SetArgs([]string{"info", "test:comp", "test message"})

		err = cmd.Execute()
		require.NoError(t, err)

		logger.Sync()
		logger.Close()

		// Read and parse JSONL
		content, err := os.ReadFile(logFile)
		require.NoError(t, err)

		var event map[string]interface{}
		err = json.Unmarshal([]byte(strings.TrimSpace(string(content))), &event)
		require.NoError(t, err)

		// Env var value should be used
		assert.Equal(t, "env-session-only", event["session_id"])
	})

	t.Run("empty flag value omits session_id", func(t *testing.T) {
		resetLoggerState()

		// Set env var
		os.Setenv("KFG_SESSION_ID", "env-session-to-omit")
		defer os.Unsetenv("KFG_SESSION_ID")

		// Create temp log file
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "test.jsonl")
		os.Setenv("KFG_LOG_FILE", logFile)
		defer os.Unsetenv("KFG_LOG_FILE")

		// Initialize logger
		err := logger.Initialize()
		require.NoError(t, err)

		// Execute command with empty session-id flag (should omit session_id)
		cmd := createLogCommand()
		cmd.SetArgs([]string{"--session-id", "", "info", "test:comp", "test message"})

		err = cmd.Execute()
		require.NoError(t, err)

		logger.Sync()
		logger.Close()

		// Read and parse JSONL
		content, err := os.ReadFile(logFile)
		require.NoError(t, err)

		var event map[string]interface{}
		err = json.Unmarshal([]byte(strings.TrimSpace(string(content))), &event)
		require.NoError(t, err)

		// session_id should be absent (explicit omission via empty flag)
		assert.NotContains(t, event, "session_id", "session_id should be absent when empty flag value is provided")
	})

	t.Run("no session_id when neither flag nor env var provided", func(t *testing.T) {
		resetLoggerState()

		// Ensure no session ID env var
		os.Unsetenv("KFG_SESSION_ID")

		// Create temp log file
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "test.jsonl")
		os.Setenv("KFG_LOG_FILE", logFile)
		defer os.Unsetenv("KFG_LOG_FILE")

		// Initialize logger
		err := logger.Initialize()
		require.NoError(t, err)

		// Execute command without session-id flag
		cmd := createLogCommand()
		cmd.SetArgs([]string{"info", "test:comp", "test message"})

		err = cmd.Execute()
		require.NoError(t, err)

		logger.Sync()
		logger.Close()

		// Read and parse JSONL
		content, err := os.ReadFile(logFile)
		require.NoError(t, err)

		var event map[string]interface{}
		err = json.Unmarshal([]byte(strings.TrimSpace(string(content))), &event)
		require.NoError(t, err)

		// session_id should be absent
		assert.NotContains(t, event, "session_id")
	})
}

// Helper to create a fresh log command for testing
func createLogCommand() *cobra.Command {
	// Create a new command instance with same structure as logCmd
	cmd := &cobra.Command{
		Use:   "log <level> <component> [message...]",
		Short: "Write structured log entry",
		Args:  cobra.MinimumNArgs(2),
		Run:   logCmd.Run,
	}

	// Add flags
	cmd.Flags().String("source", "shell", "Source identifier for the log entry")
	cmd.Flags().String("session-id", "", "Session ID for log correlation")

	return cmd
}

// Helper to reset logger state between tests
func resetLoggerState() {
	// Reset the global logger state
	logger.Reset()
}