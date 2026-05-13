package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.NotNil(t, cfg)
	assert.Equal(t, 1, cfg.Verbose)
	assert.Equal(t, "", cfg.LogFile)
	assert.Equal(t, "", cfg.LogDir)
	assert.Equal(t, "auto", cfg.LogColor)
}

func TestInitialize(t *testing.T) {
	// Reset viper and ensure clean environment for default value checks
	viper.Reset()
	os.Unsetenv("KFG_VERBOSE")
	os.Unsetenv("KFG_LOG_FILE")
	os.Unsetenv("KFG_LOG_DIR")
	os.Unsetenv("KFG_LOG_COLOR")

	err := Initialize()
	assert.NoError(t, err)

	// Verify default values are set
	assert.Equal(t, 1, viper.GetInt("verbose"))
	assert.Equal(t, "", viper.GetString("log_file"))
	assert.Equal(t, "", viper.GetString("log_dir"))
	assert.Equal(t, "auto", viper.GetString("log_color"))
}

func TestLoad(t *testing.T) {
	// Reset viper and ensure clean environment
	viper.Reset()
	os.Unsetenv("KFG_VERBOSE")

	Initialize()

	cfg := Load()
	assert.NotNil(t, cfg)
	assert.Equal(t, 1, cfg.Verbose)
}

func TestGetVerbose(t *testing.T) {
	// Reset viper and ensure clean environment
	viper.Reset()
	os.Unsetenv("KFG_VERBOSE")
	Initialize()

	// Test default (verbose=1)
	assert.Equal(t, 1, GetVerbose())

	// Test with environment variable
	t.Setenv("KFG_VERBOSE", "2")
	viper.BindEnv("verbose", "KFG_VERBOSE")
	assert.Equal(t, 2, GetVerbose())

	// Test with viper set (string format)
	viper.Reset()
	Initialize()
	viper.Set("verbose", "3")
	assert.Equal(t, 3, GetVerbose())

	// Test with invalid value (should default to 0)
	viper.Set("verbose", "invalid")
	assert.Equal(t, 0, GetVerbose())

	// Test with value > 3 (should clamp to 3)
	viper.Set("verbose", "10")
	assert.Equal(t, 3, GetVerbose())

	// Test with negative value (should default to 0)
	viper.Set("verbose", "-1")
	assert.Equal(t, 0, GetVerbose())
}

func TestGetLogFile(t *testing.T) {
	// Reset viper and ensure clean environment
	viper.Reset()
	os.Unsetenv("KFG_LOG_FILE")
	Initialize()

	// Test default
	assert.Equal(t, "", GetLogFile())

	// Test with environment variable
	t.Setenv("KFG_LOG_FILE", "/tmp/test.jsonl")
	viper.BindEnv("log_file", "KFG_LOG_FILE")
	assert.Equal(t, "/tmp/test.jsonl", GetLogFile())

	// Test with viper set
	viper.Set("log_file", "/var/log/nixai.jsonl")
	assert.Equal(t, "/var/log/nixai.jsonl", GetLogFile())
}

func TestGetLogDir(t *testing.T) {
	// Reset viper and ensure clean environment
	viper.Reset()
	os.Unsetenv("KFG_LOG_DIR")
	Initialize()

	// Test default
	assert.Equal(t, "", GetLogDir())

	// Test with environment variable
	t.Setenv("KFG_LOG_DIR", "/tmp/logs")
	viper.BindEnv("log_dir", "KFG_LOG_DIR")
	assert.Equal(t, "/tmp/logs", GetLogDir())

	// Test with viper set
	viper.Set("log_dir", "/var/log")
	assert.Equal(t, "/var/log", GetLogDir())
}

func TestGetLogColor(t *testing.T) {
	// Reset viper and ensure clean environment
	viper.Reset()
	os.Unsetenv("KFG_LOG_COLOR")
	Initialize()

	// Test default
	assert.Equal(t, "auto", GetLogColor())

	// Test with environment variable
	t.Setenv("KFG_LOG_COLOR", "always")
	viper.BindEnv("log_color", "KFG_LOG_COLOR")
	assert.Equal(t, "always", GetLogColor())

	// Test with viper set
	viper.Set("log_color", "never")
	assert.Equal(t, "never", GetLogColor())
}

func TestGetVerboseFromEnv(t *testing.T) {
	// Ensure clean environment before testing defaults
	os.Unsetenv("KFG_VERBOSE")

	// Test default (no env set)
	assert.Equal(t, 1, GetVerboseFromEnv())

	// Test with valid value
	t.Setenv("KFG_VERBOSE", "2")
	assert.Equal(t, 2, GetVerboseFromEnv())

	// Test with invalid value
	t.Setenv("KFG_VERBOSE", "invalid")
	assert.Equal(t, 0, GetVerboseFromEnv())

	// Test with value > 3
	t.Setenv("KFG_VERBOSE", "10")
	assert.Equal(t, 3, GetVerboseFromEnv())
}

func TestGetLogColorFromEnv(t *testing.T) {
	// Ensure clean environment before testing defaults
	os.Unsetenv("KFG_LOG_COLOR")

	// Test default (no env set)
	assert.Equal(t, "auto", GetLogColorFromEnv())

	// Test with valid value
	t.Setenv("KFG_LOG_COLOR", "always")
	assert.Equal(t, "always", GetLogColorFromEnv())

	// Test with never
	t.Setenv("KFG_LOG_COLOR", "never")
	assert.Equal(t, "never", GetLogColorFromEnv())
}

func TestGetLogFileFromEnv(t *testing.T) {
	// Ensure clean environment before testing defaults
	os.Unsetenv("KFG_LOG_FILE")

	// Test default (no env set)
	assert.Equal(t, "", GetLogFileFromEnv())

	// Test with value
	t.Setenv("KFG_LOG_FILE", "/tmp/test.jsonl")
	assert.Equal(t, "/tmp/test.jsonl", GetLogFileFromEnv())
}

func TestGetLogDirFromEnv(t *testing.T) {
	// Ensure clean environment before testing defaults
	os.Unsetenv("KFG_LOG_DIR")

	// Test default (no env set)
	assert.Equal(t, "", GetLogDirFromEnv())

	// Test with value
	t.Setenv("KFG_LOG_DIR", "/tmp/logs")
	assert.Equal(t, "/tmp/logs", GetLogDirFromEnv())
}

func TestIsDebugMode(t *testing.T) {
	// Reset viper and ensure clean environment
	viper.Reset()
	os.Unsetenv("KFG_VERBOSE")
	Initialize()

	// Test default (verbose=1, so debug mode is false)
	assert.False(t, IsDebugMode())

	// Test with verbose=3 (debug mode is true)
	viper.Set("verbose", "3")
	assert.True(t, IsDebugMode())

	// Test with verbose=2 (debug mode is false)
	viper.Set("verbose", "2")
	assert.False(t, IsDebugMode())
}

func TestSetDebugMode(t *testing.T) {
	// Reset viper and ensure clean environment
	viper.Reset()
	os.Unsetenv("KFG_VERBOSE")
	Initialize()

	// Test setting debug mode to true (should set verbose=3)
	SetDebugMode(true)
	assert.Equal(t, "3", viper.GetString("verbose"))
	assert.True(t, IsDebugMode())

	// Test setting debug mode to false (should set verbose=0)
	SetDebugMode(false)
	assert.Equal(t, "0", viper.GetString("verbose"))
	assert.False(t, IsDebugMode())
}

func TestGetKPath(t *testing.T) {
	// Reset viper and ensure clean environment
	viper.Reset()
	os.Unsetenv("KFG_KPATH")
	Initialize()

	// Test default (empty)
	assert.Equal(t, "", GetKPath())

	// Test with environment variable
	t.Setenv("KFG_KPATH", "./manifests")
	viper.BindEnv("kpath", "KFG_KPATH")
	assert.Equal(t, "./manifests", GetKPath())

	// Test with GitHub URL
	t.Setenv("KFG_KPATH", "https://github.com/owner/repo//manifests")
	viper.BindEnv("kpath", "KFG_KPATH")
	assert.Equal(t, "https://github.com/owner/repo//manifests", GetKPath())

	// Test with viper set
	viper.Set("kpath", "/path/to/manifests")
	assert.Equal(t, "/path/to/manifests", GetKPath())
}

func TestGetKPathFromEnv(t *testing.T) {
	// Ensure clean environment before testing defaults
	os.Unsetenv("KFG_KPATH")

	// Test default (no env set)
	assert.Equal(t, "", GetKPathFromEnv())

	// Test with local path
	t.Setenv("KFG_KPATH", "./manifests")
	assert.Equal(t, "./manifests", GetKPathFromEnv())

	// Test with GitHub URL
	t.Setenv("KFG_KPATH", "https://github.com/owner/repo//manifests")
	assert.Equal(t, "https://github.com/owner/repo//manifests", GetKPathFromEnv())
}
