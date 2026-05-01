# Project Structure Specification

## Purpose

kfg is a declarative shell compiler that transforms YAML manifests into bash functions. This specification defines the canonical directory layout and organizational structure for the project.

## Requirements

### Requirement: Project Root Structure

The project MUST maintain a clear separation between configuration, documentation, and manifests.

#### Scenario: Core configuration files
- **WHEN** the project is initialized
- **THEN** `flake.nix` SHALL define Nix packaging and development shells
- AND `flake.lock` SHALL pin Nix inputs
- AND configuration files SHALL reside at the project root

#### Scenario: Documentation location
- **WHEN** documentation is organized
- **THEN** `docs/` SHALL contain all documentation
- AND `docs/AGENTS.md` SHALL provide agent operating context
- AND `docs/CHANGELOG.md` SHALL document version history
- AND `docs/DEVELOPMENT.md` SHALL provide development guidelines
- AND `README.md` SHALL reside at the project root (universal convention)

### Requirement: OpenSpec Structure

OpenSpec artifacts MUST follow a consistent directory structure.

#### Scenario: OpenSpec root location
- **WHEN** the OpenSpec system stores specs and changes
- **THEN** `docs/context/openspec/` SHALL be the root directory
- AND `docs/context/openspec/config.yaml` SHALL define schema and context
- AND `docs/context/openspec/README.md` SHALL provide navigation

#### Scenario: Durable specs location
- **WHEN** durable capability specifications are defined
- **THEN** `docs/context/openspec/specs/<capability>/spec.md` SHALL contain each spec
- AND specs SHALL use normative language (MUST, SHALL)

#### Scenario: Active changes location
- **WHEN** implementation changes are tracked
- **THEN** `docs/context/openspec/changes/<change-name>/` SHALL contain each change
- AND each change SHALL include `proposal.md`, `design.md`, and `tasks.md`

### Requirement: Manifest Storage

Project-local manifests MUST be stored in a designated directory.

#### Scenario: Local manifests location
- **WHEN** project-specific command sets are defined
- **THEN** `./.kfg/manifests/` SHALL contain project-local YAML manifests

### Requirement: Directory Tree

The directory tree MUST follow the canonical structure.

#### Scenario: Canonical structure
- **WHEN** the project is organized
- **THEN** the directory tree SHALL follow this structure:

```
kfg_v2/
├── flake.nix                         # Nix flake definition
├── flake.lock                        # Nix flake lock
├── go.mod                            # Go module definition
├── go.sum                            # Go module checksum
├── README.md                         # Project documentation
├── src/
│   ├── cmd/kfg/                    # CLI commands
│   │   ├── main.go                   # Application entry point
│   │   ├── root.go                   # Root cobra command
│   │   └── shell.go                  # Shell subcommand
│   └── internal/                     # Internal packages
│       ├── config/                   # Configuration management
│       ├── logger/                   # Structured logging
│       ├── generate/                 # Shell code generation
│       │   └── templates/            # Template files
│       ├── manifest/                 # Manifest loading and parsing
│       ├── merge/                    # Layer merging
│       ├── resolve/                  # Dependency resolution
│       └── validate/                 # Manifest validation
├── bin/                              # Compiled binaries (gitignored)
├── docs/
│   ├── AGENTS.md                     # Agent context file
│   ├── CHANGELOG.md                  # Version history
│   ├── DEVELOPMENT.md                # Development guidelines
│   └── context/
│       └── openspec/
│           ├── config.yaml           # OpenSpec configuration
│           ├── README.md             # OpenSpec navigation
│           ├── specs/                # Durable specs
│           │   ├── project-structure/
│           │   │   └── spec.md
│           │   ├── cli-conventions/
│           │   │   └── spec.md
│           │   ├── shell-integration/
│           │   │   └── spec.md
│           │   └── manifest-model/
│           │       └── spec.md
│           └── changes/              # Active changes
│               └── pivot-kfg-to-shell-compiler/
│                   ├── proposal.md
│                   ├── design.md
│                   └── tasks.md
└── .kfg/
    └── manifests/                    # Project-local manifests
        ├── steps/
        ├── commands/
        └── sets/
```

### Requirement: Source Code Organization

Source code SHALL be organized in a language-appropriate structure when implementation exists.

#### Scenario: Go source code structure
- **WHEN** Go source code is added to the project
- **THEN** all source code SHALL be placed in `src/`
- AND `src/cmd/` SHALL contain CLI entry points
- AND `src/internal/` SHALL contain internal packages
- AND import paths SHALL use `kfg/src/internal/...` for internal packages

#### Scenario: Binary output location
- **WHEN** Go binaries are compiled
- **THEN** binaries SHALL be output to `bin/`
- AND `bin/` SHALL be excluded from version control (gitignored)

#### Scenario: Go module files location
- **WHEN** Go modules are configured
- **THEN** `go.mod` and `go.sum` SHALL reside at project root
- AND the module name SHALL be `kfg`

#### Scenario: Language-specific structure
- **WHEN** source code is added to the project
- **THEN** it SHALL follow conventions appropriate to the chosen language
- AND build commands SHALL be documented in `README.md` or `docs/DEVELOPMENT.md`
- AND no source code SHALL be placed in `docs/`

### Requirement: No Mixed Locations

Configuration, documentation, and manifests MUST NOT be mixed in inappropriate directories.

#### Scenario: Manifests not in docs
- **WHEN** manifest YAML files are placed
- **THEN** they MUST NOT be in `docs/` or subdirectories of `docs/`
- AND manifests SHALL only exist under configured manifest paths

#### Scenario: Documentation not in manifests
- **WHEN** documentation files are placed
- **THEN** they MUST NOT be in `.kfg/manifests/` or configured manifest paths
- AND documentation SHALL only exist under `docs/` or at root level