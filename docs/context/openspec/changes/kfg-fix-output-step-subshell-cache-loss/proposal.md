## Why

The generated engine runtime currently captures `spec.output` by running the entire Step body inside command substitution. That breaks a core runtime guarantee: output-producing Steps can no longer persist shell-side effects such as artifact registration into the parent execution context, which in turn breaks cache persistence and restore for real install Steps like ctx7.

## What Changes

- Change the generated runtime so Steps with `spec.output` execute in the parent shell while still capturing stdout as the output value.
- Preserve artifact registration and other runtime side effects from output-producing Steps.
- Ensure cache store/restore for output-producing Steps persists both output metadata and dynamically registered artifacts.
- Add regression tests for the non-subshell execution path.

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- `shell-runtime-api`: output-producing Steps must not require a subshell that discards runtime side effects.

## Impact

- Affects generated runtime templates and engine-side tests.
- May affect any package Step that combines `spec.output` with artifact registration.
- Requires Bats coverage for cache behavior with output-producing Steps.
