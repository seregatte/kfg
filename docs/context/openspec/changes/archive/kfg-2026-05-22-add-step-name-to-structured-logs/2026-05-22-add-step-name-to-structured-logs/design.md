## Context

The current logging contract separates level and component, but it does not provide a first-class field for Step identity. As a result, generated runtime code and package manifests have drifted into multiple conventions: some calls log with a generic component such as `cache` and embed the Step reference in the message text, while others pass `step:<name>` as the component, and many older Step calls still use message-only helper invocations. This inconsistency makes it hard to answer a simple debugging question such as "which Step emitted this log?" without parsing free-form strings.

This change crosses engine, framework, and domain layers. The engine owns JSONL shape, CLI log behavior, and generated runtime helpers. Framework and domain packages consume that API and currently encode Step identity in legacy ways. The design therefore has to add structured Step attribution without breaking already-shipped manifests and generated wrappers.

## Goals / Non-Goals

**Goals:**
- Add a dedicated `step_name` field to Step-originated log events.
- Preserve working behavior for existing `__kfg_log_*` call sites, including one-argument message-only usage and legacy `step:<name>` components.
- Keep `component` focused on subsystem identity rather than Step identity.
- Standardize the message style used by new runtime and Go log messages.
- Migrate framework and AI-agent manifests to the new structured logging contract.

**Non-Goals:**
- Redesign the full log storage format beyond the addition of `step_name` and related normalization.
- Change public Go logger call signatures.
- Require every shell log call to provide an explicit component.
- Introduce breaking manifest changes for existing packages or overlays.

## Decisions

### Decision: Add `step_name` as structured log context

The logger will enrich events with an optional `step_name` field when `KFG_STEP_NAME` is present in the environment, and the shell runtime will export that variable for the duration of each Step execution.

Why this approach:
- It makes Step attribution queryable without parsing `component` or `msg`.
- It fits the existing logger enrichment model used for `workflow_name`, `kustomization_name`, and `session_id`.
- It avoids changing Go logger call sites for non-Step events.

Alternatives considered:
- Encode Step identity only in `component`: rejected because it overloads `component` and keeps filtering brittle.
- Encode Step identity only in `msg`: rejected because it is unstructured and not reliably machine-readable.

### Decision: Preserve helper compatibility by supporting both current shell call forms

The generated `__kfg_log_*` helpers will continue to accept the existing shapes used in manifests today: a single message argument and a `component + message` pair. Message-only calls executed within Step context will default to `component="step"`, while two-argument legacy calls continue to work.

Why this approach:
- Existing framework and domain manifests already rely on both call forms.
- It avoids a flag day migration before the engine change lands.
- It lets new manifests adopt cleaner logging incrementally.

Alternatives considered:
- Require all manifests to switch to a new helper signature immediately: rejected because it would break existing Steps.
- Keep helper behavior ambiguous and undocumented: rejected because it would preserve the current drift.

### Decision: Normalize legacy `step:<name>` components at the logging backend

When shell logging receives a component in the legacy `step:<name>` form, the backend will interpret it as Step identity, emit `step_name=<name>`, and normalize `component` to `step` unless another subsystem component is explicitly provided.

Why this approach:
- It makes old manifests automatically produce the new structured shape.
- It keeps old logs meaningful during the migration window.
- It creates a single compatibility layer instead of scattering parsing logic across manifests.

Alternatives considered:
- Leave `step:<name>` untouched forever: rejected because it preserves the overloaded component model.
- Rewrite only manifests and skip backend normalization: rejected because users may have older generated shells or local manifests.

### Decision: Scope Step context to active Step execution and restore prior state

Generated Step wrappers will save any previous `KFG_STEP_NAME`, export the active Step reference name while the Step runs, and restore the previous value on exit.

Why this approach:
- It avoids leaking Step context into non-Step runtime logs.
- It supports nested helper usage and future composition safely.
- It aligns with the existing runtime practice of managing execution-scoped environment state.

Alternatives considered:
- Set `KFG_STEP_NAME` once for the whole wrapper invocation: rejected because command-level logs would be misattributed to the last Step.
- Never restore the previous value: rejected because later logs could inherit stale Step context.

### Decision: Standardize new messages to sentence case with clear verbs

New and touched log messages will use en-US sentence case, omit trailing periods for simple events, and place dynamic values in quotes where they improve readability.

Why this approach:
- It makes human-readable stderr output consistent across Go and shell logs.
- It reduces stylistic drift as new Steps are added.
- It works cleanly with structured fields carrying identity data such as `component` and `step_name`.

Alternatives considered:
- Preserve mixed historical styles: rejected because the change is already touching the logging contract and manifests.

## Risks / Trade-offs

- [Legacy helper call parsing may misclassify unusual one-argument component-only calls] -> Mitigation: constrain the compatibility logic to Step runtime usage, keep explicit two-argument calls authoritative, and add tests for the current manifest patterns.
- [Normalizing `step:<name>` changes the stored `component` value for migrated events] -> Mitigation: document the change in specs and keep legacy manifests functioning by translating them automatically.
- [Step context restoration bugs could leak wrong `step_name` values] -> Mitigation: save and restore previous values in generated wrappers and cover sequential and nested execution in tests.
- [Manifest refactors may miss some Step log call sites] -> Mitigation: update all current `packages/framework` and `packages/domains/ai-agents` Step manifests in the same change and validate with Bats.

## Migration Plan

1. Extend logger enrichment and `kfg sys log` handling to support `step_name` and legacy `step:<name>` normalization.
2. Update generated shell runtime helpers to export scoped `KFG_STEP_NAME` and preserve both current helper call forms.
3. Add or update unit and generator tests for structured Step attribution and compatibility behavior.
4. Refactor framework and AI-agent Step manifests to rely on automatic Step attribution and standardized messages.
5. Update OpenSpec deltas and user-facing docs to describe the new logging contract.
6. Run unit and Bats coverage for engine, framework, and domain logging scenarios.

## Open Questions

- None. The migration strategy intentionally preserves existing manifest behavior while converging on the structured `step_name` model.
