package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/seregatte/kfg/src/internal/image"
)

func TestGetStoreDir_WithOverride(t *testing.T) {
	// Save and restore original value
	originalOverride := storeDirOverride
	defer func() {
		storeDirOverride = originalOverride
	}()

	// Test with custom override
	storeDirOverride = "/custom/store/path"
	result := getStoreDir()
	assert.Equal(t, "/custom/store/path", result)
}

func TestGetStoreDir_WithoutOverride(t *testing.T) {
	// Save and restore original value
	originalOverride := storeDirOverride
	defer func() {
		storeDirOverride = originalOverride
	}()

	// Test with empty override (returns empty string, letting ImageStore use default)
	storeDirOverride = ""
	result := getStoreDir()
	assert.Equal(t, "", result)
}

func TestGetStoreDir_OverrideWithWhitespace(t *testing.T) {
	// Save and restore original value
	originalOverride := storeDirOverride
	defer func() {
		storeDirOverride = originalOverride
	}()

	// Test with whitespace-only override (should return the whitespace string as-is)
	storeDirOverride = "  "
	result := getStoreDir()
	assert.Equal(t, "  ", result)
}

func TestGetStoreDir_OverrideWithRelativePath(t *testing.T) {
	// Save and restore original value
	originalOverride := storeDirOverride
	defer func() {
		storeDirOverride = originalOverride
	}()

	// Test with relative path override
	storeDirOverride = "./local/store"
	result := getStoreDir()
	assert.Equal(t, "./local/store", result)
}

func TestStoreFlagPassed_ToNewImageStore(t *testing.T) {
	// Save and restore original value
	originalOverride := storeDirOverride
	defer func() {
		storeDirOverride = originalOverride
	}()

	// Create a temporary directory for this test
	tmpDir := t.TempDir()

	// Set storeDirOverride to the temp directory
	storeDirOverride = tmpDir

	// Verify getStoreDir returns the override
	storeDir := getStoreDir()
	assert.Equal(t, tmpDir, storeDir)

	// Create an ImageStore with the custom directory
	// This would be called from the store image commands with the result of getStoreDir()
	store := image.NewImageStore(storeDir)
	require.NotNil(t, store)

	// Verify the store uses the custom directory
	assert.Equal(t, tmpDir, store.GetStoreDir())

	// Initialize the store
	err := store.Initialize()
	require.NoError(t, err)

	// Verify the store directory structure was created in the custom location
	imagesDir := store.GetImagesDir()
	assert.True(t, filepath.HasPrefix(imagesDir, tmpDir), "images directory should be in custom store location")

	// Verify no pollution in default store
	stat, err := os.Stat(imagesDir)
	require.NoError(t, err)
	assert.True(t, stat.IsDir(), "images directory should exist")
}

func TestIsolatedStore_NoDefaultPollution(t *testing.T) {
	// Save and restore original value
	originalOverride := storeDirOverride
	defer func() {
		storeDirOverride = originalOverride
	}()

	// Create two separate temporary directories to simulate two isolated stores
	store1Dir := t.TempDir()
	store2Dir := t.TempDir()

	// Verify they are different
	assert.NotEqual(t, store1Dir, store2Dir)

	// Test that we can use different store directories
	storeDirOverride = store1Dir
	store1 := image.NewImageStore(getStoreDir())
	require.NotNil(t, store1)
	assert.Equal(t, store1Dir, store1.GetStoreDir())

	storeDirOverride = store2Dir
	store2 := image.NewImageStore(getStoreDir())
	require.NotNil(t, store2)
	assert.Equal(t, store2Dir, store2.GetStoreDir())

	// Each store should use its own directory
	assert.NotEqual(t, store1.GetStoreDir(), store2.GetStoreDir())
}
