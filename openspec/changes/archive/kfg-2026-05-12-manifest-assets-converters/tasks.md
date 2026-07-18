## 1. Directory Structure Setup

- [x] 1.1 Create `.manifests/base/agents/` directory with placeholder YAML files for claude, gemini, opencode, pi
- [x] 1.2 Create `.manifests/base/cmds/` directory with `agents.yaml` Cmd definitions
- [x] 1.3 Create `.manifests/base/extensions/self/assets/{commands,subagents,mcp}/` directories
- [x] 1.4 Create `.manifests/base/extensions/self/converters/{commands,mcp,subagents}/` directories
- [x] 1.5 Create `.manifests/base/extensions/self/steps/` directory
- [x] 1.6 Create `.manifests/base/extensions/{ctx7,playwright,chrome-devtools,openspec,gws,notebooklm,ccr}/` with appropriate subdirectories
- [x] 1.7 Create `.manifests/base/steps/` directory
- [x] 1.8 Create `.manifests/overlay/dev/` directory
- [x] 1.9 Create `.manifests/tests/` directory

## 2. Base Agent Assets

- [x] 2.1 Write `.manifests/base/agents/claude.yaml` with settings.json and .mcp.json data (using `{env:HOME}` placeholder)
- [x] 2.2 Write `.manifests/base/agents/gemini.yaml` with settings.json data
- [x] 2.3 Write `.manifests/base/agents/opencode.yaml` with opencode.json data
- [x] 2.4 Write `.manifests/base/agents/pi.yaml` with settings.json data
- [x] 2.5 Write `.manifests/base/cmds/agents.yaml` with 4 Cmd resources (opencode, gemini, pi, claude)

## 3. Self Extension — Assets

- [x] 3.1 Write 4 command Assets in `.manifests/base/extensions/self/assets/commands/` (git-commit, refactor-pure, review-code, review-search)
- [x] 3.2 Write 1 subagent Asset in `.manifests/base/extensions/self/assets/subagents/` (review-minimal)
- [x] 3.3 Write 3 MCP Assets in `.manifests/base/extensions/self/assets/mcp/` (context7, chrome-devtools, playwright)

## 4. Self Extension — Converters

- [x] 4.1 Write 4 command Converters in `.manifests/base/extensions/self/converters/commands/` (claude, gemini, opencode, pi)
- [x] 4.2 Write 3 MCP Converters in `.manifests/base/extensions/self/converters/mcp/` (claude, gemini, opencode)
- [x] 4.3 Write 2 subagent Converters in `.manifests/base/extensions/self/converters/subagents/` (claude, opencode)
- [x] 4.4 Verify converter expressions produce correct output using `kfg apply --convert` CLI

## 5. Self Extension — Steps

- [x] 5.1 Write `kfg.generate-commands` Step — parameterized with ASSET_PREFIX, CONVERTER_PREFIX, OUTPUT_DIR, FILE_EXT
- [x] 5.2 Write `kfg.install-mcp` Step — parameterized with MCP_ASSET, CONVERTER_PREFIX, OUTPUT_FILE
- [x] 5.3 Write `kfg.generate-subagents` Step — parameterized with ASSET_PREFIX, CONVERTER_PREFIX, OUTPUT_DIR

## 6. Other Extension Assets

- [x] 6.1 Write `.manifests/base/extensions/ctx7/assets/mcp.yaml`
- [x] 6.2 Write `.manifests/base/extensions/playwright/assets/mcp.yaml`
- [x] 6.3 Write `.manifests/base/extensions/chrome-devtools/assets/mcp.yaml`
- [x] 6.4 Write `.manifests/base/extensions/openspec/assets/skill.yaml`
- [x] 6.5 Write `.manifests/base/extensions/gws/assets/skill.yaml`
- [x] 6.6 Write `.manifests/base/extensions/notebooklm/assets/skill.yaml`
- [x] 6.7 Write `.manifests/base/extensions/ccr/assets/config.yaml`

## 7. Core Steps

- [x] 7.1 Write `kfg.detect-agent` Step
- [x] 7.2 Write `kfg.copy-context` Step
- [x] 7.3 Write `kfg.materialize-scaffold` Step — reads Assets and writes files to filesystem
- [x] 7.4 Write `kfg.cleanup-paths` Step
- [x] 7.5 Write `kfg.cleanup-workspace` Step
- [x] 7.6 Write `kfg.cleanup` Step — removes `$KFG_BUILD_RESULT_FILE`
- [x] 7.7 Write `kfg.install-skill` Step
- [x] 7.8 Write `kfg.setup-ccr` Step

## 8. Extension-Specific Steps

- [x] 8.1 Write `kfg.extension.ctx7.install-skills` Step — parameterized with CTX7_FLAG, AGENT_HOME, SKILLS_OUTPUT_DIR
- [x] 8.2 Verify step has zero conditional logic (no `case`, `if` for agent routing)

## 9. Kustomization & Overlay

- [x] 9.1 Write `.manifests/base/kustomization.yaml` referencing all base resources
- [x] 9.2 Write `.manifests/overlay/dev/kustomization.yaml` referencing `../../base`
- [x] 9.3 Write `.manifests/overlay/dev/cmds.yaml` with openspec Cmd
- [x] 9.4 Write `.manifests/overlay/dev/agents-workflow.yaml` CmdWorkflow with all phases and `when` conditions
- [x] 9.5 Validate kustomize loads successfully: `kfg build .manifests/overlay/dev`

## 10. Integration Testing

- [x] 10.1 Verify `kfg build .manifests/overlay/dev` produces valid shell with all step functions
- [x] 10.2 Verify converter CLI works: `kfg apply -k .manifests/overlay/dev --convert kfg.extension.self.commands.git-commit --use kfg.convert.self.command.claude`
- [x] 10.3 Verify workflow generates correct `when` condition shell code for agent detection
- [x] 10.4 Verify all YAML files parse without errors

## 11. Bats Tests

- [x] 11.1 Write `.manifests/tests/test_helper.bash`
- [x] 11.2 Write `manifest-loading.bats` — validates kustomize load, resource count, kind distribution
- [x] 11.3 Write `converters.bats` — validates converter output for each agent (claude, gemini, opencode, pi)
- [x] 11.4 Write `workflow.bats` — validates CmdWorkflow shell generation (step functions, when conditions)
- [x] 11.5 Add `make test-manifests` target to Makefile
- [x] 11.6 Add `make test-all` target to Makefile (runs test, test-bats, test-manifests)
