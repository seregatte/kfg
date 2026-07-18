## Context

The current engine logging implementation already enriches events with workflow, kustomization, and session metadata, but it does not have a dedicated field for Step identity. That gap has pushed Step-aware runtime and manifest code toward overloaded `component` strings such as `step:ctx7.install` or toward free-form message text such as `Cache hit for step: ...`. The result is inconsistent logs that are harder to query and harder to standardize.

The engine also has spec drift around logging: the canonical CLI path is `kfg sys log`, runtime helpers are `__kfg_log_*`, and the actual JSONL and verbosity behavior differ from some older docs. Adding `step_name` is therefore best handled alongside a cleanup of the engine logging contract so the implementation and docs converge on the same shape.

## Goals / Non-Goals

**Goals:**
- Add `step_name` to shell log events when Step context is available.
- Preserve current shell helper call compatibility.
- Normalize legacy `step:<name>` shell components into structured Step attribution.
- Standardize touched Go log messages to the recommended message style.
- Align engine specs and docs with the real `kfg sys log` contract.

**Non-Goals:**
- Redesign the full JSONL schema beyond the targeted logging corrections.
- Change Go logger method signatures.
- Introduce breaking changes in shell helper names or required arguments.

## Decisions

### Decision: Reuse environment enrichment for `step_name`

The logger will treat `KFG_STEP_NAME` as another enriched context field and map it to `step_name` in JSONL output.

Why this approach:
- It matches the existing logger architecture.
- It keeps Step attribution orthogonal to `component`.
- It works for both normal and explicit-session CLI logging paths.

Alternatives considered:
- Add dedicated `step_name` parameters everywhere: rejected because it would require broader API churn.

### Decision: Normalize legacy Step components in the CLI logging path

The CLI logging command will interpret shell components matching `step:<name>` as legacy Step identity and rewrite them to `component="step"` with `step_name=<name>`.

Why this approach:
- Existing manifests keep working immediately.
- Structured data becomes consistent even before all manifests are migrated.

Alternatives considered:
- Leave legacy components untouched: rejected because it preserves the ambiguity this change is meant to solve.

### Decision: Standardize message text while touching logging call sites

Engine log messages updated by this change will follow sentence case, omit trailing periods for simple events, and quote dynamic values when useful.

Why this approach:
- The user request includes establishing a pattern for future work.
- The logging contract change is the right time to codify message consistency.

Alternatives considered:
- Defer message cleanup to a future change: rejected because the current work is already centered on logging consistency.

## Risks / Trade-offs

- [Spec cleanup may uncover more legacy logging drift than this change directly fixes] -> Mitigation: scope the engine spec changes to `step_name`, `kfg sys log`, and touched logger behavior while leaving unrelated cleanup out.
- [Legacy component normalization could surprise downstream parsers that read raw `component`] -> Mitigation: document the normalized shape in the log-command delta and retain compatibility at the call site level.
- [Message normalization may require test updates that compare literal output] -> Mitigation: update tests in the same change and keep the wording changes limited to touched call sites.

## Migration Plan

1. Update engine logger enrichment and CLI logging normalization.
2. Update generated runtime support for scoped `KFG_STEP_NAME`.
3. Update engine tests and docs to the new contract.
4. Validate with unit and Bats coverage before package-level manifest migration lands.

## Open Questions

- None.
