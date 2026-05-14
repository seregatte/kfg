package manifest

import (
	"fmt"
	"strings"
)

// ============================================================================
// Core Types
// ============================================================================

// ResourceIdentity uniquely identifies a resource within a kind.
type ResourceIdentity struct {
	APIVersion string
	Kind       string
	Name       string
}

// String returns a human-readable representation of the identity.
func (id ResourceIdentity) String() string {
	return fmt.Sprintf("%s/%s:%s", id.APIVersion, id.Kind, id.Name)
}

// Metadata contains the resource's metadata.
type Metadata struct {
	Name        string `yaml:"name"`        // Required for all resources (namespace format)
	CommandName string `yaml:"commandName"` // Required for Cmd resources (bash function name)
	Shell       string `yaml:"shell"`       // Optional for CmdWorkflow (default: "bash")
}

// ============================================================================
// Kind-Specific Types
// ============================================================================

// Step represents a Step resource.
// Steps are reusable shell snippets that can be referenced in workflows.
type Step struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       StepSpec `yaml:"spec"`
}

// Identity returns the step's unique identity.
func (s Step) Identity() ResourceIdentity {
	return ResourceIdentity{
		APIVersion: s.APIVersion,
		Kind:       s.Kind,
		Name:       s.Metadata.Name,
	}
}

// Output defines an output that a Step can produce.
type Output struct {
	Name string `yaml:"name"` // Name of the output variable
	Type string `yaml:"type"` // Type of the output (string, boolean, etc.)
}

// StepSpec is the spec for Step resources.
type StepSpec struct {
	Run       string            `yaml:"run"`       // Required: shell script to execute
	Output    *Output           `yaml:"output"`    // Optional: output capture configuration
	Artifacts []string          `yaml:"artifacts"` // Optional: artifacts produced by this step
	Env       map[string]string `yaml:"env"`       // Optional: environment variables for this step
}

// Cmd represents a Cmd resource (pure shell function, no orchestration).
type Cmd struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       CmdSpec  `yaml:"spec"`
}

// Identity returns the cmd's unique identity.
func (c Cmd) Identity() ResourceIdentity {
	return ResourceIdentity{
		APIVersion: c.APIVersion,
		Kind:       c.Kind,
		Name:       c.Metadata.Name,
	}
}

// CmdSpec is the spec for Cmd resources.
// Cmds are pure shell functions WITHOUT orchestration (before/after moved to CmdWorkflow).
type CmdSpec struct {
	Run       string            `yaml:"run"`       // Required: shell script to execute
	Artifacts []string          `yaml:"artifacts"` // Optional: artifacts produced by this command
	Env       map[string]string `yaml:"env"`       // Optional: environment variables for this command
}

// CmdWorkflow represents a CmdWorkflow resource (orchestration for Cmds).
// Shell type is in CmdWorkflow.Metadata.Shell.
type CmdWorkflow struct {
	APIVersion string          `yaml:"apiVersion"`
	Kind       string          `yaml:"kind"`
	Metadata   Metadata        `yaml:"metadata"`
	Spec       CmdWorkflowSpec `yaml:"spec"`
}

// Identity returns the workflow's unique identity.
func (w CmdWorkflow) Identity() ResourceIdentity {
	return ResourceIdentity{
		APIVersion: w.APIVersion,
		Kind:       w.Kind,
		Name:       w.Metadata.Name,
	}
}

// CmdWorkflowSpec is the spec for CmdWorkflow resources.
// CmdWorkflows handle orchestration for multiple Cmds with before/after steps.
type CmdWorkflowSpec struct {
	Cmds   []string        `yaml:"cmds"`   // Required: list of Cmd names to orchestrate
	Before []StepReference `yaml:"before"` // Optional: steps to run before ALL cmds
	After  []StepReference `yaml:"after"`  // Optional: steps to run after ALL cmds
}

type StepReference struct {
	Step          string            `yaml:"step"`          // Required: Step metadata.name
	When          *WhenClause       `yaml:"when"`          // Optional: contextual condition
	FailurePolicy string            `yaml:"failurePolicy"` // Optional: "Fail" (default) or "Ignore"
	Env           map[string]string `yaml:"env"`           // Optional: environment variable overrides
	Artifacts     []string          `yaml:"artifacts"`     // Optional: additional artifacts produced by this step reference
}

// WhenClause defines conditional execution.
type WhenClause struct {
	Output *OutputCondition `yaml:"output"` // Optional: condition based on output
	AllOf  []WhenClause     `yaml:"allOf"`  // Optional: all conditions must match
	AnyOf  []WhenClause     `yaml:"anyOf"`  // Optional: any condition must match
	Not    *WhenClause      `yaml:"not"`    // Optional: negate condition
}

// OutputCondition defines a condition based on a step output.
type OutputCondition struct {
	Step     string   `yaml:"step"`     // Required: step that produced the output
	Name     string   `yaml:"name"`     // Required: output variable name
	Equals   string   `yaml:"equals"`   // Optional: exact match
	In       []string `yaml:"in"`       // Optional: value in list
	Contains string   `yaml:"contains"` // Optional: substring match
	Matches  string   `yaml:"matches"`  // Optional: regex match
}

// ============================================================================
// Validation
// ============================================================================

// ValidateName checks if the resource name is valid (namespace format).
func ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("metadata.name is required")
	}
	if !isValidNamespaceName(name) {
		return fmt.Errorf("metadata.name '%s' is not a valid namespace identifier (must contain only lowercase alphanumeric, hyphens, and dots, and not start with a digit)", name)
	}
	return nil
}

// isValidNamespaceName checks if a name is a valid namespace identifier.
func isValidNamespaceName(name string) bool {
	if name == "" {
		return false
	}
	if name[0] >= '0' && name[0] <= '9' {
		return false
	}
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '.') {
			return false
		}
	}
	return true
}

// ValidateCommandName checks if metadata.commandName is valid (bash function name).
// Applies to Cmd resources.
func ValidateCommandName(commandName string) error {
	if commandName == "" {
		return fmt.Errorf("metadata.commandName is required for Cmd resources")
	}
	if !isValidBashIdentifier(commandName) {
		return fmt.Errorf("metadata.commandName '%s' is not a valid bash identifier (must start with letter or underscore, contain only alphanumeric, underscores, and hyphens)", commandName)
	}
	return nil
}

// isValidBashIdentifier checks if a name is a valid bash function identifier.
func isValidBashIdentifier(name string) bool {
	if name == "" {
		return false
	}
	c := name[0]
	if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_') {
		return false
	}
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '-') {
			return false
		}
	}
	return true
}

// ValidateShell checks if metadata.shell is valid.
func ValidateShell(shell string) error {
	if shell == "" {
		return nil // Optional, defaults to "bash"
	}
	validShells := []string{"bash", "zsh", "fish", "sh"}
	for _, s := range validShells {
		if shell == s {
			return nil
		}
	}
	return fmt.Errorf("metadata.shell '%s' is not valid (must be one of: bash, zsh, fish, sh)", shell)
}

// ValidateStep validates a Step resource.
func (s Step) Validate() error {
	if err := s.ValidateAPIVersion(); err != nil {
		return err
	}
	if err := s.ValidateKind(); err != nil {
		return err
	}
	if err := ValidateName(s.Metadata.Name); err != nil {
		return err
	}
	if s.Spec.Run == "" {
		return fmt.Errorf("spec.run is required for Step resources")
	}
	return nil
}

// ValidateCmd validates a Cmd resource.
func (c Cmd) Validate() error {
	if err := c.ValidateAPIVersion(); err != nil {
		return err
	}
	if err := c.ValidateKind(); err != nil {
		return err
	}
	if err := ValidateName(c.Metadata.Name); err != nil {
		return err
	}
	if err := ValidateCommandName(c.Metadata.CommandName); err != nil {
		return err
	}
	if c.Spec.Run == "" {
		return fmt.Errorf("spec.run is required for Cmd resources")
	}
	return nil
}

// ValidateCmdWorkflow validates a CmdWorkflow resource.
func (w CmdWorkflow) Validate() error {
	if err := w.ValidateAPIVersion(); err != nil {
		return err
	}
	if err := w.ValidateKind(); err != nil {
		return err
	}
	if err := ValidateName(w.Metadata.Name); err != nil {
		return err
	}
	if err := ValidateShell(w.Metadata.Shell); err != nil {
		return err
	}
	if len(w.Spec.Cmds) == 0 && len(w.Spec.Before) == 0 && len(w.Spec.After) == 0 {
		return fmt.Errorf("spec must have at least one of: cmds, before, after for CmdWorkflow resources")
	}
	return nil
}

// ValidateAPIVersion checks if the API version is correct.
func (s Step) ValidateAPIVersion() error {
	if s.APIVersion != APIVersion {
		return fmt.Errorf("apiVersion must be %s, got %s", APIVersion, s.APIVersion)
	}
	return nil
}

func (c Cmd) ValidateAPIVersion() error {
	if c.APIVersion != APIVersion {
		return fmt.Errorf("apiVersion must be %s, got %s", APIVersion, c.APIVersion)
	}
	return nil
}

func (w CmdWorkflow) ValidateAPIVersion() error {
	if w.APIVersion != APIVersion {
		return fmt.Errorf("apiVersion must be %s, got %s", APIVersion, w.APIVersion)
	}
	return nil
}

// ValidateKind checks if the kind is correct.
func (s Step) ValidateKind() error {
	if s.Kind != "Step" {
		return fmt.Errorf("kind must be Step, got %s", s.Kind)
	}
	return nil
}

func (c Cmd) ValidateKind() error {
	if c.Kind != "Cmd" {
		return fmt.Errorf("kind must be Cmd, got %s", c.Kind)
	}
	return nil
}

func (w CmdWorkflow) ValidateKind() error {
	if w.Kind != "CmdWorkflow" {
		return fmt.Errorf("kind must be CmdWorkflow, got %s", w.Kind)
	}
	return nil
}

// ============================================================================
// Error Types
// ============================================================================

// ParseError represents an error during parsing.
type ParseError struct {
	File    string
	Line    int
	Message string
}

// Error returns the error message.
func (e ParseError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("%s:%d: %s", e.File, e.Line, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.File, e.Message)
}

// ValidationError represents an error during validation.
type ValidationError struct {
	Identity ResourceIdentity
	File     string
	Message  string
	Hint     string
}

// Error returns the error message.
func (e ValidationError) Error() string {
	msg := fmt.Sprintf("%s: %s", e.Identity, e.Message)
	if e.File != "" {
		msg = fmt.Sprintf("%s (File: %s)", msg, e.File)
	}
	if e.Hint != "" {
		msg = fmt.Sprintf("%s\nHint: %s", msg, e.Hint)
	}
	return msg
}

// ============================================================================
// Constants
// ============================================================================

// APIVersion is the expected API version for NixAI manifests.
const APIVersion = "kfg.dev/v1alpha1"

// SupportedKinds are the supported resource kinds.
var SupportedKinds = []string{
	"Step",
	"Cmd",
	"CmdWorkflow",
	"Assets",
	"Converter",
}

// SupportedInputFormats are the supported input formats for Assets and Converters.
var SupportedInputFormats = []string{
	"yaml", "json", "xml", "props", "csv", "tsv", "toml",
	"hcl", "lua", "ini", "shell", "base64", "uri", "kyaml",
}

// SupportedOutputFormats are the supported output formats for Converters.
var SupportedOutputFormats = []string{
	"yaml", "json", "xml", "props", "csv", "tsv", "toml",
	"hcl", "lua", "ini", "shell", "base64", "uri", "kyaml",
	"raw",
}

// ============================================================================
// Assets and Converter Types
// ============================================================================

// Assets represents an Assets resource kind for declaring data payloads.
type Assets struct {
	APIVersion string     `yaml:"apiVersion"`
	Kind       string     `yaml:"kind"`
	Metadata   Metadata   `yaml:"metadata"`
	Spec       AssetsSpec `yaml:"spec"`
}

// Identity returns the assets' unique identity.
func (a Assets) Identity() ResourceIdentity {
	return ResourceIdentity{
		APIVersion: a.APIVersion,
		Kind:       a.Kind,
		Name:       a.Metadata.Name,
	}
}

// AssetsSpec is the spec for Assets resources.
type AssetsSpec struct {
	Input InputSpec `yaml:"input"` // Input format configuration
	Data  any       `yaml:"data"`  // Data payload (map for YAML, string for others)
}

// InputSpec defines input format configuration.
type InputSpec struct {
	Format string `yaml:"format"` // Data format (default: yaml)
}

// Converter represents a Converter resource kind for declaring transformations.
type Converter struct {
	APIVersion string        `yaml:"apiVersion"`
	Kind       string        `yaml:"kind"`
	Metadata   Metadata      `yaml:"metadata"`
	Spec       ConverterSpec `yaml:"spec"`
}

// Identity returns the converter's unique identity.
func (c Converter) Identity() ResourceIdentity {
	return ResourceIdentity{
		APIVersion: c.APIVersion,
		Kind:       c.Kind,
		Name:       c.Metadata.Name,
	}
}

// ConverterSpec is the spec for Converter resources.
type ConverterSpec struct {
	Input  InputSpec  `yaml:"input"`  // Input format configuration
	Engine EngineSpec `yaml:"engine"` // Transformation engine configuration
	Output OutputSpec `yaml:"output"` // Output format configuration
}

// EngineSpec defines the transformation engine configuration.
type EngineSpec struct {
	Expression string `yaml:"expression"` // yq-go expression to evaluate
}

// OutputSpec defines output format configuration.
type OutputSpec struct {
	Format string `yaml:"format"` // Output format (default: yaml)
}

// ============================================================================
// Validation
// ============================================================================

// ValidateAPIVersion checks if the apiVersion is correct for Assets.
func (a Assets) ValidateAPIVersion() error {
	if a.APIVersion != APIVersion {
		return fmt.Errorf("apiVersion must be %s, got %s", APIVersion, a.APIVersion)
	}
	return nil
}

// ValidateKind checks if the kind is correct for Assets.
func (a Assets) ValidateKind() error {
	if a.Kind != "Assets" {
		return fmt.Errorf("kind must be Assets, got %s", a.Kind)
	}
	return nil
}

// ValidateAPIVersion checks if the apiVersion is correct for Converter.
func (c Converter) ValidateAPIVersion() error {
	if c.APIVersion != APIVersion {
		return fmt.Errorf("apiVersion must be %s, got %s", APIVersion, c.APIVersion)
	}
	return nil
}

// ValidateKind checks if the kind is correct for Converter.
func (c Converter) ValidateKind() error {
	if c.Kind != "Converter" {
		return fmt.Errorf("kind must be Converter, got %s", c.Kind)
	}
	return nil
}

// ValidateAssets validates an Assets resource.
func (a Assets) Validate() error {
	if err := a.ValidateAPIVersion(); err != nil {
		return err
	}
	if err := a.ValidateKind(); err != nil {
		return err
	}
	if err := ValidateName(a.Metadata.Name); err != nil {
		return err
	}
	// Validate input format
	format := a.Spec.Input.Format
	if format == "" {
		format = "yaml"
	}
	if !isSupportedInputFormat(format) {
		return fmt.Errorf("spec.input.format '%s' is not supported (supported: %s)", format, strings.Join(SupportedInputFormats, ", "))
	}
	// Validate data is present
	if a.Spec.Data == nil {
		return fmt.Errorf("spec.data is required for Assets resources")
	}
	return nil
}

// ValidateConverter validates a Converter resource.
func (c Converter) Validate() error {
	if err := c.ValidateAPIVersion(); err != nil {
		return err
	}
	if err := c.ValidateKind(); err != nil {
		return err
	}
	if err := ValidateName(c.Metadata.Name); err != nil {
		return err
	}
	// Validate input format
	inputFormat := c.Spec.Input.Format
	if inputFormat == "" {
		inputFormat = "yaml"
	}
	if !isSupportedInputFormat(inputFormat) {
		return fmt.Errorf("spec.input.format '%s' is not supported (supported: %s)", inputFormat, strings.Join(SupportedInputFormats, ", "))
	}
	// Validate output format
	outputFormat := c.Spec.Output.Format
	if outputFormat == "" {
		outputFormat = "yaml"
	}
	if !isSupportedOutputFormat(outputFormat) {
		return fmt.Errorf("spec.output.format '%s' is not supported (supported: %s)", outputFormat, strings.Join(SupportedOutputFormats, ", "))
	}
	// Validate expression is present
	if c.Spec.Engine.Expression == "" {
		return fmt.Errorf("spec.engine.expression is required for Converter resources")
	}
	return nil
}

// isSupportedInputFormat checks if a format is a supported input format.
func isSupportedInputFormat(format string) bool {
	for _, f := range SupportedInputFormats {
		if f == format {
			return true
		}
	}
	return false
}

// isSupportedOutputFormat checks if a format is a supported output format.
func isSupportedOutputFormat(format string) bool {
	for _, f := range SupportedOutputFormats {
		if f == format {
			return true
		}
	}
	return false
}
