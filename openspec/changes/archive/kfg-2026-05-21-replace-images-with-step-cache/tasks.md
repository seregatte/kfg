## 1. Manifest and resolver model

- [x] 1.1 Add `cache` configuration types to `Step.spec` and `StepReference` in `src/internal/manifest/types.go` and validate supported fields.
- [x] 1.2 Extend resolved workflow data in `src/internal/resolve/resolve.go` so each `ResolvedStep` carries merged cache configuration with StepReference precedence.
- [x] 1.3 Update manifest validation and tests to cover cache field acceptance, StepReference overrides, and invocation-specific cache identity inputs.

## 2. Shell runtime and cache behavior

- [x] 2.1 Add shell runtime helpers for cache identity, artifact snapshotting, artifact restore, output serialization, and output restore in `src/internal/generate/templates/bash_helper.tmpl`.
- [x] 2.2 Integrate cache hit/miss behavior into generated Step execution in `src/internal/generate/templates/bash_step.tmpl` and `src/internal/generate/generate.go`.
- [x] 2.3 Persist the union of declarative artifacts, StepReference artifacts, and runtime artifact delta when a cacheable Step executes.
- [x] 2.4 Restore cached outputs into `__kfg_outputs` so `when` conditions and `$kfg.output(...)` continue to work on cache hits.

## 3. CLI refresh controls and GC commands

- [x] 3.1 Add `--refresh` support to `kfg run` and propagate `KFG_REFRESH` into generated shell execution.
- [x] 3.2 Add `--refresh` support to `kfg apply` so generated shell code carries refresh state.
- [x] 3.3 Implement `kfg sys gc ls`, `inspect`, `rm`, `prune`, and `du` for persisted Step cache entries under `KFG_STORE_DIR`.
- [x] 3.4 Add unit and integration coverage for refresh behavior and `sys gc` command output.

## 4. Logging helper normalization and manifest migration

- [x] 4.1 Rename generated shell logging helpers from `_kfg.log.*` to `__kfg_log_*` without changing the `kfg sys log` backend.
- [x] 4.2 Update framework and domain manifests to call the renamed `__kfg_log_*` helpers.
- [x] 4.3 Add generator and integration tests that validate the new helper names and successful Step execution through them.

## 5. Install Step caching rollout

- [x] 5.1 Add cache configuration to the slow install Steps in `packages/domains/ai-agents/manifests/*/steps/install.yaml` with stable per-step cache keys.
- [x] 5.2 Update install Step tests or add new coverage for cache hits, restored outputs, and refresh-driven reexecution.
- [x] 5.3 Validate that downstream workflow consumers such as ctx7 injection still work when the install Step is restored from cache.

## 6. Remove image and workspace systems

- [x] 6.1 Remove `src/cmd/kfg/image.go`, `src/cmd/kfg/workspace.go`, and any root command wiring or tests that reference them.
- [x] 6.2 Remove `src/internal/image/*`, `src/internal/imagefile/*`, and dependent tests or helpers.
- [x] 6.3 Remove or update Bats and integration coverage that expects image/workspace command behavior.

## 7. Specs, docs, and help text

- [x] 7.1 Update CLI help text for `kfg`, `build`, `apply`, and `run` to document public environment variables including `KFG_REFRESH`.
- [x] 7.2 Update engine-level OpenSpec specs and developer-facing documentation to describe Step cache, `sys gc`, and the removal of image/workspace features.
- [x] 7.3 Verify all updated repository-facing text remains in en-US and reflects the breaking removal of image/workspace commands.
