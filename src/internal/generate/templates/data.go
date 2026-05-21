package templates

// HeaderData represents data for the bash_header template.
type HeaderData struct {
	SetName string // Old field for backward compat
	Shell   string
	// New fields for kustomization model
	WorkflowName      string
	KustomizationName string
}

// StepData represents data for the bash_step template.
type StepData struct {
	StepName      string
	WhenCondition string
	HasOutput     bool
	RunScript     string
	OutputName    string
	IsMultiLine   bool
	RunLines      []string
	IgnoreFailure bool
	Artifacts     []string          // NEW: artifacts produced by this step
	Env           map[string]string // NEW: environment variables for this step
	// Cache configuration
	CacheEnabled bool   // Whether caching is enabled for this step
	CacheKey     string // User-provided cache key
	ScriptHash   string // Hash of spec.run for cache identity
}

// BeforeStepData represents a before step in a command wrapper.
// Weight field removed - use explicit YAML order instead.
type BeforeStepData struct {
	StepName      string
	IgnoreFailure bool
	WhenCondition string
	Weight        int // DEPRECATED: kept for backward compat, unused in new model
}

// AfterStepData represents an after step in a command wrapper.
// Weight field removed - use explicit YAML order instead.
type AfterStepData struct {
	StepName      string
	IgnoreFailure bool
	WhenCondition string
	Weight        int // DEPRECATED: kept for backward compat, unused in new model
}

// CommandData represents data for the bash_command template.
// Used for old Command type (deprecated).
type CommandData struct {
	CommandName    string
	HasBeforeSteps bool
	BeforeSteps    []BeforeStepData
	MainRun        string
	HasAfterSteps  bool
	AfterSteps     []AfterStepData
}

// CmdData represents data for a Cmd wrapper (new type).
// Cmds are pure functions WITHOUT orchestration.
type CmdData struct {
	CmdName     string // Bash function name (metadata.commandName)
	RunScript   string
	Artifacts   []string // Artifacts produced by this cmd
	IsMultiLine bool
	RunLines    []string
	Env         map[string]string // NEW: environment variables for this cmd
}

// WorkflowStepData represents a step in workflow context.
// Used for global workflow before/after steps.
// StepRefName is the StepReference.name (runtime execution identity) used for output addressing.
type WorkflowStepData struct {
	StepRefName   string // StepReference.name (runtime execution identity)
	StepName      string // Step metadata.name (for function lookup)
	IgnoreFailure bool
	WhenCondition string
	Env           map[string]string // NEW: merged env for this step invocation
	// Cache configuration
	CacheEnabled bool   // Whether caching is enabled for this step invocation
	CacheKey     string // User-provided cache key for this invocation
	ScriptHash   string // Hash of spec.run for cache identity
	HasOutput    bool   // Whether this step has an output
	OutputName   string // Output name if HasOutput is true
	// Declarative artifacts
	StepArtifacts []string // Artifacts declared in Step.Spec.Artifacts
	RefArtifacts  []string // Artifacts declared in StepReference.Artifacts
}

// WorkflowCmdData represents a cmd in workflow context.
// Includes global workflow steps and per-cmd steps.
type WorkflowCmdData struct {
	CmdName           string
	RunScript         string
	Artifacts         []string
	HasGlobalBefore   bool
	GlobalBeforeSteps []WorkflowStepData
	HasGlobalAfter    bool
	GlobalAfterSteps  []WorkflowStepData
	HasCmdBefore      bool
	CmdBeforeSteps    []WorkflowStepData
	HasCmdAfter       bool
	CmdAfterSteps     []WorkflowStepData
	Env               map[string]string // NEW: environment variables for this cmd
}

// WorkflowData represents data for a complete workflow.
type WorkflowData struct {
	WorkflowName string
	Shell        string
	Cmds         []WorkflowCmdData
	Steps        []StepData
	HasBefore    bool
	BeforeSteps  []WorkflowStepData
	HasAfter     bool
	AfterSteps   []WorkflowStepData
}

// KustomizationData represents data for the full kustomization output.
type KustomizationData struct {
	KustomizationName string
	WorkflowName      string
	Shell             string
	Workflow          WorkflowData
}

// TemplateData represents all data needed to generate shell code (old).
// DEPRECATED: Use KustomizationData instead.
type TemplateData struct {
	Header   HeaderData
	Steps    []StepData
	Commands []CommandData
}
