package cache

import (
	"os"
	"path/filepath"
	"testing"
)

func TestComputeIdentity(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple name",
			input:    "test-step",
			expected: "a1b2c3d4e5f6", // Will be computed
		},
		{
			name:     "dotted name",
			input:    "ctx7.steps.install",
			expected: "", // Will be computed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ComputeIdentity(tt.input)
			if len(result) != 64 { // SHA256 hex is 64 chars
				t.Errorf("ComputeIdentity(%q) returned %d chars, want 64", tt.input, len(result))
			}
			// Same input should produce same output
			result2 := ComputeIdentity(tt.input)
			if result != result2 {
				t.Errorf("ComputeIdentity(%q) not deterministic: %s != %s", tt.input, result, result2)
			}
		})
	}

	// Different inputs should produce different outputs
	hash1 := ComputeIdentity("step1")
	hash2 := ComputeIdentity("step2")
	if hash1 == hash2 {
		t.Error("Different inputs should produce different hashes")
	}
}

func TestGetCacheDir(t *testing.T) {
	t.Run("with KFG_STORE_DIR", func(t *testing.T) {
		tmpDir := t.TempDir()
		os.Setenv("KFG_STORE_DIR", tmpDir)
		defer os.Unsetenv("KFG_STORE_DIR")

		result := GetCacheDir()
		expected := filepath.Join(tmpDir, "cache")
		if result != expected {
			t.Errorf("GetCacheDir() = %s, want %s", result, expected)
		}
	})

	t.Run("without KFG_STORE_DIR", func(t *testing.T) {
		os.Unsetenv("KFG_STORE_DIR")
		result := GetCacheDir()
		if !filepath.IsAbs(result) {
			t.Errorf("GetCacheDir() should return absolute path, got %s", result)
		}
	})
}

func TestExists(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("KFG_STORE_DIR", tmpDir)
	defer os.Unsetenv("KFG_STORE_DIR")

	t.Run("nonexistent entry", func(t *testing.T) {
		if Exists("nonexistent-step") {
			t.Error("Exists() should return false for nonexistent entry")
		}
	})

	t.Run("existing entry", func(t *testing.T) {
		// Create a cache entry
		entryPath := GetEntryPath("test-step")
		os.MkdirAll(entryPath, 0755)
		os.WriteFile(filepath.Join(entryPath, "metadata.yaml"), []byte("stepRefName: test-step\n"), 0644)

		if !Exists("test-step") {
			t.Error("Exists() should return true for existing entry")
		}
	})
}

func TestMetadataReadWrite(t *testing.T) {
	tmpDir := t.TempDir()

	metadata := &CacheMetadata{
		StepRefName: "test-step",
		Timestamp:   "2024-01-15T10:30:00Z",
		Artifacts:   []string{"file1.txt", "dir/file2.txt"},
		Output: &OutputMetadata{
			Name:         "result",
			ValueEncoded: "dGVzdA==", // base64("test")
		},
	}

	// Write
	if err := WriteMetadata(tmpDir, metadata); err != nil {
		t.Fatalf("WriteMetadata() error: %v", err)
	}

	// Read
	read, err := ReadMetadata(tmpDir)
	if err != nil {
		t.Fatalf("ReadMetadata() error: %v", err)
	}

	if read.StepRefName != metadata.StepRefName {
		t.Errorf("StepRefName = %s, want %s", read.StepRefName, metadata.StepRefName)
	}
	if len(read.Artifacts) != len(metadata.Artifacts) {
		t.Errorf("Artifacts count = %d, want %d", len(read.Artifacts), len(metadata.Artifacts))
	}
	if read.Output == nil {
		t.Fatal("Output should not be nil")
	}
	if read.Output.Name != "result" {
		t.Errorf("Output.Name = %s, want result", read.Output.Name)
	}
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
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := FormatSize(tt.bytes)
			if result != tt.expected {
				t.Errorf("FormatSize(%d) = %s, want %s", tt.bytes, result, tt.expected)
			}
		})
	}
}

func TestSnapshotDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test structure
	os.MkdirAll(filepath.Join(tmpDir, "dir1"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "dir1", "file2.txt"), []byte("content"), 0644)

	paths, err := SnapshotDirectory(tmpDir)
	if err != nil {
		t.Fatalf("SnapshotDirectory() error: %v", err)
	}

	expected := []string{"dir1", "dir1/file2.txt", "file1.txt"}
	if len(paths) != len(expected) {
		t.Errorf("SnapshotDirectory() returned %d paths, want %d", len(paths), len(expected))
	}

	for i, exp := range expected {
		if i >= len(paths) || paths[i] != exp {
			t.Errorf("paths[%d] = %s, want %s", i, paths[i], exp)
		}
	}
}

func TestDiffSnapshots(t *testing.T) {
	before := []string{"a.txt", "b.txt"}
	after := []string{"a.txt", "b.txt", "c.txt"}

	result := DiffSnapshots(before, after)
	if len(result) != 1 || result[0] != "c.txt" {
		t.Errorf("DiffSnapshots() = %v, want [c.txt]", result)
	}
}

func TestComputeDelta(t *testing.T) {
	before := []string{"a.txt", "b.txt"}
	after := []string{"a.txt", "b.txt", "c.txt", "d.txt"}

	delta := computeDelta(before, after)
	if len(delta) != 2 {
		t.Errorf("computeDelta() returned %d items, want 2", len(delta))
	}
}

func TestMergeArtifacts(t *testing.T) {
	a := []string{"a.txt", "b.txt"}
	b := []string{"b.txt", "c.txt"}

	merged := mergeArtifacts(a, b)
	if len(merged) != 3 {
		t.Errorf("mergeArtifacts() returned %d items, want 3", len(merged))
	}
}
