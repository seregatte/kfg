## MODIFIED Requirements

### Requirement: Unified Bats execution target
Repository Bats entrypoints MUST run against the canonical engine test root and any package-owned test roots.

#### Scenario: Running the canonical Bats target
- **WHEN** repository Bats tests are invoked through `make test-bats`
- **THEN** the target SHALL execute suites from `tests/bats/`
- **AND** the target SHALL also execute suites from `packages/framework/tests/` and any `packages/domains/*/tests/` directories

### Requirement: Engine Bats test root
Engine and integration Bats suites MUST remain under `tests/bats/`.

#### Scenario: Discovering engine Bats tests
- **WHEN** a contributor searches for engine or integration Bats coverage
- **THEN** supported engine Bats files SHALL reside under `tests/bats/`
- **AND** repository-wide helpers SHALL continue to live under `tests/bats/helpers/`

### Requirement: Package-owned Bats roots
Framework and domain packages MAY own package-local Bats suites under their package directories.

#### Scenario: Framework package test mapping
- **WHEN** a framework manifest primitive is validated by Bats
- **THEN** its suite SHALL reside under `packages/framework/tests/`

#### Scenario: Domain package test mapping
- **WHEN** a domain manifest or package-local overlay is validated by Bats
- **THEN** its suite SHALL reside under `packages/domains/<domain>/tests/`

### Requirement: Shared Bats helpers by concern
The repository MUST provide repository-wide Bats helpers under `tests/bats/helpers/`, and package-local helpers MAY extend them without depending on legacy manifest roots.

#### Scenario: Loading common Bats helpers
- **WHEN** a Bats suite needs repository root or binary bootstrap logic
- **THEN** it SHALL load shared helper code from `tests/bats/helpers/`

#### Scenario: Loading package-aware manifest helpers
- **WHEN** a package-local Bats suite executes a manifest resource
- **THEN** it MAY use a package-local helper module
- **AND** that helper SHALL resolve paths relative to the package or explicit repository root configuration rather than `.manifests/`
