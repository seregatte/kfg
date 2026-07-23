## 1. Create ai-agents Extension Structure

- [x] 1.1 Create `.manifests/base/extensions/ai-agents/` directory with subdirectories: `agents/claude/`, `agents/gemini/`, `agents/opencode/`, `agents/pi/`, `cmds/`, `steps/`, `prompts/`, `subagents/`, `converters/`
- [x] 1.2 Create per-agent subdirectories: `assets/` and `converters/` under each agent
- [x] 1.3 Create all kustomization.yaml files (19 total)

## 2. Move and Rename Agent Assets

- [x] 2.1 Move `base/agents/claude.yaml` -> `extensions/ai-agents/agents/claude/assets/settings.yaml`; rename `metadata.name` to `ai.claude.asset.settings`
- [x] 2.2 Move `base/agents/gemini.yaml` -> `extensions/ai-agents/agents/gemini/assets/settings.yaml`; rename `metadata.name` to `ai.gemini.asset.settings`
- [x] 2.3 Move `base/agents/opencode.yaml` -> `extensions/ai-agents/agents/opencode/assets/settings.yaml`; rename `metadata.name` to `ai.opencode.asset.settings`
- [x] 2.4 Move `base/agents/pi.yaml` -> `extensions/ai-agents/agents/pi/assets/settings.yaml`; rename `metadata.name` to `ai.pi.asset.settings`

## 3. Move and Rename Agent Converters

- [x] 3.1 Move `extensions/self/converters/commands/claude.yaml` -> `agents/claude/converters/command.yaml`; rename to `ai.claude.conv.command`
- [x] 3.2 Move `extensions/self/converters/commands/gemini.yaml` -> `agents/gemini/converters/command.yaml`; rename to `ai.gemini.conv.command`
- [x] 3.3 Move `extensions/self/converters/commands/opencode.yaml` -> `agents/opencode/converters/command.yaml`; rename to `ai.opencode.conv.command`
- [x] 3.4 Move `extensions/self/converters/commands/pi.yaml` -> `agents/pi/converters/command.yaml`; rename to `ai.pi.conv.command`
- [x] 3.5 Move `extensions/self/converters/mcp/claude.yaml` -> `agents/claude/converters/mcp.yaml`; rename to `ai.claude.conv.mcp`
- [x] 3.6 Move `extensions/self/converters/mcp/gemini.yaml` -> `agents/gemini/converters/mcp.yaml`; rename to `ai.gemini.conv.mcp`
- [x] 3.7 Move `extensions/self/converters/mcp/opencode.yaml` -> `agents/opencode/converters/mcp.yaml`; rename to `ai.opencode.conv.mcp`
- [x] 3.8 Move `extensions/self/converters/subagents/claude.yaml` -> `agents/claude/converters/subagent.yaml`; rename to `ai.claude.conv.subagent`
- [x] 3.9 Move `extensions/self/converters/subagents/opencode.yaml` -> `agents/opencode/converters/subagent.yaml`; rename to `ai.opencode.conv.subagent`
- [x] 3.10 Move `extensions/self/converters/config/opencode.yaml` -> `agents/opencode/converters/cfg.yaml`; rename to `ai.opencode.conv.cfg`

## 4. Move and Rename Shared Resources

- [x] 4.1 Move `base/agents/converters/to-json.yaml` -> `extensions/ai-agents/converters/to-json.yaml`; rename to `ai.conv.to-json`
- [x] 4.2 Move `base/cmds/agents.yaml` -> `extensions/ai-agents/cmds/agents.yaml`; rename cmds to `ai.claude.cmd.main`, `ai.gemini.cmd.main`, `ai.opencode.cmd.main`, `ai.pi.cmd.main`
- [x] 4.3 Move `base/cmds/openspec.yaml` -> `extensions/ai-agents/cmds/openspec.yaml`; rename to `ai.cmds.openspec`
- [x] 4.4 Move `base/steps/detect-agent.yaml` -> `extensions/ai-agents/steps/detect.yaml`; rename to `ai.steps.detect`
- [x] 4.5 Move `extensions/self/assets/commands/git-commit.yaml` -> `extensions/ai-agents/prompts/git-commit.yaml`; rename to `ai.prompts.git-commit`
- [x] 4.6 Move `extensions/self/assets/commands/refactor-pure.yaml` -> `extensions/ai-agents/prompts/refactor-pure.yaml`; rename to `ai.prompts.refactor-pure`
- [x] 4.7 Move `extensions/self/assets/commands/review-code.yaml` -> `extensions/ai-agents/prompts/review-code.yaml`; rename to `ai.prompts.review-code`
- [x] 4.8 Move `extensions/self/assets/commands/review-search.yaml` -> `extensions/ai-agents/prompts/review-search.yaml`; rename to `ai.prompts.review-search`
- [x] 4.9 Move `extensions/self/assets/subagents/review-minimal.yaml` -> `extensions/ai-agents/subagents/review-minimal.yaml`; rename to `ai.subagents.review-minimal`

## 5. Move inject-ctx7-context to ctx7 Extension

- [x] 5.1 Move `base/steps/inject-ctx7-context.yaml` -> `extensions/ctx7/steps/inject-ctx7-context.yaml`; rename to `ctx7.steps.inject`
- [x] 5.2 Update `extensions/ctx7/steps/kustomization.yaml` to include `inject-ctx7-context.yaml`

## 6. Rename Non-AI Extension Resources

- [x] 6.1 Rename `extensions/ctx7/steps/install.yaml` metadata.name from `kfg.extension.ctx7.install` to `ctx7.steps.install`
- [x] 6.2 Rename `extensions/ctx7/assets/mcp.yaml` metadata.name from `kfg.extension.ctx7.mcp` to `ctx7.assets.mcp`
- [x] 6.3 Rename `extensions/openspec/steps/install.yaml` metadata.name from `kfg.extension.openspec.install` to `openspec.steps.install`
- [x] 6.4 Rename `extensions/chrome-devtools/assets/mcp.yaml` metadata.name from `kfg.extension.chrome-devtools.mcp` to `chrome.assets.mcp`
- [x] 6.5 Rename `extensions/playwright/assets/mcp.yaml` metadata.name from `kfg.extension.playwright.mcp` to `playwright.assets.mcp`

## 7. Update Kustomization References

- [x] 7.1 Update `base/kustomization.yaml`: remove `agents` and `cmds` refs
- [x] 7.2 Update `base/steps/kustomization.yaml`: remove `detect-agent.yaml` and `inject-ctx7-context.yaml` refs
- [x] 7.3 Update `base/extensions/kustomization.yaml`: add `ai-agents`, remove `self`
- [x] 7.4 Update `base/extensions/ctx7/steps/kustomization.yaml`: added `inject-ctx7-context.yaml`
- [x] 7.5 Remove `extensions/self/` directory (removed in phase 9)

## 8. Update Overlay Workflow

- [x] 8.1 Update `overlay/dev/agents-workflow.yaml`: replace all resource references with new names
- [x] 8.2 Verify `overlay/dev/kustomization.yaml` requires no changes — confirmed, no changes needed

## 9. Remove Old Directories

- [x] 9.1 Remove `base/agents/` directory and all contents
- [x] 9.2 Remove `base/cmds/` directory and all contents
- [x] 9.3 Remove `extensions/self/` directory and all contents

## 10. Update Documentation

- [x] 10.1 Update `docs/manifest-model.md`: replace AI-specific examples with generic examples (myapp.deploy, build, test)
- [x] 10.2 Update `docs/AGENTS.md`: add section about ai-agents extension structure
- [x] 10.3 Create `.manifests/base/extensions/ai-agents/README.md` documenting the extension
- [x] 10.4 Create `.manifests/base/extensions/ai-agents/agents/README.md` documenting per-agent structure

## 11. Bats Tests

- [x] 11.1 Create `tests/bats/manifests/base/extensions/ai-agents/ai-agents-structure.bats`: validate all expected files exist in the new structure
- [x] 11.2 Create `tests/bats/manifests/base/extensions/ai-agents/ai-agents-naming.bats`: validate all resource metadata.name follow the new convention (`ai.<agent>.asset.settings`, `ai.<agent>.cmd.main`, `ai.<agent>.conv.<type>`, `ai.cmds.<name>`, `ai.steps.<name>`, `ai.conv.<name>`, `ai.prompts.<name>`, `ai.subagents.<name>`)
- [x] 11.3 Create `tests/bats/manifests/base/extensions/ai-agents/ai-agents-kustomize.bats`: validate `kustomize build` succeeds for the base and overlay
- [x] 11.4 Create `tests/bats/manifests/base/extensions/ctx7/ctx7-naming.bats`: validate ctx7 resource names follow `ctx7.<kind>.<name>` convention
- [x] 11.5 Create `tests/bats/manifests/base/extensions/openspec/openspec-naming.bats`: validate openspec resource names follow `openspec.<kind>.<name>` convention
- [x] 11.6 Create `tests/bats/manifests/base/extensions/chrome-devtools/chrome-naming.bats`: validate chrome-devtools resource names follow `chrome.<kind>.<name>` convention
- [x] 11.7 Create `tests/bats/manifests/base/extensions/playwright/playwright-naming.bats`: validate playwright resource names follow `playwright.<kind>.<name>` convention

## 12. Validation

- [ ] 12.1 Run `kustomize build .manifests/base` and verify no errors
- [ ] 12.2 Run `kustomize build .manifests/overlay/dev` and verify no errors
- [x] 12.3 Run `make build && ./bin/kfg build .manifests/overlay/dev` and verify no errors
- [x] 12.4 Run `make test-bats` and verify all tests pass
- [x] 12.5 Run `make test` and verify Go tests still pass
