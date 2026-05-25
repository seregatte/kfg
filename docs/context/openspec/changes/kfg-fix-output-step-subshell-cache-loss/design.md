## Context

Engine-generated wrappers for `spec.output` Steps currently use command substitution to capture stdout, which runs the Step body in a subshell. That execution model is incompatible with runtime helpers that mutate shell state, especially artifact registration used by cache persistence.

## Goals / Non-Goals

**Goals:**
- Execute output-producing Steps in the parent shell.
- Capture stdout without losing runtime side effects.
- Add regression coverage for cacheable output-producing Steps.

**Non-Goals:**
- Change manifest syntax or broaden runtime refactors beyond this execution bug.

## Decisions

### Decision: Use parent-shell stdout redirection for `spec.output`

The wrapper will redirect stdout to a temporary file, run the Step in the parent shell, then read the file back into the output store.

Alternatives considered:
- Subshell plus state synchronization: rejected as more fragile than removing the subshell.

## Risks / Trade-offs

- [Temporary output files need reliable cleanup] -> Mitigation: add wrapper cleanup on both success and failure paths.
