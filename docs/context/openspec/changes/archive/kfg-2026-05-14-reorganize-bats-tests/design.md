## Context

The repository currently maintains Bats coverage in two unrelated locations: `tests/bats/` for CLI and runtime-oriented tests, and `.manifests/tests/` for tests that exercise concrete manifest resources. Each subtree has its own helper conventions and its own assumptions about relative paths, which increases maintenance cost and makes it difficult to reason about ownership, discoverability, and invocation.

The desired end state is a single Bats root under `tests/bats/` with two clear testing modes inside it: generic workflow/runtime tests and manifest-resource tests. Manifest-resource tests must preserve a path structure that mirrors the resource paths under `.manifests/base/` and `.manifests/overlay/`, so contributors can infer the tested resource directly from the test path.

## Goals / Non-Goals

**Goals:**
- Establish `tests/bats/` as the only repository root for Bats tests.
- Mirror `.manifests/base/` and `.manifests/overlay/` under `tests/bats/manifests/` for resource-oriented Bats coverage.
- Replace duplicated helper implementations with a shared helper layout under `tests/bats/helpers/`.
- Keep generic engine/runtime coverage distinct from manifest-resource coverage while preserving a single Bats entrypoint.
- Update repository test targets to run the unified Bats tree without relying on `.manifests/tests/`.

**Non-Goals:**
- Changing manifest semantics, shell generation logic, or command behavior.
- Rewriting non-Bats unit or integration tests.
- Expanding Bats coverage beyond the workflow/step and manifest-resource concerns already discussed.

## Decisions

### Decision: Use a single canonical Bats root

All Bats suites will live under `tests/bats/`. This gives contributors one place to discover shell tests and lets `make test-bats` invoke a single tree.

Alternative considered:
- Keep `tests/bats/` and `.manifests/tests/` as separate roots and document the split.
- Rejected because the split is the source of the discoverability and helper-duplication problems.

### Decision: Mirror manifest resource paths under `tests/bats/manifests/`

Manifest-resource tests will be organized by the tested resource path rather than by test type. For example:

- `.manifests/base/steps/convert.yaml` -> `tests/bats/manifests/base/steps/convert.bats`
- `.manifests/base/agents/steps/settings.yaml` -> `tests/bats/manifests/base/agents/steps/settings.bats`
- `.manifests/overlay/dev/agents-workflow.yaml` -> `tests/bats/manifests/overlay/dev/agents-workflow.bats`

This preserves resource locality while still moving all Bats files into the unified tree.

Alternative considered:
- Group manifest tests by behavior such as `steps/`, `workflow/`, or `converters/`.
- Rejected because it loses the one-to-one mapping between manifest resources and test files.

### Decision: Split helpers by concern under a shared helpers directory

Shared path/bootstrap logic belongs in `tests/bats/helpers/common.bash`. Manifest execution helpers belong in `tests/bats/helpers/manifests.bash`. Workflow-shell helpers belong in `tests/bats/helpers/workflow_runtime.bash`.

This keeps helpers reusable without forcing every suite to load unrelated functions.

Alternative considered:
- Collapse all helpers into a single `test_helper.bash`.
- Rejected because the current duplication already shows that one large helper becomes hard to evolve cleanly.

### Decision: Keep generic runtime tests outside the mirrored manifest tree

Tests that validate generic `kfg apply` runtime behavior without targeting a specific checked-in manifest resource will remain under a separate subtree such as `tests/bats/workflows/`. This avoids conflating engine contracts with tests for concrete resources in `.manifests/`.

Alternative considered:
- Place all workflow tests under `tests/bats/manifests/overlay/dev/`.
- Rejected because many runtime tests are fixture-based and do not belong to a repository manifest resource.

### Decision: Remove dedicated `.manifests/tests` entrypoints

The repository will no longer maintain a separate Bats root or Make target for `.manifests/tests/`. Existing tests will move into `tests/bats/manifests/...`, and the unified Bats target will become the supported execution path.

Alternative considered:
- Keep `test-manifests` as a separate target pointing at the new subtree.
- Acceptable as a temporary compatibility alias, but not as the primary structure.

## Risks / Trade-offs

- [Relative-path helper breakage] -> Centralize root resolution in `helpers/common.bash` and update all `load` statements to explicit relative helper paths.
- [Large file move obscures review] -> Move files with a predictable resource-to-test mapping and keep content changes minimal during relocation.
- [Contributor confusion during transition] -> Update the Make targets and any repository documentation that references `.manifests/tests/`.
- [Overcoupling helpers to one test style] -> Keep helpers split by concern so manifest and workflow suites can evolve independently.

## Migration Plan

1. Create the new helper layout under `tests/bats/helpers/`.
2. Move `.manifests/tests/*.bats` into mirrored paths under `tests/bats/manifests/...`.
3. Update helper loads and root/path resolution for all moved tests.
4. Reorganize workflow/runtime Bats tests under dedicated subtrees within `tests/bats/`.
5. Update `Makefile` targets so the canonical execution path is `bats tests/bats`.
6. Run the unified Bats suite and fix any path or helper regressions.

Rollback is straightforward: the change is limited to repository layout and test harness paths, so the repository can restore the prior file layout if the unified tree proves unstable.

## Open Questions

- Whether `test-manifests` should remain as a compatibility alias that delegates to `tests/bats/manifests/`.
- Whether vendored Bats helper libraries should stay under `tests/bats/test_helper/` or move under `tests/bats/helpers/vendor/` as part of the same change.
