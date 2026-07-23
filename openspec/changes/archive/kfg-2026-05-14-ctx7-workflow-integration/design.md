## Context

The kfg overlay dev workflow (`kfg.workflow.agents`) currently has 6 phases: Detection, Scaffold, Context, Commands, MCP, and Subagent. The nixai manifests define additional phases (gitignore, ctx7 install, ctx7 inject, per-agent cleanup) that are missing from the kfg repo. This design integrates those phases while preserving the existing workflow structure.

Current state:
- `.manifests/base/steps/` has 6 steps (aggregate-mcp, cleanup, convert, copy-context, detect-agent, materialize-scaffold)
- `.manifests/base/extensions/ctx7/` has `kfg.extension.ctx7.install` with hardcoded `.claude` defaults
- `.manifests/overlay/dev/agents-workflow.yaml` has 6 phases in `before`, 1 cleanup in `after`
- No `ensure-gitignore` or `inject-ctx7-context` steps exist in kfg repo

## Goals / Non-Goals

**Goals:**
- Add `kfg.ensure-gitignore` step for automatic `.gitignore` management
- Add `kfg.inject-ctx7-context` step for injecting ctx7 documentation into agent context files
- Update `kfg.extension.ctx7.install` to use generic defaults (empty AGENT_HOME/OUTPUT_DIR)
- Merge 3 new phases into existing workflow at correct weight positions
- Fix bug: stray quote in `agents.mcp.gemini` name field
- Per-agent cleanup in after phase (replace single blanket cleanup)

**Non-Goals:**
- Adding image-start/image-stop (store image materialization) — deferred
- Changing the manifest model or schema
- Modifying base steps that already work (detect-agent, copy-context, etc.)

## Decisions

### D1: Create steps in `.manifests/base/steps/` (not overlay)

Both `ensure-gitignore` and `inject-ctx7-context` are generic steps usable by any overlay, not just dev. Placing them in `base/steps/` follows the existing pattern where reusable steps live in the base layer.

Alternative: Define in overlay dev only. Rejected because these steps could be reused by other overlays (e.g., lifeos).

### D2: Use `$AGENT` env var (not `$NIXAI_AGENT`)

The kfg repo's `detect-agent` step reads `$AGENT` and outputs `AGENT`. The new `inject-ctx7-context` MUST follow this convention. The nixai version uses `$NIXAI_AGENT` which is a nixai-specific convention.

### D3: Weight-based ordering with gaps

Use weights with gaps (-90, -70, -65, -60, -55, -50, -45, -40, -35, 90) to allow future insertions between phases. This follows the existing pattern in the workflow.

### D4: ctx7.install before inject-ctx7-context

`ctx7.install` runs `ctx7 setup --cli --project` which generates skills and documentation. `inject-ctx7-context` reads the generated docs and injects them into AGENTS.md. The install MUST complete before inject.

### D5: Empty defaults for ctx7.install AGENT_HOME/OUTPUT_DIR

The current step defaults to `.claude` paths. Changing to empty forces explicit per-invocation values, preventing incorrect behavior when used with non-Claude agents. Each workflow invocation sets agent-specific paths via `env`.

### D6: Per-agent cleanup replaces blanket cleanup

The current `after` cleanup deletes all agent directories unconditionally. The nixai approach uses per-agent cleanup in `before` (weight 90) to remove other agents' artifacts. This is more targeted. The `after` cleanup remains as a final sweep.

## Risks / Trade-offs

- **[Risk] ctx7 CLI not available** → `ctx7.install` uses `failurePolicy: Ignore` so workflow continues if ctx7 is missing. Context injection gracefully skips if ctx7-agents.md doesn't exist.
- **[Risk] inject-ctx7-context idempotency** → Uses `<!-- context7 -->` markers for upsert semantics. Re-running replaces existing section.
- **[Trade-off] 4 invocations of ctx7.install** → Verbose but explicit. Each agent gets its own `when` condition and env vars. Alternative: parameterized step with AGENT_NAME. Rejected because the step schema doesn't support dynamic path derivation from a single env var.
- **[Trade-off] More workflow steps = slower** → Steps with `failurePolicy: Ignore` fail fast. Weight ordering ensures correct execution sequence.

## Migration Plan

1. Create new step files (ensure-gitignore, inject-ctx7-context)
2. Update kustomization.yaml to include new steps
3. Update ctx7.install defaults
4. Rewrite agents-workflow.yaml with merged phases
5. Run `make test-bats` to verify generated shell output
6. No rollback needed — backward compatibility is not a requirement
