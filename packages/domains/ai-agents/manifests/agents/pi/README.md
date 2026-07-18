# Pi Agent

## Purpose

Pi is an alternative AI agent for kfg development. It provides settings materialization
and command conversion with a simpler converter surface than OpenCode.

## Prerequisites

- Pi CLI installed and available in PATH
- kfg CLI for manifest materialization

## Entrypoints

| Entrypoint | Type | Description |
|------------|------|-------------|
| `agents/pi/kustomization.yaml` | Selective | Pi agent assets and converters |

## Exported Resources

| Resource Name | Kind | Description |
|---------------|------|-------------|
| `ai.pi.asset.settings` | Assets | Pi agent configuration (theme, thinking level, compaction, retry) |
| `ai.pi.conv.command` | Converter | Convert YAML prompts to Pi command format |

## Agent Compatibility

| Agent | Supported | Notes |
|-------|-----------|-------|
| OpenCode | No | Not applicable |
| Pi | Yes | Primary target |

## Configuration

| Field | Default | Description |
|-------|---------|-------------|
| `defaultProvider` | `""` | Default AI provider |
| `defaultModel` | `""` | Default model identifier |
| `defaultThinkingLevel` | `medium` | Thinking depth (low, medium, high) |
| `theme` | `dark` | UI theme |
| `compaction.enabled` | `true` | Enable context compaction |
| `compaction.reserveTokens` | `16384` | Tokens reserved for compaction |
| `compaction.keepRecentTokens` | `20000` | Recent tokens to preserve |
| `retry.enabled` | `true` | Enable retry on failure |
| `retry.maxRetries` | `3` | Maximum retry attempts |
| `retry.baseDelayMs` | `2000` | Base retry delay in milliseconds |
| `retry.maxDelayMs` | `60000` | Maximum retry delay in milliseconds |
| `skills` | `["./skills"]` | Skill directories |
| `prompts` | `["./prompts"]` | Prompt directories |

## JSON Patch Examples

### Change thinking level

```yaml
# patch.yaml
- op: replace
  path: /spec/data/defaultThinkingLevel
  value: high
```

### Disable compaction

```yaml
# patch.yaml
- op: replace
  path: /spec/data/compaction/enabled
  value: false
```

### Add custom model

```yaml
# patch.yaml
- op: replace
  path: /spec/data/defaultModel
  value: claude-sonnet-4-20250514
```

## Workflow Usage

This capability is used in:
- `overlays/ai/` — Wizard materializes Pi settings
- `overlays/dev/` — Development workflow materializes Pi settings

## Generated Artifacts

| Artifact | Path | Description |
|----------|------|-------------|
| Settings | `.pi/settings.json` (or agent-specific path) | Pi configuration |

## Cache and Cleanup

- **Cache**: No — settings are regenerated on each run
- **Cleanup**: Generated settings files are removed during workflow cleanup phase

## Limitations

- Pi has only a command converter (`ai.pi.conv.command`); no skill, MCP, subagent, or config converters
- The command converter produces the same Markdown frontmatter format as OpenCode's command converter
- Pi does not support MCP server configuration through kfg manifests
- Pi does not support subagent definitions through kfg manifests
- Settings include `packages`, `extensions`, and `themes` arrays that are always empty in the default manifest; consumers must provide values via JSON Patch

## Compatibility Gaps

- **No MCP converter**: Pi cannot consume MCP server definitions. Consumers using Pi with MCP-dependent integrations (Chrome DevTools, Context7, Playwright) must configure MCP separately.
- **No subagent converter**: Pi cannot consume subagent definitions. The `review-minimal` subagent is OpenCode-only.
- **No config converter**: Pi settings are emitted as raw YAML/JSON; there is no `pi.conv.cfg` equivalent.

## Validation

```bash
# Build Pi agent entrypoint
nix develop .#dev --command kustomize build packages/domains/ai-agents/manifests/agents/pi/

# Verify settings resource exists
nix develop .#dev --command kustomize build packages/domains/ai-agents/manifests/agents/pi/ | grep "ai.pi.asset.settings"
```

## Canonical Specs

- [kfg-domain-package-contract](../../docs/context/openspec/specs/kfg-domain-package-contract/spec.md)
- [domain-ai-agents-capability-documentation](../../docs/context/openspec/specs/domain-ai-agents-capability-documentation/spec.md)
