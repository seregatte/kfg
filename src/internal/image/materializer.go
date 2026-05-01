// Package image provides image materialization for workspace integration.
// This file implements the workspace materialization and restoration operations.
package image

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/seregatte/kfg/src/internal/logger"
)

const (
	// WorkspaceSubdir is the subdirectory for workspace instance metadata
	WorkspaceSubdir = ".workspace"

	// BackupSubdir is the subdirectory for backup archives
	BackupSubdir = "backup"

	// InstanceFile is the filename for instance metadata
	InstanceFile = "instance.json"
)

// InstanceMetadata tracks active workspace instances.
type InstanceMetadata struct {
	// Name is the unique instance name
	Name string `json:"name"`

	// ImageRef is the image reference (name:tag)
	ImageRef string `json:"image_ref"`

	// WorkspaceRoot is the absolute path to the workspace
	WorkspaceRoot string `json:"workspace_root"`

	// StartedAt is the ISO 8601 timestamp when the instance was started
	StartedAt string `json:"started_at"`

	// ImageDigest is the digest of the materialized image
	ImageDigest string `json:"image_digest"`

	// MaterializedPaths are the paths that were materialized from the image
	MaterializedPaths []string `json:"materialized_paths"`
}

// Materializer handles image materialization to workspace.
type Materializer struct {
	imageStore *ImageStore
}

// NewMaterializer creates a new Materializer instance.
func NewMaterializer(storeDir string) *Materializer {
	return &Materializer{
		imageStore: NewImageStore(storeDir),
	}
}

// Start materializes an image into the workspace with backup.
func (m *Materializer) Start(imageRef string, workspaceRoot string, instanceName string) error {
	// Validate required parameters
	if imageRef == "" {
		return fmt.Errorf("image reference is required")
	}
	
	if workspaceRoot == "" {
		return fmt.Errorf("workspace root directory is required")
	}
	
	if instanceName == "" {
		return fmt.Errorf("instance name is required")
	}
	
	// Validate instance name format (alphanumeric, dashes, underscores)
	if !isValidInstanceName(instanceName) {
		return fmt.Errorf("instance name must contain only alphanumeric characters, dashes, and underscores")
	}

	// Resolve workspace root (default to current directory)
	if workspaceRoot == "" {
		workspaceRoot = "."
	}

	// Resolve image reference
	name, tag := ResolveRef(imageRef)

	// Check if image exists
	metadata, imageDir, err := m.imageStore.LoadImage(imageRef)
	if err != nil {
		return fmt.Errorf("image %s:%s not found in store: %w", name, tag, err)
	}

	// Check if instance already exists
	instanceDir := filepath.Join(m.imageStore.GetStoreDir(), WorkspaceSubdir, instanceName)
	if _, err := os.Stat(instanceDir); err == nil {
		// Instance record exists - check if it's active
		instanceFile := filepath.Join(instanceDir, InstanceFile)
		if _, err := os.Stat(instanceFile); err == nil {
			return fmt.Errorf("instance '%s' already exists - use a different name or stop the existing instance first", instanceName)
		}
	}

	// Create instance directory early (needed for backup location)
	if err := os.MkdirAll(instanceDir, 0755); err != nil {
		return fmt.Errorf("failed to create instance directory: %w", err)
	}

	// Compute artifact paths from image metadata
	artifactPaths, err := m.computeArtifactPaths(imageDir)
	if err != nil {
		return fmt.Errorf("failed to compute artifact paths: %w", err)
	}

	// Find which artifact paths conflict with existing workspace files
	conflictingPaths, err := m.findConflictingPaths(workspaceRoot, artifactPaths)
	if err != nil {
		return fmt.Errorf("failed to find conflicting paths: %w", err)
	}

	// Create scoped backup only for conflicting paths
	if len(conflictingPaths) > 0 {
		logger.Info("workspace:start", fmt.Sprintf("Creating backup of %d conflicting file(s)...", len(conflictingPaths)))
		backupCount, err := m.createScopedBackup(workspaceRoot, instanceName, conflictingPaths)
		if err != nil {
			return fmt.Errorf("backup creation failed - aborting to preserve data safety: %w", err)
		}
		logger.Info("workspace:start", fmt.Sprintf("Backed up %d conflicting file(s)", backupCount))
	} else {
		logger.Info("workspace:start", "No conflicting files - backup skipped")
	}

	// Materialize image files to workspace
	logger.Info("workspace:start", fmt.Sprintf("Materializing image %s:%s to workspace...", metadata.Name, metadata.Tag))
	materializedPaths, err := m.materializeFiles(imageDir, workspaceRoot)
	if err != nil {
		return fmt.Errorf("materialization failed: %w", err)
	}

	// Create instance metadata with materialized paths
	instance := InstanceMetadata{
		Name:              instanceName,
		ImageRef:          imageRef,
		WorkspaceRoot:     workspaceRoot,
		StartedAt:         time.Now().UTC().Format(time.RFC3339),
		ImageDigest:       metadata.ImageDigest,
		MaterializedPaths: materializedPaths,
	}

	// Save instance metadata
	instanceFile := filepath.Join(instanceDir, InstanceFile)
	if err := m.saveInstance(instance, instanceFile); err != nil {
		return fmt.Errorf("failed to save instance metadata: %w", err)
	}

	logger.Info("workspace:start", fmt.Sprintf("Instance '%s' started successfully", instanceName))
	logger.Info("workspace:start", fmt.Sprintf("Image: %s:%s (digest: %s)", metadata.Name, metadata.Tag, metadata.GetShortDigest()))
	logger.Info("workspace:start", fmt.Sprintf("Workspace: %s", workspaceRoot))
	logger.Info("workspace:start", fmt.Sprintf("Materialized %d file(s)", len(materializedPaths)))

	return nil
}

// Stop restores workspace from backup and cleans up instance.
func (m *Materializer) Stop(instanceName string) error {
	// Validate inputs
	if instanceName == "" {
		return fmt.Errorf("instance name is required")
	}

	// Check if instance exists
	instanceDir := filepath.Join(m.imageStore.GetStoreDir(), WorkspaceSubdir, instanceName)
	instanceFile := filepath.Join(instanceDir, InstanceFile)

	instance, err := m.loadInstance(instanceFile)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Info("workspace:stop", fmt.Sprintf("Instance '%s' not found - nothing to stop (idempotent success)", instanceName))
			return nil
		}
		return fmt.Errorf("failed to load instance metadata: %w", err)
	}

	// Remove materialized paths before restoring backup
	removedCount := 0
	for _, path := range instance.MaterializedPaths {
		fullPath := filepath.Join(instance.WorkspaceRoot, path)
		if err := m.removeMaterializedPath(fullPath); err != nil {
			// Log but continue - idempotent cleanup
			logger.Info("workspace:stop", fmt.Sprintf("%s already removed", path))
		} else {
			logger.Info("workspace:stop", fmt.Sprintf("Removing materialized file: %s", path))
			removedCount++
		}
	}
	logger.Info("workspace:stop", fmt.Sprintf("Removed %d materialized artifact(s)", removedCount))

	// Restore from backup if it exists
	backupDir := filepath.Join(instanceDir, BackupSubdir)
	backupDataDir := filepath.Join(backupDir, "data")

	if _, err := os.Stat(backupDataDir); err == nil {
		logger.Info("workspace:stop", "Restoring workspace from backup...")
		if err := m.restoreBackup(backupDataDir, instance.WorkspaceRoot); err != nil {
			return fmt.Errorf("restore failed: %w", err)
		}
		logger.Info("workspace:stop", "Workspace restored from backup")

		// Delete backup after restore (backup is consumed)
		if err := os.RemoveAll(backupDir); err != nil {
			logger.Warn("workspace:stop", fmt.Sprintf("Failed to delete backup directory: %v", err))
		}
	} else {
		logger.Info("workspace:stop", "No backup found - cleanup only")
	}

	// Remove instance record
	if err := os.RemoveAll(instanceDir); err != nil {
		return fmt.Errorf("failed to cleanup instance directory: %w", err)
	}

	logger.Info("workspace:stop", fmt.Sprintf("Instance '%s' stopped and cleaned up", instanceName))

	return nil
}

// removeMaterializedPath removes a single materialized path with idempotent handling.
func (m *Materializer) removeMaterializedPath(fullPath string) error {
	// Check if path exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// Path doesn't exist - idempotent success
		return fmt.Errorf("path does not exist")
	}

	// Remove the path (file or directory)
	if err := os.RemoveAll(fullPath); err != nil {
		return err
	}

	// Clean up empty parent directories
	m.cleanupEmptyParentDirs(fullPath)

	return nil
}

// cleanupEmptyParentDirs removes empty parent directories after file removal.
func (m *Materializer) cleanupEmptyParentDirs(fullPath string) {
	// Walk up the directory tree and remove empty directories
	parentDir := filepath.Dir(fullPath)
	for {
		// Check if directory is empty
		isEmpty, err := isDirectoryEmpty(parentDir)
		if err != nil || !isEmpty {
			// Directory not empty or error - stop cleaning up
			break
		}

		// Remove empty directory
		if err := os.Remove(parentDir); err != nil {
			// Couldn't remove - stop
			break
		}

		// Move to next parent
		parentDir = filepath.Dir(parentDir)
	}
}

// createBackup creates a tar archive of the workspace.
func (m *Materializer) createBackup(workspaceRoot string, instanceName string) error {
	backupDir := filepath.Join(m.imageStore.GetStoreDir(), WorkspaceSubdir, instanceName, BackupSubdir)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// For MVP, we'll use a simple recursive copy approach
	backupDataDir := filepath.Join(backupDir, "data")
	if err := copyDirectory(workspaceRoot, backupDataDir); err != nil {
		return fmt.Errorf("failed to backup workspace: %w", err)
	}

	// Rename data directory to indicate it's a backup
	// (In production, this would be a tar archive)
	// For MVP, we keep it as a directory for easier restoration

	return nil
}

// restoreBackup restores workspace from backup archive.
func (m *Materializer) restoreBackup(backupDataDir string, workspaceRoot string) error {
	// For MVP, backup is stored as a directory instead of tar
	if _, err := os.Stat(backupDataDir); err != nil {
		if os.IsNotExist(err) {
			logger.Info("workspace:stop", "No backup data found - nothing to restore")
			return nil
		}
		return fmt.Errorf("failed to access backup: %w", err)
	}

	// Copy backup data to workspace
	if err := copyDirectory(backupDataDir, workspaceRoot); err != nil {
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	return nil
}

// materializeFiles copies image files to workspace and returns list of materialized paths.
func (m *Materializer) materializeFiles(imageDir string, workspaceRoot string) ([]string, error) {
	materializedPaths := []string{}

	// Walk image directory and copy all files (except metadata.json)
	err := filepath.Walk(imageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip metadata.json (not part of image content)
		if info.Name() == "metadata.json" {
			return nil
		}

		// Calculate relative path from image directory
		relPath, err := filepath.Rel(imageDir, path)
		if err != nil {
			return err
		}

		// Destination path in workspace
		destPath := filepath.Join(workspaceRoot, relPath)

		if info.IsDir() {
			// Create directory in workspace
			return os.MkdirAll(destPath, info.Mode())
		}

		// Copy file to workspace
		if err := copyFile(path, destPath); err != nil {
			return err
		}

		// Track materialized file path
		materializedPaths = append(materializedPaths, relPath)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return materializedPaths, nil
}

// saveInstance saves instance metadata to file.
func (m *Materializer) saveInstance(instance InstanceMetadata, filePath string) error {
	data, err := json.MarshalIndent(instance, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal instance metadata: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write instance file: %w", err)
	}

	return nil
}

// loadInstance loads instance metadata from file.
func (m *Materializer) loadInstance(filePath string) (*InstanceMetadata, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var instance InstanceMetadata
	if err := json.Unmarshal(data, &instance); err != nil {
		return nil, fmt.Errorf("failed to parse instance metadata: %w", err)
	}

	// Ensure MaterializedPaths is never nil (backward compatibility)
	if instance.MaterializedPaths == nil {
		instance.MaterializedPaths = []string{}
	}

	return &instance, nil
}

// isDirectoryEmpty checks if a directory is empty.
func isDirectoryEmpty(dir string) (bool, error) {
	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return true, nil // Non-existent directory is empty
	}

	// Open directory
	f, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// Read directory entries
	_, err = f.Readdirnames(1)
	if err == nil {
		return false, nil // Directory has at least one entry
	}

	// io.EOF means directory is empty
	return true, nil
}

// isValidInstanceName validates instance name format.
func isValidInstanceName(name string) bool {
	if name == "" {
		return false
	}
	
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return false
		}
	}
	
	return true
}

// Note: copyFile and copyDirectory are defined in store.go in the same package

// computeArtifactPaths returns the list of artifact paths from image metadata.
func (m *Materializer) computeArtifactPaths(imageDir string) ([]string, error) {
	// Load image metadata
	metadata, err := LoadMetadataFromDir(imageDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load image metadata: %w", err)
	}

	// Extract paths from Files manifest
	paths := make([]string, 0, len(metadata.Files))
	for path := range metadata.Files {
		paths = append(paths, path)
	}

	return paths, nil
}

// findConflictingPaths returns which artifact paths already exist in the workspace.
func (m *Materializer) findConflictingPaths(workspaceRoot string, artifactPaths []string) ([]string, error) {
	conflicting := []string{}

	for _, path := range artifactPaths {
		fullPath := filepath.Join(workspaceRoot, path)
		if _, err := os.Stat(fullPath); err == nil {
			// Path exists in workspace
			conflicting = append(conflicting, path)
		}
	}

	return conflicting, nil
}

// backupPath backs up a single path to the backup directory, preserving directory structure.
func (m *Materializer) backupPath(workspaceRoot string, relativePath string, backupDataDir string) error {
	srcPath := filepath.Join(workspaceRoot, relativePath)
	destPath := filepath.Join(backupDataDir, relativePath)

	// Get file info
	info, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("failed to stat path %s: %w", srcPath, err)
	}

	if info.IsDir() {
		// Copy directory recursively
		return copyDirectory(srcPath, destPath)
	}

	// Ensure parent directory exists in backup
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory %s: %w", destDir, err)
	}

	// Copy file
	return copyFile(srcPath, destPath)
}

// createScopedBackup creates backup of only the specified paths.
func (m *Materializer) createScopedBackup(workspaceRoot string, instanceName string, paths []string) (int, error) {
	if len(paths) == 0 {
		return 0, nil
	}

	backupDir := filepath.Join(m.imageStore.GetStoreDir(), WorkspaceSubdir, instanceName, BackupSubdir)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return 0, fmt.Errorf("failed to create backup directory: %w", err)
	}

	backupDataDir := filepath.Join(backupDir, "data")
	if err := os.MkdirAll(backupDataDir, 0755); err != nil {
		return 0, fmt.Errorf("failed to create backup data directory: %w", err)
	}

	for _, path := range paths {
		if err := m.backupPath(workspaceRoot, path, backupDataDir); err != nil {
			return 0, fmt.Errorf("failed to backup path %s: %w", path, err)
		}
		logger.Info("workspace:start", fmt.Sprintf("Backing up %s (conflicts with image artifact)", path))
	}

	return len(paths), nil
}