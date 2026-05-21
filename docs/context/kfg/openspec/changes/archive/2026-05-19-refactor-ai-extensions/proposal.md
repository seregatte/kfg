## Why

All extensions under `.manifests/base/extensions/` (ai-agents, chrome-devtools, ctx7, gws, notebooklm, openspec, playwright) represent AI agent tooling and configuration. The current flat structure treats them as separate top-level extensions, but they are conceptually a single unified "ai" extension with multiple sub-components. This creates unnecessary nesting and makes the manifest structure harder to navigate.

Additionally, naming is inconsistent: `ctx7` and `openspec` use the short convention (`ctx7.steps.install`, `openspec.steps.install`) while `chrome-devtools`, `gws`, `notebooklm`, and `playwright` still use the old `kfg.extension.*` prefix (`kfg.extension.chrome-devtools.install`).

## What Changes

- Rename `ai-agents/` to `ai/` and move all other extensions inside it as subdirectories
- New structure: `extensions/ai/{agents,cmds,steps,converters,prompts,subagents,chrome-devtools,ctx7,gws,notebooklm,openspec,playwright}`
- Update `extensions/kustomization.yaml` to reference only `ai`
- Move each extension's `kustomization.yaml` resources into the parent `ai/kustomization.yaml`
- Rename all resources to short convention: `chrome-devtools.steps.install`, `gws.steps.install`, `notebooklm.steps.install`, `playwright.steps.install`
- No README changes

## Non-Goals

- Change Go code or CLI behavior
- Add new extensions or agents
- Change `apiVersion` from `kfg.dev/v1alpha1`
- Rename core generic steps (`kfg.cleanup`, `kfg.materialize`, etc.)

## Capabilities

### Modified Capabilities

- `ai_extension_consolidation`: consolidate all 7 extensions into single `ai/` directory
- `manifest_naming_convention`: standardize all resource names to `<ext>.<kind>.<name>` format
- `kustomization_hierarchy`: restructure extension kustomization references

## Impact

- Affected: all extension directories under `.manifests/base/extensions/`
- Affected: `base/extensions/kustomization.yaml`, all extension kustomization files
- Affected: `overlay/dev/agents-workflow.yaml` (references to renamed resources)
- User-facing: resource names change, workflow behavior identical
