# Domain Package Contract Specification

## Purpose

Domain packages (such as AI agents) are the kfg content layer. This specification defines the public contract that domain packages MUST expose: their public kustomization entrypoints, overlay structure, OpenSpec roots, and test coverage expectations.

## Requirements

### Requirement: Public domain kustomization entrypoint

Each domain package MUST expose a stable public kustomization entrypoint.

#### Scenario: Domain public API
- **WHEN** downstream consumers or the project need domain-specific functionality
- **THEN** `packages/domains/<domain>/kustomization.yaml` SHALL be the public entrypoint
- **AND** consumers MUST NOT reference internal domain paths like `packages/domains/<domain>/manifests/<name>.yaml` directly
- **AND** the public kustomization SHALL compose the framework package and domain-specific manifests

#### Scenario: Domain composition
- **WHEN** a domain package's public kustomization is applied
- **THEN** it SHALL automatically include the framework package as a base
- **AND** domain-specific manifests, steps, or extensions SHALL be layered on top
- **AND** the composition order MUST be: engine → framework → domain

#### Scenario: Public resource guarantees
- **WHEN** a consumer references a domain kustomization
- **THEN** resources exported by that kustomization SHALL remain available across updates
- **AND** breaking changes to exported resources MUST be documented in advance
- **AND** internal manifest paths MAY change without notice

### Requirement: Domain directory structure

Each domain package MUST follow a consistent internal directory layout.

#### Scenario: Domain content organization
- **WHEN** domain content is organized
- **THEN** the structure SHALL be:
  ```
  packages/domains/<domain>/
  ├── kustomization.yaml            # Public entrypoint
  ├── manifests/
  │   ├── kustomization.yaml         # Domain-specific manifests
  │   ├── resources/                 # Domain resources (YAML files)
  │   ├── converters/                # Domain-specific converters
  │   ├── assets/                    # Domain-specific assets
  │   └── ...
  ├── overlays/
  │   ├── dev/                       # Development overlay
  │   │   └── kustomization.yaml
  │   └── ...                        # Additional overlays
  ├── openspec/
  │   ├── config.yaml                # Package-local OpenSpec config
  │   ├── specs/                     # Domain capability specs
  │   └── changes/                   # Domain-specific changes
  └── tests/
      └── ...                        # Domain-specific Bats tests
  ```
- **AND** internal paths MAY change without notice
- **AND** only the public `kustomization.yaml` is guaranteed to remain stable

### Requirement: Domain overlay structure

Domain packages MUST support overlays for different deployment scenarios.

#### Scenario: Development overlay
- **WHEN** a domain package is used in development
- **THEN** `packages/domains/<domain>/overlays/dev/` SHALL provide a development overlay
- **AND** the development overlay SHALL be composable with the base domain kustomization
- **AND** `./bin/kfg build packages/domains/<domain>/overlays/dev` MUST succeed

#### Scenario: Additional overlays
- **WHEN** additional deployment scenarios are needed
- **THEN** additional overlays MAY be added under `packages/domains/<domain>/overlays/<scenario>/`
- **AND** each overlay MUST be composable with the base domain kustomization
- **AND** overlays MAY reference resources from the framework package

### Requirement: Domain-to-framework composition

Domain packages MUST explicitly compose the framework package.

#### Scenario: Framework base composition
- **WHEN** a domain package's public kustomization is defined
- **THEN** it SHALL reference `packages/framework/` as a base using Kustomize `bases` or `resources`
- **AND** domain resources SHALL be layered after the framework composition
- **AND** domain packages MAY extend or override framework steps only where documented

#### Scenario: Framework step visibility
- **WHEN** a domain package applies the framework base
- **THEN** all exported framework steps SHALL be available within the domain package
- **AND** the domain package MAY reference framework steps in its workflows without duplication

### Requirement: Domain OpenSpec root

Each domain package MUST have its own OpenSpec root for domain-specific specifications and changes.

#### Scenario: Domain OpenSpec location
- **WHEN** domain-specific capability specs are defined
- **THEN** `packages/domains/<domain>/openspec/config.yaml` SHALL be the package-local OpenSpec root
- **AND** domain specs SHALL reside under `packages/domains/<domain>/openspec/specs/`
- **AND** domain changes SHALL be tracked under `packages/domains/<domain>/openspec/changes/`

#### Scenario: Domain spec ownership
- **WHEN** a capability spec applies specifically to a domain
- **THEN** it MUST reside in the domain OpenSpec root
- **AND** it MUST NOT be duplicated in the engine-level or framework OpenSpec roots
- **AND** changes to domain behavior MUST be proposed and tracked in the domain OpenSpec

### Requirement: Domain test coverage

Each domain package MUST define Bats test coverage for its resources and overlays.

#### Scenario: Domain test location
- **WHEN** tests are written for domain functionality
- **THEN** they SHALL reside under `packages/domains/<domain>/tests/`
- **AND** they SHALL be discoverable by the canonical `make test-bats` target
- **AND** they SHALL validate domain manifest behavior and overlay composition

#### Scenario: Overlay testing
- **WHEN** domain overlays are added or modified
- **THEN** corresponding Bats coverage SHALL be added under `packages/domains/<domain>/tests/`
- **AND** tests SHALL validate overlay application and kustomization correctness

### Requirement: Stability and independence

Domain packages MUST maintain stability and operate independently where feasible.

#### Scenario: Domain API stability
- **WHEN** downstream systems depend on a domain package
- **THEN** the public `kustomization.yaml` and overlay entrypoints SHALL remain available
- **AND** exported domain behavior SHALL not change in breaking ways without warning
- **AND** the framework API that domain packages depend on SHALL remain stable

#### Scenario: Domain independence
- **WHEN** a domain package is deployed
- **THEN** it SHALL function correctly when paired with any compatible framework version
- **AND** domain-specific specs and changes SHALL not require changes to the engine
- **AND** domain packages SHALL NOT hardcode paths to other domains

### Requirement: Cross-domain considerations

Multiple domain packages in the same repository MUST be able to coexist.

#### Scenario: Multiple domains
- **WHEN** more than one domain package exists in the repository
- **THEN** each SHALL have a distinct path under `packages/domains/`
- **AND** they SHALL NOT reference each other's internal resources directly
- **AND** they SHALL share only the framework package as a common base

#### Scenario: Sibling changes coordination
- **WHEN** a change affects multiple domain packages (e.g., framework update)
- **THEN** each affected domain MUST have a corresponding change in its own OpenSpec root
- **AND** all sibling changes SHALL use the same slug for traceability
- **AND** they MUST be coordinated to land together
