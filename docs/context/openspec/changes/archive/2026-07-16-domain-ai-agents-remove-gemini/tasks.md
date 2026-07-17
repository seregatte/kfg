## 1. Manifest Cleanup

- [x] 1.1 Remove `ai.gemini.cmd.main` from `packages/domains/ai-agents/manifests/cmds/agents.yaml`
- [x] 1.2 Remove `- gemini` from `packages/domains/ai-agents/manifests/agents/kustomization.yaml`
- [x] 1.3 Delete entire `packages/domains/ai-agents/manifests/agents/gemini/` directory
- [x] 1.4 Remove gemini branch from `packages/domains/ai-agents/manifests/ctx7/steps/install.yaml`
- [x] 1.5 Verify `kfg build packages/domains/ai-agents/` succeeds without gemini

## 2. Engine Cleanup

- [x] 2.1 Remove `pkgs.gemini-cli-bin` from `flake.nix` kfg-bundle
- [x] 2.2 Remove `/.gemini/` from `.gitignore`
- [x] 2.3 Remove gemini references from `packages/domains/ai-agents/overlays/dev/agents-workflow.yaml`

## 3. Documentation

- [x] 3.1 Remove gemini references from `packages/domains/ai-agents/manifests/README.md`
- [x] 3.2 Remove gemini references from `packages/domains/ai-agents/manifests/agents/README.md`

## 4. Test Updates

- [x] 4.1 Remove gemini test cases from `packages/framework/tests/materialize-scaffold.bats`
- [x] 4.2 Remove gemini test cases from `packages/domains/ai-agents/tests/base/extensions/ai-agents-naming.bats`
- [x] 4.3 Remove gemini test cases from `packages/domains/ai-agents/tests/base/extensions/ai-agents-structure.bats`
- [x] 4.4 Remove gemini references from `packages/domains/ai-agents/tests/base/extensions/ai-agents-kustomize.bats`
- [x] 4.5 Remove gemini fixtures from `src/internal/generate/generate_test.go`
- [x] 4.6 Remove gemini fixtures from `src/cmd/kfg/run_test.go`
- [x] 4.7 Update comment references in `src/internal/manifest/types.go`
- [x] 4.8 Update comment references in `src/internal/resolve/resolve_test.go`

## 5. Validation

- [ ] 5.1 Run `make test` — all Go tests pass
- [ ] 5.2 Run `make test-bats` — all Bats tests pass
- [ ] 5.3 Run `make build` — binary builds successfully
