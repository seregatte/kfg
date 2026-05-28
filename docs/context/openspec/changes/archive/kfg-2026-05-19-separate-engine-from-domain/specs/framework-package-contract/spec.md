## ADDED Requirements

### Requirement: Framework package public entrypoint
The repository SHALL expose the shared framework manifests through `packages/framework/kustomization.yaml` as the public package entrypoint.

#### Scenario: Referencing the framework package
- **WHEN** a domain package or downstream consumer needs shared manifest primitives
- **THEN** it SHALL reference `packages/framework/kustomization.yaml`
- **AND** it SHALL NOT depend on framework internal subpaths as the public API

### Requirement: Framework package OpenSpec root
The framework package SHALL own a package-local OpenSpec root.

#### Scenario: Storing framework specs and changes
- **WHEN** framework behavior is specified or changed
- **THEN** the package SHALL use `packages/framework/openspec/`
- **AND** that root SHALL contain `config.yaml`, `specs/`, and `changes/`

### Requirement: Framework package test root
The framework package SHALL own package-local Bats coverage for shared manifest primitives.

#### Scenario: Running framework Bats coverage
- **WHEN** shared steps are validated by Bats
- **THEN** their suites SHALL reside under `packages/framework/tests/`

### Requirement: Framework package exported primitives
The framework package SHALL own the reusable `kfg.*` step primitives that are shared across domains.

#### Scenario: Exporting shared steps
- **WHEN** the framework package is built through its public entrypoint
- **THEN** it SHALL expose shared steps including `kfg.cleanup`, `kfg.materialize`, `kfg.materialize-scaffold`, `kfg.ensure-gitignore`, and `kfg.copy-context`
