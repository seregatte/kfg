## MODIFIED Requirements

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

### Requirement: Context management helpers

The engine runtime MUST provide context management primitives for steps that modify environment or state.

#### Scenario: Output helpers
- **WHEN** a step needs to produce output for the generated shell to consume
- **THEN** output helpers SHALL be available to facilitate structured output
- **AND** helpers SHALL integrate with the artifact and cache collection system

#### Scenario: Cache helpers preserve relative paths
- **WHEN** internal cache helpers persist or restore artifacts
- **THEN** they SHALL preserve the original relative artifact paths rather than reducing them to basenames

#### Scenario: Cache helper diagnostics
- **WHEN** internal cache helpers process cache identity, hits, misses, store, or restore
- **THEN** they SHALL emit runtime detail or debug logs suitable for diagnosing cache behavior
