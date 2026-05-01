package kustomize

import (
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
	"sigs.k8s.io/kustomize/kyaml/openapi"
)

// Loader wraps kustomize's krusty.Kustomizer for loading kustomizations.
// It integrates NixAI's embedded OpenAPI schema to enable strategic merge
// patches on custom resource types.
type Loader struct {
	kustomizer *krusty.Kustomizer
	fSys       filesys.FileSystem
}

// LoaderOptions configures the kustomize loader.
type LoaderOptions struct {
	// LoadRestrictions controls how resources are loaded.
	// Default is LoadRestrictionsNone (allows any path).
	LoadRestrictions types.LoadRestrictions
}

// DefaultLoaderOptions returns the default loader options.
func DefaultLoaderOptions() *LoaderOptions {
	return &LoaderOptions{
		LoadRestrictions: types.LoadRestrictionsNone,
	}
}

// NewLoader creates a new kustomize loader with the embedded OpenAPI schema.
// The loader uses the on-disk filesystem by default.
func NewLoader(opts *LoaderOptions) *Loader {
	if opts == nil {
		opts = DefaultLoaderOptions()
	}

	// Initialize OpenAPI schema for NixAI types before creating the kustomizer.
	// This enables strategic merge patches to work correctly on custom arrays.
	initOpenAPI()

	// Create filesystem abstraction
	fSys := filesys.MakeFsOnDisk()

	// Configure kustomizer options
	kOpts := &krusty.Options{
		LoadRestrictions: opts.LoadRestrictions,
		PluginConfig:     types.DisabledPluginConfig(),
	}

	// Create the kustomizer
	k := krusty.MakeKustomizer(kOpts)

	return &Loader{
		kustomizer: k,
		fSys:       fSys,
	}
}

// Load reads a kustomization from the given path and returns the merged ResMap.
// The path should point to a directory containing a kustomization.yaml file,
// or to the kustomization.yaml file directly.
// Uses the default filesystem (on-disk).
//
// Example:
//
//	loader := NewLoader(nil)
//	resMap, err := loader.Load(".nixai/overlay/dev")
func (l *Loader) Load(path string) (resmap.ResMap, error) {
	return l.kustomizer.Run(l.fSys, path)
}

// LoadWithFS reads a kustomization from the given path using a custom filesystem.
// This is useful for testing with in-memory filesystems.
//
// Example:
//
//	fSys := filesys.MakeFsInMemory()
//	fSys.WriteFile("/test/kustomization.yaml", []byte("..."))
//	resMap, err := loader.LoadWithFS(fSys, "/test")
func (l *Loader) LoadWithFS(fSys filesys.FileSystem, path string) (resmap.ResMap, error) {
	return l.kustomizer.Run(fSys, path)
}

// initOpenAPI initializes the OpenAPI schema with NixAI custom type definitions.
// This must be called before running the kustomizer to enable strategic merge
// patches on arrays like CmdWorkflow's before/after steps.
func initOpenAPI() {
	// The openapi package manages schema definitions for custom resources.
	// We add NixAI's schema so kustomize knows how to merge our custom arrays.
	openapi.AddSchema(GetOpenAPISchema())
}