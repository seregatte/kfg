# AI Agents Domain Inventory

## Agents

| Agent | Resource Name | Entrypoint | Supported Converters | Settings Target |
|-------|---------------|------------|---------------------|-----------------|
| OpenCode | `ai.opencode.cmd.main` | `agents/opencode/` | command, skill, mcp, subagent, cfg | `AGENTS.md` |
| Pi | `ai.pi.cmd.main` | `agents/pi/` | command | `./skills`, `./prompts` |

## Commands

| Resource Name | Entrypoint | Description |
|---------------|------------|-------------|
| `ai.opencode.cmd.main` | `agents/opencode/` | Runs OpenCode binary |
| `ai.pi.cmd.main` | `agents/pi/` | Runs Pi binary |
| `ai.cmds.openspec` | `cmds/openspec.yaml` | Wraps OpenSpec CLI |

## Steps

| Resource Name | Entrypoint | Supported Agents | Cache | Cleanup |
|---------------|------------|-----------------|-------|---------|
| `ai.steps.detect` | `steps/detect.yaml` | All | No | No |
| `ai.steps.deprecation-warn` | `steps/deprecation-warn.yaml` | All | No | No |
| `ctx7.steps.install` | `ctx7/steps/install.yaml` | opencode, pi | Yes | Yes |
| `ctx7.steps.inject` | `ctx7/steps/inject-ctx7-context.yaml` | opencode, pi | No | No |
| `chrome-devtools.steps.install` | `chrome-devtools/steps/install.yaml` | opencode, pi | Yes | Yes |
| `playwright.steps.install` | `playwright/steps/install.yaml` | opencode, pi | Yes | Yes |
| `notebooklm.steps.install` | `notebooklm/steps/install.yaml` | opencode, pi | Yes | Yes |
| `gws.steps.install` | `gws/steps/install.yaml` | opencode, pi | Yes | Yes |
| `openspec.steps.install` | `openspec/steps/install.yaml` | opencode, pi | Yes | Yes |

## Prompts

| Resource Name | Entrypoint | Description |
|---------------|------------|-------------|
| `ai.prompts.git-commit` | `prompts/git-commit.yaml` | Generate conventional commit |
| `ai.prompts.review-code` | `prompts/review-code.yaml` | Review code quality |
| `ai.prompts.review-search` | `prompts/review-search.yaml` | Search and review code |
| `ai.prompts.refactor-pure` | `prompts/refactor-pure.yaml` | Pure functional refactor |

## Subagents

| Resource Name | Entrypoint | Model |
|---------------|------------|-------|
| `ai.subagents.review-minimal` | `subagents/review-minimal.yaml` | claude-sonnet-4-20250514 |

## Converters

| Resource Name | Entrypoint | Input | Output Format |
|---------------|------------|-------|---------------|
| `ai.conv.to-json` | `converters/to-json/to-json.yaml` | YAML | JSON |
| `ai.opencode.conv.command` | `agents/opencode/converters/command.yaml` | YAML prompts | Markdown frontmatter |
| `ai.opencode.conv.skill` | `agents/opencode/converters/skill.yaml` | YAML prompts | Markdown frontmatter |
| `ai.opencode.conv.mcp` | `agents/opencode/converters/mcp.yaml` | YAML MCP | OpenCode JSON |
| `ai.opencode.conv.subagent` | `agents/opencode/converters/subagent.yaml` | YAML subagents | Markdown frontmatter |
| `ai.opencode.conv.cfg` | `agents/opencode/converters/cfg.yaml` | YAML settings | OpenCode config JSON |
| `ai.pi.conv.command` | `agents/pi/converters/command.yaml` | YAML prompts | Markdown frontmatter |

## Integrations

| Name | Has MCP Assets | Has Steps | Resource Names |
|------|---------------|-----------|----------------|
| Chrome DevTools | Yes | Yes | `chrome.assets.mcp`, `chrome-devtools.steps.install` |
| Context7 | Yes | Yes | `ctx7.assets.mcp`, `ctx7.steps.install`, `ctx7.steps.inject` |
| Playwright | Yes | Yes | `playwright.assets.mcp`, `playwright.steps.install` |
| Google Workspace | No | Yes | `gws.steps.install` |
| NotebookLM | No | Yes | `notebooklm.steps.install` |
| OpenSpec | No | Yes | `openspec.steps.install` |

## Kustomization Entrypoints

### Complete Package
- `kustomization.yaml` (domain root) — All agents, commands, steps, prompts, subagents, converters, integrations

### Category Aggregates
- `agents/kustomization.yaml` — All agents (opencode, pi)
- `cmds/kustomization.yaml` — All commands
- `steps/kustomization.yaml` — All shared steps
- `prompts/kustomization.yaml` — All prompts
- `subagents/kustomization.yaml` — All subagents
- `converters/kustomization.yaml` — All converters

### Selective Agent Entrypoints
- `agents/opencode/kustomization.yaml` — OpenCode agent (assets + converters)
- `agents/pi/kustomization.yaml` — Pi agent (assets + converters)

### Selective Integration Entrypoints
- `chrome-devtools/kustomization.yaml` — Chrome DevTools (assets + steps)
- `ctx7/kustomization.yaml` — Context7 (assets + steps)
- `playwright/kustomization.yaml` — Playwright (assets + steps)
- `notebooklm/kustomization.yaml` — NotebookLM (steps only)
- `gws/kustomization.yaml` — Google Workspace (steps only)
- `openspec/kustomization.yaml` — OpenSpec (steps only)

### Overlay Entrypoints
- `overlays/ai/kustomization.yaml` — AI wizard (active)
- `overlays/dev/kustomization.yaml` — Development (deprecated)
