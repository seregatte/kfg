package main

import (
	"encoding/base64"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/seregatte/kfg/src/internal/logger"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// CacheMetadata represents the metadata.yaml content in a cache entry
type CacheMetadata struct {
	StepRefName string    `yaml:"stepRefName"`
	Timestamp   string    `yaml:"timestamp"`
	Output      *struct {
		Name        string `yaml:"name"`
		ValueEncoded string `yaml:"valueEncoded"`
	} `yaml:"output,omitempty"`
	Artifacts   []string `yaml:"artifacts,omitempty"` // Artifact paths stored in metadata
}

// CacheEntry represents a cache entry with its metadata and path
type CacheEntry struct {
	ID           string        // The cache entry directory name (hash)
	Path         string        // Full path to the cache entry directory
	StepRefName  string        // Step reference name from metadata
	Timestamp    time.Time     // Parsed timestamp from metadata
	Size         int64         // Total size in bytes
	ArtifactsDir string        // Path to artifacts directory
	HasOutput    bool          // Whether the entry has an output
	OutputName   string        // Output name if present
	OutputValue  string        // Decoded output value if present
	Artifacts    []string      // Artifact paths from metadata
}

// gcCmd represents the gc command group for cache garbage collection
var gcCmd = &cobra.Command{
	Use:   "gc",
	Short: "Garbage collection commands for Step cache",
	Long: `Garbage collection commands for managing persisted Step cache entries.

These commands operate on cache entries stored under KFG_STORE_DIR/cache
(defaults to ~/.kfg/store/cache).

Subcommands:
  ls      List cache entries with metadata
  inspect Show detailed metadata for a cache entry
  rm      Remove specific cache entries
  prune   Remove old or unused cache entries
  du      Show disk usage of cache entries

Examples:
  kfg sys gc ls
  kfg sys gc inspect abc123
  kfg sys gc rm abc123
  kfg sys gc prune
  kfg sys gc du`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// gcLsCmd lists cache entries
var gcLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List cache entries with metadata",
	Long: `List all cache entries with their stable identifiers and metadata.

Each cache entry is displayed with:
  - ID: The stable identifier (hash) for the cache entry
  - Step Ref Name: The workflow step reference name
  - Timestamp: When the cache entry was created
  - Size: Disk usage in bytes

Examples:
  kfg sys gc ls`,
	Run: func(cmd *cobra.Command, args []string) {
		entries, err := listCacheEntries()
		if err != nil {
			logger.Error("sys:gc:ls", err.Error())
			os.Exit(1)
		}

		if len(entries) == 0 {
			fmt.Println("No cache entries found")
			return
		}

		// Sort entries by timestamp (most recent first)
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Timestamp.After(entries[j].Timestamp)
		})

		// Print header
		fmt.Printf("%-16s %-30s %-25s %10s\n", "ID", "STEP REF NAME", "TIMESTAMP", "SIZE")
		fmt.Println(strings.Repeat("-", 85))

		// Print entries
		for _, entry := range entries {
			timestamp := entry.Timestamp.Format("2006-01-02T15:04:05Z")
			size := formatSize(entry.Size)
			fmt.Printf("%-16s %-30s %-25s %10s\n", entry.ID, entry.StepRefName, timestamp, size)
		}

		// Print summary
		totalSize := int64(0)
		for _, entry := range entries {
			totalSize += entry.Size
		}
		fmt.Println()
		fmt.Printf("Total: %d entries, %s\n", len(entries), formatSize(totalSize))
	},
}

// gcInspectCmd shows detailed metadata for a cache entry
var gcInspectCmd = &cobra.Command{
	Use:   "inspect <id>",
	Short: "Show detailed metadata for a cache entry",
	Long: `Show detailed metadata for a specific cache entry.

Displays:
  - Cache entry ID and path
  - Step reference name
  - Timestamp
  - Disk usage
  - Artifacts list
  - Output metadata (if present)

Examples:
  kfg sys gc inspect abc123`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entryID := args[0]
		entry, err := inspectCacheEntry(entryID)
		if err != nil {
			logger.Error("sys:gc:inspect", err.Error())
			os.Exit(1)
		}

		// Print detailed metadata
		fmt.Printf("Cache Entry: %s\n", entry.ID)
		fmt.Printf("Path: %s\n", entry.Path)
		fmt.Println(strings.Repeat("-", 60))
		fmt.Printf("Step Ref Name: %s\n", entry.StepRefName)
		fmt.Printf("Timestamp: %s\n", entry.Timestamp.Format("2006-01-02T15:04:05Z"))
		fmt.Printf("Size: %s\n", formatSize(entry.Size))
		fmt.Println()

		// Print artifacts from metadata
		if len(entry.Artifacts) > 0 {
			fmt.Printf("Artifacts (%d):\n", len(entry.Artifacts))
			for _, artifact := range entry.Artifacts {
				fmt.Printf("  - %s\n", artifact)
			}
		} else {
			fmt.Println("Artifacts: (none)")
		}

		// Print output
		fmt.Println()
		if entry.HasOutput {
			fmt.Printf("Output:\n")
			fmt.Printf("  Name: %s\n", entry.OutputName)
			// Show first 100 chars of value, or full value if short
			valuePreview := entry.OutputValue
			if len(valuePreview) > 100 {
				valuePreview = valuePreview[:100] + "..."
			}
			fmt.Printf("  Value: %s\n", valuePreview)
		} else {
			fmt.Println("Output: (none)")
		}
	},
}

// gcRmCmd removes specific cache entries
var gcRmCmd = &cobra.Command{
	Use:   "rm <id> [<id>...]",
	Short: "Remove specific cache entries",
	Long: `Remove one or more cache entries from storage.

Each entry is identified by its stable ID (hash).

Examples:
  kfg sys gc rm abc123
  kfg sys gc rm abc123 def456`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, entryID := range args {
			err := removeCacheEntry(entryID)
			if err != nil {
				logger.Error("sys:gc:rm", fmt.Sprintf("Failed to remove %s: %v", entryID, err))
				os.Exit(1)
			}
			fmt.Printf("Removed cache entry: %s\n", entryID)
		}
	},
}

// gcPruneCmd removes old or unused cache entries
var gcPruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove old or unused cache entries",
	Long: `Remove cache entries according to the prune policy.

Currently implemented prune policy: remove entries older than 30 days.

Examples:
  kfg sys gc prune`,
	Run: func(cmd *cobra.Command, args []string) {
		entries, err := listCacheEntries()
		if err != nil {
			logger.Error("sys:gc:prune", err.Error())
			os.Exit(1)
		}

		// Prune policy: entries older than 30 days
		cutoff := time.Now().AddDate(0, 0, -30)
		pruned := 0
		prunedSize := int64(0)

		for _, entry := range entries {
			if entry.Timestamp.Before(cutoff) {
				err := removeCacheEntry(entry.ID)
				if err != nil {
					logger.Warn("sys:gc:prune", fmt.Sprintf("Failed to remove %s: %v", entry.ID, err))
					continue
				}
				pruned++
				prunedSize += entry.Size
				fmt.Printf("Pruned: %s (%s, %s)\n", entry.ID, entry.StepRefName, entry.Timestamp.Format("2006-01-02"))
			}
		}

		if pruned == 0 {
			fmt.Println("No entries to prune")
		} else {
			fmt.Printf("\nPruned %d entries, freed %s\n", pruned, formatSize(prunedSize))
		}
	},
}

// gcDuCmd shows disk usage of cache entries
var gcDuCmd = &cobra.Command{
	Use:   "du",
	Short: "Show disk usage of cache entries",
	Long: `Report disk usage for persisted cache entries.

Shows:
  - Per-entry disk usage
  - Total disk usage
  - Cache directory location

Examples:
  kfg sys gc du`,
	Run: func(cmd *cobra.Command, args []string) {
		cacheDir := getCacheDir()
		fmt.Printf("Cache Directory: %s\n", cacheDir)
		fmt.Println(strings.Repeat("-", 60))

		entries, err := listCacheEntries()
		if err != nil {
			logger.Error("sys:gc:du", err.Error())
			os.Exit(1)
		}

		if len(entries) == 0 {
			fmt.Println("No cache entries found")
			return
		}

		// Sort by size (largest first)
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Size > entries[j].Size
		})

		// Print per-entry usage
		fmt.Printf("%-16s %-30s %10s\n", "ID", "STEP REF NAME", "SIZE")
		fmt.Println(strings.Repeat("-", 60))
		for _, entry := range entries {
			fmt.Printf("%-16s %-30s %10s\n", entry.ID, entry.StepRefName, formatSize(entry.Size))
		}

		// Print total
		totalSize := int64(0)
		for _, entry := range entries {
			totalSize += entry.Size
		}
		fmt.Println()
		fmt.Printf("Total: %d entries, %s\n", len(entries), formatSize(totalSize))
	},
}

func init() {
	// Add gcCmd to sysCmd
	sysCmd.AddCommand(gcCmd)

	// Add subcommands to gcCmd
	gcCmd.AddCommand(gcLsCmd)
	gcCmd.AddCommand(gcInspectCmd)
	gcCmd.AddCommand(gcRmCmd)
	gcCmd.AddCommand(gcPruneCmd)
	gcCmd.AddCommand(gcDuCmd)
}

// getCacheDir returns the cache directory path
func getCacheDir() string {
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

// listCacheEntries reads all cache entries from the cache directory
func listCacheEntries() ([]CacheEntry, error) {
	cacheDir := getCacheDir()

	// Check if cache directory exists
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		return []CacheEntry{}, nil
	}

	// Read cache directory
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache directory: %w", err)
	}

	var cacheEntries []CacheEntry
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		entryPath := filepath.Join(cacheDir, entry.Name())
		cacheEntry, err := readCacheEntry(entryPath)
		if err != nil {
			// Skip entries with invalid metadata
			logger.Warn("sys:gc", fmt.Sprintf("Skipping invalid cache entry %s: %v", entry.Name(), err))
			continue
		}

		cacheEntries = append(cacheEntries, cacheEntry)
	}

	return cacheEntries, nil
}

// readCacheEntry reads a single cache entry from its directory
func readCacheEntry(entryPath string) (CacheEntry, error) {
	entryID := filepath.Base(entryPath)

	// Read metadata.yaml
	metadataPath := filepath.Join(entryPath, "metadata.yaml")
	metadataData, err := os.ReadFile(metadataPath)
	if err != nil {
		return CacheEntry{}, fmt.Errorf("failed to read metadata.yaml: %w", err)
	}

	var metadata CacheMetadata
	err = yaml.Unmarshal(metadataData, &metadata)
	if err != nil {
		return CacheEntry{}, fmt.Errorf("failed to parse metadata.yaml: %w", err)
	}

	// Parse timestamp
	timestamp, err := time.Parse("2006-01-02T15:04:05Z", metadata.Timestamp)
	if err != nil {
		// Try alternative format without 'Z'
		timestamp, err = time.Parse("2006-01-02T15:04:05", metadata.Timestamp)
		if err != nil {
			timestamp = time.Time{} // Use zero time if parsing fails
		}
	}

	// Calculate size
	size, err := calculateDirSize(entryPath)
	if err != nil {
		size = 0
	}

	// Check for artifacts directory
	artifactsDir := filepath.Join(entryPath, "artifacts")
	if _, err := os.Stat(artifactsDir); os.IsNotExist(err) {
		artifactsDir = ""
	}

	// Parse output
	var hasOutput bool
	var outputName string
	var outputValue string
	if metadata.Output != nil {
		hasOutput = true
		outputName = metadata.Output.Name
		// Decode base64 value
		if metadata.Output.ValueEncoded != "" {
			decoded, err := base64.StdEncoding.DecodeString(metadata.Output.ValueEncoded)
			if err != nil {
				outputValue = "(decode error)"
			} else {
				outputValue = string(decoded)
			}
		}
	}

	// Read artifacts from metadata or fallback to artifact_paths.txt (legacy)
	var artifacts []string
	if len(metadata.Artifacts) > 0 {
		// New format: read from metadata.yaml Artifacts field
		artifacts = metadata.Artifacts
	} else {
		// Legacy fallback: read from artifact_paths.txt
		artifactPathsFile := filepath.Join(entryPath, "artifact_paths.txt")
		if data, err := os.ReadFile(artifactPathsFile); err == nil {
			// Split by newline and filter empty lines
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if line != "" {
					artifacts = append(artifacts, line)
				}
			}
		}
	}

	return CacheEntry{
		ID:           entryID,
		Path:         entryPath,
		StepRefName:  metadata.StepRefName,
		Timestamp:    timestamp,
		Size:         size,
		ArtifactsDir: artifactsDir,
		Artifacts:    artifacts,
		HasOutput:    hasOutput,
		OutputName:   outputName,
		OutputValue:  outputValue,
	}, nil
}

// inspectCacheEntry inspects a specific cache entry by ID
func inspectCacheEntry(entryID string) (CacheEntry, error) {
	cacheDir := getCacheDir()
	entryPath := filepath.Join(cacheDir, entryID)

	// Check if entry exists
	if _, err := os.Stat(entryPath); os.IsNotExist(err) {
		return CacheEntry{}, fmt.Errorf("cache entry not found: %s", entryID)
	}

	return readCacheEntry(entryPath)
}

// removeCacheEntry removes a cache entry by ID
func removeCacheEntry(entryID string) error {
	cacheDir := getCacheDir()
	entryPath := filepath.Join(cacheDir, entryID)

	// Check if entry exists
	if _, err := os.Stat(entryPath); os.IsNotExist(err) {
		return fmt.Errorf("cache entry not found: %s", entryID)
	}

	// Remove entry directory
	return os.RemoveAll(entryPath)
}

// calculateDirSize calculates the total size of a directory
func calculateDirSize(dirPath string) (int64, error) {
	var size int64
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			size += info.Size()
		}
		return nil
	})
	return size, err
}



// formatSize formats a size in bytes to a human-readable string
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}