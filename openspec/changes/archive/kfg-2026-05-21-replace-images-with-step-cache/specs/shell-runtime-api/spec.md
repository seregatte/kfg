## MODIFIED Requirements

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
- **THEN** `KFG_ARTIFACTS` SHALL contain the registered paths for the current invocation
- **AND** steps SHALL be able to iterate over `KFG_ARTIFACTS`

### Requirement: Context management helpers

The engine runtime MUST provide context management primitives for steps that modify environment or state.

#### Scenario: Context reset
- **WHEN** a step needs to clean up its modifications
- **THEN** `__kfg_ctx_reset` SHALL provide a mechanism to restore prior context
- **AND** the mechanism SHALL be step-aware for runtime output state

#### Scenario: Output helpers
- **WHEN** a step needs to produce output for the generated shell to consume
- **THEN** output helpers SHALL be available to facilitate structured output
- **AND** helpers SHALL integrate with the artifact and cache collection system

#### Scenario: Cache helpers
- **WHEN** a cacheable Step executes within generated shell code
- **THEN** the runtime SHALL provide internal helpers to compute cache identity, persist artifacts and outputs, and restore them on cache hit

### Requirement: Structured logging API

The engine runtime MUST provide structured logging that framework steps can call.

#### Scenario: Step logging
- **WHEN** a framework step executes
- **THEN** helpers matching the pattern `__kfg_log_*` SHALL be available
- **AND** steps SHALL use these helpers instead of unstructured echo or printf
- **AND** log output SHALL be properly tagged with component and level information

#### Scenario: Log levels
- **WHEN** a step calls a logging helper
- **THEN** the runtime SHALL support at least `info`, `warn`, and `error` levels
- **AND** verbosity settings SHALL control which log levels appear in output

#### Scenario: Logging backend compatibility
- **WHEN** a generated logging helper is invoked
- **THEN** it MAY delegate to `kfg sys log`
- **AND** the helper naming contract exposed to Steps SHALL remain `__kfg_log_*`
