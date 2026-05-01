package image

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewMetadata(t *testing.T) {
	metadata := NewMetadata("claude-base", "v2")

	if metadata.Name != "claude-base" {
		t.Errorf("expected name 'claude-base', got '%s'", metadata.Name)
	}
	if metadata.Tag != "v2" {
		t.Errorf("expected tag 'v2', got '%s'", metadata.Tag)
	}
	if metadata.ImageDigest != "" {
		t.Error("expected empty digest for new metadata")
	}
	if metadata.CreatedAt == "" {
		t.Error("expected created_at to be set")
	}
	if len(metadata.SourceImages) != 0 {
		t.Errorf("expected empty source_images, got %d", len(metadata.SourceImages))
	}
	if metadata.Files == nil {
		t.Error("expected files map to be initialized")
	}
	if len(metadata.Files) != 0 {
		t.Errorf("expected empty files map, got %d files", len(metadata.Files))
	}
	if metadata.FormatVersion != "metadata.v1" {
		t.Errorf("expected format_version 'metadata.v1', got '%s'", metadata.FormatVersion)
	}
}

func TestSetDigest(t *testing.T) {
	metadata := NewMetadata("test", "latest")
	digest := "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	metadata.SetDigest(digest)

	if metadata.ImageDigest != digest {
		t.Errorf("expected digest '%s', got '%s'", digest, metadata.ImageDigest)
	}
}

func TestAddSourceImage(t *testing.T) {
	metadata := NewMetadata("test", "latest")

	metadata.AddSourceImage("claude-base:v2", "sha256:eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")

	if len(metadata.SourceImages) != 1 {
		t.Errorf("expected 1 source image, got %d", len(metadata.SourceImages))
	}

	source := metadata.SourceImages[0]
	if source.Ref != "claude-base:v2" {
		t.Errorf("expected ref 'claude-base:v2', got '%s'", source.Ref)
	}
	if source.ResolvedDigest != "sha256:eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
		t.Errorf("expected digest 'sha256:eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee', got '%s'", source.ResolvedDigest)
	}
}

func TestSetRecipe(t *testing.T) {
	metadata := NewMetadata("test", "latest")
	content := "FROM scratch\nCOPY file.txt ./"

	metadata.SetRecipe("./Imagefile", content)

	if metadata.Recipe.SourcePath != "./Imagefile" {
		t.Errorf("expected source_path './Imagefile', got '%s'", metadata.Recipe.SourcePath)
	}
	if metadata.Recipe.Content != content {
		t.Errorf("expected content '%s', got '%s'", content, metadata.Recipe.Content)
	}
	if metadata.Recipe.Format != "imagefile.v1" {
		t.Errorf("expected format 'imagefile.v1', got '%s'", metadata.Recipe.Format)
	}
}

func TestAddFile(t *testing.T) {
	metadata := NewMetadata("test", "latest")

	metadata.AddFile("AGENTS.md", "stage:base")
	metadata.AddFile(".claude/", "stage:claude")

	if len(metadata.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(metadata.Files))
	}
	if metadata.Files["AGENTS.md"] != "stage:base" {
		t.Errorf("expected 'stage:base', got '%s'", metadata.Files["AGENTS.md"])
	}
	if metadata.Files[".claude/"] != "stage:claude" {
		t.Errorf("expected 'stage:claude', got '%s'", metadata.Files[".claude/"])
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		metadata    *ImageMetadata
		expectValid bool
		errorCount  int
	}{
		{
			name:        "empty metadata",
			metadata:    &ImageMetadata{},
			expectValid: false,
			errorCount:  7, // name, tag, digest, created_at, format_version, recipe.content, recipe.format
		},
		{
			name: "missing name",
			metadata: &ImageMetadata{
				Tag:           "v1",
				ImageDigest:   "sha256:1111111111111111111111111111111111111111111111111111111111111111",
				CreatedAt:     "2024-01-01T00:00:00Z",
				FormatVersion: "metadata.v1",
				Recipe:        Recipe{Content: "FROM scratch", Format: "imagefile.v1"},
			},
			expectValid: false,
			errorCount:  1,
		},
		{
			name: "missing tag",
			metadata: &ImageMetadata{
				Name:          "test",
				ImageDigest:   "sha256:1111111111111111111111111111111111111111111111111111111111111111",
				CreatedAt:     "2024-01-01T00:00:00Z",
				FormatVersion: "metadata.v1",
				Recipe:        Recipe{Content: "FROM scratch", Format: "imagefile.v1"},
			},
			expectValid: false,
			errorCount:  1,
		},
		{
			name: "invalid digest format",
			metadata: &ImageMetadata{
				Name:          "test",
				Tag:           "v1",
				ImageDigest:   "invalid-digest",
				CreatedAt:     "2024-01-01T00:00:00Z",
				FormatVersion: "metadata.v1",
				Recipe:        Recipe{Content: "FROM scratch", Format: "imagefile.v1"},
			},
			expectValid: false,
			errorCount:  1,
		},
		{
			name: "valid metadata",
			metadata: &ImageMetadata{
				Name:          "test",
				Tag:           "v1",
				ImageDigest:   "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				CreatedAt:     "2024-01-01T00:00:00Z",
				FormatVersion: "metadata.v1",
				Recipe:        Recipe{Content: "FROM scratch", Format: "imagefile.v1"},
				Files:         make(map[string]string),
			},
			expectValid: true,
			errorCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validation := tt.metadata.Validate()

			if tt.expectValid && !validation.IsValid() {
				t.Errorf("expected valid, but got errors: %s", validation.Error())
			}

			if !tt.expectValid && validation.IsValid() {
				t.Error("expected invalid, but validation passed")
			}

			if !tt.expectValid && len(validation.Errors) != tt.errorCount {
				t.Errorf("expected %d errors, got %d: %v", tt.errorCount, len(validation.Errors), validation.Errors)
			}
		})
	}
}

func TestIsValidDigest(t *testing.T) {
	tests := []struct {
		digest   string
		expected bool
	}{
		{"sha256:1111111111111111111111111111111111111111111111111111111111111111", true},                                                        // valid 64 chars
		{"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa789012345678901234567890123456789012345678901234567890", false}, // too long (>64 chars)
		{"sha256:dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd90", false},                                                     // too long (66 chars)
		{"sha256:", false},       // empty hex part
		{"sha512:abc123", false}, // wrong algorithm
		{"abc123", false},        // no algorithm prefix
		{"sha256:1111111111111111111111111111111111111111111111111111111111111111xyz", false}, // non-hex characters
	}

	for _, tt := range tests {
		t.Run(tt.digest, func(t *testing.T) {
			result := isValidDigest(tt.digest)
			if result != tt.expected {
				t.Errorf("isValidDigest('%s') = %v, expected %v", tt.digest, result, tt.expected)
			}
		})
	}
}

func TestSerializeDeserialize(t *testing.T) {
	// Create test metadata
	metadata := NewMetadata("test-image", "v1.0")
	metadata.SetDigest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	metadata.SetRecipe("./Imagefile", "FROM scratch\nCOPY file.txt ./")
	metadata.AddFile("file.txt", "workspace")
	metadata.AddSourceImage("base:v1", "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	// Create temp file
	tmpDir := t.TempDir()
	metadataPath := filepath.Join(tmpDir, "metadata.json")

	// Serialize
	err := metadata.Serialize(metadataPath)
	if err != nil {
		t.Fatalf("serialization failed: %v", err)
	}

	// Check file exists
	if _, err := os.Stat(metadataPath); err != nil {
		t.Fatalf("metadata file not created: %v", err)
	}

	// Deserialize
	loaded, err := DeserializeMetadata(metadataPath)
	if err != nil {
		t.Fatalf("deserialization failed: %v", err)
	}

	// Verify content
	if loaded.Name != metadata.Name {
		t.Errorf("expected name '%s', got '%s'", metadata.Name, loaded.Name)
	}
	if loaded.Tag != metadata.Tag {
		t.Errorf("expected tag '%s', got '%s'", metadata.Tag, loaded.Tag)
	}
	if loaded.ImageDigest != metadata.ImageDigest {
		t.Errorf("expected digest '%s', got '%s'", metadata.ImageDigest, loaded.ImageDigest)
	}
	if loaded.Recipe.Content != metadata.Recipe.Content {
		t.Errorf("expected recipe content '%s', got '%s'", metadata.Recipe.Content, loaded.Recipe.Content)
	}
	if len(loaded.Files) != len(metadata.Files) {
		t.Errorf("expected %d files, got %d", len(metadata.Files), len(loaded.Files))
	}
	if len(loaded.SourceImages) != len(metadata.SourceImages) {
		t.Errorf("expected %d source images, got %d", len(metadata.SourceImages), len(loaded.SourceImages))
	}
}

func TestSerializeInvalidMetadata(t *testing.T) {
	// Create invalid metadata (missing required fields)
	metadata := &ImageMetadata{
		Name: "test",
		// Missing tag, digest, etc.
	}

	tmpDir := t.TempDir()
	metadataPath := filepath.Join(tmpDir, "metadata.json")

	// Attempt to serialize - should fail validation
	err := metadata.Serialize(metadataPath)
	if err == nil {
		t.Error("expected serialization to fail for invalid metadata")
	}

	// File should not be created
	if _, err := os.Stat(metadataPath); err == nil {
		t.Error("metadata file should not be created for invalid metadata")
	}
}

func TestDeserializeInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	metadataPath := filepath.Join(tmpDir, "metadata.json")

	// Write invalid JSON
	os.WriteFile(metadataPath, []byte("invalid json"), 0644)

	// Attempt to deserialize - should fail
	_, err := DeserializeMetadata(metadataPath)
	if err == nil {
		t.Error("expected deserialization to fail for invalid JSON")
	}
}

func TestLoadMetadataFromDir(t *testing.T) {
	// Create test metadata
	metadata := NewMetadata("test", "v1")
	metadata.SetDigest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	metadata.SetRecipe("./Imagefile", "FROM scratch")

	// Create temp directory with metadata
	tmpDir := t.TempDir()
	imageDir := filepath.Join(tmpDir, "test:v1")

	err := metadata.SaveMetadataToDir(imageDir)
	if err != nil {
		t.Fatalf("failed to save metadata: %v", err)
	}

	// Load from directory
	loaded, err := LoadMetadataFromDir(imageDir)
	if err != nil {
		t.Fatalf("failed to load metadata: %v", err)
	}

	if loaded.Name != metadata.Name {
		t.Errorf("expected name '%s', got '%s'", metadata.Name, loaded.Name)
	}
}

func TestToJSONFromJSON(t *testing.T) {
	metadata := NewMetadata("test", "v1")
	metadata.SetDigest("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	metadata.SetRecipe("./Imagefile", "FROM scratch")

	// Convert to JSON
	jsonStr, err := metadata.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	// Parse from JSON
	loaded, err := FromJSON(jsonStr)
	if err != nil {
		t.Fatalf("FromJSON failed: %v", err)
	}

	if loaded.Name != metadata.Name {
		t.Errorf("expected name '%s', got '%s'", metadata.Name, loaded.Name)
	}
}

func TestGetFullName(t *testing.T) {
	metadata := NewMetadata("claude-base", "v2")

	fullName := metadata.GetFullName()
	expected := "claude-base:v2"

	if fullName != expected {
		t.Errorf("expected '%s', got '%s'", expected, fullName)
	}
}

func TestGetFileCount(t *testing.T) {
	metadata := NewMetadata("test", "v1")

	if metadata.GetFileCount() != 0 {
		t.Errorf("expected 0 files, got %d", metadata.GetFileCount())
	}

	metadata.AddFile("file1.txt", "workspace")
	metadata.AddFile("file2.txt", "stage:base")

	if metadata.GetFileCount() != 2 {
		t.Errorf("expected 2 files, got %d", metadata.GetFileCount())
	}
}

func TestHasSourceImage(t *testing.T) {
	metadata := NewMetadata("test", "v1")
	metadata.AddSourceImage("base:v1", "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	if !metadata.HasSourceImage("base:v1") {
		t.Error("expected to find source image 'base:v1'")
	}

	if metadata.HasSourceImage("other:v2") {
		t.Error("expected not to find source image 'other:v2'")
	}
}

func TestGetSourceDigest(t *testing.T) {
	metadata := NewMetadata("test", "v1")
	metadata.AddSourceImage("base:v1", "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	digest := metadata.GetSourceDigest("base:v1")
	if digest != "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" {
		t.Errorf("expected digest 'sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb', got '%s'", digest)
	}

	digest = metadata.GetSourceDigest("other:v2")
	if digest != "" {
		t.Errorf("expected empty digest for unknown source, got '%s'", digest)
	}
}

func TestRecipeDisplay(t *testing.T) {
	metadata := NewMetadata("test", "v1")
	metadata.SetRecipe("./Imagefile", "FROM scratch\nCOPY file.txt ./")

	display := metadata.RecipeDisplay()
	expected := "# Recipe: ./Imagefile\n# Format: imagefile.v1\n\nFROM scratch\nCOPY file.txt ./"

	if display != expected {
		t.Errorf("expected recipe display:\n%s\n\ngot:\n%s", expected, display)
	}
}

func TestMetadataValidationError(t *testing.T) {
	validation := &MetadataValidation{Errors: []string{"error1", "error2"}}

	if validation.IsValid() {
		t.Error("expected validation to be invalid")
	}

	errorMsg := validation.Error()
	if errorMsg == "" {
		t.Error("expected non-empty error message")
	}

	if !contains(errorMsg, "error1") || !contains(errorMsg, "error2") {
		t.Errorf("expected error message to contain both errors, got: %s", errorMsg)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
