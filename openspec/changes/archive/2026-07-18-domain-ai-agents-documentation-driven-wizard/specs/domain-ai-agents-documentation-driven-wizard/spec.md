## ADDED Requirements

### Requirement: Wizard discovers the domain catalog

The AI wizard SHALL discover available capabilities from the installed AI agents domain documentation instead of relying on a hardcoded resource catalog.

#### Scenario: Domain index discovery
- **WHEN** the wizard starts capability discovery
- **THEN** it SHALL locate the domain through `KFG_DIR`
- **AND** it SHALL read `packages/domains/ai-agents/README.md`
- **AND** it SHALL follow links only for capabilities relevant to the user's requirements

#### Scenario: Newly documented capability
- **WHEN** a capability is added to the domain index with an adjacent README and public Kustomization
- **THEN** the wizard SHALL be able to discover it without changing the wizard's embedded capability list

### Requirement: Wizard verifies documented capabilities

The wizard MUST verify selected capability documentation against its adjacent Kustomization before proposing a configuration.

#### Scenario: Documented resource is available
- **WHEN** a README advertises a resource or selective entrypoint
- **THEN** the wizard SHALL inspect the adjacent Kustomization
- **AND** it SHALL use only resource names and entrypoints present in the composed manifests

#### Scenario: Documentation and manifests disagree
- **WHEN** an advertised resource or entrypoint cannot be verified
- **THEN** the wizard MUST report the inconsistency
- **AND** it MUST NOT invent a replacement resource or silently use an unrelated aggregate

### Requirement: Wizard uses compatibility-aware discovery

The wizard SHALL ask capability questions based on documented agent compatibility and prerequisites.

#### Scenario: Supported combination
- **WHEN** the user selects an agent and a capability documented as compatible
- **THEN** the wizard SHALL include the documented entrypoint and workflow pattern in its proposal

#### Scenario: Unsupported combination
- **WHEN** the user selects an agent and capability combination that is not documented as supported
- **THEN** the wizard SHALL explain the incompatibility
- **AND** it SHALL ask the user to choose a supported alternative or approve custom development

### Requirement: Wizard generates selective overlays

The wizard SHALL prefer selective capability entrypoints and JSON Patch customization in generated project overlays.

#### Scenario: Existing capability selected
- **WHEN** a requested capability has a selective entrypoint
- **THEN** the generated Kustomization SHALL reference that entrypoint
- **AND** it SHALL NOT load the complete AI agents domain unless the user explicitly requests the full package

#### Scenario: Capability customization
- **WHEN** user requirements differ from documented defaults
- **THEN** the wizard SHALL generate JSON Patch against the documented resource name and path
- **AND** it SHALL NOT duplicate the base resource in the generated project

### Requirement: Wizard previews operational impact

The wizard MUST present the complete generated composition before writing files.

#### Scenario: Pre-write confirmation
- **WHEN** discovery and configuration are complete
- **THEN** the wizard SHALL list selected entrypoints, JSON Patches, prerequisites, generated artifacts, cleanup behavior, unsupported combinations, and output paths
- **AND** it MUST receive explicit user confirmation before writing

### Requirement: Wizard validates generated configuration

The wizard SHALL validate generated Kustomizations after writing them.

#### Scenario: Generated overlay succeeds
- **WHEN** the user confirms generation and files are written
- **THEN** the wizard SHALL build the generated Kustomization
- **AND** it SHALL report the validation result and usable next command

#### Scenario: Generated overlay fails
- **WHEN** the generated Kustomization does not build
- **THEN** the wizard SHALL report the validation error
- **AND** it SHALL preserve enough context for the user to review the generated files
- **AND** it SHALL NOT claim successful configuration

### Requirement: Wizard prompt avoids catalog duplication

The wizard prompt SHALL encode the discovery protocol and safety rules without embedding the current domain resource inventory.

#### Scenario: Wizard prompt validation
- **WHEN** domain tests inspect the wizard prompt
- **THEN** the prompt SHALL require reading the domain and selected capability READMEs
- **AND** it SHALL require adjacent Kustomization verification
- **AND** it SHALL NOT contain a hardcoded list presented as the authoritative capability catalog
