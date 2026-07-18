## Why

The kfg `.manifests/base/` structure was recently refactored to consolidate all AI-related extensions under a unified `extensions/ai/` directory with new naming conventions (`ai.*` instead of `kfg.agent.*`, `<ext>.*` instead of `kfg.extension.*`). The lifeos repository's `.manifests/overlay/dev/lifeos-workflow.yaml` still references the old naming patterns and deprecated steps (`kfg.convert`, `kfg.aggregate-mcp`, `kfg.agents.steps.settings`), making it incompatible with the current kfg base.

## What Changes

- Update all step references in `lifeos-workflow.yaml` to match new kfg naming conventions
- Replace deprecated `kfg.convert` and `kfg.aggregate-mcp` steps with unified `kfg.materialize` step
- Update converter references from `kfg.convert.self.*` to `ai.<agent>.conv.*` pattern
- Update asset references from `kfg.agent.*` to `ai.<agent>.asset.settings`
- Update extension install references from `kfg.extension.*.install` to `<ext>.steps.install`
- Update command references from `kfg.agent.cmd.*` to `ai.*.cmd.main`
- Update `kfg.extension.self.commands.git-commit` to `ai.prompts.git-commit`
- Local lifeos assets (commands, MCPs, subagents) maintain their current naming structure

## Capabilities

### New Capabilities
- `lifeos-workflow-refactor`: Updated workflow steps and references compatible with refactored kfg base

### Modified Capabilities
- None (lifeos has no existing specs)

## Impact

- **Affected files**: `lifeos/.manifests/overlay/dev/lifeos-workflow.yaml` (~80+ reference updates)
- **Compatible with**: kfg refactored manifests (commits `e41ab21`, `44f186a`, `19fa1e8`)
- **Validation**: `kfg apply -k .manifests/overlay/dev` must succeed after changes
- **No changes to**: `base/kustomization.yaml`, `overlay/dev/kustomization.yaml`, or local lifeos assets
