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

#### Scenario: Filesystem snapshot helpers
- **WHEN** a Step needs to discover filesystem changes relative to a directory root
- **THEN** the runtime SHALL expose helpers that delegate to `kfg sys fs snapshot` and `kfg sys fs diff`
- **AND** those helpers SHALL be usable together with `__kfg_add_artifact` to register dynamically discovered artifacts

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

#### Scenario: Quiet internal kfg subprocesses
- **WHEN** runtime helpers invoke nested internal `kfg` subprocesses
- **THEN** those subprocesses SHALL execute with child-scoped human log verbosity disabled
- **AND** the parent invocation's verbosity SHALL remain unchanged

## ADDED Requirements

### Requirement: Internal filesystem command backend

The internal CLI SHALL provide filesystem snapshot and diff commands that generated runtime helpers can call on any supported platform.

#### Scenario: Snapshot with depth limit
- **WHEN** user runs `kfg sys fs snapshot <path> --maxdepth 1`
- **THEN** the command SHALL print normalized relative child paths beneath `<path>`
- **AND** it SHALL include only entries reachable within depth `1`
- **AND** it SHALL sort output deterministically

#### Scenario: Snapshot with unlimited depth
- **WHEN** user runs `kfg sys fs snapshot <path> --maxdepth 0`
- **THEN** the command SHALL traverse the full subtree under `<path>`
- **AND** it SHALL print normalized relative paths for all discovered entries

#### Scenario: Diff returns only new paths
- **WHEN** user runs `kfg sys fs diff --before <snapshot> --after <snapshot>`
- **THEN** the command SHALL print only paths present in `after` and absent in `before`
- **AND** it SHALL preserve deterministic ordering
