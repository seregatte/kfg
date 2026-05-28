## MODIFIED Requirements

### Requirement: Project Root Structure

The project MUST maintain a clear separation between engine implementation, package-owned manifests, documentation, and test assets.

#### Scenario: Core configuration files
- **WHEN** the project is initialized
- **THEN** `flake.nix` SHALL define Nix packaging and development shells
- **AND** `flake.lock` SHALL pin Nix inputs
- **AND** configuration files SHALL reside at the project root

#### Scenario: Documentation location
- **WHEN** documentation is organized
- **THEN** `docs/` SHALL contain repository-wide documentation
- **AND** `docs/AGENTS.md` SHALL provide agent operating context
- **AND** `README.md` SHALL reside at the project root

#### Scenario: Test asset location
- **WHEN** repository shell tests are organized
- **THEN** engine and integration Bats tests SHALL reside under `tests/bats/`
- **AND** shared repository-wide Bats helpers SHALL reside under `tests/bats/helpers/`
- **AND** package-specific Bats suites MAY reside under `packages/*/tests/`

### Requirement: OpenSpec Structure

OpenSpec artifacts MUST follow a consistent directory structure for each ownership boundary.

#### Scenario: Engine OpenSpec root location
- **WHEN** engine, CLI, or project-level specs and changes are stored
- **THEN** `docs/context/openspec/` SHALL be the OpenSpec root for those concerns
- **AND** `docs/context/openspec/config.yaml` SHALL define schema and context

#### Scenario: Package OpenSpec root location
- **WHEN** a package owns framework or domain-specific behavior
- **THEN** the package SHALL provide its own OpenSpec root at `packages/<package>/openspec/`
- **AND** that root SHALL contain its own `config.yaml`, `specs/`, and `changes/` directories

### Requirement: Directory Tree

The directory tree MUST follow the canonical layered structure.

#### Scenario: Canonical structure
- **WHEN** the repository is organized
- **THEN** the directory tree SHALL include `src/` for engine code, `docs/context/openspec/` for engine and CLI specs, `packages/framework/` for shared framework manifests, and `packages/domains/*/` for domain packages
- **AND** engine and integration Bats tests SHALL remain under `tests/bats/`
- **AND** package-specific tests SHALL live under the owning package's `tests/` directory

### Requirement: Manifest Storage

Repository-owned manifests MUST be stored in package directories with explicit public entrypoints.

#### Scenario: Framework package entrypoint
- **WHEN** shared manifest primitives are exposed
- **THEN** `packages/framework/kustomization.yaml` SHALL be the public framework entrypoint

#### Scenario: Domain package entrypoint
- **WHEN** a domain package is exposed
- **THEN** `packages/domains/<domain>/kustomization.yaml` SHALL be the public domain entrypoint
- **AND** domain-specific overlays SHALL live under `packages/domains/<domain>/overlays/`

### Requirement: No Mixed Locations

Engine, framework, domain, and test assets MUST NOT depend on legacy repository-local manifest layout.

#### Scenario: Legacy manifest root removed
- **WHEN** repository-owned manifests are stored
- **THEN** they MUST NOT be stored under `.manifests/`
- **AND** consumers SHALL use package public entrypoints instead of internal legacy paths
