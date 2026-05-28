## Why

Step cache refresh currently bypasses cache reuse but does not rebuild the stored entry, which leaves the implementation out of sync with the documented cache contract. At the same time, cacheable Steps still need verbose per-reference artifact declarations for directory-shaped outputs, and nested `kfg` subprocesses can re-emit startup logs when parent verbosity is high.

## What Changes

- Make refresh semantics rebuild and overwrite cache entries after re-running a cacheable Step instead of only bypassing cache restore.
- Add a platform-neutral internal filesystem API, `kfg sys fs`, for stable snapshot and diff operations with configurable depth.
- Extend the generated shell runtime with thin wrappers for filesystem snapshot/diff and for nested internal `kfg` execution.
- Keep declarative `Step.spec.artifacts` and `StepReference.artifacts` working while allowing Steps to register additional artifacts discovered through the runtime filesystem API.
- Update framework Steps that shell out to `kfg` so nested invocations suppress human startup logs without affecting the parent invocation's verbosity.
- Migrate the ctx7 install Step to discover new skill artifacts from `OUTPUT_DIR` instead of requiring repeated workflow-level artifact lists.
- Update specs, docs, and tests across engine, framework, and AI-agent roots for refresh rebuild behavior, filesystem helpers, and quiet internal subcommands.

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- `step-cache`: refresh-driven reexecution overwrites the stored cache entry completely instead of only bypassing cache reuse.
- `shell-runtime-api`: generated runtime exposes filesystem snapshot/diff helpers and a quiet internal `kfg` execution wrapper in addition to existing artifact registration helpers.
- `cli-conventions`: the internal CLI surface includes `kfg sys fs` and refresh-oriented help text reflects bypass-plus-rebuild behavior.
- `reusable-framework-steps`: framework Steps that invoke nested `kfg` commands use the runtime's internal execution wrapper so parent logs stay readable.
- `materialize-step`: nested conversion subcommands execute through the runtime wrapper without leaking child startup logs into the parent step output.
- `ctx7-install`: ctx7 installation discovers new skill artifacts from `OUTPUT_DIR` dynamically so cacheable workflow references do not need repeated artifact declarations.

## Impact

- Affects generated runtime templates in `src/internal/generate/templates/` and cache-related shell generation logic.
- Affects CLI command handling in `src/cmd/kfg/` for the new `sys fs` group and refresh help text.
- Affects framework manifests under `packages/framework/manifests/steps/`, especially `kfg.materialize`.
- Affects the AI-agent ctx7 install manifest and development overlay under `packages/domains/ai-agents/`.
- Requires unit, generator, and Bats coverage for filesystem snapshot/diff, cache overwrite semantics, quiet internal subcommands, and ctx7 artifact discovery.
