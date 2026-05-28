package cache

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// CopyArtifact copies a file or directory from source to the target directory.
// Preserves the relative path structure within the target.
func CopyArtifact(sourcePath string, targetDir string) error {
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

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
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}

	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	targetFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, sourceInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", err)
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
}

// CopyDirectory copies a directory tree from source to target.
// Preserves directory structure and file permissions.
func CopyDirectory(sourcePath string, targetPath string) error {
	targetParent := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetParent, 0755); err != nil {
		return fmt.Errorf("failed to create target parent directory: %w", err)
	}

	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to get source directory info: %w", err)
	}

	if err := os.MkdirAll(targetPath, sourceInfo.Mode()); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	return filepath.WalkDir(sourcePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		if relPath == "." {
			return nil
		}

		targetItemPath := filepath.Join(targetPath, relPath)

		if d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return fmt.Errorf("failed to get directory info: %w", err)
			}
			return os.MkdirAll(targetItemPath, info.Mode())
		}

		return CopyFile(path, targetItemPath)
	})
}

// CopyArtifacts copies multiple artifacts from a source root to a target directory.
// Each artifact path is relative to the source root.
// Returns list of successfully copied artifacts.
func CopyArtifacts(sourceRoot string, artifacts []string, targetDir string) ([]string, error) {
	var copied []string

	for _, artifact := range artifacts {
		sourcePath := filepath.Join(sourceRoot, artifact)
		if err := copyArtifactWithPath(sourcePath, artifact, targetDir); err != nil {
			return copied, fmt.Errorf("failed to copy %s: %w", artifact, err)
		}
		copied = append(copied, artifact)
	}

	return copied, nil
}

// copyArtifactWithPath copies an artifact preserving its relative path structure.
func copyArtifactWithPath(sourcePath string, relPath string, targetDir string) error {
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	targetPath := filepath.Join(targetDir, relPath)

	if sourceInfo.IsDir() {
		return CopyDirectory(sourcePath, targetPath)
	}

	return CopyFile(sourcePath, targetPath)
}

// RestoreArtifacts copies artifacts from cache to the workdir.
// Returns the list of restored artifact paths.
func RestoreArtifacts(entryPath string, workdir string, artifacts []string) ([]string, error) {
	artifactsDir := filepath.Join(entryPath, "artifacts")
	var restored []string

	for _, artifact := range artifacts {
		sourcePath := filepath.Join(artifactsDir, artifact)
		targetPath := filepath.Join(workdir, artifact)

		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			continue
		}

		// Ensure target directory exists
		targetDir := filepath.Dir(targetPath)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return restored, fmt.Errorf("failed to create directory for %s: %w", artifact, err)
		}

		sourceInfo, err := os.Stat(sourcePath)
		if err != nil {
			return restored, fmt.Errorf("failed to stat %s: %w", artifact, err)
		}

		if sourceInfo.IsDir() {
			if err := CopyDirectory(sourcePath, targetPath); err != nil {
				return restored, fmt.Errorf("failed to restore directory %s: %w", artifact, err)
			}
		} else {
			if err := CopyFile(sourcePath, targetPath); err != nil {
				return restored, fmt.Errorf("failed to restore file %s: %w", artifact, err)
			}
		}

		restored = append(restored, artifact)
	}

	return restored, nil
}
