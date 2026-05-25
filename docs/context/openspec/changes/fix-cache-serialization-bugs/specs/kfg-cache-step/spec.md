## MODIFIED Requirements

### Requirement: Cache Identity and Invalidation
The runtime SHALL identify cache entries per resolved Step invocation using only the `StepReference.name`.

#### Scenario: Invocation-specific cache key
- **WHEN** a cacheable workflow Step executes
- **THEN** the cache identity SHALL be identified by `StepReference.name`
- **AND** SHALL NOT include any user-supplied key string
- **AND** SHALL NOT include the Step `spec.run` content or any hash of it