package main

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestSnapshotDirectory(t *testing.T) {
	// Create a temporary directory structure for testing
	tmpDir := t.TempDir()

	// Create test structure:
	// tmpDir/
	//   file1.txt
	//   dir1/
	//     file2.txt
	//     subdir/
	//       file3.txt
	//   dir2/
	//     file4.txt

	dir1 := filepath.Join(tmpDir, "dir1")
	subdir := filepath.Join(dir1, "subdir")
	dir2 := filepath.Join(tmpDir, "dir2")

	os.MkdirAll(dir1, 0755)
	os.MkdirAll(subdir, 0755)
	os.MkdirAll(dir2, 0755)

	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(dir1, "file2.txt"), []byte("content2"), 0644)
	os.WriteFile(filepath.Join(subdir, "file3.txt"), []byte("content3"), 0644)
	os.WriteFile(filepath.Join(dir2, "file4.txt"), []byte("content4"), 0644)

	// Test unlimited depth (maxdepth 0)
	t.Run("unlimited depth", func(t *testing.T) {
		paths, err := snapshotDirectory(tmpDir, 0)
		if err != nil {
			t.Fatalf("snapshotDirectory() error: %v", err)
		}

		// Should contain all paths
		expectedPaths := []string{
			"dir1",
			"dir1/file2.txt",
			"dir1/subdir",
			"dir1/subdir/file3.txt",
			"dir2",
			"dir2/file4.txt",
			"file1.txt",
		}

		// Sort both for comparison
		sort.Strings(expectedPaths)
		sort.Strings(paths)

		if len(paths) != len(expectedPaths) {
			t.Errorf("snapshotDirectory() returned %d paths, want %d", len(paths), len(expectedPaths))
		}

		for i, expected := range expectedPaths {
			if i >= len(paths) || paths[i] != expected {
				t.Errorf("paths[%d] = %s, want %s", i, paths[i], expected)
			}
		}
	})

	// Test depth 1 (immediate children)
	t.Run("depth 1", func(t *testing.T) {
		paths, err := snapshotDirectory(tmpDir, 1)
		if err != nil {
			t.Fatalf("snapshotDirectory() error: %v", err)
		}

		// Should only contain immediate children (depth 1)
		expectedPaths := []string{
			"dir1",
			"dir2",
			"file1.txt",
		}

		// Sort both for comparison
		sort.Strings(expectedPaths)
		sort.Strings(paths)

		if len(paths) != len(expectedPaths) {
			t.Errorf("snapshotDirectory() returned %d paths, want %d", len(paths), len(expectedPaths))
			t.Errorf("paths: %v", paths)
			t.Errorf("expected: %v", expectedPaths)
		}

		for i, expected := range expectedPaths {
			if i >= len(paths) || paths[i] != expected {
				t.Errorf("paths[%d] = %s, want %s", i, paths[i], expected)
			}
		}
	})

	// Test depth 2
	t.Run("depth 2", func(t *testing.T) {
		paths, err := snapshotDirectory(tmpDir, 2)
		if err != nil {
			t.Fatalf("snapshotDirectory() error: %v", err)
		}

		// Should contain depth 1 and 2 (but not depth 3)
		expectedPaths := []string{
			"dir1",
			"dir1/file2.txt",
			"dir1/subdir",
			"dir2",
			"dir2/file4.txt",
			"file1.txt",
		}

		// Sort both for comparison
		sort.Strings(expectedPaths)
		sort.Strings(paths)

		if len(paths) != len(expectedPaths) {
			t.Errorf("snapshotDirectory() returned %d paths, want %d", len(paths), len(expectedPaths))
			t.Errorf("paths: %v", paths)
			t.Errorf("expected: %v", expectedPaths)
		}

		for i, expected := range expectedPaths {
			if i >= len(paths) || paths[i] != expected {
				t.Errorf("paths[%d] = %s, want %s", i, paths[i], expected)
			}
		}
	})

	// Test empty directory
	t.Run("empty directory", func(t *testing.T) {
		emptyDir := filepath.Join(tmpDir, "empty")
		os.MkdirAll(emptyDir, 0755)

		paths, err := snapshotDirectory(emptyDir, 0)
		if err != nil {
			t.Fatalf("snapshotDirectory() error: %v", err)
		}

		// Empty directory should return no paths
		if len(paths) != 0 {
			t.Errorf("snapshotDirectory() returned %d paths for empty directory, want 0", len(paths))
		}
	})
}

func TestPathNormalization(t *testing.T) {
	// Test that paths are normalized to forward slashes
	// This is important for cross-platform consistency

	// Create a temporary directory with nested structure
	tmpDir := t.TempDir()
	nestedDir := filepath.Join(tmpDir, "dir1", "subdir", "nested")
	os.MkdirAll(nestedDir, 0755)
	os.WriteFile(filepath.Join(nestedDir, "file.txt"), []byte("content"), 0644)

	paths, err := snapshotDirectory(tmpDir, 0)
	if err != nil {
		t.Fatalf("snapshotDirectory() error: %v", err)
	}

	// All paths should use forward slashes (not backslashes on Windows)
	for _, path := range paths {
		if filepath.Separator == '\\' {
			// On Windows, verify paths don't contain backslashes
			for _, c := range path {
				if c == '\\' {
					t.Errorf("path %s contains backslash, should use forward slash", path)
				}
			}
		}
	}
}

func TestPathDepth(t *testing.T) {
	tests := []struct {
		path     string
		expected int
	}{
		{"file.txt", 1},
		{"dir/file.txt", 2},
		{"dir/subdir/file.txt", 3},
		{"dir", 1},
		{"dir/subdir", 2},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := pathDepth(tt.path)
			if result != tt.expected {
				t.Errorf("pathDepth(%s) = %d, want %d", tt.path, result, tt.expected)
			}
		})
	}
}

func TestDiffSnapshots(t *testing.T) {
	tests := []struct {
		name     string
		before   []string
		after    []string
		expected []string
	}{
		{
			name:     "no change",
			before:   []string{"file1.txt", "file2.txt"},
			after:    []string{"file1.txt", "file2.txt"},
			expected: []string{},
		},
		{
			name:     "one new file",
			before:   []string{"file1.txt"},
			after:    []string{"file1.txt", "file2.txt"},
			expected: []string{"file2.txt"},
		},
		{
			name:     "multiple new files",
			before:   []string{"file1.txt"},
			after:    []string{"file1.txt", "file2.txt", "file3.txt"},
			expected: []string{"file2.txt", "file3.txt"},
		},
		{
			name:     "empty before",
			before:   []string{},
			after:    []string{"file1.txt", "file2.txt"},
			expected: []string{"file1.txt", "file2.txt"},
		},
		{
			name:     "empty after",
			before:   []string{"file1.txt", "file2.txt"},
			after:    []string{},
			expected: []string{},
		},
		{
			name:     "both empty",
			before:   []string{},
			after:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Sort inputs (as the function expects sorted input)
			sort.Strings(tt.before)
			sort.Strings(tt.after)

			result := diffSnapshots(tt.before, tt.after)

			// Sort expected and result for comparison
			sort.Strings(tt.expected)
			sort.Strings(result)

			if len(result) != len(tt.expected) {
				t.Errorf("diffSnapshots() returned %d paths, want %d", len(result), len(tt.expected))
				t.Errorf("result: %v", result)
				t.Errorf("expected: %v", tt.expected)
			}

			for i, expected := range tt.expected {
				if i >= len(result) || result[i] != expected {
					t.Errorf("result[%d] = %s, want %s", i, result[i], expected)
				}
			}
		})
	}
}

func TestReadSnapshotFile(t *testing.T) {
	// Create a temporary snapshot file
	tmpDir := t.TempDir()
	snapshotFile := filepath.Join(tmpDir, "snapshot.txt")

	// Test with normal content
	t.Run("normal content", func(t *testing.T) {
		content := "file1.txt\nfile2.txt\nfile3.txt\n"
		os.WriteFile(snapshotFile, []byte(content), 0644)

		paths, err := readSnapshotFile(snapshotFile)
		if err != nil {
			t.Fatalf("readSnapshotFile() error: %v", err)
		}

		expected := []string{"file1.txt", "file2.txt", "file3.txt"}
		sort.Strings(expected)

		if len(paths) != len(expected) {
			t.Errorf("readSnapshotFile() returned %d paths, want %d", len(paths), len(expected))
		}

		for i, exp := range expected {
			if i >= len(paths) || paths[i] != exp {
				t.Errorf("paths[%d] = %s, want %s", i, paths[i], exp)
			}
		}
	})

	// Test with empty lines and whitespace
	t.Run("empty lines and whitespace", func(t *testing.T) {
		content := "file1.txt\n  \n\nfile2.txt\n  file3.txt  \n"
		os.WriteFile(snapshotFile, []byte(content), 0644)

		paths, err := readSnapshotFile(snapshotFile)
		if err != nil {
			t.Fatalf("readSnapshotFile() error: %v", err)
		}

		// Should filter out empty lines and trim whitespace
		expected := []string{"file1.txt", "file2.txt", "file3.txt"}
		sort.Strings(expected)

		if len(paths) != len(expected) {
			t.Errorf("readSnapshotFile() returned %d paths, want %d", len(paths), len(expected))
		}
	})

	// Test with empty file
	t.Run("empty file", func(t *testing.T) {
		os.WriteFile(snapshotFile, []byte(""), 0644)

		paths, err := readSnapshotFile(snapshotFile)
		if err != nil {
			t.Fatalf("readSnapshotFile() error: %v", err)
		}

		if len(paths) != 0 {
			t.Errorf("readSnapshotFile() returned %d paths for empty file, want 0", len(paths))
		}
	})

	// Test with file that doesn't exist
	t.Run("file does not exist", func(t *testing.T) {
		nonExistent := filepath.Join(tmpDir, "does-not-exist.txt")
		_, err := readSnapshotFile(nonExistent)
		if err == nil {
			t.Errorf("readSnapshotFile() should return error for non-existent file")
		}
	})
}

func TestSnapshotDeterministicOrdering(t *testing.T) {
	// Test that snapshot returns paths in deterministic (sorted) order
	tmpDir := t.TempDir()

	// Create multiple files (not in alphabetical order)
	os.WriteFile(filepath.Join(tmpDir, "z-file.txt"), []byte("z"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "a-file.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "m-file.txt"), []byte("m"), 0644)

	paths, err := snapshotDirectory(tmpDir, 0)
	if err != nil {
		t.Fatalf("snapshotDirectory() error: %v", err)
	}

	// Verify paths are sorted
	expected := []string{"a-file.txt", "m-file.txt", "z-file.txt"}
	for i, exp := range expected {
		if i >= len(paths) || paths[i] != exp {
			t.Errorf("paths[%d] = %s, want %s", i, paths[i], exp)
		}
	}
}

func TestDiffPreservesOrdering(t *testing.T) {
	// Test that diff returns paths in sorted order
	before := []string{"a.txt", "z.txt"}
	after := []string{"a.txt", "m.txt", "z.txt"}

	result := diffSnapshots(before, after)

	// Result should be sorted
	expected := []string{"m.txt"}
	for i, exp := range expected {
		if i >= len(result) || result[i] != exp {
			t.Errorf("result[%d] = %s, want %s", i, result[i], exp)
		}
	}
}

// Note: Invalid depth validation (negative values) is handled at the CLI level,
// not in the snapshotDirectory function. The CLI command Run function checks
// maxdepth >= 0 and exits with error for negative values.
// This behavior is tested via integration tests or manual CLI invocation.