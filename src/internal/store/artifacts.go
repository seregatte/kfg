// Package store provides named configuration entry storage.
package store

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// CopyArtifact copies a file or directory from source to the target directory.
// Preserves the relative path structure within the target.
// The artifact path (relative or absolute) is preserved as a subdirectory.
func CopyArtifact(sourcePath string, targetDir string) error {
	// Check if source exists
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	// Get the base name of the source (preserves the artifact name)
	baseName := filepath.Base(sourcePath)
	targetPath := filepath.Join(targetDir, baseName)

	if sourceInfo.IsDir() {
		return CopyDirectory(sourcePath, targetPath)
	}

	return CopyFile(sourcePath, targetPath)
}

// CopyFile copies a single file from source to target.
// Preserves file permissions.
func CopyFile(sourcePath string, targetPath string) error {
	// Open source file
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// Get source file info for permissions
	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}

	// Ensure target directory exists
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Create target file with same permissions
	targetFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, sourceInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", err)
	}
	defer targetFile.Close()

	// Copy content
	_, err = io.Copy(targetFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
}

// CopyDirectory copies a directory tree from source to target.
// Preserves directory structure and file permissions.
func CopyDirectory(sourcePath string, targetPath string) error {
	// Ensure target parent directory exists
	targetParent := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetParent, 0755); err != nil {
		return fmt.Errorf("failed to create target parent directory: %w", err)
	}

	// Get source directory info for permissions
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to get source directory info: %w", err)
	}

	// Create target directory with same permissions
	if err := os.MkdirAll(targetPath, sourceInfo.Mode()); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Walk source directory and copy each item
	err = filepath.WalkDir(sourcePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Get relative path from source
		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Skip the root directory itself
		if relPath == "." {
			return nil
		}

		// Target path
		targetItemPath := filepath.Join(targetPath, relPath)

		if d.IsDir() {
			// Get directory info for permissions
			info, err := d.Info()
			if err != nil {
				return fmt.Errorf("failed to get directory info: %w", err)
			}

			// Create directory with same permissions
			return os.MkdirAll(targetItemPath, info.Mode())
		}

		// Copy file
		return CopyFile(path, targetItemPath)
	})

	if err != nil {
		return fmt.Errorf("failed to copy directory: %w", err)
	}

	return nil
}

// CopyArtifactsFromRoot copies multiple artifacts from a root directory to target.
// Each artifact path is relative to root.
// Returns list of successfully copied artifacts and any errors.
func CopyArtifactsFromRoot(root string, artifacts []string, targetDir string) ([]string, []error) {
	copied := []string{}
	errors := []error{}

	for _, artifact := range artifacts {
		sourcePath := filepath.Join(root, artifact)
		if err := CopyArtifact(sourcePath, targetDir); err != nil {
			errors = append(errors, fmt.Errorf("failed to copy %s: %w", artifact, err))
		} else {
			copied = append(copied, artifact)
		}
	}

	return copied, errors
}

// ArtifactExists checks if an artifact exists at the given path.
func ArtifactExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GetArtifactSize returns the size of a file or directory.
func GetArtifactSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	if info.IsDir() {
		return GetDirectorySize(path)
	}

	return info.Size(), nil
}

// GetDirectorySize returns the total size of a directory tree.
func GetDirectorySize(path string) (int64, error) {
	var totalSize int64

	err := filepath.WalkDir(path, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			totalSize += info.Size()
		}

		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to calculate directory size: %w", err)
	}

	return totalSize, nil
}

// DetectFilesToBackup scans workspace and store entry to identify files that would be overwritten.
// Returns list of relative paths for files that exist in workspace and will be replaced by pull.
func DetectFilesToBackup(workspaceRoot string, storeArtifactsPath string) ([]string, error) {
	filesToBackup := []string{}

	// Walk the store artifacts directory to see what files would be pulled
	err := filepath.WalkDir(storeArtifactsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Get relative path from artifacts directory
		relPath, err := filepath.Rel(storeArtifactsPath, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Skip the root directory itself
		if relPath == "." {
			return nil
		}

		// Only check files (directories don't need backup)
		if !d.IsDir() {
			// Check if file exists in workspace
			workspacePath := filepath.Join(workspaceRoot, relPath)
			if ArtifactExists(workspacePath) {
				filesToBackup = append(filesToBackup, relPath)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to detect files to backup: %w", err)
	}

	return filesToBackup, nil
}
