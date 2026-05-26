## Why

The cache layer never worked reliably. All store/restore logic lives in shell helpers embedded in generated templates, making it untestable, fragile, and hard to debug. Metadata serialization has recurring bugs (multi-line base64, double slashes in paths). The `kfg sys gc` command group uses opaque hash IDs instead of human-readable step names. The `kfg sys fs` subsystem exists only to support cache and adds unnecessary complexity. This change moves all cache logic into Go, exposes it through a unified `kfg sys cache` command namespace, and makes the shell a thin client that delegates to Go subcommands.

## What Changes

- **BREAKING**: Remove all `__kfg_cache_*` shell helper functions from `bash_helper.tmpl` — replaced by thin wrappers that call `kfg sys cache` subcommands
- **BREAKING**: Remove `__kfg_fs_snapshot`, `__kfg_fs_diff`, `__kfg_add_diff_artifacts` shell helpers — filesystem diff moves into Go
- **BREAKING**: Remove `kfg sys gc` command group — replaced by `kfg sys cache`
- **BREAKING**: Remove `kfg sys fs` command group — absorbed into Go `cache` package
- **BREAKING**: Remove `src/internal/store/` package (legacy v1 store, dead code)
- **BREAKING**: All `kfg sys cache` operations use `StepReference.name` as identifier (no hash IDs exposed)
- Create `src/internal/cache/` Go package with identity, metadata, artifacts, fsdiff, store, restore modules
- Create `kfg sys cache` command group: `exists`, `restore`, `store`, `ls`, `inspect`, `rm`, `prune`, `du`
- `kfg sys cache store` reads JSON from stdin (artifacts before/after, declarative artifacts, output) and does fs diff internally
- `kfg sys cache restore` emits shell eval-safe output (calls to `__kfg_output_set` and `__kfg_add_artifact`)
- All output subcommands (`ls`, `inspect`, `du`, `prune`) support `--json` and `--yaml` flags
- Simplify `ctx7/steps/install.yaml` and `openspec/steps/install.yaml` — remove manual fs snapshot/diff from step scripts
- Rewrite `tests/bats/workflows/step-cache-isolation.bats` for new architecture

## Capabilities

### New Capabilities

- `kfg-cache-sys-cache-command`: The unified `kfg sys cache` command group for all cache operations — runtime subcommands (`exists`, `restore`, `store`) consumed by shell wrappers and admin subcommands (`ls`, `inspect`, `rm`, `prune`, `du`) for cache management. All operations use `StepReference.name` as identifier. Output subcommands support `--json`/`--yaml` flags.

### Modified Capabilities

- `kfg-cache-step`: Cache persistence and restore semantics now execute through Go subcommands instead of shell helpers. The shell template calls `kfg sys cache exists/restore/store` via thin wrappers. Artifact detection uses hybrid approach (registered artifacts delta + Go fs diff). Requirements for identity, atomic write, diagnostics, and persistence format remain unchanged.
- `kfg-cache-internal-atomic-write`: Atomic write responsibility moves entirely to Go. The `store` subcommand writes to `.tmp` directory and renames on completion. No shell-side file operations for cache entries.

## Impact

**Affected code:**
- `src/internal/generate/templates/bash_helper.tmpl` — remove ~150 lines of shell cache/fs helpers, replace with thin Go wrappers
- `src/internal/generate/templates/bash_step.tmpl` — cache block delegates to shell wrappers (minimal change)
- `src/cmd/kfg/sys_gc.go` — removed entirely
- `src/cmd/kfg/sys_gc_test.go` — removed entirely
- `src/cmd/kfg/sys_fs.go` — removed entirely
- `src/cmd/kfg/sys.go` — register `sys_cache` instead of `sys_gc` and `sys_fs`
- `src/internal/store/` — removed entirely (dead code)
- New `src/internal/cache/` package (~6 files)
- New `src/cmd/kfg/sys_cache.go`

**Affected manifests:**
- `packages/domains/ai-agents/manifests/ctx7/steps/install.yaml` — remove fs snapshot/diff code
- `packages/domains/ai-agents/manifests/openspec/steps/install.yaml` — remove fs snapshot/diff code

**Affected tests:**
- `tests/bats/workflows/step-cache-isolation.bats` — rewrite 42 tests for new architecture

**Affected specs:**
- `docs/context/openspec/specs/kfg-cache-step/spec.md` — update for Go runtime
- `docs/context/openspec/specs/kfg-cache-sys-gc-command/spec.md` — removed, replaced by new spec
- `docs/context/openspec/specs/kfg-cache-internal-atomic-write/spec.md` — update for Go responsibility
