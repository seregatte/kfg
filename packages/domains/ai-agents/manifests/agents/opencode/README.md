# OpenCode Agent

## Purpose

OpenCode is the primary AI agent for kfg development. It provides settings materialization,
command and skill conversion, MCP server configuration, subagent definitions, and config
file generation.

## Prerequisites

- OpenCode CLI installed and available in PATH
- kfg CLI for manifest materialization

## Entrypoints

| Entrypoint | Type | Description |
|------------|------|-------------|
| `agents/opencode/kustomization.yaml` | Selective | OpenCode agent assets and converters |

## Exported Resources

| Resource Name | Kind | Description |
|---------------|------|-------------|
| `ai.opencode.asset.settings` | Assets | OpenCode agent configuration (instructions, models) |
| `ai.opencode.conv.command` | Converter | Convert YAML prompts to OpenCode command format |
| `ai.opencode.conv.skill` | Converter | Convert YAML prompts to OpenCode skill format |
| `ai.opencode.conv.mcp` | Converter | Convert MCP server configs to OpenCode JSON format |
| `ai.opencode.conv.subagent` | Converter | Convert subagent definitions to OpenCode format |
| `ai.opencode.conv.cfg` | Converter | Convert settings to OpenCode config JSON |

## Agent Compatibility

| Agent | Supported | Notes |
|-------|-----------|-------|
| OpenCode | Yes | Primary target |
| Pi | No | Not applicable |

## Configuration

| Field | Default | Description |
|-------|---------|-------------|
| `instructions` | `["AGENTS.md"]` | List of instruction files |
| `mode.plan.model` | `opencode/mimo-v2.5-free` | Model for planning tasks |
| `mode.build.model` | `opencode/deepseek-v4-flash-free` | Model for build tasks |

## JSON Patch Examples

### Change plan model

```yaml
# patch.yaml
- op: replace
  path: /spec/data/mode/plan/model
  value: opencode/claude-sonnet-4-20250514
```

### Add instruction file

```yaml
# patch.yaml
- op: add
  path: /spec/data/instructions/-
  value: CONTRIBUTING.md
```

## Workflow Usage

This capability is used in:
- `overlays/ai/` — Wizard materializes OpenCode settings and commands
- `overlays/dev/` — Development workflow materializes settings, AGENTS.md, commands, and MCP

## Generated Artifacts

| Artifact | Path | Description |
|----------|------|-------------|
| Settings | `.opencode/config.json` or agent-specific path | OpenCode configuration |
| Commands | `.opencode/commands/*.md` | Materialized command files |
| Skills | `.opencode/skills/*/SKILL.md` | Materialized skill files |
| MCP config | `.opencode/mcp.json` | MCP server configuration |
| Subagents | `.opencode/agents/*.md` | Materialized subagent files |

## Cache and Cleanup

- **Cache**: No — settings are regenerated on each run
- **Cleanup**: Generated files in `.opencode/` are removed during workflow cleanup phase

## Limitations

- Pi converter (`ai.pi.conv.command`) is separate; OpenCode converters are agent-specific
- MCP converter emits only `type`, `command`, `environment`, and `enabled` fields; does not support `tools`, `permissions`, or other MCP metadata
- Subagent converter produces Markdown frontmatter only; does not validate model availability
- Config converter emits `$schema`, `instructions`, `model`, and `mode` only; does not include all OpenCode config options

## Validation

```bash
# Build OpenCode agent entrypoint
nix develop .#dev --command kustomize build packages/domains/ai-agents/manifests/agents/opencode/

# Verify settings resource exists
nix develop .#dev --command kustomize build packages/domains/ai-agents/manifests/agents/opencode/ | grep "ai.opencode.asset.settings"
```

## Canonical Specs

- [kfg-domain-package-contract](../../docs/context/openspec/specs/kfg-domain-package-contract/spec.md)
- [domain-ai-agents-capability-documentation](../../docs/context/openspec/specs/domain-ai-agents-capability-documentation/spec.md)
