## Why

The dev overlay workflow (`agents-workflow.yaml`) is missing three functional phases that exist in the nixai manifests: automatic `.gitignore` management, ctx7 skill installation, and ctx7 context injection into agent files. Without these, developers must manually manage `.gitignore` entries, install ctx7 skills per-agent, and inject documentation context. This change brings kfg's overlay dev to parity with nixai's workflow while keeping the existing scaffold, settings, commands, MCP, and subagent phases.

## What Changes

- **New step `kfg.ensure-gitignore`**: Adds entries to `.gitignore` if not already present. Configurable via `GITIGNORE_ENTRIES` env var.
- **New step `kfg.inject-ctx7-context`**: Reads `.$AGENT/ctx7-agents.md`, extracts content between `<!-- context7 -->` markers, and injects it into the target file (AGENTS.md or CLAUDE.md). Uses `$AGENT` env var (not `$NIXAI_AGENT`).
- **Modified step `kfg.extension.ctx7.install`**: Change default `AGENT_HOME` and `OUTPUT_DIR` from hardcoded `.claude` paths to empty strings, requiring explicit values at invocation time.
- **Modified workflow `kfg.workflow.agents`**: Merge three new phases into the existing 6-phase workflow, ordered by weight:
  - Gitignore (-90) → Detection (-70) → Scaffold (-65) → Context (-60) → ctx7 install (-55) → ctx7 inject (-50) → Commands (-45) → MCP (-40) → Subagent (-35) → Cleanup (90)
- **Bug fix**: Fix stray quote in `agents.mcp.gemini` name field (line 172).

## Capabilities

### New Capabilities

- `ensure-gitignore`: Step that ensures specified entries exist in `.gitignore`, creating the file if needed. Idempotent — skips entries already present.
- `inject-ctx7-context`: Step that injects ctx7 documentation sections from per-agent ctx7 files into AGENTS.md or CLAUDE.md, using `<!-- context7 -->` markers for upsert semantics.

### Modified Capabilities

- `ctx7-install`: Default env vars (`AGENT_HOME`, `OUTPUT_DIR`) changed from hardcoded `.claude` paths to empty, requiring per-invocation values.
- `dev-workflow`: Overlay workflow expanded from 6 to 9 phases with new gitignore, ctx7 install, and ctx7 inject phases integrated at correct weight positions.

## Impact

- **Manifest files**: 2 new step YAML files in `.manifests/base/steps/`, 1 modified step in `.manifests/base/extensions/ctx7/steps/`, 1 modified kustomization, 1 rewritten workflow
- **Shell generation**: Generated shell functions will include gitignore management and ctx7 context injection steps
- **Agent UX**: Developers get automatic `.gitignore` setup, ctx7 skills installed per-agent, and ctx7 documentation injected into context files without manual intervention
- **Dependencies**: Requires `ctx7` CLI available in PATH for the install step
