# Materialize Step Specification

## Purpose

The `kfg.materialize` step is the canonical manifest-level primitive for converter-driven artifact generation. It provides a unified contract for both per-item and aggregate materialization patterns, replacing the fragmented `kfg.convert`, `kfg.aggregate-mcp`, and `kfg.agents.steps.settings` steps.

## Requirements

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

### Requirement: Aggregate Mode SHALL Merge Converted Assets into One Output

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

### Requirement: Materialize Step SHALL Validate Required Inputs Strictly

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

## Data Contract

The step contract is:

```yaml
env:
  MODE: ""        # per-item | aggregate (required)
  ASSETS: ""      # colon-separated asset names (required)
  CONVERTER: ""   # converter metadata.name (required)
  OUTPUTS: ""     # colon-separated output paths; exactly one in aggregate mode (required)
  WRAP_KEY: ""    # optional, aggregate mode only
```

## Examples

### Per-item mode for command materialization

```yaml
- name: agents.commands.claude
  step: kfg.materialize
  weight: -45
  env:
    MODE: "per-item"
    ASSETS: "kfg.extension.self.commands.git-commit"
    CONVERTER: "kfg.convert.self.command.claude"
    OUTPUTS: ".claude/commands/git-commit.md"
  when:
    output:
      step: kfg.detect-agent
      name: AGENT
      equals: "claude"
```

### Per-item mode for multi-item materialization

```yaml
- name: agents.commands.grouped
  step: kfg.materialize
  weight: -45
  env:
    MODE: "per-item"
    ASSETS: "kfg.extension.self.commands.git-commit:kfg.extension.self.commands.pr-review"
    CONVERTER: "kfg.convert.self.command.claude"
    OUTPUTS: ".claude/commands/git-commit.md:.claude/commands/pr-review.md"
```

### Aggregate mode for MCP aggregation

```yaml
- name: agents.mcp.claude
  step: kfg.materialize
  weight: -40
  env:
    MODE: "aggregate"
    ASSETS: "kfg.extension.ctx7.mcp:kfg.extension.chrome-devtools.mcp:kfg.extension.playwright.mcp"
    CONVERTER: "kfg.convert.self.mcp.claude"
    OUTPUTS: ".mcp.json"
    WRAP_KEY: "mcpServers"
  when:
    output:
      step: kfg.detect-agent
      name: AGENT
      equals: "claude"
```

### Aggregate mode without wrapper key

```yaml
- name: agents.settings.aggregate
  step: kfg.materialize
  weight: -63
  env:
    MODE: "aggregate"
    ASSETS: "asset.one:asset.two"
    CONVERTER: "some.converter"
    OUTPUTS: "output.yaml"
```