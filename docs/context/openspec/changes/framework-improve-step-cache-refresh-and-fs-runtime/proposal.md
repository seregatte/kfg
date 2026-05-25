## Why

Framework Steps such as `kfg.materialize` shell out to `kfg` internally, and those child invocations can re-emit startup logs when parent verbosity is high. That noise makes parent Step output harder to read and undermines the value of the structured runtime logs the framework already depends on.

## What Changes

- Adopt the runtime's quiet internal `kfg` execution wrapper for framework Steps that invoke nested engine commands.
- Update `kfg.materialize` to run its nested conversion calls through that wrapper without changing its conversion behavior or artifact registration contract.
- Keep framework Steps compatible with declarative and runtime-discovered artifacts while preserving existing output behavior.
- Add or update framework tests and documentation for quiet nested `kfg` execution.

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- `reusable-framework-steps`: framework Steps use the runtime's internal execution wrapper when they invoke nested `kfg` subprocesses.
- `materialize-step`: materialize performs nested conversions without leaking child startup logs into parent stderr output.

## Impact

- Affects manifests under `packages/framework/manifests/steps/`, especially `materialize.yaml`.
- Depends on the engine runtime exposing the internal execution helper.
- Requires framework-level validation that materialize output stays unchanged while nested logs are suppressed.
