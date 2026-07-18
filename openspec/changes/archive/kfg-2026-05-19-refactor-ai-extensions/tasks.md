## 1. Create ai Extension Structure

- [x] 1.1 Create `extensions/ai/` directory
- [x] 1.2 Create `extensions/ai/kustomization.yaml` referencing all 13 subdirectories: `agents`, `cmds`, `steps`, `prompts`, `subagents`, `converters`, `chrome-devtools`, `ctx7`, `gws`, `notebooklm`, `openspec`, `playwright`

## 2. Move ai-agents Content

- [x] 2.1 Move `extensions/ai-agents/agents/` -> `extensions/ai/agents/`
- [x] 2.2 Move `extensions/ai-agents/cmds/` -> `extensions/ai/cmds/`
- [x] 2.3 Move `extensions/ai-agents/steps/` -> `extensions/ai/steps/`
- [x] 2.4 Move `extensions/ai-agents/prompts/` -> `extensions/ai/prompts/`
- [x] 2.5 Move `extensions/ai-agents/subagents/` -> `extensions/ai/subagents/`
- [x] 2.6 Move `extensions/ai-agents/converters/` -> `extensions/ai/converters/`

## 3. Move Extension Directories

- [x] 3.1 Move `extensions/chrome-devtools/` -> `extensions/ai/chrome-devtools/`
- [x] 3.2 Move `extensions/ctx7/` -> `extensions/ai/ctx7/`
- [x] 3.3 Move `extensions/gws/` -> `extensions/ai/gws/`
- [x] 3.4 Move `extensions/notebooklm/` -> `extensions/ai/notebooklm/`
- [x] 3.5 Move `extensions/openspec/` -> `extensions/ai/openspec/`
- [x] 3.6 Move `extensions/playwright/` -> `extensions/ai/playwright/`

## 4. Update Kustomization Files

- [x] 4.1 Update `extensions/kustomization.yaml`: remove all 7 refs, add only `ai`
- [x] 4.2 Remove old `extensions/ai-agents/` directory (content already moved)

## 5. Rename Resources (Naming Audit)

- [x] 5.1 Rename `extensions/ai/chrome-devtools/steps/install.yaml`: `kfg.extension.chrome-devtools.install` -> `chrome-devtools.steps.install`
- [x] 5.2 Rename `extensions/ai/gws/steps/install.yaml`: `kfg.extension.gws.install` -> `gws.steps.install`
- [x] 5.3 Rename `extensions/ai/notebooklm/steps/install.yaml`: `kfg.extension.notebooklm.install` -> `notebooklm.steps.install`
- [x] 5.4 Rename `extensions/ai/playwright/steps/install.yaml`: `kfg.extension.playwright.install` -> `playwright.steps.install`

## 6. Validation

- [x] 6.1 Run `kustomize build .manifests/base` and verify no errors
- [x] 6.2 Run `kustomize build .manifests/overlay/dev` and verify no errors
- [x] 6.3 Run `make build && ./bin/kfg build .manifests/overlay/dev` and verify no errors
- [x] 6.4 Run `make test-bats` and verify all tests pass
