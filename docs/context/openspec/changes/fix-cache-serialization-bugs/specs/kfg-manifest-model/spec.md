## MODIFIED Requirements

### Requirement: Optional cache field
- GIVEN a Step resource
- WHEN `spec.cache` is specified
- THEN the schema SHALL accept `enabled`
- AND the cache configuration SHALL define the default cache behavior for workflow references to that Step