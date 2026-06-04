package cache

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// StoreInput represents the JSON input for the store subcommand.
type StoreInput struct {
	Before      []string     `json:"before"`
	After       []string     `json:"after"`
	Declarative []string     `json:"declarative"`
	Output      *StoreOutput `json:"output,omitempty"`
}

// StoreOutput represents the output to cache.
type StoreOutput struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// StoreResult represents the result of a store operation.
type StoreResult struct {
	StepRefName string
	EntryPath   string
	Artifacts   []string
	HasOutput   bool
}

// Store persists Step execution results to the cache.
// It reads JSON from stdin, computes artifact delta, performs fs diff,
// copies artifacts, writes metadata, and commits atomically.
func Store(stepRefName string, workdir string, input *StoreInput) (*StoreResult, error) {
	cacheDir := GetCacheDir()
	identity := ComputeIdentity(stepRefName)
	entryPath := filepath.Join(cacheDir, identity)
	tmpPath := entryPath + ".tmp"

	// Clean up any stale temp directory
	os.RemoveAll(tmpPath)

	// Create temp directory for atomic write
	if err := os.MkdirAll(filepath.Join(tmpPath, "artifacts"), 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Compute artifact delta (after - before)
	deltaArtifacts := computeDelta(input.Before, input.After)

	// Merge with declarative artifacts
	allArtifacts := mergeArtifacts(deltaArtifacts, input.Declarative)

	// Perform fs diff to detect unregistered artifacts
	fsNewFiles, err := detectFsNewFiles(workdir, input.Before)
	if err != nil {
		// Log but don't fail - fs diff is best-effort
		fmt.Fprintf(os.Stderr, "warning: fs diff failed: %v\n", err)
	}

	// Filter out cache directory from fs diff results
	fsNewFiles = filterOutCacheDir(workdir, cacheDir, fsNewFiles)

	// Merge fs-detected artifacts
	allArtifacts = mergeArtifacts(allArtifacts, fsNewFiles)

	// Filter to only existing files
	allArtifacts = FilterExistingPaths(workdir, allArtifacts)

	// Copy artifacts to cache
	for _, artifact := range allArtifacts {
		sourcePath := filepath.Join(workdir, artifact)
		if err := copyArtifactWithPath(sourcePath, artifact, filepath.Join(tmpPath, "artifacts")); err != nil {
			os.RemoveAll(tmpPath)
			return nil, fmt.Errorf("failed to copy artifact %s: %w", artifact, err)
		}
	}

	// Build metadata
	metadata := &CacheMetadata{
		StepRefName: stepRefName,
		Timestamp:   time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		Artifacts:   allArtifacts,
	}

	// Add output if present
	if input.Output != nil && input.Output.Name != "" {
		metadata.Output = &OutputMetadata{
			Name:         input.Output.Name,
			ValueEncoded: base64.StdEncoding.EncodeToString([]byte(input.Output.Value)),
		}
	}

	// Write metadata
	if err := WriteMetadata(tmpPath, metadata); err != nil {
		os.RemoveAll(tmpPath)
		return nil, fmt.Errorf("failed to write metadata: %w", err)
	}

	// Atomically rename temp to final
	if err := os.Rename(tmpPath, entryPath); err != nil {
		os.RemoveAll(tmpPath)
		return nil, fmt.Errorf("failed to commit cache entry: %w", err)
	}

	return &StoreResult{
		StepRefName: stepRefName,
		EntryPath:   entryPath,
		Artifacts:   allArtifacts,
		HasOutput:   input.Output != nil && input.Output.Name != "",
	}, nil
}

// ParseStoreInput parses JSON from a byte slice into StoreInput.
func ParseStoreInput(data []byte) (*StoreInput, error) {
	var input StoreInput
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, fmt.Errorf("failed to parse store input: %w", err)
	}
	return &input, nil
}

// computeDelta returns items in after that are not in before.
func computeDelta(before, after []string) []string {
	beforeSet := make(map[string]bool, len(before))
	for _, item := range before {
		beforeSet[item] = true
	}

	var delta []string
	for _, item := range after {
		if !beforeSet[item] {
			delta = append(delta, item)
		}
	}
	return delta
}

// mergeArtifacts merges two artifact lists, removing duplicates.
func mergeArtifacts(a, b []string) []string {
	seen := make(map[string]bool)
	var merged []string

	for _, item := range a {
		if !seen[item] {
			seen[item] = true
			merged = append(merged, item)
		}
	}

	for _, item := range b {
		if !seen[item] {
			seen[item] = true
			merged = append(merged, item)
		}
	}

	return merged
}

// detectFsNewFiles detects new files in workdir by comparing before/after snapshots.
func detectFsNewFiles(workdir string, beforePaths []string) ([]string, error) {
	afterPaths, err := SnapshotDirectory(workdir)
	if err != nil {
		return nil, err
	}

	return DiffSnapshots(beforePaths, afterPaths), nil
}

// filterOutCacheDir removes any paths that are within the cache directory.
// This prevents the cache directory itself from being cached as an artifact.
func filterOutCacheDir(workdir string, cacheDir string, paths []string) []string {
	// Get relative path of cache dir from workdir
	relCacheDir, err := filepath.Rel(workdir, cacheDir)
	if err != nil {
		return paths
	}

	// Also get parent directory of cache (e.g., .kfg-store from .kfg-store/cache)
	relCacheParent := filepath.Dir(relCacheDir)

	var filtered []string
	for _, path := range paths {
		// Skip if path is the cache dir, cache parent, or within them
		if path == relCacheDir || path == relCacheParent {
			continue
		}
		// Skip if path starts with cache dir or parent dir
		if len(path) > len(relCacheDir) && path[:len(relCacheDir)+1] == relCacheDir+"/" {
			continue
		}
		if len(path) > len(relCacheParent) && path[:len(relCacheParent)+1] == relCacheParent+"/" {
			continue
		}
		filtered = append(filtered, path)
	}
	return filtered
}
