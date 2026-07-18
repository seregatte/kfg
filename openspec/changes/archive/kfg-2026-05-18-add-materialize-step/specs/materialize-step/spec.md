## ADDED Requirements

### Requirement: Shared materialize step contract

The repository SHALL provide a shared step named `kfg.materialize` as the canonical manifest-level primitive for converter-driven artifact generation.

#### Scenario: Step exists as shared primitive
- **WHEN** shared base manifests are loaded
- **THEN** a Step resource named `kfg.materialize` SHALL be available for workflow references
- **AND** workflows SHALL use it instead of specialized settings, convert, or aggregate materialization steps

### Requirement: Per-item mode SHALL map assets to outputs positionally

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

### Requirement: Aggregate mode SHALL merge converted assets into one output

When `MODE` is `aggregate`, `kfg.materialize` SHALL convert each asset named in `ASSETS`, deep-merge the converted documents in input order, optionally wrap the merged object under `WRAP_KEY`, and write one output file.

#### Scenario: Aggregate output with wrapper key
- **WHEN** `kfg.materialize` runs with `MODE="aggregate"`
- **AND** `ASSETS` contains `kfg.extension.ctx7.mcp:kfg.extension.chrome-devtools.mcp`
- **AND** `OUTPUTS` contains `.mcp.json`
- **AND** `WRAP_KEY` is `mcpServers`
- **THEN** the converted documents SHALL be deep-merged
- **AND** the merged object SHALL be wrapped under `mcpServers`
- **AND** the result SHALL be written to `.mcp.json`

#### Scenario: Aggregate mode merges with existing output
- **WHEN** `kfg.materialize` runs with `MODE="aggregate"`
- **AND** the target output file already exists
- **THEN** the existing file content SHALL be deep-merged with the newly generated aggregate object before the file is updated

#### Scenario: Aggregate mode requires one output path
- **WHEN** `kfg.materialize` runs with `MODE="aggregate"`
- **AND** `OUTPUTS` contains zero paths or more than one path
- **THEN** the step SHALL fail
- **AND** the error SHALL indicate that aggregate mode requires exactly one output path

### Requirement: Materialize step SHALL validate required inputs strictly

`kfg.materialize` SHALL fail when required inputs are absent or when mode-specific fields are used incorrectly.

#### Scenario: Missing required input fails
- **WHEN** `kfg.materialize` runs without `MODE`, `ASSETS`, `CONVERTER`, or `OUTPUTS`
- **THEN** the step SHALL fail
- **AND** the error SHALL identify the missing required input

#### Scenario: Invalid mode fails
- **WHEN** `kfg.materialize` runs with a `MODE` value other than `per-item` or `aggregate`
- **THEN** the step SHALL fail
- **AND** the error SHALL indicate the supported mode values

#### Scenario: Wrapper key is ignored outside aggregate mode is not allowed
- **WHEN** `kfg.materialize` runs with `MODE="per-item"`
- **AND** `WRAP_KEY` is set
- **THEN** the step SHALL fail
- **AND** the error SHALL indicate that `WRAP_KEY` is only valid in aggregate mode
