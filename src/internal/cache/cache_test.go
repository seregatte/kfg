package cache

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestComputeIdentityDeterministic(t *testing.T) {
	inputs := []string{
		"test-step",
		"ctx7.steps.install",
		"openspec.steps.install",
		"a",
		"",
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			result := ComputeIdentity(input)
			if len(result) != 64 { // SHA256 hex is 64 chars
				t.Errorf("ComputeIdentity(%q) returned %d chars, want 64", input, len(result))
			}
			// Same input should produce same output (deterministic)
			result2 := ComputeIdentity(input)
			if result != result2 {
				t.Errorf("ComputeIdentity(%q) not deterministic: %s != %s", input, result, result2)
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

func TestComputeIdentityVersionInvalidation(t *testing.T) {
	name := "openspec.steps.install"

	// Current version (v2) produces a known hash
	currentHash := ComputeIdentity(name)
	if len(currentHash) != 64 {
		t.Fatalf("ComputeIdentity(%q) returned %d chars, want 64", name, len(currentHash))
	}

	// Manually compute what the old v1 identity would have been.
	// v1 used raw SHA256(name) without a version prefix.
	// If the old identity ever collides with the new one, the engine
	// would silently restore stale artifacts from a legacy entry.
	oldKey := name
	oldHash := sha256Hex(oldKey)
	if currentHash == oldHash {
		t.Errorf("v2 identity %s must not collide with v1 identity %s", currentHash, oldHash)
	}

	// Manually compute what a zero-version identity would have been.
	zeroKey := "\x00" + name
	zeroHash := sha256Hex(zeroKey)
	if currentHash == zeroHash {
		t.Errorf("v2 identity %s must not collide with zero-version identity %s", currentHash, zeroHash)
	}
}

func TestComputeIdentityDifferentVersions(t *testing.T) {
	// Verify that changing the version constant changes all identities.
	// We test this by comparing the current ComputeIdentity output with
	// a manual SHA256("v1\x00"+name) — if the constant is bumped to v3,
	// this test will need updating, which is the intended safety net.
	name := "test.step"
	v2Hash := ComputeIdentity(name)

	v1Key := "v1\x00" + name
	v1Hash := sha256Hex(v1Key)

	if v2Hash == v1Hash {
		t.Errorf("v2 identity must differ from v1 identity for %q", name)
	}
}

// sha256Hex returns the hex-encoded SHA256 hash of data.
// Used only in tests to verify version-invalidation properties.
func sha256Hex(data string) string {
	h := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", h)
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

func TestSnapshotDirectoryLeafOnly(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test structure with files and directories
	if err := os.MkdirAll(filepath.Join(tmpDir, "dir1"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tmpDir, "dir1", "dir2"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "dir1", "file2.txt"), []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "dir1", "dir2", "file3.txt"), []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	paths, err := SnapshotDirectory(tmpDir)
	if err != nil {
		t.Fatalf("SnapshotDirectory() error: %v", err)
	}

	// Directories (dir1, dir1/dir2) must NOT appear in the snapshot.
	expected := []string{"dir1/dir2/file3.txt", "dir1/file2.txt", "file1.txt"}
	if len(paths) != len(expected) {
		t.Errorf("SnapshotDirectory() returned %d paths, want %d", len(paths), len(expected))
		for _, p := range paths {
			t.Logf("  got: %s", p)
		}
	}

	for i, exp := range expected {
		if i >= len(paths) || paths[i] != exp {
			t.Errorf("paths[%d] = %s, want %s", i, paths[i], exp)
		}
	}
}

func TestSnapshotDirectoryEmpty(t *testing.T) {
	tmpDir := t.TempDir()

	paths, err := SnapshotDirectory(tmpDir)
	if err != nil {
		t.Fatalf("SnapshotDirectory() error: %v", err)
	}

	if len(paths) != 0 {
		t.Errorf("SnapshotDirectory() returned %d paths on empty dir, want 0", len(paths))
	}
}

func TestSnapshotDirectorySymlinks(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file and a symlink to it
	if err := os.WriteFile(filepath.Join(tmpDir, "real.txt"), []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("real.txt", filepath.Join(tmpDir, "link.txt")); err != nil {
		t.Fatal(err)
	}
	// Create a directory symlink (should be excluded)
	if err := os.MkdirAll(filepath.Join(tmpDir, "realdir"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("realdir", filepath.Join(tmpDir, "linkdir")); err != nil {
		t.Fatal(err)
	}

	paths, err := SnapshotDirectory(tmpDir)
	if err != nil {
		t.Fatalf("SnapshotDirectory() error: %v", err)
	}

	// link.txt (symlink to file) should appear; linkdir (symlink to dir) should not
	expected := []string{"link.txt", "real.txt"}
	if len(paths) != len(expected) {
		t.Errorf("SnapshotDirectory() returned %d paths, want %d", len(paths), len(expected))
		for _, p := range paths {
			t.Logf("  got: %s", p)
		}
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

func TestDiffSnapshotsDirectoryExclusion(t *testing.T) {
	// Verify that directory paths never appear in snapshot diffs,
	// even if they exist in the directory tree. This is the core
	// property that prevents openspec/ from being treated as an artifact.
	tmpDir := t.TempDir()

	// Initial state: empty directory
	before, err := SnapshotDirectory(tmpDir)
	if err != nil {
		t.Fatalf("SnapshotDirectory() error: %v", err)
	}
	if len(before) != 0 {
		t.Fatalf("Expected empty snapshot, got %d paths", len(before))
	}

	// Step creates a directory and files
	if err := os.MkdirAll(filepath.Join(tmpDir, "openspec", "changes", "foo"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "openspec", "changes", "foo", "proposal.md"), []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	after, err := SnapshotDirectory(tmpDir)
	if err != nil {
		t.Fatalf("SnapshotDirectory() error: %v", err)
	}

	// Snapshot should only contain files, not directories
	for _, p := range after {
		if p == "openspec" || p == "openspec/changes" || p == "openspec/changes/foo" {
			t.Errorf("Directory %q should not appear in leaf-only snapshot", p)
		}
	}

	// Diff should only show new files, not directories
	delta := DiffSnapshots(before, after)
	for _, p := range delta {
		if p == "openspec" || p == "openspec/changes" || p == "openspec/changes/foo" {
			t.Errorf("Directory %q should not appear in diff delta", p)
		}
	}

	// Only files should be in the delta
	expectedFiles := map[string]bool{
		"file.txt":                              false,
		"openspec/changes/foo/proposal.md":      false,
	}
	if len(delta) != len(expectedFiles) {
		t.Errorf("Delta has %d items, want %d", len(delta), len(expectedFiles))
		for _, p := range delta {
			t.Logf("  delta: %s", p)
		}
	}
	for _, p := range delta {
		if _, ok := expectedFiles[p]; !ok {
			t.Errorf("Unexpected file in delta: %s", p)
		}
	}
}
