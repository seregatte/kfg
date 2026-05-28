## MODIFIED Requirements

### Requirement: Project Root Structure

The project MUST maintain a clear separation between configuration, documentation, manifests, and test assets.

#### Scenario: Core configuration files
- **WHEN** the project is initialized
- **THEN** `flake.nix` SHALL define Nix packaging and development shells
- **AND** `flake.lock` SHALL pin Nix inputs
- **AND** configuration files SHALL reside at the project root

#### Scenario: Documentation location
- **WHEN** documentation is organized
- **THEN** `docs/` SHALL contain all documentation
- **AND** `docs/AGENTS.md` SHALL provide agent operating context
- **AND** `docs/CHANGELOG.md` SHALL document version history
- **AND** `docs/DEVELOPMENT.md` SHALL provide development guidelines
- **AND** `README.md` SHALL reside at the project root (universal convention)

#### Scenario: Test asset location
- **WHEN** repository shell tests are organized
- **THEN** Bats tests SHALL reside under `tests/bats/`
- **AND** shared Bats helpers SHALL reside under `tests/bats/helpers/`
- **AND** manifest-resource Bats tests SHALL be organized under `tests/bats/manifests/`

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
│   ├── cmd/kfg/                      # CLI commands
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
│           └── changes/              # Active changes
├── tests/
│   └── bats/
│       ├── helpers/                  # Shared Bats helpers
│       ├── manifests/                # Mirrored tests for .manifests resources
│       │   ├── base/
│       │   └── overlay/
│       └── workflows/                # Generic workflow/runtime Bats tests
└── .kfg/
    └── manifests/                    # Project-local manifests
```

### Requirement: No Mixed Locations

Configuration, documentation, manifests, and test assets MUST NOT be mixed in inappropriate directories.

#### Scenario: Manifests not in docs
- **WHEN** manifest YAML files are placed
- **THEN** they MUST NOT be in `docs/` or subdirectories of `docs/`
- **AND** manifests SHALL only exist under configured manifest paths

#### Scenario: Documentation not in manifests
- **WHEN** documentation files are placed
- **THEN** they MUST NOT be in `.kfg/manifests/` or configured manifest paths
- **AND** documentation SHALL only exist under `docs/` or at root level

#### Scenario: Bats tests not in manifest directories
- **WHEN** repository Bats tests are placed
- **THEN** they MUST NOT be stored under `.manifests/tests/`
- **AND** supported Bats test files SHALL only exist under `tests/bats/`
