// Package image implements the image build lifecycle, store persistence,
// and workspace integration for NixAI's image layer system.
package image

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/seregatte/kfg/src/internal/imagefile"
	"github.com/seregatte/kfg/src/internal/logger"
)

// BuildOptions configures the image build process.
type BuildOptions struct {
	Imagefile string // Path to Imagefile (default: "./Imagefile")
	Root      string // Root directory for file resolution (default: current directory)
	Output    string // Output directory for candidate (default: $TMPDIR/<name>/<tag>)
}

// BuildResult represents the result of an image build.
type BuildResult struct {
	Name      string            // Image name from TAG instruction
	Tag       string            // Image tag from TAG instruction
	Digest    string            // SHA256 digest of the build
	Candidate string            // Path to candidate directory
	Recipe    string            // Original Imagefile content
	Files     map[string]string // File manifest (path -> source stage)
}

// Builder constructs configuration images from Imagefile manifests.
type Builder struct {
	options    BuildOptions
	AST        *imagefile.AST
	stages     map[string]*stageContext
	stageDirs  map[string]string // NEW: per-stage directories for isolation
	imageStore *ImageStore
}

// stageContext holds the resolved state of a build stage.
type stageContext struct {
	name    string
	files   map[string]string // path -> content source
	env     map[string]string // environment variables
	workDir string            // working directory (default: "/")
}

// NewBuilder creates a new image builder.
func NewBuilder(options BuildOptions) *Builder {
	// Set defaults
	if options.Imagefile == "" {
		options.Imagefile = "./Imagefile"
	}
	if options.Root == "" {
		options.Root = "."
	}

	return &Builder{
		options:    options,
		stages:     make(map[string]*stageContext),
		imageStore: NewImageStore(""), // Use default store directory
	}
}

// Build executes the build process and returns the result.
func (b *Builder) Build() (*BuildResult, error) {
	// Step 1: Parse Imagefile
	recipeContent, err := b.parseImagefile()
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	// Step 2: Determine base output directory
	baseCandidateDir := b.determineOutputDir()
	logger.Info("build", fmt.Sprintf("Cleaning existing build directory: %s", baseCandidateDir))
	os.RemoveAll(baseCandidateDir)
	if err := os.MkdirAll(baseCandidateDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create candidate directory: %w", err)
	}

	// Initialize stageDirs map for per-stage isolation
	b.stageDirs = make(map[string]string)

	// Step 3: Process each stage in ISOLATED directory
	for _, stage := range b.AST.Stages {
		stageDir := filepath.Join(baseCandidateDir, "stages", stage.Name)
		if err := os.MkdirAll(stageDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create stage directory for %s: %w", stage.Name, err)
		}
		b.stageDirs[stage.Name] = stageDir

		if err := b.processStage(stage, stageDir); err != nil {
			// Leave candidate for inspection on error
			return nil, fmt.Errorf("stage %s failed: %w", stage.Name, err)
		}
	}

	// Step 4: Final stage output becomes the image (copy to artifacts/)
	finalStage := b.AST.Stages[len(b.AST.Stages)-1]
	finalStageDir := b.stageDirs[finalStage.Name]
	finalImageDir := filepath.Join(baseCandidateDir, "artifacts")

	if err := b.copyDir(finalStageDir, finalImageDir); err != nil {
		return nil, fmt.Errorf("failed to prepare final image: %w", err)
	}

	// Step 5: Compute digest from final image only (artifacts directory)
	digest, err := b.computeDigest(recipeContent, finalImageDir)
	if err != nil {
		return nil, fmt.Errorf("digest computation failed: %w", err)
	}

	// Step 6: Extract name and tag from final stage
	name, tag := b.extractNameTag(finalStage)

	// Step 7: Build file manifest from artifacts directory only
	files := b.buildFileManifest(finalImageDir)

	return &BuildResult{
		Name:      name,
		Tag:       tag,
		Digest:    digest,
		Candidate: finalImageDir, // Return artifacts directory as candidate
		Recipe:    recipeContent,
		Files:     files,
	}, nil
}

// parseImagefile reads and parses the Imagefile.
func (b *Builder) parseImagefile() (string, error) {
	imagefilePath := filepath.Join(b.options.Root, b.options.Imagefile)

	file, err := os.Open(imagefilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open Imagefile %s: %w", imagefilePath, err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read Imagefile: %w", err)
	}

	parser := imagefile.NewParser(strings.NewReader(string(content)))
	ast, err := parser.Parse()
	if err != nil {
		return "", err
	}

	b.AST = ast
	return string(content), nil
}

// determineOutputDir determines the candidate output directory.
func (b *Builder) determineOutputDir() string {
	if b.options.Output != "" {
		return b.options.Output
	}

	// Default: $TMPDIR/<name>/<tag>
	tmpdir := os.Getenv("TMPDIR")
	if tmpdir == "" {
		tmpdir = "/tmp"
	}

	// Use placeholder name/tag if not available yet
	name := "unknown"
	tag := "latest"

	if b.AST != nil && len(b.AST.Stages) > 0 {
		finalStage := b.AST.Stages[len(b.AST.Stages)-1]
		if finalStage.Tag != nil {
			name = finalStage.Tag.Name
			tag = finalStage.Tag.Tag
		}
	}

	return filepath.Join(tmpdir, name, tag)
}

// processStage processes a single build stage.
func (b *Builder) processStage(stage *imagefile.Stage, candidateDir string) error {
	logger.Info("build:stage", fmt.Sprintf("Processing stage: %s", stage.Name))

	// Create stage context with default workDir of "" (image root)
	ctx := &stageContext{
		name:    stage.Name,
		files:   make(map[string]string),
		env:     make(map[string]string),
		workDir: "", // Default to image root (empty string after / stripping)
	}

	// Handle FROM instruction
	if stage.From != nil {
		if err := b.handleFrom(stage.From, ctx, candidateDir); err != nil {
			return err
		}
	}

	// Process instructions in order
	for _, instr := range stage.Instructions {
		switch i := instr.(type) {
		case *imagefile.CopyInstruction:
			if err := b.handleCopy(i, ctx, candidateDir); err != nil {
				return err
			}
		case *imagefile.EnvInstruction:
			if err := b.handleEnv(i, ctx); err != nil {
				return err
			}
		case *imagefile.RunInstruction:
			if err := b.handleRun(i, ctx, candidateDir); err != nil {
				return err
			}
		case *imagefile.WorkdirInstruction:
			if err := b.handleWorkdir(i, ctx, candidateDir); err != nil {
				return err
			}
		}
	}

	// Store stage context
	b.stages[stage.Name] = ctx

	return nil
}

// handleFrom handles FROM instruction (stage initialization).
func (b *Builder) handleFrom(from *imagefile.FromInstruction, ctx *stageContext, candidateDir string) error {
	// Build FROM message with optional AS clause
	fromMsg := fmt.Sprintf("FROM %s", from.ImageRef)
	if from.StageName != "" {
		fromMsg = fmt.Sprintf("FROM %s AS %s", from.ImageRef, from.StageName)
	}
	logger.Info("build:from", fromMsg)

	// Handle scratch (empty stage)
	if from.ImageRef == "scratch" {
		logger.Info("build:from", fmt.Sprintf("Stage %s initialized as empty (FROM scratch)", ctx.name))
		return nil
	}

	// Load image from store
	metadata, imageDir, err := b.imageStore.LoadImage(from.ImageRef)
	if err != nil {
		return fmt.Errorf("FROM %s: %w", from.ImageRef, err)
	}

	logger.Info("build:from", fmt.Sprintf("Loaded image %s:%s (digest: %s)", metadata.Name, metadata.Tag, metadata.GetShortDigest()))

	// Copy all files from image directory to candidate
	// This implements "last write wins" for file composition
	err = filepath.Walk(imageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip metadata.json (it's not part of the image content)
		if info.Name() == "metadata.json" {
			return nil
		}

		// Calculate relative path from image directory
		relPath, err := filepath.Rel(imageDir, path)
		if err != nil {
			return err
		}

		// Destination path in candidate
		destPath := filepath.Join(candidateDir, relPath)

		if info.IsDir() {
			// Create directory in candidate
			return os.MkdirAll(destPath, info.Mode())
		}

		// Copy file to candidate
		if err := b.copyFile(path, destPath); err != nil {
			return fmt.Errorf("failed to copy %s: %w", relPath, err)
		}

		// Track file in stage context
		ctx.files[relPath] = fmt.Sprintf("stage:%s:image:%s", ctx.name, from.ImageRef)

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to copy image files: %w", err)
	}

	logger.Info("build:from", fmt.Sprintf("Stage %s initialized with %d files from image %s", ctx.name, len(ctx.files), from.ImageRef))

	return nil
}

// expandSources expands glob patterns in source paths.
// Returns expanded absolute paths and their relative equivalents.
// Returns error if a glob pattern matches zero files.
func expandSources(root string, sources []string) ([]string, []string, error) {
	var absPaths []string
	var relPaths []string

	for _, src := range sources {
		srcPath := filepath.Join(root, src)

		// Check if source contains glob pattern characters
		if strings.ContainsAny(src, "*?[") {
			// Expand glob pattern
			matches, err := filepath.Glob(srcPath)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid glob pattern %s: %w", src, err)
			}
			if len(matches) == 0 {
				return nil, nil, fmt.Errorf("glob pattern %s matched no files", src)
			}

			// Add each matched file
			for _, match := range matches {
				absPaths = append(absPaths, match)
				rel, err := filepath.Rel(root, match)
				if err != nil {
					return nil, nil, fmt.Errorf("failed to get relative path for %s: %w", match, err)
				}
				relPaths = append(relPaths, rel)
			}
		} else {
			// Non-glob source: add as-is
			absPaths = append(absPaths, srcPath)
			relPaths = append(relPaths, src)
		}
	}

	return absPaths, relPaths, nil
}

// determineDestPath resolves the final destination path based on Docker COPY semantics.
// Takes candidateDir, dest path, src path, number of sources, and workDir.
// Returns the final destination path and whether it's a directory destination.
// Handles absolute path normalization by stripping leading / from dest.
// Resolves relative dest paths against workDir (WORKDIR instruction).
func determineDestPath(candidateDir, dest, src string, numSources int, workDir string) (string, bool) {
	// Track if original dest was "." or empty (directory semantics)
	wasCurrentDir := dest == "." || dest == ""

	// Normalize absolute paths: strip leading / to join with candidateDir
	// This handles Docker semantics where / refers to image root, not host filesystem
	if strings.HasPrefix(dest, "/") {
		dest = strings.TrimPrefix(dest, "/")
	} else if workDir != "" {
		// Relative destination: resolve against WORKDIR
		// Special case: "." means current directory in workdir context
		if dest == "." || dest == "" {
			dest = workDir
		} else {
			dest = filepath.Join(workDir, dest)
		}
	}

	// Handle edge case where dest was "//" or just "//" - becomes empty after strip
	// Treat empty dest as current directory (".")
	if dest == "" || dest == "/" {
		dest = "."
	}

	// Check if dest ends with "/" (explicit directory marker)
	isExplicitDir := strings.HasSuffix(dest, "/")

	// Preserve directory semantics if original dest was "." or empty
	// This ensures COPY file.txt . copies INTO the directory, not TO the directory path
	isCurrentDir := wasCurrentDir || dest == "."

	// Multiple sources always force directory semantics
	isMultipleSources := numSources > 1

	// Determine if destination should be treated as directory
	isDirDest := isExplicitDir || isCurrentDir || isMultipleSources

	// Resolve base destination path
	baseDest := filepath.Join(candidateDir, strings.TrimSuffix(dest, "/"))

	// If directory destination, append source filename
	if isDirDest {
		baseDest = filepath.Join(baseDest, filepath.Base(src))
	}

	return baseDest, isDirDest
}

// copyPathWithDestInfo copies a file or directory using resolved destination info.
func (b *Builder) copyPathWithDestInfo(src, dest string, isDirDest bool) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		return b.copyDir(src, dest)
	}

	return b.copyFile(src, dest)
}

// handleCopy handles COPY instruction.
func (b *Builder) handleCopy(copy *imagefile.CopyInstruction, ctx *stageContext, candidateDir string) error {
	// Build COPY message with optional --from clause
	copyMsg := fmt.Sprintf("COPY %v -> %s", copy.Sources, copy.Dest)
	if copy.FromStage != "" {
		copyMsg = fmt.Sprintf("COPY %v -> %s (from stage: %s)", copy.Sources, copy.Dest, copy.FromStage)
	}
	logger.Info("build:copy", copyMsg)

	// Handle --from=<stage> (copy from another stage)
	if copy.FromStage != "" {
		sourceStageDir, ok := b.stageDirs[copy.FromStage]
		if !ok {
			return fmt.Errorf("COPY --from=%s: stage not found", copy.FromStage)
		}

		// Expand glob patterns in sources from source stage's directory
		absPaths, relPaths, err := expandSources(sourceStageDir, copy.Sources)
		if err != nil {
			return fmt.Errorf("COPY --from=%s: %w", copy.FromStage, err)
		}

		// Copy expanded files from SOURCE stage's isolated directory to current stage
		for i, srcPath := range absPaths {
			relPath := relPaths[i]

			// Verify source exists in the source stage directory
			if _, err := os.Stat(srcPath); err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("COPY --from=%s: file %s not found in stage %s", copy.FromStage, relPath, copy.FromStage)
				}
				return fmt.Errorf("COPY --from=%s: failed to access %s: %w", copy.FromStage, relPath, err)
			}

			// Resolve destination path in CURRENT stage's directory using Docker semantics
			// Pass total expanded sources count for directory semantics
			destPath, isDirDest := determineDestPath(candidateDir, copy.Dest, relPath, len(absPaths), ctx.workDir)

			// Ensure destination directory exists
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return fmt.Errorf("COPY --from=%s: failed to create destination directory: %w", copy.FromStage, err)
			}

			// Copy file or directory
			if err := b.copyPathWithDestInfo(srcPath, destPath, isDirDest); err != nil {
				return fmt.Errorf("COPY --from=%s: failed to copy %s: %w", copy.FromStage, relPath, err)
			}

			// Track file in current stage context
			ctx.files[copy.Dest] = fmt.Sprintf("stage:%s:copy-from:%s", ctx.name, copy.FromStage)
			logger.Info("build:copy", fmt.Sprintf("Copied %s from stage %s to %s", relPath, copy.FromStage, copy.Dest))
		}

		return nil
	}

	// Copy from workspace (default behavior)
	// Expand glob patterns in sources from workspace root
	absPaths, relPaths, err := expandSources(b.options.Root, copy.Sources)
	if err != nil {
		return fmt.Errorf("COPY: %w", err)
	}

	// Copy expanded files from workspace to current stage
	for i, srcPath := range absPaths {
		relPath := relPaths[i]

		// Resolve destination path using Docker semantics
		// Pass total expanded sources count for directory semantics
		destPath, isDirDest := determineDestPath(candidateDir, copy.Dest, relPath, len(absPaths), ctx.workDir)

		// Ensure destination directory exists
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("COPY: failed to create destination directory: %w", err)
		}

		// Copy file or directory
		if err := b.copyPathWithDestInfo(srcPath, destPath, isDirDest); err != nil {
			return fmt.Errorf("COPY %s to %s failed: %w", relPath, copy.Dest, err)
		}

		// Record in file manifest
		ctx.files[copy.Dest] = fmt.Sprintf("stage:%s:workspace", ctx.name)
	}

	return nil
}

// copyPath copies a file or directory from source to destination.
func (b *Builder) copyPath(src, dest string) error {
	// Check if source is a directory
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	if info.IsDir() {
		// Copy directory recursively
		return b.copyDir(src, dest)
	}

	// Copy single file
	return b.copyFile(src, dest)
}

// copyFile copies a single file preserving permissions.
func (b *Builder) copyFile(src, dest string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Get source file info for permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Create destination file with same permissions
	destFile, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy content
	_, err = io.Copy(destFile, srcFile)
	return err
}

// copyDir copies a directory recursively.
func (b *Builder) copyDir(src, dest string) error {
	// Ensure destination directory exists
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	// Walk source directory
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dest, relPath)

		// Create directory or copy file
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return b.copyFile(path, destPath)
	})

	return err
}

// handleEnv handles ENV instruction.
func (b *Builder) handleEnv(env *imagefile.EnvInstruction, ctx *stageContext) error {
	logger.Info("build:env", fmt.Sprintf("ENV %v", env.Vars))

	// Merge environment variables
	for key, value := range env.Vars {
		ctx.env[key] = value
	}

	return nil
}

// handleWorkdir handles WORKDIR instruction.
func (b *Builder) handleWorkdir(workdir *imagefile.WorkdirInstruction, ctx *stageContext, candidateDir string) error {
	logger.Info("build:workdir", fmt.Sprintf("WORKDIR %s", workdir.Path))

	// Normalize the path: strip leading / to make it relative to candidateDir
	workPath := workdir.Path
	if strings.HasPrefix(workPath, "/") {
		// Absolute path: normalize and set directly
		workPath = strings.TrimPrefix(workPath, "/")
		ctx.workDir = workPath
	} else {
		// Relative path: combine with previous workDir if set
		if ctx.workDir != "" {
			ctx.workDir = filepath.Join(ctx.workDir, workPath)
		} else {
			ctx.workDir = workPath
		}
	}

	// Create the working directory inside candidateDir if it doesn't exist
	fullWorkDir := filepath.Join(candidateDir, ctx.workDir)
	if err := os.MkdirAll(fullWorkDir, 0755); err != nil {
		return fmt.Errorf("WORKDIR: failed to create directory %s: %w", ctx.workDir, err)
	}

	logger.Info("build:workdir", fmt.Sprintf("Working directory set to %s", ctx.workDir))
	return nil
}

// handleRun handles RUN instruction.
func (b *Builder) handleRun(run *imagefile.RunInstruction, ctx *stageContext, candidateDir string) error {
	logger.Info("build:run", fmt.Sprintf("RUN %s", run.Command))

	// Resolve working directory for RUN command
	// Use workDir if set, otherwise use candidateDir (image root)
	var runDir string
	if ctx.workDir != "" {
		runDir = filepath.Join(candidateDir, ctx.workDir)
	} else {
		runDir = candidateDir
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command("sh", "-c", run.Command)
	cmd.Dir = runDir
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Set environment variables
	cmd.Env = os.Environ() // Start with current environment
	for key, value := range ctx.env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	// Run command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("RUN command failed: %w", err)
	}

	// Log stdout lines
	if stdoutBuf.Len() > 0 {
		for _, line := range strings.Split(stdoutBuf.String(), "\n") {
			if line != "" {
				logger.LogWithSource("detail", "go", "build:run:stdout", line, nil)
			}
		}
	}

	// Log stderr lines
	if stderrBuf.Len() > 0 {
		for _, line := range strings.Split(stderrBuf.String(), "\n") {
			if line != "" {
				logger.LogWithSource("warn", "go", "build:run:stderr", line, nil)
			}
		}
	}

	logger.Info("build:run", "RUN completed successfully")
	return nil
}

// computeDigest computes SHA256 digest of recipe + sorted files.
func (b *Builder) computeDigest(recipe string, candidateDir string) (string, error) {
	hasher := sha256.New()

	// Add recipe content
	hasher.Write([]byte(recipe))

	// Collect and sort all file paths
	var files []string
	err := filepath.Walk(candidateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(candidateDir, path)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	// Sort files for deterministic hash
	sort.Strings(files)

	// Add each file content to hash
	for _, file := range files {
		filePath := filepath.Join(candidateDir, file)
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", err
		}
		hasher.Write([]byte(file))
		hasher.Write(data)
	}

	// Return hex-encoded digest
	digest := hex.EncodeToString(hasher.Sum(nil))
	return fmt.Sprintf("sha256:%s", digest), nil
}

// extractNameTag extracts name and tag from final stage.
func (b *Builder) extractNameTag(stage *imagefile.Stage) (string, string) {
	if stage.Tag != nil {
		return stage.Tag.Name, stage.Tag.Tag
	}

	// Default values if no TAG instruction
	return "unknown", "latest"
}

// buildFileManifest creates a manifest of all files in the candidate.
func (b *Builder) buildFileManifest(candidateDir string) map[string]string {
	files := make(map[string]string)

	filepath.Walk(candidateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(candidateDir, path)
			if err != nil {
				return err
			}
			files[relPath] = "build"
		}
		return nil
	})

	return files
}
