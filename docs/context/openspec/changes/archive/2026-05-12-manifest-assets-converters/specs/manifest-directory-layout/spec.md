# Manifest Directory Layout Specification

## Purpose

This specification defines the `.manifests/` directory structure — a self-contained manifest package at the project root that uses base/overlay composition for reusable agent configurations.

## ADDED Requirements

### Requirement: .manifests/ Root Location

The `.manifests/` directory SHALL reside at the project root and serve as a self-contained manifest package.

#### Scenario: Project root location
- **WHEN** `.manifests/` is created
- **THEN** it SHALL be located at `<project-root>/.manifests/`
- AND it SHALL be a sibling to `src/`, `docs/`, and `flake.nix`

#### Scenario: Git tracking
- **WHEN** the project is versioned
- **THEN** `.manifests/` SHALL be tracked by Git
- AND it SHALL NOT be gitignored

### Requirement: Base Directory Structure

The `base/` subdirectory SHALL contain reusable manifest resources organized by category.

#### Scenario: Base agents
- **WHEN** agent configurations are defined
- **THEN** `base/agents/` SHALL contain one YAML file per agent (claude, gemini, opencode, pi)
- AND each file SHALL define an Assets resource with agent-specific settings and scaffold data

#### Scenario: Base extensions
- **WHEN** extension data is defined
- **THEN** `base/extensions/` SHALL contain subdirectories named after each extension
- AND each extension subdirectory SHALL contain `assets/` and optionally `converters/` and `steps/`
- AND extension names SHALL be: `self`, `ctx7`, `playwright`, `chrome-devtools`, `openspec`, `gws`, `notebooklm`, `ccr`

#### Scenario: Base steps
- **WHEN** core reusable steps are defined
- **THEN** `base/steps/` SHALL contain one YAML file per step
- AND steps SHALL include: detect-agent, copy-context, materialize-scaffold, cleanup-paths, cleanup-workspace, cleanup, install-skill, setup-ccr

### Requirement: Overlay Directory Structure

The `overlay/` subdirectory SHALL contain profile-specific configurations that compose base resources.

#### Scenario: Overlay dev
- **WHEN** a dev profile is defined
- **THEN** `overlay/dev/` SHALL exist
- AND it SHALL contain a `kustomization.yaml` referencing `../../base`
- AND it SHALL contain profile-specific Cmds and CmdWorkflows

#### Scenario: Overlay kustomization
- **WHEN** an overlay kustomization is loaded
- **THEN** it SHALL reference `../../base` as a resource
- AND it SHALL include local `cmds.yaml` and `agents-workflow.yaml`
- AND the overlay layer SHALL have higher precedence than base for duplicate resource identities

### Requirement: Co-located Tests

Tests SHALL be co-located within `.manifests/` alongside the manifests they validate.

#### Scenario: Tests directory
- **WHEN** manifest tests are defined
- **THEN** `.manifests/tests/` SHALL exist
- AND it SHALL contain `.bats` test files and optionally a `test_helper.bash`

#### Scenario: Test runner
- **WHEN** `make test-manifests` is executed
- **THEN** bats SHALL run all `.bats` files in `.manifests/tests/`
- AND the kfg binary SHALL be built before tests run

### Requirement: Canonical Directory Tree

The `.manifests/` directory SHALL follow the canonical structure.

#### Scenario: Canonical structure
- **WHEN** the directory tree is inspected
- **THEN** it SHALL follow this structure:

```
.manifests/
├── base/
│   ├── agents/
│   │   ├── claude.yaml
│   │   ├── gemini.yaml
│   │   ├── opencode.yaml
│   │   └── pi.yaml
│   ├── cmds/
│   │   └── agents.yaml
│   ├── extensions/
│   │   ├── self/
│   │   │   ├── assets/
│   │   │   │   ├── commands/
│   │   │   │   ├── subagents/
│   │   │   │   └── mcp/
│   │   │   ├── converters/
│   │   │   │   ├── commands/
│   │   │   │   ├── mcp/
│   │   │   │   └── subagents/
│   │   │   └── steps/
│   │   ├── ctx7/
│   │   │   ├── assets/
│   │   │   └── steps/
│   │   ├── playwright/
│   │   │   └── assets/
│   │   ├── chrome-devtools/
│   │   │   └── assets/
│   │   ├── openspec/
│   │   │   └── assets/
│   │   ├── gws/
│   │   │   └── assets/
│   │   ├── notebooklm/
│   │   │   └── assets/
│   │   └── ccr/
│   │       └── assets/
│   └── steps/
├── overlay/
│   └── dev/
│       ├── kustomization.yaml
│       ├── cmds.yaml
│       └── agents-workflow.yaml
└── tests/
    ├── test_helper.bash
    ├── manifest-loading.bats
    ├── converters.bats
    └── workflow.bats
```

### Requirement: Makefile Integration

The Makefile SHALL provide targets for running manifest tests.

#### Scenario: test-manifests target
- **WHEN** `make test-manifests` is run
- **THEN** the kfg binary SHALL be built
- AND bats SHALL execute `.manifests/tests/*.bats`
- AND `BATS_LIB_PATH` SHALL be set

#### Scenario: test-all target
- **WHEN** `make test-all` is run
- **THEN** `make test`, `make test-bats`, AND `make test-manifests` SHALL all execute
