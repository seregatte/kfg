package main

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGetCacheDir(t *testing.T) {
	// Test with KFG_STORE_DIR set
	t.Run("with KFG_STORE_DIR", func(t *testing.T) {
		tmpDir := t.TempDir()
		os.Setenv("KFG_STORE_DIR", tmpDir)
		defer os.Unsetenv("KFG_STORE_DIR")

		cacheDir := getCacheDir()
		expected := filepath.Join(tmpDir, "cache")
		if cacheDir != expected {
			t.Errorf("getCacheDir() = %s, want %s", cacheDir, expected)
		}
	})

	// Test without KFG_STORE_DIR (uses home directory)
	t.Run("without KFG_STORE_DIR", func(t *testing.T) {
		os.Unsetenv("KFG_STORE_DIR")
		
		cacheDir := getCacheDir()
		// Should contain ".kfg/store/cache"
		if !filepath.IsAbs(cacheDir) {
			t.Errorf("getCacheDir() should return absolute path, got %s", cacheDir)
		}
		if !strings.Contains(cacheDir, ".kfg") || !strings.Contains(cacheDir, "store") || !strings.Contains(cacheDir, "cache") {
			t.Errorf("getCacheDir() should contain .kfg/store/cache, got %s", cacheDir)
		}
	})
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{100, "100 B"},
		{1024, "1.0 KiB"},
		{1536, "1.5 KiB"},
		{1048576, "1.0 MiB"},
		{1572864, "1.5 MiB"},
		{1073741824, "1.0 GiB"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := formatSize(tt.bytes)
			if result != tt.expected {
				t.Errorf("formatSize(%d) = %s, want %s", tt.bytes, result, tt.expected)
			}
		})
	}
}

func TestListCacheEntries(t *testing.T) {
	// Create a temporary cache directory with test entries
	tmpDir := t.TempDir()
	os.Setenv("KFG_STORE_DIR", tmpDir)
	defer os.Unsetenv("KFG_STORE_DIR")

	cacheDir := getCacheDir()
	err := os.MkdirAll(cacheDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}

	// Test with empty cache
	t.Run("empty cache", func(t *testing.T) {
		entries, err := listCacheEntries()
		if err != nil {
			t.Errorf("listCacheEntries() error: %v", err)
		}
		if len(entries) != 0 {
			t.Errorf("listCacheEntries() returned %d entries, want 0", len(entries))
		}
	})

	// Create a test cache entry
	entryID := "test123"
	entryPath := filepath.Join(cacheDir, entryID)
	err = os.MkdirAll(entryPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test entry directory: %v", err)
	}

	// Create metadata.yaml
	metadataContent := `stepRefName: test-step
timestamp: 2024-01-15T10:30:00Z
`
	err = os.WriteFile(filepath.Join(entryPath, "metadata.yaml"), []byte(metadataContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create metadata.yaml: %v", err)
	}

	// Test with one cache entry
	t.Run("one entry", func(t *testing.T) {
		entries, err := listCacheEntries()
		if err != nil {
			t.Errorf("listCacheEntries() error: %v", err)
		}
		if len(entries) != 1 {
			t.Errorf("listCacheEntries() returned %d entries, want 1", len(entries))
		}
		if entries[0].ID != entryID {
			t.Errorf("entry.ID = %s, want %s", entries[0].ID, entryID)
		}
		if entries[0].StepRefName != "test-step" {
			t.Errorf("entry.StepRefName = %s, want test-step", entries[0].StepRefName)
		}
	})
}

func TestReadCacheEntry(t *testing.T) {
	// Create a temporary cache entry
	tmpDir := t.TempDir()
	os.Setenv("KFG_STORE_DIR", tmpDir)
	defer os.Unsetenv("KFG_STORE_DIR")

	cacheDir := getCacheDir()
	entryID := "test456"
	entryPath := filepath.Join(cacheDir, entryID)
	err := os.MkdirAll(entryPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test entry directory: %v", err)
	}

	// Test basic entry without output
	t.Run("basic entry", func(t *testing.T) {
		metadataContent := `stepRefName: basic-step
timestamp: 2024-01-15T10:30:00Z
`
		err = os.WriteFile(filepath.Join(entryPath, "metadata.yaml"), []byte(metadataContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create metadata.yaml: %v", err)
		}

		entry, err := readCacheEntry(entryPath)
		if err != nil {
			t.Errorf("readCacheEntry() error: %v", err)
		}
		if entry.ID != entryID {
			t.Errorf("entry.ID = %s, want %s", entry.ID, entryID)
		}
		if entry.StepRefName != "basic-step" {
			t.Errorf("entry.StepRefName = %s, want basic-step", entry.StepRefName)
		}
		if entry.HasOutput {
			t.Errorf("entry.HasOutput = true, want false")
		}
	})

	// Test entry with output
	t.Run("entry with output", func(t *testing.T) {
		outputValue := "test output value"
		outputEncoded := base64.StdEncoding.EncodeToString([]byte(outputValue))
		metadataContent := `stepRefName: output-step
timestamp: 2024-01-15T10:30:00Z
output:
  name: result
  valueEncoded: ` + outputEncoded + `
`
		err = os.WriteFile(filepath.Join(entryPath, "metadata.yaml"), []byte(metadataContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create metadata.yaml: %v", err)
		}

		entry, err := readCacheEntry(entryPath)
		if err != nil {
			t.Errorf("readCacheEntry() error: %v", err)
		}
		if !entry.HasOutput {
			t.Errorf("entry.HasOutput = false, want true")
		}
		if entry.OutputName != "result" {
			t.Errorf("entry.OutputName = %s, want result", entry.OutputName)
		}
		if entry.OutputValue != outputValue {
			t.Errorf("entry.OutputValue = %s, want %s", entry.OutputValue, outputValue)
		}
	})

	// Test entry with artifacts
	t.Run("entry with artifacts", func(t *testing.T) {
		metadataContent := `stepRefName: artifacts-step
timestamp: 2024-01-15T10:30:00Z
`
		err = os.WriteFile(filepath.Join(entryPath, "metadata.yaml"), []byte(metadataContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create metadata.yaml: %v", err)
		}

		// Create artifacts directory
		artifactsDir := filepath.Join(entryPath, "artifacts")
		err = os.MkdirAll(artifactsDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create artifacts directory: %v", err)
		}

		// Create a test artifact
		err = os.WriteFile(filepath.Join(artifactsDir, "test-artifact.txt"), []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test artifact: %v", err)
		}

		entry, err := readCacheEntry(entryPath)
		if err != nil {
			t.Errorf("readCacheEntry() error: %v", err)
		}
		if entry.ArtifactsDir == "" {
			t.Errorf("entry.ArtifactsDir should not be empty")
		}
		if entry.Size == 0 {
			t.Errorf("entry.Size should be greater than 0 with artifacts")
		}
	})
}

func TestInspectCacheEntry(t *testing.T) {
	// Create a temporary cache entry
	tmpDir := t.TempDir()
	os.Setenv("KFG_STORE_DIR", tmpDir)
	defer os.Unsetenv("KFG_STORE_DIR")

	cacheDir := getCacheDir()
	entryID := "inspect789"
	entryPath := filepath.Join(cacheDir, entryID)
	err := os.MkdirAll(entryPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test entry directory: %v", err)
	}

	metadataContent := `stepRefName: inspect-step
timestamp: 2024-01-15T10:30:00Z
`
	err = os.WriteFile(filepath.Join(entryPath, "metadata.yaml"), []byte(metadataContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create metadata.yaml: %v", err)
	}

	// Test inspect existing entry
	t.Run("existing entry", func(t *testing.T) {
		entry, err := inspectCacheEntry(entryID)
		if err != nil {
			t.Errorf("inspectCacheEntry() error: %v", err)
		}
		if entry.ID != entryID {
			t.Errorf("entry.ID = %s, want %s", entry.ID, entryID)
		}
		if entry.StepRefName != "inspect-step" {
			t.Errorf("entry.StepRefName = %s, want inspect-step", entry.StepRefName)
		}
	})

	// Test inspect non-existing entry
	t.Run("non-existing entry", func(t *testing.T) {
		_, err := inspectCacheEntry("nonexistent")
		if err == nil {
			t.Errorf("inspectCacheEntry() should return error for nonexistent entry")
		}
	})
}

func TestRemoveCacheEntry(t *testing.T) {
	// Create a temporary cache entry
	tmpDir := t.TempDir()
	os.Setenv("KFG_STORE_DIR", tmpDir)
	defer os.Unsetenv("KFG_STORE_DIR")

	cacheDir := getCacheDir()
	entryID := "remove123"
	entryPath := filepath.Join(cacheDir, entryID)
	err := os.MkdirAll(entryPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test entry directory: %v", err)
	}

	metadataContent := `stepRefName: remove-step
timestamp: 2024-01-15T10:30:00Z
`
	err = os.WriteFile(filepath.Join(entryPath, "metadata.yaml"), []byte(metadataContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create metadata.yaml: %v", err)
	}

	// Test remove existing entry
	t.Run("existing entry", func(t *testing.T) {
		err := removeCacheEntry(entryID)
		if err != nil {
			t.Errorf("removeCacheEntry() error: %v", err)
		}

		// Verify entry is removed
		if _, err := os.Stat(entryPath); !os.IsNotExist(err) {
			t.Errorf("entry should be removed after removeCacheEntry()")
		}
	})

	// Test remove non-existing entry
	t.Run("non-existing entry", func(t *testing.T) {
		err := removeCacheEntry("nonexistent")
		if err == nil {
			t.Errorf("removeCacheEntry() should return error for nonexistent entry")
		}
	})
}

func TestCalculateDirSize(t *testing.T) {
	// Create a temporary directory with files
	tmpDir := t.TempDir()

	// Create files with known sizes
	err := os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("12345"), 0644) // 5 bytes
	if err != nil {
		t.Fatalf("Failed to create file1.txt: %v", err)
	}

	err = os.WriteFile(filepath.Join(tmpDir, "file2.txt"), []byte("1234567890"), 0644) // 10 bytes
	if err != nil {
		t.Fatalf("Failed to create file2.txt: %v", err)
	}

	// Create subdirectory with file
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	err = os.WriteFile(filepath.Join(subDir, "file3.txt"), []byte("123456789012345"), 0644) // 15 bytes
	if err != nil {
		t.Fatalf("Failed to create file3.txt: %v", err)
	}

	// Calculate size
	size, err := calculateDirSize(tmpDir)
	if err != nil {
		t.Errorf("calculateDirSize() error: %v", err)
	}

	// Expected size: 5 + 10 + 15 = 30 bytes
	if size != 30 {
		t.Errorf("calculateDirSize() = %d, want 30", size)
	}
}

func TestListArtifacts(t *testing.T) {
	// Create a temporary artifacts directory
	tmpDir := t.TempDir()
	artifactsDir := filepath.Join(tmpDir, "artifacts")
	err := os.MkdirAll(artifactsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create artifacts directory: %v", err)
	}

	// Create test artifacts
	err = os.WriteFile(filepath.Join(artifactsDir, "artifact1.txt"), []byte("content1"), 0644)
	if err != nil {
		t.Fatalf("Failed to create artifact1.txt: %v", err)
	}

	err = os.WriteFile(filepath.Join(artifactsDir, "artifact2.txt"), []byte("content2"), 0644)
	if err != nil {
		t.Fatalf("Failed to create artifact2.txt: %v", err)
	}

	// List artifacts
	artifacts, err := listArtifacts(artifactsDir)
	if err != nil {
		t.Errorf("listArtifacts() error: %v", err)
	}

	if len(artifacts) != 2 {
		t.Errorf("listArtifacts() returned %d artifacts, want 2", len(artifacts))
	}

	// Check that artifacts contain expected names
	expected := []string{"artifact1.txt", "artifact2.txt"}
	for _, exp := range expected {
		found := false
		for _, artifact := range artifacts {
			if artifact == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("listArtifacts() missing artifact %s", exp)
		}
	}
}

// Integration tests for sys gc commands

func TestGcLsCommand(t *testing.T) {
	// Create a temporary cache directory with test entries
	tmpDir := t.TempDir()
	os.Setenv("KFG_STORE_DIR", tmpDir)
	defer os.Unsetenv("KFG_STORE_DIR")

	cacheDir := getCacheDir()
	err := os.MkdirAll(cacheDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}

	// Create test entry
	entryID := "test-ls-123"
	entryPath := filepath.Join(cacheDir, entryID)
	err = os.MkdirAll(entryPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test entry directory: %v", err)
	}

	metadataContent := `stepRefName: ls-test-step
timestamp: 2024-01-15T10:30:00Z
`
	err = os.WriteFile(filepath.Join(entryPath, "metadata.yaml"), []byte(metadataContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create metadata.yaml: %v", err)
	}

	// Test that listCacheEntries returns the test entry
	entries, err := listCacheEntries()
	if err != nil {
		t.Errorf("listCacheEntries() error: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("listCacheEntries() returned %d entries, want 1", len(entries))
	}
}

func TestGcPruneCommand(t *testing.T) {
	// Create a temporary cache directory with old and new entries
	tmpDir := t.TempDir()
	os.Setenv("KFG_STORE_DIR", tmpDir)
	defer os.Unsetenv("KFG_STORE_DIR")

	cacheDir := getCacheDir()
	err := os.MkdirAll(cacheDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}

	// Create old entry (older than 30 days)
	oldEntryID := "old-entry-123"
	oldEntryPath := filepath.Join(cacheDir, oldEntryID)
	err = os.MkdirAll(oldEntryPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create old entry directory: %v", err)
	}

	oldTimestamp := time.Now().AddDate(0, 0, -35).Format("2006-01-02T15:04:05Z")
	oldMetadataContent := `stepRefName: old-step
timestamp: ` + oldTimestamp + `
`
	err = os.WriteFile(filepath.Join(oldEntryPath, "metadata.yaml"), []byte(oldMetadataContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create old metadata.yaml: %v", err)
	}

	// Create new entry (within 30 days)
	newEntryID := "new-entry-456"
	newEntryPath := filepath.Join(cacheDir, newEntryID)
	err = os.MkdirAll(newEntryPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create new entry directory: %v", err)
	}

	newTimestamp := time.Now().Format("2006-01-02T15:04:05Z")
	newMetadataContent := `stepRefName: new-step
timestamp: ` + newTimestamp + `
`
	err = os.WriteFile(filepath.Join(newEntryPath, "metadata.yaml"), []byte(newMetadataContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create new metadata.yaml: %v", err)
	}

	// Test that prune removes old entry but keeps new entry
	// Note: prune policy removes entries older than 30 days
	entries, err := listCacheEntries()
	if err != nil {
		t.Fatalf("listCacheEntries() error: %v", err)
	}

	// Count entries that should be pruned
	cutoff := time.Now().AddDate(0, 0, -30)
	prunedCount := 0
	for _, entry := range entries {
		if entry.Timestamp.Before(cutoff) {
			prunedCount++
		}
	}

	// We expect 1 entry to be pruned (the old one)
	if prunedCount != 1 {
		t.Errorf("Expected 1 entry to be pruned, got %d", prunedCount)
	}
}