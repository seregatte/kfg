## ADDED Requirements

### Requirement: Single OpenSpec Root
The project SHALL maintain a single OpenSpec root at `docs/context/openspec/` covering all layers (engine, framework, domain).

#### Scenario: Root has config and structure
- **WHEN** the OpenSpec root is initialized
- **THEN** `docs/context/openspec/config.yaml` SHALL define schema, context, and rules for all layers
- **AND** `docs/context/openspec/specs/` SHALL contain all durable capability specifications
- **AND** `docs/context/openspec/changes/` SHALL contain all active and archived changes
- **AND** no other OpenSpec roots SHALL exist under `docs/context/`

#### Scenario: OPENSPEC_ROOT_DIR points at the root
- **WHEN** the development shell is initialized via `flake.nix`
- **THEN** `OPENSPEC_ROOT_DIR` SHALL be set to `docs/context/openspec`

### Requirement: Layer Subfolders in Specs
Specs SHALL be organized into layer subfolders within `docs/context/openspec/specs/`.

#### Scenario: Engine specs location
- **WHEN** engine or CLI capability specs are defined
- **THEN** they SHALL reside under `docs/context/openspec/specs/kfg/<capability>/spec.md`

#### Scenario: Framework specs location
- **WHEN** framework package capability specs are defined
- **THEN** they SHALL reside under `docs/context/openspec/specs/framework/<capability>/spec.md`

#### Scenario: Domain specs location
- **WHEN** domain package capability specs are defined
- **THEN** they SHALL reside under `docs/context/openspec/specs/domain-ai-agents/<capability>/spec.md`

### Requirement: Layer Prefix on Change Slugs
Change slugs SHALL be prefixed with their originating layer.

#### Scenario: Engine change slug
- **WHEN** a change primarily affects the engine or CLI
- **THEN** its directory SHALL be named `kfg-<slug>/` under `changes/` or `changes/archive/`

#### Scenario: Framework change slug
- **WHEN** a change primarily affects the framework package
- **THEN** its directory SHALL be named `framework-<slug>/` under `changes/` or `changes/archive/`

#### Scenario: Domain change slug
- **WHEN** a change primarily affects a domain package
- **THEN** its directory SHALL be named `domain-<domain-name>-<slug>/` under `changes/` or `changes/archive/`

#### Scenario: Cross-layer change slug
- **WHEN** a change spans multiple layers
- **THEN** its slug SHALL use the prefix of the primary layer it originates from
- **AND** its internal `specs/` tree SHALL contain delta specs under the appropriate layer subfolders

### Requirement: Consolidated Config
The single `docs/context/openspec/config.yaml` SHALL cover context and rules for all layers.

#### Scenario: Config covers all layers
- **WHEN** the config is read by the OpenSpec tool
- **THEN** the `context:` section SHALL document engine, framework, and domain layer details
- **AND** the `rules:` section SHALL apply across all layers with per-layer notes where behavior differs

### Requirement: No Sibling Change Roots
Cross-layer changes SHALL NOT be duplicated across multiple OpenSpec roots.

#### Scenario: Single entry for cross-layer change
- **WHEN** a change touches both engine and framework specs
- **THEN** a single change directory SHALL exist under `docs/context/openspec/changes/`
- **AND** its `specs/` tree SHALL contain delta files under both `specs/kfg/` and `specs/framework/`
- **AND** no sibling change directories SHALL exist in separate per-layer roots
