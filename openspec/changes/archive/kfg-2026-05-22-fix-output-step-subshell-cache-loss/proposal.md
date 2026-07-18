## Why

Steps with `spec.output` currently execute their full `run` script inside command substitution so the wrapper can capture stdout. That execution model drops shell side effects such as `__kfg_add_artifact` registrations and can also lose or corrupt cached output semantics, which makes cacheable output-producing Steps unreliable.

## What Changes

- Replace command-substitution execution for output-producing Steps with an output-capture mechanism that runs the Step in the parent shell.
- Preserve runtime side effects from output-producing Steps, including `__kfg_add_artifact`, `__kfg_output_set`, and any other shell state updates performed by the Step.
- Ensure cache store logic persists both dynamically registered artifacts and captured output values for Steps with `spec.output`.
- Add regression coverage for output-producing Steps that also register artifacts and use cache restore.

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- `shell-runtime-api`: output-producing Steps execute without a subshell that discards runtime side effects, while still exposing captured output to the generated shell context.
- `step-cache`: cacheable Steps with `spec.output` persist both captured output values and runtime artifact registrations from the same invocation.

## Impact

- Affects generated Step wrappers in `src/internal/generate/templates/bash_step.tmpl`.
- Affects cache store/restore behavior in `src/internal/generate/templates/bash_helper.tmpl`.
- Requires unit, golden, and Bats regression coverage for output-producing cacheable Steps.
