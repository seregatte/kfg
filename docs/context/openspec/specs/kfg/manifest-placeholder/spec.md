# placeholder-resolution Specification

## Purpose
Placeholder resolution enables dynamic configuration values using environment variables. The resolver transforms `{env:VAR_NAME}` placeholders to shell variable syntax `$VAR_NAME`, allowing configurations to adapt to different environments without hardcoding sensitive values.

## Requirements

### Requirement: Placeholder pattern recognition

The resolver SHALL recognize `{env:VAR_NAME}` placeholders in string values where `VAR_NAME` matches the pattern `[A-Za-z_][A-Za-z0-9_]*`.

#### Scenario: Simple placeholder in string value
- **WHEN** data contains `"apiKey": "{env:EXP_API_KEY}"`
- **THEN** resolver SHALL transform to `"apiKey": "$EXP_API_KEY"`

#### Scenario: Placeholder with underscores and numbers
- **WHEN** data contains `"key": "{env:MY_VAR_123}"`
- **THEN** resolver SHALL transform to `"key": "$MY_VAR_123"`

#### Scenario: Invalid placeholder name
- **WHEN** data contains `"key": "{env:123_INVALID}"`
- **THEN** resolver SHALL NOT transform (pattern doesn't match)

### Requirement: Recursive resolution in nested structures

The resolver SHALL recursively process placeholders in nested maps and arrays.

#### Scenario: Placeholder in nested map
- **WHEN** data contains:
  ```yaml
  servers:
    main:
      apiKey: "{env:SERVER_API_KEY}"
  ```
- **THEN** resolver SHALL transform nested value to `"$SERVER_API_KEY"`

#### Scenario: Placeholder in array element
- **WHEN** data contains:
  ```yaml
  args:
    - "--api-key"
    - "{env:API_KEY}"
  ```
- **THEN** resolver SHALL transform array element to `"$API_KEY"`

#### Scenario: Mixed content in string
- **WHEN** data contains `"url": "https://{env:HOST}:8080"`
- **THEN** resolver SHALL transform to `"url": "https://$HOST:8080"`

### Requirement: Missing env var handling

The resolver SHALL resolve to empty string when env var is not set, and SHALL log a warning.

#### Scenario: Missing env var
- **WHEN** env var `MISSING_VAR` is not set
- **AND** data contains `"key": "{env:MISSING_VAR}"`
- **THEN** resolver SHALL transform to `"key": ""`
- **AND** resolver SHALL log warning: "env var MISSING_VAR not set"

#### Scenario: Empty env var
- **WHEN** env var `EMPTY_VAR` is set to empty string
- **AND** data contains `"key": "{env:EMPTY_VAR}"`
- **THEN** resolver SHALL transform to `"key": ""`

### Requirement: Resolution preserves non-placeholder content

The resolver SHALL NOT modify strings that do not contain placeholders.

#### Scenario: No placeholder in string
- **WHEN** data contains `"name": "server-name"`
- **THEN** resolver SHALL return unchanged `"name": "server-name"`

#### Scenario: Similar but invalid pattern
- **WHEN** data contains `"key": "{var:API_KEY}"` (not `env:`)
- **THEN** resolver SHALL NOT transform (wrong prefix)

### Requirement: Integration with converter engine

The converter engine SHALL apply placeholder resolution to Assets.data before template execution.

#### Scenario: Placeholder in MCP asset
- **WHEN** MCP asset contains `{env:EXP_CONTEXT7_API_KEY}` in args
- **AND** converter template uses `{{ .args }}`
- **THEN** converter SHALL receive args with `$EXP_CONTEXT7_API_KEY`
- **AND** output SHALL contain shell variable syntax

#### Scenario: YQ engine with placeholder
- **WHEN** asset data contains placeholder
- **AND** converter uses yq engine
- **THEN** yq SHALL receive resolved data

### Requirement: Integration with shell generation

Shell generation SHALL apply placeholder resolution to Step/Cmd env values.

#### Scenario: Step with env placeholder
- **WHEN** Step manifest contains:
  ```yaml
  spec:
    env:
      API_KEY: "{env:EXP_API_KEY}"
  ```
- **THEN** generated shell SHALL contain `export API_KEY="$EXP_API_KEY"`

#### Scenario: Cmd with env placeholder
- **WHEN** Cmd manifest contains:
  ```yaml
  spec:
    env:
      MODEL: "{env:KFG_MODEL}"
  ```
- **THEN** generated wrapper SHALL contain `export MODEL="$KFG_MODEL"`