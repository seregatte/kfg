// Package image provides image storage and retrieval operations for KFG.
// This file implements the store persistence layer for managing image artifacts.
package image

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/seregatte/kfg/src/internal/logger"
)

const (
	// DefaultStoreDir is the default location for the image store
	DefaultStoreDir = "~/.config/kfg/store"

	// ImagesSubdir is the subdirectory for stored images
	ImagesSubdir = "images"

	// MetadataFile is the filename for image metadata
	MetadataFile = "metadata.json"
)

// StoreConfig holds configuration for the store.
type StoreConfig struct {
	// StoreDir is the root directory for the store
	StoreDir string
}

// ImageStore manages stored images in the NixAI store.
type ImageStore struct {
	config StoreConfig
}

// ImageInfo represents a summarized image entry for listing.
type ImageInfo struct {
	Name        string `json:"name"`
	Tag         string `json:"tag"`
	Digest      string `json:"digest"`       // Full digest
	ShortDigest string `json:"short_digest"` // First 12 chars
	CreatedAt   string `json:"created_at"`
	FileCount   int    `json:"file_count"`
}

// NewImageStore creates a new ImageStore instance.
func NewImageStore(storeDir string) *ImageStore {
	if storeDir == "" {
		storeDir = DefaultStoreDir
	}

	// Expand ~ to home directory
	if strings.HasPrefix(storeDir, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// Fall back to current directory if home dir unavailable
			storeDir = ".nixai/store"
		} else {
			storeDir = filepath.Join(homeDir, strings.TrimPrefix(storeDir, "~"))
		}
	}

	return &ImageStore{
		config: StoreConfig{StoreDir: storeDir},
	}
}

// GetStoreDir returns the store directory path.
func (s *ImageStore) GetStoreDir() string {
	return s.config.StoreDir
}

// GetImagesDir returns the images directory path.
func (s *ImageStore) GetImagesDir() string {
	return filepath.Join(s.config.StoreDir, ImagesSubdir)
}

// Initialize creates the store directory structure if it doesn't exist.
func (s *ImageStore) Initialize() error {
	imagesDir := s.GetImagesDir()
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		return fmt.Errorf("failed to create images directory: %w", err)
	}
	return nil
}

// ResolveRef resolves an image reference to name and tag.
// If tag is omitted, defaults to "latest".
func ResolveRef(ref string) (name, tag string) {
	// Check if reference contains a tag
	if strings.Contains(ref, ":") {
		// Split on last colon (digests have multiple colons)
		parts := strings.SplitN(ref, ":", 2)
		if len(parts) == 2 {
			// Check if this looks like a digest (sha256:...)
			if parts[1] == "sha256" || strings.HasPrefix(parts[1], "sha256:") {
				// This is a digest reference, not a tag
				// For now, treat as name with no tag (defaults to :latest)
				return ref, "latest"
			}
			return parts[0], parts[1]
		}
	}
	// No tag, default to latest
	return ref, "latest"
}

// GetImageDir returns the directory path for a specific image.
func (s *ImageStore) GetImageDir(name, tag string) string {
	return filepath.Join(s.GetImagesDir(), name, tag)
}

// PushImage persists a built image candidate to the store.
// The candidate directory should contain the image files and metadata.
func (s *ImageStore) PushImage(candidateDir string, keepBuild bool) error {
	// Ensure store is initialized
	if err := s.Initialize(); err != nil {
		return err
	}

	// Load metadata from candidate
	metadata, err := LoadMetadataFromDir(candidateDir)
	if err != nil {
		return fmt.Errorf("failed to load candidate metadata: %w", err)
	}

	// Validate metadata
	validation := metadata.Validate()
	if !validation.IsValid() {
		return fmt.Errorf("candidate metadata is invalid: %s", validation.Error())
	}

	// Check if image already exists (immutable)
	imageDir := s.GetImageDir(metadata.Name, metadata.Tag)
	if _, err := os.Stat(imageDir); err == nil {
		return fmt.Errorf("image %s:%s already exists in store (images are immutable)", metadata.Name, metadata.Tag)
	}

	// Create image directory
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		return fmt.Errorf("failed to create image directory: %w", err)
	}

	// Copy all files from candidate to store
	err = copyDirectory(candidateDir, imageDir)
	if err != nil {
		// Clean up on failure
		os.RemoveAll(imageDir)
		return fmt.Errorf("failed to copy image files: %w", err)
	}

	// Clean up candidate unless keepBuild is true
	if !keepBuild {
		if err := os.RemoveAll(candidateDir); err != nil {
			// Log warning but don't fail
			logger.Warn("store:push", fmt.Sprintf("Failed to clean up build directory: %v", err))
		}
	}

	return nil
}

// LoadImage loads an image from the store for use in FROM instructions.
func (s *ImageStore) LoadImage(ref string) (*ImageMetadata, string, error) {
	// Resolve reference
	name, tag := ResolveRef(ref)

	// Get image directory
	imageDir := s.GetImageDir(name, tag)

	// Check if image exists
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		return nil, "", fmt.Errorf("image %s:%s not found in store", name, tag)
	}

	// Load metadata
	metadata, err := LoadMetadataFromDir(imageDir)
	if err != nil {
		return nil, "", fmt.Errorf("failed to load image metadata: %w", err)
	}

	return metadata, imageDir, nil
}

// ListImages returns all images in the store.
func (s *ImageStore) ListImages() ([]ImageInfo, error) {
	// Ensure store directory exists
	imagesDir := s.GetImagesDir()
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		return []ImageInfo{}, nil // Empty store
	}

	// Scan for images
	var images []ImageInfo

	// Walk through name directories
	nameDirs, err := os.ReadDir(imagesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read images directory: %w", err)
	}

	for _, nameEntry := range nameDirs {
		if !nameEntry.IsDir() {
			continue
		}

		name := nameEntry.Name()
		namePath := filepath.Join(imagesDir, name)

		// Walk through tag directories
		tagDirs, err := os.ReadDir(namePath)
		if err != nil {
			continue // Skip invalid directories
		}

		for _, tagEntry := range tagDirs {
			if !tagEntry.IsDir() {
				continue
			}

			tag := tagEntry.Name()
			imageDir := filepath.Join(namePath, tag)

			// Load metadata
			metadata, err := LoadMetadataFromDir(imageDir)
			if err != nil {
				// Skip images with invalid metadata
				continue
			}

			// Create ImageInfo
			info := ImageInfo{
				Name:        metadata.Name,
				Tag:         metadata.Tag,
				Digest:      metadata.ImageDigest,
				ShortDigest: shortenDigest(metadata.ImageDigest),
				CreatedAt:   metadata.CreatedAt,
				FileCount:   metadata.GetFileCount(),
			}

			images = append(images, info)
		}
	}

	// Sort by name, then by tag
	sort.Slice(images, func(i, j int) bool {
		if images[i].Name != images[j].Name {
			return images[i].Name < images[j].Name
		}
		return images[i].Tag < images[j].Tag
	})

	return images, nil
}

// InspectImage returns detailed metadata for an image.
func (s *ImageStore) InspectImage(ref string) (*ImageMetadata, error) {
	// Resolve reference
	name, tag := ResolveRef(ref)

	// Get image directory
	imageDir := s.GetImageDir(name, tag)

	// Check if image exists
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		// Try to provide helpful suggestions
		suggestions := s.findSimilarImages(name)
		if len(suggestions) > 0 {
			return nil, fmt.Errorf("image %s:%s not found. Similar images: %s", name, tag, strings.Join(suggestions, ", "))
		}
		return nil, fmt.Errorf("image %s:%s not found in store", name, tag)
	}

	// Load metadata
	metadata, err := LoadMetadataFromDir(imageDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load image metadata: %w", err)
	}

	return metadata, nil
}

// RemoveImage deletes an image from the store.
func (s *ImageStore) RemoveImage(ref string) error {
	// Resolve reference
	name, tag := ResolveRef(ref)

	// Get image directory
	imageDir := s.GetImageDir(name, tag)

	// Check if image exists
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		return fmt.Errorf("image %s:%s not found in store", name, tag)
	}

	// Remove image directory
	if err := os.RemoveAll(imageDir); err != nil {
		return fmt.Errorf("failed to remove image: %w", err)
	}

	// Clean up empty parent directory if no tags remain
	nameDir := filepath.Join(s.GetImagesDir(), name)
	remainingTags, err := os.ReadDir(nameDir)
	if err == nil && len(remainingTags) == 0 {
		os.Remove(nameDir)
	}

	return nil
}

// findSimilarImages returns images with similar names for helpful error messages.
func (s *ImageStore) findSimilarImages(name string) []string {
	imagesDir := s.GetImagesDir()
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		return []string{}
	}

	var suggestions []string

	// Look for names that contain the search name
	nameDirs, err := os.ReadDir(imagesDir)
	if err != nil {
		return []string{}
	}

	for _, nameEntry := range nameDirs {
		if !nameEntry.IsDir() {
			continue
		}

		foundName := nameEntry.Name()
		if strings.Contains(foundName, name) || strings.Contains(name, foundName) {
			// Get all tags for this name
			namePath := filepath.Join(imagesDir, foundName)
			tagDirs, err := os.ReadDir(namePath)
			if err != nil {
				continue
			}

			for _, tagEntry := range tagDirs {
				if tagEntry.IsDir() {
					suggestions = append(suggestions, fmt.Sprintf("%s:%s", foundName, tagEntry.Name()))
				}
			}
		}
	}

	sort.Strings(suggestions)
	return suggestions
}

// shortenDigest returns the first 12 characters of a digest.
func shortenDigest(digest string) string {
	// Remove sha256: prefix if present
	hexPart := strings.TrimPrefix(digest, "sha256:")
	if len(hexPart) >= 12 {
		return hexPart[:12]
	}
	return hexPart
}

// copyDirectory copies all files from source to destination.
func copyDirectory(src, dst string) error {
	// Get file info
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create destination directory
	if err := os.MkdirAll(dst, info.Mode()); err != nil {
		return err
	}

	// Read source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy directory
			if err := copyDirectory(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file from source to destination.
func copyFile(src, dst string) error {
	// Read source file
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Get file info for permissions
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Write to destination
	if err := os.WriteFile(dst, data, info.Mode()); err != nil {
		return err
	}

	return nil
}

// FormatListTable formats images as a human-readable table.
func FormatListTable(images []ImageInfo) string {
	if len(images) == 0 {
		return "No images found"
	}

	// Header
	output := "NAME\tTAG\tDIGEST\tCREATED\tFILES\n"

	// Rows
	for _, img := range images {
		created := formatTime(img.CreatedAt)
		output += fmt.Sprintf("%s\t%s\t%s\t%s\t%d\n",
			img.Name, img.Tag, img.ShortDigest, created, img.FileCount)
	}

	return output
}

// FormatListJSON formats images as JSON array.
func FormatListJSON(images []ImageInfo) (string, error) {
	data, err := json.MarshalIndent(images, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format JSON: %w", err)
	}
	return string(data), nil
}

// FormatInspectJSON formats metadata as JSON.
func FormatInspectJSON(metadata *ImageMetadata) (string, error) {
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format JSON: %w", err)
	}
	return string(data), nil
}

// FormatInspectHuman formats metadata for human reading.
func FormatInspectHuman(metadata *ImageMetadata) string {
	output := fmt.Sprintf("Name: %s\n", metadata.Name)
	output += fmt.Sprintf("Tag: %s\n", metadata.Tag)
	output += fmt.Sprintf("Digest: %s\n", metadata.ImageDigest)
	output += fmt.Sprintf("Created: %s\n", formatTime(metadata.CreatedAt))
	output += fmt.Sprintf("Files: %d\n", metadata.GetFileCount())

	if len(metadata.SourceImages) > 0 {
		output += "\nSource Images:\n"
		for _, source := range metadata.SourceImages {
			output += fmt.Sprintf("  - %s (resolved: %s)\n", source.Ref, source.ResolvedDigest)
		}
	}

	return output
}

// FormatRecipeOnly returns only the Imagefile content without metadata.
func FormatRecipeOnly(metadata *ImageMetadata) string {
	return metadata.Recipe.Content
}

// FormatFilesList returns file paths as a newline-separated list, sorted alphabetically.
func FormatFilesList(metadata *ImageMetadata) string {
	if len(metadata.Files) == 0 {
		return "No files"
	}

	paths := make([]string, 0, len(metadata.Files))
	for path := range metadata.Files {
		paths = append(paths, path)
	}

	sort.Strings(paths)

	return strings.Join(paths, "\n")
}

func FormatFilesListJSON(metadata *ImageMetadata) (string, error) {
	if len(metadata.Files) == 0 {
		return "[]", nil
	}

	paths := make([]string, 0, len(metadata.Files))
	for path := range metadata.Files {
		paths = append(paths, path)
	}

	sort.Strings(paths)

	data, err := json.MarshalIndent(paths, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format JSON: %w", err)
	}
	return string(data), nil
}

// formatTime formats an ISO 8601 timestamp for human reading.
func formatTime(timestamp string) string {
	// Parse ISO 8601 timestamp
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return timestamp // Return original if parse fails
	}

	// Format as relative or absolute
	now := time.Now()
	duration := now.Sub(t)

	if duration < time.Hour {
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	} else if duration < 7*24*time.Hour {
		return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
	}

	// Older than a week, show date
	return t.Format("2006-01-02")
}
