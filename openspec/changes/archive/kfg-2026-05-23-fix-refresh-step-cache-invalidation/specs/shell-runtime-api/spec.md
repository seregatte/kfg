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

#### Scenario: Active step metadata
- **WHEN** a generated Step wrapper begins executing a referenced Step
- **THEN** it SHALL export `KFG_STEP_NAME` containing the active Step reference name
- **AND** logs emitted during that Step execution SHALL be attributable through the `step_name` field
- **AND** the prior `KFG_STEP_NAME` value, if any, SHALL be restored when the Step exits

#### Scenario: Refresh control
- **WHEN** shell runtime is generated for apply or run
- **THEN** `KFG_REFRESH` SHALL cause cacheable Steps to invalidate and rebuild their own cache entries rather than reuse restored cache results

### Requirement: Artifact collection helpers

The engine runtime MUST provide helpers that allow steps to register build artifacts.

#### Scenario: Registering artifacts
- **WHEN** a step executes and produces output artifacts
- **THEN** it SHALL call `__kfg_add_artifact <artifact_path>` to register each artifact
- **AND** registered artifacts SHALL be available to downstream steps via `KFG_ARTIFACTS`

#### Scenario: Artifact availability
- **WHEN** a downstream step needs to reference upstream artifacts
- **THEN** `KFG_ARTIFACTS` SHALL contain the registered paths for the current invocation
- **AND** cache persistence helpers SHALL distinguish artifacts newly produced by the current Step from artifacts registered earlier in the invocation

#### Scenario: Registering diff-based artifacts under a root
- **WHEN** a Step needs to register artifacts discovered through `__kfg_fs_diff`
- **THEN** the runtime SHALL provide a helper that accepts the snapshot root path plus before/after snapshot files
- **AND** that helper SHALL prefix each diff path with the provided root before calling `__kfg_add_artifact`
- **AND** it SHALL ignore diff paths that no longer exist in the workspace
