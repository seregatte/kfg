## Context

The current runtime already computes cache identity per workflow step invocation and writes cache entries after successful execution. The open question is not identity scope but refresh semantics: the user wants refresh to mean that the currently targeted step cache entry is invalidated and rebuilt, not merely described as a bypass. The same change also exposes repeated shell code in Steps such as `ctx7.steps.install`, where directory snapshots are diffed and each discovered path is manually prefixed and registered as an artifact.

The user clarified two important constraints. First, refresh invalidation must remain scoped to the current step cache entry rather than wiping related workflow entries. Second, the workspace can be treated as clean at the beginning of a new run because `kfg.cleanup` is responsible for removing registered artifacts in `after` steps. That lets the runtime invalidate the current step cache entry before execution without needing workflow-wide store coordination.

## Goals / Non-Goals

**Goals:**
- Make refresh semantics explicit: invalidate the current step cache entry, execute the step, and rebuild the cache entry from the new result.
- Keep refresh invalidation scoped to the current cacheable step invocation.
- Add a reusable runtime helper that registers artifacts from filesystem snapshot diffs rooted at a caller-supplied directory.
- Replace open-coded diff-to-artifact loops in ctx7 install with the shared runtime helper.
- Update logs, docs, and specs so refresh is described as invalidation plus rebuild rather than bypass-only behavior.

**Non-Goals:**
- Redesign cache identity or widen invalidation to all steps in a workflow.
- Change the broader cache persistence model away from step-scoped observable results.
- Introduce a generic artifact mutation framework beyond the narrow diff-registration helper needed here.
- Rework unrelated agent install steps unless they need the new helper immediately.

## Decisions

### Decision: Invalidate the current step cache entry before refresh execution

When `KFG_REFRESH` is set for a cacheable step, the generated wrapper will compute `__cache_path`, emit an explicit refresh invalidation log, remove that path, skip restore, and then execute the step. After a successful execution it will emit a rebuild log and store the new cache entry as usual.

Why this approach:
- It matches the requested semantics exactly: refresh invalidates the targeted step cache entry and rebuilds it.
- It preserves current cache identity and keeps the operation scoped to the step rather than the workflow.
- It relies on the agreed workspace-clean assumption instead of inventing workflow-level cleanup logic inside the runtime.

Alternatives considered:
- Leave refresh as restore bypass only: rejected because it keeps the language and behavior underspecified.
- Remove all workflow-related cache entries before execution: rejected because cache identity is step-scoped and the user explicitly wants step scope preserved.

### Decision: Make refresh diagnostics explicit about invalidation and rebuild

Refresh-related logs and help text will use explicit terms such as `invalidating cache` and `rebuilding cache` instead of `bypass` as the primary description.

Why this approach:
- It matches the intended mental model and avoids the ambiguity that started this change.
- It makes generated runtime logs easier to interpret when diagnosing refresh behavior.

Alternatives considered:
- Keep `bypass` wording and rely on secondary text to explain rebuilds: rejected because it keeps the misleading shorthand in the most visible places.

### Decision: Add a root-aware diff artifact registration helper to the runtime

The generated runtime will add a helper such as `__kfg_add_diff_artifacts <root> <before_snapshot> <after_snapshot>`. The helper will run `__kfg_fs_diff`, prefix each returned relative path with the provided root, verify existence, and register the resulting artifact path with `__kfg_add_artifact`.

Why this approach:
- `__kfg_fs_diff` returns paths relative to the snapshot root, so the helper must know that root to reconstruct real workspace paths.
- It avoids fragile shell APIs based on command substitution and word-splitting such as `__kfg_add_artifacts $(...)`.
- It centralizes repeated loop logic in the runtime where other Steps can reuse it.

Alternatives considered:
- Use `__kfg_add_artifacts $(__kfg_fs_diff ...)`: rejected because the helper would not know which root to prefix and would rely on brittle shell splitting semantics.
- Keep open-coded loops in each Step: rejected because the same pattern will be needed in additional Steps.

### Decision: Migrate ctx7 install first and keep the helper narrowly scoped

`ctx7.steps.install` will adopt the new diff-based helper immediately. The helper remains narrowly focused on converting snapshot diffs into artifact registrations rather than becoming a broad multi-purpose artifact API.

Why this approach:
- It removes the concrete duplication already identified by the user.
- It gives the runtime helper an immediate production consumer without broadening scope unnecessarily.

Alternatives considered:
- Delay migration until multiple Steps need it: rejected because ctx7 already contains the exact duplicated pattern.

## Risks / Trade-offs

- [Pre-execution invalidation removes the old step cache entry before the new execution succeeds] -> Mitigation: rely on the agreed clean-workspace model and cover refresh failure/rebuild flows in tests so the semantics stay intentional and visible.
- [Diff helper callers could pass the wrong root and register invalid artifact paths] -> Mitigation: keep the helper root parameter explicit, verify existence before registration, and test with nested relative paths.
- [Docs or tests may continue using bypass terminology in older text] -> Mitigation: update all touched help, specs, and runtime log assertions in the same change.

## Migration Plan

1. Update generated step wrappers so refresh invalidates only the current step cache entry before execution and logs invalidation/rebuild explicitly.
2. Add the diff-based artifact registration helper to the shell runtime.
3. Replace the ctx7 install loop with the shared helper.
4. Update docs, specs, and tests to match the new refresh wording and helper contract.
5. Run unit and Bats coverage to confirm refreshed cache entries and diff-based artifacts are rebuilt correctly.

## Open Questions

- None. Step scope, explicit wording, and root-aware diff registration were all decided during discussion.
