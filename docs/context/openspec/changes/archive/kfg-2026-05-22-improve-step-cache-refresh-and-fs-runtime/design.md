## Context

The engine generates shell runtime helpers that already own cache identity, artifact registration, and restore/store behavior, so the most durable place to fix refresh and dynamic filesystem tracking is the engine layer. Two current gaps drive this change: refresh skips cache storage entirely, and there is no Go-owned internal API for filesystem snapshot/diff that can survive a future Windows backend.

## Goals / Non-Goals

**Goals:**
- Make refresh rebuild and replace cached Step entries.
- Add a portable internal `kfg sys fs` CLI for runtime filesystem snapshot/diff.
- Provide quiet nested internal `kfg` execution helpers for generated runtime code.

**Non-Goals:**
- Replace generated shell execution with a native Go Step runner.
- Change cache identity rules or public `sys gc` semantics.

## Decisions

### Decision: `kfg sys fs` is an internal CLI backend, not a shell implementation detail

The engine will own snapshot/diff logic in Go and surface it through `kfg sys fs` so the generated shell only orchestrates calls and consumes text output.

Alternatives considered:
- Bash-only helpers: rejected because they entrench platform-specific behavior.

### Decision: Refresh rewrites cache entries after successful execution

The generated Step wrapper will skip restore on refresh but still store the new result. Cache store logic will clear any existing entry before writing the replacement.

Alternatives considered:
- Refresh bypass-only: rejected because it conflicts with the cache specification.

### Decision: Runtime wrappers use a quiet internal `kfg` subprocess helper

Nested engine subprocesses will execute with child-scoped `KFG_VERBOSE=0` so human startup logs from the child do not pollute the parent Step output.

Alternatives considered:
- Let each Step manage child verbosity manually: rejected because it duplicates policy across manifests and templates.

## Risks / Trade-offs

- [Internal CLI calls add subprocess overhead] -> Mitigation: keep command output simple and use the API only where runtime portability matters.
- [Clearing cache entries before rewrite can remove the old entry if the new write fails] -> Mitigation: limit the change to successful Step executions and cover rewrite behavior in tests.
