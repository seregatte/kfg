## MODIFIED Requirements

### Requirement: Openspec base Cmd

The system SHALL provide a Cmd resource that wraps the `openspec` CLI binary as a shell function and delegates root resolution to the CLI.

#### Scenario: Openspec Cmd generates shell function
- **WHEN** the workflow includes `ai.cmds.openspec`
- **THEN** the generated shell SHALL contain a function `openspec()` that delegates to the `openspec` binary
- **AND** the wrapper SHALL NOT change directories based on `OPENSPEC_ROOT_DIR`

#### Scenario: Openspec Cmd passes arguments
- **WHEN** the user invokes `openspec list`
- **THEN** the function SHALL execute `command openspec "list"`
- **AND** all arguments, including `--store`, SHALL be forwarded unchanged

#### Scenario: Native root resolution
- **WHEN** the wrapper is invoked inside a project or standalone store
- **THEN** the OpenSpec CLI SHALL resolve an explicit store, nearest local root, declared pointer, or configured fallback according to its native precedence
- **AND** the wrapper SHALL NOT override the resolved root

### Requirement: Openspec Cmd follows naming convention

The Cmd resource SHALL use `metadata.name: ai.cmds.openspec` with `commandName: openspec`.

#### Scenario: Cmd metadata
- **WHEN** the manifest is parsed
- **THEN** the Cmd SHALL have `metadata.name: ai.cmds.openspec`
- **AND** it SHALL have `metadata.commandName: openspec`
