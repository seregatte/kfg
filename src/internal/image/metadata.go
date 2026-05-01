// Package image provides metadata structures and operations for NixAI images.
// This file implements the metadata handling for image persistence and inspection.
package image

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ImageMetadata represents the complete metadata for a stored image.
type ImageMetadata struct {
	// Name is the image name (e.g., "claude-base")
	Name string `json:"name"`

	// Tag is the image tag (e.g., "v2")
	Tag string `json:"tag"`

	// ImageDigest is the SHA256 digest of the image content
	ImageDigest string `json:"image_digest"`

	// CreatedAt is the ISO 8601 timestamp when the image was created
	CreatedAt string `json:"created_at"`

	// SourceImages tracks all parent images referenced in FROM instructions
	SourceImages []SourceImage `json:"source_images"`

	// Recipe stores the original Imagefile content
	Recipe Recipe `json:"recipe"`

	// Files is the manifest of all files in the image (path -> source)
	Files map[string]string `json:"files"`

	// FormatVersion identifies the metadata schema version
	FormatVersion string `json:"format_version"`
}

// SourceImage represents a parent image reference with resolved digest.
type SourceImage struct {
	// Ref is the original reference (e.g., "claude-base:v2")
	Ref string `json:"ref"`

	// ResolvedDigest is the actual digest at build time
	ResolvedDigest string `json:"resolved_digest"`
}

// Recipe stores the original Imagefile content for reproducibility.
type Recipe struct {
	// SourcePath is the relative path to the Imagefile (e.g., "./Imagefile")
	SourcePath string `json:"source_path"`

	// Content is the full Imagefile text
	Content string `json:"content"`

	// Format identifies the Imagefile syntax version (e.g., "imagefile.v1")
	Format string `json:"format"`
}

// MetadataValidation represents validation errors for metadata.
type MetadataValidation struct {
	Errors []string
}

// NewMetadata creates a new ImageMetadata instance with defaults.
func NewMetadata(name, tag string) *ImageMetadata {
	return &ImageMetadata{
		Name:          name,
		Tag:           tag,
		ImageDigest:   "",
		CreatedAt:     time.Now().UTC().Format(time.RFC3339),
		SourceImages:  []SourceImage{},
		Recipe:        Recipe{},
		Files:         make(map[string]string),
		FormatVersion: "metadata.v1",
	}
}

// SetDigest sets the image digest.
func (m *ImageMetadata) SetDigest(digest string) {
	m.ImageDigest = digest
}

// AddSourceImage adds a source image reference.
func (m *ImageMetadata) AddSourceImage(ref, digest string) {
	m.SourceImages = append(m.SourceImages, SourceImage{
		Ref:            ref,
		ResolvedDigest: digest,
	})
}

// SetRecipe sets the recipe (Imagefile) information.
func (m *ImageMetadata) SetRecipe(sourcePath, content string) {
	m.Recipe = Recipe{
		SourcePath: sourcePath,
		Content:    content,
		Format:     "imagefile.v1",
	}
}

// AddFile adds a file to the manifest.
func (m *ImageMetadata) AddFile(path, source string) {
	m.Files[path] = source
}

// Validate checks that metadata has all required fields.
func (m *ImageMetadata) Validate() *MetadataValidation {
	validation := &MetadataValidation{Errors: []string{}}

	// Check required fields
	if m.Name == "" {
		validation.Errors = append(validation.Errors, "name is required")
	}
	if m.Tag == "" {
		validation.Errors = append(validation.Errors, "tag is required")
	}
	if m.ImageDigest == "" {
		validation.Errors = append(validation.Errors, "image_digest is required")
	}
	if m.CreatedAt == "" {
		validation.Errors = append(validation.Errors, "created_at is required")
	}
	if m.FormatVersion == "" {
		validation.Errors = append(validation.Errors, "format_version is required")
	}

	// Validate digest format
	if m.ImageDigest != "" && !isValidDigest(m.ImageDigest) {
		validation.Errors = append(validation.Errors, "image_digest must be in format 'sha256:<hex>'")
	}

	// Validate recipe
	if m.Recipe.Content == "" {
		validation.Errors = append(validation.Errors, "recipe.content is required")
	}
	if m.Recipe.Format == "" {
		validation.Errors = append(validation.Errors, "recipe.format is required")
	}

	// Validate source images
	for i, source := range m.SourceImages {
		if source.Ref == "" {
			validation.Errors = append(validation.Errors, fmt.Sprintf("source_images[%d].ref is required", i))
		}
		if source.ResolvedDigest == "" {
			validation.Errors = append(validation.Errors, fmt.Sprintf("source_images[%d].resolved_digest is required", i))
		}
		if source.ResolvedDigest != "" && !isValidDigest(source.ResolvedDigest) {
			validation.Errors = append(validation.Errors, fmt.Sprintf("source_images[%d].resolved_digest must be in format 'sha256:<hex>'", i))
		}
	}

	return validation
}

// IsValid returns true if validation has no errors.
func (v *MetadataValidation) IsValid() bool {
	return len(v.Errors) == 0
}

// Error returns a formatted error message.
func (v *MetadataValidation) Error() string {
	if len(v.Errors) == 0 {
		return ""
	}
	return fmt.Sprintf("metadata validation failed: %s", strings.Join(v.Errors, "; "))
}

// isValidDigest checks if a digest string is in the correct format.
func isValidDigest(digest string) bool {
	// Must start with "sha256:" and have 64 hex characters
	if !strings.HasPrefix(digest, "sha256:") {
		return false
	}
	hexPart := strings.TrimPrefix(digest, "sha256:")
	if len(hexPart) != 64 {
		return false
	}
	// Check that all characters are hex
	for _, c := range hexPart {
		if !isHexChar(c) {
			return false
		}
	}
	return true
}

// isHexChar checks if a character is a valid hex digit.
func isHexChar(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

// Serialize writes metadata to JSON file.
func (m *ImageMetadata) Serialize(path string) error {
	// Validate before serialization
	validation := m.Validate()
	if !validation.IsValid() {
		return fmt.Errorf("cannot serialize invalid metadata: %s", validation.Error())
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}

// Deserialize reads metadata from JSON file.
func DeserializeMetadata(path string) (*ImageMetadata, error) {
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	// Unmarshal JSON
	var metadata ImageMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	// Validate after deserialization
	validation := metadata.Validate()
	if !validation.IsValid() {
		return nil, fmt.Errorf("loaded metadata is invalid: %s", validation.Error())
	}

	return &metadata, nil
}

// LoadMetadataFromDir loads metadata from an image directory.
func LoadMetadataFromDir(imageDir string) (*ImageMetadata, error) {
	metadataPath := filepath.Join(imageDir, "metadata.json")
	return DeserializeMetadata(metadataPath)
}

// SaveMetadataToDir saves metadata to an image directory.
func (m *ImageMetadata) SaveMetadataToDir(imageDir string) error {
	// Ensure directory exists
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		return fmt.Errorf("failed to create image directory: %w", err)
	}

	metadataPath := filepath.Join(imageDir, "metadata.json")
	return m.Serialize(metadataPath)
}

// ToJSON returns metadata as JSON string.
func (m *ImageMetadata) ToJSON() (string, error) {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal metadata: %w", err)
	}
	return string(data), nil
}

// FromJSON parses metadata from JSON string.
func FromJSON(jsonStr string) (*ImageMetadata, error) {
	var metadata ImageMetadata
	if err := json.Unmarshal([]byte(jsonStr), &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	// Validate
	validation := metadata.Validate()
	if !validation.IsValid() {
		return nil, fmt.Errorf("parsed metadata is invalid: %s", validation.Error())
	}

	return &metadata, nil
}

// GetFullName returns the full image reference (name:tag).
func (m *ImageMetadata) GetFullName() string {
	return fmt.Sprintf("%s:%s", m.Name, m.Tag)
}

// GetShortDigest returns the first 12 characters of the digest hex.
func (m *ImageMetadata) GetShortDigest() string {
	// Remove sha256: prefix and take first 12 chars
	shortDigest := shortenDigestHelper(m.ImageDigest)
	return shortDigest
}

// shortenDigestHelper removes the sha256: prefix and returns first 12 chars.
func shortenDigestHelper(digest string) string {
	// Remove sha256: prefix if present
	hexPart := strings.TrimPrefix(digest, "sha256:")
	if len(hexPart) >= 12 {
		return hexPart[:12]
	}
	return hexPart
}

// GetFileCount returns the number of files in the manifest.
func (m *ImageMetadata) GetFileCount() int {
	return len(m.Files)
}

// HasSourceImage checks if a specific source image is tracked.
func (m *ImageMetadata) HasSourceImage(ref string) bool {
	for _, source := range m.SourceImages {
		if source.Ref == ref {
			return true
		}
	}
	return false
}

// GetSourceDigest returns the resolved digest for a source reference.
func (m *ImageMetadata) GetSourceDigest(ref string) string {
	for _, source := range m.SourceImages {
		if source.Ref == ref {
			return source.ResolvedDigest
		}
	}
	return ""
}

// RecipeDisplay returns the recipe content for display.
func (m *ImageMetadata) RecipeDisplay() string {
	return fmt.Sprintf("# Recipe: %s\n# Format: %s\n\n%s",
		m.Recipe.SourcePath, m.Recipe.Format, m.Recipe.Content)
}