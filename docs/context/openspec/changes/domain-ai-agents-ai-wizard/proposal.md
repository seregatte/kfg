## Why

The kfg domain `ai-agents` currently relies on projects inheriting a shared overlay (`overlays/dev/`) and patching it for customization. This creates a fragile inheritance chain where projects must understand overlay composition and kustomize patching patterns. The barrier to entry is high: new users struggle with manifest inheritance, and existing users duplicate patch logic across projects. An interactive AI wizard removes this friction by generating complete, project-specific manifests tailored to each user's needs through conversation.

## What Changes

- **New CLI subcommand**: `kfg ai` — an interactive configuration wizard
- **New overlay**: `overlays/ai/` — lightweight workspace for the wizard agent (NOT a project template)
- **New step**: `ai.steps.deprecation-warn` — deprecation notice for `overlays/dev/`
- **Deprecation**: `overlays/dev/` logs a deprecation warning when used
- **New skill**: wizard skill prompt in `assets/prompts/wizard.yaml` — teaches the agent to compose building blocks, only creating custom manifests when nothing in the domain catalog covers the need

## Capabilities

### New Capabilities

- `ai-wizard-command`: The `kfg ai` CLI command and its interaction model (env var KFG_AI_AGENT, pass-through args, lifecycle)
- `ai-wizard-overlay`: The AI overlay that serves as the wizard agent workspace, including ctx7 integration and wizard command materialization
- `ai-wizard-skill`: The wizard skill prompt that instructs agents on how to compose kfg configurations from domain building blocks

### Modified Capabilities

- `domain-ai-agents-skill-installation-steps`: The skill installation path gains a new consumer (the wizard overlay) and the ctx7 steps mechanism must accommodate the wizard workflow

## Impact

- **Domain `ai-agents`**: New `overlays/ai/` directory; new step in `manifests/steps/`; new asset in `overlays/ai/assets/prompts/`
- **Framework**: No changes (wizard uses existing steps: materialize, ensure-gitignore, cleanup, materialize-scaffold)
- **Engine**: New CLI subcommand `kfg ai` (new `src/cmd/kfg/ai.go`) that wraps `kfg run`
- **Backward compatibility**: `overlays/dev/` is NOT removed — only deprecated with a warning. All existing projects continue to work unchanged.
