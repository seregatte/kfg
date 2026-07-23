## Why

The wizard skill command is currently named `/kfg-wizard`, which is unnecessarily verbose since it only exists within a kfg-managed agent environment. Renaming to `/wizard` is shorter and more intuitive for users.

## What Changes

- **Rename**: `ai.prompts.kfg-wizard` asset `spec.data.name` from `kfg-wizard` to `wizard`
- **Rename file**: `assets/prompts/kfg-wizard.yaml` → `assets/prompts/wizard.yaml`
- **Update**: Materialize OUTPUTS from `kfg-wizard.md` to `wizard.md` in `ai-workflow.yaml`
- **Update**: Resource reference in `overlays/ai/kustomization.yaml`
- **Update**: All OpenSpec spec files and design docs referencing the old name
- **Update**: `docs/AGENTS.md` file path reference
- **Update**: Bats test references in `ai-wizard-overlay.bats`
- **Clean up**: Remove stale `claude`/`gemini` references in wizard prompt catalog

## Capabilities

### Modified Capabilities

- `ai-wizard-command`: The wizard command name changes from `/kfg-wizard` to `/wizard`
- `ai-wizard-overlay`: Output paths change from `kfg-wizard.md` to `wizard.md`
- `ai-wizard-skill`: Skill prompt asset name and file change

## Impact

- **Domain `ai-agents`**: File rename in `overlays/ai/assets/prompts/`; kustomization update
- **Framework**: No changes
- **Engine**: No changes (`kfg ai` CLI command unchanged)
- **Backward compatibility**: Minimal — users who manually reference `/kfg-wizard` in their agent configs will need to update to `/wizard`.
