## 1. Cleanup Step Contract

- [x] 1.1 Update `.manifests/base/steps/cleanup.yaml` to remove only paths recorded in `KFG_ARTIFACTS`
- [x] 1.2 Remove the `PATHS` environment contract from the shared `kfg.cleanup` step
- [x] 1.3 Audit all `kfg.cleanup` call sites in `.manifests/overlay/dev/agents-workflow.yaml` and decide which ones must migrate away from `PATHS`

## 2. Workflow Migration

- [x] 2.1 Update per-agent cleanup usage so tracked artifact cleanup still removes the intended generated files
- [x] 2.2 Update after-phase cleanup usage so it no longer depends on `PATHS`
- [x] 2.3 Verify no remaining workflow behavior depends on the removed `PATHS` contract

## 3. Validation

- [x] 3.1 Add Bats coverage for `kfg.cleanup` removing tracked files, tracked directories, and no-op behavior when `KFG_ARTIFACTS` is empty
- [x] 3.2 Add or update integration coverage for workflow cleanup behavior if manifest call sites change materially
- [x] 3.3 Run the relevant Bats and Go test suites to confirm cleanup behavior remains deterministic

## 4. Documentation

- [x] 4.1 Align manifest or developer documentation with the artifact-tracked cleanup contract if implementation changes user-visible workflow behavior
