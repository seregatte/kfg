## 1. Manifest Cleanup

- [x] 1.1 Remove `ai.claude.cmd.main` from `packages/domains/ai-agents/manifests/cmds/agents.yaml`
- [x] 1.2 Remove `- claude` from `packages/domains/ai-agents/manifests/agents/kustomization.yaml`
- [x] 1.3 Delete entire `packages/domains/ai-agents/manifests/agents/claude/` directory
- [x] 1.4 Remove claude branch from `packages/domains/ai-agents/manifests/ctx7/steps/install.yaml`
- [x] 1.5 Verify `kfg build packages/domains/ai-agents/` succeeds without claude

## 2. Engine Cleanup

- [x] 2.1 Remove `pkgs.claude-code` from `flake.nix` kfg-bundle
- [x] 2.2 Remove `/.claude/` and `/CLAUDE.md` from `.gitignore`
- [x] 2.3 Remove claude from `packages/domains/ai-agents/overlays/dev/agents-workflow.yaml`

## 3. Documentation

- [x] 3.1 Remove claude references from `packages/domains/ai-agents/manifests/README.md`
- [x] 3.2 Remove claude references from `packages/domains/ai-agents/manifests/agents/README.md`

## 4. Wizard Prompt Update

- [x] 4.1 Remove `claude` from agent catalog in `packages/domains/ai-agents/overlays/ai/assets/prompts/kfg-wizard.yaml`
- [x] 4.2 Update example conversations in wizard prompt

## 5. Test Updates

- [x] 5.1 Remove claude test cases from `packages/framework/tests/materialize-scaffold.bats`
- [x] 5.2 Remove claude test cases from `packages/domains/ai-agents/tests/base/extensions/ai-agents-naming.bats`
- [x] 5.3 Remove claude test cases from `packages/domains/ai-agents/tests/base/extensions/ai-agents-structure.bats`
- [x] 5.4 Remove claude references from `packages/domains/ai-agents/tests/base/extensions/ai-agents-kustomize.bats`
- [x] 5.5 Remove claude fixtures from `src/internal/generate/generate_test.go`
- [x] 5.6 Remove claude fixtures from `src/cmd/kfg/run_test.go`
- [x] 5.7 Remove claude references from `src/internal/kustomize/integration_test.go`
- [x] 5.8 Update comment references in `src/internal/manifest/types.go`
- [x] 5.9 Update comment references in `src/internal/resolve/resolve_test.go`
- [x] 5.10 Remove claude reference from `src/cmd/kfg/apply.go` help text
- [x] 5.11 Remove claude reference from `src/cmd/kfg/log.go` help text
- [x] 5.12 Remove claude reference from `tests/bats/cli/run-command-examples.bats`
- [x] 5.13 Remove claude references from `tests/bats/workflows/stepref-output-addressing.bats`

## 6. Spec & Doc Updates

- [x] 6.1 Update `docs/context/openspec/specs/` claude references to opencode
- [x] 6.2 Update `docs/context/openspec/config.yaml` claude reference
- [x] 6.3 Update `docs/manifest-model.md` claude references

## 7. Validation

- [x] 7.1 Run `make test` — all Go tests pass
- [x] 7.2 Run `make build` — binary builds successfully
