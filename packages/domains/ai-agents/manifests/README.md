# AI Agents Extension

This extension consolidates all AI agent resources into a single, organized structure under `extensions/ai-agents/`.

## Structure

```
ai-agents/
├── agents/                    # Per-agent resources
│   ├── claude/
│   │   ├── assets/            # Agent settings (ai.claude.asset.settings)
│   │   └── converters/        # Per-agent converters (ai.claude.conv.*)
│   ├── gemini/
│   ├── opencode/
│   └── pi/
├── cmds/                      # Shared command wrappers (ai.cmds.*)
├── steps/                     # Shared workflow steps (ai.steps.*)
├── prompts/                   # Shared prompt templates (ai.prompts.*)
├── subagents/                 # Shared subagent definitions (ai.subagents.*)
└── converters/                # Shared converters (ai.conv.*)
```

## Naming Convention

All resources follow a short, consistent naming convention:

| Pattern | Example | Description |
|---------|---------|-------------|
| `ai.<agent>.asset.settings` | `ai.claude.asset.settings` | Agent settings |
| `ai.<agent>.cmd.main` | `ai.opencode.cmd.main` | Agent command wrapper |
| `ai.<agent>.conv.<type>` | `ai.claude.conv.mcp` | Agent converter |
| `ai.cmds.<name>` | `ai.cmds.openspec` | Shared command |
| `ai.steps.<name>` | `ai.steps.detect` | Shared step |
| `ai.conv.<name>` | `ai.conv.to-json` | Shared converter |
| `ai.prompts.<name>` | `ai.prompts.git-commit` | Shared prompt |
| `ai.subagents.<name>` | `ai.subagents.review-minimal` | Shared subagent |

## Adding a New Agent

1. Create a new directory under `agents/<name>/`
2. Add `assets/settings.yaml` with `metadata.name: ai.<name>.asset.settings`
3. Add `converters/` with converter files named `ai.<name>.conv.<type>`
4. Update `agents/kustomization.yaml` to include the new agent
5. Add a Cmd entry in `cmds/agents.yaml` with `metadata.name: ai.<name>.cmd.main`
