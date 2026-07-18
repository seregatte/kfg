## 1. Establish the unified Bats layout

- [x] 1.1 Create the canonical `tests/bats/` subtrees for `helpers/`, `manifests/base/`, `manifests/overlay/`, and workflow/runtime suites.
- [x] 1.2 Decide whether vendored Bats helper libraries stay in place or move under the new helper tree, and update references consistently.
- [x] 1.3 Update repository path/bootstrap helpers so all Bats suites resolve `PROJECT_ROOT` and `KFG_BIN` through shared helper code.

## 2. Migrate manifest-resource Bats suites

- [x] 2.1 Move each existing `.manifests/tests/*.bats` file into a mirrored path under `tests/bats/manifests/` based on the resource it validates.
- [x] 2.2 Split manifest-specific helper functions into shared modules under `tests/bats/helpers/` and remove duplicated helper roots.
- [x] 2.3 Update manifest Bats suites to load the new helpers and resolve manifest resource paths relative to the repository root.

## 3. Reorganize generic Bats suites

- [x] 3.1 Move workflow and runtime-oriented Bats tests into dedicated subtrees under `tests/bats/` without mixing them into the mirrored manifest tree.
- [x] 3.2 Keep CLI-only Bats suites under their own subtree and update any helper loads affected by the directory move.
- [x] 3.3 Remove or retire the old `.manifests/tests/` root once all supported suites have been migrated.

## 4. Update test entrypoints and documentation

- [x] 4.1 Update `Makefile` Bats targets so the canonical entrypoint runs `bats tests/bats`.
- [x] 4.2 Decide whether `test-manifests` remains as a compatibility alias and implement that decision consistently.
- [x] 4.3 Update repository documentation or contributor guidance that points to `.manifests/tests/` or old helper locations.

## 5. Validate the reorganized suite

- [x] 5.1 Run the unified Bats suite and fix path, helper, or fixture regressions introduced by the move.
- [x] 5.2 Run targeted workflow/runtime Bats checks to confirm shell output and `CmdWorkflow`/`Step` behavior still validate correctly.
- [x] 5.3 Run targeted manifest-resource Bats checks to confirm mirrored tests still execute the intended `.manifests/base/` and `.manifests/overlay/` resources.
- [x] 5.4 Run the relevant Go unit and integration tests that guard Bats-adjacent behavior to ensure the repository remains deterministic after the layout change.
