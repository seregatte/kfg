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

### Requirement: Structured logging API

The engine runtime MUST provide structured logging that framework steps can call.

#### Scenario: Step logging
- **WHEN** a framework step executes
- **THEN** helpers matching the pattern `_kfg.log.*` SHALL be available
- **AND** steps SHALL use these helpers instead of unstructured echo or printf
- **AND** log output SHALL be properly tagged with component and level information

#### Scenario: Log levels
- **WHEN** a step calls a logging helper
- **THEN** the runtime SHALL support at least `info`, `warn`, and `error` levels
- **AND** verbosity settings SHALL control which log levels appear in output

### Requirement: Conditional execution

The engine runtime MUST provide helpers for conditional step execution based on runtime state.

#### Scenario: Conditional step execution
- **WHEN** a step has a `when` condition
- **THEN** `__kfg_when_*` helpers SHALL evaluate conditions before executing the step
- **AND** conditions SHALL be able to reference environment and artifact state

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
