## Context

Current structure has 7 flat extensions under `extensions/`:
```
extensions/
├── ai-agents/          (agents, cmds, steps, prompts, subagents, converters)
├── chrome-devtools/    (assets, steps)
├── ctx7/               (assets, steps)
├── gws/                (steps)
├── notebooklm/         (steps)
├── openspec/           (steps)
├── playwright/         (assets, steps)
└── kustomization.yaml  (references all 7)
```

Naming inconsistency: 4 extensions use old `kfg.extension.*` prefix while 2 use the short convention.

## Goals / Non-Goals

**Goals:**
- Consolidate all 7 extensions into single `ai/` directory
- Standardize all resource names to `<ext>.<kind>.<name>` format
- Update kustomization references
- Maintain all existing functionality

**Non-Goals:**
- Change Go code or CLI behavior
- Add new agents or functionality
- Change `apiVersion` from `kfg.dev/v1alpha1`
- Add or update READMEs

## Decisions

### Consolidate into `extensions/ai/`

All extension directories become subdirectories of `ai/`:
```
extensions/
└── ai/
    ├── agents/          (claude, gemini, opencode, pi - each with assets/converters)
    ├── cmds/            (agents.yaml, openspec.yaml)
    ├── steps/           (detect.yaml)
    ├── prompts/         (git-commit, refactor-pure, review-code, review-search)
    ├── subagents/       (review-minimal)
    ├── converters/      (to-json)
    ├── chrome-devtools/ (assets, steps)
    ├── ctx7/            (assets, steps)
    ├── gws/             (steps)
    ├── notebooklm/      (steps)
    ├── openspec/        (steps)
    └── playwright/      (assets, steps)
```

Parent `ai/kustomization.yaml` references all subdirectories.
`extensions/kustomization.yaml` references only `ai`.

### Naming Convention

| Current Name | New Name |
|---|---|
| `kfg.extension.chrome-devtools.install` | `chrome-devtools.steps.install` |
| `kfg.extension.gws.install` | `gws.steps.install` |
| `kfg.extension.notebooklm.install` | `notebooklm.steps.install` |
| `kfg.extension.playwright.install` | `playwright.steps.install` |

Already correct: `ctx7.steps.install`, `ctx7.steps.inject`, `openspec.steps.install`, `chrome.assets.mcp`, `playwright.assets.mcp`, all `ai.*` resources.

## Data Contract

### Kustomization hierarchy

```yaml
# extensions/kustomization.yaml
resources:
  - ai

# extensions/ai/kustomization.yaml
resources:
  - agents
  - cmds
  - steps
  - prompts
  - subagents
  - converters
  - chrome-devtools
  - ctx7
  - gws
  - notebooklm
  - openspec
  - playwright
```

## Risks / Trade-offs

- [Larger single directory] -> `ai/` has 13 subdirectories; acceptable since they are logical groupings
- [Deep nesting] -> Max depth increases by 1; mitigated by organization
- [Breaking existing overlays] -> Only dev overlay exists; audited — step references to renamed resources are not in workflow, only asset refs `chrome.assets.mcp` and `playwright.assets.mcp` are already correct

## Migration Plan

1. Create `extensions/ai/` with kustomization.yaml
2. Move ai-agents content into `ai/`
3. Move each extension directory into `ai/`
4. Update `extensions/kustomization.yaml`
5. Rename resource metadata names
6. Remove old `ai-agents/` directory
7. Validate with kustomize build and kfg build
