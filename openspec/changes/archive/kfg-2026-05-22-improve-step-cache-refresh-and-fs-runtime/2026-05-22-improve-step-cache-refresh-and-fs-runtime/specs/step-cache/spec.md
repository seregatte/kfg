## MODIFIED Requirements

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

#### Scenario: Cache rewrite removes stale artifacts
- **WHEN** a cacheable Step rewrites an existing cache entry with a smaller artifact set
- **THEN** the stored cache entry SHALL no longer contain artifacts omitted from the new result
