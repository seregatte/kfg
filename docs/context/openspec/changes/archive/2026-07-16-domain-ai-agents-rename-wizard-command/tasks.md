## 1. Asset Rename

- [x] 1.1 Rename file `packages/domains/ai-agents/overlays/ai/assets/prompts/kfg-wizard.yaml` → `wizard.yaml`
- [x] 1.2 Change `spec.data.name` from `kfg-wizard` to `wizard` in `wizard.yaml`

## 2. Overlay Updates

- [x] 2.1 Update resource reference in `packages/domains/ai-agents/overlays/ai/kustomization.yaml` from `assets/prompts/kfg-wizard.yaml` to `assets/prompts/wizard.yaml`
- [x] 2.2 Update OUTPUTS for opencode from `.opencode/commands/kfg-wizard.md` to `.opencode/commands/wizard.md` in `ai-workflow.yaml`
- [x] 2.3 Update OUTPUTS for pi from `.pi/prompts/kfg-wizard.md` to `.pi/prompts/wizard.md` in `ai-workflow.yaml`

## 3. Wizard Prompt Cleanup

- [x] 3.1 Remove `claude` and `gemini` from agent catalog in `wizard.yaml`
- [x] 3.2 Remove `ai.claude.cmd.main` and `ai.gemini.cmd.main` from Cmds catalog in `wizard.yaml`
- [x] 3.3 Update example conversations in `wizard.yaml` to not reference claude/gemini

## 4. Documentation

- [x] 4.1 Update file path reference in `docs/AGENTS.md` from `kfg-wizard.yaml` to `wizard.yaml`
- [x] 4.2 Update `docs/context/openspec/changes/domain-ai-agents-ai-wizard/proposal.md` references
- [x] 4.3 Update `docs/context/openspec/changes/domain-ai-agents-ai-wizard/design.md` references
- [x] 4.4 Update `docs/context/openspec/changes/domain-ai-agents-ai-wizard/tasks.md` references
- [x] 4.5 Update `docs/context/openspec/changes/domain-ai-agents-ai-wizard/specs/ai-wizard-skill/spec.md`
- [x] 4.6 Update `docs/context/openspec/changes/domain-ai-agents-ai-wizard/specs/ai-wizard-overlay/spec.md`

## 5. Test Updates

- [x] 5.1 Update file path and name references in `packages/domains/ai-agents/tests/base/extensions/ai-wizard-overlay.bats`

## 6. Validation

- [x] 6.1 Verify `kfg build packages/domains/ai-agents/overlays/ai/` produces valid output with new paths
- [x] 6.2 Run `make test-bats` — all Bats tests pass
- [x] 6.3 Run `make build` — binary builds successfully
