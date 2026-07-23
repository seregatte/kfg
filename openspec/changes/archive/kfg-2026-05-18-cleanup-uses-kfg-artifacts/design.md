## Context

`kfg.cleanup` is a shared manifest step, but its current contract is split across two mechanisms: artifact-producing steps append paths to `KFG_ARTIFACTS`, while `kfg.cleanup` still deletes paths from `PATHS`. This leaves generated shell state and cleanup behavior out of sync and forces workflow authors to keep a second cleanup list that is independent from actual produced artifacts.

The current dev workflow also contains `kfg.cleanup` call sites that pass `PATHS` in `agents-workflow.yaml`. If `kfg.cleanup` becomes artifact-driven only, those call sites must either stop using the shared step or register cleanup targets through the tracked artifact mechanism before invoking cleanup.

## Goals / Non-Goals

**Goals:**
- Make `kfg.cleanup` consume the same `KFG_ARTIFACTS` contract used by generated shell code.
- Remove the `PATHS` environment contract from the base cleanup step.
- Keep cleanup safe when no artifacts were recorded.
- Document the tracked-artifact cleanup contract in OpenSpec and cover it with tests.

**Non-Goals:**
- Introduce a new environment variable or alternate cleanup API.
- Export `KFG_ARTIFACTS` outside the current shell session.
- Preserve `PATHS`-based behavior for existing workflow call sites.

## Decisions

### Make `kfg.cleanup` iterate `KFG_ARTIFACTS` directly

The base step will become single-purpose: iterate the bash array and remove each recorded path with correct quoting. This keeps cleanup aligned with the shell generator contract and avoids the ambiguous word-splitting behavior of `rm -rf $PATHS`.

Alternative considered: continue supporting both `KFG_ARTIFACTS` and `PATHS`. Rejected because it preserves two overlapping contracts and keeps workflow behavior harder to reason about.

### Treat `PATHS` call sites as a migration concern, not a compatibility requirement

Workflow manifests that currently pass `PATHS` need to be updated alongside the step. The implementation should audit `kfg.cleanup` references and either:
- replace static path sweeps with tracked artifact registration plus `kfg.cleanup`, or
- inline static cleanup where the behavior is intentionally not artifact-driven.

Alternative considered: silently ignore `PATHS` without touching workflow manifests. Rejected because it would create a hidden regression in existing cleanup steps.

### Cover cleanup behavior at the manifest-step level

The highest-value regression tests are Bats tests that execute the cleanup step with populated `KFG_ARTIFACTS`, empty `KFG_ARTIFACTS`, and mixed file/directory targets. This validates the shell behavior actually consumed by generated workflows.

Alternative considered: rely only on generator or unit tests. Rejected because the bug is in shell-step behavior, not just Go-side generation.

## Risks / Trade-offs

- [Existing workflows rely on `PATHS`] -> Audit and update all `kfg.cleanup` call sites in the same implementation change.
- [Tracked artifacts may not cover all static cleanup paths] -> Decide case by case whether a path should become tracked or should be handled by a different explicit cleanup mechanism.
- [Array iteration may remove unexpected paths if artifacts are registered incorrectly] -> Keep tests focused on exact recorded paths and avoid broad glob-style cleanup in `kfg.cleanup`.

## Migration Plan

1. Change the base `kfg.cleanup` step to remove only `KFG_ARTIFACTS`.
2. Audit workflow manifest call sites that currently pass `PATHS`.
3. Convert those call sites to tracked-artifact cleanup or explicit non-shared cleanup behavior.
4. Add Bats coverage for the base step and any workflow behavior that changes as part of the migration.

## Open Questions

- Which current `PATHS`-based cleanup call sites are truly artifact-driven versus static stale-file cleanup?
- Should static pre-run cleanup remain in the dev workflow, or move to dedicated steps separate from `kfg.cleanup`?
