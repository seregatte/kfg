## Context

Framework Steps are consumers of the engine runtime API. `kfg.materialize` is the main affected Step because it shells out to `kfg apply` repeatedly during per-item and aggregate conversion flows, which can currently emit child startup logs when parent verbosity is high.

## Goals / Non-Goals

**Goals:**
- Route nested framework `kfg` subprocesses through the engine's quiet internal execution wrapper.
- Preserve materialize behavior, outputs, and artifact registration while removing child startup noise.

**Non-Goals:**
- Change the materialize data contract or converter semantics.
- Silence third-party tools that are not `kfg` subprocesses.

## Decisions

### Decision: Framework manifests consume the runtime helper instead of embedding verbosity policy

Framework Steps will call the engine-provided internal execution helper rather than adding `KFG_VERBOSE=0` inline around each nested `kfg` call.

Alternatives considered:
- Inline `KFG_VERBOSE=0` in each call site: rejected because it is repetitive and couples framework manifests to engine subprocess policy.

## Risks / Trade-offs

- [If the runtime helper changes shape, framework Steps must stay aligned] -> Mitigation: capture the helper contract in specs and cover materialize behavior in framework tests.
