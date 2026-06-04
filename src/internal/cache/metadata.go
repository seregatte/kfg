package cache

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// CacheMetadata represents the metadata.yaml content in a cache entry.
type CacheMetadata struct {
	StepRefName string          `yaml:"stepRefName"`
	Timestamp   string          `yaml:"timestamp"`
	Output      *OutputMetadata `yaml:"output,omitempty"`
	Artifacts   []string        `yaml:"artifacts,omitempty"`
}

// OutputMetadata represents cached output data.
type OutputMetadata struct {
	Name         string `yaml:"name"`
	ValueEncoded string `yaml:"valueEncoded"`
}

// CacheEntry represents a parsed cache entry with metadata and computed fields.
type CacheEntry struct {
	Identity    string
	Path        string
	StepRefName string
	Timestamp   time.Time
	Size        int64
	HasOutput   bool
	OutputName  string
	OutputValue string
	Artifacts   []string
}

// ReadMetadata reads and parses metadata.yaml from a cache entry directory.
func ReadMetadata(entryPath string) (*CacheMetadata, error) {
	metadataPath := filepath.Join(entryPath, "metadata.yaml")
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata.yaml: %w", err)
	}

	var metadata CacheMetadata
	if err := yaml.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata.yaml: %w", err)
	}

	return &metadata, nil
}

// WriteMetadata writes metadata.yaml to a cache entry directory.
func WriteMetadata(entryPath string, metadata *CacheMetadata) error {
	data, err := yaml.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	metadataPath := filepath.Join(entryPath, "metadata.yaml")
	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata.yaml: %w", err)
	}

	return nil
}

// ReadCacheEntry reads a complete cache entry including metadata and computed fields.
func ReadCacheEntry(entryPath string) (*CacheEntry, error) {
	metadata, err := ReadMetadata(entryPath)
	if err != nil {
		return nil, err
	}

	// Parse timestamp
	timestamp, err := time.Parse("2006-01-02T15:04:05Z", metadata.Timestamp)
	if err != nil {
		// Try alternative format without 'Z'
		timestamp, err = time.Parse("2006-01-02T15:04:05", metadata.Timestamp)
		if err != nil {
			timestamp = time.Time{}
		}
	}

	// Calculate size
	size, err := calculateDirSize(entryPath)
	if err != nil {
		size = 0
	}

	// Parse output
	var hasOutput bool
	var outputName string
	var outputValue string
	if metadata.Output != nil {
		hasOutput = true
		outputName = metadata.Output.Name
		if metadata.Output.ValueEncoded != "" {
			decoded, err := base64.StdEncoding.DecodeString(metadata.Output.ValueEncoded)
			if err != nil {
				outputValue = "(decode error)"
			} else {
				outputValue = string(decoded)
			}
		}
	}

	return &CacheEntry{
		Identity:    filepath.Base(entryPath),
		Path:        entryPath,
		StepRefName: metadata.StepRefName,
		Timestamp:   timestamp,
		Size:        size,
		HasOutput:   hasOutput,
		OutputName:  outputName,
		OutputValue: outputValue,
		Artifacts:   metadata.Artifacts,
	}, nil
}

// calculateDirSize calculates the total size of a directory.
func calculateDirSize(dirPath string) (int64, error) {
	var size int64
	err := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
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

// FormatSize formats a size in bytes to a human-readable string.
func FormatSize(bytes int64) string {
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
