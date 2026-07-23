# Playwright

## Purpose

Playwright E2E testing integration with MCP server for browser automation.

## Prerequisites

- Node.js and `npx` available on `$PATH`
- `@playwright/cli` (installed via `npx --yes @playwright/cli`)
- Playwright browsers installed (`npx playwright install`)

## Selective Entrypoints

| Entrypoint | Description |
|------------|-------------|
| `kustomization.yaml` | Playwright package (selective) |

## Aggregate Entrypoints

| Entrypoint | Description |
|------------|-------------|
| `assets/kustomization.yaml` | All Playwright assets |
| `steps/kustomization.yaml` | All Playwright steps |

## Exported Resources

| Resource Name | Kind | Description |
|---------------|------|-------------|
| `playwright.assets.mcp` | Assets | MCP server configuration for Playwright browser automation |
| `playwright.steps.install` | Step | Installs Playwright CLI skills into the output directory |

## Agent Compatibility

| Agent | Supported | Notes |
|-------|-----------|-------|
| opencode | Yes | MCP config supported via converter |
| pi | Yes | MCP config supported |

## Configuration

### MCP Server (`playwright.assets.mcp`)

| Field | Default | Description |
|-------|---------|-------------|
| `server.command` | `npx` | MCP server command |
| `server.args` | `["-y", "@playwright/mcp@latest", "--extension"]` | MCP server arguments |
| `server.env` | `{}` | Additional environment variables |
| `enabled` | `true` | Whether the MCP server is active |

### Step Inputs (`playwright.steps.install`)

| Field | Default | Description |
|-------|---------|-------------|
| `OUTPUT_DIR` | `.opencode/skills/` | Output directory for generated skill files |

## JSON Patch Examples

### Change MCP server args

```yaml
# kustomization.yaml
patches:
  - target:
      kind: Assets
      name: playwright.assets.mcp
    patch: |
      - op: replace
        path: /spec/data/server/args
        value:
          - -y
          - "@playwright/mcp@latest"
```

### Disable MCP

```yaml
# kustomization.yaml
patches:
  - target:
      kind: Assets
      name: playwright.assets.mcp
    patch: |
      - op: replace
        path: /spec/data/enabled
        value: false
```

## Workflow Usage

Available but not wired into either the `dev` or `ai` overlay by default.
To include it, add the step to a workflow YAML:

```yaml
steps:
  - name: playwright.steps.install
    weight: 0
```

## Generated Artifacts

- **MCP config**: Playwright MCP server configuration for agent integration
- **Skills**: Playwright CLI skill files copied to `$OUTPUT_DIR` (typically `.opencode/skills/playwright-cli/`)

## Cache and Cleanup

| Behavior | Detail |
|----------|--------|
| Cache | Enabled (`cache.enabled: true`) |
| Cleanup | Removes temporary `.opencode` directory after skill copy |

## Limitations

- Requires Playwright browsers to be installed separately (`npx playwright install`).
- MCP tools are browser-specific; not all Playwright features are exposed via MCP.
- The step performs no validation of Playwright version or browser availability.

## Validation

```bash
kustomize build packages/domains/ai-agents/manifests/playwright/
```

## Canonical Specs

- [kfg-domain-package-contract](../../../../docs/context/openspec/specs/kfg-domain-package-contract/spec.md)
