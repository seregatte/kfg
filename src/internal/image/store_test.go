package image

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewImageStore(t *testing.T) {
	store := NewImageStore("")

	if store == nil {
		t.Fatal("expected store to be created")
	}

	// Default store directory should be set
	storeDir := store.GetStoreDir()
	if storeDir == "" {
		t.Error("expected default store directory to be set")
	}
}

func TestImageStoreInitialize(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewImageStore(tmpDir)

	err := store.Initialize()
	if err != nil {
		t.Fatalf("initialization failed: %v", err)
	}

	// Check images directory exists
	imagesDir := store.GetImagesDir()
	if _, err := os.Stat(imagesDir); err != nil {
		t.Errorf("images directory not created: %v", err)
	}
}

func TestResolveRef(t *testing.T) {
	tests := []struct {
		input        string
		expectedName string
		expectedTag  string
	}{
		{"claude-base:v2", "claude-base", "v2"},
		{"my-config:latest", "my-config", "latest"},
		{"no-tag", "no-tag", "latest"},
		{"test:1.0.0", "test", "1.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			name, tag := ResolveRef(tt.input)
			if name != tt.expectedName {
				t.Errorf("expected name '%s', got '%s'", tt.expectedName, name)
			}
			if tag != tt.expectedTag {
				t.Errorf("expected tag '%s', got '%s'", tt.expectedTag, tag)
			}
		})
	}
}

func TestPushAndLoadImage(t *testing.T) {
	// Create temp store and candidate
	tmpDir := t.TempDir()
	store := NewImageStore(tmpDir)

	// Initialize store
	if err := store.Initialize(); err != nil {
		t.Fatalf("failed to initialize store: %v", err)
	}

	// Create a test image candidate
	candidateDir := filepath.Join(tmpDir, "candidate")
	if err := os.MkdirAll(candidateDir, 0755); err != nil {
		t.Fatalf("failed to create candidate directory: %v", err)
	}

	// Create metadata
	metadata := NewMetadata("test-image", "v1")
	metadata.SetDigest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	metadata.SetRecipe("./Imagefile", "FROM scratch\nTAG test-image:v1")
	metadata.AddFile("test.txt", "workspace")

	// Save metadata to candidate
	if err := metadata.SaveMetadataToDir(candidateDir); err != nil {
		t.Fatalf("failed to save metadata: %v", err)
	}

	// Create a test file in candidate
	testFile := filepath.Join(candidateDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Push image to store
	err := store.PushImage(candidateDir, false)
	if err != nil {
		t.Fatalf("push failed: %v", err)
	}

	// Verify image directory exists in store
	imageDir := store.GetImageDir("test-image", "v1")
	if _, err := os.Stat(imageDir); err != nil {
		t.Errorf("image directory not found in store: %v", err)
	}

	// Load image from store
	loadedMetadata, loadedDir, err := store.LoadImage("test-image:v1")
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loadedMetadata.Name != metadata.Name {
		t.Errorf("expected name '%s', got '%s'", metadata.Name, loadedMetadata.Name)
	}

	if loadedDir != imageDir {
		t.Errorf("expected loadedDir '%s', got '%s'", imageDir, loadedDir)
	}
}

func TestPushExistingImage(t *testing.T) {
	// Test that pushing existing image fails (immutability)
	tmpDir := t.TempDir()
	store := NewImageStore(tmpDir)
	store.Initialize()

	// Create candidate with same name:tag twice
	candidateDir := filepath.Join(tmpDir, "candidate1")
	os.MkdirAll(candidateDir, 0755)

	metadata := NewMetadata("duplicate", "v1")
	metadata.SetDigest("sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	metadata.SetRecipe("./Imagefile", "FROM scratch\nTAG duplicate:v1")
	metadata.SaveMetadataToDir(candidateDir)

	// First push should succeed
	err := store.PushImage(candidateDir, false)
	if err != nil {
		t.Fatalf("first push failed: %v", err)
	}

	// Create second candidate with same name:tag
	candidateDir2 := filepath.Join(tmpDir, "candidate2")
	os.MkdirAll(candidateDir2, 0755)
	metadata2 := NewMetadata("duplicate", "v1")
	metadata2.SetDigest("sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc")
	metadata2.SetRecipe("./Imagefile", "FROM scratch\nTAG duplicate:v1")
	metadata2.SaveMetadataToDir(candidateDir2)

	// Second push should fail (immutability)
	err = store.PushImage(candidateDir2, false)
	if err == nil {
		t.Error("expected second push to fail (image already exists)")
	}
}

func TestListImages(t *testing.T) {
	// Test listing images
	tmpDir := t.TempDir()
	store := NewImageStore(tmpDir)
	store.Initialize()

	// Initially should have no images
	images, err := store.ListImages()
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(images) != 0 {
		t.Errorf("expected 0 images initially, got %d", len(images))
	}

	// Push an image
	candidateDir := filepath.Join(tmpDir, "candidate")
	os.MkdirAll(candidateDir, 0755)
	metadata := NewMetadata("list-test", "v1")
	metadata.SetDigest("sha256:dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd")
	metadata.SetRecipe("./Imagefile", "FROM scratch\nTAG list-test:v1")
	metadata.SaveMetadataToDir(candidateDir)

	store.PushImage(candidateDir, false)

	// List should now have 1 image
	images, err = store.ListImages()
	if err != nil {
		t.Fatalf("list after push failed: %v", err)
	}
	if len(images) != 1 {
		t.Errorf("expected 1 image, got %d", len(images))
	}

	// Verify image info
	if images[0].Name != "list-test" {
		t.Errorf("expected name 'list-test', got '%s'", images[0].Name)
	}
	if images[0].Tag != "v1" {
		t.Errorf("expected tag 'v1', got '%s'", images[0].Tag)
	}
}

func TestRemoveImage(t *testing.T) {
	// Test removing images
	tmpDir := t.TempDir()
	store := NewImageStore(tmpDir)
	store.Initialize()

	// Push an image
	candidateDir := filepath.Join(tmpDir, "candidate")
	os.MkdirAll(candidateDir, 0755)
	metadata := NewMetadata("remove-test", "v1")
	metadata.SetDigest("sha256:eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")
	metadata.SetRecipe("./Imagefile", "FROM scratch\nTAG remove-test:v1")
	metadata.SaveMetadataToDir(candidateDir)

	store.PushImage(candidateDir, false)

	// Remove the image
	err := store.RemoveImage("remove-test:v1")
	if err != nil {
		t.Fatalf("remove failed: %v", err)
	}

	// Verify image is removed
	images, err := store.ListImages()
	if err != nil {
		t.Fatalf("list after remove failed: %v", err)
	}
	if len(images) != 0 {
		t.Errorf("expected 0 images after removal, got %d", len(images))
	}
}

func TestRemoveNonExistentImage(t *testing.T) {
	// Test removing image that doesn't exist
	tmpDir := t.TempDir()
	store := NewImageStore(tmpDir)
	store.Initialize()

	err := store.RemoveImage("nonexistent:v1")
	if err == nil {
		t.Error("expected error when removing non-existent image")
	}
}

func TestFormatRecipeOnly(t *testing.T) {
	// Test 2.1: FormatRecipeOnly returns only recipe content
	metadata := NewMetadata("test-image", "v1")
	metadata.SetDigest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	metadata.SetRecipe("./Imagefile", "FROM scratch\nCOPY file.txt ./\nTAG test-image:v1")

	output := FormatRecipeOnly(metadata)
	expected := "FROM scratch\nCOPY file.txt ./\nTAG test-image:v1"
	if output != expected {
		t.Errorf("expected '%s', got '%s'", expected, output)
	}

	// Test 2.3: Output contains no metadata fields
	if strings.Contains(output, "Name:") || strings.Contains(output, "Tag:") {
		t.Error("recipe output should not contain metadata fields like 'Name:' or 'Tag:'")
	}
	if strings.Contains(output, "Digest:") || strings.Contains(output, "Created:") {
		t.Error("recipe output should not contain metadata fields like 'Digest:' or 'Created:'")
	}
}

func TestFormatRecipeOnlyEmpty(t *testing.T) {
	// Test 2.2: FormatRecipeOnly with empty recipe returns empty string
	metadata := NewMetadata("test-image", "v1")
	metadata.SetDigest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	// Recipe is empty by default in NewMetadata
	output := FormatRecipeOnly(metadata)
	if output != "" {
		t.Errorf("expected empty string for empty recipe, got '%s'", output)
	}
}

func TestFormatRecipeOnlyNoMetadata(t *testing.T) {
	// Test 2.3: Verify recipe output contains NO metadata whatsoever
	metadata := NewMetadata("test-image", "v1")
	metadata.SetDigest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	metadata.SetRecipe("./Imagefile", "# Simple Imagefile\nFROM scratch\nCOPY artifact.txt ./")

	output := FormatRecipeOnly(metadata)

	// Check that metadata fields are not present
	metadataFields := []string{
		"Name:", "Tag:", "Digest:", "Created:", "Files:",
		"Source Images:", "Recipe:", "Format:",
		"test-image", "v1", "sha256:",
	}

	for _, field := range metadataFields {
		if strings.Contains(output, field) {
			t.Errorf("recipe output should not contain metadata field '%s'", field)
		}
	}

	// Verify it's just the recipe content
	expected := "# Simple Imagefile\nFROM scratch\nCOPY artifact.txt ./"
	if output != expected {
		t.Errorf("expected '%s', got '%s'", expected, output)
	}
}

func TestFormatFilesListMultipleFiles(t *testing.T) {
	// Test 4.1: FormatFilesList with multiple files outputs one path per line
	metadata := NewMetadata("test-image", "v1")
	metadata.SetDigest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	metadata.AddFile("CLAUDE.md", "workspace")
	metadata.AddFile(".pi/config.json", "workspace")
	metadata.AddFile("README.md", "workspace")

	output := FormatFilesList(metadata)

	// Verify output is newline-separated
	lines := strings.Split(output, "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}

	// Verify each path is present
	expectedPaths := []string{".pi/config.json", "CLAUDE.md", "README.md"}
	for _, expected := range expectedPaths {
		if !strings.Contains(output, expected) {
			t.Errorf("expected path '%s' not found in output", expected)
		}
	}
}

func TestFormatFilesListEmptyFiles(t *testing.T) {
	// Test 4.2: FormatFilesList with empty files outputs "No files"
	metadata := NewMetadata("test-image", "v1")
	metadata.SetDigest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	output := FormatFilesList(metadata)

	if output != "No files" {
		t.Errorf("expected 'No files', got '%s'", output)
	}
}

func TestFormatFilesListSortedAlphabetically(t *testing.T) {
	// Test 4.3: FormatFilesList outputs paths sorted alphabetically
	metadata := NewMetadata("test-image", "v1")
	metadata.SetDigest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	// Add files in non-alphabetical order
	metadata.AddFile("z-file.txt", "workspace")
	metadata.AddFile("a-file.txt", "workspace")
	metadata.AddFile("m-file.txt", "workspace")
	metadata.AddFile("b-file.txt", "workspace")

	output := FormatFilesList(metadata)

	// Verify alphabetical order
	lines := strings.Split(output, "\n")
	expected := []string{"a-file.txt", "b-file.txt", "m-file.txt", "z-file.txt"}
	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("expected line %d to be '%s', got '%s'", i, expected[i], line)
		}
	}
}

func TestFormatFilesListNoMetadata(t *testing.T) {
	// Test 4.4: FormatFilesList output contains no metadata fields
	metadata := NewMetadata("test-image", "v1")
	metadata.SetDigest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	metadata.AddFile("CLAUDE.md", "workspace")
	metadata.AddFile("config.json", "workspace")

	output := FormatFilesList(metadata)

	// Check that metadata fields are not present
	metadataFields := []string{
		"Name:", "Tag:", "Digest:", "Created:", "Files:",
		"Source Images:", "Recipe:",
		"test-image", "v1", "sha256:",
	}

	for _, field := range metadataFields {
		if strings.Contains(output, field) {
			t.Errorf("files output should not contain metadata field '%s', got '%s'", field, output)
		}
	}

	// Verify output is just paths
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line != "CLAUDE.md" && line != "config.json" {
			t.Errorf("unexpected line '%s' in files output", line)
		}
	}
}

func TestFormatFilesListJSON_Empty(t *testing.T) {
	metadata := NewMetadata("test-image", "v1")
	metadata.SetDigest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	output, err := FormatFilesListJSON(metadata)
	if err != nil {
		t.Fatalf("FormatFilesListJSON failed: %v", err)
	}

	if output != "[]" {
		t.Errorf("expected '[]', got '%s'", output)
	}
}

func TestFormatFilesListJSON_Multiple(t *testing.T) {
	metadata := NewMetadata("test-image", "v1")
	metadata.SetDigest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	metadata.AddFile("CLAUDE.md", "workspace")
	metadata.AddFile(".pi/config.json", "workspace")
	metadata.AddFile("README.md", "workspace")

	output, err := FormatFilesListJSON(metadata)
	if err != nil {
		t.Fatalf("FormatFilesListJSON failed: %v", err)
	}

	if !strings.Contains(output, "CLAUDE.md") {
		t.Error("expected output to contain 'CLAUDE.md'")
	}
	if !strings.Contains(output, ".pi/config.json") {
		t.Error("expected output to contain '.pi/config.json'")
	}
	if !strings.Contains(output, "README.md") {
		t.Error("expected output to contain 'README.md'")
	}
	if !strings.HasPrefix(output, "[") || !strings.HasSuffix(output, "]") {
		t.Error("expected JSON array format starting with '[' and ending with ']'")
	}
}

func TestFormatFilesListJSON_Sorted(t *testing.T) {
	metadata := NewMetadata("test-image", "v1")
	metadata.SetDigest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	metadata.AddFile("z-file.txt", "workspace")
	metadata.AddFile("a-file.txt", "workspace")
	metadata.AddFile("m-file.txt", "workspace")
	metadata.AddFile("b-file.txt", "workspace")

	output, err := FormatFilesListJSON(metadata)
	if err != nil {
		t.Fatalf("FormatFilesListJSON failed: %v", err)
	}

	var paths []string
	if err := json.Unmarshal([]byte(output), &paths); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	expectedOrder := []string{"a-file.txt", "b-file.txt", "m-file.txt", "z-file.txt"}
	if len(paths) != len(expectedOrder) {
		t.Errorf("expected %d paths, got %d", len(expectedOrder), len(paths))
	}

	for i, actual := range paths {
		if actual != expectedOrder[i] {
			t.Errorf("expected sorted order at position %d to be '%s', got '%s'", i, expectedOrder[i], actual)
		}
	}
}

func TestFormatListJSON_Empty(t *testing.T) {
	images := []ImageInfo{}

	output, err := FormatListJSON(images)
	if err != nil {
		t.Fatalf("FormatListJSON failed: %v", err)
	}

	if output != "[]" {
		t.Errorf("expected '[]' for empty images, got '%s'", output)
	}
}
