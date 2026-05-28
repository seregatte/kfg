package cache

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// SnapshotDirectory walks a directory and returns normalized relative paths.
// Paths are sorted deterministically for consistent comparison.
func SnapshotDirectory(rootPath string) ([]string, error) {
	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	var paths []string

	err = filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == absRoot {
			return nil
		}

		relPath, err := filepath.Rel(absRoot, path)
		if err != nil {
			return fmt.Errorf("failed to compute relative path: %w", err)
		}

		relPath = filepath.ToSlash(relPath)
		paths = append(paths, relPath)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	sort.Strings(paths)
	return paths, nil
}

// DiffSnapshots returns paths present in after but absent from before.
// Both inputs must be sorted.
func DiffSnapshots(beforePaths, afterPaths []string) []string {
	beforeSet := make(map[string]bool, len(beforePaths))
	for _, path := range beforePaths {
		beforeSet[path] = true
	}

	var newPaths []string
	for _, path := range afterPaths {
		if !beforeSet[path] {
			newPaths = append(newPaths, path)
		}
	}

	return newPaths
}

// DetectNewFiles compares a directory before and after an operation
// and returns the list of newly created files (relative paths).
func DetectNewFiles(workdir string, beforePaths []string) ([]string, error) {
	afterPaths, err := SnapshotDirectory(workdir)
	if err != nil {
		return nil, fmt.Errorf("failed to snapshot after state: %w", err)
	}

	newPaths := DiffSnapshots(beforePaths, afterPaths)
	return newPaths, nil
}

// FilterExistingPaths filters paths to only those that exist on disk.
func FilterExistingPaths(root string, paths []string) []string {
	var existing []string
	for _, path := range paths {
		fullPath := filepath.Join(root, path)
		if _, err := os.Stat(fullPath); err == nil {
			existing = append(existing, path)
		}
	}
	return existing
}

// NormalizePath normalizes a path to use forward slashes.
func NormalizePath(path string) string {
	return filepath.ToSlash(path)
}

// JoinPaths joins root and relative path, normalizing separators.
func JoinPaths(root string, relPath string) string {
	// Remove trailing separator from root and leading from relPath
	root = strings.TrimRight(root, "/\\")
	relPath = strings.TrimLeft(relPath, "/\\")
	return root + "/" + relPath
}

// ReadFsSnapshot reads a filesystem snapshot file (one path per line).
// Returns sorted list of paths.
func ReadFsSnapshot(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshot file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var paths []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			paths = append(paths, line)
		}
	}

	sort.Strings(paths)
	return paths, nil
}
