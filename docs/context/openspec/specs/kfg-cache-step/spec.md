## Purpose

Define how cacheable Steps persist and restore their runtime results, including artifact isolation semantics, path preservation, and diagnostic logging.
## Requirements
### Requirement: Cache Identity and Invalidation
The runtime SHALL identify cache entries per resolved Step invocation using only the `StepReference.name`.

#### Scenario: Invocation-specific cache key
- **WHEN** a cacheable workflow Step executes
- **THEN** the cache identity SHALL be indentified by `StepReference.name`
- **AND** SHALL NOT include any user-supplied key string
- **AND** SHALL NOT include the Step `spec.run` content or any hash of it

#### Scenario: Refresh rebuild cache
- **WHEN** runtime refresh is enabled through `KFG_REFRESH`
- **THEN** the runtime SHALL rebuild any existing cache entry for the Step invocation
- **AND** SHALL execute the Step again
- **AND** SHALL overwrite or replace the stored cache entry with the new result

### Requirement: Cache Persistence Semantics
The runtime SHALL persist only the observable Step results needed for downstream reuse.

#### Scenario: Declarative artifacts are cached
- **WHEN** a cacheable Step declares `spec.artifacts`
- **THEN** those artifact paths SHALL be included in the persisted cache entry for that Step invocation

#### Scenario: StepReference artifacts are cached
- **WHEN** a cacheable workflow Step reference declares `artifacts`
- **THEN** those artifact paths SHALL be included in the persisted cache entry for that Step invocation

#### Scenario: Runtime artifact registrations are cached
- **WHEN** a cacheable Step calls `__kfg_add_artifact` during execution
- **THEN** the runtime SHALL persist the artifact paths newly registered by that Step invocation
- **AND** SHALL NOT persist unrelated artifact paths that were already registered before the Step started

#### Scenario: Step outputs are cached automatically
- **WHEN** a cacheable Step declares `spec.output`
- **THEN** the runtime SHALL persist the captured output value as part of the cache entry
- **AND** no extra manifest field SHALL be required to opt in output persistence

#### Scenario: Output-producing step caches output and runtime artifacts together
- **WHEN** a cacheable Step declares `spec.output` and registers artifacts during execution
- **THEN** the runtime SHALL persist both the captured output value and the Step-local runtime artifact registrations from the same invocation

### Requirement: Cache Restore Semantics
Cache hits SHALL preserve workflow semantics for downstream Steps.

#### Scenario: Cached artifacts are restored to original relative paths
- **WHEN** a cacheable Step finds a matching cache entry
- **THEN** the runtime SHALL restore each persisted artifact to the same relative path it had when the cache entry was created

#### Scenario: Cached outputs are restored into runtime context
- **WHEN** a cacheable Step with `spec.output` is restored from cache
- **THEN** the runtime SHALL repopulate the Step output in the generated shell context
- **AND** later `when` conditions and `$kfg.output(...)` lookups SHALL observe the restored value

### Requirement: Cache Storage and Operations
The CLI SHALL expose internal operational management for cached Step entries.

#### Scenario: List cache entries
- **WHEN** user runs `kfg sys gc ls`
- **THEN** the CLI SHALL list cached Step entries with stable identifiers and operational metadata

#### Scenario: Inspect cache entry
- **WHEN** user runs `kfg sys gc inspect <id>`
- **THEN** the CLI SHALL show the stored metadata for that cache entry
- **AND** SHALL include the full relative paths of all persisted artifacts as recorded in `metadata.yaml`
- **AND** SHALL include output metadata

#### Scenario: Remove cache entry
- **WHEN** user runs `kfg sys gc rm <id>`
- **THEN** the CLI SHALL remove the specified cache entry from storage

#### Scenario: Prune cache entries
- **WHEN** user runs `kfg sys gc prune`
- **THEN** the CLI SHALL remove cache entries according to the implemented prune policy

#### Scenario: Show cache disk usage
- **WHEN** user runs `kfg sys gc du`
- **THEN** the CLI SHALL report disk usage for persisted cache entries

### Requirement: Cache Diagnostics
The runtime SHALL emit diagnostic logs for cache behavior.

#### Scenario: Cache hit and miss logging
- **WHEN** a cacheable Step checks for an existing cache entry
- **THEN** the runtime SHALL emit a detail log indicating whether the Step saw a cache hit or a cache miss

#### Scenario: Refresh rebuild logging
- **WHEN** cache rebuild occurs because refresh is enabled
- **THEN** the runtime SHALL emit a detail log explaining that cache was bypassed

#### Scenario: Store and restore path logging
- **WHEN** artifacts are stored into or restored from cache
- **THEN** the runtime SHALL emit debug logs that identify the artifact paths involved

### Requirement: Cache Persistence Format
The runtime SHALL store all cache entry metadata, including artifact paths, in a single `metadata.yaml` file within the cache entry directory.

#### Scenario: Artifact paths stored in metadata.yaml
- **WHEN** a cacheable Step completes execution and its cache entry is written
- **THEN** the persisted `metadata.yaml` SHALL include an `artifacts:` YAML list containing the relative path of every cached artifact
- **AND** the runtime SHALL NOT create any separated file for new cache entries

#### Scenario: Declarative artifacts with spaces in paths
- **WHEN** a Step or StepReference declares an artifact whose path contains spaces
- **THEN** the path SHALL be preserved correctly through the store and restore cycle without corruption

