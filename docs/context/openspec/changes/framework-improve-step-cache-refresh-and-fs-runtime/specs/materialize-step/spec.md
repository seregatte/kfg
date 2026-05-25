## MODIFIED Requirements

### Requirement: Shared materialize Step Contract

The repository SHALL provide a shared step named `kfg.materialize` as the canonical manifest-level primitive for converter-driven artifact generation.

#### Scenario: Step exists as shared primitive
- **WHEN** shared base manifests are loaded
- **THEN** a Step resource named `kfg.materialize` SHALL be available for workflow references
- **AND** workflows SHALL use it instead of specialized settings, convert, or aggregate materialization steps

### Requirement: Per-item Mode SHALL Map Assets to Outputs Positionally

When `MODE` is `per-item`, `kfg.materialize` SHALL convert each asset named in `ASSETS` with the specified `CONVERTER` and write each converted result to the corresponding output path in `OUTPUTS`.

#### Scenario: Single-item per-item materialization
- **WHEN** `kfg.materialize` runs with `MODE="per-item"`
- **AND** `ASSETS` contains `kfg.extension.self.commands.git-commit`
- **AND** `OUTPUTS` contains `.claude/commands/git-commit.md`
- **THEN** the step SHALL write the converted result to `.claude/commands/git-commit.md`
- **AND** it SHALL register that path as an artifact

#### Scenario: Multi-item per-item materialization
- **WHEN** `kfg.materialize` runs with `MODE="per-item"`
- **AND** `ASSETS` contains `a:b:c`
- **AND** `OUTPUTS` contains `x:y:z`
- **THEN** asset `a` SHALL be converted into `x`
- **AND** asset `b` SHALL be converted into `y`
- **AND** asset `c` SHALL be converted into `z`

#### Scenario: Per-item count mismatch fails
- **WHEN** `kfg.materialize` runs with `MODE="per-item"`
- **AND** `ASSETS` and `OUTPUTS` contain different numbers of items
- **THEN** the step SHALL fail
- **AND** the error SHALL indicate that per-item input and output counts must match

#### Scenario: Nested conversions use quiet internal kfg execution
- **WHEN** `kfg.materialize` performs nested `kfg` conversion calls
- **THEN** it SHALL invoke those commands through the runtime's internal execution wrapper
- **AND** child startup logs SHALL NOT appear in the parent step's human stderr output
