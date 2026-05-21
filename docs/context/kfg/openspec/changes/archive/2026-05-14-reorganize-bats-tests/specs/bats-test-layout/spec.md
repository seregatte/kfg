## ADDED Requirements

### Requirement: Unified Bats test root
The repository MUST store all Bats tests under a single canonical root at `tests/bats/`.

#### Scenario: Discovering Bats tests
- **WHEN** a contributor searches for repository Bats coverage
- **THEN** all supported Bats test files SHALL reside under `tests/bats/`
- **AND** no supported Bats test root SHALL remain under `.manifests/tests/`

### Requirement: Mirrored manifest test layout
Manifest-resource Bats tests MUST mirror the resource paths they validate under `.manifests/base/` and `.manifests/overlay/`.

#### Scenario: Base resource test mapping
- **WHEN** a manifest resource exists at `.manifests/base/<path>/<name>.yaml`
- **THEN** its Bats coverage SHALL reside at `tests/bats/manifests/base/<path>/<name>.bats`

#### Scenario: Overlay resource test mapping
- **WHEN** a manifest resource exists at `.manifests/overlay/<overlay>/<name>.yaml`
- **THEN** its Bats coverage SHALL reside at `tests/bats/manifests/overlay/<overlay>/<name>.bats`

### Requirement: Shared Bats helpers by concern
The repository MUST provide shared Bats helpers under `tests/bats/helpers/` with concern-specific entrypoints.

#### Scenario: Loading common Bats helpers
- **WHEN** a Bats suite needs repository root or binary bootstrap logic
- **THEN** it SHALL load shared helper code from `tests/bats/helpers/`
- **AND** repository-wide helper behavior SHALL NOT be duplicated across multiple helper roots

#### Scenario: Loading manifest execution helpers
- **WHEN** a Bats suite executes a manifest resource such as a Step or overlay workflow
- **THEN** it SHALL use helper functions from the shared manifest helper module
- **AND** those helpers SHALL resolve manifest paths relative to the repository root rather than the suite directory

### Requirement: Separate workflow runtime coverage
Generic workflow and shell-runtime Bats tests MUST remain distinct from manifest-resource tests even though they share the same root.

#### Scenario: Organizing generic workflow tests
- **WHEN** a Bats test validates generic `CmdWorkflow` or `Step` runtime behavior using ad hoc fixtures
- **THEN** the test SHALL live under a non-manifest subtree within `tests/bats/`
- **AND** it SHALL NOT be placed in the mirrored manifest-resource tree unless it targets a checked-in manifest resource

### Requirement: Unified Bats execution target
Repository Bats entrypoints MUST run against the unified `tests/bats/` tree.

#### Scenario: Running the canonical Bats target
- **WHEN** repository Bats tests are invoked through the canonical test target
- **THEN** the target SHALL execute suites from `tests/bats/`
- **AND** contributors SHALL NOT need a separate manifest-specific Bats root to run supported shell tests
