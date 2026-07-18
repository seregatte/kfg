## MODIFIED Requirements

### Requirement: OpenSpec Structure

OpenSpec artifacts MUST use the standard repository-root directory structure.

#### Scenario: OpenSpec root location
- **WHEN** the OpenSpec system stores specs and changes for kfg
- **THEN** `openspec/` SHALL be the only canonical repository OpenSpec root
- **AND** `openspec/config.yaml` SHALL define schema, context, and artifact rules
- **AND** `openspec/README.md` SHALL provide navigation
- **AND** `docs/context/openspec/` SHALL NOT remain as a second root, compatibility symlink, or fallback

#### Scenario: Durable specs location
- **WHEN** durable capability specifications are defined
- **THEN** `openspec/specs/<capability>/spec.md` SHALL contain each spec
- **AND** specs SHALL use normative language such as MUST and SHALL

#### Scenario: Active changes location
- **WHEN** implementation changes are tracked
- **THEN** `openspec/changes/<change-name>/` SHALL contain each active change
- **AND** each change SHALL contain the artifacts required by its OpenSpec schema
- **AND** completed changes SHALL be archived under `openspec/changes/archive/`

### Requirement: Package Structure

The project SHALL use a package-oriented structure with explicit layers for engine, framework, and domain packages plus consolidated OpenSpec specifications.

#### Scenario: Framework package location
- **WHEN** shared manifest primitives and reusable steps are defined
- **THEN** `packages/framework/` SHALL contain the framework package
- **AND** `packages/framework/manifests/steps/` SHALL contain shared steps
- **AND** `packages/framework/kustomization.yaml` SHALL be the public entrypoint
- **AND** framework capability specifications SHALL use the `framework-` prefix under `openspec/specs/`
- **AND** `packages/framework/tests/` SHALL contain package-specific Bats suites

#### Scenario: Domain package location
- **WHEN** domain-specific manifests are defined
- **THEN** `packages/domains/<domain>/` SHALL contain each domain package
- **AND** `packages/domains/<domain>/manifests/` SHALL contain domain-specific manifests
- **AND** `packages/domains/<domain>/overlays/` SHALL contain domain-specific overlays
- **AND** `packages/domains/<domain>/kustomization.yaml` SHALL be the complete-package public entrypoint
- **AND** `packages/domains/<domain>/README.md` SHALL index documented selective entrypoints
- **AND** domain capability specifications SHALL use the domain prefix under `openspec/specs/`
- **AND** `packages/domains/<domain>/tests/` SHALL contain package-specific Bats suites

### Requirement: Directory Tree

The directory tree MUST reflect the standard OpenSpec root and documented domain capability layout.

#### Scenario: Canonical structure
- **WHEN** the project is organized
- **THEN** the directory tree SHALL include this structure:

```text
kfg/
├── flake.nix
├── flake.lock
├── go.mod
├── go.sum
├── README.md
├── openspec/
│   ├── config.yaml
│   ├── README.md
│   ├── specs/
│   └── changes/
│       └── archive/
├── src/
│   ├── cmd/kfg/
│   └── internal/
├── packages/
│   ├── framework/
│   │   ├── kustomization.yaml
│   │   ├── manifests/
│   │   └── tests/
│   └── domains/
│       └── ai-agents/
│           ├── README.md
│           ├── kustomization.yaml
│           ├── manifests/
│           │   ├── kustomization.yaml
│           │   └── <capability>/
│           │       ├── README.md
│           │       └── kustomization.yaml
│           ├── overlays/
│           └── tests/
├── docs/
│   ├── AGENTS.md
│   ├── CHANGELOG.md
│   └── DEVELOPMENT.md
├── tests/
│   └── bats/
└── .kfg/
    └── manifests/
```

- **AND** `openspec/` SHALL be tracked by Git
- **AND** package and capability READMEs SHALL reside adjacent to the public entrypoints they document
