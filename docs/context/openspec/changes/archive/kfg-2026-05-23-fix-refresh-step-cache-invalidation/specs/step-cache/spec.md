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

#### Scenario: Refresh invalidates and rebuilds the current step cache entry
- **WHEN** runtime refresh is enabled through `KFG_REFRESH`
- **THEN** the runtime SHALL remove any existing cache entry for that Step invocation before executing the Step
- **AND** SHALL execute the Step again without restoring cached results
- **AND** SHALL rebuild the stored cache entry from the new result after successful execution

### Requirement: Cache Diagnostics

The runtime SHALL emit diagnostic logs for cache behavior.

#### Scenario: Cache hit and miss logging
- **WHEN** a cacheable Step checks for an existing cache entry
- **THEN** the runtime SHALL emit a detail log indicating whether the Step saw a cache hit or a cache miss

#### Scenario: Refresh invalidation logging
- **WHEN** refresh is enabled for a cacheable Step
- **THEN** the runtime SHALL emit a detail log that the current Step cache entry is being invalidated before execution
- **AND** SHALL emit a detail log that the cache entry is being rebuilt after successful execution

#### Scenario: Store and restore path logging
- **WHEN** artifacts are stored into or restored from cache
- **THEN** the runtime SHALL emit debug logs that identify the artifact paths involved
