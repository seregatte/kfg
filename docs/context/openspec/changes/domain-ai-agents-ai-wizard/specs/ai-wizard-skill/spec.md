## ADDED Requirements

### Requirement: WIZARD SKILL PROMPT

The wizard skill SHALL be defined as a `kind: Assets` resource with `metadata.name: ai.prompts.kfg-wizard`.

The prompt content SHALL instruct the agent to:
1. Learn the user's requirements through interactive questions (never assume)
2. Compose configurations from existing domain building blocks before creating custom ones
3. Present a summary of what will be generated before writing any files
4. Ask for confirmation before writing files
5. Explore the domain manifests directory if uncertain about available building blocks

#### Scenario: Skill defined as Assets resource
- **WHEN** the asset `ai.prompts.kfg-wizard` is loaded
- **THEN** it SHALL have `apiVersion: kfg.dev/v1alpha1`
- **AND** `kind: Assets`
- **AND** `spec.input.format: yaml`
- **AND** `spec.data.name: kfg-wizard`
- **AND** `spec.data.description` SHALL not be empty
- **AND** `spec.data.prompt` SHALL contain agent instructions

### Requirement: WIZARD BEHAVIOR — INTERACTIVE BY DEFAULT

The wizard SHALL operate in interactive mode by default, asking questions until it has a complete understanding of the user's requirements.

- The agent SHALL ask about: project type, language/stack, required AI agents, required skills/domains, directory structure preferences
- The agent SHALL NOT write files until it explicitly asks and receives user confirmation

#### Scenario: Interactive discovery
- **WHEN** user starts the wizard with `kfg ai`
- **THEN** the wizard SHALL prompt the user with questions about their project requirements

### Requirement: WIZARD BEHAVIOR — COMPOSE FIRST

Before creating any custom manifest, the wizard SHALL verify that the required capability does not already exist in the domain catalog.

- The agent SHALL maintain a mental catalog of available building blocks (steps, cmds, converters, assets)
- If the required capability exists, the wizard SHALL reference it via kustomize `resources:` rather than creating a new one
- Only when nothing in the domain catalog matches, the wizard SHALL create custom manifest resources

#### Scenario: Existing building block used
- **WHEN** user requests ctx7 support
- **THEN** the wizard SHALL reference `ctx7.steps.install` from the domain manifests
- **AND** SHALL NOT create a custom install step for ctx7

#### Scenario: Custom manifest created
- **WHEN** user requests integration with a service not in the domain catalog
- **THEN** the wizard SHALL create custom manifest resources for that service
- **AND** SHALL inform the user that a new custom step/asset was created

### Requirement: WIZARD BEHAVIOR — CONFIRMATION BEFORE WRITING

The agent SHALL present a summary of what will be generated before writing any files.

- The summary SHALL list all files to be created with their paths
- The wizard SHALL ask "Shall I proceed?" and wait for user input
- On confirmation ("yes"), the wizard SHALL create all files
- On rejection ("no"), the wizard SHALL ask what to change

#### Scenario: Confirmation prompt
- **WHEN** the wizard has gathered all requirements
- **THEN** the wizard SHALL display a summary of files to be created
- **AND** SHALL ask for user confirmation before writing
