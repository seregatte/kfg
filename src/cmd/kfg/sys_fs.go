package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/seregatte/kfg/src/internal/logger"
	"github.com/spf13/cobra"
)

// fsCmd represents the fs command group for filesystem operations
var fsCmd = &cobra.Command{
	Use:   "fs",
	Short: "Internal filesystem commands",
	Long: `Internal filesystem commands for snapshot and diff operations.

These commands are intended for internal runtime use and provide
portable filesystem inspection capabilities.

Subcommands:
  snapshot  List normalized relative paths under a directory
  diff      Report paths newly present in an after snapshot

Examples:
  kfg sys fs snapshot /path/to/dir --maxdepth 1
  kfg sys fs diff --before before.txt --after after.txt`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// fsSnapshotCmd represents the snapshot command
var fsSnapshotCmd = &cobra.Command{
	Use:   "snapshot <path> [--maxdepth N]",
	Short: "List normalized relative paths under a directory",
	Long: `Print normalized relative paths rooted at the given directory.

Paths are:
  - Normalized to use forward slashes
  - Sorted deterministically
  - Relative to the given root path

The --maxdepth flag controls traversal depth:
  - 0: unlimited depth (traverse full subtree)
  - N > 0: limit to N levels below the root
  - Negative values are rejected with an error

Examples:
  kfg sys fs snapshot /path/to/dir
  kfg sys fs snapshot /path/to/dir --maxdepth 0  # unlimited
  kfg sys fs snapshot /path/to/dir --maxdepth 1  # immediate children`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rootPath := args[0]
		maxDepth, _ := cmd.Flags().GetInt("maxdepth")

		// Validate maxdepth
		if maxDepth < 0 {
			logger.Error("sys:fs:snapshot", "maxdepth must be >= 0")
			os.Exit(1)
		}

		// Verify path exists
		if _, err := os.Stat(rootPath); os.IsNotExist(err) {
			logger.Error("sys:fs:snapshot", fmt.Sprintf("path does not exist: %s", rootPath))
			os.Exit(1)
		}

		// Snapshot the directory
		paths, err := snapshotDirectory(rootPath, maxDepth)
		if err != nil {
			logger.Error("sys:fs:snapshot", err.Error())
			os.Exit(1)
		}

		// Print paths (one per line, sorted)
		for _, path := range paths {
			fmt.Println(path)
		}
	},
}

// fsDiffCmd represents the diff command
var fsDiffCmd = &cobra.Command{
	Use:   "diff --before <snapshot> --after <snapshot>",
	Short: "Report paths newly present in the after snapshot",
	Long: `Print only paths that are present in the after snapshot
but absent from the before snapshot.

Paths are read from snapshot files (one path per line).
Output preserves deterministic ordering.

Examples:
  kfg sys fs diff --before before.txt --after after.txt`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		beforeFile, _ := cmd.Flags().GetString("before")
		afterFile, _ := cmd.Flags().GetString("after")

		// Read before snapshot
		beforePaths, err := readSnapshotFile(beforeFile)
		if err != nil {
			logger.Error("sys:fs:diff", fmt.Sprintf("failed to read before snapshot: %v", err))
			os.Exit(1)
		}

		// Read after snapshot
		afterPaths, err := readSnapshotFile(afterFile)
		if err != nil {
			logger.Error("sys:fs:diff", fmt.Sprintf("failed to read after snapshot: %v", err))
			os.Exit(1)
		}

		// Compute diff (paths in after but not in before)
		newPaths := diffSnapshots(beforePaths, afterPaths)

		// Print new paths (one per line, sorted)
		for _, path := range newPaths {
			fmt.Println(path)
		}
	},
}

func init() {
	// Add fsCmd to sysCmd
	sysCmd.AddCommand(fsCmd)

	// Add subcommands to fsCmd
	fsCmd.AddCommand(fsSnapshotCmd)
	fsCmd.AddCommand(fsDiffCmd)

	// Add flags for snapshot
	fsSnapshotCmd.Flags().IntP("maxdepth", "d", 0, "Maximum traversal depth (0 = unlimited)")

	// Add flags for diff
	fsDiffCmd.Flags().StringP("before", "b", "", "Before snapshot file (required")
	fsDiffCmd.Flags().StringP("after", "a", "", "After snapshot file (required)")
	fsDiffCmd.MarkFlagRequired("before")
	fsDiffCmd.MarkFlagRequired("after")
}

// snapshotDirectory walks a directory and returns normalized relative paths
func snapshotDirectory(rootPath string, maxDepth int) ([]string, error) {
	var paths []string

	// Normalize root path
	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Walk the directory
	err = filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == absRoot {
			return nil
		}

		// Compute relative path
		relPath, err := filepath.Rel(absRoot, path)
		if err != nil {
			return fmt.Errorf("failed to compute relative path: %w", err)
		}

		// Normalize to forward slashes (for cross-platform consistency)
		relPath = filepath.ToSlash(relPath)

		// Check depth constraint
		if maxDepth > 0 {
			depth := pathDepth(relPath)
			if depth > maxDepth {
				// Skip this entry and its children if it's a directory
				if d.IsDir() {
					return fs.SkipDir
				}
				return nil
			}
		}

		// Add path to results
		paths = append(paths, relPath)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	// Sort paths deterministically
	sort.Strings(paths)

	return paths, nil
}

// pathDepth computes the depth of a relative path
func pathDepth(relPath string) int {
	// Count the number of path separators
	return len(strings.Split(relPath, "/"))
}

// readSnapshotFile reads a snapshot file (one path per line)
func readSnapshotFile(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Split into lines and filter empty lines
	lines := strings.Split(string(data), "\n")
	var paths []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			paths = append(paths, line)
		}
	}

	// Sort for consistent comparison
	sort.Strings(paths)

	return paths, nil
}

// diffSnapshots returns paths present in after but absent from before
func diffSnapshots(beforePaths, afterPaths []string) []string {
	// Build a set of before paths for fast lookup
	beforeSet := make(map[string]bool)
	for _, path := range beforePaths {
		beforeSet[path] = true
	}

	// Find paths in after that are not in before
	var newPaths []string
	for _, path := range afterPaths {
		if !beforeSet[path] {
			newPaths = append(newPaths, path)
		}
	}

	// Already sorted since afterPaths was sorted
	return newPaths
}