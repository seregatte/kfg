## MODIFIED Requirements

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

#### Scenario: Filesystem snapshot wrappers
- **WHEN** a Step needs to discover new paths beneath a directory root
- **THEN** runtime helpers for filesystem snapshot and diff SHALL be available
- **AND** those helpers SHALL delegate to the internal `kfg sys fs` command surface

### Requirement: Structured logging API

The engine runtime MUST provide structured logging that framework steps can call.

#### Scenario: Step logging
- **WHEN** a framework step executes
- **THEN** helpers matching the pattern `_kfg.log.*` SHALL be available
- **AND** steps SHALL use these helpers instead of unstructured echo or printf
- **AND** nested internal `kfg` subprocesses triggered through runtime wrappers SHALL NOT emit child human startup logs by default

## ADDED Requirements

### Requirement: Internal command wrappers

The generated shell runtime MUST provide wrappers for internal engine subprocesses.

#### Scenario: Quiet internal kfg execution
- **WHEN** runtime code invokes a nested internal `kfg` subprocess through the dedicated wrapper
- **THEN** that subprocess SHALL execute with child-scoped `KFG_VERBOSE=0`
- **AND** the parent shell environment SHALL keep its original `KFG_VERBOSE` value
