# AI Agents Manifests

This directory contains all AI agents domain manifests, including agents, commands,
steps, prompts, subagents, converters, and integrations.

## Package Structure

```
manifests/
├── agents/                    # Per-agent resources
│   ├── opencode/
│   └── pi/
├── cmds/                      # Shared command wrappers (ai.cmds.*)
├── steps/                     # Shared workflow steps (ai.steps.*)
├── prompts/                   # Shared prompt templates (ai.prompts.*)
├── subagents/                 # Shared subagent definitions (ai.subagents.*)
├── converters/                # Shared converters (ai.conv.*)
├── chrome-devtools/           # Chrome DevTools MCP integration
├── ctx7/                      # Context7 library documentation integration
├── gws/                       # Google Workspace integration
├── notebooklm/                # NotebookLM integration
├── openspec/                  # OpenSpec planning integration
└── playwright/                # Playwright E2E testing integration
```

## Naming Convention

All resources follow a consistent naming convention:

| Pattern | Example | Description |
|---------|---------|-------------|
| `ai.<agent>.asset.settings` | `ai.opencode.asset.settings` | Agent settings |
| `ai.<agent>.cmd.main` | `ai.opencode.cmd.main` | Agent command wrapper |
| `ai.<agent>.conv.<type>` | `ai.opencode.conv.mcp` | Agent converter |
| `ai.cmds.<name>` | `ai.cmds.openspec` | Shared command |
| `ai.steps.<name>` | `ai.steps.detect` | Shared step |
| `ai.conv.<name>` | `ai.conv.to-json` | Shared converter |
| `ai.prompts.<name>` | `ai.prompts.git-commit` | Shared prompt |
| `ai.subagents.<name>` | `ai.subagents.review-minimal` | Shared subagent |
| `<integration>.assets.mcp` | `ctx7.assets.mcp` | Integration MCP assets |
| `<integration>.steps.<name>` | `ctx7.steps.install` | Integration step |

## Capability Documentation

Each independently consumable capability has an adjacent `README.md` documenting
its purpose, prerequisites, entrypoints, exported resources, agent compatibility,
configuration, JSON Patch examples, workflow usage, generated artifacts, cache and
cleanup behavior, limitations, validation, and canonical specs.

| Capability | Documentation |
|------------|--------------|
| OpenCode agent | [agents/opencode/README.md](agents/opencode/README.md) |
| Pi agent | [agents/pi/README.md](agents/pi/README.md) |
| Chrome DevTools | [chrome-devtools/README.md](chrome-devtools/README.md) |
| Context7 | [ctx7/README.md](ctx7/README.md) |
| Google Workspace | [gws/README.md](gws/README.md) |
| NotebookLM | [notebooklm/README.md](notebooklm/README.md) |
| OpenSpec | [openspec/README.md](openspec/README.md) |
| Playwright | [playwright/README.md](playwright/README.md) |

## Adding a New Agent

1. Create a new directory under `agents/<name>/`
2. Add `assets/settings.yaml` with `metadata.name: ai.<name>.asset.settings`
3. Add `converters/` with converter files as needed
4. Update `agents/kustomization.yaml` to include the new agent
5. Add a Cmd entry in `cmds/agents.yaml` with `metadata.name: ai.<name>.cmd.main`
6. Add an adjacent `README.md` following the capability README template

## Adding a New Integration

1. Create a new directory under `<integration-name>/`
2. Add `assets/mcp.yaml` if the integration provides MCP tools
3. Add `steps/install.yaml` for skill installation
4. Add a `kustomization.yaml` entrypoint
5. Add an adjacent `README.md` following the capability README template
