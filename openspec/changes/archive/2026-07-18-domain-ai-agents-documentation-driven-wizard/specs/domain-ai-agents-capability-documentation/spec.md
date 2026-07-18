## ADDED Requirements

### Requirement: Domain capability index

The AI agents domain SHALL provide a root README that acts as the public capability index for human users and agents.

#### Scenario: Consumer discovers domain capabilities
- **WHEN** a consumer reads `packages/domains/ai-agents/README.md`
- **THEN** the index SHALL distinguish the complete package, the domain-only aggregate, selective capability entrypoints, and operational overlays
- **AND** it SHALL link to the usage documentation for every public capability
- **AND** it SHALL identify deprecated entrypoints

### Requirement: Capability-specific usage documentation

Every independently consumable capability in the AI agents domain MUST provide adjacent usage documentation.

#### Scenario: Capability README contract
- **WHEN** a capability exposes a public Kustomization entrypoint
- **THEN** an adjacent `README.md` MUST document its purpose, prerequisites, entrypoints, exported resource names and kinds, supported agents, configurable fields and defaults, workflow usage, generated artifacts, cache and cleanup behavior, known limitations, validation commands, and canonical specifications

#### Scenario: Capability composition example
- **WHEN** a capability README documents consumption
- **THEN** it MUST include a minimal `resources:` composition example
- **AND** it MUST include JSON Patch examples for its supported customization points
- **AND** every patch target MUST identify an exported resource by stable `kind` and `metadata.name`

### Requirement: Selective capability composition

Every independently useful capability SHALL expose a Kustomization entrypoint that can be composed without loading unrelated domain capabilities.

#### Scenario: Selective entrypoint builds
- **WHEN** a consumer builds an advertised selective capability entrypoint
- **THEN** the build MUST succeed
- **AND** it SHALL include the resources documented by the adjacent README
- **AND** it SHALL NOT include unrelated capability resources

#### Scenario: Aggregate entrypoints remain available
- **WHEN** existing consumers use the domain or category aggregate Kustomizations
- **THEN** those entrypoints SHALL continue to build
- **AND** they SHALL aggregate the corresponding selective capability resources

### Requirement: Variants reflect semantic composition

The domain MUST use separate Kustomizations only when a choice changes resource composition or operational semantics.

#### Scenario: Value-only customization
- **WHEN** a consumer changes a model, enabled flag, command argument, output path, store ID, remote, or another manifest value
- **THEN** the capability documentation SHALL direct the consumer to JSON Patch
- **AND** the domain SHALL NOT require a named variant solely for that value

#### Scenario: Semantic variant
- **WHEN** two supported modes export different resources or change where an external tool reads or writes state
- **THEN** each mode MAY expose a separate documented Kustomization
- **AND** the README MUST document whether the modes are exclusive or composable

### Requirement: Documentation changes with public behavior

Changes to a public AI agents capability MUST update its adjacent documentation in the same change.

#### Scenario: Public capability changes
- **WHEN** a change modifies an exported resource name, default, configurable field, compatibility rule, generated artifact, cleanup behavior, workflow usage, or Kustomization entrypoint
- **THEN** the capability's README MUST be updated in the same change
- **AND** `docs/AGENTS.md` MUST instruct contributors and agents to enforce this requirement

### Requirement: Documentation contract validation

The AI agents domain SHALL automatically validate its capability documentation contract.

#### Scenario: Domain documentation tests
- **WHEN** `make test-bats` runs
- **THEN** domain tests SHALL verify README coverage and required sections for public capabilities
- **AND** they SHALL verify advertised entrypoints build
- **AND** they SHALL verify documented resource names and representative JSON Patch targets exist
