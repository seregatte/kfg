# Project Structure Specification

## Purpose

kfg is a declarative shell compiler that transforms YAML manifests into bash functions. This specification defines the canonical directory layout and organizational structure for the project.

## MODIFIED Requirements

### Requirement: Manifest Storage

Project-local manifests MAY be stored in one of two designated directories.

#### Scenario: Local manifests in .kfg/
- **WHEN** project-specific command sets are defined
- **THEN** `./.kfg/manifests/` MAY contain project-local YAML manifests
- **AND** this path SHALL be included in `KFG_MANIFEST_PATH` by default

#### Scenario: Local manifests in .manifests/
- **WHEN** a self-contained manifest package is defined
- **THEN** `./.manifests/` MAY contain a base/overlay manifest structure
- **AND** this path SHALL be included in `KFG_MANIFEST_PATH` by default
- **AND** the directory SHALL include `base/` for reusable resources and `overlay/` for profile-specific compositions

### Requirement: Directory Tree

The directory tree SHALL support both `.kfg/` and `.manifests/` layouts.

#### Scenario: Canonical structure with .manifests/
- **WHEN** the project uses `.manifests/`
- **THEN** the directory tree SHALL include:

```
kfg/
├── flake.nix
├── flake.lock
├── go.mod
├── go.sum
├── README.md
├── .manifests/
│   ├── base/
│   │   ├── agents/
│   │   ├── cmds/
│   │   ├── extensions/
│   │   └── steps/
│   ├── overlay/
│   │   └── dev/
│   │       ├── kustomization.yaml
│   │       ├── cmds.yaml
│   │       └── agents-workflow.yaml
│   └── tests/
│       └── *.bats
├── src/
│   └── ...
└── docs/
    └── ...
```

#### Scenario: Both layouts coexist
- **WHEN** both `.kfg/` and `.manifests/` exist
- **THEN** both SHALL be loaded via `KFG_MANIFEST_PATH`
- **AND** rightmost paths SHALL have higher precedence for duplicate resource identities
