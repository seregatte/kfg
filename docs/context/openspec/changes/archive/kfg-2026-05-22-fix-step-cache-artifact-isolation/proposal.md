## Why

The initial Step cache implementation is functionally incorrect in three important ways: it stores artifacts from the whole invocation instead of the current Step, it restores cached artifacts into the current directory instead of their original relative paths, and it provides too little debug visibility to explain cache behavior. These issues make cache hits unsafe for real workflows and make failures hard to diagnose.

## What Changes

- Restrict Step cache persistence to artifacts produced by the current Step invocation instead of the full `KFG_ARTIFACTS` array.
- Preserve relative artifact paths when storing and restoring cache entries so cached directories and files return to their original locations.
- Add cache-specific `detail` and `debug` logging for refresh bypass, cache hit/miss, store, and restore flows.
- Update tests to cover artifact isolation, path-preserving restore, and cache logging.

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- `step-cache`: cache entries only persist Step-local artifacts, restore those artifacts to original relative paths, and emit diagnostic logging for cache behavior.
- `shell-runtime-api`: generated runtime helpers expose the corrected cache persistence semantics and cache-path logging.

## Impact

- Affects cache helpers and generated Step execution in `src/internal/generate/templates/`.
- Affects workflow wrapper artifact registration in `src/internal/generate/generate.go`.
- Requires generator, unit, and integration coverage for cache isolation and restore behavior.
