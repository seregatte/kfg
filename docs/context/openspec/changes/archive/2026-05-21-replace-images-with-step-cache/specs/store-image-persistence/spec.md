## REMOVED Requirements

### Requirement: Image persistence in store
**Reason**: Immutable image persistence is no longer part of the engine feature set.
**Migration**: Persist reusable Step artifacts and outputs through the Step cache store rooted under `KFG_STORE_DIR`.
