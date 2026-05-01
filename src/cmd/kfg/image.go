package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/seregatte/kfg/src/internal/image"
	"github.com/seregatte/kfg/src/internal/logger"
)

var imageFile string
var imageRoot string
var imageOutput string

var keepBuild bool
var imageBuildPush bool

var imageListJSON bool

// storeDirOverride is shared across image and workspace commands for store directory override
var storeDirOverride string

var imageInspectJSON bool
var imageInspectRecipe bool
var imageInspectFiles bool

// imageCmd is the parent command for image operations
var imageCmd = &cobra.Command{
	Use:     "image",
	Aliases: []string{"img"},
	Short:   "Manage stored configuration images",
	Long: `Manage stored configuration images in the KFG store.

Images are immutable configuration snapshots built from Imagefiles.
Each image has a name, tag, and SHA256 digest for identification.

Subcommands:
  build    Build an image from an Imagefile
  push     Persist a built image to the store
  list     List all stored images
  inspect  Display image metadata
  remove   Delete an image from the store

Examples:
  kfg image build
  kfg img push /tmp/.kfg/build/myimage/v1
  kfg image list
  kfg img inspect claude-base:v2 --recipe
  kfg image remove old-config:latest`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// imageBuildCmd builds an image from an Imagefile
var imageBuildCmd = &cobra.Command{
	Use:   "build [-f <path>] [--root <dir>] [--output <dir>]",
	Short: "Build an image from an Imagefile",
	Long: `Build a local candidate image from an Imagefile.

The build process:
1. Parses the Imagefile
2. Resolves FROM stages (loads images from store)
3. Copies files from stages and workspace
4. Executes RUN commands
5. Computes SHA256 digest
6. Outputs candidate to build directory

The candidate is NOT automatically pushed to the store unless --push is used.
Use 'kfg image push' to persist the built image.

Flags:
  -f, --file <path>      Imagefile path (default: ./Imagefile)
  --root <dir>           Root directory for file resolution (default: current directory)
  --output <dir>         Output directory for candidate (default: $TMPDIR/.kfg/build/<name>/<tag>)
  --push                 Automatically push image after successful build
  --keep-build           Preserve build directory after push (use with --push)

Examples:
  kfg image build
  kfg image build -f ./Custom.dockerfile
  kfg img build --root /path/to/project --output /tmp/mybuild
  kfg image build --push
  kfg img build --push --keep-build`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Set defaults
		imagefile := imageFile
		if imagefile == "" {
			imagefile = "./Imagefile"
		}

		root := imageRoot
		if root == "" {
			root = "."
		}

		// Create build options
		options := image.BuildOptions{
			Imagefile: imagefile,
			Root:      root,
			Output:    imageOutput,
		}

		// Create builder
		builder := image.NewBuilder(options)

		// Execute build
		result, err := builder.Build()
		if err != nil {
			logger.Error("image:build", fmt.Sprintf("Build failed: %v", err))
			os.Exit(1)
		}

		// Create metadata from build result
		metadata := image.NewMetadata(result.Name, result.Tag)
		metadata.SetDigest(result.Digest)
		metadata.SetRecipe(imagefile, result.Recipe)

		// Add files to manifest
		for path, source := range result.Files {
			metadata.AddFile(path, source)
		}

		// Validate metadata
		validation := metadata.Validate()
		if !validation.IsValid() {
			logger.Error("image:build", fmt.Sprintf("Metadata validation failed: %s", validation.Error()))
			os.Exit(1)
		}

		// Save metadata to candidate directory
		if err := metadata.SaveMetadataToDir(result.Candidate); err != nil {
			logger.Error("image:build", fmt.Sprintf("Failed to save metadata: %v", err))
			os.Exit(1)
		}

		// Report success
		logger.Info("image:build", fmt.Sprintf("Build succeeded: %s:%s", result.Name, result.Tag))
		fmt.Printf("Candidate directory: %s\n", result.Candidate)
		fmt.Printf("Digest: %s\n", result.Digest)
		fmt.Printf("Files: %d\n", len(result.Files))

		if imageBuildPush {
			store := image.NewImageStore(getStoreDir())
			err := store.PushImage(result.Candidate, keepBuild)
			if err != nil {
				logger.Error("image:build", fmt.Sprintf("Push failed: %v", err))
				os.Exit(1)
			}
			logger.Info("image:build", fmt.Sprintf("Image pushed successfully: %s:%s", result.Name, result.Tag))
		} else {
			fmt.Printf("\nTo persist this image, run:\n  kfg image push %s\n", result.Candidate)
		}
	},
}

// imagePushCmd persists a built image to the store
var imagePushCmd = &cobra.Command{
	Use:   "push <build-dir> [--keep-build]",
	Short: "Persist a built image to the store",
	Long: `Persist a built image candidate to the store.

The candidate directory should contain the image files and metadata.json.
After push, the image becomes immutable and can be referenced by name:tag.

By default, the build directory is cleaned up after successful push.
Use --keep-build to preserve the candidate for inspection.

Images are immutable: pushing an image with an existing name:tag fails.

Flags:
  --keep-build    Preserve build directory after push (default: false)

Examples:
  kfg image push /tmp/.kfg/build/myimage/v1
  kfg img push /tmp/.kfg/build/myimage/v1 --keep-build`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		buildDir := args[0]

		// Check if build directory exists
		if _, err := os.Stat(buildDir); err != nil {
			if os.IsNotExist(err) {
				logger.Error("image:push", fmt.Sprintf("Build directory not found: %s", buildDir))
				os.Exit(1)
			}
			logger.Error("image:push", fmt.Sprintf("Failed to access build directory: %v", err))
			os.Exit(1)
		}

		// Create image store
		store := image.NewImageStore(getStoreDir())

		// Push image
		err := store.PushImage(buildDir, keepBuild)
		if err != nil {
			logger.Error("image:push", fmt.Sprintf("Push failed: %v", err))
			os.Exit(1)
		}

		logger.Info("image:push", fmt.Sprintf("Image pushed successfully from: %s", buildDir))
	},
}

// imageListCmd lists all stored images
var imageListCmd = &cobra.Command{
	Use:   "list [--json]",
	Short: "List all stored images",
	Long: `List all stored images in the KFG store.

Shows: NAME, TAG, DIGEST (first 12 chars), CREATED, FILES.

With --json flag, outputs JSON array of image objects.

Examples:
  kfg image list
  kfg img ls --json`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		store := image.NewImageStore(getStoreDir())

		images, err := store.ListImages()
		if err != nil {
			logger.Error("image:list", fmt.Sprintf("Failed to list images: %v", err))
			os.Exit(1)
		}

		if imageListJSON {
			jsonOutput, err := image.FormatListJSON(images)
			if err != nil {
				logger.Error("image:list", fmt.Sprintf("Failed to format JSON: %v", err))
				os.Exit(1)
			}
			fmt.Println(jsonOutput)
			return
		}

		if len(images) == 0 {
			fmt.Println("No images found")
			return
		}

		tableOutput := image.FormatListTable(images)
		fmt.Println(tableOutput)
	},
}

// imageInspectCmd displays image metadata
var imageInspectCmd = &cobra.Command{
	Use:   "inspect <name[:tag]> [--json] [--recipe] [--files]",
	Short: "Display image metadata",
	Long: `Display detailed metadata for a stored image.

Shows: name, tag, digest, created date, file count, source images.

With --json flag, outputs full metadata as JSON.
With --recipe flag, outputs only the original Imagefile content (no metadata).
With --files flag, outputs only the file paths (one per line, sorted alphabetically).

If tag is omitted, defaults to :latest.

Examples:
  kfg image inspect claude-base:v2
  kfg img inspect myconfig --json
  kfg image inspect claude-base:v2 --recipe
  kfg img inspect claude-base:v2 --files`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ref := args[0]

		if imageInspectRecipe && imageInspectFiles {
			logger.Error("image:inspect", "flags --recipe and --files are mutually exclusive")
			os.Exit(1)
		}

		if imageInspectRecipe && imageInspectJSON {
			logger.Error("image:inspect", "flags --recipe and --json are mutually exclusive")
			os.Exit(1)
		}

		store := image.NewImageStore(getStoreDir())

		metadata, err := store.InspectImage(ref)
		if err != nil {
			logger.Error("image:inspect", fmt.Sprintf("Failed to inspect image: %v", err))
			os.Exit(1)
		}

		if imageInspectFiles {
			if imageInspectJSON {
				jsonOutput, err := image.FormatFilesListJSON(metadata)
				if err != nil {
					logger.Error("image:inspect", fmt.Sprintf("Failed to format JSON: %v", err))
					os.Exit(1)
				}
				fmt.Println(jsonOutput)
			} else {
				output := image.FormatFilesList(metadata)
				fmt.Println(output)
			}
		} else if imageInspectJSON {
			jsonOutput, err := image.FormatInspectJSON(metadata)
			if err != nil {
				logger.Error("image:inspect", fmt.Sprintf("Failed to format JSON: %v", err))
				os.Exit(1)
			}
			fmt.Println(jsonOutput)
		} else if imageInspectRecipe {
			output := image.FormatRecipeOnly(metadata)
			fmt.Println(output)
		} else {
			output := image.FormatInspectHuman(metadata)
			fmt.Println(output)
		}
	},
}

// imageRemoveCmd deletes an image from the store
var imageRemoveCmd = &cobra.Command{
	Use:   "remove <name[:tag]>",
	Short: "Delete an image from the store",
	Long: `Delete an image from the KFG store.

Alias: rm

If tag is omitted, defaults to :latest.

Examples:
  kfg image remove old-config:v1
  kfg img rm myimage:latest`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ref := args[0]

		// Create image store
		store := image.NewImageStore(getStoreDir())

		// Remove image
		err := store.RemoveImage(ref)
		if err != nil {
			logger.Error("image:remove", fmt.Sprintf("Failed to remove image: %v", err))
			os.Exit(1)
		}

		logger.Info("image:remove", fmt.Sprintf("Image removed: %s", ref))
	},
}

// imageRmCmd is an alias for remove
var imageRmCmd = &cobra.Command{
	Use:   "rm <name[:tag]>",
	Short: "Alias for 'remove'",
	Long:  `Alias for 'image remove'. See 'kfg image remove --help' for details.`,
	Args:  cobra.ExactArgs(1),
	Run:   imageRemoveCmd.Run,
}

func init() {
	// Add imageCmd directly to rootCmd (promoted from store)
	rootCmd.AddCommand(imageCmd)

	// Persistent flags for image command (available to all subcommands)
	imageCmd.PersistentFlags().StringVar(&storeDirOverride, "store", "", "Store directory override (default: $KFG_STORE_DIR or ~/.config/kfg/store)")

	// Add image subcommands to imageCmd
	imageCmd.AddCommand(imageBuildCmd)
	imageCmd.AddCommand(imagePushCmd)
	imageCmd.AddCommand(imageListCmd)
	imageCmd.AddCommand(imageInspectCmd)
	imageCmd.AddCommand(imageRemoveCmd)
	imageCmd.AddCommand(imageRmCmd)

	// Flags for build command
	imageBuildCmd.Flags().StringVarP(&imageFile, "file", "f", "", "Imagefile path (default: ./Imagefile)")
	imageBuildCmd.Flags().StringVar(&imageRoot, "root", "", "Root directory for file resolution (default: current directory)")
	imageBuildCmd.Flags().StringVar(&imageOutput, "output", "", "Output directory for candidate (default: $TMPDIR/.kfg/build/<name>/<tag>)")
	imageBuildCmd.Flags().BoolVar(&imageBuildPush, "push", false, "Automatically push image after successful build")
	imageBuildCmd.Flags().BoolVar(&keepBuild, "keep-build", false, "Preserve build directory after push (use with --push)")

	// Flags for push command
	imagePushCmd.Flags().BoolVar(&keepBuild, "keep-build", false, "Preserve build directory after push")

	// Flags for list command
	imageListCmd.Flags().BoolVar(&imageListJSON, "json", false, "Output in JSON format")

	// Flags for inspect command
	imageInspectCmd.Flags().BoolVar(&imageInspectJSON, "json", false, "Output in JSON format")
	imageInspectCmd.Flags().BoolVar(&imageInspectRecipe, "recipe", false, "Output only the original Imagefile content")
	imageInspectCmd.Flags().BoolVar(&imageInspectFiles, "files", false, "Output only the file paths (one per line)")
}

// getStoreDir returns the store directory path (shared with v1 store commands)
func getStoreDir() string {
	// Use the same storeDirOverride logic from v1 store commands
	if storeDirOverride != "" {
		return storeDirOverride
	}
	return "" // Will use default from ImageStore constructor
}
