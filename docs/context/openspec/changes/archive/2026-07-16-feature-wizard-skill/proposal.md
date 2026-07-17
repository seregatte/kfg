## Why

The AI wizard prompt is currently materialized only as a command for opencode (`.opencode/commands/wizard.md`) and a prompt for pi (`.pi/prompts/wizard.md`). Users must explicitly invoke `/wizard` to start it. Making it also available as a skill allows opencode and pi to discover it contextually, providing a richer workflow that can be auto-triggered when the agent detects the user needs configuration generation.

## What Changes

Add two new `kfg.materialize` steps in `ai-workflow.yaml` (weight -44) that output the same wizard asset (`ai.prompts.kfg-wizard`) to skill directories:

1. `.opencode/skills/wizard/SKILL.md` for opencode
2. `.pi/skills/wizard/SKILL.md` for pi

## Impact

- **Single file change**: Only `packages/domains/ai-agents/overlays/ai/ai-workflow.yaml`
- No new Assets, Converters, or scaffolding needed
- Compatible with both agents' existing settings
