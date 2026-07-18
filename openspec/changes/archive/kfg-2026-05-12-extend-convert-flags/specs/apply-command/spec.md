## ADDED Requirements

### Requirement: Conversion mode with --convert and --use

The CLI SHALL provide a conversion mode accessible via `--convert` and `--use` flags that transforms Asset data using a Converter resource.

#### Scenario: Basic conversion
- **WHEN** user runs `kfg apply -f manifest.yaml --convert my-asset --use my-converter`
- **THEN** the Asset with `metadata.name: my-asset` is found
- **AND** the Converter with `metadata.name: my-converter` is found
- **AND** the Converter's expression is applied to the Asset data
- **AND** result is output to stdout in the Converter's specified format

#### Scenario: Conversion with output file
- **WHEN** user runs `kfg apply -f manifest.yaml --convert my-asset --use my-converter -o output.json`
- **THEN** conversion result is written to `output.json`
- **AND** nothing is printed to stdout

#### Scenario: Conversion mode mutual exclusivity with shell generation
- **WHEN** user runs `kfg apply -f manifest.yaml --convert my-asset --use my-converter -w my-workflow`
- **THEN** exit code 2 (usage error)
- **AND** error message indicates conversion mode cannot be used with workflow selection

### Requirement: Raw string input for --convert

The `--convert` flag SHALL accept a raw string literal (JSON/YAML) as input when no Asset with the matching name exists in the loaded manifests.

#### Scenario: Raw JSON string input
- **WHEN** user runs `kfg apply -f manifest.yaml --convert '{"server":{"command":"npx","args":["-y","pkg"]}}' --use my-converter`
- **AND** no Asset with matching name exists
- **THEN** the string is parsed as JSON input
- **AND** the Converter's expression is applied
- **AND** result is output

### Requirement: Inline expression with --with flag

The CLI SHALL provide a `--with` flag for specifying a raw yq expression without requiring a Converter resource.

#### Scenario: --with with Asset lookup
- **WHEN** user runs `kfg apply -f manifest.yaml --convert my-asset --with '.data | {"key": .value}'`
- **THEN** the Asset is loaded by name
- **AND** the inline expression is applied
- **AND** output is returned in YAML format

#### Scenario: --with incompatible with --use
- **WHEN** user runs `kfg apply -f manifest.yaml --convert my-asset --use my-converter --with '.expr'`
- **THEN** exit code 2 (usage error)

### Requirement: Stdin pipeline with -f - and --with

When `-f -` is used with `--with`, stdin content SHALL be passed directly to the conversion engine without manifest parsing.

#### Scenario: Stdin data with inline expression
- **WHEN** user runs `echo '{"a":1}---{"b":2}' | kfg apply -f - --with 'select(fi == 0) * select(fi == 1)'`
- **THEN** stdin content is used as multi-document input
- **AND** the expression is applied
- **AND** merged result is output
