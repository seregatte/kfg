# Chrome DevTools

## Purpose

Chrome DevTools MCP integration for browser automation and debugging.

## Prerequisites

- Node.js (with `npx` available)
- Chrome or Chromium installed on the system

## Entrypoints

| Entrypoint | Type | Description |
|------------|------|-------------|
| `kustomization.yaml` | Selective | Exports all chrome-devtools assets and steps |

## Exported Resources

| Resource Name | Kind | Description |
|---------------|------|-------------|
| `chrome.assets.mcp` | Assets | MCP server configuration for Chrome DevTools |
| `chrome-devtools.steps.install` | Step | Install Chrome DevTools skills via npx |

## Agent Compatibility

| Agent | Supported | Notes |
|-------|-----------|-------|
| OpenCode | Yes | MCP config materialized via `ai.opencode.conv.mcp` |
| Pi | Yes | MCP config supported |

## Configuration

The MCP asset defines the server configuration and available tools.

### Server Configuration

| Field | Default | Description |
|-------|---------|-------------|
| `spec.data.server.command` | `npx` | MCP server command |
| `spec.data.server.args` | `["-y", "chrome-devtools-mcp@latest"]` | Command arguments |
| `spec.data.enabled` | `true` | Enable/disable the integration |

### Available Tools

| Tool | Description |
|------|-------------|
| `navigate` | Navigate to a URL |
| `screenshot` | Take a screenshot of the current page |
| `evaluate` | Execute JavaScript in the browser |

## JSON Patch Examples

### Change MCP server args

```yaml
# patch.yaml
- op: replace
  path: /spec/data/server/args
  value:
    - -y
    - chrome-devtools-mcp@0.2.0
```

### Disable the integration

```yaml
# patch.yaml
- op: replace
  path: /spec/data/enabled
  value: false
```

## Workflow Usage

Used in `overlays/dev/agents-workflow.yaml` during MCP aggregation (Phase 10, weight -40). The `chrome.assets.mcp` asset is aggregated alongside `ctx7.assets.mcp` and `playwright.assets.mcp`, then converted by `ai.opencode.conv.mcp` into agent-specific MCP configuration.

## Generated Artifacts

MCP config files are materialized per-agent:

| Agent | Output Path |
|-------|-------------|
| OpenCode | `opencode.json` (mcp section) |

## Cache and Cleanup

| Aspect | Behavior |
|--------|----------|
| Install step cache | Enabled — `chrome-devtools.steps.install` caches to avoid repeated npx downloads |
| Cleanup | `kfg.cleanup` runs after workflow completion, removing temporary files |

## Limitations

- Requires Chrome or Chromium installed on the host system
- MCP tools are browser-specific — only browser automation and debugging operations are available
- The `npx` command requires network access to download the MCP server package on first run

## Validation

```bash
nix develop .#dev --command kustomize build packages/domains/ai-agents/manifests/chrome-devtools
```

## Canonical Specs

- [kfg-domain-package-contract](../../../../../../docs/context/openspec/specs/kfg-domain-package-contract/spec.md)
