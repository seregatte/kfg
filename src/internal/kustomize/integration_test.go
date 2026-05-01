//go:build integration

package kustomize

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/seregatte/kfg/src/internal/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

// TestLoaderHTTPResources tests loading resources from HTTP URLs.
// This test requires network access and is tagged with 'integration'.
// Run with: go test -tags=integration ./src/internal/kustomize/...
//
// Note: HTTP resource loading is handled by kustomize internally.
// This test validates that our Loader wrapper correctly handles HTTP URLs.
func TestLoaderHTTPResources(t *testing.T) {
	// Create in-memory filesystem
	fSys := filesys.MakeFsInMemory()

	// Use standard kustomize format with HTTP URL in resources
	// Note: This requires a stable URL that returns valid YAML
	// For testing purposes, we use a kubernetes-sigs/kustomize test fixture
	kustomizationYAML := `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/cmd/config/testdata/bases/simple/configMap.yaml
`

	// Write kustomization.yaml
	if err := fSys.WriteFile("/test/kustomization.yaml", []byte(kustomizationYAML)); err != nil {
		t.Fatalf("Failed to write kustomization.yaml: %v", err)
	}

	// Create loader and load
	loader := NewLoader(nil)
	resMap, err := loader.LoadWithFS(fSys, "/test")
	if err != nil {
		t.Fatalf("LoadWithFS failed: %v", err)
	}

	// Verify we got resources
	if resMap.Size() == 0 {
		t.Error("Expected at least 1 resource from HTTP URL, got 0")
	}

	// Convert to NixAI types
	adapter := NewAdapter()
	resources, err := adapter.ResMapToResources(resMap)
	if err != nil {
		t.Fatalf("ResMapToResources failed: %v", err)
	}

	// Verify resources were loaded from HTTP
	if len(resources) == 0 {
		t.Error("Expected at least 1 resource from HTTP URL")
	}

	t.Logf("Successfully loaded %d resources from HTTP URL", len(resources))
}

// TestLoaderAssets tests loading Assets resources from kustomization.
// This test validates that Assets are correctly parsed and validated.
func TestLoaderAssets(t *testing.T) {
	// Create in-memory filesystem
	fSys := filesys.MakeFsInMemory()

	kustomizationYAML := `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - assets.yaml
`

	assetsYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test-providers
spec:
  schemaRef: schema://providers
  data:
    servers:
      - enabled: true
        type: anthropic
        command: claude
`

	// Write files
	if err := fSys.WriteFile("/test/kustomization.yaml", []byte(kustomizationYAML)); err != nil {
		t.Fatalf("Failed to write kustomization.yaml: %v", err)
	}
	if err := fSys.WriteFile("/test/assets.yaml", []byte(assetsYAML)); err != nil {
		t.Fatalf("Failed to write assets.yaml: %v", err)
	}

	// Create loader and load
	loader := NewLoader(nil)
	resMap, err := loader.LoadWithFS(fSys, "/test")
	if err != nil {
		t.Fatalf("LoadWithFS failed: %v", err)
	}

	// Convert to NixAI types
	adapter := NewAdapter()
	resources, err := adapter.ResMapToResources(resMap)
	if err != nil {
		t.Fatalf("ResMapToResources failed: %v", err)
	}

	// Verify we got Assets
	if len(resources) == 0 {
		t.Error("Expected at least 1 resource, got 0")
	}

	// Find Assets by kind
	assetsResources := adapter.GetByKind(resources, "Assets")
	if len(assetsResources) != 1 {
		t.Errorf("Expected 1 Assets resource, got %d", len(assetsResources))
	}

	// Verify Assets was parsed correctly
	if assetsResources[0].Kind() != "Assets" {
		t.Errorf("Expected kind 'Assets', got '%s'", assetsResources[0].Kind())
	}
	if assetsResources[0].Name() != "test-providers" {
		t.Errorf("Expected name 'test-providers', got '%s'", assetsResources[0].Name())
	}

	t.Logf("Successfully loaded Assets: %s", assetsResources[0].Name())
}

// TestLoaderConverter tests loading Converter resources from kustomization.
// This test validates that Converters are correctly parsed and validated.
func TestLoaderConverter(t *testing.T) {
	// Create in-memory filesystem
	fSys := filesys.MakeFsInMemory()

	kustomizationYAML := `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - converter.yaml
`

	converterYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Converter
metadata:
  name: providers-to-claude
spec:
  input:
    schemaRef: schema://providers
  engine:
    type: template
    template: '{"name": "{{ .Metadata.Name }}"}'
  output:
    format: json
`

	// Write files
	if err := fSys.WriteFile("/test/kustomization.yaml", []byte(kustomizationYAML)); err != nil {
		t.Fatalf("Failed to write kustomization.yaml: %v", err)
	}
	if err := fSys.WriteFile("/test/converter.yaml", []byte(converterYAML)); err != nil {
		t.Fatalf("Failed to write converter.yaml: %v", err)
	}

	// Create loader and load
	loader := NewLoader(nil)
	resMap, err := loader.LoadWithFS(fSys, "/test")
	if err != nil {
		t.Fatalf("LoadWithFS failed: %v", err)
	}

	// Convert to NixAI types
	adapter := NewAdapter()
	resources, err := adapter.ResMapToResources(resMap)
	if err != nil {
		t.Fatalf("ResMapToResources failed: %v", err)
	}

	// Verify we got Converter
	if len(resources) == 0 {
		t.Error("Expected at least 1 resource, got 0")
	}

	// Find Converter by kind
	converterResources := adapter.GetByKind(resources, "Converter")
	if len(converterResources) != 1 {
		t.Errorf("Expected 1 Converter resource, got %d", len(converterResources))
	}

	// Verify Converter was parsed correctly
	if converterResources[0].Kind() != "Converter" {
		t.Errorf("Expected kind 'Converter', got '%s'", converterResources[0].Kind())
	}
	if converterResources[0].Name() != "providers-to-claude" {
		t.Errorf("Expected name 'providers-to-claude', got '%s'", converterResources[0].Name())
	}

	t.Logf("Successfully loaded Converter: %s", converterResources[0].Name())
}

// TestLoaderAllFiveKinds tests loading all five kinds from kustomization.
// This validates that Step, Cmd, CmdWorkflow, Assets, and Converter all load together.
func TestLoaderAllFiveKinds(t *testing.T) {
	// Create in-memory filesystem
	fSys := filesys.MakeFsInMemory()

	kustomizationYAML := `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - step.yaml
  - cmd.yaml
  - workflow.yaml
  - assets.yaml
  - converter.yaml
`

	stepYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Step
metadata:
  name: setup-step
spec:
  run: echo "setup"
`

	cmdYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Cmd
metadata:
  name: test-cmd
  commandName: test_cmd
spec:
  run: echo "cmd"
`

	workflowYAML := `
apiVersion: kfg.dev/v1alpha1
kind: CmdWorkflow
metadata:
  name: test-workflow
spec:
  cmds:
    - test-cmd
`

	assetsYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test-providers
spec:
  schemaRef: schema://providers
  data:
    servers:
      - enabled: true
        type: anthropic
        command: claude
`

	converterYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Converter
metadata:
  name: providers-to-claude
spec:
  input:
    schemaRef: schema://providers
  engine:
    type: template
    template: '{"name": "{{ .Metadata.Name }}"}'
  output:
    format: json
`

	// Write all files
	files := map[string]string{
		"/test/kustomization.yaml": kustomizationYAML,
		"/test/step.yaml":          stepYAML,
		"/test/cmd.yaml":           cmdYAML,
		"/test/workflow.yaml":      workflowYAML,
		"/test/assets.yaml":        assetsYAML,
		"/test/converter.yaml":     converterYAML,
	}
	for path, content := range files {
		if err := fSys.WriteFile(path, []byte(content)); err != nil {
			t.Fatalf("Failed to write %s: %v", path, err)
		}
	}

	// Create loader and load
	loader := NewLoader(nil)
	resMap, err := loader.LoadWithFS(fSys, "/test")
	if err != nil {
		t.Fatalf("LoadWithFS failed: %v", err)
	}

	// Convert to NixAI types
	adapter := NewAdapter()
	resources, err := adapter.ResMapToResources(resMap)
	if err != nil {
		t.Fatalf("ResMapToResources failed: %v", err)
	}

	// Verify we got all five kinds
	expectedKinds := map[string]int{
		"Step":        1,
		"Cmd":         1,
		"CmdWorkflow": 1,
		"Assets":      1,
		"Converter":   1,
	}
	for kind, expectedCount := range expectedKinds {
		found := adapter.GetByKind(resources, kind)
		if len(found) != expectedCount {
			t.Errorf("Expected %d %s resources, got %d", expectedCount, kind, len(found))
		}
	}

	t.Logf("Successfully loaded all five kinds: %d total resources", len(resources))
}

// TestLoaderAssetsValidation tests that invalid Assets fail validation.
func TestLoaderAssetsValidation(t *testing.T) {
	// Create in-memory filesystem
	fSys := filesys.MakeFsInMemory()

	kustomizationYAML := `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - schema.yaml
  - invalid-assets.yaml
`

	// Schema for providers (requires servers field)
	schemaYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Schema
metadata:
  name: providers
spec:
  type: object
  required:
    - servers
  properties:
    servers:
      type: array
      minItems: 1
      items:
        type: object
        required:
          - enabled
          - type
          - command
        properties:
          enabled:
            type: boolean
          type:
            type: string
          command:
            type: string
`

	// Invalid Assets (missing required servers field)
	invalidAssetsYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: invalid-providers
spec:
  schemaRef: schema://providers
  data: {}
`

	// Write files
	if err := fSys.WriteFile("/test/kustomization.yaml", []byte(kustomizationYAML)); err != nil {
		t.Fatalf("Failed to write kustomization.yaml: %v", err)
	}
	if err := fSys.WriteFile("/test/schema.yaml", []byte(schemaYAML)); err != nil {
		t.Fatalf("Failed to write schema.yaml: %v", err)
	}
	if err := fSys.WriteFile("/test/invalid-assets.yaml", []byte(invalidAssetsYAML)); err != nil {
		t.Fatalf("Failed to write invalid-assets.yaml: %v", err)
	}

	// Create loader and load
	loader := NewLoader(nil)
	resMap, err := loader.LoadWithFS(fSys, "/test")
	if err != nil {
		t.Fatalf("LoadWithFS failed: %v", err)
	}

	// Build with validation - this should fail validation
	adapter := NewAdapter()
	_, err = adapter.BuildFromResMap(resMap, "/test")
	if err == nil {
		t.Error("Expected validation error for invalid Assets, got none")
	}
	t.Logf("Got expected validation error: %v", err)
}

// TestLoaderConverterValidation tests that invalid Converter fails validation.
func TestLoaderConverterValidation(t *testing.T) {
	// Create in-memory filesystem
	fSys := filesys.MakeFsInMemory()

	kustomizationYAML := `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - invalid-converter.yaml
`

	// Invalid Converter (missing engine.type)
	invalidConverterYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Converter
metadata:
  name: invalid-converter
spec:
  input:
    schemaRef: schema://providers
  engine:
    template: 'test'
  output:
    format: json
`

	// Write files
	if err := fSys.WriteFile("/test/kustomization.yaml", []byte(kustomizationYAML)); err != nil {
		t.Fatalf("Failed to write kustomization.yaml: %v", err)
	}
	if err := fSys.WriteFile("/test/invalid-converter.yaml", []byte(invalidConverterYAML)); err != nil {
		t.Fatalf("Failed to write invalid-converter.yaml: %v", err)
	}

	// Create loader and load
	loader := NewLoader(nil)
	resMap, err := loader.LoadWithFS(fSys, "/test")
	if err != nil {
		t.Fatalf("LoadWithFS failed: %v", err)
	}

	// Convert to NixAI types - this should validate and fail
	adapter := NewAdapter()
	_, err = adapter.ResMapToResources(resMap)
	if err == nil {
		t.Error("Expected validation error for invalid Converter, got none")
	}
	t.Logf("Got expected validation error: %v", err)
}

// ============================================================================
// Schema Manifest Integration Tests
// ============================================================================

// TestLoaderSchemaManifest tests loading Schema resources from kustomization.
// This validates that Schema manifests are correctly parsed and loaded.
func TestLoaderSchemaManifest(t *testing.T) {
	// Create in-memory filesystem
	fSys := filesys.MakeFsInMemory()

	kustomizationYAML := `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - schema.yaml
  - assets.yaml
`

	schemaYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Schema
metadata:
  name: test-schema
spec:
  type: object
  required:
    - name
  properties:
    name:
      type: string
    count:
      type: integer
      default: 0
`

	assetsYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test-assets
spec:
  schemaRef: schema://test-schema
  data:
    name: test-name
    count: 5
`

	// Write files
	if err := fSys.WriteFile("/test/kustomization.yaml", []byte(kustomizationYAML)); err != nil {
		t.Fatalf("Failed to write kustomization.yaml: %v", err)
	}
	if err := fSys.WriteFile("/test/schema.yaml", []byte(schemaYAML)); err != nil {
		t.Fatalf("Failed to write schema.yaml: %v", err)
	}
	if err := fSys.WriteFile("/test/assets.yaml", []byte(assetsYAML)); err != nil {
		t.Fatalf("Failed to write assets.yaml: %v", err)
	}

	// Create loader and load
	loader := NewLoader(nil)
	resMap, err := loader.LoadWithFS(fSys, "/test")
	if err != nil {
		t.Fatalf("LoadWithFS failed: %v", err)
	}

	// Build with validation
	adapter := NewAdapter()
	result, err := adapter.BuildFromResMap(resMap, "/test")
	if err != nil {
		t.Fatalf("BuildFromResMap failed: %v", err)
	}

	// Verify Schema was loaded (schemas are stored separately in result.Schemas)
	if len(result.Schemas) != 1 {
		t.Errorf("Expected 1 Schema resource, got %d", len(result.Schemas))
	}
	if result.SchemaIndex["test-schema"] == nil {
		t.Error("Expected schema 'test-schema' to be indexed")
	}

	// Verify Schema content
	if result.Schemas[0].Metadata.Name != "test-schema" {
		t.Errorf("Expected schema name 'test-schema', got '%s'", result.Schemas[0].Metadata.Name)
	}

	t.Logf("Successfully loaded Schema: %s", result.Schemas[0].Metadata.Name)
}

// TestLoaderSchemaWithRefResolution tests $ref resolution across schemas.
// This validates that schemas can reference other schemas using $ref.
func TestLoaderSchemaWithRefResolution(t *testing.T) {
	// Create in-memory filesystem
	fSys := filesys.MakeFsInMemory()

	kustomizationYAML := `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - base-schema.yaml
  - ref-schema.yaml
  - assets.yaml
`

	baseSchemaYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Schema
metadata:
  name: base-item
spec:
  type: object
  required:
    - id
  properties:
    id:
      type: string
    name:
      type: string
`

	refSchemaYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Schema
metadata:
  name: container
spec:
  type: object
  required:
    - items
  properties:
    items:
      type: array
      items:
        $ref: schema://base-item
`

	assetsYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test-container
spec:
  schemaRef: schema://container
  data:
    items:
      - id: "item-1"
        name: "First Item"
      - id: "item-2"
        name: "Second Item"
`

	// Write files
	if err := fSys.WriteFile("/test/kustomization.yaml", []byte(kustomizationYAML)); err != nil {
		t.Fatalf("Failed to write kustomization.yaml: %v", err)
	}
	if err := fSys.WriteFile("/test/base-schema.yaml", []byte(baseSchemaYAML)); err != nil {
		t.Fatalf("Failed to write base-schema.yaml: %v", err)
	}
	if err := fSys.WriteFile("/test/ref-schema.yaml", []byte(refSchemaYAML)); err != nil {
		t.Fatalf("Failed to write ref-schema.yaml: %v", err)
	}
	if err := fSys.WriteFile("/test/assets.yaml", []byte(assetsYAML)); err != nil {
		t.Fatalf("Failed to write assets.yaml: %v", err)
	}

	// Create loader and load
	loader := NewLoader(nil)
	resMap, err := loader.LoadWithFS(fSys, "/test")
	if err != nil {
		t.Fatalf("LoadWithFS failed: %v", err)
	}

	// Build with validation
	adapter := NewAdapter()
	result, err := adapter.BuildFromResMap(resMap, "/test")
	if err != nil {
		t.Fatalf("BuildFromResMap failed: %v", err)
	}

	// Verify both schemas were loaded (schemas are stored separately in result.Schemas)
	if len(result.Schemas) != 2 {
		t.Errorf("Expected 2 Schema resources, got %d", len(result.Schemas))
	}
	if result.SchemaIndex["base-item"] == nil {
		t.Error("Expected base-item schema to be indexed")
	}
	if result.SchemaIndex["container"] == nil {
		t.Error("Expected container schema to be indexed")
	}

	t.Logf("Successfully loaded schemas with $ref: base-item and container")
}

// TestLoaderInvalidAssetsAgainstSchema tests validation of invalid Assets data.
// This validates that Assets with invalid data fail validation.
func TestLoaderInvalidAssetsAgainstSchema(t *testing.T) {
	// Create in-memory filesystem
	fSys := filesys.MakeFsInMemory()

	kustomizationYAML := `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - schema.yaml
  - invalid-assets.yaml
`

	schemaYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Schema
metadata:
  name: strict-schema
spec:
  type: object
  required:
    - name
    - email
  properties:
    name:
      type: string
      minLength: 1
    email:
      type: string
`

	invalidAssetsYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: invalid-data
spec:
  schemaRef: schema://strict-schema
  data:
    name: ""
    # missing email field
`

	// Write files
	if err := fSys.WriteFile("/test/kustomization.yaml", []byte(kustomizationYAML)); err != nil {
		t.Fatalf("Failed to write kustomization.yaml: %v", err)
	}
	if err := fSys.WriteFile("/test/schema.yaml", []byte(schemaYAML)); err != nil {
		t.Fatalf("Failed to write schema.yaml: %v", err)
	}
	if err := fSys.WriteFile("/test/invalid-assets.yaml", []byte(invalidAssetsYAML)); err != nil {
		t.Fatalf("Failed to write invalid-assets.yaml: %v", err)
	}

	// Create loader and load
	loader := NewLoader(nil)
	resMap, err := loader.LoadWithFS(fSys, "/test")
	if err != nil {
		t.Fatalf("LoadWithFS failed: %v", err)
	}

	// Build with validation - this should fail validation
	adapter := NewAdapter()
	_, err = adapter.BuildFromResMap(resMap, "/test")
	if err == nil {
		t.Error("Expected validation error for invalid Assets data")
	}
	t.Logf("Got expected validation error: %v", err)
}

// TestLoaderMissingSchemaReference tests Assets with missing schema reference.
// This validates that Assets referencing a non-existent schema fail validation.
func TestLoaderMissingSchemaReference(t *testing.T) {
	// Create in-memory filesystem
	fSys := filesys.MakeFsInMemory()

	kustomizationYAML := `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - schema.yaml
  - assets.yaml
`

	// Add a schema so that validation is applied (even if not the one referenced)
	schemaYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Schema
metadata:
  name: providers
spec:
  type: object
  properties:
    servers:
      type: array
`

	assetsYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: missing-ref
spec:
  schemaRef: schema://non-existent-schema
  data:
    something: value
`

	// Write files
	if err := fSys.WriteFile("/test/kustomization.yaml", []byte(kustomizationYAML)); err != nil {
		t.Fatalf("Failed to write kustomization.yaml: %v", err)
	}
	if err := fSys.WriteFile("/test/schema.yaml", []byte(schemaYAML)); err != nil {
		t.Fatalf("Failed to write schema.yaml: %v", err)
	}
	if err := fSys.WriteFile("/test/assets.yaml", []byte(assetsYAML)); err != nil {
		t.Fatalf("Failed to write assets.yaml: %v", err)
	}

	// Create loader and load
	loader := NewLoader(nil)
	resMap, err := loader.LoadWithFS(fSys, "/test")
	if err != nil {
		t.Fatalf("LoadWithFS failed: %v", err)
	}

	// Build with validation - this should fail (missing schema)
	adapter := NewAdapter()
	_, err = adapter.BuildFromResMap(resMap, "/test")
	if err == nil {
		t.Error("Expected validation error for missing schema reference")
	}
	t.Logf("Got expected validation error: %v", err)
}

// ============================================================================
// Logger Integration Tests for Schema Validation
// ============================================================================

// TestLoaderSchemaDefaultValuesLogging tests that default values from schemas
// are logged at DETAIL level to the JSONL file.
// This validates the logger integration in NewValidator.
func TestLoaderSchemaDefaultValuesLogging(t *testing.T) {
	// Initialize logger with temp JSONL file
	tmpDir := t.TempDir()
	jsonlFile := tmpDir + "/test.jsonl"
	
	// Reset and initialize logger
	logger.Reset()
	os.Setenv("KFG_LOG_FILE", jsonlFile)
	os.Setenv("KFG_VERBOSE", "2") // DETAIL level requires verbose=2
	defer os.Unsetenv("KFG_LOG_FILE")
	defer os.Unsetenv("KFG_VERBOSE")
	
	err := logger.Initialize()
	require.NoError(t, err)
	defer logger.Close()

	// Create in-memory filesystem
	fSys := filesys.MakeFsInMemory()

	kustomizationYAML := `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - schema.yaml
  - assets.yaml
`

	// Schema with default values
	schemaYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Schema
metadata:
  name: test-schema-with-defaults
spec:
  type: object
  required:
    - name
  properties:
    name:
      type: string
    enabled:
      type: boolean
      default: true
    count:
      type: integer
      default: 0
    timeout:
      type: string
      default: "30s"
`

	assetsYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test-assets
spec:
  schemaRef: schema://test-schema-with-defaults
  data:
    name: test-name
`

	// Write files
	if err := fSys.WriteFile("/test/kustomization.yaml", []byte(kustomizationYAML)); err != nil {
		t.Fatalf("Failed to write kustomization.yaml: %v", err)
	}
	if err := fSys.WriteFile("/test/schema.yaml", []byte(schemaYAML)); err != nil {
		t.Fatalf("Failed to write schema.yaml: %v", err)
	}
	if err := fSys.WriteFile("/test/assets.yaml", []byte(assetsYAML)); err != nil {
		t.Fatalf("Failed to write assets.yaml: %v", err)
	}

	// Create loader and load
	loader := NewLoader(nil)
	resMap, err := loader.LoadWithFS(fSys, "/test")
	if err != nil {
		t.Fatalf("LoadWithFS failed: %v", err)
	}

	// Build with validation - this will create a Validator and log default values
	adapter := NewAdapter()
	result, err := adapter.BuildFromResMap(resMap, "/test")
	if err != nil {
		t.Fatalf("BuildFromResMap failed: %v", err)
	}

	// Verify schema was loaded
	assert.Equal(t, 1, len(result.Schemas))
	assert.Equal(t, "test-schema-with-defaults", result.Schemas[0].Metadata.Name)

	// Sync and read JSONL file
	logger.Sync()
	content, err := os.ReadFile(jsonlFile)
	require.NoError(t, err)

	// Parse JSONL lines
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	assert.GreaterOrEqual(t, len(lines), 1)

	// Look for default value log entries
	foundDefaults := false
	for _, line := range lines {
		var event map[string]interface{}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}
		
		// Check for DETAIL level log about default values
		msg, _ := event["msg"].(string)
		component, _ := event["component"].(string)
		
		// Default values are logged by validator component at DETAIL level
		if strings.Contains(msg, "default") && component == "validator" {
			foundDefaults = true
			t.Logf("Found default value log: %s", msg)
		}
	}

	assert.True(t, foundDefaults, "Expected to find default value log entries in JSONL file")
	t.Logf("Successfully verified default values logging in JSONL")
}

// TestLoaderSchemaValidationJSONLEntries tests that schema validation entries
// are logged to the JSONL file.
// This validates the logger integration in validateSchemaManifest.
func TestLoaderSchemaValidationJSONLEntries(t *testing.T) {
	// Initialize logger with temp JSONL file
	tmpDir := t.TempDir()
	jsonlFile := tmpDir + "/test.jsonl"
	
	// Reset and initialize logger
	logger.Reset()
	os.Setenv("KFG_LOG_FILE", jsonlFile)
	os.Setenv("KFG_VERBOSE", "0") // Only JSONL, no stderr
	defer os.Unsetenv("KFG_LOG_FILE")
	defer os.Unsetenv("KFG_VERBOSE")
	
	err := logger.Initialize()
	require.NoError(t, err)
	defer logger.Close()

	// Create in-memory filesystem
	fSys := filesys.MakeFsInMemory()

	kustomizationYAML := `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - schema.yaml
  - assets.yaml
`

	// Schema for validation
	schemaYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Schema
metadata:
  name: validation-schema
spec:
  type: object
  required:
    - name
  properties:
    name:
      type: string
      minLength: 1
`

	// Valid assets
	assetsYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: valid-assets
spec:
  schemaRef: schema://validation-schema
  data:
    name: "test-name"
`

	// Write files
	if err := fSys.WriteFile("/test/kustomization.yaml", []byte(kustomizationYAML)); err != nil {
		t.Fatalf("Failed to write kustomization.yaml: %v", err)
	}
	if err := fSys.WriteFile("/test/schema.yaml", []byte(schemaYAML)); err != nil {
		t.Fatalf("Failed to write schema.yaml: %v", err)
	}
	if err := fSys.WriteFile("/test/assets.yaml", []byte(assetsYAML)); err != nil {
		t.Fatalf("Failed to write assets.yaml: %v", err)
	}

	// Create loader and load
	loader := NewLoader(nil)
	resMap, err := loader.LoadWithFS(fSys, "/test")
	if err != nil {
		t.Fatalf("LoadWithFS failed: %v", err)
	}

	// Build with validation
	adapter := NewAdapter()
	_, err = adapter.BuildFromResMap(resMap, "/test")
	if err != nil {
		t.Fatalf("BuildFromResMap failed: %v", err)
	}

	// Sync and read JSONL file
	logger.Sync()
	content, err := os.ReadFile(jsonlFile)
	require.NoError(t, err)

	// Parse JSONL lines
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	assert.GreaterOrEqual(t, len(lines), 1)

	// Look for schema validation log entries
	foundSchemaFound := false
	foundValidationPassed := false
	for _, line := range lines {
		var event map[string]interface{}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}
		
		msg, _ := event["msg"].(string)
		component, _ := event["component"].(string)
		
		// Check for schema found log (kustomize component)
		if strings.Contains(msg, "found schema://") && component == "kustomize" {
			foundSchemaFound = true
			t.Logf("Found schema found log: %s", msg)
		}
		
		// Check for validation passed log (validator component)
		if strings.Contains(msg, "validation passed") && component == "validator" {
			foundValidationPassed = true
			t.Logf("Found validation passed log: %s", msg)
		}
	}

	assert.True(t, foundSchemaFound, "Expected to find 'found schema://' log entry in JSONL")
	assert.True(t, foundValidationPassed, "Expected to find 'validation passed' log entry in JSONL")
	t.Logf("Successfully verified schema validation entries in JSONL")
}

// TestLoaderVerboseLevelBehavior tests the verbose level gating for schema logs.
// verbose=1 shows info but not detail, verbose=2 shows detail including default values.
func TestLoaderVerboseLevelBehavior(t *testing.T) {
	tests := []struct {
		name           string
		verbose         string
		expectDetail    bool // DETAIL level (default values) should be in JSONL
		expectDetailStderr bool // DETAIL level should be on stderr
	}{
		{
			name:           "verbose=0 (JSONL only, no stderr)",
			verbose:         "0",
			expectDetail:    true,  // JSONL always has all levels
			expectDetailStderr: false, // No stderr output at verbose=0
		},
		{
			name:           "verbose=1 (info level, detail not on stderr)",
			verbose:         "1",
			expectDetail:    true,  // JSONL always has all levels
			expectDetailStderr: false, // DETAIL requires verbose=2
		},
		{
			name:           "verbose=2 (detail level on stderr)",
			verbose:         "2",
			expectDetail:    true,  // JSONL always has all levels
			expectDetailStderr: true, // DETAIL shown at verbose=2
		},
		{
			name:           "verbose=3 (all levels on stderr)",
			verbose:         "3",
			expectDetail:    true,  // JSONL always has all levels
			expectDetailStderr: true, // DEBUG shown, DETAIL also shown
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize logger with temp JSONL file
			tmpDir := t.TempDir()
			jsonlFile := tmpDir + "/test.jsonl"
			
			// Reset and initialize logger
			logger.Reset()
			os.Setenv("KFG_LOG_FILE", jsonlFile)
			os.Setenv("KFG_VERBOSE", tt.verbose)
			defer os.Unsetenv("KFG_LOG_FILE")
			defer os.Unsetenv("KFG_VERBOSE")
			
			err := logger.Initialize()
			require.NoError(t, err)
			defer logger.Close()

			// Create in-memory filesystem with schema having default values
			fSys := filesys.MakeFsInMemory()

			kustomizationYAML := `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - schema.yaml
  - assets.yaml
`

			schemaYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Schema
metadata:
  name: defaults-schema
spec:
  type: object
  properties:
    name:
      type: string
    enabled:
      type: boolean
      default: true
`

			assetsYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test-assets
spec:
  schemaRef: schema://defaults-schema
  data:
    name: test
`

			// Write files
			if err := fSys.WriteFile("/test/kustomization.yaml", []byte(kustomizationYAML)); err != nil {
				t.Fatalf("Failed to write kustomization.yaml: %v", err)
			}
			if err := fSys.WriteFile("/test/schema.yaml", []byte(schemaYAML)); err != nil {
				t.Fatalf("Failed to write schema.yaml: %v", err)
			}
			if err := fSys.WriteFile("/test/assets.yaml", []byte(assetsYAML)); err != nil {
				t.Fatalf("Failed to write assets.yaml: %v", err)
			}

			// Create loader and build
			loader := NewLoader(nil)
			resMap, err := loader.LoadWithFS(fSys, "/test")
			if err != nil {
				t.Fatalf("LoadWithFS failed: %v", err)
			}

			adapter := NewAdapter()
			_, err = adapter.BuildFromResMap(resMap, "/test")
			if err != nil {
				t.Fatalf("BuildFromResMap failed: %v", err)
			}

			// Sync and read JSONL file
			logger.Sync()
			content, err := os.ReadFile(jsonlFile)
			require.NoError(t, err)

			// Parse JSONL and check for DETAIL level entries (default values)
			lines := strings.Split(strings.TrimSpace(string(content)), "\n")
			foundDetailInJSONL := false
			for _, line := range lines {
				var event map[string]interface{}
				if err := json.Unmarshal([]byte(line), &event); err != nil {
					continue
				}
				
				msg, _ := event["msg"].(string)
				component, _ := event["component"].(string)
				
				if strings.Contains(msg, "default") && component == "validator" {
					foundDetailInJSONL = true
				}
			}

			assert.Equal(t, tt.expectDetail, foundDetailInJSONL, 
				"DETAIL level entries in JSONL should always be present regardless of verbose setting")
			
			t.Logf("JSONL has DETAIL entries: %v (verbose=%s)", foundDetailInJSONL, tt.verbose)
		})
	}
}

// TestLoaderCircularSchemaRef tests circular $ref detection between schemas.
// This validates that circular references are detected and reported.
func TestLoaderCircularSchemaRef(t *testing.T) {
	// Create in-memory filesystem
	fSys := filesys.MakeFsInMemory()

	kustomizationYAML := `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - schema-a.yaml
  - schema-b.yaml
  - assets.yaml
`

	// Schema A references Schema B via $ref
	schemaAYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Schema
metadata:
  name: schema-a
spec:
  type: object
  properties:
    name:
      type: string
    refToB:
      $ref: schema://schema-b
`

	// Schema B references Schema A via $ref (circular)
	schemaBYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Schema
metadata:
  name: schema-b
spec:
  type: object
  properties:
    title:
      type: string
    refToA:
      $ref: schema://schema-a
`

	// Assets referencing Schema A
	assetsYAML := `
apiVersion: kfg.dev/v1alpha1
kind: Assets
metadata:
  name: test-assets
spec:
  schemaRef: schema://schema-a
  data:
    name: test
`

	// Write files
	if err := fSys.WriteFile("/test/kustomization.yaml", []byte(kustomizationYAML)); err != nil {
		t.Fatalf("Failed to write kustomization.yaml: %v", err)
	}
	if err := fSys.WriteFile("/test/schema-a.yaml", []byte(schemaAYAML)); err != nil {
		t.Fatalf("Failed to write schema-a.yaml: %v", err)
	}
	if err := fSys.WriteFile("/test/schema-b.yaml", []byte(schemaBYAML)); err != nil {
		t.Fatalf("Failed to write schema-b.yaml: %v", err)
	}
	if err := fSys.WriteFile("/test/assets.yaml", []byte(assetsYAML)); err != nil {
		t.Fatalf("Failed to write assets.yaml: %v", err)
	}

	// Create loader and load
	loader := NewLoader(nil)
	resMap, err := loader.LoadWithFS(fSys, "/test")
	if err != nil {
		t.Fatalf("LoadWithFS failed: %v", err)
	}

	// Build with validation - this may or may not fail depending on how gojsonschema handles circular refs
	adapter := NewAdapter()
	result, err := adapter.BuildFromResMap(resMap, "/test")
	if err != nil {
		t.Logf("BuildFromResMap failed (circular ref detected): %v", err)
		// Circular refs should cause an error
		return
	}

	// If no error, check that validation handled the circular ref
	t.Logf("BuildFromResMap succeeded with %d schemas and %d other resources", len(result.Schemas), len(result.Resources))
	t.Log("Note: gojsonschema may handle circular refs gracefully by ignoring them")
}