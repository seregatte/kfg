## MODIFIED Requirements

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

### Requirement: Cache Restore Semantics
Cache hits SHALL preserve workflow semantics for downstream Steps.

#### Scenario: Cached artifacts are restored to original relative paths
- **WHEN** a cacheable Step finds a matching cache entry
- **THEN** the runtime SHALL restore each persisted artifact to the same relative path it had when the cache entry was created

#### Scenario: Cached outputs are restored into runtime context
- **WHEN** a cacheable Step with `spec.output` is restored from cache
- **THEN** the runtime SHALL repopulate the Step output in the generated shell context
- **AND** later `when` conditions and `$kfg.output(...)` lookups SHALL observe the restored value

### Requirement: Cache Diagnostics
The runtime SHALL emit diagnostic logs for cache behavior.

#### Scenario: Cache hit and miss logging
- **WHEN** a cacheable Step checks for an existing cache entry
- **THEN** the runtime SHALL emit a detail log indicating whether the Step saw a cache hit or a cache miss

#### Scenario: Refresh bypass logging
- **WHEN** cache bypass occurs because refresh is enabled
- **THEN** the runtime SHALL emit a detail log explaining that cache was bypassed

#### Scenario: Store and restore path logging
- **WHEN** artifacts are stored into or restored from cache
- **THEN** the runtime SHALL emit debug logs that identify the artifact paths involved
