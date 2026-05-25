## Why

The repository currently models manifest-driven materialization through three overlapping steps: `kfg.agents.steps.settings`, `kfg.convert`, and `kfg.aggregate-mcp`. Each step performs the same core operation of converting one or more `Assets` with a `Converter` and writing generated artifacts, but each exposes a different contract for inputs, outputs, validation, and post-processing.

This overlap makes workflow authoring harder to reason about and encourages more special-purpose steps whenever a new materialization pattern appears. The project has already decided that backward compatibility is not required here, so this is a good point to replace the fragmented contracts with a single stable primitive.

## What Changes

- Introduce a new shared step `kfg.materialize` as the single manifest-level primitive for asset conversion and artifact materialization.
- Support two explicit modes in `kfg.materialize`: `per-item` for positional `ASSETS[i] -> OUTPUTS[i]` materialization and `aggregate` for `ASSETS[*] -> merged output` materialization.
- Use a single output field, `OUTPUTS`, for both modes, with one output path required in aggregate mode.
- Make validation strict: required fields must be present, `per-item` counts must match, and unsupported mode-specific combinations must fail.
- Migrate workflow usage away from `kfg.agents.steps.settings`, `kfg.convert`, and `kfg.aggregate-mcp` to `kfg.materialize`.
- Add unit, integration, and Bats coverage for the new shared contract and update OpenSpec specs/documentation to describe the new durable behavior.

## Non-Goals

- Preserve the old step interfaces or provide compatibility shims for `ASSET`, `OUTPUT`, or `TARGET`.
- Expose arbitrary reducer expressions or a general-purpose transformation DSL in the step contract.
- Change the `kfg apply --convert/--use` CLI conversion contract itself.

## Capabilities

### New Capabilities

- `materialize-step`: define a single reusable manifest step contract for per-item and aggregate asset materialization.

### Modified Capabilities

- `dev-workflow`: use `kfg.materialize` as the shared primitive for settings, command, subagent, and MCP materialization phases.
- `manifest-model`: document the durable manifest-step contract that treats materialization as a first-class reusable step pattern.

## Impact

- Affected manifests: `.manifests/base/steps/convert.yaml`, `.manifests/base/steps/aggregate-mcp.yaml`, `.manifests/base/agents/steps/settings.yaml`, and `.manifests/overlay/dev/agents-workflow.yaml`
- Affected tests: Bats tests for base steps and workflow behavior, plus Go tests where step contracts are documented or validated indirectly
- Affected specs: new `materialize-step` capability plus `dev-workflow` and `manifest-model`
- User-facing shell UX: workflow authors get one shared step contract for all converter-driven materialization cases
