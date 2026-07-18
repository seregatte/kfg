## Context

The `ai-agents` domain currently supports 4 agents: claude, gemini, opencode, and pi. The wizard overlay only supports pi and opencode. Gemini CLI requires `gemini-cli-bin` in the nix bundle but has no consumer in any overlay or project. Removing it is a cleanup that reduces the supported agent surface.

## Goals / Non-Goals

**Goals:**
- Completely remove all Gemini agent manifest resources, Cmd definitions, converters, and assets
- Remove `gemini-cli-bin` from the nix `kfg-bundle`
- Update all test files to not reference gemini
- Clean up documentation references

**Non-Goals:**
- Not modifying the agent detection step (`ai.steps.detect`) — it doesn't reference specific agents
- Not changing how other agents work
- Not touching the wizard overlay (already doesn't use gemini)

## Decisions

1. **Full removal (not deprecation)**: Since the wizard overlay already excludes gemini and no production project depends on it, we remove entirely rather than deprecate.
2. **Order after claude removal**: Both removals share a dependency chain (agents.yaml, agents/kustomization.yaml, flake.nix). Execute sequentially to avoid merge conflicts.
3. **Test updates inline**: Remove gemini test cases rather than replacing with another agent — reduces test maintenance.

## Risks / Trade-offs

- **[No migration path]** Users with custom overlays referencing `ai.gemini.cmd.main` or gemini assets will break → Mitigation: gemini was never used in any production project overlay. If needed, users can pin to a previous version.
