## ADDED Requirements

### Requirement: Assets resource kind

The manifest model MUST support an Assets resource kind for declaring data payloads.

#### Scenario: Assets resource structure
- GIVEN a manifest file with `kind: Assets`
- WHEN the resource is parsed
- THEN `apiVersion` SHALL be `kfg.dev/v1alpha1`
- AND `kind` SHALL be `Assets`
- AND `metadata.name` SHALL be present
- AND `spec.data` SHALL be present

#### Scenario: Assets with YAML data (default)
- GIVEN an Assets resource without `spec.input.format`
- WHEN the resource is parsed
- THEN `spec.input.format` SHALL default to `yaml`
- AND `spec.data` SHALL be parsed as a YAML map

#### Scenario: Assets with non-YAML data
- GIVEN an Assets resource with `spec.input.format: json`
- WHEN the resource is parsed
- THEN `spec.data` SHALL be parsed as a string
- AND the string SHALL contain valid JSON

#### Scenario: Assets with TOML data
- GIVEN an Assets resource with `spec.input.format: toml`
- WHEN the resource is parsed
- THEN `spec.data` SHALL be parsed as a string
- AND the string SHALL contain valid TOML

### Requirement: Supported input formats

The system MUST support all yq-go input formats for Assets.

#### Scenario: All formats supported
- GIVEN an Assets resource
- WHEN `spec.input.format` is specified
- THEN the format MUST be one of: yaml, json, xml, props, csv, tsv, toml, hcl, lua, ini, shell, base64, uri, kyaml

#### Scenario: Invalid format
- GIVEN an Assets resource with unsupported `spec.input.format`
- WHEN the resource is validated
- THEN a validation error SHALL be raised
- AND the error SHALL list supported formats

### Requirement: Converter resource kind

The manifest model MUST support a Converter resource kind for declaring transformations.

#### Scenario: Converter resource structure
- GIVEN a manifest file with `kind: Converter`
- WHEN the resource is parsed
- THEN `apiVersion` SHALL be `kfg.dev/v1alpha1`
- AND `kind` SHALL be `Converter`
- AND `metadata.name` SHALL be present
- AND `spec.input.format` SHALL be present
- AND `spec.engine.expression` SHALL be present
- AND `spec.output.format` SHALL be present

#### Scenario: Converter input format
- GIVEN a Converter resource
- WHEN `spec.input.format` is specified
- THEN it SHALL be one of the supported input formats
- AND the default SHALL be `yaml`

#### Scenario: Converter output format
- GIVEN a Converter resource
- WHEN `spec.output.format` is specified
- THEN it SHALL be one of the supported output formats
- AND the default SHALL be `yaml`

### Requirement: Supported output formats

The system MUST support all yq-go output formats plus `raw`.

#### Scenario: All formats supported
- GIVEN a Converter resource
- WHEN `spec.output.format` is specified
- THEN the format MUST be one of: yaml, json, xml, props, csv, tsv, toml, hcl, lua, ini, shell, base64, uri, kyaml, raw

#### Scenario: Invalid format
- GIVEN a Converter resource with unsupported `spec.output.format`
- WHEN the resource is validated
- THEN a validation error SHALL be raised
- AND the error SHALL list supported formats

### Requirement: Assets validation

Assets resources MUST be validated.

#### Scenario: Missing name
- GIVEN an Assets resource without `metadata.name`
- WHEN the resource is validated
- THEN a validation error SHALL be raised

#### Scenario: Missing data
- GIVEN an Assets resource without `spec.data`
- WHEN the resource is validated
- THEN a validation error SHALL be raised

#### Scenario: Invalid format
- GIVEN an Assets resource with unsupported format
- WHEN the resource is validated
- THEN a validation error SHALL be raised

### Requirement: Converter validation

Converter resources MUST be validated.

#### Scenario: Missing name
- GIVEN a Converter resource without `metadata.name`
- WHEN the resource is validated
- THEN a validation error SHALL be raised

#### Scenario: Missing expression
- GIVEN a Converter resource without `spec.engine.expression`
- WHEN the resource is validated
- THEN a validation error SHALL be raised

#### Scenario: Invalid input format
- GIVEN a Converter resource with unsupported input format
- WHEN the resource is validated
- THEN a validation error SHALL be raised

#### Scenario: Invalid output format
- GIVEN a Converter resource with unsupported output format
- WHEN the resource is validated
- THEN a validation error SHALL be raised

### Requirement: Source kinds skipped in resolution

Assets and Converter MUST be skipped during resolution and shell generation.

#### Scenario: Assets not indexed
- GIVEN manifests containing Assets resources
- WHEN the resolver creates an index
- THEN Assets SHALL NOT be included in the index

#### Scenario: Converter not indexed
- GIVEN manifests containing Converter resources
- WHEN the resolver creates an index
- THEN Converter SHALL NOT be included in the index

#### Scenario: Shell generation unaffected
- GIVEN manifests containing Assets and Converter
- WHEN shell code is generated
- THEN Assets and Converter SHALL NOT affect shell output
