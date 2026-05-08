## 1. Directory Structure Setup

- [ ] 1.1 Create `.manifests/base/agents/` directory with placeholder YAML files for claude, gemini, opencode, pi
- [ ] 1.2 Create `.manifests/base/cmds/` directory with `agents.yaml` Cmd definitions
- [ ] 1.3 Create `.manifests/base/extensions/self/assets/{commands,subagents,mcp}/` directories
- [ ] 1.4 Create `.manifests/base/extensions/self/converters/{commands,mcp,subagents}/` directories
- [ ] 1.5 Create `.manifests/base/extensions/self/steps/` directory
- [ ] 1.6 Create `.manifests/base/extensions/{ctx7,playwright,chrome-devtools,openspec,gws,notebooklm,ccr}/` with appropriate subdirectories
- [ ] 1.7 Create `.manifests/base/steps/` directory
- [ ] 1.8 Create `.manifests/overlay/dev/` directory
- [ ] 1.9 Create `.manifests/tests/` directory

## 2. Base Agent Assets

- [ ] 2.1 Write `.manifests/base/agents/claude.yaml` with settings.json and .mcp.json data (using `{env:HOME}` placeholder)
- [ ] 2.2 Write `.manifests/base/agents/gemini.yaml` with settings.json data
- [ ] 2.3 Write `.manifests/base/agents/opencode.yaml` with opencode.json data
- [ ] 2.4 Write `.manifests/base/agents/pi.yaml` with settings.json data
- [ ] 2.5 Write `.manifests/base/cmds/agents.yaml` with 4 Cmd resources (opencode, gemini, pi, claude)

## 3. Self Extension — Assets

- [ ] 3.1 Write 4 command Assets in `.manifests/base/extensions/self/assets/commands/` (git-commit, refactor-pure, review-code, review-search)
- [ ] 3.2 Write 1 subagent Asset in `.manifests/base/extensions/self/assets/subagents/` (review-minimal)
- [ ] 3.3 Write 3 MCP Assets in `.manifests/base/extensions/self/assets/mcp/` (context7, chrome-devtools, playwright)

## 4. Self Extension — Converters

- [ ] 4.1 Write 4 command Converters in `.manifests/base/extensions/self/converters/commands/` (claude, gemini, opencode, pi)
- [ ] 4.2 Write 3 MCP Converters in `.manifests/base/extensions/self/converters/mcp/` (claude, gemini, opencode)
- [ ] 4.3 Write 2 subagent Converters in `.manifests/base/extensions/self/converters/subagents/` (claude, opencode)
- [ ] 4.4 Verify converter expressions produce correct output using `kfg apply --convert` CLI

## 5. Self Extension — Steps

- [ ] 5.1 Write `kfg.generate-commands` Step — parameterized with ASSET_PREFIX, CONVERTER_PREFIX, OUTPUT_DIR, FILE_EXT
- [ ] 5.2 Write `kfg.install-mcp` Step — parameterized with MCP_ASSET, CONVERTER_PREFIX, OUTPUT_FILE
- [ ] 5.3 Write `kfg.generate-subagents` Step — parameterized with ASSET_PREFIX, CONVERTER_PREFIX, OUTPUT_DIR

## 6. Other Extension Assets

- [ ] 6.1 Write `.manifests/base/extensions/ctx7/assets/mcp.yaml`
- [ ] 6.2 Write `.manifests/base/extensions/playwright/assets/mcp.yaml`
- [ ] 6.3 Write `.manifests/base/extensions/chrome-devtools/assets/mcp.yaml`
- [ ] 6.4 Write `.manifests/base/extensions/openspec/assets/skill.yaml`
- [ ] 6.5 Write `.manifests/base/extensions/gws/assets/skill.yaml`
- [ ] 6.6 Write `.manifests/base/extensions/notebooklm/assets/skill.yaml`
- [ ] 6.7 Write `.manifests/base/extensions/ccr/assets/config.yaml`

## 7. Core Steps

- [ ] 7.1 Write `kfg.detect-agent` Step
- [ ] 7.2 Write `kfg.copy-context` Step
- [ ] 7.3 Write `kfg.materialize-scaffold` Step — reads Assets and writes files to filesystem
- [ ] 7.4 Write `kfg.cleanup-paths` Step
- [ ] 7.5 Write `kfg.cleanup-workspace` Step
- [ ] 7.6 Write `kfg.cleanup` Step — removes `$KFG_BUILD_RESULT_FILE`
- [ ] 7.7 Write `kfg.install-skill` Step
- [ ] 7.8 Write `kfg.setup-ccr` Step

## 8. Extension-Specific Steps

- [ ] 8.1 Write `kfg.extension.ctx7.install-skills` Step — parameterized with CTX7_FLAG, AGENT_HOME, SKILLS_OUTPUT_DIR
- [ ] 8.2 Verify step has zero conditional logic (no `case`, `if` for agent routing)

## 9. Kustomization & Overlay

- [ ] 9.1 Write `.manifests/base/kustomization.yaml` referencing all base resources
- [ ] 9.2 Write `.manifests/overlay/dev/kustomization.yaml` referencing `../../base`
- [ ] 9.3 Write `.manifests/overlay/dev/cmds.yaml` with openspec Cmd
- [ ] 9.4 Write `.manifests/overlay/dev/agents-workflow.yaml` CmdWorkflow with all phases and `when` conditions
- [ ] 9.5 Validate kustomize loads successfully: `kfg build -k .manifests/`

## 10. Integration Testing

- [ ] 10.1 Verify `kfg build -k .manifests/` produces valid shell with all step functions
- [ ] 10.2 Verify converter CLI works: `kfg apply -k .manifests/ --convert kfg.extension.self.commands.git-commit --use kfg.convert.self.command.claude`
- [ ] 10.3 Verify workflow generates correct `when` condition shell code for agent detection
- [ ] 10.4 Verify all 43 YAML files parse without errors

## 11. Bats Tests

- [ ] 11.1 Write `.manifests/tests/test_helper.bash`
- [ ] 11.2 Write `manifest-loading.bats` — validates kustomize load, resource count, kind distribution
- [ ] 11.3 Write `converters.bats` — validates converter output for each agent (claude, gemini, opencode, pi)
- [ ] 11.4 Write `workflow.bats` — validates CmdWorkflow shell generation (step functions, when conditions)
- [ ] 11.5 Add `make test-manifests` target to Makefile
- [ ] 11.6 Add `make test-all` target to Makefile (runs test, test-bats, test-manifests)
