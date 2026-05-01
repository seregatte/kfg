package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/seregatte/kfg/src/internal/image"
	"github.com/seregatte/kfg/src/internal/logger"
)

var startImageRef string
var startName string
var startRoot string

var stopName string

// workspaceCmd is the parent command for workspace operations
var workspaceCmd = &cobra.Command{
	Use:     "workspace",
	Aliases: []string{"ws"},
	Short:   "Manage workspace instances from stored images",
	Long: `Manage workspace instances by materializing and restoring stored images.

Workspace operations allow you to:
  - Materialize stored images into your workspace (start)
  - Restore workspace from automatic backups (stop)

The workspace instance system provides:
  - Automatic backup before materialization (prevents data loss)
  - Named instances for tracking multiple projects
  - Idempotent stop (succeeds even if backup missing)

Subcommands:
  start  Materialize an image into the workspace
  stop   Restore workspace from backup and cleanup instance

Examples:
  kfg workspace start claude-base:v2 --name myproject
  kfg ws stop --name myproject`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// startCmd materializes an image into the workspace
var workspaceStartCmd = &cobra.Command{
	Use:   "start <image[:tag]> --name <name> [--root <dir>]",
	Short: "Materialize an image into the workspace",
	Long: `Materialize a stored image into the workspace with automatic backup.

This command:
1. Validates the image exists in store
2. Creates backup of existing workspace files (if any)
3. Copies image files to workspace
4. Creates instance record for later cleanup

The --name flag is required and uniquely identifies this instance.
Use the same name with 'stop' to restore from backup.

The --root flag specifies the workspace directory (default: current directory).

Backup behavior:
  - Empty workspace: No backup created, materialization proceeds
  - Non-empty workspace: Backup created before materialization
  - Repeated start: Second start overwrites previous backup

Instance tracking:
  - Instance metadata stored in $KFG_STORE_DIR/.workspace/<name>/
  - Records: instance name, image ref, timestamp, workspace root
  - Instance names must be unique (global namespace)

If tag is omitted, defaults to :latest.

Examples:
  kfg workspace start claude-base:v2 --name myproject
  kfg ws start my-config --name proj1 --root /path/to/project
  kfg workspace start opencode:latest --name opencode-default`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		imageRef := args[0]

		// Validate required flags
		if startName == "" {
			logger.Error("workspace:start", "Instance name is required (--name)")
			os.Exit(1)
		}

		// Set default root
		root := startRoot
		if root == "" {
			root = "."
		}

		// Create materializer
		materializer := image.NewMaterializer(getStoreDir())

		// Execute start
		err := materializer.Start(imageRef, root, startName)
		if err != nil {
			logger.Error("workspace:start", fmt.Sprintf("Failed to start: %v", err))
			os.Exit(1)
		}
	},
}

// stopCmd restores workspace from backup
var workspaceStopCmd = &cobra.Command{
	Use:   "stop --name <name>",
	Short: "Restore workspace from backup and cleanup instance",
	Long: `Restore workspace from automatic backup and cleanup instance.

This command:
1. Locates instance by name
2. Restores workspace from backup (if exists)
3. Removes instance record and backup

The --name flag identifies which instance to stop (must match a previous start).

Idempotent behavior:
  - If instance not found: succeeds with message (no error)
  - If backup missing: succeeds, workspace left as-is

After stop:
  - Backup is consumed (deleted after restoration)
  - Instance metadata removed
  - Workspace restored to pre-start state

Examples:
  kfg workspace stop --name myproject
  kfg ws stop --name proj1`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Validate required flags
		if stopName == "" {
			logger.Error("workspace:stop", "Instance name is required (--name)")
			os.Exit(1)
		}

		// Create materializer
		materializer := image.NewMaterializer(getStoreDir())

		// Execute stop
		err := materializer.Stop(stopName)
		if err != nil {
			logger.Error("workspace:stop", fmt.Sprintf("Failed to stop: %v", err))
			os.Exit(1)
		}
	},
}

func init() {
	// Add workspaceCmd directly to rootCmd (promoted from store)
	rootCmd.AddCommand(workspaceCmd)

	// Persistent flags for workspace command (available to all subcommands)
	workspaceCmd.PersistentFlags().StringVar(&storeDirOverride, "store", "", "Store directory override (default: $KFG_STORE_DIR or ~/.config/kfg/store)")

	// Add subcommands to workspaceCmd
	workspaceCmd.AddCommand(workspaceStartCmd)
	workspaceCmd.AddCommand(workspaceStopCmd)

	// Flags for start command
	workspaceStartCmd.Flags().StringVar(&startName, "name", "", "Instance name (required, unique identifier)")
	workspaceStartCmd.Flags().StringVar(&startRoot, "root", "", "Workspace directory (default: current directory)")
	workspaceStartCmd.MarkFlagRequired("name")

	// Flags for stop command
	workspaceStopCmd.Flags().StringVar(&stopName, "name", "", "Instance name to stop (required)")
	workspaceStopCmd.MarkFlagRequired("name")
}