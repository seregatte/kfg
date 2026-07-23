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
- **AND** logs emitted during that Step execution SHALL use that value for structured Step attribution
- **AND** any previous `KFG_STEP_NAME` value SHALL be restored when the Step exits

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
- **THEN** the runtime SHALL accept the call as a structured log event
- **AND** it SHALL default the component to `step`
- **AND** it SHALL attach `step_name` when Step context is available

#### Scenario: Legacy step component compatibility
- **WHEN** a Step calls a logging helper with a component matching `step:<name>`
- **THEN** the runtime SHALL normalize the component to `step`
- **AND** it SHALL emit `step_name` with the legacy Step name value
- **AND** the log call SHALL remain successful without requiring manifest changes
