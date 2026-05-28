## Purpose

Define how cacheable Steps persist and restore their runtime results, including artifact isolation semantics, path preservation, and diagnostic logging.

## Requirements

### Requirement: Cache Identity and Invalidation
The runtime SHALL identify cache entries per resolved Step invocation using only the `StepReference.name`. The cache identity computation and existence check SHALL be performed by Go subcommands, with shell helpers acting as thin wrappers.

#### Scenario: Invocation-specific cache key
- **WHEN** a cacheable workflow Step executes
- **THEN** the cache identity SHALL be identified by `StepReference.name`
- **AND** SHALL NOT include any user-supplied key string
- **AND** SHALL NOT include the Step `spec.run` content or any hash of it

#### Scenario: Cache existence check via Go subcommand
- **WHEN** a cacheable Step checks for an existing cache entry
- **THEN** the shell helper `__kfg_cache_exists` SHALL delegate to `kfg sys cache exists <step-ref-name>`
- **AND** the subcommand SHALL exit 0 for hit and 1 for miss

#### Scenario: Refresh rebuild cache
- **WHEN** runtime refresh is enabled through `KFG_REFRESH`
- **THEN** the runtime SHALL rebuild any existing cache entry for the Step invocation
- **AND** SHALL execute the Step again
- **AND** SHALL overwrite or replace the stored cache entry with the new result

### Requirement: Cache Persistence Semantics
The runtime SHALL persist only the observable Step results needed for downstream reuse. Artifact detection uses a hybrid approach: registered artifact delta plus filesystem diff performed by the Go `store` subcommand.

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
- **AND** the shell wrapper SHALL pass the before/after artifact lists to `kfg sys cache store` via JSON on stdin

#### Scenario: Filesystem diff detects unregistered artifacts
- **WHEN** a cacheable Step creates files in the workdir but does not call `__kfg_add_artifact` for them
- **THEN** the Go `store` subcommand SHALL detect those files via filesystem diff and include them in the cache entry

#### Scenario: Step outputs are cached automatically
- **WHEN** a cacheable Step declares `spec.output`
- **THEN** the runtime SHALL persist the captured output value as part of the cache entry
- **AND** no extra manifest field SHALL be required to opt in output persistence

#### Scenario: Output-producing step caches output and runtime artifacts together
- **WHEN** a cacheable Step declares `spec.output` and registers artifacts during execution
- **THEN** the runtime SHALL persist both the captured output value and the Step-local runtime artifact registrations from the same invocation

### Requirement: Cache Restore Semantics
Cache hits SHALL preserve workflow semantics for downstream Steps. The restore is performed by a Go subcommand that emits shell eval-safe output.

#### Scenario: Cached artifacts are restored to original relative paths
- **WHEN** a cacheable Step finds a matching cache entry
- **THEN** the Go `restore` subcommand SHALL copy each persisted artifact to the same relative path it had when the cache entry was created
- **AND** SHALL emit `__kfg_add_artifact` calls to stdout for the shell to eval

#### Scenario: Cached outputs are restored into runtime context
- **WHEN** a cacheable Step with `spec.output` is restored from cache
- **THEN** the Go `restore` subcommand SHALL emit an `__kfg_output_set` call to stdout
- **AND** the shell wrapper SHALL eval the output to populate the runtime context
- **AND** later `when` conditions and `$kfg.output(...)` lookups SHALL observe the restored value

#### Scenario: Shell wrapper uses eval for restore
- **WHEN** the shell helper `__kfg_cache_restore` is called
- **THEN** it SHALL execute `eval "$(kfg sys cache restore <step-ref> --workdir "$PWD")"` to apply the restored state

### Requirement: Cache Storage and Operations
The CLI SHALL expose internal operational management for cached Step entries through `kfg sys cache` command group.

#### Scenario: List cache entries
- **WHEN** user runs `kfg sys cache ls`
- **THEN** the CLI SHALL list cached Step entries with step reference names and operational metadata

#### Scenario: Inspect cache entry
- **WHEN** user runs `kfg sys cache inspect <step-ref>`
- **THEN** the CLI SHALL show the stored metadata for that cache entry
- **AND** SHALL include the full relative paths of all persisted artifacts as recorded in `metadata.yaml`
- **AND** SHALL include output metadata

#### Scenario: Remove cache entry
- **WHEN** user runs `kfg sys cache rm <step-ref>`
- **THEN** the CLI SHALL remove the specified cache entry from storage

#### Scenario: Prune cache entries
- **WHEN** user runs `kfg sys cache prune`
- **THEN** the CLI SHALL remove cache entries according to the implemented prune policy

#### Scenario: Show cache disk usage
- **WHEN** user runs `kfg sys cache du`
- **THEN** the CLI SHALL report disk usage for persisted cache entries

### Requirement: Cache Diagnostics
The runtime SHALL emit diagnostic logs for cache behavior. Logging remains in the shell wrapper since it operates in the user's shell context.

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
The runtime SHALL store all cache entry metadata, including artifact paths, in a single `metadata.yaml` file within the cache entry directory. The format is unchanged from the current implementation.

#### Scenario: Artifact paths stored in metadata.yaml
- **WHEN** a cacheable Step completes execution and its cache entry is written
- **THEN** the persisted `metadata.yaml` SHALL include an `artifacts:` YAML list containing the relative path of every cached artifact
- **AND** the runtime SHALL NOT create any separated file for new cache entries

#### Scenario: Declarative artifacts with spaces in paths
- **WHEN** a Step or StepReference declares an artifact whose path contains spaces
- **THEN** the path SHALL be preserved correctly through the store and restore cycle without corruption

## REMOVED Requirements

### Requirement: Shell-side cache identity computation
**Reason**: Cache identity computation moves from shell (`sha256sum` in `__kfg_cache_identity`) to Go (`kfg sys cache exists` subcommand handles hashing internally).
**Migration**: The `__kfg_cache_identity` shell helper is removed. Shell calls `kfg sys cache exists <step-ref-name>` directly.

### Requirement: Shell-side filesystem snapshot and diff
**Reason**: `__kfg_fs_snapshot`, `__kfg_fs_diff`, and `__kfg_add_diff_artifacts` shell helpers are removed. Filesystem diff is performed internally by the Go `kfg sys cache store` subcommand.
**Migration**: Domain manifests (`ctx7/steps/install.yaml`, `openspec/steps/install.yaml`) no longer need manual snapshot/diff code. The Go `store` subcommand detects new files automatically.

