// Package logger provides structured logging using Zerolog.
// All logs persist to JSONL file, with optional human-readable stderr output
// controlled by KFG_VERBOSE (0-3).
package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// Level constants for logging.
type Level int

const (
	LevelError  Level = 0
	LevelWarn   Level = 1
	LevelInfo   Level = 2
	LevelDetail Level = 3
	LevelDebug  Level = 4
)

// ColorMode constants for ANSI color output.
type ColorMode string

const (
	ColorAuto   ColorMode = "auto"
	ColorAlways ColorMode = "always"
	ColorNever  ColorMode = "never"
)

// ANSI color codes for level output.
var levelColors = map[Level]string{
	LevelError:  "\x1b[31m", // red
	LevelWarn:   "\x1b[33m", // orange/yellow
	LevelInfo:   "\x1b[32m", // green
	LevelDetail: "\x1b[36m", // cyan
	LevelDebug:  "\x1b[35m", // magenta
}

const colorReset = "\x1b[0m"
const colorWhite = "\x1b[37m"
const colorYellow = "\x1b[33m"

// Global logger state.
var (
	globalLogger *zerolog.Logger
	jsonlFile    *os.File
	verbose      int
	colorMode    ColorMode
	isTTY        bool
	mu           sync.Mutex
	initialized  bool
)

// Context fields from environment.
var contextFields = []string{
	"KFG_WORKFLOW_NAME",
	"KFG_KUSTOMIZATION_NAME",
	"KFG_SESSION_ID",
}

// Initialize sets up the global logger with JSONL file and stderr output.
// It reads configuration directly from environment variables.
func Initialize() error {
	mu.Lock()
	defer mu.Unlock()

	if initialized {
		return nil
	}

	// Configure Zerolog field names to match spec
	zerolog.TimestampFieldName = "ts"
	zerolog.MessageFieldName = "msg"
	zerolog.LevelFieldName = "level"

	// Get configuration from environment
	verbose = getVerboseFromEnv()
	colorMode = getColorModeFromEnv()
	isTTY = isStderrTTY()

	// Set up JSONL file
	filePath := getLogFilePath()
	if err := setupJSONLFile(filePath); err != nil {
		// Fallback: log to stderr only (no persistence)
		globalLogger = setupStderrOnlyLogger()
		initialized = true
		return nil
	}

	// Set up multi-output logger (JSONL + stderr)
	globalLogger = setupMultiOutputLogger()
	initialized = true
	return nil
}

// getVerboseFromEnv returns KFG_VERBOSE value (0-3).
// Default is 1 (error visible).
func getVerboseFromEnv() int {
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

// getColorModeFromEnv returns KFG_LOG_COLOR value.
func getColorModeFromEnv() ColorMode {
	val := os.Getenv("KFG_LOG_COLOR")
	if val == "" {
		return ColorAuto
	}
	switch strings.ToLower(val) {
	case "always":
		return ColorAlways
	case "never":
		return ColorNever
	default:
		return ColorAuto
	}
}

// isStderrTTY returns true if stderr is a terminal.
func isStderrTTY() bool {
	fi, err := os.Stderr.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// getLogFilePath returns the JSONL log file path.
func getLogFilePath() string {
	// Check KFG_LOG_FILE first
	if path := os.Getenv("KFG_LOG_FILE"); path != "" {
		return path
	}

	// Check KFG_LOG_DIR
	if dir := os.Getenv("KFG_LOG_DIR"); dir != "" {
		return filepath.Join(dir, "kfg.log")
	}

	// Default: XDG_STATE_HOME or ~/.local/state
	stateHome := os.Getenv("XDG_STATE_HOME")
	if stateHome == "" {
		home := os.Getenv("HOME")
		if home == "" {
			home = os.Getenv("USERPROFILE") // Windows
		}
		if home != "" {
			stateHome = filepath.Join(home, ".local", "state")
		}
	}

	if stateHome != "" {
		return filepath.Join(stateHome, "kfg", "logs", "kfg.log")
	}

	// Fallback: temp directory
	return filepath.Join(os.TempDir(), "kfg.log")
}

// setupJSONLFile creates the JSONL log file.
func setupJSONLFile(path string) error {
	// Create directory if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Open file for append-only
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}

	jsonlFile = file
	return nil
}

// setupStderrOnlyLogger creates a logger that only writes to stderr.
func setupStderrOnlyLogger() *zerolog.Logger {
	consoleWriter := zerolog.ConsoleWriter{
		Out:         os.Stderr,
		TimeFormat:  time.RFC3339,
		NoColor:     !shouldUseColor(),
		FormatLevel: formatLevelForConsole,
	}
	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()
	return &logger
}

// setupMultiOutputLogger creates a logger that writes to both JSONL and stderr.
func setupMultiOutputLogger() *zerolog.Logger {
	// Create multi-writer that writes to JSONL file and stderr
	writer := newMultiWriter(jsonlFile, newHumanWriter())

	// Create logger with context enrichment (source is added per log call)
	logger := zerolog.New(writer).With().
		Timestamp().
		Int("pid", os.Getpid()).
		Logger()

	// Add context enrichment from environment
	logger = enrichWithContext(logger)

	return &logger
}

// enrichWithContext adds environment context fields to the logger.
func enrichWithContext(logger zerolog.Logger) zerolog.Logger {
	ctx := logger.With()
	for _, envVar := range contextFields {
		val := os.Getenv(envVar)
		if val != "" {
			// Extract field name from env var (remove KFG_ prefix)
			fieldName := strings.ToLower(strings.TrimPrefix(envVar, "KFG_"))
			ctx = ctx.Str(fieldName, val)
		}
	}
	return ctx.Logger()
}

// shouldUseColor returns true if colors should be used.
func shouldUseColor() bool {
	switch colorMode {
	case ColorAlways:
		return true
	case ColorNever:
		return false
	default: // ColorAuto
		return isTTY
	}
}

// formatLevelForConsole formats the level for console output.
func formatLevelForConsole(i interface{}) string {
	if ll, ok := i.(string); ok {
		switch ll {
		case "error":
			if shouldUseColor() {
				return fmt.Sprintf("%sERROR%s", levelColors[LevelError], colorReset)
			}
			return "ERROR"
		case "warn":
			if shouldUseColor() {
				return fmt.Sprintf("%sWARN%s", levelColors[LevelWarn], colorReset)
			}
			return "WARN"
		case "info":
			if shouldUseColor() {
				return fmt.Sprintf("%sINFO%s", levelColors[LevelInfo], colorReset)
			}
			return "INFO"
		case "detail":
			if shouldUseColor() {
				return fmt.Sprintf("%sDETAIL%s", levelColors[LevelDetail], colorReset)
			}
			return "DETAIL"
		case "debug":
			if shouldUseColor() {
				return fmt.Sprintf("%sDEBUG%s", levelColors[LevelDebug], colorReset)
			}
			return "DEBUG"
		}
	}
	return strings.ToUpper(fmt.Sprintf("%v", i))
}

// humanWriter is a custom writer for human-readable stderr output.
type humanWriter struct{}

func newHumanWriter() *humanWriter {
	return &humanWriter{}
}

func (w *humanWriter) Write(p []byte) (n int, err error) {
	// Parse JSONL event
	var event map[string]interface{}
	if err := json.Unmarshal(p, &event); err != nil {
		// Not valid JSON, write raw
		return os.Stderr.Write(p)
	}

	// Check verbose gating
	levelStr, ok := event["level"].(string)
	if !ok {
		return os.Stderr.Write(p)
	}

	level := parseLevelString(levelStr)
	if !shouldShowLevel(level) {
		return len(p), nil // consumed but not displayed
	}

	// Format human output: [LEVEL][component] message
	component, _ := event["component"].(string)
	msg, _ := event["msg"].(string)

	var output string
	if shouldUseColor() {
		output = formatColoredOutput(levelStr, component, msg)
	} else {
		output = formatPlainOutput(levelStr, component, msg)
	}

	return fmt.Fprintln(os.Stderr, output)
}

// multiWriter writes to both JSONL file and stderr.
type multiWriter struct {
	jsonlWriter  io.Writer
	stderrWriter io.Writer
}

func newMultiWriter(jsonl io.Writer, stderr io.Writer) *multiWriter {
	return &multiWriter{
		jsonlWriter:  jsonl,
		stderrWriter: stderr,
	}
}

func (w *multiWriter) Write(p []byte) (n int, err error) {
	// Always write to JSONL file
	if w.jsonlWriter != nil {
		if _, err := w.jsonlWriter.Write(p); err != nil {
			// Log file write failed, but continue
		}
	}

	// Write to stderr writer (conditional display logic inside humanWriter)
	if w.stderrWriter != nil {
		return w.stderrWriter.Write(p)
	}

	return len(p), nil
}

// parseLevelString converts a level string to Level constant.
func parseLevelString(s string) Level {
	switch s {
	case "error":
		return LevelError
	case "warn":
		return LevelWarn
	case "info":
		return LevelInfo
	case "detail":
		return LevelDetail
	case "debug":
		return LevelDebug
	default:
		return LevelInfo
	}
}

// shouldShowLevel returns true if the level should be shown at current verbose.
func shouldShowLevel(level Level) bool {
	// verbose=0: no human output (JSONL only)
	// verbose=1: error only
	// verbose=2: error + warn + info
	// verbose=3: all levels
	switch verbose {
	case 0:
		return false
	case 1:
		return level == LevelError
	case 2:
		return level <= LevelInfo
	case 3:
		return true
	default:
		return false
	}
}

// formatColoredOutput formats output with ANSI colors.
func formatColoredOutput(level, component, msg string) string {
	levelColor := levelColors[parseLevelString(level)]
	return fmt.Sprintf("[%s%s%s][%s%s%s] %s%s%s",
		levelColor, strings.ToUpper(level), colorReset,
		colorWhite, component, colorReset,
		colorYellow, msg, colorReset)
}

// formatPlainOutput formats output without colors.
func formatPlainOutput(level, component, msg string) string {
	return fmt.Sprintf("[%s][%s] %s",
		strings.ToUpper(level), component, msg)
}

// Close closes the JSONL file.
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if jsonlFile != nil {
		err := jsonlFile.Close()
		jsonlFile = nil
		return err
	}
	return nil
}

// Reset resets the global logger state. This is primarily for testing.
// After calling Reset, Initialize() will re-initialize the logger.
func Reset() {
	mu.Lock()
	defer mu.Unlock()

	if jsonlFile != nil {
		jsonlFile.Close()
		jsonlFile = nil
	}
	globalLogger = nil
	verbose = 0
	colorMode = ColorAuto
	isTTY = false
	initialized = false
}

// Reinitialize resets and re-initializes the logger with current environment variables.
// This is useful when environment variables change after initial initialization,
// such as when the --verbose flag is parsed.
func Reinitialize() error {
	Reset()
	return Initialize()
}

// Sync flushes any buffered log entries.
func Sync() error {
	mu.Lock()
	defer mu.Unlock()

	if jsonlFile != nil {
		return jsonlFile.Sync()
	}
	return nil
}

// Get returns the global logger instance.
// If Initialize hasn't been called, it returns a no-op logger.
func Get() *zerolog.Logger {
	mu.Lock()
	defer mu.Unlock()

	if globalLogger == nil {
		// Return a no-op logger if not initialized
		nop := zerolog.Nop()
		return &nop
	}
	return globalLogger
}

// log writes a log entry at the specified level.
// Go code logs automatically get a "core:" prefix on the component.
func log(level zerolog.Level, levelName string, component string, msg string) {
	logger := Get()

	// Add "core:" prefix for Go logs to distinguish from shell logs
	prefixedComponent := "core:" + component

	// Map zerolog levels to our custom levels
	// zerolog: trace=-1, debug=0, info=1, warn=2, error=3
	// We need to handle "detail" as a custom level between info and debug
	switch levelName {
	case "error":
		logger.Error().Str("source", "go").Str("component", prefixedComponent).Msg(msg)
	case "warn":
		logger.Warn().Str("source", "go").Str("component", prefixedComponent).Msg(msg)
	case "info":
		logger.Info().Str("source", "go").Str("component", prefixedComponent).Msg(msg)
	case "detail":
		// Zerolog doesn't have detail, use debug with detail marker
		logger.Debug().Str("source", "go").Str("component", prefixedComponent).Str("level", "detail").Msg(msg)
	case "debug":
		logger.Debug().Str("source", "go").Str("component", prefixedComponent).Msg(msg)
	}
}

// Error logs an error message.
func Error(component string, msg string) {
	log(zerolog.ErrorLevel, "error", component, msg)
}

// Warn logs a warning message.
func Warn(component string, msg string) {
	log(zerolog.WarnLevel, "warn", component, msg)
}

// Info logs an info message.
func Info(component string, msg string) {
	log(zerolog.InfoLevel, "info", component, msg)
}

// Detail logs a detail message (between info and debug).
func Detail(component string, msg string) {
	log(zerolog.DebugLevel, "detail", component, msg)
}

// Debug logs a debug message.
func Debug(component string, msg string) {
	log(zerolog.DebugLevel, "debug", component, msg)
}

// LogWithSource logs a message with a custom source field.
// This is used by the `kfg sys log` CLI command to set source="shell".
// If sessionID is provided (non-nil), it overrides the KFG_SESSION_ID environment variable.
// If sessionID is an empty string, it explicitly omits the session_id field.
func LogWithSource(level string, source string, component string, msg string, sessionID *string) {
	logger := Get()

	// Handle session ID:
	// - nil: use enriched value from environment variable (no override needed)
	// - non-nil: use explicit session_id control (override or omit)

	if sessionID != nil {
		// Flag was provided - use explicit session_id control
		// This ensures the flag value takes precedence over any env var
		var explicitSessionID *string
		if *sessionID != "" {
			explicitSessionID = sessionID
		}
		// If sessionID is empty, explicitSessionID is nil -> omit session_id
		// If sessionID is non-empty, explicitSessionID has the value -> use it
		logWithExplicitSessionID(logger, level, source, component, msg, explicitSessionID)
		return
	}

	// Normal case (flag not provided): build context with source, inherit enriched session_id
	ctx := logger.With().Str("source", source)
	newLogger := ctx.Logger()

	// Log based on level
	switch level {
	case "error":
		newLogger.Error().Str("component", component).Msg(msg)
	case "warn":
		newLogger.Warn().Str("component", component).Msg(msg)
	case "info":
		newLogger.Info().Str("component", component).Msg(msg)
	case "detail":
		newLogger.Debug().Str("component", component).Str("level", "detail").Msg(msg)
	case "debug":
		newLogger.Debug().Str("component", component).Msg(msg)
	}
}

// logWithExplicitSessionID logs with explicit control over session_id.
// When sessionID is nil, session_id is omitted entirely.
// When sessionID is non-nil, it's used as the session_id value.
func logWithExplicitSessionID(baseLogger *zerolog.Logger, level string, source string, component string, msg string, sessionID *string) {
	mu.Lock()
	writer := newMultiWriter(jsonlFile, newHumanWriter())
	mu.Unlock()

	// Build a fresh context with the fields we want, excluding session_id from enrichment
	ctx := zerolog.New(writer).With().
		Timestamp().
		Int("pid", os.Getpid()).
		Str("source", source)

	// Add other context fields from environment (excluding session_id)
	for _, envVar := range []string{"KFG_WORKFLOW_NAME", "KFG_KUSTOMIZATION_NAME"} {
		val := os.Getenv(envVar)
		if val != "" {
			fieldName := strings.ToLower(strings.TrimPrefix(envVar, "KFG_"))
			ctx = ctx.Str(fieldName, val)
		}
	}

	// Add session_id only if provided (non-nil and non-empty)
	if sessionID != nil && *sessionID != "" {
		ctx = ctx.Str("session_id", *sessionID)
	}

	newLogger := ctx.Logger()

	// Log based on level
	switch level {
	case "error":
		newLogger.Error().Str("component", component).Msg(msg)
	case "warn":
		newLogger.Warn().Str("component", component).Msg(msg)
	case "info":
		newLogger.Info().Str("component", component).Msg(msg)
	case "detail":
		newLogger.Debug().Str("component", component).Str("level", "detail").Msg(msg)
	case "debug":
		newLogger.Debug().Str("component", component).Msg(msg)
	}
}

// GetVerbose returns the current verbose level.
func GetVerbose() int {
	return verbose
}

// IsColorEnabled returns true if color output is enabled.
func IsColorEnabled() bool {
	return shouldUseColor()
}

// GetJSONLPath returns the path to the JSONL file.
func GetJSONLPath() string {
	if jsonlFile != nil {
		return jsonlFile.Name()
	}
	return ""
}
