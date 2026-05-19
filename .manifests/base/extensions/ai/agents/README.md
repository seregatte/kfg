# Per-Agent Structure

Each agent follows a consistent directory layout with separate `assets/` and `converters/` subdirectories.

## Directory Layout

```
agents/<name>/
├── kustomization.yaml
├── assets/
│   ├── kustomization.yaml
│   └── settings.yaml              # ai.<name>.asset.settings
└── converters/
    ├── kustomization.yaml
    ├── command.yaml               # ai.<name>.conv.command
    ├── mcp.yaml                   # ai.<name>.conv.mcp (optional)
    ├── subagent.yaml              # ai.<name>.conv.subagent (optional)
    └── cfg.yaml                   # ai.<name>.conv.cfg (optional)
```

## Resources

### Assets (`ai.<name>.asset.settings`)

Agent configuration data — permissions, models, UI settings, etc.

### Converters

Per-agent format converters that transform shared Assets/Converter resources into agent-specific output:

| Converter | Name Pattern | Purpose |
|-----------|-------------|---------|
| `command.yaml` | `ai.<name>.conv.command` | Convert prompts to agent command format |
| `mcp.yaml` | `ai.<name>.conv.mcp` | Convert MCP server configs to agent format |
| `subagent.yaml` | `ai.<name>.conv.subagent` | Convert subagent definitions to agent format |
| `cfg.yaml` | `ai.<name>.conv.cfg` | Convert to agent config file format |

Not all agents need all converters — only include what the agent supports.
