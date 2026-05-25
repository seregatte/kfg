## MODIFIED Requirements

### Requirement: Context management helpers

The engine runtime MUST provide context management primitives for steps that modify environment or state.

#### Scenario: Context reset
- **WHEN** a step needs to clean up its modifications
- **THEN** `__kfg_ctx_reset` SHALL provide a mechanism to restore prior context
- **AND** the mechanism SHALL be step-aware (resetting only the current step's modifications)

#### Scenario: Output helpers
- **WHEN** a step needs to produce output for the generated shell to consume
- **THEN** output helpers SHALL be available to facilitate structured output
- **AND** helpers SHALL integrate with the artifact and result collection system

#### Scenario: Output-producing step avoids subshell side-effect loss
- **WHEN** a Step declares `spec.output`
- **THEN** the generated runtime SHALL capture the output value without executing the Step body in a subshell that loses shell-side effects
- **AND** runtime state changes such as artifact registration SHALL remain visible after the Step finishes
