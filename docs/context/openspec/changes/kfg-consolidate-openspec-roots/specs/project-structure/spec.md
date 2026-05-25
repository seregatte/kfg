## MODIFIED Requirements

### Requirement: OpenSpec Structure

OpenSpec artifacts MUST follow a consistent directory structure.

#### Scenario: OpenSpec root location
- **WHEN** the OpenSpec system stores specs and changes
- **THEN** `docs/context/openspec/` SHALL be the single OpenSpec root directory for all layers
- **AND** `docs/context/openspec/config.yaml` SHALL define schema, context, and rules for all layers
- **AND** no separate per-layer OpenSpec roots SHALL exist under `docs/context/`

#### Scenario: Durable specs location
- **WHEN** durable capability specifications are defined
- **THEN** engine specs SHALL reside at `docs/context/openspec/specs/kfg/<capability>/spec.md`
- **AND** framework specs SHALL reside at `docs/context/openspec/specs/framework/<capability>/spec.md`
- **AND** domain specs SHALL reside at `docs/context/openspec/specs/domain-<name>/<capability>/spec.md`
- **AND** specs SHALL use normative language (MUST, SHALL)

#### Scenario: Active changes location
- **WHEN** implementation changes are tracked
- **THEN** `docs/context/openspec/changes/<layer-prefix>-<change-name>/` SHALL contain each change
- **AND** each change SHALL include `proposal.md`, `design.md`, and `tasks.md`
- **AND** cross-layer changes SHALL contain delta specs under the relevant layer subfolders in their `specs/` tree

### Requirement: Package Structure

The project SHALL use a package-oriented structure with explicit layers for engine, framework, and domain packages.

#### Scenario: Framework package location
- **WHEN** shared manifest primitives and reusable steps are defined
- **THEN** `packages/framework/` SHALL contain the framework package
- **AND** `packages/framework/manifests/steps/` SHALL contain shared steps
- **AND** `packages/framework/kustomization.yaml` SHALL be the public entrypoint
- **AND** framework OpenSpec specs SHALL reside at `docs/context/openspec/specs/framework/`
- **AND** `packages/framework/tests/` SHALL contain package-specific Bats suites

#### Scenario: Domain package location
- **WHEN** domain-specific manifests are defined
- **THEN** `packages/domains/<domain>/` SHALL contain each domain package
- **AND** `packages/domains/<domain>/manifests/` SHALL contain domain-specific manifests
- **AND** `packages/domains/<domain>/overlays/` SHALL contain domain-specific overlays
- **AND** `packages/domains/<domain>/kustomization.yaml` SHALL be the public domain entrypoint
- **AND** domain OpenSpec specs SHALL reside at `docs/context/openspec/specs/domain-<domain>/`
- **AND** `packages/domains/<domain>/tests/` SHALL contain package-specific Bats suites

### Requirement: Directory Tree

The directory tree MUST follow the canonical structure.

#### Scenario: Canonical structure
- **WHEN** the project is organized
- **THEN** the directory tree SHALL follow this structure:

```
kfg/
├── flake.nix                         # Nix flake definition
├── flake.lock                        # Nix flake lock
├── go.mod                            # Go module definition
├── go.sum                            # Go module checksum
├── README.md                         # Project documentation
├── src/
│   ├── cmd/kfg/                      # CLI commands
│   └── internal/                     # Internal packages
├── bin/                              # Compiled binaries (gitignored)
├── packages/
│   ├── framework/                    # Shared framework package
│   │   ├── kustomization.yaml        # Framework public entrypoint
│   │   ├── manifests/steps/          # Shared reusable steps
│   │   └── tests/                    # Framework-specific Bats suites
│   └── domains/
│       └── ai-agents/                # AI agents domain package
│           ├── kustomization.yaml    # Domain public entrypoint
│           ├── manifests/            # Domain-specific manifests
│           ├── overlays/dev/         # Development overlay
│           └── tests/                # Domain-specific Bats suites
├── docs/
│   ├── AGENTS.md                     # Agent context file
│   └── context/
│       └── openspec/                 # Single OpenSpec root (all layers)
│           ├── config.yaml           # Unified OpenSpec configuration
│           ├── specs/
│           │   ├── kfg/              # Engine and CLI capability specs
│           │   ├── framework/        # Framework capability specs
│           │   └── domain-ai-agents/ # AI agents domain specs
│           └── changes/              # All changes (layer-prefixed slugs)
│               └── archive/          # Archived changes
├── tests/
│   └── bats/
│       ├── helpers/                  # Shared Bats helpers
│       ├── cli/                      # Engine CLI command tests
│       └── workflows/                # Engine workflow and runtime tests
└── .kfg/
    └── manifests/                    # Project-local manifests
```
