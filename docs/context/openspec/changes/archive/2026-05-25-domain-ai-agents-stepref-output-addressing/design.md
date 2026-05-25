## Context

`kfg` currently treats a Step resource name as both the reusable implementation identity and the runtime execution identity. That works only while each workflow references a given Step once. In the agent workflow, the same Step resource is invoked multiple times with different `env` and different intended output consumers, so output lookups and `when.output` conditions become ambiguous.

The current shell generator also isolates `StepReference.env` by wrapping step calls in a subshell. That keeps env values from leaking, but it also causes any output written during the step invocation to be lost when the step itself produces `spec.output`. The ctx7 install and inject flow works around this by sharing a generated file path instead of sharing step output directly.

This change is cross-cutting because it touches manifest schema, workflow validation, resolved execution metadata, shell generation, and the ctx7 workflow content layer.

## Goals / Non-Goals

**Goals:**
- Make workflow step execution identity explicit with required `StepReference.name` values.
- Resolve output conditions and output reads by workflow step-reference name, not by reusable Step resource name.
- Allow workflow `env` values to read prior step outputs through `$kfg.output(<step-reference-name>)`.
- Preserve step-reference env isolation without losing produced outputs.
- Replace the ctx7 file-coupled handoff with output-to-env flow inside the workflow.

**Non-Goals:**
- Preserve compatibility with workflows that reference outputs by `Step.metadata.name`.
- Add multi-output Step resources or a new manifest field such as `inputs`.
- Change Cmd resource semantics or introduce persistent outputs across command invocations.

## Decisions

### Decision: Use `StepReference.name` as the only runtime output identity

Workflow step references already carry human-meaningful names in YAML, but the model ignores them. This change formalizes that field and makes it the only valid identifier for runtime outputs and `when.output.step` references.

Why this approach:
- It disambiguates multiple invocations of the same Step resource in the same workflow.
- It aligns the visible workflow structure with runtime behavior.
- It fixes the existing structural bug in `when.output`, not just the new env expansion use case.

Alternatives considered:
- Keep using `Step.metadata.name`: rejected because identical Step resources cannot be distinguished per invocation.
- Introduce a separate `outputAs` field: rejected because it duplicates the existing step-reference identity and creates two names for the same execution.

### Decision: Extend workflow `env` syntax instead of adding a new schema field

Dynamic output reads will use `$kfg.output(<step-reference-name>)` inside `StepReference.env` values. The generator will recognize this form and emit runtime lookup code.

Why this approach:
- It keeps the manifest model small.
- It reuses the existing env override mechanism that already scopes data per invocation.
- It matches the user's goal of avoiding a separate `inputs` feature.

Alternatives considered:
- Add `spec.inputs`: rejected as additional schema surface for a single data-flow need.
- Keep file-based handoff: rejected because it preserves hidden coupling and does not solve ambiguous step identity.

### Decision: Keep one reusable step function per Step resource and pass reference identity at invocation time

Generated shell will continue to define one function per Step resource, but each invocation will pass the executing `StepReference.name` into that function. Output storage and runtime lookups will key off the invocation identity.

Why this approach:
- It preserves deduplicated shell generation.
- It avoids generating duplicate functions for every workflow entry.
- It localizes the model change to runtime invocation metadata and output addressing.

Alternatives considered:
- Generate one function per StepReference: rejected because it duplicates code and makes multi-workflow output harder to reason about.

### Decision: Replace subshell-based step-reference env isolation with per-call inline environment assignment

The generator currently wraps step-reference env overrides in a subshell. The new implementation will instead emit inline environment assignment for the function call so that the invoked step sees isolated values while any output storage remains in the caller shell.

Why this approach:
- It preserves invocation-local env behavior.
- It avoids losing output writes performed by the step.
- It supports dynamic env expansion from earlier step outputs.

Alternatives considered:
- Keep subshells and try to synchronize output state back: rejected as brittle and more complex than needed.
- Export env globally before each step and manually restore: rejected because it increases leakage risk and rollback complexity.

### Decision: Move ctx7 handoff to output-to-env flow

`ctx7.steps.install` will emit the extracted context block as its step output. `ctx7.steps.inject` will consume that content through `CTX7_CONTEXT` supplied from workflow env expansion, and the step will no longer read `.$AGENT/ctx7-agents.md` itself.

Why this approach:
- It exercises the new output-addressing design on a real workflow need.
- It removes a hidden file dependency between two separate steps.
- It keeps install responsible for extracting ctx7-generated content and inject responsible for file upsert semantics.

Alternatives considered:
- Merge install and inject into one Step: rejected because it collapses separate responsibilities and makes the general output-flow feature unnecessary.

## Risks / Trade-offs

- [All existing `when.output.step` references must be updated] → Mitigation: validate strictly against `StepReference.name` and migrate the agent workflow in the same change.
- [Inline env assignment changes shell generation behavior] → Mitigation: add unit and Bats coverage for env isolation, output persistence, and sequential repeated invocations.
- [Dynamic env parsing could accept malformed expressions silently] → Mitigation: validate `$kfg.output(...)` syntax during generation or resolution and fail with explicit reference errors.
- [Multi-workflow generation may still deduplicate Step functions incorrectly if invocation metadata is not passed] → Mitigation: keep function deduplication but store invocation identity in resolved workflow step data and pass it on every generated call.

## Migration Plan

1. Update manifest schema and validation to require `StepReference.name` and to resolve output references only through that name.
2. Update resolver and shell generator to propagate step-reference identity and inline env assignment.
3. Add dynamic `$kfg.output(...)` env expansion and update `when.output` resolution.
4. Migrate the agent workflow to use named step references everywhere outputs are referenced.
5. Update ctx7 install/inject steps to use output-to-env handoff.
6. Run unit and Bats coverage for the changed execution model.

Rollback is straightforward during development because backward compatibility is not required; the change should ship as one coherent update of schema, generator, and manifests. Partial rollout is not desirable because old workflows will fail validation once the new model is enforced.

## Open Questions

- None. The change intentionally chooses strict step-reference addressing with no compatibility mode.
