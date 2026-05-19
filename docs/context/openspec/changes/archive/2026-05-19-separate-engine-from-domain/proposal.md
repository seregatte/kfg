## Why

kfg is already a generic manifest engine in production code, but the repository structure still mixes engine concerns with shared framework manifests and the AI agents domain. Separating these layers now reduces future extraction cost, makes ownership clearer, and avoids repeating path-coupling regressions in downstream repositories.

## What Changes

- Reorganize the repository around three explicit layers: engine, shared framework package, and domain packages.
- Move shared manifest primitives into `packages/framework/` with a public `kustomization.yaml` entrypoint.
- Move AI agents manifests into `packages/domains/ai-agents/` with a public `kustomization.yaml` entrypoint and `overlays/` layout.
- Keep engine, CLI, and project-level OpenSpec specs in `docs/context/openspec/`.
- Add package-local OpenSpec roots for `packages/framework/` and `packages/domains/ai-agents/`.
- Move package-specific Bats suites into `packages/<name>/tests/` while keeping engine and integration suites under `tests/bats/`.
- **BREAKING** Remove the repository-local `.manifests/` layout in favor of package entrypoints under `packages/`.
- **BREAKING** Remove the `make test-manifests` alias and redefine `make test-bats` to run engine and package Bats suites.
- Clean up stale NixAI and outdated terminology where it conflicts with the new repository model.

## Capabilities

### New Capabilities
- `shell-runtime-api`: define the generated shell runtime contract that framework steps consume.
- `framework-package-contract`: define the public entrypoints, exported primitives, OpenSpec root, and tests for the shared framework package.
- `domain-package-contract`: define the public entrypoints, overlays, OpenSpec root, and tests for domain packages.

### Modified Capabilities
- `project-structure`: change the canonical repository layout to use `packages/framework/` and `packages/domains/*/` instead of `.manifests/`, and allow package-local OpenSpec and test roots.
- `bats-test-layout`: change Bats layout rules so engine tests stay under `tests/bats/` while package-specific suites live under `packages/*/tests/` and are discovered by the canonical `make test-bats` target.

## Impact

- Affected manifests: current `.manifests/base/` and `.manifests/overlay/dev/` content moves into `packages/framework/` and `packages/domains/ai-agents/`.
- Affected tests: manifest Bats suites and helpers must be relocated and rewritten to stop assuming `.manifests/` and `tests/bats/`-only roots.
- Affected documentation: AGENTS docs, README examples, project structure docs, and path references must be updated.
- Affected build/test entrypoints: `make test-bats` changes semantics to run multiple roots; `make test-manifests` is removed.
- Affected downstream consumers: any repository depending on internal `.manifests/` paths must migrate to package public entrypoints.
