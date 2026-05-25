## Context

The AI-agent domain's biggest artifact-verbosity issue is the ctx7 install flow in the development overlay, where each agent-specific reference repeats the final `find-docs` path even though the Step already knows the output directory it writes into.

## Goals / Non-Goals

**Goals:**
- Let ctx7 install discover its produced skill artifacts dynamically from `OUTPUT_DIR`.
- Remove redundant workflow-level ctx7 artifact declarations after the Step owns registration.

**Non-Goals:**
- Rewrite all install Steps in this change.
- Change ctx7 install's user-facing output or install command contract.

## Decisions

### Decision: Observe `OUTPUT_DIR` at depth 1 before and after ctx7 install

The ctx7 Step will snapshot `OUTPUT_DIR` with `--maxdepth 1`, diff the before/after results, and register the newly created children as artifacts.

Alternatives considered:
- Register `OUTPUT_DIR` wholesale: rejected because it is too broad and can capture unrelated sibling outputs.
- Keep workflow-level artifact declarations: rejected because they are the verbosity problem this change is meant to remove.

## Risks / Trade-offs

- [ctx7 could change its output layout in the future] -> Mitigation: keep the engine API depth-configurable and cover the current output shape in domain tests.
