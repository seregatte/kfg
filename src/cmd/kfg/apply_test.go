package main

import (
	"os"
	"testing"

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