# Framework Package Contract Specification

## Purpose

The framework package is the kfg shared manifest layer. This specification defines the public contract that framework packages MUST expose: their exported primitives, public kustomization entrypoints, OpenSpec roots, and test coverage expectations.

## Requirements

### Requirement: Public kustomization entrypoint

The framework package MUST expose a stable public kustomization entrypoint.

#### Scenario: Framework public API
- **WHEN** downstream consumers or domain packages need shared framework steps
- **THEN** `packages/framework/kustomization.yaml` SHALL be the public entrypoint
- **AND** consumers MUST NOT reference internal framework paths like `packages/framework/manifests/steps/<name>.yaml` directly
- **AND** the public kustomization SHALL compose all exported framework resources

#### Scenario: Public resource guarantees
- **WHEN** a consumer references the framework kustomization
- **THEN** resources exported by that kustomization SHALL remain available across framework updates
- **AND** breaking changes to exported resources MUST be documented in advance
- **AND** deprecated resources SHALL have a removal timeline

### Requirement: Framework directory structure

The framework package MUST follow a consistent internal directory layout.

#### Scenario: Framework content organization
- **WHEN** framework content is organized
- **THEN** the structure SHALL be:
  ```
  packages/framework/
  ├── kustomization.yaml           # Public entrypoint
  ├── manifests/
  │   ├── kustomization.yaml        # Manifest composition
  │   └── steps/                    # Reusable steps
  │       ├── cleanup-step.yaml
  │       ├── ensure-gitignore.yaml
  │       ├── materialize-step.yaml
  │       └── ...
  ├── overlays/                     # Framework overlays (if any)
  ├── openspec/
  │   ├── config.yaml               # Package-local OpenSpec config
  │   ├── specs/                    # Framework capability specs
  │   └── changes/                  # Framework-specific changes
  └── tests/
      └── ...                       # Framework-specific Bats tests
  ```
- **AND** internal paths MAY change without notice
- **AND** only the public `kustomization.yaml` is guaranteed to remain stable

### Requirement: Exported reusable steps

The framework package MUST export a documented set of reusable manifest steps.

#### Scenario: Framework step exports
- **WHEN** framework capabilities are defined
- **THEN** the following reusable steps MUST be exported:
  - `kfg.materialize`: Generate shell from manifests
  - `kfg.cleanup`: Clean up generated artifacts
  - `kfg.ensure-gitignore`: Manage gitignore entries
  - `kfg.copy-context`: Copy context files into generated artifacts
  - `kfg.materialize-scaffold`: Generate scaffolding from templates
- **AND** each exported step SHALL have a documented specification in `packages/framework/openspec/specs/`

#### Scenario: New framework step addition
- **WHEN** new framework steps are added
- **THEN** they SHALL be spec'd in the framework OpenSpec
- **AND** they SHALL be exported through the public kustomization
- **AND** they SHALL be backward-compatible or follow a deprecation timeline

### Requirement: Framework OpenSpec root

The framework package MUST have its own OpenSpec root for package-specific specifications and changes.

#### Scenario: Framework OpenSpec location
- **WHEN** framework-specific capability specs are defined
- **THEN** `packages/framework/openspec/config.yaml` SHALL be the package-local OpenSpec root
- **AND** framework specs SHALL reside under `packages/framework/openspec/specs/`
- **AND** framework changes SHALL be tracked under `packages/framework/openspec/changes/`

#### Scenario: Framework spec ownership
- **WHEN** a capability spec applies specifically to framework behavior
- **THEN** it MUST reside in the framework OpenSpec root
- **AND** it MUST NOT be duplicated in the engine-level OpenSpec root
- **AND** changes to framework behavior MUST be proposed and tracked in the framework OpenSpec

### Requirement: Framework test coverage

The framework package MUST define Bats test coverage for exported primitives.

#### Scenario: Framework test location
- **WHEN** tests are written for framework functionality
- **THEN** they SHALL reside under `packages/framework/tests/`
- **AND** they SHALL be discoverable by the canonical `make test-bats` target
- **AND** they SHALL validate exported step behavior and shell runtime contracts

#### Scenario: Exported step testing
- **WHEN** an exported step is added or modified
- **THEN** corresponding Bats coverage SHALL be added under `packages/framework/tests/`
- **AND** tests SHALL validate both the step's manifest behavior and its shell generation

### Requirement: Stability commitment

The framework package MUST maintain a stability commitment to downstream packages.

#### Scenario: API stability guarantee
- **WHEN** downstream packages or domain consumers depend on framework exports
- **THEN** the public `kustomization.yaml` and exported step names SHALL remain available
- **AND** exported step behavior SHALL not change in breaking ways
- **AND** the `shell-runtime-api` that framework steps depend on SHALL remain stable

#### Scenario: Deprecation process
- **WHEN** a framework feature must be removed or fundamentally changed
- **THEN** a deprecation period MUST be announced in advance
- **AND** migration documentation SHALL be provided
- **AND** a release timeline SHALL be published
