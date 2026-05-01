package image

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewMaterializer(t *testing.T) {
	materializer := NewMaterializer("")

	if materializer == nil {
		t.Fatal("expected materializer to be created")
	}
}

func TestStartEmptyWorkspace(t *testing.T) {
	// Test starting in empty workspace (no backup created)
	tmpDir := t.TempDir()
	storeDir := filepath.Join(tmpDir, "store")
	workspaceDir := filepath.Join(tmpDir, "workspace")

	// Create store and image
	store := NewImageStore(storeDir)
	store.Initialize()

	candidateDir := filepath.Join(tmpDir, "candidate")
	os.MkdirAll(candidateDir, 0755)

	metadata := NewMetadata("start-test", "v1")
	metadata.SetDigest("sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	metadata.SetRecipe("./Imagefile", "FROM scratch\nCOPY test.txt ./\nTAG start-test:v1")
	metadata.AddFile("test.txt", "workspace")
	metadata.SaveMetadataToDir(candidateDir)

	// Add a file to candidate
	testFile := filepath.Join(candidateDir, "test.txt")
	os.WriteFile(testFile, []byte("test content"), 0644)

	store.PushImage(candidateDir, false)

	// Create empty workspace
	os.MkdirAll(workspaceDir, 0755)

	// Start image in empty workspace
	materializer := NewMaterializer(storeDir)
	err := materializer.Start("start-test:v1", workspaceDir, "test-instance")

	if err != nil {
		t.Fatalf("start failed: %v", err)
	}

	// Verify file was materialized
	materializedFile := filepath.Join(workspaceDir, "test.txt")
	if _, err := os.Stat(materializedFile); err != nil {
		t.Errorf("expected file to be materialized: %v", err)
	}

	// Verify instance metadata exists
	instanceDir := filepath.Join(storeDir, ".workspace", "test-instance")
	instanceFile := filepath.Join(instanceDir, "instance.json")
	if _, err := os.Stat(instanceFile); err != nil {
		t.Errorf("instance metadata not created: %v", err)
	}
}

func TestStartWithExistingFiles(t *testing.T) {
	// Test starting in workspace with conflicting files (scoped backup created)
	tmpDir := t.TempDir()
	storeDir := filepath.Join(tmpDir, "store")
	workspaceDir := filepath.Join(tmpDir, "workspace")

	// Create store and image
	store := NewImageStore(storeDir)
	store.Initialize()

	candidateDir := filepath.Join(tmpDir, "candidate")
	os.MkdirAll(candidateDir, 0755)

	metadata := NewMetadata("backup-test", "v1")
	metadata.SetDigest("sha256:0000000000000000000000000000000000000000000000000000000000000000")
	metadata.SetRecipe("./Imagefile", "FROM scratch\nCOPY new.txt ./\nTAG backup-test:v1")
	metadata.AddFile("new.txt", "workspace") // Track file in manifest
	metadata.SaveMetadataToDir(candidateDir)

	os.WriteFile(filepath.Join(candidateDir, "new.txt"), []byte("new content"), 0644)
	store.PushImage(candidateDir, false)

	// Create workspace with existing file that CONFLICTS with image
	os.MkdirAll(workspaceDir, 0755)
	os.WriteFile(filepath.Join(workspaceDir, "new.txt"), []byte("existing content"), 0644) // Conflict!
	os.WriteFile(filepath.Join(workspaceDir, "unrelated.txt"), []byte("unrelated"), 0644) // No conflict

	// Start image (should create scoped backup for conflicting file only)
	materializer := NewMaterializer(storeDir)
	err := materializer.Start("backup-test:v1", workspaceDir, "backup-instance")

	if err != nil {
		t.Fatalf("start with backup failed: %v", err)
	}

	// Verify backup was created for conflicting file
	backupDir := filepath.Join(storeDir, ".workspace", "backup-instance", "backup")
	backupDataDir := filepath.Join(backupDir, "data")
	if _, err := os.Stat(backupDataDir); err != nil {
		t.Errorf("backup not created: %v", err)
	}

	// Verify only conflicting file was backed up (new.txt)
	backedUpFile := filepath.Join(backupDataDir, "new.txt")
	if _, err := os.Stat(backedUpFile); err != nil {
		t.Errorf("conflicting file not backed up: %v", err)
	}

	// Verify unrelated file was NOT backed up (scoped backup)
	unrelatedBackup := filepath.Join(backupDataDir, "unrelated.txt")
	if _, err := os.Stat(unrelatedBackup); err == nil {
		t.Error("unrelated file should not be backed up (scoped backup)")
	}
}

func TestStopWithBackup(t *testing.T) {
	// Test stopping and restoring from backup
	tmpDir := t.TempDir()
	storeDir := filepath.Join(tmpDir, "store")
	workspaceDir := filepath.Join(tmpDir, "workspace")

	// Setup store and workspace
	store := NewImageStore(storeDir)
	store.Initialize()

	candidateDir := filepath.Join(tmpDir, "candidate")
	os.MkdirAll(candidateDir, 0755)

	metadata := NewMetadata("stop-test", "v1")
	metadata.SetDigest("sha256:1111111111111111111111111111111111111111111111111111111111111111")
	metadata.SetRecipe("./Imagefile", "FROM scratch\nCOPY original.txt ./\nTAG stop-test:v1")
	metadata.AddFile("original.txt", "workspace") // Track file in manifest
	metadata.SaveMetadataToDir(candidateDir)

	os.WriteFile(filepath.Join(candidateDir, "original.txt"), []byte("image content"), 0644)
	store.PushImage(candidateDir, false)

	// Create workspace with existing file that CONFLICTS with image
	os.MkdirAll(workspaceDir, 0755)
	originalContent := []byte("original content")
	os.WriteFile(filepath.Join(workspaceDir, "original.txt"), originalContent, 0644) // Conflict!
	os.WriteFile(filepath.Join(workspaceDir, "unrelated.txt"), []byte("unrelated content"), 0644) // No conflict

	// Start instance (creates backup for conflicting file)
	materializer := NewMaterializer(storeDir)
	materializer.Start("stop-test:v1", workspaceDir, "stop-instance")

	// Verify original.txt was overwritten by image content
	content, _ := os.ReadFile(filepath.Join(workspaceDir, "original.txt"))
	if string(content) != "image content" {
		t.Fatalf("expected 'image content' after start, got '%s'", string(content))
	}

	// Stop instance (should restore original.txt)
	err := materializer.Stop("stop-instance")
	if err != nil {
		t.Fatalf("stop failed: %v", err)
	}

	// Verify original content was restored
	restoredFile := filepath.Join(workspaceDir, "original.txt")
	content, err = os.ReadFile(restoredFile)
	if err != nil {
		t.Fatalf("failed to read restored file: %v", err)
	}

	if string(content) != "original content" {
		t.Errorf("expected 'original content', got '%s'", string(content))
	}

	// Verify unrelated file remains untouched
	unrelatedContent, err := os.ReadFile(filepath.Join(workspaceDir, "unrelated.txt"))
	if err != nil {
		t.Fatalf("failed to read unrelated file: %v", err)
	}
	if string(unrelatedContent) != "unrelated content" {
		t.Errorf("unrelated file should be preserved: got '%s'", string(unrelatedContent))
	}

	// Verify instance is cleaned up
	instanceDir := filepath.Join(storeDir, ".workspace", "stop-instance")
	if _, err := os.Stat(instanceDir); err == nil {
		t.Error("instance directory should be removed after stop")
	}
}

func TestStopNoBackup(t *testing.T) {
	// Test stopping when no backup exists (idempotent)
	tmpDir := t.TempDir()
	storeDir := filepath.Join(tmpDir, "store")

	materializer := NewMaterializer(storeDir)

	// Stop instance that never existed (should succeed)
	err := materializer.Stop("nonexistent-instance")
	if err != nil {
		t.Errorf("stop should succeed for nonexistent instance (idempotent): %v", err)
	}
}

func TestDuplicateInstanceName(t *testing.T) {
	// Test that duplicate instance names fail
	tmpDir := t.TempDir()
	storeDir := filepath.Join(tmpDir, "store")
	workspaceDir := filepath.Join(tmpDir, "workspace")

	store := NewImageStore(storeDir)
	store.Initialize()

	candidateDir := filepath.Join(tmpDir, "candidate")
	os.MkdirAll(candidateDir, 0755)

	metadata := NewMetadata("duplicate-test", "v1")
	metadata.SetDigest("sha256:2222222222222222222222222222222222222222222222222222222222222222")
	metadata.SetRecipe("./Imagefile", "FROM scratch\nTAG duplicate-test:v1")
	metadata.SaveMetadataToDir(candidateDir)

	store.PushImage(candidateDir, false)

	os.MkdirAll(workspaceDir, 0755)

	materializer := NewMaterializer(storeDir)

	// First start should succeed
	err := materializer.Start("duplicate-test:v1", workspaceDir, "duplicate-name")
	if err != nil {
		t.Fatalf("first start failed: %v", err)
	}

	// Second start with same name should fail
	err = materializer.Start("duplicate-test:v1", workspaceDir, "duplicate-name")
	if err == nil {
		t.Error("expected second start with duplicate name to fail")
	}
}

func TestMissingImageRef(t *testing.T) {
	// Test starting with image that doesn't exist
	tmpDir := t.TempDir()
	storeDir := filepath.Join(tmpDir, "store")
	workspaceDir := filepath.Join(tmpDir, "workspace")

	store := NewImageStore(storeDir)
	store.Initialize()

	os.MkdirAll(workspaceDir, 0755)

	materializer := NewMaterializer(storeDir)

	err := materializer.Start("nonexistent:v1", workspaceDir, "test-instance")
	if err == nil {
		t.Error("expected error when starting non-existent image")
	}
}

func TestInstanceNameRequired(t *testing.T) {
	// Test that instance name is required
	tmpDir := t.TempDir()

	materializer := NewMaterializer(tmpDir)

	err := materializer.Start("test:v1", tmpDir, "")
	if err == nil {
		t.Error("expected error when instance name is empty")
	}
}

func TestInstanceMetadataWithMaterializedPaths(t *testing.T) {
	// Test that MaterializedPaths is serialized correctly
	tmpDir := t.TempDir()
	storeDir := filepath.Join(tmpDir, "store")
	os.MkdirAll(storeDir, 0755)

	materializer := NewMaterializer(storeDir)

	// Create instance metadata with materialized paths
	instance := InstanceMetadata{
		Name:              "test-instance",
		ImageRef:          "test:v1",
		WorkspaceRoot:     "/workspace",
		StartedAt:         "2024-01-01T00:00:00Z",
		ImageDigest:       "sha256:abcdef",
		MaterializedPaths: []string{"CLAUDE.md", ".pi/config.json"},
	}

	// Save instance
	instanceFile := filepath.Join(storeDir, "instance.json")
	err := materializer.saveInstance(instance, instanceFile)
	if err != nil {
		t.Fatalf("failed to save instance: %v", err)
	}

	// Load instance and verify
	loaded, err := materializer.loadInstance(instanceFile)
	if err != nil {
		t.Fatalf("failed to load instance: %v", err)
	}

	if len(loaded.MaterializedPaths) != 2 {
		t.Errorf("expected 2 materialized paths, got %d", len(loaded.MaterializedPaths))
	}

	if loaded.MaterializedPaths[0] != "CLAUDE.md" {
		t.Errorf("expected first path to be CLAUDE.md, got %s", loaded.MaterializedPaths[0])
	}

	if loaded.MaterializedPaths[1] != ".pi/config.json" {
		t.Errorf("expected second path to be .pi/config.json, got %s", loaded.MaterializedPaths[1])
	}
}

func TestInstanceMetadataBackwardCompatibility(t *testing.T) {
	// Test loading old metadata without MaterializedPaths field
	tmpDir := t.TempDir()
	storeDir := filepath.Join(tmpDir, "store")
	os.MkdirAll(storeDir, 0755)

	materializer := NewMaterializer(storeDir)

	// Create old-style metadata without materialized_paths
	oldMetadata := `{"name":"old-instance","image_ref":"old:v1","workspace_root":"/workspace","started_at":"2024-01-01T00:00:00Z","image_digest":"sha256:old"}`
	instanceFile := filepath.Join(storeDir, "instance.json")
	if err := os.WriteFile(instanceFile, []byte(oldMetadata), 0644); err != nil {
		t.Fatalf("failed to write old metadata: %v", err)
	}

	// Load instance and verify MaterializedPaths defaults to empty slice
	loaded, err := materializer.loadInstance(instanceFile)
	if err != nil {
		t.Fatalf("failed to load old instance metadata: %v", err)
	}

	if loaded.MaterializedPaths == nil {
		t.Error("MaterializedPaths should not be nil for backward compatibility")
	}

	if len(loaded.MaterializedPaths) != 0 {
		t.Errorf("expected empty MaterializedPaths, got %d items", len(loaded.MaterializedPaths))
	}
}

func TestInstanceMetadataEmptyPathsSerialization(t *testing.T) {
	// Test that empty MaterializedPaths slice serializes as [] not null
	tmpDir := t.TempDir()
	storeDir := filepath.Join(tmpDir, "store")
	os.MkdirAll(storeDir, 0755)

	materializer := NewMaterializer(storeDir)

	// Create instance with empty materialized paths
	instance := InstanceMetadata{
		Name:              "empty-paths-instance",
		ImageRef:          "test:v1",
		WorkspaceRoot:     "/workspace",
		StartedAt:         "2024-01-01T00:00:00Z",
		ImageDigest:       "sha256:empty",
		MaterializedPaths: []string{}, // Empty slice, not nil
	}

	// Save instance
	instanceFile := filepath.Join(storeDir, "instance.json")
	err := materializer.saveInstance(instance, instanceFile)
	if err != nil {
		t.Fatalf("failed to save instance: %v", err)
	}

	// Read raw file and verify it contains empty array
	data, err := os.ReadFile(instanceFile)
	if err != nil {
		t.Fatalf("failed to read instance file: %v", err)
	}

	// Check that materialized_paths is serialized as [] not null
	if !containsString(string(data), `"materialized_paths": []`) {
		t.Errorf("expected materialized_paths to be empty array, got: %s", string(data))
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStringHelper(s, substr))
}

func containsStringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
// --- Section 2: Scoped Backup Unit Tests (Tasks 2.8-2.11) ---

func TestComputeArtifactPaths(t *testing.T) {
	// Task 2.8: Test computeArtifactPaths from image metadata
	tmpDir := t.TempDir()
	storeDir := filepath.Join(tmpDir, "store")
	store := NewImageStore(storeDir)
	store.Initialize()

	candidateDir := filepath.Join(tmpDir, "candidate")
	os.MkdirAll(candidateDir, 0755)

	metadata := NewMetadata("artifact-path-test", "v1")
	metadata.SetDigest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	metadata.SetRecipe("./Imagefile", "FROM scratch\nCOPY CLAUDE.md ./\nCOPY .pi/config.json ./\nCOPY README.md ./\nTAG artifact-path-test:v1")
	metadata.AddFile("CLAUDE.md", "workspace")
	metadata.AddFile(".pi/config.json", "workspace")
	metadata.AddFile("README.md", "workspace")
	metadata.SaveMetadataToDir(candidateDir)

	os.WriteFile(filepath.Join(candidateDir, "CLAUDE.md"), []byte("claude content"), 0644)
	os.MkdirAll(filepath.Join(candidateDir, ".pi"), 0755)
	os.WriteFile(filepath.Join(candidateDir, ".pi", "config.json"), []byte("config content"), 0644)
	os.WriteFile(filepath.Join(candidateDir, "README.md"), []byte("readme content"), 0644)

	store.PushImage(candidateDir, false)

	m := NewMaterializer(storeDir)
	_, imageDir, err := store.LoadImage("artifact-path-test:v1")
	if err != nil {
		t.Fatalf("failed to load image: %v", err)
	}

	paths, err := m.computeArtifactPaths(imageDir)
	if err != nil {
		t.Fatalf("failed to compute artifact paths: %v", err)
	}

	if len(paths) != 3 {
		t.Errorf("expected 3 artifact paths, got %d", len(paths))
	}

	expectedPaths := []string{"CLAUDE.md", ".pi/config.json", "README.md"}
	for _, expected := range expectedPaths {
		found := false
		for _, path := range paths {
			if path == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected path %s not found in artifact paths", expected)
		}
	}
}

func TestFindConflictingPathsPartial(t *testing.T) {
	// Task 2.9: Test findConflictingPaths with partial conflicts
	tmpDir := t.TempDir()
	workspaceDir := filepath.Join(tmpDir, "workspace")
	os.MkdirAll(workspaceDir, 0755)

	os.WriteFile(filepath.Join(workspaceDir, "CLAUDE.md"), []byte("existing claude"), 0644)
	os.WriteFile(filepath.Join(workspaceDir, "README.md"), []byte("existing readme"), 0644)

	m := NewMaterializer(tmpDir)
	artifactPaths := []string{"CLAUDE.md", ".pi/config.json", "README.md"}

	conflicts, err := m.findConflictingPaths(workspaceDir, artifactPaths)
	if err != nil {
		t.Fatalf("failed to find conflicting paths: %v", err)
	}

	if len(conflicts) != 2 {
		t.Errorf("expected 2 conflicting paths, got %d", len(conflicts))
	}
}

func TestScopedBackupPreservingDirectoryStructure(t *testing.T) {
	// Task 2.10: Test scoped backup preserving directory structure
	tmpDir := t.TempDir()
	storeDir := filepath.Join(tmpDir, "store")
	store := NewImageStore(storeDir)
	store.Initialize()

	workspaceDir := filepath.Join(tmpDir, "workspace")
	os.MkdirAll(workspaceDir, 0755)
	os.MkdirAll(filepath.Join(workspaceDir, ".pi"), 0755)
	os.WriteFile(filepath.Join(workspaceDir, ".pi", "config.json"), []byte("existing config"), 0644)

	m := NewMaterializer(storeDir)
	conflictingPaths := []string{".pi/config.json"}
	count, err := m.createScopedBackup(workspaceDir, "test-backup", conflictingPaths)
	if err != nil {
		t.Fatalf("scoped backup failed: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 backed up path, got %d", count)
	}

	backupDir := filepath.Join(storeDir, ".workspace", "test-backup", "backup", "data")
	backedUpFile := filepath.Join(backupDir, ".pi", "config.json")
	if _, err := os.Stat(backedUpFile); err != nil {
		t.Errorf("backup file not found at correct path: %v", err)
	}

	content, err := os.ReadFile(backedUpFile)
	if err != nil {
		t.Fatalf("failed to read backed up file: %v", err)
	}
	if string(content) != "existing config" {
		t.Errorf("backup content mismatch: expected 'existing config', got '%s'", string(content))
	}
}

func TestBackupSkipWhenNoConflicts(t *testing.T) {
	// Task 2.11: Test backup skip when no conflicts
	tmpDir := t.TempDir()
	storeDir := filepath.Join(tmpDir, "store")
	store := NewImageStore(storeDir)
	store.Initialize()

	workspaceDir := filepath.Join(tmpDir, "workspace")
	os.MkdirAll(workspaceDir, 0755)
	os.WriteFile(filepath.Join(workspaceDir, "unrelated.txt"), []byte("unrelated content"), 0644)

	m := NewMaterializer(storeDir)
	conflictingPaths := []string{}
	count, err := m.createScopedBackup(workspaceDir, "test-skip-backup", conflictingPaths)
	if err != nil {
		t.Fatalf("scoped backup failed: %v", err)
	}

	if count != 0 {
		t.Errorf("expected 0 backed up paths, got %d", count)
	}

	backupDir := filepath.Join(storeDir, ".workspace", "test-skip-backup", "backup", "data")
	if _, err := os.Stat(backupDir); err == nil {
		t.Error("backup directory should not exist when no conflicts")
	}
}

// --- Section 3: Materialized Paths Tracking Tests (Tasks 3.4-3.6) ---

func TestMaterializedPathsTrackingDuringStart(t *testing.T) {
	// Task 3.4: Test materialized paths tracking during start
	tmpDir := t.TempDir()
	storeDir := filepath.Join(tmpDir, "store")
	store := NewImageStore(storeDir)
	store.Initialize()

	workspaceDir := filepath.Join(tmpDir, "workspace")
	os.MkdirAll(workspaceDir, 0755)

	candidateDir := filepath.Join(tmpDir, "candidate")
	os.MkdirAll(candidateDir, 0755)

	metadata := NewMetadata("tracking-test", "v1")
	metadata.SetDigest("sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	metadata.SetRecipe("./Imagefile", "FROM scratch\nCOPY CLAUDE.md ./\nCOPY GEMINI.md ./\nTAG tracking-test:v1")
	metadata.AddFile("CLAUDE.md", "workspace")
	metadata.AddFile("GEMINI.md", "workspace")
	metadata.SaveMetadataToDir(candidateDir)

	os.WriteFile(filepath.Join(candidateDir, "CLAUDE.md"), []byte("claude content"), 0644)
	os.WriteFile(filepath.Join(candidateDir, "GEMINI.md"), []byte("gemini content"), 0644)

	store.PushImage(candidateDir, false)

	m := NewMaterializer(storeDir)
	err := m.Start("tracking-test:v1", workspaceDir, "tracking-instance")
	if err != nil {
		t.Fatalf("start failed: %v", err)
	}

	instanceFile := filepath.Join(storeDir, ".workspace", "tracking-instance", "instance.json")
	instance, err := m.loadInstance(instanceFile)
	if err != nil {
		t.Fatalf("failed to load instance: %v", err)
	}

	if len(instance.MaterializedPaths) != 2 {
		t.Errorf("expected 2 materialized paths, got %d", len(instance.MaterializedPaths))
	}
}

func TestNestedPathTracking(t *testing.T) {
	// Task 3.5: Test nested path tracking
	tmpDir := t.TempDir()
	storeDir := filepath.Join(tmpDir, "store")
	store := NewImageStore(storeDir)
	store.Initialize()

	workspaceDir := filepath.Join(tmpDir, "workspace")
	os.MkdirAll(workspaceDir, 0755)

	candidateDir := filepath.Join(tmpDir, "candidate")
	os.MkdirAll(filepath.Join(candidateDir, ".pi", "subdir"), 0755)

	metadata := NewMetadata("nested-test", "v1")
	metadata.SetDigest("sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc")
	metadata.SetRecipe("./Imagefile", "FROM scratch\nCOPY .pi/subdir/file.txt ./\nTAG nested-test:v1")
	metadata.AddFile(".pi/subdir/file.txt", "workspace")
	metadata.SaveMetadataToDir(candidateDir)

	os.WriteFile(filepath.Join(candidateDir, ".pi", "subdir", "file.txt"), []byte("nested content"), 0644)

	store.PushImage(candidateDir, false)

	m := NewMaterializer(storeDir)
	err := m.Start("nested-test:v1", workspaceDir, "nested-instance")
	if err != nil {
		t.Fatalf("start failed: %v", err)
	}

	instanceFile := filepath.Join(storeDir, ".workspace", "nested-instance", "instance.json")
	instance, err := m.loadInstance(instanceFile)
	if err != nil {
		t.Fatalf("failed to load instance: %v", err)
	}

	if len(instance.MaterializedPaths) != 1 {
		t.Errorf("expected 1 materialized path, got %d", len(instance.MaterializedPaths))
	}

	if instance.MaterializedPaths[0] != ".pi/subdir/file.txt" {
		t.Errorf("expected nested path '.pi/subdir/file.txt', got '%s'", instance.MaterializedPaths[0])
	}
}

func TestDirectoryVsFilePathDistinction(t *testing.T) {
	// Task 3.6: Test directory vs file path distinction
	tmpDir := t.TempDir()
	storeDir := filepath.Join(tmpDir, "store")
	store := NewImageStore(storeDir)
	store.Initialize()

	workspaceDir := filepath.Join(tmpDir, "workspace")
	os.MkdirAll(workspaceDir, 0755)

	candidateDir := filepath.Join(tmpDir, "candidate")
	os.MkdirAll(filepath.Join(candidateDir, ".pi"), 0755)

	metadata := NewMetadata("file-vs-dir-test", "v1")
	metadata.SetDigest("sha256:dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd")
	metadata.SetRecipe("./Imagefile", "FROM scratch\nCOPY .pi/config.json ./\nTAG file-vs-dir-test:v1")
	metadata.AddFile(".pi/config.json", "workspace")
	metadata.SaveMetadataToDir(candidateDir)

	os.WriteFile(filepath.Join(candidateDir, ".pi", "config.json"), []byte("config content"), 0644)

	store.PushImage(candidateDir, false)

	m := NewMaterializer(storeDir)
	err := m.Start("file-vs-dir-test:v1", workspaceDir, "file-vs-dir-instance")
	if err != nil {
		t.Fatalf("start failed: %v", err)
	}

	instanceFile := filepath.Join(storeDir, ".workspace", "file-vs-dir-instance", "instance.json")
	instance, err := m.loadInstance(instanceFile)
	if err != nil {
		t.Fatalf("failed to load instance: %v", err)
	}

	for _, path := range instance.MaterializedPaths {
		if path == ".pi" {
			t.Error("should not track directory '.pi', only files inside it")
		}
	}
}

// --- Section 4: Scoped Cleanup Tests (Tasks 4.8-4.12) ---

func TestScopedCleanupRemovingOnlyMaterializedPaths(t *testing.T) {
	// Task 4.8: Test scoped cleanup removing only materialized paths
	tmpDir := t.TempDir()
	storeDir := filepath.Join(tmpDir, "store")
	store := NewImageStore(storeDir)
	store.Initialize()

	workspaceDir := filepath.Join(tmpDir, "workspace")
	os.MkdirAll(workspaceDir, 0755)

	candidateDir := filepath.Join(tmpDir, "candidate")
	os.MkdirAll(candidateDir, 0755)

	metadata := NewMetadata("cleanup-test", "v1")
	metadata.SetDigest("sha256:eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")
	metadata.SetRecipe("./Imagefile", "FROM scratch\nCOPY CLAUDE.md ./\nTAG cleanup-test:v1")
	metadata.AddFile("CLAUDE.md", "workspace")
	metadata.SaveMetadataToDir(candidateDir)

	os.WriteFile(filepath.Join(candidateDir, "CLAUDE.md"), []byte("claude content"), 0644)
	store.PushImage(candidateDir, false)

	os.WriteFile(filepath.Join(workspaceDir, "README.md"), []byte("unrelated readme"), 0644)

	m := NewMaterializer(storeDir)
	err := m.Start("cleanup-test:v1", workspaceDir, "cleanup-instance")
	if err != nil {
		t.Fatalf("start failed: %v", err)
	}

	err = m.Stop("cleanup-instance")
	if err != nil {
		t.Fatalf("stop failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(workspaceDir, "CLAUDE.md")); err == nil {
		t.Error("CLAUDE.md should be removed after stop")
	}
	if _, err := os.Stat(filepath.Join(workspaceDir, "README.md")); err != nil {
		t.Error("README.md should remain after stop (it was not materialized)")
	}
}

func TestCleanupIdempotency(t *testing.T) {
	// Task 4.9: Test cleanup idempotency
	tmpDir := t.TempDir()
	storeDir := filepath.Join(tmpDir, "store")
	store := NewImageStore(storeDir)
	store.Initialize()

	workspaceDir := filepath.Join(tmpDir, "workspace")
	os.MkdirAll(workspaceDir, 0755)

	candidateDir := filepath.Join(tmpDir, "candidate")
	os.MkdirAll(candidateDir, 0755)

	metadata := NewMetadata("idempotent-test", "v1")
	metadata.SetDigest("sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	metadata.SetRecipe("./Imagefile", "FROM scratch\nCOPY test.txt ./\nTAG idempotent-test:v1")
	metadata.AddFile("test.txt", "workspace")
	metadata.SaveMetadataToDir(candidateDir)

	os.WriteFile(filepath.Join(candidateDir, "test.txt"), []byte("test content"), 0644)
	store.PushImage(candidateDir, false)

	m := NewMaterializer(storeDir)
	err := m.Start("idempotent-test:v1", workspaceDir, "idempotent-instance")
	if err != nil {
		t.Fatalf("start failed: %v", err)
	}

	os.Remove(filepath.Join(workspaceDir, "test.txt"))

	err = m.Stop("idempotent-instance")
	if err != nil {
		t.Fatalf("stop should succeed even with missing path: %v", err)
	}

	instanceDir := filepath.Join(storeDir, ".workspace", "idempotent-instance")
	if _, err := os.Stat(instanceDir); err == nil {
		t.Error("instance directory should be removed after stop")
	}
}

func TestEmptyDirectoryRemovalAfterCleanup(t *testing.T) {
	// Task 4.11: Test empty directory removal after cleanup
	tmpDir := t.TempDir()
	storeDir := filepath.Join(tmpDir, "store")
	store := NewImageStore(storeDir)
	store.Initialize()

	workspaceDir := filepath.Join(tmpDir, "workspace")
	os.MkdirAll(workspaceDir, 0755)

	candidateDir := filepath.Join(tmpDir, "candidate")
	os.MkdirAll(filepath.Join(candidateDir, ".pi"), 0755)

	metadata := NewMetadata("empty-dir-test", "v1")
	metadata.SetDigest("sha256:111111111111111111111111111111111111111111111111111111111111111a")
	metadata.SetRecipe("./Imagefile", "FROM scratch\nCOPY .pi/config.json ./\nTAG empty-dir-test:v1")
	metadata.AddFile(".pi/config.json", "workspace")
	metadata.SaveMetadataToDir(candidateDir)

	os.WriteFile(filepath.Join(candidateDir, ".pi", "config.json"), []byte("config content"), 0644)
	store.PushImage(candidateDir, false)

	m := NewMaterializer(storeDir)
	err := m.Start("empty-dir-test:v1", workspaceDir, "empty-dir-instance")
	if err != nil {
		t.Fatalf("start failed: %v", err)
	}

	err = m.Stop("empty-dir-instance")
	if err != nil {
		t.Fatalf("stop failed: %v", err)
	}

	piDir := filepath.Join(workspaceDir, ".pi")
	if _, err := os.Stat(piDir); err == nil {
		t.Error("empty .pi directory should be removed after cleanup")
	}
}

func TestPreserveNonMaterializedFilesDuringCleanup(t *testing.T) {
	// Task 4.12: Test preserving non-materialized files during cleanup
	tmpDir := t.TempDir()
	storeDir := filepath.Join(tmpDir, "store")
	store := NewImageStore(storeDir)
	store.Initialize()

	workspaceDir := filepath.Join(tmpDir, "workspace")
	os.MkdirAll(workspaceDir, 0755)
	os.MkdirAll(filepath.Join(workspaceDir, ".pi"), 0755)

	candidateDir := filepath.Join(tmpDir, "candidate")
	os.MkdirAll(filepath.Join(candidateDir, ".pi"), 0755)

	metadata := NewMetadata("preserve-test", "v1")
	metadata.SetDigest("sha256:111111111111111111111111111111111111111111111111111111111111111b")
	metadata.SetRecipe("./Imagefile", "FROM scratch\nCOPY .pi/config.json ./\nTAG preserve-test:v1")
	metadata.AddFile(".pi/config.json", "workspace")
	metadata.SaveMetadataToDir(candidateDir)

	os.WriteFile(filepath.Join(candidateDir, ".pi", "config.json"), []byte("config content"), 0644)
	store.PushImage(candidateDir, false)

	os.WriteFile(filepath.Join(workspaceDir, ".pi", "notes.txt"), []byte("user notes"), 0644)

	m := NewMaterializer(storeDir)
	err := m.Start("preserve-test:v1", workspaceDir, "preserve-instance")
	if err != nil {
		t.Fatalf("start failed: %v", err)
	}

	err = m.Stop("preserve-instance")
	if err != nil {
		t.Fatalf("stop failed: %v", err)
	}

	configFile := filepath.Join(workspaceDir, ".pi", "config.json")
	if _, err := os.Stat(configFile); err == nil {
		t.Error("config.json should be removed")
	}

	notesFile := filepath.Join(workspaceDir, ".pi", "notes.txt")
	if _, err := os.Stat(notesFile); err != nil {
		t.Error("notes.txt should be preserved (not materialized from image)")
	}

	piDir := filepath.Join(workspaceDir, ".pi")
	if _, err := os.Stat(piDir); err != nil {
		t.Error(".pi directory should remain since it contains user-added notes.txt")
	}
}

// --- Section 5: Additional Integration Tests ---

func TestRepeatedStartOverwritesBackup(t *testing.T) {
	// Task 5.4: Test repeated start overwrites backup correctly
	tmpDir := t.TempDir()
	storeDir := filepath.Join(tmpDir, "store")
	store := NewImageStore(storeDir)
	store.Initialize()

	workspaceDir := filepath.Join(tmpDir, "workspace")
	os.MkdirAll(workspaceDir, 0755)

	candidateDir := filepath.Join(tmpDir, "candidate")
	os.MkdirAll(candidateDir, 0755)

	metadata := NewMetadata("repeat-test", "v1")
	metadata.SetDigest("sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	metadata.SetRecipe("./Imagefile", "FROM scratch\nCOPY file.txt ./\nTAG repeat-test:v1")
	metadata.AddFile("file.txt", "workspace")
	metadata.SaveMetadataToDir(candidateDir)

	os.WriteFile(filepath.Join(candidateDir, "file.txt"), []byte("image content"), 0644)
	store.PushImage(candidateDir, false)

	// Create workspace with conflicting file
	os.WriteFile(filepath.Join(workspaceDir, "file.txt"), []byte("original content"), 0644)

	m := NewMaterializer(storeDir)

	// First start: creates backup of "original content"
	err := m.Start("repeat-test:v1", workspaceDir, "repeat-instance-1")
	if err != nil {
		t.Fatalf("first start failed: %v", err)
	}

	// Stop first instance
	err = m.Stop("repeat-instance-1")
	if err != nil {
		t.Fatalf("stop failed: %v", err)
	}

	// Modify workspace after stop
	os.WriteFile(filepath.Join(workspaceDir, "file.txt"), []byte("modified content"), 0644)

	// Second start: should backup "modified content"
	err = m.Start("repeat-test:v1", workspaceDir, "repeat-instance-2")
	if err != nil {
		t.Fatalf("second start failed: %v", err)
	}

	// Verify backup contains "modified content"
	backupDir := filepath.Join(storeDir, ".workspace", "repeat-instance-2", "backup", "data")
	backedUpFile := filepath.Join(backupDir, "file.txt")
	content, err := os.ReadFile(backedUpFile)
	if err != nil {
		t.Fatalf("failed to read backup: %v", err)
	}

	if string(content) != "modified content" {
		t.Errorf("expected backup to contain 'modified content', got '%s'", string(content))
	}

	// Stop second instance and verify restoration
	err = m.Stop("repeat-instance-2")
	if err != nil {
		t.Fatalf("second stop failed: %v", err)
	}

	// Verify "modified content" is restored
	restoredContent, err := os.ReadFile(filepath.Join(workspaceDir, "file.txt"))
	if err != nil {
		t.Fatalf("failed to read restored file: %v", err)
	}

	if string(restoredContent) != "modified content" {
		t.Errorf("expected 'modified content' after restoration, got '%s'", string(restoredContent))
	}
}
