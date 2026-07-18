## Why

Structured logs currently make Step attribution inconsistent because some runtime and manifest call sites encode the Step identity in the `component` string while others mention the Step only in the message text. That makes Step-scoped debugging, filtering, and future log analysis unreliable, and it leaves package manifests with multiple incompatible logging styles.

## What Changes

- Add an optional structured `step_name` field to runtime and CLI log events so Step identity is recorded separately from `component`.
- Export Step execution context from generated shell wrappers so `kfg sys log` can enrich Step-originated events automatically.
- Keep existing `__kfg_log_*` helper call patterns working, including one-argument message-only calls and legacy two-argument calls that use `step:<name>` as the component.
- Normalize legacy `step:<name>` shell log components into `component="step"` plus `step_name="<name>"` without breaking existing Steps.
- Establish a consistent message style for new logging work and fix known out-of-pattern Go log messages.
- Refactor framework and AI-agent manifests to stop encoding Step identity in `component` strings and rely on structured Step attribution instead.
- Update specs, docs, and tests for the new structured logging contract and compatibility behavior.

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- `shell-runtime-api`: generated shell runtime logs include Step execution identity through a dedicated `step_name` field while preserving existing helper compatibility.
- `cli-conventions`: `kfg sys log` accepts and normalizes Step-scoped runtime context instead of requiring Step identity to be encoded in the component string.
- `reusable-framework-steps`: framework Steps use the structured logging contract with automatic Step attribution and standardized message patterns.
- `skill-installation-steps`: AI-agent install Steps use the structured logging contract with automatic Step attribution and retain working behavior during migration.
- `inject-ctx7-context`: ctx7 injection logs use structured Step attribution instead of the legacy `step:<name>` component pattern.

## Impact

- Affects logger enrichment and CLI log handling in `src/internal/logger/` and `src/cmd/kfg/log.go`.
- Affects generated shell runtime templates in `src/internal/generate/templates/` and related generator tests.
- Affects framework manifests under `packages/framework/manifests/steps/`.
- Affects AI-agent manifests under `packages/domains/ai-agents/manifests/`.
- Requires unit, generator, and Bats coverage for Step log enrichment, legacy compatibility, and migrated manifest behavior.
