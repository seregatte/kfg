package main

import (
	"os"
	"testing"

	"github.com/seregatte/kfg/src/internal/manifest"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestApplyCommandKPathFallback(t *testing.T) {
	// Reset viper for each test
	viper.Reset()

	// Test 1: KFG_KPATH is set, no -k or -f flag provided
	os.Setenv("KFG_KPATH", "./test-manifests")
	viper.BindEnv("kpath", "KFG_KPATH")

	// The Run function should use GetKPath() when no -k or -f is provided
	// We can verify the config getter works
	assert.Equal(t, "./test-manifests", viper.GetString("kpath"))
	os.Unsetenv("KFG_KPATH")

	// Test 2: KFG_KPATH is not set, -k flag provided
	viper.Reset()
	assert.Equal(t, "", viper.GetString("kpath"))

	// Test 3: KFG_KPATH with GitHub URL
	os.Setenv("KFG_KPATH", "https://github.com/owner/repo//manifests")
	viper.BindEnv("kpath", "KFG_KPATH")
	assert.Equal(t, "https://github.com/owner/repo//manifests", viper.GetString("kpath"))
	os.Unsetenv("KFG_KPATH")
}

func TestApplyCommandArgs(t *testing.T) {
	// Test that MaximumNArgs(1) is used (allows 0 or 1 args)
	assert.NotNil(t, applyCmd)

	// Test with 0 args (should pass with KFG_KPATH set)
	err := applyCmd.Args(applyCmd, []string{})
	assert.NoError(t, err)

	// Test with 1 arg (should pass)
	err = applyCmd.Args(applyCmd, []string{"./manifests"})
	assert.NoError(t, err)

	// Test with 2 args (should fail)
	err = applyCmd.Args(applyCmd, []string{"./manifests", "./other"})
	assert.Error(t, err)
}

func TestApplyCommandLongDescription(t *testing.T) {
	// Verify the Long description mentions KFG_KPATH and GitHub URLs
	assert.Contains(t, applyCmd.Long, "KFG_KPATH")
	assert.Contains(t, applyCmd.Long, "github.com")
	assert.Contains(t, applyCmd.Long, "https://github.com/owner/repo//path")
}

func TestApplyCommandExamples(t *testing.T) {
	// Verify the examples include GitHub URL and KFG_KPATH usage
	assert.Contains(t, applyCmd.Long, "kfg apply -k https://github.com/owner/repo//manifests")
	assert.Contains(t, applyCmd.Long, "KFG_KPATH=./manifests kfg apply")
	assert.Contains(t, applyCmd.Long, "KFG_KPATH=https://github.com")
}

// --- Tasks 3.1-3.8: Apply command conversion flag tests ---

func TestValidateWithFlagMutualExclusivity(t *testing.T) {
	// 3.5: --with and --use mutual exclusivity
	err := validateWithFlag(".key", "my-asset", "my-converter", "manifest.yaml", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mutually exclusive")

	// No error when --with is used without --use
	err = validateWithFlag(".key", "my-asset", "", "manifest.yaml", "", "")
	assert.NoError(t, err)
}

func TestValidateWithFlagRequiresConvertOrStdin(t *testing.T) {
	// 3.8: --with without --convert and without stdin (exit code 2)
	err := validateWithFlag(".key", "", "", "manifest.yaml", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--with requires --convert or -f -")

	// No error when --with is used with --convert
	err = validateWithFlag(".key", "my-asset", "", "manifest.yaml", "", "")
	assert.NoError(t, err)

	// No error when --with is used with -f - (stdin)
	err = validateWithFlag(".key", "", "", "-", "", "")
	assert.NoError(t, err)
}

func TestValidateWithFlagNoWorkflowOrCmds(t *testing.T) {
	// --with cannot be used with -w/--workflow
	err := validateWithFlag(".key", "my-asset", "", "manifest.yaml", "my-workflow", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--workflow")

	// --with cannot be used with -c/--cmds
	err = validateWithFlag(".key", "my-asset", "", "manifest.yaml", "", "my-cmd")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--cmds")
}

func TestValidateWithFlagEmpty(t *testing.T) {
	// No error when --with is not specified
	err := validateWithFlag("", "", "", "", "", "")
	assert.NoError(t, err)
}

func TestRunConversionWithAssetName(t *testing.T) {
	// 3.1: --convert with Asset name (existing behavior unchanged)
	// Create test manifest with Asset and Converter
	manifestYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test-asset
spec:
  input:
    format: json
  data: '{"key":"value","nested":{"a":1}}'
---
apiVersion: kfg.dev/v1alpha1
kind: Converter
metadata:
  name: test-converter
spec:
  input:
    format: json
  output:
    format: json
  engine:
    expression: .key
`
	parser := manifest.NewParser()
	resources, err := parser.ParseData("test.yaml", []byte(manifestYAML))
	if err != nil {
		t.Fatalf("failed to parse manifest: %v", err)
	}

	// Run conversion with Asset name
	err = runConversion(resources, "test-asset", "test-converter", "", "")
	assert.NoError(t, err)
}

func TestRunConversionWithRawJSONString(t *testing.T) {
	// 3.2: --convert with raw JSON string (fallback when no Asset match)
	manifestYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Converter
metadata:
  name: test-converter
spec:
  input:
    format: json
  output:
    format: json
  engine:
    expression: .key
`
	parser := manifest.NewParser()
	resources, err := parser.ParseData("test.yaml", []byte(manifestYAML))
	if err != nil {
		t.Fatalf("failed to parse manifest: %v", err)
	}

	// Run conversion with raw JSON string (no matching Asset)
	err = runConversion(resources, `{"key":"hello"}`, "test-converter", "", "")
	assert.NoError(t, err)
}

func TestRunConversionWithInlineExpression(t *testing.T) {
	// 3.3: --with inline expression with Asset lookup
	manifestYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test-asset
spec:
  input:
    format: yaml
  data:
    key: value
    nested:
      a: 1
`
	parser := manifest.NewParser()
	resources, err := parser.ParseData("test.yaml", []byte(manifestYAML))
	if err != nil {
		t.Fatalf("failed to parse manifest: %v", err)
	}

	// Run conversion with Asset name and inline expression
	err = runConversion(resources, "test-asset", "", ".key", "")
	assert.NoError(t, err)
}

func TestRunConversionWithRawStringAndInlineExpression(t *testing.T) {
	// 3.4: --with with raw string input (no Asset)
	manifestYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: unrelated-asset
spec:
  input:
    format: yaml
  data:
    key: other
`
	parser := manifest.NewParser()
	resources, err := parser.ParseData("test.yaml", []byte(manifestYAML))
	if err != nil {
		t.Fatalf("failed to parse manifest: %v", err)
	}

	// Run conversion with raw JSON and inline expression (no matching Asset)
	err = runConversion(resources, `{"key":"world"}`, "", ".key", "")
	assert.NoError(t, err)
}

func TestRunStdinConversion(t *testing.T) {
	// 3.6: -f - with --with stdin pipeline
	err := runStdinConversion(`{"key":"stdin-value"}`, ".key", "")
	assert.NoError(t, err)
}

func TestRunStdinConversionWithOutput(t *testing.T) {
	// 3.7: -f - with --with and -o output file
	tmpFile := t.TempDir() + "/output.json"
	err := runStdinConversion(`{"a":1,"b":2}`, ".", tmpFile)
	assert.NoError(t, err)

	// Verify file was written
	data, err := os.ReadFile(tmpFile)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "a")
	assert.Contains(t, string(data), "b")

	os.Remove(tmpFile)
}

func TestRunConversionAssetNotFoundNoFallback(t *testing.T) {
	// When --convert value doesn't match any Asset and is not valid JSON/YAML,
	// it should return an error listing available assets
	manifestYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: my-asset
spec:
  input:
    format: yaml
  data:
    key: value
`
	parser := manifest.NewParser()
	resources, err := parser.ParseData("test.yaml", []byte(manifestYAML))
	if err != nil {
		t.Fatalf("failed to parse manifest: %v", err)
	}

	// No converter or --with provided, asset not found, input not valid JSON
	err = runConversion(resources, "not-valid-json-or-asset", "", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Contains(t, err.Error(), "my-asset") // lists available assets
}

func TestWriteOutputToFile(t *testing.T) {
	tmpFile := t.TempDir() + "/test-output.txt"
	err := writeOutput("hello world", tmpFile)
	assert.NoError(t, err)

	data, err := os.ReadFile(tmpFile)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", string(data))

	os.Remove(tmpFile)
}

func TestWriteOutputToStdout(t *testing.T) {
	// No error when writing to stdout (empty outputFile)
	err := writeOutput("test output", "")
	assert.NoError(t, err)
}

func TestApplyCommandRefreshFlag(t *testing.T) {
	// Test that --refresh flag is registered
	flags := applyCmd.Flags()

	refreshFlag := flags.Lookup("refresh")
	assert.NotNil(t, refreshFlag)
	assert.Equal(t, "r", refreshFlag.Shorthand)
	assert.Equal(t, "false", refreshFlag.DefValue)
}
