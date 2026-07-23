## Context

The AI overlay's workflow (`ai-workflow.yaml`) orchestrates setup, detection, scaffolding, ctx7 install, and wizard command materialization. Adding skill output follows the same pattern: same asset, same converter, different target path.

## Architecture

```
asset: ai.prompts.kfg-wizard
      │
      ├── converter: ai.opencode.conv.command
      │     ├── output: .opencode/commands/wizard.md        (weight -45, existing)
      │     └── output: .opencode/skills/wizard/SKILL.md    (weight -44, new)
      │
      └── converter: ai.pi.conv.command
            ├── output: .pi/prompts/wizard.md               (weight -45, existing)
            └── output: .pi/skills/wizard/SKILL.md          (weight -44, new)
```

## Key Design Decisions

1. **Same asset, same converter**: The wizard prompt content is identical for command and skill roles. The differentiation is purely about file location — commands and skills in opencode/pi share the same markdown+YAML format.
2. **Weight -44**: Runs immediately after the command materialization (weight -45). Order doesn't matter in practice (they're independent steps), but -44 keeps the sequence logical.
3. **mkdir -p handles subdirs**: `kfg.materialize` calls `mkdir -p "$(dirname "$output")"`, so `skills/wizard/` is auto-created. No scaffold update needed.

## Non-Goals

- Not changing the wizard prompt content
- Not adding new converters or assets
- Not modifying pi or opencode settings
