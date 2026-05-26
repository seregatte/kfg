// Package cache provides Step cache operations for kfg.
// It handles identity computation, metadata management, artifact storage,
// filesystem diff, and atomic cache writes.
package cache

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
)

// ComputeIdentity computes the cache identity hash for a StepReference.name.
// The identity is SHA256(name) used as the directory name on disk.
func ComputeIdentity(stepRefName string) string {
	hash := sha256.Sum256([]byte(stepRefName))
	return fmt.Sprintf("%x", hash)
}

// GetCacheDir returns the cache directory path.
// It uses KFG_STORE_DIR/cache or defaults to ~/.kfg/store/cache.
func GetCacheDir() string {
	storeDir := os.Getenv("KFG_STORE_DIR")
	if storeDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = "/tmp"
		}
		storeDir = filepath.Join(homeDir, ".kfg", "store")
	}
	return filepath.Join(storeDir, "cache")
}

// GetEntryPath returns the full path to a cache entry for a given StepReference.name.
func GetEntryPath(stepRefName string) string {
	identity := ComputeIdentity(stepRefName)
	return filepath.Join(GetCacheDir(), identity)
}

// Exists checks if a valid cache entry exists for the given StepReference.name.
// A valid entry has both the directory and metadata.yaml file.
func Exists(stepRefName string) bool {
	entryPath := GetEntryPath(stepRefName)
	metadataPath := filepath.Join(entryPath, "metadata.yaml")

	if _, err := os.Stat(entryPath); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		return false
	}
	return true
}
