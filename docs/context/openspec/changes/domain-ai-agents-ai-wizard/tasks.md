## 1. Deprecation Step

- [x] 1.1 Create `packages/domains/ai-agents/manifests/steps/deprecation-warn.yaml` with `kind: Step` and `metadata.name: ai.steps.deprecation-warn` containing a `__kfg_log_warn` message directing users to `kfg ai`
- [x] 1.2 Register the new step in `packages/domains/ai-agents/manifests/steps/kustomization.yaml` by adding `deprecation-warn.yaml` to `resources`
- [x] 1.3 Add step reference `ai.steps.deprecation-warn` to `packages/domains/ai-agents/overlays/dev/agents-workflow.yaml` at weight `-100` as the first before phase
- [x] 1.4 Verify deprecation warning appears when running `kfg apply -k packages/domains/ai-agents/overlays/dev` — validate Bats test

## 2. Wizard Overlay

- [x] 2.1 Create `packages/domains/ai-agents/overlays/ai/kustomization.yaml` referencing framework base `../../../../framework/`, domain manifests `../../manifests/`, the wizard asset, and the wizard workflow
- [x] 2.2 Create `packages/domains/ai-agents/overlays/ai/ai-workflow.yaml` with CmdWorkflow `kfg.ai.workflow` supporting `ai.pi.cmd.main` and `ai.opencode.cmd.main`, with before phase steps: ensure-gitignore (-90), detect-agent (-70), scaffold per-agent (-65), ctx7 install per-agent (-55), wizard materialize per-agent (-45), and after phase: cleanup
- [x] 2.3 Create `packages/domains/ai-agents/overlays/ai/assets/prompts/wizard.yaml` with `kind: Assets`, `metadata.name: ai.prompts.wizard` containing the wizard skill prompt with interactive rules and domain building block catalog
- [x] 2.4 Verify `kfg build packages/domains/ai-agents/overlays/ai/` produces valid YAML with `kfg.ai.workflow`

## 3. CLI Command

- [x] 3.1 Create `src/cmd/kfg/ai.go` with cobra command `kfg ai`, reading `KFG_AI_AGENT` env var (default `pi`), constructing and executing `kfg run -k packages/domains/ai-agents/overlays/ai <agent> -- <args>` via subprocess
- [x] 3.2 Set up `--help` text documenting purpose, KFG_AI_AGENT env var, supported agents (pi, opencode), and usage examples
- [x] 3.3 Register command with `rootCmd.AddCommand(aiCmd)` in init()
- [x] 3.4 Verify `kfg ai --help` output includes env var documentation and usage examples

## 4. Testing

- [x] 4.1 Add Bats test verifying `kfg build packages/domains/ai-agents/overlays/ai/` produces valid YAML output
- [x] 4.2 Add Bats test verifying steps directory includes `deprecation-warn.yaml`
- [x] 4.3 Add Bats test verifying `kfg ai --help` output includes KFG_AI_AGENT documentation
- [x] 4.4 Run `make test-bats` to ensure no regressions

## 5. Documentation

- [x] 5.1 Update `docs/AGENTS.md` with new `kfg ai` command usage
- [x] 5.2 Add KFG_AI_AGENT to environment variable documentation
