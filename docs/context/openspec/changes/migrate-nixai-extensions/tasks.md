## 1. Kustomization files — leaf directories

- [ ] 1.1 Create `manifests/base/agents/converters/kustomization.yaml` (references: `to-json.yaml`)
- [ ] 1.2 Create `manifests/base/agents/steps/kustomization.yaml` (references: `settings.yaml`)
- [ ] 1.3 Create `manifests/base/cmds/kustomization.yaml` (references: `agents.yaml`)
- [ ] 1.4 Create `manifests/base/steps/kustomization.yaml` (references: `aggregate-mcp.yaml`, `cleanup.yaml`, `convert.yaml`, `copy-context.yaml`, `detect-agent.yaml`, `materialize-scaffold.yaml`)

## 2. Kustomization files — extensions

- [ ] 2.1 Create `manifests/base/extensions/self/assets/kustomization.yaml` (references: `commands`, `mcp`, `subagents`)
- [ ] 2.2 Create `manifests/base/extensions/self/converters/kustomization.yaml` (references: `commands`, `config`, `mcp`, `subagents`)
- [ ] 2.3 Create `manifests/base/extensions/self/kustomization.yaml` (references: `assets`, `converters`)
- [ ] 2.4 Create `manifests/base/extensions/ctx7/kustomization.yaml` (references: `assets`)
- [ ] 2.5 Create `manifests/base/extensions/chrome-devtools/kustomization.yaml` (references: `assets`)
- [ ] 2.6 Create `manifests/base/extensions/playwright/kustomization.yaml` (references: `assets`)
- [ ] 2.7 Create `manifests/base/extensions/gws/kustomization.yaml` (empty resources)
- [ ] 2.8 Create `manifests/base/extensions/notebooklm/kustomization.yaml` (empty resources)
- [ ] 2.9 Create `manifests/base/extensions/openspec/kustomization.yaml` (empty resources)
- [ ] 2.10 Create `manifests/base/extensions/kustomization.yaml` (references: `self`, `ctx7`, `chrome-devtools`, `playwright`, `gws`, `notebooklm`, `openspec`)

## 3. Kustomization files — parent directories

- [ ] 3.1 Create `manifests/base/agents/kustomization.yaml` (references: `claude.yaml`, `gemini.yaml`, `opencode.yaml`, `pi.yaml`, `converters`, `steps`)
- [ ] 3.2 Update `manifests/base/kustomization.yaml` (references: `agents`, `cmds`, `extensions`, `steps`)

## 4. Install Steps

- [ ] 4.1 Create `manifests/base/extensions/ctx7/steps/install.yaml` — generic ctx7 setup Step
- [ ] 4.2 Create `manifests/base/extensions/chrome-devtools/steps/install.yaml` — npx skills add Step
- [ ] 4.3 Create `manifests/base/extensions/playwright/steps/install.yaml` — npx skills add Step
- [ ] 4.4 Create `manifests/base/extensions/gws/steps/install.yaml` — npx skills add Step
- [ ] 4.5 Create `manifests/base/extensions/notebooklm/steps/install.yaml` — notebooklm skill install Step
- [ ] 4.6 Create `manifests/base/extensions/openspec/steps/install.yaml` — openspec init Step

## 5. Update extension kustomizations to include steps

- [ ] 5.1 Update `manifests/base/extensions/ctx7/kustomization.yaml` (add `steps`)
- [ ] 5.2 Update `manifests/base/extensions/chrome-devtools/kustomization.yaml` (add `steps`)
- [ ] 5.3 Update `manifests/base/extensions/playwright/kustomization.yaml` (add `steps`)
- [ ] 5.4 Update `manifests/base/extensions/gws/kustomization.yaml` (add `steps`)
- [ ] 5.5 Update `manifests/base/extensions/notebooklm/kustomization.yaml` (add `steps`)
- [ ] 5.6 Update `manifests/base/extensions/openspec/kustomization.yaml` (add `steps`)

## 6. Validation

- [ ] 6.1 Run `make build` to verify Go code compiles
- [ ] 6.2 Run `make test` to verify unit tests pass
- [ ] 6.3 Run `kfg build` with the reorganized manifests to verify loading
- [ ] 6.4 Run `make test-bats` to verify integration tests pass
