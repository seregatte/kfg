## 1. Manifest and resolver model

- [x] 1.1 Add `name` to `StepReference` in `src/internal/manifest/types.go` and validate it as required for workflow step references.
- [x] 1.2 Validate `StepReference.name` uniqueness within each `CmdWorkflow` and fail on duplicates with explicit errors.
- [x] 1.3 Update output-reference validation so `when.output.step` resolves only through `StepReference.name` in the same workflow.
- [x] 1.4 Validate `$kfg.output(<step-reference-name>)` env references against existing named workflow step references and referenced Step outputs.

## 2. Shell generation and runtime output addressing

- [x] 2.1 Extend resolved workflow step data to carry `StepReference.name` as the runtime execution identity.
- [x] 2.2 Update generated step functions to store outputs under `StepReference.name` instead of the underlying Step resource name.
- [x] 2.3 Replace subshell-based step-reference env handling with per-invocation inline env assignment that preserves output writes.
- [x] 2.4 Add generator support for expanding `$kfg.output(<step-reference-name>)` inside workflow step-reference `env` values.
- [x] 2.5 Update `when.output` code generation to look up outputs by `StepReference.name`.

## 3. Workflow and ctx7 migration

- [x] 3.1 Migrate `packages/domains/ai-agents/overlays/dev/agents-workflow.yaml` so all workflow step references have required names and all `when.output.step` values use those names.
- [x] 3.2 Update `packages/domains/ai-agents/manifests/ctx7/steps/install.yaml` to emit extracted ctx7 context as the Step output.
- [x] 3.3 Update `packages/domains/ai-agents/manifests/ctx7/steps/inject-ctx7-context.yaml` to consume `CTX7_CONTEXT` from env instead of rereading `.$AGENT/ctx7-agents.md`.
- [x] 3.4 Wire ctx7 inject workflow entries to pass `CTX7_CONTEXT` via `$kfg.output(<step-reference-name>)`.

## 4. Tests and documentation

- [x] 4.1 Add or update unit tests for manifest validation covering required step-reference names, duplicate names, invalid `when.output.step`, and invalid `$kfg.output(...)` references.
- [x] 4.2 Add or update unit tests for resolver and generator behavior covering per-reference output identity, dynamic env expansion, and preserved outputs with step-reference env overrides.
- [x] 4.3 Add or update Bats or integration tests covering repeated use of the same Step resource in one workflow, output-based `when` resolution by step-reference name, and ctx7 output-to-env handoff.
- [x] 4.4 Update any relevant developer-facing documentation or comments that describe output addressing and workflow env behavior.
