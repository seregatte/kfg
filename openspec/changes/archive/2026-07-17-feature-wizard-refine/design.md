## Context

Following the OpenSpec pattern (research via ctx7), commands and skills have distinct roles:
- **Commands**: Explicit entry points, minimal front matter, action-oriented
- **Skills**: Contextual discovery, rich metadata, knowledge-oriented

The current wizard materializes the same asset to both locations with the same converter,
causing the skill to be empty (or contain wrong format).

## Architecture

### Assets
- `ai.prompts.wizard-command` → short prompt: "load the wizard skill"
- `ai.prompts.kfg-wizard` → full wizard instructions (unchanged)

### Converters
- `ai.opencode.conv.command` → Command format (existing)
- `ai.opencode.conv.skill` → SKILL.md with rich metadata (new)

Expression:
```
"---\nname: " + .name + "\ndescription: |\n  " + (.description | gsub("\n"; "\n  ")) + "\nlicense: MIT\ncompatibility: Requires kfg CLI.\nmetadata:\n  author: kfg\n  version: \"1.0\"\n---\n\n" + .prompt
```

### Workflow Steps
- Phase 5a (-45): Command materialize (ai.prompts.wizard-command → .opencode/commands/wizard.md)
- Phase 5b (-44): Skill materialize (ai.prompts.kfg-wizard → .opencode/skills/wizard/SKILL.md)

## Non-Goals
- Not changing the wizard prompt content
- Not modifying pi agent (same pattern would apply)
