## Why

The current workflow model addresses step outputs by `Step.metadata.name`, which makes outputs and `when.output` ambiguous when the same Step is referenced multiple times in a single workflow. This blocks passing output-derived values through workflow `env` and forces file-based coupling such as the current ctx7 install and inject flow.

## What Changes

- Require every workflow `StepReference` to have a stable `name` and use that name as the execution identity for outputs and output-based conditions.
- Change `when.output.step` to reference `StepReference.name` instead of the underlying Step resource name.
- Extend workflow `env` values to support reading a referenced step output with `$kfg.output(<step-reference-name>)`.
- Update generated shell execution so step reference `env` remains isolated without using a subshell that discards step outputs.
- Update ctx7 install and inject behavior so install emits the extracted ctx7 context as step output and inject consumes it from workflow `env` instead of rereading an agent-local file.
- **BREAKING** Existing workflows that use `when.output.step` with `Step.metadata.name` must migrate to `StepReference.name`.

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- `manifest-model`: workflow step references gain required names and output conditions resolve by step-reference identity instead of Step resource identity.
- `manifest-env`: workflow step-reference env values can consume prior step outputs through `$kfg.output(<step-reference-name>)`, and step-reference env isolation must preserve produced outputs.
- `inject-ctx7-context`: ctx7 context injection consumes supplied env content instead of rereading `.$AGENT/ctx7-agents.md` directly.

## Impact

- Affects manifest validation and resolution in `src/internal/manifest` and `src/internal/resolve`.
- Affects shell generation and runtime helpers in `src/internal/generate`.
- Requires workflow manifest updates under `packages/domains/ai-agents/overlays/dev/`.
- Requires ctx7 step updates under `packages/domains/ai-agents/manifests/ctx7/steps/`.
- Requires unit and integration coverage for output addressing, dynamic env expansion, and ctx7 workflow behavior.
