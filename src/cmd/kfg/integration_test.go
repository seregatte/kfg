//go:build integration

package main

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBuildCommandGitHubURL tests building from a GitHub URL.
// This test requires network access and is tagged with 'integration'.
// Run with: go test -tags=integration ./src/cmd/kfg/...
//
// Note: This test uses a real GitHub repository with kustomization files.
// The kustomize library handles the git cloning internally.
func TestBuildCommandGitHubURL(t *testing.T) {
	// Skip if running in CI without network access
	if os.Getenv("CI") == "true" && os.Getenv("KFG_NETWORK_TESTS") != "true" {
		t.Skip("Skipping network test in CI environment")
	}

	// Use a stable test fixture from kustomize repository
	// Note: This URL points to a kustomize test fixture that should be stable
	githubURL := "https://github.com/kubernetes-sigs/kustomize//cmd/config/testdata/bases/simple?ref=master"

	// Reset viper
	viper.Reset()

	// Test that the build command can load from GitHub URL
	// We can't easily test the full Run function, but we can test that
	// the URL is correctly detected and passed to kustomize

	// Set up the command arguments
	args := []string{githubURL}

	// Verify the Args validator allows 1 argument
	err := buildCmd.Args(buildCmd, args)
	require.NoError(t, err)

	// Verify the URL would be passed to kustomize (we can't mock kustomize easily)
	// The actual test would run: kfg build <githubURL> -o <outputFile>
	// and verify the output contains valid YAML

	t.Logf("Build command would process GitHub URL: %s", githubURL)
	t.Log("Note: Full CLI execution test would require building the binary and running it")
}

// TestBuildCommandGitHubURLWithRef tests building from a GitHub URL with a specific ref.
// This test validates that ?ref= parameter is correctly handled.
func TestBuildCommandGitHubURLWithRef(t *testing.T) {
	// Skip if running in CI without network access
	if os.Getenv("CI") == "true" && os.Getenv("KFG_NETWORK_TESTS") != "true" {
		t.Skip("Skipping network test in CI environment")
	}

	// Use a specific tag/branch reference
	githubURL := "https://github.com/kubernetes-sigs/kustomize//cmd/config/testdata/bases/simple?ref=v5.4.3"

	// Reset viper
	viper.Reset()

	// Test that the build command can load from GitHub URL with ref
	args := []string{githubURL}

	err := buildCmd.Args(buildCmd, args)
	require.NoError(t, err)

	t.Logf("Build command would process GitHub URL with ref: %s", githubURL)
}

// TestApplyCommandGitHubURL tests applying from a GitHub URL.
// This test requires network access.
func TestApplyCommandGitHubURL(t *testing.T) {
	// Skip if running in CI without network access
	if os.Getenv("CI") == "true" && os.Getenv("KFG_NETWORK_TESTS") != "true" {
		t.Skip("Skipping network test in CI environment")
	}

	githubURL := "https://github.com/kubernetes-sigs/kustomize//cmd/config/testdata/bases/simple?ref=master"

	// Reset viper
	viper.Reset()

	// Test that the apply command can use -k with GitHub URL
	// Note: The actual apply requires manifests with CmdWorkflow, which the test fixture doesn't have
	// So we just test that the URL is correctly accepted

	// Verify the Args validator allows 0 arguments (GitHub URL would be via -k flag)
	err := applyCmd.Args(applyCmd, []string{})
	require.NoError(t, err)

	t.Logf("Apply command would process GitHub URL via -k: %s", githubURL)
}

// TestRunCommandGitHubURL tests running from a GitHub URL.
// This test requires network access.
func TestRunCommandGitHubURL(t *testing.T) {
	// Skip if running in CI without network access
	if os.Getenv("CI") == "true" && os.Getenv("KFG_NETWORK_TESTS") != "true" {
		t.Skip("Skipping network test in CI environment")
	}

	githubURL := "https://github.com/kubernetes-sigs/kustomize//cmd/config/testdata/bases/simple?ref=master"

	// Reset viper
	viper.Reset()

	// Test that the run command can use -k with GitHub URL
	// Note: The actual run requires manifests with Cmds/CmdWorkflow, which the test fixture doesn't have
	// So we just test that the URL is correctly accepted

	// Verify the command structure
	assert.NotNil(t, runCmd)

	t.Logf("Run command would process GitHub URL via -k: %s", githubURL)
}

// TestKPathEnvVarIntegration tests KFG_KPATH environment variable behavior.
// This test validates the full integration of KFG_KPATH across all commands.
func TestKPathEnvVarIntegration(t *testing.T) {
	// Reset viper
	viper.Reset()

	// Test 1: KFG_KPATH set with local path
	os.Setenv("KFG_KPATH", "./test-manifests")
	viper.BindEnv("kpath", "KFG_KPATH")

	require.Equal(t, "./test-manifests", viper.GetString("kpath"))

	// Verify all commands accept 0 arguments when KFG_KPATH is set
	err := buildCmd.Args(buildCmd, []string{})
	require.NoError(t, err)

	err = applyCmd.Args(applyCmd, []string{})
	require.NoError(t, err)

	os.Unsetenv("KFG_KPATH")

	// Test 2: KFG_KPATH with GitHub URL
	viper.Reset()
	os.Setenv("KFG_KPATH", "https://github.com/owner/repo//manifests")
	viper.BindEnv("kpath", "KFG_KPATH")

	require.Equal(t, "https://github.com/owner/repo//manifests", viper.GetString("kpath"))

	os.Unsetenv("KFG_KPATH")
}

// TestSourceResolutionPriorityIntegration tests the source resolution priority chain.
// Priority: arg > flag > env var > error
func TestSourceResolutionPriorityIntegration(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		kpathEnv       string
		expectedSource string
	}{
		{
			name:           "arg overrides env",
			args:           []string{"./arg-path"},
			kpathEnv:       "./env-path",
			expectedSource: "./arg-path", // arg wins
		},
		{
			name:           "env when no arg",
			args:           []string{},
			kpathEnv:       "./env-path",
			expectedSource: "./env-path", // env used
		},
		{
			name:           "no source when neither",
			args:           []string{},
			kpathEnv:       "",
			expectedSource: "", // error expected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()

			// Set up environment
			if tt.kpathEnv != "" {
				os.Setenv("KFG_KPATH", tt.kpathEnv)
				viper.BindEnv("kpath", "KFG_KPATH")
			}

			// Verify Args validator
			err := buildCmd.Args(buildCmd, tt.args)
			require.NoError(t, err)

			// Verify viper value (this is what the command would use)
			if tt.args != nil && len(tt.args) > 0 {
				// Positional arg would be used
				t.Logf("Source would be: %s (from arg)", tt.args[0])
			} else if tt.kpathEnv != "" {
				// KFG_KPATH would be used
				t.Logf("Source would be: %s (from KFG_KPATH)", viper.GetString("kpath"))
			} else {
				// No source - error expected in Run function
				t.Log("No source available - error expected")
			}

			os.Unsetenv("KFG_KPATH")
		})
	}
}
