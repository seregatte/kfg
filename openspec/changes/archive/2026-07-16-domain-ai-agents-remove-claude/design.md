## Context

The `ai-agents` domain currently supports 4 agents: claude, gemini, opencode, and pi. Claude requires the most converters (command, mcp, subagent) and a permissions settings asset. The wizard overlay only supports pi and opencode. Removing claude is cleanup that reduces the supported agent surface to only the two agents the project actively maintains.

## Goals / Non-Goals

**Goals:**
- Completely remove all Claude agent manifest resources, Cmd definitions, converters, and assets
- Remove `claude-code` from the nix `kfg-bundle`
- Update all test files to not reference claude
- Clean up documentation and wizard prompt references

**Non-Goals:**
- Not modifying the agent detection step (`ai.steps.detect`)
- Not changing how opencode/pi work
- Not touching the wizard overlay structure (only prompt content)

## Decisions

1. **Full removal (not deprecation)**: Same rationale as gemini — no production project depends on claude via kfg.
2. **Sequential execution**: Execute after gemini removal to avoid merge conflicts in shared files (agents.yaml, agents/kustomization.yaml, flake.nix).
3. **Wizard prompt update**: Remove `claude` from the wizard skill's catalog list to avoid suggesting unsupported agents.

## Risks / Trade-offs

- **[Wider breakage]** Claude has more converters than gemini (subagent converter) so removal touches more files → Mitigation: comprehensive test coverage validates nothing is orphaned.
- **[Wizard staleness]** Wizard prompt may still reference claude in example conversations → Mitigation: included in task list.
