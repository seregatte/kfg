# Converter Execution Specification

## Purpose

This specification defines how Converter resources are applied to Asset data to produce transformed output — both via the CLI (`kfg apply --convert`) and within Step execution via `$KFG_BUILD_RESULT_FILE`.

## ADDED Requirements

### Requirement: Step Applies Converter Expression

Steps SHALL apply Converter expressions to Asset data using yq at runtime.

#### Scenario: Step extracts converter expression
- **WHEN** a Step needs to apply a Converter
- **THEN** it SHALL extract the expression from `$KFG_BUILD_RESULT_FILE` using `yq 'select(.kind == "Converter" and .metadata.name == "name") | .spec.engine.expression'`
- AND the expression SHALL be stored in a shell variable

#### Scenario: Step extracts asset data
- **WHEN** a Step needs Asset data
- **THEN** it SHALL extract the data from `$KFG_BUILD_RESULT_FILE` using `yq 'select(.kind == "Assets" and .metadata.name == "name") | .spec.data'`
- AND the data SHALL be stored in a shell variable

#### Scenario: Step applies expression to data
- **WHEN** a Step pipes asset data through a converter expression
- **THEN** it SHALL use the pattern `echo "$data" | yq "$expression"`
- AND the output SHALL be written to the target file path

#### Scenario: Missing converter error
- **WHEN** the extracted expression is empty or `null`
- **THEN** the Step SHALL log an error via `_kfg.log.error`
- AND the Step SHALL exit with non-zero status

### Requirement: Converter Expression Patterns

Converter expressions SHALL follow specific patterns for different transformation types.

#### Scenario: Command transformation to markdown
- **WHEN** a Converter transforms command data to agent-specific markdown
- **THEN** the expression SHALL use string concatenation with yq
- AND it SHALL include frontmatter markers (`---`)
- AND it SHALL conditionally include optional fields via `if .field then ... else "" end`
- AND the output format SHALL be `raw`

#### Scenario: MCP transformation to JSON
- **WHEN** a Converter transforms MCP data to agent-specific JSON
- **THEN** the expression SHALL use yq `map()`, `select()`, and dynamic keys `{(.name): {...}}`
- AND the output format SHALL be `json`
- AND disabled MCP entries SHALL be filtered via `select(.enabled == true)`

#### Scenario: Subagent transformation to markdown
- **WHEN** a Converter transforms subagent data to agent-specific markdown
- **THEN** the expression SHALL use string concatenation with yq
- AND it SHALL include frontmatter with name, description, and optional model/tools/permission
- AND complex fields SHALL be serialized via `@json`
- AND the output format SHALL be `raw`

### Requirement: Step Parameterization

Steps SHALL receive all configuration via environment variables passed from workflow references.

#### Scenario: Step receives asset prefix
- **WHEN** a Step processes multiple Assets
- **THEN** it SHALL receive `ASSET_PREFIX` via env
- AND it SHALL filter Assets by name prefix using `grep "^${ASSET_PREFIX}\."`

#### Scenario: Step receives converter prefix
- **WHEN** a Step applies a Converter
- **THEN** it SHALL receive `CONVERTER_PREFIX` via env
- AND it SHALL construct the Converter name as `${CONVERTER_PREFIX}.${NIXAI_AGENT}`

#### Scenario: Step receives output configuration
- **WHEN** a Step writes output
- **THEN** it SHALL receive `OUTPUT_DIR` (for commands/subagents) or `OUTPUT_FILE` (for MCP) via env
- AND it SHALL receive `FILE_EXT` and `NAME_PREFIX` for filename construction

### Requirement: Placeholder Resolution in Converter Expressions

Converter expressions SHALL support `{env:VAR}` placeholders.

#### Scenario: env placeholder in expression
- **WHEN** a Converter expression contains `{env:VAR_NAME}`
- **THEN** it SHALL be resolved to `$VAR_NAME` at generation time
- AND it SHALL be expanded by the shell at runtime

#### Scenario: Missing environment variable
- **WHEN** `{env:VAR_NAME}` references an unset environment variable
- **THEN** it SHALL resolve to an empty string at runtime
- AND a warning SHALL be logged
