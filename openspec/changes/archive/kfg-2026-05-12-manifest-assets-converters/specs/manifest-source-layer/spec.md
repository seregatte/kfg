# Manifest Source Layer Specification

## Purpose

This specification defines the Assets and Converter resource kinds — the "Source Layer" of the manifest model. These kinds declare data payloads and transformation rules but are not part of the shell execution pipeline. They are consumed by Steps via `$KFG_BUILD_RESULT_FILE` or by the CLI via `kfg apply --convert`.

## ADDED Requirements

### Requirement: Assets Resource Kind

The manifest model SHALL support an `Assets` resource kind for declaring structured data payloads.

#### Scenario: Assets schema validation
- **WHEN** an Assets resource is validated
- **THEN** `apiVersion` SHALL be `kfg.dev/v1alpha1`
- AND `kind` SHALL be `Assets`
- AND `metadata.name` SHALL be present and follow namespace naming rules
- AND `spec.data` SHALL be present

#### Scenario: Assets input format
- **WHEN** an Assets resource is parsed
- **THEN** `spec.input.format` SHALL default to `yaml` if not specified
- AND `spec.input.format` SHALL be one of the supported formats: yaml, json, xml, props, csv, tsv, toml, hcl, lua, ini, shell, base64, uri, kyaml

#### Scenario: Assets data payload
- **WHEN** an Assets resource is loaded
- **THEN** `spec.data` SHALL contain the full data payload
- AND for `yaml` input format, `spec.data` SHALL be a map structure
- AND for non-YAML input formats, `spec.data` SHALL be a string

#### Scenario: Assets are skipped during resolution
- **WHEN** the resolution pipeline indexes parsed resources
- **THEN** Assets resources SHALL be excluded from the execution index
- **AND** Assets SHALL be included in the YAML build result (`$KFG_BUILD_RESULT_FILE`)
- **AND** Assets SHALL be accessible to Steps via `yq` queries against `$KFG_BUILD_RESULT_FILE`

#### Scenario: Assets naming convention
- **WHEN** an Assets resource name is assigned
- **THEN** it SHALL follow the pattern `kfg.<layer>.<category>.<name>`
- **AND** examples include: `kfg.agent.claude`, `kfg.extension.ctx7.mcp`, `kfg.extension.self.commands.git-commit`

### Requirement: Converter Resource Kind

The manifest model SHALL support a `Converter` resource kind for declaring data transformations using yq-go expressions.

#### Scenario: Converter schema validation
- **WHEN** a Converter resource is validated
- **THEN** `apiVersion` SHALL be `kfg.dev/v1alpha1`
- AND `kind` SHALL be `Converter`
- AND `metadata.name` SHALL be present and follow namespace naming rules
- AND `spec.engine.expression` SHALL be present and non-empty
- AND `spec.input.format` SHALL be present
- AND `spec.output.format` SHALL be present

#### Scenario: Converter input/output format validation
- **WHEN** a Converter resource is parsed
- **THEN** `spec.input.format` SHALL be one of the supported input formats
- AND `spec.output.format` SHALL be one of the supported output formats
- AND `spec.output.format` MAY be `raw` for plain text output
- AND `raw` SHALL be output-only (not valid as input format)

#### Scenario: Converter expression evaluation
- **WHEN** a Converter's expression is applied to an Asset's data
- **THEN** the yq-go engine SHALL evaluate the expression against the input data
- AND the output format SHALL be applied to the result (json, yaml, raw, etc.)
- AND `raw` output SHALL return the expression result as plain text

#### Scenario: Converters are skipped during resolution
- **WHEN** the resolution pipeline indexes parsed resources
- **THEN** Converter resources SHALL be excluded from the execution index
- **AND** Converters SHALL be included in the YAML build result
- **AND** Converters SHALL be accessible to Steps via `yq` queries against `$KFG_BUILD_RESULT_FILE`

#### Scenario: Converter naming convention
- **WHEN** a Converter resource name is assigned
- **THEN** it SHALL follow the pattern `kfg.convert.<extension>.<type>.<agent>`
- **AND** examples include: `kfg.convert.self.command.claude`, `kfg.convert.self.mcp.opencode`

### Requirement: Cross-Resource Reference (Asset to Converter)

Steps SHALL be able to reference both Assets and Converters from the build result.

#### Scenario: Step reads Asset data from build result
- **WHEN** a Step executes with `$KFG_BUILD_RESULT_FILE` available
- **THEN** the Step SHALL be able to extract Asset data via `yq 'select(.kind == "Assets" and .metadata.name == "name") | .spec.data'`
- AND the extracted data SHALL be in the format declared by `spec.input.format`

#### Scenario: Step reads Converter expression from build result
- **WHEN** a Step executes with `$KFG_BUILD_RESULT_FILE` available
- **THEN** the Step SHALL be able to extract a Converter expression via `yq 'select(.kind == "Converter" and .metadata.name == "name") | .spec.engine.expression'`
- AND the extracted expression SHALL be a valid yq-go expression string

#### Scenario: Step applies Converter to Asset
- **WHEN** a Step pipes Asset data through a Converter expression
- **THEN** the command `echo "$data" | yq "$expression"` SHALL produce transformed output
- AND the output format SHALL match the Converter's `spec.output.format`

### Requirement: Assets Supported Input Formats

The following input formats SHALL be supported for Assets resources.

#### Scenario: YAML format (default)
- **WHEN** `spec.input.format` is `yaml` or omitted
- **THEN** `spec.data` SHALL be parsed as a YAML map
- AND the data SHALL be queryable via yq dot notation

#### Scenario: JSON format
- **WHEN** `spec.input.format` is `json`
- **THEN** `spec.data` SHALL be a JSON string
- AND it SHALL be converted to YAML for internal processing

#### Scenario: TOML format
- **WHEN** `spec.input.format` is `toml`
- **THEN** `spec.data` SHALL be a TOML string
- AND it SHALL be converted to YAML for internal processing

### Requirement: Converters Supported Output Formats

The following output formats SHALL be supported for Converter resources.

#### Scenario: JSON output
- **WHEN** `spec.output.format` is `json`
- **THEN** the Converter result SHALL be serialized as JSON

#### Scenario: YAML output
- **WHEN** `spec.output.format` is `yaml`
- **THEN** the Converter result SHALL be serialized as YAML

#### Scenario: Raw output
- **WHEN** `spec.output.format` is `raw`
- **THEN** the Converter result SHALL be returned as plain text
- AND no serialization SHALL be applied
- AND the expression result SHALL be written directly to stdout
