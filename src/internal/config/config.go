// Package config provides configuration management using viper.
package config

import (
	"os"
	"strconv"

	"github.com/spf13/viper"
)

// Config holds application configuration.
type Config struct {
	// Verbose controls logging verbosity (0-3).
	// 0 = errors only to JSONL, no human output
	// 1 = errors, warnings, info to JSONL and human output
	// 2 = + detail level
	// 3 = + debug level
	Verbose int

	// LogFile is the path to the JSONL log file.
	// If empty, uses default location.
	LogFile string

	// LogDir is the directory for the log file.
	// If empty, uses default location.
	LogDir string

	// LogColor controls ANSI color output.
	// "auto" = enable if stderr is TTY
	// "always" = always enable colors
	// "never" = always disable colors
	LogColor string

	// StoreDir is the directory for the artifact cache store.
	// If empty, uses default location (~/.config/kfg/store).
	StoreDir string
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Verbose:   1,
		LogFile:   "",
		LogDir:    "",
		LogColor:  "auto",
		StoreDir:  "",
	}
}

// Initialize sets up viper with default values and environment variable bindings.
func Initialize() error {
	// Set default values
	viper.SetDefault("verbose", 1)
	viper.SetDefault("log_file", "")
	viper.SetDefault("log_dir", "")
	viper.SetDefault("log_color", "auto")
	viper.SetDefault("store_dir", "") // Will be computed dynamically if empty

	// Bind environment variables
	viper.BindEnv("verbose", "KFG_VERBOSE")
	viper.BindEnv("log_file", "KFG_LOG_FILE")
	viper.BindEnv("log_dir", "KFG_LOG_DIR")
	viper.BindEnv("log_color", "KFG_LOG_COLOR")
	viper.BindEnv("debug", "KFG_DEBUG")
	viper.BindEnv("store_dir", "KFG_STORE_DIR")

	// Set environment variable prefix
	viper.SetEnvPrefix("KFG")

	// Allow viper to read environment variables with prefix
	viper.AutomaticEnv()

	return nil
}

// Load returns the current configuration from viper.
func Load() *Config {
	return &Config{
		Verbose:   GetVerbose(),
		LogFile:   viper.GetString("log_file"),
		LogDir:    viper.GetString("log_dir"),
		LogColor:  viper.GetString("log_color"),
		StoreDir:  viper.GetString("store_dir"),
	}
}

// GetVerbose returns the verbose level (0-3).
// Parses KFG_VERBOSE as an integer, defaults to 1.
func GetVerbose() int {
	str := viper.GetString("verbose")
	if str == "" {
		return 1
	}
	val, err := strconv.Atoi(str)
	if err != nil || val < 0 {
		return 0
	}
	if val > 3 {
		return 3
	}
	return val
}

// GetLogFile returns the log file path configuration.
func GetLogFile() string {
	return viper.GetString("log_file")
}

// GetLogDir returns the log directory configuration.
func GetLogDir() string {
	return viper.GetString("log_dir")
}

// GetLogColor returns the log color mode configuration.
func GetLogColor() string {
	return viper.GetString("log_color")
}

// IsDebugMode returns whether debug mode is enabled.
// Deprecated: Use GetVerbose() >= 3 instead.
func IsDebugMode() bool {
	return GetVerbose() >= 3
}

// SetDebugMode sets the debug mode flag.
// Deprecated: Use SetVerbose() instead.
func SetDebugMode(enabled bool) {
	if enabled {
		viper.Set("verbose", "3")
	} else {
		viper.Set("verbose", "0")
	}
}

// GetVerboseFromEnv returns the KFG_VERBOSE value directly from environment.
// This is useful before viper is initialized.
// Default is 1 (error visible).
func GetVerboseFromEnv() int {
	str := os.Getenv("KFG_VERBOSE")
	if str == "" {
		return 1
	}
	val, err := strconv.Atoi(str)
	if err != nil || val < 0 {
		return 0
	}
	if val > 3 {
		return 3
	}
	return val
}

// GetLogColorFromEnv returns the KFG_LOG_COLOR value directly from environment.
// This is useful before viper is initialized.
func GetLogColorFromEnv() string {
	val := os.Getenv("KFG_LOG_COLOR")
	if val == "" {
		return "auto"
	}
	return val
}

// GetLogFileFromEnv returns the KFG_LOG_FILE value directly from environment.
func GetLogFileFromEnv() string {
	return os.Getenv("KFG_LOG_FILE")
}

// GetLogDirFromEnv returns the KFG_LOG_DIR value directly from environment.
func GetLogDirFromEnv() string {
	return os.Getenv("KFG_LOG_DIR")
}

// GetStoreDir returns the store directory configuration.
func GetStoreDir() string {
	return viper.GetString("store_dir")
}

// GetStoreDirFromEnv returns the KFG_STORE_DIR value directly from environment.
func GetStoreDirFromEnv() string {
	return os.Getenv("KFG_STORE_DIR")
}