package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/seregatte/kfg/src/internal/kustomize"
	"github.com/seregatte/kfg/src/internal/logger"
)

var (
	// Flags for the build command
	buildOutput string
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build <path>",
	Short: "Build kustomization and output YAML",
	Long: `Build a kustomization directory and output the resulting YAML.

This command uses kustomize to process a kustomization.yaml file and outputs
the resulting YAML manifests. It supports HTTP resources, strategic merge
patches, and overlays.

Examples:
  kfg build .nixai/overlay/dev
  kfg build .nixai/base -o output.yaml
  kfg build https://example.com/kustomization.yaml`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		// Create kustomize loader
		loader := kustomize.NewLoader(nil)

		// Load kustomization
		resMap, err := loader.Load(path)
		if err != nil {
			logger.Error("build", fmt.Sprintf("Failed to load kustomization: %v", err))
			printBuildError(err, path)
			os.Exit(1)
		}

		// Output YAML
		yamlOutput, err := resMap.AsYaml()
		if err != nil {
			logger.Error("build", fmt.Sprintf("Failed to convert to YAML: %v", err))
			os.Exit(1)
		}

		// Write to output file or stdout
		if buildOutput != "" {
			err = os.WriteFile(buildOutput, yamlOutput, 0644)
			if err != nil {
				logger.Error("build", fmt.Sprintf("Failed to write output file: %v", err))
				os.Exit(1)
			}
			logger.Info("build", fmt.Sprintf("Wrote output to %s", buildOutput))
		} else {
			fmt.Print(string(yamlOutput))
		}
	},
}

// kustomizeCmd is an alias for the build command
var kustomizeCmd = &cobra.Command{
	Use:   "kustomize <path>",
	Short: "Alias for 'build' command",
	Long:  `Alias for the 'build' command. See 'kfg build --help' for details.`,
	Args:  cobra.ExactArgs(1),
	Run:   buildCmd.Run,
}

func init() {
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(kustomizeCmd)

	// Add flags for the build command
	buildCmd.Flags().StringVarP(&buildOutput, "output", "o", "", "Output file path (default: stdout)")

	// Copy flags to kustomize alias
	kustomizeCmd.Flags().StringVarP(&buildOutput, "output", "o", "", "Output file path (default: stdout)")
}

func printBuildError(err error, path string) {
	logger.Error("build", fmt.Sprintf("Failed to build kustomization at %s: %v", path, err))

	// Check for common errors
	if os.IsNotExist(err) {
		logger.Error("build", "The specified path does not exist. Make sure the kustomization directory or file exists.")
	}
}