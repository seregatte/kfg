## MODIFIED Requirements

### Requirement: Log Command

The CLI MUST provide `kfg sys log` for structured logging used by generated shell helpers.

#### Scenario: Shell helper logging
- **GIVEN** a generated shell helper emits a structured log
- **WHEN** it delegates to the CLI logging command
- **THEN** `kfg sys log` SHALL accept the log event and persist it using the configured logger behavior

#### Scenario: Step context enrichment
- **GIVEN** `KFG_STEP_NAME` is set in the environment
- **WHEN** `kfg sys log` records a shell log event
- **THEN** the event SHALL include `step_name`
- **AND** the field value SHALL match the active Step name from the environment

#### Scenario: Legacy step component normalization
- **WHEN** `kfg sys log` receives a shell component matching `step:<name>`
- **THEN** it SHALL record `component` as `step`
- **AND** it SHALL record `step_name` as `<name>`
- **AND** the command SHALL succeed without requiring the caller to change arguments
