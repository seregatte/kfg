# Bats Test Layout Specification

## Purpose

This specification defines the repository's Bats test layout including multi-root test discovery, shared helpers, and canonical execution targets.
## Requirements

### Requirement: Multi-root Bats test discovery
The repository MUST discover Bats tests from engine and package-local roots.

#### Scenario: Discovering engine Bats tests
- **WHEN** a contributor searches for engine tests
- **THEN** engine and integration tests SHALL reside under `tests/bats/`
- **AND** the canonical `make test-bats` target SHALL execute engine tests

#### Scenario: Discovering package Bats tests
- **WHEN** a package defines package-specific Bats suites
- **THEN** package tests SHALL reside under `packages/<package>/tests/`
- **AND** the canonical `make test-bats` target SHALL discover and execute package tests

#### Scenario: Repository root detection
- **WHEN** a Bats suite from any root needs to discover the repository root
- **THEN** it SHALL NOT assume it lies under `tests/bats/`
- **AND** it SHALL use an explicit root detection mechanism independent of path position

### Requirement: Engine and package test organization
Bats tests MUST be organized around content ownership and responsibility.

#### Scenario: Engine test location
- **WHEN** a Bats test validates engine CLI commands or workflow runtime behavior
- **THEN** the test SHALL reside under `tests/bats/`
- **AND** tests MAY be organized into subdirectories such as `cli/` or `workflows/`

#### Scenario: Package test location
- **WHEN** a Bats test validates package-specific resources, steps, or overlays
- **THEN** the test SHALL reside under `packages/<package>/tests/`
- **AND** tests MAY mirror the package structure they validate

### Requirement: Shared Bats helpers by concern
The repository MUST provide shared Bats helpers under `tests/bats/helpers/` with concern-specific entrypoints.

#### Scenario: Loading common Bats helpers
- **WHEN** a Bats suite from engine or any package needs repository root or binary bootstrap logic
- **THEN** it SHALL source shared helper code from `tests/bats/helpers/`
- **AND** repository-wide helper behavior SHALL NOT be duplicated across multiple helper roots

#### Scenario: Package-aware manifest helpers
- **WHEN** a Bats suite executes a manifest resource such as a Step or overlay workflow
- **THEN** it SHALL use helper functions that accept explicit package paths rather than hardcoding `.manifests/`
- **AND** helpers SHALL resolve manifest paths dynamically based on package entrypoint location

### Requirement: Canonical Bats execution target
Repository Bats entrypoints MUST discover and run tests from all roots.

#### Scenario: Running the canonical Bats target
- **WHEN** repository Bats tests are invoked through the canonical `make test-bats` target
- **THEN** the target SHALL execute suites from `tests/bats/` and all package roots under `packages/*/tests/`
- **AND** results SHALL aggregate across all roots
- **AND** contributors SHALL NOT need to know individual package root locations