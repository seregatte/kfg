## REMOVED Requirements

### Requirement: Declarative Step Cache Configuration
**Reason**: The `cache.key` field is removed from the manifest model. Cache identity is now determined solely by `StepReference.name`, making the `key` field redundant. The `enabled` field is retained.
**Migration**: Remove `key:` from any `cache:` block in Step or StepReference manifests. Cache identity is now `StepReference.name` only.

## MODIFIED Requirements

### Requirement: Cache Identity and Invalidation
The runtime SHALL identify cache entries per resolved Step invocation using only the `StepReference.name`.

#### Scenario: Invocation-specific cache key
- **WHEN** a cacheable workflow Step executes
- **THEN** the cache identity SHALL be computed as `SHA256(StepReference.name)`
- **AND** SHALL NOT include any user-supplied key string
- **AND** SHALL NOT include the Step `spec.run` content or any hash of it

#### Scenario: Cache is NOT automatically invalidated on script change
- **WHEN** a Step `spec.run` changes
- **THEN** previously stored cache entries SHALL remain valid and SHALL be reused on the next execution
- **AND** the user SHALL use `--refresh` or change the `StepReference.name` to force re-execution

#### Scenario: Refresh bypasses cache
- **WHEN** runtime refresh is enabled through `KFG_REFRESH`
- **THEN** the runtime SHALL skip any existing cache entry for the Step invocation
- **AND** SHALL execute the Step again
- **AND** SHALL overwrite or replace the stored cache entry with the new result

### Requirement: Cache Persistence Format
The runtime SHALL store all cache entry metadata, including artifact paths, in a single `metadata.yaml` file within the cache entry directory.

#### Scenario: Artifact paths stored in metadata.yaml
- **WHEN** a cacheable Step completes execution and its cache entry is written
- **THEN** the persisted `metadata.yaml` SHALL include an `artifacts:` YAML list containing the relative path of every cached artifact
- **AND** the runtime SHALL NOT create a separate `artifact_paths.txt` file for new cache entries

#### Scenario: Backward compatible restore from artifact_paths.txt
- **WHEN** a cacheable Step finds a cache entry that contains `artifact_paths.txt` but no `artifacts:` key in `metadata.yaml`
- **THEN** the runtime SHALL read artifact paths from `artifact_paths.txt` and restore them correctly

#### Scenario: Declarative artifacts with spaces in paths
- **WHEN** a Step or StepReference declares an artifact whose path contains spaces
- **THEN** the path SHALL be preserved correctly through the store and restore cycle without corruption

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
