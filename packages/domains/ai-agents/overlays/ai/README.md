# AI Wizard Overlay

The active overlay for the `kfg ai` interactive configuration wizard.

## Composition

This overlay composes:
- Framework base package
- Domain manifests (agents, commands, steps, prompts, subagents, converters, integrations)
- AI workflow (`ai-workflow.yaml`)
- Wizard prompts (`assets/prompts/wizard.yaml`, `assets/prompts/wizard-command.yaml`)

## Entrypoint

| Entrypoint | Description |
|------------|-------------|
| `kustomization.yaml` | Complete AI wizard overlay |

## Workflow

The `kfg.ai.workflow` CmdWorkflow orchestrates the wizard lifecycle:

1. **Ensure gitignore** (weight -90) — Adds agent directories to `.gitignore`
2. **Detect agent** (weight -70) — Identifies active agent (opencode or pi)
3. **Scaffold** (weight -65) — Creates agent directories per agent
4. **Install ctx7** (weight -55) — Installs Context7 library documentation skills
5. **Materialize wizard command** (weight -45) — Generates wizard entry point
6. **Materialize wizard skill** (weight -44) — Generates wizard skill file
7. **Cleanup** — Removes temporary files

## Generated Artifacts

| Artifact | Agent | Path |
|----------|-------|------|
| Gitignore entries | All | `.gitignore` (appended) |
| Agent directories | OpenCode | `.opencode/`, `.opencode/skills/`, `.opencode/commands/` |
| Agent directories | Pi | `.pi/`, `.pi/skills/`, `.pi/prompts/` |
| Wizard command | OpenCode | `.opencode/commands/wizard.md` |
| Wizard command | Pi | `.pi/prompts/wizard.md` |
| Wizard skill | OpenCode | `.opencode/skills/wizard/SKILL.md` |
| Wizard skill | Pi | `.pi/skills/wizard/SKILL.md` |
| ctx7 skills | All | Agent-specific skills directory |

## Cleanup

The `kfg.cleanup` step runs in the `after` phase and removes:
- Temporary ctx7 installation files
- Any intermediate artifacts from the workflow

## Migration

This overlay is the recommended way to generate project-specific kfg configurations.
Use `kfg ai` to start the wizard.
