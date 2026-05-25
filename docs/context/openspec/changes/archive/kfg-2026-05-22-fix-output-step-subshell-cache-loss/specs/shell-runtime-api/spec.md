## MODIFIED Requirements

### Requirement: Context management helpers

The engine runtime MUST provide context management primitives for steps that modify environment or state.

#### Scenario: Output helpers
- **WHEN** a step needs to produce output for the generated shell to consume
- **THEN** output helpers SHALL be available to facilitate structured output
- **AND** helpers SHALL integrate with the artifact and cache collection system

#### Scenario: Output-producing step preserves runtime side effects
- **WHEN** a Step declares `spec.output`
- **THEN** the generated runtime SHALL execute that Step without a subshell that discards runtime side effects
- **AND** artifact registrations performed during the Step SHALL remain visible to the parent shell runtime after execution
