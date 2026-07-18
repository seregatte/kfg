## ADDED Requirements

### Requirement: Declarative Step Cache Configuration
The manifest model SHALL allow Steps and workflow Step references to declare cache behavior.

#### Scenario: Step-level cache defaults
- **WHEN** a `Step` declares `spec.cache`
- **THEN** the schema SHALL accept `enabled` and `key` fields
- **AND** that cache configuration SHALL define the default behavior for every workflow reference to the Step

#### Scenario: StepReference cache override
- **WHEN** a workflow `StepReference` declares `cache`
- **THEN** the runtime SHALL use the StepReference cache values instead of the referenced Step defaults for that invocation
- **AND** the override SHALL apply only to that workflow step invocation

### Requirement: Cache Identity and Invalidation
The runtime SHALL identify cache entries per resolved Step invocation.

#### Scenario: Invocation-specific cache key
- **WHEN** a cacheable workflow Step executes
- **THEN** the cache identity SHALL include the workflow `StepReference.name`
- **AND** SHALL include the resolved cache key string
- **AND** SHALL include the Step `spec.run` content or a deterministic hash of it

#### Scenario: Script change invalidates cache
- **WHEN** a Step `spec.run` changes
- **THEN** previously stored cache entries for the old script SHALL NOT be reused

#### Scenario: Refresh bypasses cache
- **WHEN** runtime refresh is enabled through `KFG_REFRESH`
- **THEN** the runtime SHALL skip any existing cache entry for the Step invocation
- **AND** SHALL execute the Step again
- **AND** SHALL overwrite or replace the stored cache entry with the new result

### Requirement: Cache Persistence Semantics
The runtime SHALL persist the observable Step results needed for downstream reuse.

#### Scenario: Declarative artifacts are cached
- **WHEN** a cacheable Step declares `spec.artifacts`
- **THEN** those artifact paths SHALL be included in the persisted cache entry

#### Scenario: StepReference artifacts are cached
- **WHEN** a cacheable workflow Step reference declares `artifacts`
- **THEN** those artifact paths SHALL be included in the persisted cache entry for that invocation

#### Scenario: Runtime artifact registrations are cached
- **WHEN** a cacheable Step calls `__kfg_add_artifact` during execution
- **THEN** the runtime SHALL persist any newly registered artifact paths produced by that Step invocation

#### Scenario: Step outputs are cached automatically
- **WHEN** a cacheable Step declares `spec.output`
- **THEN** the runtime SHALL persist the captured output value as part of the cache entry
- **AND** no extra manifest field SHALL be required to opt in output persistence

### Requirement: Cache Restore Semantics
Cache hits SHALL preserve workflow semantics for downstream Steps.

#### Scenario: Cached artifacts are restored before downstream execution
- **WHEN** a cacheable Step finds a matching cache entry
- **THEN** the runtime SHALL restore the persisted artifact paths into the working tree before later Steps execute

#### Scenario: Cached outputs are restored into runtime context
- **WHEN** a cacheable Step with `spec.output` is restored from cache
- **THEN** the runtime SHALL repopulate the Step output in the generated shell context
- **AND** later `when` conditions and `$kfg.output(...)` lookups SHALL observe the restored value

#### Scenario: Multiline outputs survive cache round trip
- **WHEN** a Step output contains multiline or special-character content
- **THEN** the cache serialization format SHALL preserve the full value without truncation or corruption

### Requirement: Cache Storage and Operations
The CLI SHALL expose internal operational management for cached Step entries.

#### Scenario: List cache entries
- **WHEN** user runs `kfg sys gc ls`
- **THEN** the CLI SHALL list cached Step entries with stable identifiers and operational metadata

#### Scenario: Inspect cache entry
- **WHEN** user runs `kfg sys gc inspect <id>`
- **THEN** the CLI SHALL show the stored metadata for that cache entry
- **AND** SHALL include persisted artifacts and output metadata

#### Scenario: Remove cache entry
- **WHEN** user runs `kfg sys gc rm <id>`
- **THEN** the CLI SHALL remove the specified cache entry from storage

#### Scenario: Prune cache entries
- **WHEN** user runs `kfg sys gc prune`
- **THEN** the CLI SHALL remove cache entries according to the implemented prune policy

#### Scenario: Show cache disk usage
- **WHEN** user runs `kfg sys gc du`
- **THEN** the CLI SHALL report disk usage for persisted cache entries
