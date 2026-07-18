## ADDED Requirements

### Requirement: Domain package public entrypoint
Each domain package SHALL expose a root `kustomization.yaml` as its public package entrypoint.

#### Scenario: Referencing a domain package
- **WHEN** a consumer or overlay references a domain package
- **THEN** it SHALL reference `packages/domains/<domain>/kustomization.yaml`
- **AND** it SHALL NOT depend on internal manifest subpaths as the public API

### Requirement: Domain package overlays layout
Domain-specific overlays SHALL live under an `overlays/` directory in the owning domain package.

#### Scenario: Storing domain overlays
- **WHEN** a domain package defines environment-specific overlays
- **THEN** they SHALL reside under `packages/domains/<domain>/overlays/`

### Requirement: Domain package OpenSpec root
Each domain package SHALL own a package-local OpenSpec root.

#### Scenario: Storing domain specs and changes
- **WHEN** domain-specific behavior is specified or changed
- **THEN** the domain package SHALL use `packages/domains/<domain>/openspec/`
- **AND** that root SHALL contain `config.yaml`, `specs/`, and `changes/`

### Requirement: Domain package test root
Each domain package SHALL own its package-local Bats suites.

#### Scenario: Running domain Bats coverage
- **WHEN** domain manifests or overlays are validated by Bats
- **THEN** their suites SHALL reside under `packages/domains/<domain>/tests/`

### Requirement: Cross-layer changes use sibling change names
Cross-layer initiatives SHALL be represented by sibling OpenSpec changes with the same slug in each affected layer.

#### Scenario: Coordinating a cross-layer change
- **WHEN** a change affects engine, framework, and domain responsibilities
- **THEN** each affected OpenSpec root SHALL contain a change with the same slug
- **AND** each change SHALL be authoritative only for its own layer
