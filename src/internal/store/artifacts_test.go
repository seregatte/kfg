package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFile(t *testing.T) {
	// Create temp directories
	tempDir, err := os.MkdirTemp("", "nixai-store-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create source file
	sourcePath := filepath.Join(tempDir, "source.txt")
	sourceContent := []byte("test content for file copy")
	if err := os.WriteFile(sourcePath, sourceContent, 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	// Copy file
	targetPath := filepath.Join(tempDir, "target.txt")
	if err := CopyFile(sourcePath, targetPath); err != nil {
		t.Fatalf("failed to copy file: %v", err)
	}

	// Verify content
	targetContent, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("failed to read target file: %v", err)
	}

	if string(targetContent) != string(sourceContent) {
		t.Errorf("content mismatch: got %s, want %s", targetContent, sourceContent)
	}

	// Verify permissions
	sourceInfo, _ := os.Stat(sourcePath)
	targetInfo, _ := os.Stat(targetPath)
	if sourceInfo.Mode() != targetInfo.Mode() {
		t.Errorf("mode mismatch: got %v, want %v", targetInfo.Mode(), sourceInfo.Mode())
	}
}

func TestCopyFileNonExistentSource(t *testing.T) {
	err := CopyFile("/nonexistent/file.txt", "/tmp/target.txt")
	if err == nil {
		t.Error("expected error for non-existent source")
	}
}

func TestCopyDirectory(t *testing.T) {
	// Create temp directories
	tempDir, err := os.MkdirTemp("", "nixai-store-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create source directory structure
	sourceDir := filepath.Join(tempDir, "source")
	if err := os.MkdirAll(filepath.Join(sourceDir, "subdir"), 0755); err != nil {
		t.Fatalf("failed to create source directory: %v", err)
	}

	// Create files in directory
	file1 := filepath.Join(sourceDir, "file1.txt")
	file2 := filepath.Join(sourceDir, "subdir", "file2.txt")
	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatalf("failed to create file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
		t.Fatalf("failed to create file2: %v", err)
	}

	// Copy directory
	targetDir := filepath.Join(tempDir, "target")
	if err := CopyDirectory(sourceDir, targetDir); err != nil {
		t.Fatalf("failed to copy directory: %v", err)
	}

	// Verify files exist
	targetFile1 := filepath.Join(targetDir, "file1.txt")
	targetFile2 := filepath.Join(targetDir, "subdir", "file2.txt")

	content1, err := os.ReadFile(targetFile1)
	if err != nil {
		t.Errorf("failed to read target file1: %v", err)
	}
	if string(content1) != "content1" {
		t.Errorf("file1 content mismatch")
	}

	content2, err := os.ReadFile(targetFile2)
	if err != nil {
		t.Errorf("failed to read target file2: %v", err)
	}
	if string(content2) != "content2" {
		t.Errorf("file2 content mismatch")
	}
}

func TestCopyDirectoryPreservesPermissions(t *testing.T) {
	// Create temp directories
	tempDir, err := os.MkdirTemp("", "nixai-store-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create source directory with specific permissions
	sourceDir := filepath.Join(tempDir, "source")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("failed to create source directory: %v", err)
	}

	// Create executable file
	execFile := filepath.Join(sourceDir, "script.sh")
	if err := os.WriteFile(execFile, []byte("#!/bin/bash\necho hello"), 0755); err != nil {
		t.Fatalf("failed to create executable file: %v", err)
	}

	// Copy directory
	targetDir := filepath.Join(tempDir, "target")
	if err := CopyDirectory(sourceDir, targetDir); err != nil {
		t.Fatalf("failed to copy directory: %v", err)
	}

	// Verify executable permissions preserved
	targetExecFile := filepath.Join(targetDir, "script.sh")
	targetInfo, err := os.Stat(targetExecFile)
	if err != nil {
		t.Fatalf("failed to stat target executable: %v", err)
	}

	// Check that executable bit is preserved
	if targetInfo.Mode()&0111 == 0 {
		t.Errorf("executable permission not preserved: got %v", targetInfo.Mode())
	}
}

func TestCopyArtifact(t *testing.T) {
	// Create temp directories
	tempDir, err := os.MkdirTemp("", "nixai-store-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create source file
	sourcePath := filepath.Join(tempDir, "artifact.txt")
	if err := os.WriteFile(sourcePath, []byte("artifact content"), 0644); err != nil {
		t.Fatalf("failed to create artifact: %v", err)
	}

	// Create target directory
	targetDir := filepath.Join(tempDir, "store", "artifacts")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("failed to create target directory: %v", err)
	}

	// Copy artifact
	if err := CopyArtifact(sourcePath, targetDir); err != nil {
		t.Fatalf("failed to copy artifact: %v", err)
	}

	// Verify file exists with correct name (basename preserved)
	targetPath := filepath.Join(targetDir, "artifact.txt")
	content, err := os.ReadFile(targetPath)
	if err != nil {
		t.Errorf("artifact not found at expected path: %v", err)
	}
	if string(content) != "artifact content" {
		t.Errorf("artifact content mismatch")
	}
}

func TestCopyArtifactDirectory(t *testing.T) {
	// Create temp directories
	tempDir, err := os.MkdirTemp("", "nixai-store-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create artifact directory
	artifactDir := filepath.Join(tempDir, "output")
	if err := os.MkdirAll(filepath.Join(artifactDir, "nested"), 0755); err != nil {
		t.Fatalf("failed to create artifact directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(artifactDir, "file.txt"), []byte("content"), 0644); err != nil {
		t.Fatalf("failed to create file in artifact: %v", err)
	}

	// Create target directory
	targetDir := filepath.Join(tempDir, "store", "artifacts")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("failed to create target directory: %v", err)
	}

	// Copy artifact (directory)
	if err := CopyArtifact(artifactDir, targetDir); err != nil {
		t.Fatalf("failed to copy artifact directory: %v", err)
	}

	// Verify directory exists
	targetArtifactDir := filepath.Join(targetDir, "output")
	if _, err := os.Stat(targetArtifactDir); err != nil {
		t.Errorf("artifact directory not found: %v", err)
	}

	// Verify nested structure
	targetFile := filepath.Join(targetArtifactDir, "file.txt")
	content, err := os.ReadFile(targetFile)
	if err != nil {
		t.Errorf("nested file not found: %v", err)
	}
	if string(content) != "content" {
		t.Errorf("content mismatch")
	}
}

func TestCopyArtifactsFromRoot(t *testing.T) {
	// Create temp directories
	tempDir, err := os.MkdirTemp("", "nixai-store-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create root directory with artifacts
	rootDir := filepath.Join(tempDir, "project")
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		t.Fatalf("failed to create root directory: %v", err)
	}

	// Create artifacts
	if err := os.WriteFile(filepath.Join(rootDir, "app"), []byte("app content"), 0644); err != nil {
		t.Fatalf("failed to create app: %v", err)
	}
	if err := os.WriteFile(filepath.Join(rootDir, "config.yaml"), []byte("config content"), 0644); err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	// Create target directory
	targetDir := filepath.Join(tempDir, "store", "artifacts")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("failed to create target directory: %v", err)
	}

	// Copy multiple artifacts
	artifacts := []string{"app", "config.yaml"}
	copied, errors := CopyArtifactsFromRoot(rootDir, artifacts, targetDir)

	// Verify all copied
	if len(copied) != 2 {
		t.Errorf("expected 2 artifacts copied, got %d", len(copied))
	}
	if len(errors) != 0 {
		t.Errorf("unexpected errors: %v", errors)
	}

	// Verify files exist
	for _, artifact := range artifacts {
		targetPath := filepath.Join(targetDir, artifact)
		if _, err := os.Stat(targetPath); err != nil {
			t.Errorf("artifact %s not found: %v", artifact, err)
		}
	}
}

func TestCopyArtifactsFromRootWithMissing(t *testing.T) {
	// Create temp directories
	tempDir, err := os.MkdirTemp("", "nixai-store-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create root directory with one artifact
	rootDir := filepath.Join(tempDir, "project")
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		t.Fatalf("failed to create root directory: %v", err)
	}

	// Create one artifact, one missing
	if err := os.WriteFile(filepath.Join(rootDir, "app"), []byte("app content"), 0644); err != nil {
		t.Fatalf("failed to create app: %v", err)
	}

	// Create target directory
	targetDir := filepath.Join(tempDir, "store", "artifacts")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("failed to create target directory: %v", err)
	}

	// Try to copy both (one missing)
	artifacts := []string{"app", "missing"}
	copied, errors := CopyArtifactsFromRoot(rootDir, artifacts, targetDir)

	// Verify one copied, one error
	if len(copied) != 1 {
		t.Errorf("expected 1 artifact copied, got %d", len(copied))
	}
	if len(errors) != 1 {
		t.Errorf("expected 1 error for missing artifact, got %d", len(errors))
	}
	if copied[0] != "app" {
		t.Errorf("expected 'app' copied, got %s", copied[0])
	}
}

func TestArtifactExists(t *testing.T) {
	// Create temp file
	tempFile, err := os.CreateTemp("", "nixai-test")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempPath)

	// Test existing file
	if !ArtifactExists(tempPath) {
		t.Error("expected ArtifactExists to return true for existing file")
	}

	// Test non-existent file
	if ArtifactExists("/nonexistent/file.txt") {
		t.Error("expected ArtifactExists to return false for non-existent file")
	}
}

func TestGetArtifactSize(t *testing.T) {
	// Create temp file
	tempDir, err := os.MkdirTemp("", "nixai-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create file with known size
	filePath := filepath.Join(tempDir, "test.txt")
	content := []byte("test content 1234567890") // 20 bytes
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	// Get size
	size, err := GetArtifactSize(filePath)
	if err != nil {
		t.Fatalf("failed to get artifact size: %v", err)
	}

	if size != int64(len(content)) {
		t.Errorf("size mismatch: got %d, want %d", size, len(content))
	}
}

func TestGetDirectorySize(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nixai-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create files with known sizes
	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(tempDir, "file2.txt")
	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatalf("failed to create file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("content22"), 0644); err != nil {
		t.Fatalf("failed to create file2: %v", err)
	}

	// Expected size: 8 + 9 = 17 bytes
	expectedSize := int64(8 + 9)

	// Get directory size
	size, err := GetDirectorySize(tempDir)
	if err != nil {
		t.Fatalf("failed to get directory size: %v", err)
	}

	if size != expectedSize {
		t.Errorf("size mismatch: got %d, want %d", size, expectedSize)
	}
}