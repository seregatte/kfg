## MODIFIED Requirements

### Requirement: Step execution contract

All framework steps MUST follow the shell runtime API contract.

#### Scenario: Artifact registration
- **WHEN** a framework step generates artifacts
- **THEN** it SHALL register each artifact via `__kfg_add_artifact <path>`
- **AND** artifacts SHALL be discoverable via the `KFG_ARTIFACTS` variable

#### Scenario: Logging
- **WHEN** a framework step executes
- **THEN** it SHALL use the structured logging API (`_kfg.log.*`)
- **AND** it MUST NOT use unstructured echo or printf

#### Scenario: Nested internal kfg execution
- **WHEN** a framework step invokes a nested internal `kfg` subprocess
- **THEN** it SHALL use the runtime's internal execution wrapper
- **AND** the nested subprocess SHALL NOT emit child startup logs into the parent Step output by default

#### Scenario: Conditional execution
- **WHEN** a framework step has prerequisites or conditions
- **THEN** it SHALL use the `when` condition mechanism
- **AND** it MUST NOT implement its own conditional logic
