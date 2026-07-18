## Why

Claude Code agent support (`claude-code` nix package) adds maintenance overhead with subagent converters, permission settings, and 3 converters (command, mcp, subagent) — more than any other agent. The wizard overlay already only supports `pi` and `opencode`. Removing claude simplifies the manifest catalog, reduces test surface, and aligns with the project's direction toward opencode/pi as the primary agents.

## What Changes

- **Remove**: Claude agent directory `manifests/agents/claude/` (assets, converters, kustomization)
- **Remove**: Claude Cmd definition `ai.claude.cmd.main` from `manifests/cmds/agents.yaml`
- **Remove**: Claude branch from `ctx7.steps.install` agent detection
- **Remove**: `pkgs.claude-code` from `flake.nix` kfg-bundle
- **Remove**: `/.claude/` and `/CLAUDE.md` entries from `.gitignore`
- **Update**: Agent registry in `manifests/agents/kustomization.yaml`
- **Update**: Documentation in `manifests/README.md` and `manifests/agents/README.md`
- **Update**: All test files referencing claude (Bats + Go)
- **Update**: Comment-level references in Go source files
- **Update**: Wizard skill prompt to not mention claude in agent catalog

## Capabilities

### Removed Capabilities

- `domain-ai-agents-claude-agent`: The entire Claude agent support — no Cmd, no converters (command, mcp, subagent), no settings, no nix package

## Impact

- **Domain `ai-agents`**: Fewer agent directories; no subagent converter; cleaner agent registry
- **Engine**: One less nix dependency in `flake.nix`; fewer test fixtures
- **Framework**: Framework tests lose claude scaffolding cases
- **Backward compatibility**: Breaking — any `kustomization.yaml` referencing claude will fail. Users must migrate to opencode or pi.
