## ADDED Requirements

### Requirement: Openspec base Cmd

The system SHALL provide a Cmd resource that wraps the `openspec` CLI binary as a shell function, following the standard agent cmd pattern.

#### Scenario: Openspec Cmd generates shell function
- **WHEN** the workflow includes `kfg.agent.cmd.openspec`
- **THEN** the generated shell contains a function `openspec()` that delegates to the `openspec` binary

#### Scenario: Openspec Cmd passes arguments
- **WHEN** the user invokes `openspec list` in the generated shell
- **THEN** the function executes `command openspec "list"`

### Requirement: Openspec Cmd follows naming convention

The Cmd resource SHALL use the naming convention `kfg.agent.cmd.openspec` with `commandName: openspec`.

#### Scenario: Cmd metadata
- **WHEN** the manifest is parsed
- **THEN** the Cmd has `metadata.name: kfg.agent.cmd.openspec` and `metadata.commandName: openspec`
