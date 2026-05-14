## 1. New Steps

- [x] 1.1 Create `.manifests/base/steps/ensure-gitignore.yaml` with step `kfg.ensure-gitignore`
- [x] 1.2 Create `.manifests/base/steps/inject-ctx7-context.yaml` with step `kfg.inject-ctx7-context`
- [x] 1.3 Update `.manifests/base/steps/kustomization.yaml` to include both new step files

## 2. Modify ctx7 Install Step

- [x] 2.1 Update `.manifests/base/extensions/ctx7/steps/install.yaml` — change `AGENT_HOME` default from `.claude` to `""` and `OUTPUT_DIR` from `.claude/skills/` to `""`

## 3. Rewrite Dev Workflow

- [x] 3.1 Rewrite `.manifests/overlay/dev/agents-workflow.yaml` with merged 10-phase workflow:
  - Phase 1 (-90): ensure-gitignore
  - Phase 2 (-70): detect-agent
  - Phase 3 (-65): materialize-scaffold ×4 + settings ×4
  - Phase 4 (-60): copy-context ×2
  - Phase 5 (-55): ctx7.install ×4 (per-agent with explicit AGENT_HOME/OUTPUT_DIR)
  - Phase 6 (-50): inject-ctx7-context ×2
  - Phase 7 (-45): convert commands ×4
  - Phase 8 (-40): aggregate-mcp ×3 (fix stray quote on gemini name)
  - Phase 9 (-35): convert subagent ×2
  - Phase 10 (90): cleanup ×4 (per-agent)
  - After: final cleanup
- [x] 3.2 Fix bug: remove stray quote from `agents.mcp.gemini"` name field

## 4. Validation

- [x] 4.1 Run `make test` to verify Go unit tests pass
- [x] 4.2 Run `make test-bats` to verify integration tests pass
- [x] 4.3 Verify `kfg build` generates correct shell output for the dev workflow
