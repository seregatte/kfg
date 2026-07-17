## Why

The wizard skill folder is created but empty — the asset used is the same as the command's, but
the skill needs its own format (SKILL.md with rich metadata). Additionally, the command should
instruct the agent to load the skill rather than duplicating the full wizard prompt.

## What Changes

1. **New asset**: `ai.prompts.wizard-command` — short command referencing the skill
2. **New converter**: `ai.opencode.conv.skill` — produces SKILL.md with metadata
3. **Workflow updated**: command uses new asset, skill uses new converter + existing asset
4. **Kustomizations updated**: converter and overlay register new resources

## Impact

- 2 new files created, 2 files modified
- Same functionality, better separation of concerns
