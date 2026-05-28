package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/seregatte/kfg/src/internal/cache"
	"github.com/seregatte/kfg/src/internal/logger"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// cacheCmd represents the cache command group
var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Cache operations for Step results",
	Long: `Cache operations for managing persisted Step execution results.

These commands operate on cache entries stored under KFG_STORE_DIR/cache
(defaults to ~/.kfg/store/cache).

Runtime subcommands (consumed by shell wrappers):
  exists   Check if a cache entry exists
  store    Store Step results to cache
  restore  Restore Step results from cache

Admin subcommands:
  ls       List cache entries with metadata
  inspect  Show detailed metadata for a cache entry
  rm       Remove specific cache entries
  prune    Remove old cache entries
  du       Show disk usage of cache entries

Examples:
  kfg sys cache exists ctx7.steps.install
  kfg sys cache store ctx7.steps.install --workdir /path
  kfg sys cache restore ctx7.steps.install --workdir /path
  kfg sys cache ls
  kfg sys cache inspect ctx7.steps.install
  kfg sys cache rm ctx7.steps.install
  kfg sys cache prune
  kfg sys cache du`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// cacheExistsCmd checks if a cache entry exists
var cacheExistsCmd = &cobra.Command{
	Use:   "exists <step-ref>",
	Short: "Check if a cache entry exists",
	Long: `Check if a valid cache entry exists for the given StepReference.name.

Exits with code 0 for cache hit, 1 for cache miss.

Examples:
  kfg sys cache exists ctx7.steps.install`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stepRefName := args[0]
		if cache.Exists(stepRefName) {
			os.Exit(0)
		}
		os.Exit(1)
	},
}

// cacheStoreCmd stores Step results to cache
var cacheStoreCmd = &cobra.Command{
	Use:   "store <step-ref> --workdir <path>",
	Short: "Store Step results to cache",
	Long: `Store Step execution results (artifacts and output) to the cache.

Reads a JSON object from stdin with fields:
  - before: string array of artifact paths before step
  - after: string array of artifact paths after step
  - declarative: string array of manifest-declared artifacts
  - output: object with name (string) and value (base64 string)

The --fs-before flag overrides the "before" field with a filesystem snapshot file.

Examples:
  echo '{"before":[],"after":["file.txt"],"declarative":[],"output":null}' | kfg sys cache store ctx7.steps.install --workdir /path`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stepRefName := args[0]
		workdir, _ := cmd.Flags().GetString("workdir")
		fsBefore, _ := cmd.Flags().GetString("fs-before")

		if workdir == "" {
			workdir, _ = os.Getwd()
		}

		// Read JSON from stdin
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			logger.Error("sys:cache:store", fmt.Sprintf("failed to read stdin: %v", err))
			os.Exit(1)
		}

		// Parse input
		input, err := cache.ParseStoreInput(data)
		if err != nil {
			logger.Error("sys:cache:store", fmt.Sprintf("failed to parse input: %v", err))
			os.Exit(1)
		}

		// Override "before" with filesystem snapshot if provided
		if fsBefore != "" {
			beforePaths, err := cache.ReadFsSnapshot(fsBefore)
			if err != nil {
				logger.Error("sys:cache:store", fmt.Sprintf("failed to read fs snapshot: %v", err))
				os.Exit(1)
			}
			input.Before = beforePaths
		}

		// Store
		result, err := cache.Store(stepRefName, workdir, input)
		if err != nil {
			logger.Error("sys:cache:store", fmt.Sprintf("failed to store: %v", err))
			os.Exit(1)
		}

		logger.Detail("sys:cache:store", fmt.Sprintf("Cached %d artifacts for %s", len(result.Artifacts), stepRefName))
	},
}

// cacheRestoreCmd restores Step results from cache
var cacheRestoreCmd = &cobra.Command{
	Use:   "restore <step-ref> --workdir <path>",
	Short: "Restore Step results from cache",
	Long: `Restore cached Step results and emit shell eval-safe output.

The output lines can be eval'd in bash to restore artifacts and outputs:
  eval "$(kfg sys cache restore ctx7.steps.install --workdir /path)"

Examples:
  kfg sys cache restore ctx7.steps.install --workdir /path`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stepRefName := args[0]
		workdir, _ := cmd.Flags().GetString("workdir")

		if workdir == "" {
			workdir, _ = os.Getwd()
		}

		output, err := cache.RestoreToStdout(stepRefName, workdir)
		if err != nil {
			logger.Error("sys:cache:restore", fmt.Sprintf("failed to restore: %v", err))
			os.Exit(1)
		}

		fmt.Println(output)
	},
}

// cacheLsCmd lists cache entries
var cacheLsCmd = &cobra.Command{
	Use:   "ls [--json|--yaml]",
	Short: "List cache entries with metadata",
	Long: `List all cache entries with their step reference names and metadata.

Examples:
  kfg sys cache ls
  kfg sys cache ls --json
  kfg sys cache ls --yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonOutput, _ := cmd.Flags().GetBool("json")
		yamlOutput, _ := cmd.Flags().GetBool("yaml")

		entries, err := listCacheEntries()
		if err != nil {
			logger.Error("sys:cache:ls", err.Error())
			os.Exit(1)
		}

		if len(entries) == 0 {
			if jsonOutput {
				fmt.Println("[]")
			} else if yamlOutput {
				fmt.Println("[]")
			} else {
				fmt.Println("No cache entries found")
			}
			return
		}

		// Sort by timestamp (most recent first)
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Timestamp.After(entries[j].Timestamp)
		})

		if jsonOutput {
			outputJSON(entries)
		} else if yamlOutput {
			outputYAML(entries)
		} else {
			outputTable(entries)
		}
	},
}

// cacheInspectCmd shows detailed metadata for a cache entry
var cacheInspectCmd = &cobra.Command{
	Use:   "inspect <step-ref> [--json|--yaml]",
	Short: "Show detailed metadata for a cache entry",
	Long: `Show detailed metadata for a cache entry identified by StepReference.name.

Examples:
  kfg sys cache inspect ctx7.steps.install
  kfg sys cache inspect ctx7.steps.install --json`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stepRefName := args[0]
		jsonOutput, _ := cmd.Flags().GetBool("json")
		yamlOutput, _ := cmd.Flags().GetBool("yaml")

		entryPath := cache.GetEntryPath(stepRefName)
		entry, err := cache.ReadCacheEntry(entryPath)
		if err != nil {
			logger.Error("sys:cache:inspect", fmt.Sprintf("entry not found: %s", stepRefName))
			os.Exit(1)
		}

		if jsonOutput {
			inspectJSON(entry)
		} else if yamlOutput {
			inspectYAML(entry)
		} else {
			inspectText(entry)
		}
	},
}

// cacheRmCmd removes cache entries
var cacheRmCmd = &cobra.Command{
	Use:   "rm <step-ref> [<step-ref>...]",
	Short: "Remove cache entries by StepReference.name",
	Long: `Remove one or more cache entries identified by StepReference.name.

Examples:
  kfg sys cache rm ctx7.steps.install
  kfg sys cache rm ctx7.steps.install openspec.steps.install`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, stepRefName := range args {
			entryPath := cache.GetEntryPath(stepRefName)
			if _, err := os.Stat(entryPath); os.IsNotExist(err) {
				logger.Warn("sys:cache:rm", fmt.Sprintf("entry not found: %s", stepRefName))
				continue
			}

			if err := os.RemoveAll(entryPath); err != nil {
				logger.Error("sys:cache:rm", fmt.Sprintf("failed to remove %s: %v", stepRefName, err))
				os.Exit(1)
			}
			fmt.Printf("Removed: %s\n", stepRefName)
		}
	},
}

// cachePruneCmd removes old cache entries
var cachePruneCmd = &cobra.Command{
	Use:   "prune [--json|--yaml]",
	Short: "Remove cache entries older than 30 days",
	Long: `Remove cache entries with timestamps older than 30 days.

Examples:
  kfg sys cache prune
  kfg sys cache prune --json`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonOutput, _ := cmd.Flags().GetBool("json")
		yamlOutput, _ := cmd.Flags().GetBool("yaml")

		entries, err := listCacheEntries()
		if err != nil {
			logger.Error("sys:cache:prune", err.Error())
			os.Exit(1)
		}

		cutoff := time.Now().AddDate(0, 0, -30)
		var pruned []string
		var prunedSize int64

		for _, entry := range entries {
			if entry.Timestamp.Before(cutoff) {
				if err := os.RemoveAll(entry.Path); err != nil {
					logger.Warn("sys:cache:prune", fmt.Sprintf("failed to remove %s: %v", entry.StepRefName, err))
					continue
				}
				pruned = append(pruned, entry.StepRefName)
				prunedSize += entry.Size
			}
		}

		if jsonOutput {
			result := map[string]interface{}{
				"pruned":     pruned,
				"count":      len(pruned),
				"freedBytes": prunedSize,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
		} else if yamlOutput {
			result := map[string]interface{}{
				"pruned":     pruned,
				"count":      len(pruned),
				"freedBytes": prunedSize,
			}
			data, _ := yaml.Marshal(result)
			fmt.Println(string(data))
		} else {
			if len(pruned) == 0 {
				fmt.Println("No entries to prune")
			} else {
				for _, name := range pruned {
					fmt.Printf("Pruned: %s\n", name)
				}
				fmt.Printf("\nPruned %d entries, freed %s\n", len(pruned), cache.FormatSize(prunedSize))
			}
		}
	},
}

// cacheDuCmd shows disk usage of cache entries
var cacheDuCmd = &cobra.Command{
	Use:   "du [--json|--yaml]",
	Short: "Show disk usage of cache entries",
	Long: `Report disk usage for persisted cache entries.

Examples:
  kfg sys cache du
  kfg sys cache du --json`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonOutput, _ := cmd.Flags().GetBool("json")
		yamlOutput, _ := cmd.Flags().GetBool("yaml")

		cacheDir := cache.GetCacheDir()
		entries, err := listCacheEntries()
		if err != nil {
			logger.Error("sys:cache:du", err.Error())
			os.Exit(1)
		}

		if jsonOutput {
			duJSON(cacheDir, entries)
		} else if yamlOutput {
			duYAML(cacheDir, entries)
		} else {
			duText(cacheDir, entries)
		}
	},
}

func init() {
	// Add cacheCmd to sysCmd
	sysCmd.AddCommand(cacheCmd)

	// Add runtime subcommands
	cacheCmd.AddCommand(cacheExistsCmd)
	cacheCmd.AddCommand(cacheStoreCmd)
	cacheCmd.AddCommand(cacheRestoreCmd)

	// Add admin subcommands
	cacheCmd.AddCommand(cacheLsCmd)
	cacheCmd.AddCommand(cacheInspectCmd)
	cacheCmd.AddCommand(cacheRmCmd)
	cacheCmd.AddCommand(cachePruneCmd)
	cacheCmd.AddCommand(cacheDuCmd)

	// Add flags
	cacheStoreCmd.Flags().StringP("workdir", "w", "", "Working directory (default: $PWD)")
	cacheStoreCmd.Flags().String("fs-before", "", "Filesystem snapshot file before step execution")
	cacheRestoreCmd.Flags().StringP("workdir", "w", "", "Working directory (default: $PWD)")

	cacheLsCmd.Flags().Bool("json", false, "Output in JSON format")
	cacheLsCmd.Flags().Bool("yaml", false, "Output in YAML format")

	cacheInspectCmd.Flags().Bool("json", false, "Output in JSON format")
	cacheInspectCmd.Flags().Bool("yaml", false, "Output in YAML format")

	cachePruneCmd.Flags().Bool("json", false, "Output in JSON format")
	cachePruneCmd.Flags().Bool("yaml", false, "Output in YAML format")

	cacheDuCmd.Flags().Bool("json", false, "Output in JSON format")
	cacheDuCmd.Flags().Bool("yaml", false, "Output in YAML format")
}

// listCacheEntries reads all cache entries from the cache directory
func listCacheEntries() ([]cache.CacheEntry, error) {
	cacheDir := cache.GetCacheDir()

	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		return []cache.CacheEntry{}, nil
	}

	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache directory: %w", err)
	}

	var cacheEntries []cache.CacheEntry
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		entryPath := cache.JoinPaths(cacheDir, entry.Name())
		cacheEntry, err := cache.ReadCacheEntry(entryPath)
		if err != nil {
			logger.Warn("sys:cache", fmt.Sprintf("Skipping invalid cache entry %s: %v", entry.Name(), err))
			continue
		}

		cacheEntries = append(cacheEntries, *cacheEntry)
	}

	return cacheEntries, nil
}

// outputTable outputs entries in table format
func outputTable(entries []cache.CacheEntry) {
	fmt.Printf("%-30s %-25s %10s\n", "STEP REF NAME", "TIMESTAMP", "SIZE")
	fmt.Println(strings.Repeat("-", 70))

	for _, entry := range entries {
		timestamp := entry.Timestamp.Format("2006-01-02T15:04:05Z")
		size := cache.FormatSize(entry.Size)
		fmt.Printf("%-30s %-25s %10s\n", entry.StepRefName, timestamp, size)
	}

	totalSize := int64(0)
	for _, entry := range entries {
		totalSize += entry.Size
	}
	fmt.Println()
	fmt.Printf("Total: %d entries, %s\n", len(entries), cache.FormatSize(totalSize))
}

// outputJSON outputs entries in JSON format
func outputJSON(entries []cache.CacheEntry) {
	type entryJSON struct {
		StepRef string `json:"stepRef"`
		Time    string `json:"timestamp"`
		Size    int64  `json:"size"`
	}

	var result []entryJSON
	for _, entry := range entries {
		result = append(result, entryJSON{
			StepRef: entry.StepRefName,
			Time:    entry.Timestamp.Format("2006-01-02T15:04:05Z"),
			Size:    entry.Size,
		})
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(data))
}

// outputYAML outputs entries in YAML format
func outputYAML(entries []cache.CacheEntry) {
	type entryYAML struct {
		StepRef string `yaml:"stepRef"`
		Time    string `yaml:"timestamp"`
		Size    int64  `yaml:"size"`
	}

	var result []entryYAML
	for _, entry := range entries {
		result = append(result, entryYAML{
			StepRef: entry.StepRefName,
			Time:    entry.Timestamp.Format("2006-01-02T15:04:05Z"),
			Size:    entry.Size,
		})
	}

	data, _ := yaml.Marshal(result)
	fmt.Println(string(data))
}

// inspectText outputs entry in text format
func inspectText(entry *cache.CacheEntry) {
	fmt.Printf("Step Ref: %s\n", entry.StepRefName)
	fmt.Printf("Path: %s\n", entry.Path)
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("Timestamp: %s\n", entry.Timestamp.Format("2006-01-02T15:04:05Z"))
	fmt.Printf("Size: %s\n", cache.FormatSize(entry.Size))
	fmt.Println()

	if len(entry.Artifacts) > 0 {
		fmt.Printf("Artifacts (%d):\n", len(entry.Artifacts))
		for _, artifact := range entry.Artifacts {
			fmt.Printf("  - %s\n", artifact)
		}
	} else {
		fmt.Println("Artifacts: (none)")
	}

	fmt.Println()
	if entry.HasOutput {
		fmt.Printf("Output:\n")
		fmt.Printf("  Name: %s\n", entry.OutputName)
		fmt.Printf("  Value: %s\n", entry.OutputValue)
	} else {
		fmt.Println("Output: (none)")
	}
}

// inspectJSON outputs entry in JSON format
func inspectJSON(entry *cache.CacheEntry) {
	type outputJSON struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	result := struct {
		StepRef   string      `json:"stepRef"`
		Timestamp string      `json:"timestamp"`
		Size      int64       `json:"size"`
		Artifacts []string    `json:"artifacts"`
		Output    *outputJSON `json:"output,omitempty"`
	}{
		StepRef:   entry.StepRefName,
		Timestamp: entry.Timestamp.Format("2006-01-02T15:04:05Z"),
		Size:      entry.Size,
		Artifacts: entry.Artifacts,
	}

	if entry.HasOutput {
		result.Output = &outputJSON{
			Name:  entry.OutputName,
			Value: entry.OutputValue,
		}
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(data))
}

// inspectYAML outputs entry in YAML format
func inspectYAML(entry *cache.CacheEntry) {
	type outputYAML struct {
		Name  string `yaml:"name"`
		Value string `yaml:"value"`
	}

	result := struct {
		StepRef   string      `yaml:"stepRef"`
		Timestamp string      `yaml:"timestamp"`
		Size      int64       `yaml:"size"`
		Artifacts []string    `yaml:"artifacts"`
		Output    *outputYAML `yaml:"output,omitempty"`
	}{
		StepRef:   entry.StepRefName,
		Timestamp: entry.Timestamp.Format("2006-01-02T15:04:05Z"),
		Size:      entry.Size,
		Artifacts: entry.Artifacts,
	}

	if entry.HasOutput {
		result.Output = &outputYAML{
			Name:  entry.OutputName,
			Value: entry.OutputValue,
		}
	}

	data, _ := yaml.Marshal(result)
	fmt.Println(string(data))
}

// duText outputs disk usage in text format
func duText(cacheDir string, entries []cache.CacheEntry) {
	fmt.Printf("Cache Directory: %s\n", cacheDir)
	fmt.Println(strings.Repeat("-", 60))

	if len(entries) == 0 {
		fmt.Println("No cache entries found")
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Size > entries[j].Size
	})

	fmt.Printf("%-30s %10s\n", "STEP REF NAME", "SIZE")
	fmt.Println(strings.Repeat("-", 45))
	for _, entry := range entries {
		fmt.Printf("%-30s %10s\n", entry.StepRefName, cache.FormatSize(entry.Size))
	}

	totalSize := int64(0)
	for _, entry := range entries {
		totalSize += entry.Size
	}
	fmt.Println()
	fmt.Printf("Total: %d entries, %s\n", len(entries), cache.FormatSize(totalSize))
}

// duJSON outputs disk usage in JSON format
func duJSON(cacheDir string, entries []cache.CacheEntry) {
	type entryJSON struct {
		StepRef string `json:"stepRef"`
		Size    int64  `json:"size"`
	}

	var entryList []entryJSON
	totalSize := int64(0)
	for _, entry := range entries {
		entryList = append(entryList, entryJSON{
			StepRef: entry.StepRefName,
			Size:    entry.Size,
		})
		totalSize += entry.Size
	}

	result := struct {
		CacheDir   string      `json:"cacheDir"`
		Entries    []entryJSON `json:"entries"`
		TotalBytes int64       `json:"totalBytes"`
	}{
		CacheDir:   cacheDir,
		Entries:    entryList,
		TotalBytes: totalSize,
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(data))
}

// duYAML outputs disk usage in YAML format
func duYAML(cacheDir string, entries []cache.CacheEntry) {
	type entryYAML struct {
		StepRef string `yaml:"stepRef"`
		Size    int64  `yaml:"size"`
	}

	var entryList []entryYAML
	totalSize := int64(0)
	for _, entry := range entries {
		entryList = append(entryList, entryYAML{
			StepRef: entry.StepRefName,
			Size:    entry.Size,
		})
		totalSize += entry.Size
	}

	result := struct {
		CacheDir   string      `yaml:"cacheDir"`
		Entries    []entryYAML `yaml:"entries"`
		TotalBytes int64       `yaml:"totalBytes"`
	}{
		CacheDir:   cacheDir,
		Entries:    entryList,
		TotalBytes: totalSize,
	}

	data, _ := yaml.Marshal(result)
	fmt.Println(string(data))
}
