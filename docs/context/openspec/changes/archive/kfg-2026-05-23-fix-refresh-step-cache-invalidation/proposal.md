## Why

Step cache refresh currently describes itself as a cache bypass, but the intended behavior is step-scoped invalidation followed by a rebuild of that step's cache entry. That mismatch makes the runtime contract harder to reason about and has already led to confusion around whether refresh should preserve or replace cached artifacts.

## What Changes

- Redefine refresh semantics as invalidating the targeted step cache entry before the step runs, then rebuilding that entry from the new step result.
- Keep refresh scoped to the current step invocation rather than broadening invalidation to the whole workflow.
- Add an explicit runtime helper for registering artifacts from filesystem snapshot diffs rooted at a known directory.
- Update cache refresh diagnostics and user-facing help text to describe invalidation and rebuild behavior rather than generic bypass wording.
- Update the ctx7 install Step to use the new diff-based artifact helper instead of open-coded artifact registration loops.
- Add tests covering step-scoped refresh invalidation, refreshed cache reconstruction, explicit refresh logs, and diff-based artifact registration.

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- `step-cache`: refresh semantics invalidate the current step cache entry before execution and rebuild it from the new result, with explicit refresh diagnostics.
- `shell-runtime-api`: generated runtime exposes a helper that registers artifacts from filesystem snapshot diffs rooted at a provided directory.
- `cli-conventions`: refresh help text and command wording describe step-scoped invalidation and cache rebuild semantics.
- `apply-command`: apply-generated refresh state causes cacheable steps to invalidate and rebuild their own cache entries when executed.
- `run-command`: run-generated refresh state causes cacheable steps to invalidate and rebuild their own cache entries during that run.

## Impact

- Affects generated runtime templates in `src/internal/generate/templates/`.
- Affects CLI help and flag descriptions in `src/cmd/kfg/` and `docs/cli-reference.md`.
- Affects the ctx7 install manifest at `packages/domains/ai-agents/manifests/ctx7/steps/install.yaml`.
- Requires updates to OpenSpec delta files and tests covering runtime generation and Bats workflow behavior.
