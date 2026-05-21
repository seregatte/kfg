## Why

The current `kfg.cleanup` step still relies on a `PATHS` environment variable even though artifact-producing steps already record their outputs in `KFG_ARTIFACTS`. This creates a mismatch between generated shell behavior and cleanup behavior, so workflow cleanup can miss files that were actually produced.

## What Changes

- Change `kfg.cleanup` to remove only paths recorded in `KFG_ARTIFACTS`.
- Remove the step-level `PATHS` environment contract from `kfg.cleanup`.
- Add tests that verify cleanup removes recorded artifact files and directories and does nothing when no artifacts were recorded.
- Clarify spec coverage for artifact tracking and cleanup behavior.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `internal-shell-gen`: clarify that artifact-producing steps append cleanup targets to the shared `KFG_ARTIFACTS` array for later shell-session cleanup.
- `dev-workflow`: tighten cleanup behavior so workflow cleanup operates on tracked artifacts instead of an external `PATHS` contract.

## Impact

- Affected manifests: `.manifests/base/steps/cleanup.yaml`
- Affected tests: Bats tests for manifest steps cleanup behavior
- Affected specs: `internal-shell-gen`, `dev-workflow`
- User-facing shell UX: cleanup becomes consistent with the generated artifact tracking contract
