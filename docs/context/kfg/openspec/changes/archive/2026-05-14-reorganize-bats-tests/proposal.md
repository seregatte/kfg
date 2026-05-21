## Why

The repository currently splits Bats coverage across `tests/bats/` and `.manifests/tests/`, which makes it harder to understand ownership, evolve helpers consistently, and see which manifest resources are covered by runtime tests. We need a single canonical Bats test tree that preserves the distinction between generic workflow runtime coverage and tests for concrete manifest resources.

## What Changes

- Consolidate all Bats tests under a single `tests/bats/` root.
- Reorganize manifest-oriented Bats tests so their paths mirror the resource paths under `.manifests/base/` and `.manifests/overlay/`.
- Split shared Bats helpers into a single helper structure under `tests/bats/helpers/` and remove duplicated helper logic.
- Keep generic engine/runtime Bats coverage separate from manifest-resource Bats coverage by using distinct subtrees under the shared root.
- Update test entrypoints so repository Bats execution no longer depends on a separate `.manifests/tests/` location.

## Capabilities

### New Capabilities
- `bats-test-layout`: Defines the canonical repository layout for all Bats tests, including mirrored manifest test paths and shared helper organization.

### Modified Capabilities
- `project-structure`: Clarify the canonical location and organization rules for Bats tests and remove the split between top-level and manifest-local Bats suites.

## Impact

- Affected paths include `tests/bats/`, `.manifests/tests/`, and `Makefile` Bats targets.
- Affected shell test helpers include both existing `test_helper.bash` implementations.
- User-facing impact is limited to contributor and CI test invocation, with a simpler and more predictable shell test layout.
- Manifest model behavior does not change, but the repository will document and enforce a new mapping between manifest resources and their Bats coverage.
