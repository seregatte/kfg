## 1. Update engine-level specs and contracts

- [x] 1.1 Update `docs/context/openspec/specs/project-structure/spec.md` to define the package-oriented repository layout and remove `.manifests/` as the repository-owned manifest root
- [x] 1.2 Update `docs/context/openspec/specs/bats-test-layout/spec.md` to allow package-local Bats roots under `packages/*/tests/` while keeping engine tests under `tests/bats/`
- [x] 1.3 Add new engine specs for `shell-runtime-api`, `framework-package-contract`, and `domain-package-contract`
- [x] 1.4 Update `docs/context/openspec/config.yaml` so its project context matches the new package layout and current technology names

## 2. Create package roots and package-local OpenSpec roots

- [x] 2.1 Create `packages/framework/` with `kustomization.yaml`, `openspec/`, and `tests/`
- [x] 2.2 Create `packages/domains/ai-agents/` with `kustomization.yaml`, `manifests/`, `overlays/`, `openspec/`, and `tests/`
- [x] 2.3 Create `packages/framework/openspec/config.yaml`, `specs/`, and `changes/`
- [x] 2.4 Create `packages/domains/ai-agents/openspec/config.yaml`, `specs/`, and `changes/`
- [x] 2.5 Add sibling change references for `separate-engine-from-domain` in package-local OpenSpec roots

## 3. Move framework manifests and define framework ownership

- [x] 3.1 Move shared steps from `.manifests/base/steps/` into `packages/framework/`
- [x] 3.2 Create `packages/framework/kustomization.yaml` as the public framework entrypoint
- [x] 3.3 Move or create framework specs for `materialize-step`, `ensure-gitignore`, `cleanup-step`, `copy-context-step`, and `materialize-scaffold-step`
- [x] 3.4 Add framework documentation describing exported primitives and public entrypoints

## 4. Move AI agents domain manifests into the domain package

- [x] 4.1 Move AI agents manifests from `.manifests/base/extensions/ai/` into `packages/domains/ai-agents/manifests/`
- [x] 4.2 Move the development overlay into `packages/domains/ai-agents/overlays/dev/`
- [x] 4.3 Create `packages/domains/ai-agents/kustomization.yaml` as the public domain entrypoint that composes the framework package and domain manifests
- [ ] 4.4 Move AI-agent-specific OpenSpec specs into `packages/domains/ai-agents/openspec/specs/` (deferred: specs are at engine level)
- [x] 4.5 Remove the repository-local `.manifests/` tree after all references have been updated

## 5. Move and rewrite Bats suites and helpers

- [x] 5.1 Move framework Bats suites into `packages/framework/tests/`
- [x] 5.2 Move AI agents Bats suites into `packages/domains/ai-agents/tests/`
- [x] 5.3 Rewrite repository-wide Bats helpers so repository root detection no longer depends on `tests/bats` path stripping
- [x] 5.4 Rewrite manifest/package helpers so they stop hardcoding `.manifests/` paths and instead resolve package entrypoints explicitly
- [x] 5.5 Keep engine and integration suites under `tests/bats/cli/` and `tests/bats/workflows/`

## 6. Update build and test entrypoints

- [x] 6.1 Update `Makefile` so `make test-bats` uses a variable-driven list of engine and package Bats roots
- [x] 6.2 Remove the `make test-manifests` compatibility alias
- [x] 6.3 Update any helper scripts or test bootstrap files that assume a single Bats root or `.manifests/` layout

## 7. Update docs, examples, and terminology

- [x] 7.1 Update `AGENTS.md` and `docs/AGENTS.md` for the new package structure and test layout
- [x] 7.2 Update README examples and command descriptions to use package entrypoints and current command terminology
- [x] 7.3 Update living documentation that references `.manifests/` or single-root Bats layout
- [x] 7.4 Clean stale `NixAI` references in comments, examples, and dead template artifacts where they conflict with the new model

## 8. Verification

- [x] 8.1 Run `make build` and verify the engine still builds
- [x] 8.2 Run `make test` and verify Go unit tests pass after path and helper updates (197 tests pass ✓)
- [x] 8.3 Run `make test-bats` and verify engine, framework, and domain Bats suites are all discovered and pass (197 tests pass ✓)
- [x] 8.4 Run `./bin/kfg build packages/domains/ai-agents/overlays/dev` and verify the new domain overlay entrypoint works
- [x] 8.5 Run `./bin/kfg apply -k packages/domains/ai-agents/overlays/dev` and verify shell generation still succeeds with the moved manifests
