# Shell Runtime API Specification

## Purpose

kfg generates shell code that provides runtime facilities for framework steps and domain packages to consume. This specification defines the stable contract between the engine's shell runtime and the code it generates.

## Requirements

### Requirement: Runtime environment variables

The engine runtime MUST export metadata environment variables that steps and packages can consume.

#### Scenario: Session identification
- **WHEN** the engine generates shell code
- **THEN** it SHALL export `KFG_SESSION_ID` containing a unique session identifier
- **AND** the session ID SHALL remain constant across all step executions in a single invocation

#### Scenario: Workflow metadata
- **WHEN** a step executes within a generated shell context
- **THEN** `KFG_WORKFLOW_NAME` SHALL contain the CmdWorkflow name
- **AND** `KFG_KUSTOMIZATION_NAME` SHALL contain the base kustomization name
- **AND** `KFG_SHELL` SHALL contain the shell type (e.g., `bash`)

#### Scenario: Active step metadata
- **WHEN** a generated Step wrapper begins executing a referenced Step
- **THEN** it SHALL export `KFG_STEP_NAME` containing the active Step reference name
- **AND** logs emitted during that Step execution SHALL use that value for structured Step attribution
- **AND** any previous `KFG_STEP_NAME` value SHALL be restored when the Step exits

#### Scenario: Refresh control
- **WHEN** shell runtime is generated for apply or run
- **THEN** `KFG_REFRESH` SHALL control whether cacheable Steps bypass cached entries

### Requirement: Artifact collection helpers

The engine runtime MUST provide helpers that allow steps to register build artifacts.

#### Scenario: Registering artifacts
- **WHEN** a step executes and produces output artifacts
- **THEN** it SHALL call `__kfg_add_artifact <artifact_path>` to register each artifact
- **AND** registered artifacts SHALL be available to downstream steps via `KFG_ARTIFACTS`

#### Scenario: Artifact availability
- **WHEN** a downstream step needs to reference upstream artifacts
- **THEN** `KFG_ARTIFACTS` SHALL contain a colon-separated list of registered paths
- **AND** steps SHALL be able to iterate over `KFG_ARTIFACTS`
- **AND** cache persistence helpers SHALL distinguish artifacts newly produced by the current Step from artifacts registered earlier in the invocation

#### Scenario: Cache helpers preserve relative paths
- **WHEN** internal cache helpers persist or restore artifacts
- **THEN** they SHALL preserve the original relative artifact paths rather than reducing them to basenames

#### Scenario: Cache helper diagnostics
- **WHEN** internal cache helpers process cache identity, hits, misses, store, or restore
- **THEN** they SHALL emit runtime detail or debug logs suitable for diagnosing cache behavior

#### Scenario: Filesystem snapshot wrappers
- **WHEN** a Step needs to discover new paths beneath a directory root
- **THEN** runtime helpers for filesystem snapshot and diff SHALL be available
- **AND** those helpers SHALL delegate to the internal `kfg sys fs` command surface

### Requirement: Build result reporting

The engine runtime MUST provide helpers for steps to report build results and state.

#### Scenario: Recording build result
- **WHEN** a step completes execution
- **THEN** it MAY call `__kfg_build_result <key> <value>` to record result data
- **AND** results SHALL be persisted in `KFG_BUILD_RESULT_FILE`

#### Scenario: Result persistence
- **WHEN** multiple steps execute in sequence
- **THEN** `KFG_BUILD_RESULT_FILE` SHALL be a persistent file path
- **AND** each step's results SHALL accumulate across the workflow

### Requirement: Context management helpers

The engine runtime MUST provide context management primitives for steps that modify environment or state.

#### Scenario: Context reset
- **WHEN** a step needs to clean up its modifications
- **THEN** `__kfg_ctx_reset` SHALL provide a mechanism to restore prior context
- **AND** the mechanism SHALL be step-aware (resetting only the current step's modifications)

#### Scenario: Output helpers
- **WHEN** a step needs to produce output for the generated shell to consume
- **THEN** output helpers SHALL be available to facilitate structured output
- **AND** helpers SHALL integrate with the artifact and result collection system

#### Scenario: Output-producing step preserves runtime side effects
- **WHEN** a Step declares `spec.output`
- **THEN** the generated runtime SHALL execute that Step without a subshell that discards runtime side effects
- **AND** artifact registrations performed during the Step SHALL remain visible to the parent shell runtime after execution

### Requirement: Structured logging API

The engine runtime MUST provide structured logging that framework steps can call.

#### Scenario: Step logging
- **WHEN** a framework step executes
- **THEN** helpers matching the pattern `__kfg_log_*` SHALL be available
- **AND** steps SHALL use these helpers instead of unstructured echo or printf
- **AND** log output SHALL be properly tagged with component, level, and Step identity when Step context is available

#### Scenario: Log levels
- **WHEN** a step calls a logging helper
- **THEN** the runtime SHALL support at least `info`, `warn`, and `error` levels
- **AND** verbosity settings SHALL control which log levels appear in output

#### Scenario: Message-only helper compatibility
- **WHEN** a Step calls a logging helper with only a message argument
- **THEN** the runtime SHALL treat the call as a valid structured log event
- **AND** it SHALL default the component to `step`
- **AND** it SHALL attach `step_name` when Step context is available

#### Scenario: Legacy step component compatibility
- **WHEN** a Step calls a logging helper with a legacy component matching `step:<name>`
- **THEN** the runtime SHALL normalize the component to `step`
- **AND** it SHALL emit `step_name` with the legacy Step name value
- **AND** the log call SHALL remain successful without requiring manifest changes

#### Scenario: Logging backend compatibility
- **WHEN** a generated logging helper is invoked
- **THEN** it MAY delegate to `kfg sys log`
- **AND** the helper naming contract exposed to Steps SHALL remain `__kfg_log_*`

### Requirement: Conditional execution

The engine runtime MUST provide helpers for conditional step execution based on runtime state.

#### Scenario: Conditional step execution
- **WHEN** a step has a `when` condition
- **THEN** `__kfg_when_*` helpers SHALL evaluate conditions before executing the step
- **AND** conditions SHALL be able to reference environment and artifact state

### Requirement: Internal command wrappers

The generated shell runtime MUST provide wrappers for internal engine subprocesses.

#### Scenario: Quiet internal kfg execution
- **WHEN** runtime code invokes a nested internal `kfg` subprocess through the dedicated wrapper
- **THEN** that subprocess SHALL execute with child-scoped `KFG_VERBOSE=0`
- **AND** the parent shell environment SHALL keep its original `KFG_VERBOSE` value

### Requirement: Stability and versioning

The shell runtime API MUST remain stable to allow framework packages to depend on it without version coupling.

#### Scenario: API stability guarantee
- **WHEN** a framework package is deployed into a new engine version
- **THEN** the shell runtime API helpers and variables documented in this spec SHALL remain available
- **AND** helper signatures SHALL NOT change in breaking ways
- **AND** new helpers MAY be added without breaking existing packages

#### Scenario: API documentation for consumers
- **WHEN** a framework or domain package needs to implement a reusable step
- **THEN** it SHALL be able to reference this spec as the stable contract
- **AND** the spec SHALL be authoritative for what runtime facilities are available
