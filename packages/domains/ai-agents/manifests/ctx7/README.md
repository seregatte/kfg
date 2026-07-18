# Context7

## Purpose

Context7 provides library documentation integration for AI agents via MCP (Model Context Protocol). It fetches up-to-date library documentation and injects context into agent instruction files, ensuring agents always have current API references.

## Prerequisites

- Node.js and npx available in PATH
- `ctx7` CLI (installed via `npx @upstash/context7-mcp`)
- Internet connection for library documentation fetching

## Entrypoints

| Entrypoint | Type | Description |
|------------|------|-------------|
| `ctx7/kustomization.yaml` | Selective | Context7 assets and steps |

## Exported Resources

| Resource Name | Kind | Description |
|---------------|------|-------------|
| `ctx7.assets.mcp` | Assets | Context7 MCP server configuration |
| `ctx7.steps.install` | Step | Installs ctx7 skills and extracts context |
| `ctx7.steps.inject` | Step | Injects ctx7 context into target file |

## Agent Compatibility

| Agent | Supported | Notes |
|-------|-----------|-------|
| OpenCode | Yes | Skills installed to `.opencode/skills/`, MCP config aggregated into `opencode.json` |
| Pi | Yes | MCP config supported; inject step targets `AGENTS.md` |

## Configuration

| Field | Default | Description |
|-------|---------|-------------|
| `server.command` | `npx` | MCP server command |
| `server.args` | `["-y", "@upstash/context7-mcp"]` | MCP server arguments |
| `tools[0].name` | `get-library-docs` | Fetch documentation for a specific library |
| `tools[1].name` | `search-library-docs` | Search documentation across libraries |
| `install.env.TARGET_FILE` | `AGENTS.md` | Target file for context injection |

## JSON Patch Examples

### Change inject target file

```yaml
# patch.yaml
- op: replace
  path: /spec/env/TARGET_FILE
  value: OPENCODE_INSTRUCTIONS.md
```

### Disable MCP server

```yaml
# patch.yaml
- op: replace
  path: /spec/data/enabled
  value: false
```

## Workflow Usage

This capability is used in `overlays/dev/agents-workflow.yaml`:

| Step | Weight | Description |
|------|--------|-------------|
| `agents.ctx7.install.opencode` | -55 | Installs ctx7 skills for OpenCode agent |
| `agents.ctx7.install.pi` | -55 | Installs ctx7 skills for Pi agent |

The install step runs before MCP materialization (weight -40) and commands (weight -45). The inject step (currently commented out in the workflow) would run after install to inject context into the target file.

## Generated Artifacts

| Artifact | Path | Description |
|----------|------|-------------|
| Skills | `.opencode/skills/*/SKILL.md` or `.pi/skills/*/SKILL.md` | Installed ctx7 skill files |
| Context section | `AGENTS.md` (or configured target) | Injected `<!-- context7 -->` section with library docs |

## Cache and Cleanup

- **Cache**: Enabled — install step caches results to avoid redundant `ctx7 setup` runs
- **Cleanup**: The install step removes the temporary `.agents` directory created during setup (OpenCode agent only)

## Limitations

- Requires internet connection for library documentation fetching
- The inject step modifies the target file in-place; existing `<!-- context7 -->` sections are replaced
- The inject step is currently commented out in the default workflow

## Validation

```bash
# Build Context7 entrypoint
nix develop .#dev --command kustomize build packages/domains/ai-agents/manifests/ctx7/

# Verify resources exist
nix develop .#dev --command kustomize build packages/domains/ai-agents/manifests/ctx7/ | grep "ctx7"
```

## Canonical Specs

- [kfg-domain-package-contract](../../../docs/context/openspec/specs/kfg-domain-package-contract/spec.md)
