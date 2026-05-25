## Why

The engine already documents refresh as a cache-bypass mechanism that reexecutes the Step and replaces the stored result, but the generated runtime currently skips cache storage entirely when refresh is enabled. The engine also lacks a portable way to snapshot filesystem state for dynamic artifact discovery, and nested internal `kfg` subprocesses can flood stderr with startup logs when parent verbosity is set to `3`.

## What Changes

- Align generated refresh behavior with the cache contract by bypassing restore while still rebuilding and overwriting the cache entry after Step execution.
- Add a new internal CLI surface, `kfg sys fs`, that snapshots and diffs filesystem paths with stable normalized output and configurable depth.
- Extend the generated shell runtime with wrappers for `kfg sys fs` and for nested internal `kfg` execution.
- Scope nested internal `kfg` subprocesses to `KFG_VERBOSE=0` so child startup logs do not leak into the parent invocation's human output.
- Update engine tests and user-facing help text for refresh rebuild semantics and the new internal filesystem command surface.

## Capabilities

### New Capabilities
- `sys-fs-command`: internal CLI support for filesystem snapshot and diff operations used by generated runtime helpers.

### Modified Capabilities
- `shell-runtime-api`: runtime helpers include filesystem snapshot/diff wrappers and a quiet internal `kfg` execution wrapper.
- `cli-conventions`: the internal CLI surface documents `kfg sys fs`, and refresh wording reflects bypass-plus-rebuild behavior.

## Impact

- Affects `src/cmd/kfg/` command registration and help text.
- Affects generated runtime templates in `src/internal/generate/templates/`.
- Requires new Go helpers and tests for snapshot/diff normalization and cache overwrite behavior.
- Requires generator, unit, and Bats coverage for refresh rebuild semantics and nested `kfg` log suppression.
