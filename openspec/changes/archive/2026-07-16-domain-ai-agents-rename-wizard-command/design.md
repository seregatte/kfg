## Context

The wizard skill prompt is materialized into agent-specific locations via `kfg.materialize` with the agent's command converter. Currently it produces:
- `.opencode/commands/kfg-wizard.md` (opencode)
- `.pi/prompts/kfg-wizard.md` (pi)

The user invokes it with `/kfg-wizard`. The `kfg-` prefix is redundant since the command only exists within a kfg-managed agent environment. Shortening to `/wizard` is more ergonomic.

## Goals / Non-Goals

**Goals:**
- Rename the wizard command from `/kfg-wizard` to `/wizard`
- Rename the asset file from `kfg-wizard.yaml` to `wizard.yaml`
- Update all OUTPUTS paths in `ai-workflow.yaml`
- Update all documentation and spec references
- Clean up stale agent references (claude, gemini) in the wizard prompt

**Non-Goals:**
- Not changing the `kfg ai` CLI command name
- Not changing the wizard prompt content (beyond agent catalog cleanup)
- Not changing the materialization mechanism
- Not adding new agents to the wizard catalog

## Decisions

1. **File rename (not copy)**: `git mv` the file to preserve history.
2. **Concurrent agent catalog cleanup**: Since claude and gemini are being removed in parallel features, clean up their references in the wizard prompt as part of this change.
3. **No backward compat shim**: Old `/kfg-wizard` commands will simply not work — users will see the new `/wizard` command after re-applying.

## Risks / Trade-offs

- **[User muscle memory]** Users accustomed to `/kfg-wizard` will need to learn `/wizard` → Low risk: the wizard is invoked infrequently and the `kfg ai` CLI command is the primary entrypoint.
