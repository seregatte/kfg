## Context

The engine implementation in `src/` is already domain-agnostic: production code does not hard-code AI-agent behavior and the manifest model is generic. The repository structure does not reflect that reality. Shared framework steps, AI agents manifests, engine specs, and package-specific shell tests are still organized as if they belong to one repository-local manifest tree.

That mismatch creates three problems:

1. Ownership is unclear. It is hard to tell which files define the engine contract, which files implement reusable framework behavior, and which files belong to a specific domain.
2. Tests and docs encode old path assumptions. Bats helpers assume `.manifests/` and `tests/bats/`-only roots, while multiple docs and specs still require the old layout.
3. External consumers can depend on internal paths. Prior downstream breakage already showed that relative references to internal manifest directories are not a stable public API.

The new design keeps engine and CLI specs in `docs/context/openspec/`, but introduces package roots for framework and domains so each layer can own its manifests, tests, and specs without relying on a global `.manifests/` tree.

## Goals / Non-Goals

**Goals:**
- Make the repository structure match the actual architectural layers: engine, framework, domain.
- Define public package entrypoints that downstream consumers can depend on instead of internal paths.
- Keep engine and CLI OpenSpec artifacts in `docs/context/openspec/`.
- Give framework and domain packages their own OpenSpec roots, test roots, and public kustomize entrypoints.
- Preserve the engine runtime API that framework steps already consume.
- Make `make test-bats` continue to be the canonical Bats entrypoint while discovering tests across engine and packages.
- Eliminate stale `.manifests/`, `NixAI`, and outdated agent-centric terminology where it conflicts with the new model.

**Non-Goals:**
- Splitting the monorepo into multiple repositories now.
- Changing manifest resource kinds or shell generation semantics.
- Introducing backward-compatibility aliases for the old `.manifests/` layout.
- Changing engine runtime behavior beyond what is necessary to support the new repository structure.

## Decisions

### 1. Repository structure becomes package-oriented

The repository will use three explicit layers:

- `src/` for the engine implementation.
- `packages/framework/` for shared manifest primitives.
- `packages/domains/ai-agents/` for the AI agents domain.

The old `.manifests/` tree will be removed. Shared steps move into `packages/framework/`. AI agents manifests move into `packages/domains/ai-agents/manifests/` with package-local overlays under `packages/domains/ai-agents/overlays/`.

Rationale:
- Matches actual ownership boundaries already present in the code and naming conventions.
- Prepares future repo extraction without requiring a second structural rewrite.

Alternative considered:
- Keep `.manifests/` and only reorganize specs. Rejected because it preserves the main source of path coupling and keeps framework/domain ownership implicit.

### 2. Every package gets a public kustomization entrypoint

The framework package MUST expose `packages/framework/kustomization.yaml`.
The AI agents domain MUST expose `packages/domains/ai-agents/kustomization.yaml`.

The domain root entrypoint composes the framework package and its own manifests. Domain overlays live under `packages/domains/ai-agents/overlays/` and can reference the domain package root rather than internal framework paths.

Rationale:
- Prevents downstream consumers from depending on internal package layout.
- Makes the public API of each package explicit.

Alternative considered:
- Let consumers reference internal subdirectories directly. Rejected because previous downstream breakage came from depending on private internal paths.

### 3. OpenSpec ownership is split by layer, not centralized for cross-layer work

OpenSpec roots become:

- `docs/context/openspec/` for engine, CLI, and project concerns.
- `packages/framework/openspec/` for shared framework behavior.
- `packages/domains/ai-agents/openspec/` for AI agents domain behavior.

There is no central cross-layer change root. Cross-layer initiatives use sibling changes with the same slug in each affected OpenSpec root. Each change is authoritative only for its own layer.

Rationale:
- Keeps responsibility boundaries explicit.
- Avoids a central change record that blurs ownership back together.

Alternative considered:
- Keep all changes in `docs/context/openspec/changes/`. Rejected because the user explicitly wants ownership split by responsibility.

### 4. Engine and package tests have separate homes

Engine and integration Bats suites stay under `tests/bats/`.
Package-specific suites move to `packages/framework/tests/` and `packages/domains/ai-agents/tests/`.

`make test-bats` remains the canonical target, but it runs Bats against multiple roots via a variable-driven directory list.

Rationale:
- Tests follow content ownership.
- Engine tests stay easy to discover in the top-level test tree.
- Package extraction later does not require moving package-local tests again.

Alternative considered:
- Keep all Bats suites under one root. Rejected because it conflicts with package ownership and future extraction.

### 5. Bats helpers become package-aware instead of path-hack based

The current helper strategy assumes every suite lives under `tests/bats/` and every manifest lives under `.manifests/`. That must be replaced.

The new helper model will:

- Keep repository-wide bootstrap helpers under `tests/bats/helpers/`.
- Resolve repository root without relying on `/tests/bats` path stripping.
- Let package-local helpers compute paths relative to package roots or explicit environment variables.
- Stop hardcoding `.manifests/` in helper functions.

Rationale:
- The new layout otherwise breaks immediately.
- This also makes future package extraction simpler because helpers stop encoding monorepo-specific legacy paths.

### 6. Framework behavior gets its own explicit contract

The framework package owns reusable steps such as:

- `kfg.cleanup`
- `kfg.materialize`
- `kfg.materialize-scaffold`
- `kfg.ensure-gitignore`
- `kfg.copy-context`

Framework OpenSpec must document the functional behavior of these shared steps. `ensure-gitignore` belongs to the framework layer, not the AI agents domain. Additional framework specs should be added for steps that currently lack explicit coverage.

Rationale:
- The framework is a first-class package, not just moved files.
- Package extraction later requires stable functional documentation for exported primitives.

### 7. The engine runtime shell API remains the boundary contract

The engine continues to generate runtime helpers and metadata consumed by framework steps:

- `__kfg_add_artifact()` and `KFG_ARTIFACTS`
- `_kfg.log.*()`
- `__kfg_build_result()` and `KFG_BUILD_RESULT_FILE`
- `__kfg_when_*()`
- `__kfg_ctx_reset()` and output helpers
- `KFG_SESSION_ID`, `KFG_WORKFLOW_NAME`, `KFG_KUSTOMIZATION_NAME`, `KFG_SHELL`

This change does not move those APIs. It formalizes them as the engine-to-framework contract.

Rationale:
- The runtime boundary is already clean in the codebase.
- Preserving that contract allows repository structure to change without changing shell semantics.

## Risks / Trade-offs

- **Multiple OpenSpec roots increase coordination cost** → Use identical slugs for sibling changes and require each proposal to list related changes in other roots.
- **Helper rewrites may break Bats suites in subtle ways** → Update structure and helper specs first, then migrate helpers, then move tests.
- **Downstream users may still depend on old internal paths** → Document the new public package entrypoints and do not preserve `.manifests/` compatibility aliases.
- **The big-bang move can create noisy diffs** → Sequence the work in clear phases and verify after each structural section even if it ships as one change.
- **Framework ownership may remain incomplete if only some steps get specs** → Add missing framework specs as part of this change rather than leaving ownership implicit.

## Migration Plan

1. Update engine-level specs that currently contradict the target structure, especially `project-structure` and `bats-test-layout`.
2. Add new engine capability specs for the shell runtime API, framework package contract, and domain package contract.
3. Create package-local OpenSpec roots for `packages/framework/` and `packages/domains/ai-agents/`.
4. Move shared framework manifests into `packages/framework/` and establish its public `kustomization.yaml`.
5. Move AI agents manifests into `packages/domains/ai-agents/` and create the package root `kustomization.yaml` plus `overlays/` structure.
6. Move framework and domain Bats suites into package-local `tests/` directories.
7. Rewrite helpers and Makefile targets so `make test-bats` discovers engine and package suites.
8. Update documentation, examples, AGENTS docs, and lingering terminology/path drift.

Rollback strategy:
- Because this is a structural change with no backward-compatibility requirement, rollback is simply reverting the change before downstream migration begins.

## Open Questions

- Whether package-local test helpers should share one implementation file via `tests/bats/helpers/` sourcing, or whether each package should own a thin wrapper helper tailored to its package root.
- Whether the framework package should gain its own README-only public documentation now, or whether that should be deferred to a follow-up after the package split lands.
