## Context

The generated Step wrapper has two execution paths: one for plain Steps and one for Steps that declare `spec.output`. The output path currently captures stdout with `__output="$( {{ .RunScript }} )"`, which means the entire Step body runs inside a subshell. That makes stdout capture easy, but it also discards mutations to shell state performed by the Step body. In practice, any `__kfg_add_artifact` call made by a Step with `spec.output` updates the subshell copy of `KFG_ARTIFACTS`, not the parent shell array used by cache delta computation and later cleanup. The same structural problem risks corrupting or dropping other runtime-side effects that are expected to survive Step execution.

This is an engine-layer bug rather than a package-specific bug. Domain Steps like `ctx7.steps.install` only expose it because they both emit output and register artifacts. The correct fix belongs in the generated runtime execution model.

## Goals / Non-Goals

**Goals:**
- Run output-producing Steps in the parent shell so runtime side effects persist.
- Continue capturing the Step's stdout as the declared output value.
- Preserve compatibility for Steps without `spec.output`.
- Make cache persistence for output-producing Steps include both artifacts and output.

**Non-Goals:**
- Redesign the general Step execution contract beyond what is required to eliminate the subshell side-effect loss.
- Change the public manifest model for `spec.output`.
- Introduce a native Go executor for Steps.

## Decisions

### Decision: Capture output through a temporary file instead of command substitution

For Steps with `spec.output`, the wrapper will execute the Step body in the parent shell with stdout redirected to a temporary file. After successful execution, the wrapper will read the file content and call `__kfg_output_set` with the captured value.

Why this approach:
- It preserves shell-side effects because the Step body no longer runs in a subshell.
- It keeps the output capture contract straightforward and shell-native.
- It avoids fragile attempts to synchronize artifacts or shell arrays back from a subshell.

Alternatives considered:
- Keep command substitution and export side effects separately: rejected because it duplicates state synchronization logic and remains brittle.
- Pipe the Step body into `tee` or other process substitution forms: rejected because many of those still introduce subshell or pipeline side-effect traps.

### Decision: Treat output and artifact persistence as one invocation contract

Cache store for output-producing Steps will continue to consume the same runtime artifact delta and the output value captured after parent-shell execution.

Why this approach:
- It preserves the current cache model while fixing the underlying execution bug.
- It ensures output-producing Steps behave the same as non-output Steps with respect to artifact tracking.

Alternatives considered:
- Special-case cache logic for `spec.output`: rejected because the real bug is execution isolation, not cache schema.

## Risks / Trade-offs

- [Temporary-file output capture could change whitespace or newline handling if implemented carelessly] -> Mitigation: read the file content verbatim and cover multiline outputs in tests.
- [A failed Step could leave temporary capture files behind] -> Mitigation: use trap-based cleanup or explicit cleanup in the wrapper after each execution path.
- [Some Steps may rely on stderr/stdout interleaving behavior] -> Mitigation: preserve stderr passthrough and only redirect stdout for output capture.

## Migration Plan

1. Update generated output-producing Step wrappers to capture stdout via a temporary file in the parent shell.
2. Add regression tests for output-producing Steps that register artifacts and use cache store/restore.
3. Validate the ctx7 install path as the representative real-world case.

## Open Questions

- None. The failure mode and the preferred parent-shell capture strategy are clear.
