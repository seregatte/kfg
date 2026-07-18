## Context

The current Step cache flow already snapshots `KFG_ARTIFACTS` around cacheable Step execution and stores the union of declarative artifacts plus the Step-local runtime delta. That model is close to the intended contract, but refresh handling is incomplete: when `KFG_REFRESH` is set, the generated wrapper bypasses cache restore and also skips cache storage, so the old cache entry remains unchanged even after a successful rerun.

At the same time, reusable Steps that produce directory-shaped outputs often need workflow-level `StepReference.artifacts` declarations because the runtime has no portable way to discover newly created paths from a directory before and after execution. The immediate pain point is `ctx7.steps.install`, where each agent-specific workflow reference repeats the final skill path. A related ergonomics issue is that Steps like `kfg.materialize` shell out to `kfg` internally, and those subprocesses inherit `KFG_VERBOSE=3`, causing child startup logs to appear inside the parent Step output.

This change crosses engine, framework, and domain layers. The engine owns refresh semantics, CLI command shape, and runtime helpers. The framework consumes the runtime to perform nested `kfg` conversions. The AI-agent domain consumes the runtime for artifact registration and can simplify overlays once filesystem discovery is available. The design therefore has to centralize the portable behavior in Go while keeping the shell runtime as a thin orchestration layer.

## Goals / Non-Goals

**Goals:**
- Make refresh execute cacheable Steps and replace the stored cache entry with the new result.
- Add a Go-owned internal filesystem snapshot/diff API that the shell runtime can call on any platform.
- Keep `Step.spec.artifacts`, `StepReference.artifacts`, and runtime-discovered artifacts working together.
- Provide a single runtime wrapper for nested internal `kfg` subprocesses that suppresses child human logs without mutating parent verbosity.
- Use the new filesystem API to remove repeated ctx7 workflow artifact declarations.

**Non-Goals:**
- Redesign cache identity, `sys gc`, or the broader manifest model beyond the minimum internal helper surface needed here.
- Replace the shell runtime with a native Go Step executor.
- Remove declarative artifact support from Steps or workflow references.
- Silence arbitrary third-party subprocesses invoked by Steps; this change only scopes nested internal `kfg` calls.

## Decisions

### Decision: Separate refresh restore control from cache store control

Generated Step wrappers will continue to skip cache restore when `KFG_REFRESH` is set, but they will always perform cache storage after a successful execution when caching is enabled.

Why this approach:
- It aligns the implementation with the existing Step cache spec, which already requires refresh to overwrite the stored result.
- It keeps the refresh behavior easy to reason about: rerun now, use the rerun result next time.
- It minimizes the change to one branch in the generated wrapper rather than introducing a parallel refresh-specific cache path.

Alternatives considered:
- Keep refresh as bypass-only: rejected because it contradicts the documented behavior and leaves stale cache entries in place.
- Add a separate refresh cache namespace: rejected because refresh is an invalidation control, not a new identity dimension.

### Decision: Clear the cache entry before writing refreshed or updated content

`__kfg_cache_store` will remove any existing cache entry directory before recreating it for the new write.

Why this approach:
- It guarantees that refresh and normal rewrites replace the full cache entry instead of leaving stale artifact files behind.
- It keeps the cache store format unchanged while making overwrite semantics correct.
- It is simpler and less error-prone than trying to reconcile the previous artifact tree with the new artifact list.

Alternatives considered:
- Delete only paths not present in the new artifact list: rejected because it complicates the shell implementation and still requires full metadata reconciliation.
- Keep old files and rely on `artifact_paths.txt`: rejected because stale on-disk files would remain available to future restore logic and operational tooling.

### Decision: Introduce `kfg sys fs` as the portable filesystem snapshot backend

The engine will add a new internal CLI group, `kfg sys fs`, with `snapshot` and `diff` subcommands. The generated shell runtime will wrap those commands instead of implementing filesystem walking and diff logic in bash.

Why this approach:
- It centralizes platform-sensitive filesystem traversal and path normalization in Go, which is the right place for future Windows support.
- It keeps the shell runtime API stable and readable while still allowing the current bash backend to orchestrate Step execution.
- It mirrors the successful `kfg sys log` pattern, where shell helpers delegate durable behavior to the CLI.

Alternatives considered:
- Implement snapshot/diff purely in shell with `find`, `comm`, or `diff`: rejected because that hardens a bash-specific contract and makes future Windows support harder.
- Make Steps call Go code directly without CLI mediation: rejected because the runtime today is generated shell, so the CLI is the cleanest integration point.

### Decision: Keep `kfg sys fs` text-oriented and depth-bounded

`kfg sys fs snapshot <path> [--maxdepth N]` will print normalized relative paths, one per line, sorted deterministically. `--maxdepth 0` means unlimited depth. `kfg sys fs diff --before <snapshot> --after <snapshot>` will print only paths present in `after` and absent in `before`.

Why this approach:
- Text snapshots are easy for shell wrappers to capture and pass around without introducing temporary state files or `eval`-style APIs.
- Relative, normalized paths make snapshots portable across platforms and straightforward to prefix back onto a base directory.
- `--maxdepth` allows Steps such as ctx7 install to observe only top-level skill directories while leaving deeper traversal available for future use cases.

Alternatives considered:
- Return JSON snapshots: rejected because the runtime only needs simple line-based data and JSON would add parsing overhead in shell.
- Hardcode a single depth policy: rejected because different Steps may need immediate children versus full tree traversal.

### Decision: Keep artifact registration additive across declarative and discovered sources

The runtime will continue to treat declarative Step artifacts, StepReference artifacts, and dynamic registrations as a union. The new filesystem API only supplies additional paths for `__kfg_add_artifact`; it does not replace existing declarative behavior.

Why this approach:
- It preserves compatibility with current manifests and cache semantics.
- It lets Steps move from repeated workflow declarations to dynamic discovery incrementally.
- It avoids forcing every cacheable Step into a directory-diff pattern when some already know their output paths statically.

Alternatives considered:
- Deprecate `StepReference.artifacts` immediately: rejected because some workflows still need explicit artifact paths and the spec still supports them.
- Prefer discovered artifacts over declarative paths: rejected because declarative paths are part of the current contract and may describe files not created during the current invocation.

### Decision: Add one quiet internal `kfg` execution wrapper to the shell runtime

The runtime will expose a single helper for nested engine subprocesses, with semantics equivalent to `KFG_VERBOSE=0 kfg "$@"` scoped only to that child process. Runtime wrappers such as `__kfg_fs_snapshot` and framework Steps such as `kfg.materialize` will use that helper for internal `kfg` commands.

Why this approach:
- It solves the immediate startup-log problem for both existing nested `kfg apply` calls and the new `kfg sys fs` calls.
- It centralizes the policy so manifests do not need to repeat `KFG_VERBOSE=0` by hand.
- It preserves the parent invocation's verbosity and structured logs.

Alternatives considered:
- Silence all subprocess stderr in the Step: rejected because it would hide genuine failures from users.
- Introduce a new public environment variable for child log suppression: rejected because the behavior is internal runtime plumbing, not a public user control.

### Decision: Use ctx7 install as the first dynamic artifact discovery consumer

`ctx7.steps.install` will snapshot `OUTPUT_DIR` before and after installation with `--maxdepth 1`, diff the snapshots, and register the newly created children under `OUTPUT_DIR` as artifacts. The workflow-level ctx7 `artifacts:` lists can then be removed.

Why this approach:
- It directly addresses the most repetitive artifact declarations in the current overlay.
- The ctx7 install output shape maps naturally to top-level skill directories under `OUTPUT_DIR`.
- It validates the new engine/runtime API in a real domain Step before wider rollout.

Alternatives considered:
- Migrate `openspec.steps.install` first: rejected because it currently registers a broader directory path and is less representative of the per-skill repetition problem.
- Add manifest placeholders in `artifacts`: rejected because it would still tie the solution to manifest authoring rather than runtime discovery.

## Risks / Trade-offs

- [Refreshing a cache entry could briefly remove the old entry before the new write completes] -> Mitigation: keep the store flow single-process, recreate the entry immediately, and rely on rerun determinism for the next invocation.
- [Text snapshot transport might mishandle unusual paths if normalization is inconsistent] -> Mitigation: normalize separators in Go, sort deterministically, and cover paths with nested directories in unit tests.
- [Nested internal `kfg` commands with `KFG_VERBOSE=0` could hide child human logs users expected to see] -> Mitigation: scope the helper to internal engine subprocesses only and preserve stdout/stderr failure propagation.
- [ctx7 diffing at depth 1 could miss future nested output layouts] -> Mitigation: keep `--maxdepth` configurable and document why ctx7 intentionally uses depth 1 for current skill directories.

## Migration Plan

1. Add the engine-side `kfg sys fs` command group and Go helpers for snapshot and diff.
2. Update runtime helpers to include the quiet internal `kfg` wrapper plus `fs` snapshot/diff wrappers.
3. Fix cache refresh/store semantics and ensure cache entries are cleared before rewrite.
4. Update framework Steps that invoke nested `kfg` commands to use the internal execution wrapper.
5. Update ctx7 install to use directory snapshot/diff for artifact registration and remove redundant workflow artifact declarations.
6. Add unit, generator, and Bats coverage for refresh rebuilds, `sys fs`, quiet internal subprocesses, and ctx7 artifact discovery.

## Open Questions

- None. The command shape, refresh semantics, and ctx7 migration scope have been decided for this change.
