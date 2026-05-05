package main

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestBuildCommandArgs(t *testing.T) {
	// Test that MaximumNArgs(1) is used (allows 0 or 1 args)
	assert.NotNil(t, buildCmd)
	
	// The Args validator should allow 0 or 1 arguments
	// Test with 0 args (should pass with KFG_KPATH set)
	err := buildCmd.Args(buildCmd, []string{})
	// Should not error with MaximumNArgs(1)
	assert.NoError(t, err)
	
	// Test with 1 arg (should pass)
	err = buildCmd.Args(buildCmd, []string{"./manifests"})
	assert.NoError(t, err)
	
	// Test with 2 args (should fail)
	err = buildCmd.Args(buildCmd, []string{"./manifests", "./other"})
	assert.Error(t, err)
}

func TestBuildCommandKPathFallback(t *testing.T) {
	// Reset viper for each test
	viper.Reset()
	
	// Test 1: KFG_KPATH is set, no argument provided
	os.Setenv("KFG_KPATH", "./test-manifests")
	viper.BindEnv("kpath", "KFG_KPATH")
	
	// The Run function should use GetKPath() when no argument is provided
	// We can't easily test the full Run function without mocking kustomize,
	// but we can verify the config getter works
	assert.Equal(t, "./test-manifests", viper.GetString("kpath"))
	os.Unsetenv("KFG_KPATH")
	
	// Test 2: KFG_KPATH is not set, argument provided
	viper.Reset()
	assert.Equal(t, "", viper.GetString("kpath"))
	
	// Test 3: KFG_KPATH with GitHub URL
	os.Setenv("KFG_KPATH", "https://github.com/owner/repo//manifests")
	viper.BindEnv("kpath", "KFG_KPATH")
	assert.Equal(t, "https://github.com/owner/repo//manifests", viper.GetString("kpath"))
	os.Unsetenv("KFG_KPATH")
}

func TestBuildCommandLongDescription(t *testing.T) {
	// Verify the Long description mentions KFG_KPATH and GitHub URLs
	assert.Contains(t, buildCmd.Long, "KFG_KPATH")
	assert.Contains(t, buildCmd.Long, "github.com")
	assert.Contains(t, buildCmd.Long, "https://github.com/owner/repo//path")
}

func TestBuildCommandExamples(t *testing.T) {
	// Verify the examples include GitHub URL and KFG_KPATH usage
	assert.Contains(t, buildCmd.Long, "kfg build https://github.com/owner/repo//manifests")
	assert.Contains(t, buildCmd.Long, "KFG_KPATH=./manifests kfg build")
	assert.Contains(t, buildCmd.Long, "KFG_KPATH=https://github.com")
}

func TestKustomizeAliasCommand(t *testing.T) {
	// Test that kustomizeCmd alias uses MaximumNArgs(1)
	assert.NotNil(t, kustomizeCmd)
	assert.Equal(t, "kustomize [path-or-url]", kustomizeCmd.Use)
	
	// Test with 0 args (should pass)
	err := kustomizeCmd.Args(kustomizeCmd, []string{})
	assert.NoError(t, err)
	
	// Test with 1 arg (should pass)
	err = kustomizeCmd.Args(kustomizeCmd, []string{"./manifests"})
	assert.NoError(t, err)
	
	// Test with 2 args (should fail)
	err = kustomizeCmd.Args(kustomizeCmd, []string{"./manifests", "./other"})
	assert.Error(t, err)
}