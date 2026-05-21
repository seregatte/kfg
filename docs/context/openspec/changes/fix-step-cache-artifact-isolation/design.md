## Context

The first version of Step cache introduced runtime persistence for expensive install Steps, but the implementation currently treats `KFG_ARTIFACTS` as if it were the artifact set for the active Step. In reality, the generated command wrapper pre-registers artifacts for the command and all dependent Steps, so by the time a cacheable Step stores its result, the global array already contains unrelated entries. This causes cache entries to absorb artifacts from prior and later Steps.

The same implementation also stores artifacts under `basename(path)` and restores them into the current working directory. Any nested artifact such as `.pi/skills/foo` or two artifacts with the same basename from different directories loses its original location. Finally, the cache flow has almost no logging, so users cannot see whether a Step hit cache, missed cache, or restored the wrong paths.

This fix is scoped to correctness and diagnosability of the existing Step cache feature. It does not change the cache model, CLI surface, or refresh controls introduced by the parent cache change.

## Goals / Non-Goals

**Goals:**
- Persist only artifacts attributable to the current Step invocation.
- Restore cached artifacts to their original relative paths.
- Emit useful `detail` and `debug` logs for cache identity, hit/miss, store, restore, and refresh bypass.
- Add regression coverage for artifact isolation and path preservation.

**Non-Goals:**
- Change cache identity semantics, refresh semantics, or `sys gc` commands.
- Introduce new manifest fields.
- Redesign the broader artifact model outside the cache bugfix.

## Decisions

### Decision: Compute a per-Step artifact delta

Cache store logic will snapshot artifact state before the Step runs and compare it with artifact state after the Step finishes. The cache entry will persist only the new artifacts attributable to that Step, plus any explicit Step and StepReference artifact declarations relevant to that invocation.

Why this approach:
- It matches the intended contract of Step-local cache persistence.
- It works even when a Step dynamically registers artifacts via `__kfg_add_artifact`.
- It avoids relying on the global `KFG_ARTIFACTS` array as if it were Step-scoped.

Alternatives considered:
- Persist the entire array and filter later: rejected because it retains incorrect coupling and bloats cache entries.
- Use only declarative artifacts: rejected because several Steps register artifacts dynamically.

### Decision: Preserve full relative paths in cache storage

The cache helper will store artifacts using their original relative paths beneath the cache entry and restore them to the same paths. Restore will create parent directories before copying content back into the workspace.

Why this approach:
- It preserves nested directories exactly.
- It avoids basename collisions.
- It matches user expectations for cached generated content.

Alternatives considered:
- Continue using basenames: rejected because it is lossy and incorrect.

### Decision: Reduce wrapper-level pre-registration for Step artifacts

The command wrapper should not pre-register Step artifacts globally before the Step executes, because doing so pollutes the artifact state used by cache persistence. Step artifacts should be registered when the Step itself runs or restores from cache.

Why this approach:
- It keeps `KFG_ARTIFACTS` closer to actual runtime production order.
- It allows cache delta logic to be meaningful.

Alternatives considered:
- Keep pre-registration and try to subtract known wrapper artifacts: rejected as more fragile than not pre-registering them.

### Decision: Instrument the cache path with runtime logs

Generated cache helpers will emit `detail` logs for refresh bypass, cache hit, cache miss, store start/end, and restore start/end. They will emit `debug` logs for cache path/identity and individual artifact paths being stored or restored.

Why this approach:
- Cache behavior is otherwise opaque during debugging.
- The logs make it much easier to diagnose bad cache contents or restore errors.

Alternatives considered:
- Only log errors: rejected because cache correctness issues often present as wrong-but-successful behavior.

## Risks / Trade-offs

- [Changing wrapper artifact registration may affect cleanup behavior] -> Mitigation: add tests that verify cleanup still sees Step artifacts after execution and cache restore.
- [Path-preserving restore could overwrite unintended files if paths are wrong] -> Mitigation: keep cache paths relative, preserve only registered artifact paths, and add regression tests for nested paths and basename collisions.
- [Extra logs could be noisy] -> Mitigation: keep operational messages at `detail` and path-level messages at `debug`.

## Migration Plan

1. Update generated wrappers so Step artifacts are not globally pre-registered before execution.
2. Update cache helpers to compute per-Step artifact delta and store artifact paths with their original relative locations.
3. Update restore logic to recreate parent directories and restore outputs with accompanying logs.
4. Add regression tests for isolated artifacts, path-preserving restore, and cache logging.
5. Re-run cache-related integration scenarios for install Steps.

## Open Questions

- None. The fix is intentionally limited to correctness and observability of the existing Step cache feature.
