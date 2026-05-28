# Convert Raw Input Specification

## Purpose

This specification defines the --convert flag's raw string input behavior and the --with flag for inline yq expressions.
## Requirements

### Requirement: --convert accepts raw string input

The `--convert` flag SHALL accept both an Asset `metadata.name` for lookup and a raw string literal (JSON or YAML). When no Asset with the given name exists in the loaded manifests, the value SHALL be treated as raw input data and passed directly to the conversion engine.

#### Scenario: Asset name resolution (existing behavior)
- **WHEN** user runs `kfg apply -f manifest.yaml --convert my.asset --use my.converter`
- **THEN** the system looks up an Asset with `metadata.name: my.asset`
- **AND** converts it using the Converter with `metadata.name: my.converter`
- **AND** outputs the result

#### Scenario: Raw string input fallback
- **WHEN** user runs `kfg apply -f manifest.yaml --convert '{"key":"value"}' --use my.converter`
- **AND** no Asset with `metadata.name` matching the string exists
- **THEN** the string is treated as raw input data
- **AND** the data is parsed according to the Converter's input format
- **AND** the expression is applied and output returned

#### Scenario: Asset name takes precedence over raw string
- **WHEN** an Asset exists with `metadata.name` equal to the `--convert` value
- **THEN** the Asset is used for conversion (not treated as raw string)
- **AND** available Asset names are listed in the error message if lookup fails before raw fallback

### Requirement: --with flag for inline yq expressions

The CLI SHALL provide a `--with` flag that accepts a raw yq expression string, bypassing Converter resource lookup. When `--with` is used, no Converter resource is required in the loaded manifests.

#### Scenario: Inline expression with Asset
- **WHEN** user runs `kfg apply -f manifest.yaml --convert my.asset --with '.data | {"key": .value}'`
- **THEN** the Asset is loaded by name
- **AND** the inline expression is applied to the Asset data
- **AND** output is produced in YAML format (default)

#### Scenario: Inline expression with raw string input
- **WHEN** user runs `kfg apply -f manifest.yaml --convert '{"key":"value"}' --with '.key'`
- **THEN** the raw string is parsed as JSON input
- **AND** the inline expression is applied
- **AND** output is returned

#### Scenario: --with and --use mutual exclusivity
- **WHEN** user runs `kfg apply -f manifest.yaml --convert my.asset --use my.converter --with '.expr'`
- **THEN** exit code 2 (usage error)
- **AND** error message indicates flags are mutually exclusive

### Requirement: Stdin multi-document processing with -f and --with

When `-f -` is used with `--with`, stdin content SHALL be passed directly to the yq-go engine without manifest parsing. This enables multi-document merge operations.

#### Scenario: Multi-document JSON merge from stdin
- **WHEN** user pipes two JSON documents separated by `---` to `kfg apply -f - --with 'select(di == 0) * select(di == 1)'`
- **THEN** both documents are parsed
- **AND** the deep merge expression is applied
- **AND** merged result is output to stdout

#### Scenario: Stdin with -o output file
- **WHEN** user runs `echo '{"a":1}---{"b":2}' | kfg apply -f - --with 'select(fi == 0) * select(fi == 1)' -o merged.json`
- **THEN** the merged result is written to `merged.json`
- **AND** nothing is printed to stdout

### Requirement: --with flag validation

The `--with` flag SHALL require either `--convert` or `-f -` (stdin) to be provided.

#### Scenario: --with without --convert or stdin
- **WHEN** user runs `kfg apply -f manifest.yaml --with '.data'`
- **AND** no `--convert` flag is provided
- **THEN** exit code 2 (usage error)
- **AND** error message indicates `--with` requires `--convert` or stdin

#### Scenario: --with with -f - (stdin, no --convert)
- **WHEN** user runs `echo "data" | kfg apply -f - --with '.key'`
- **THEN** stdin content is used as input
- **AND** the expression is applied
- **AND** result is output

### Requirement: Output format control with --with

When `--with` is used, the output format SHALL default to YAML. The `-o` flag controls output destination but not format. Format inference SHALL be based on the input type.

#### Scenario: JSON input with --with outputs JSON-like YAML
- **WHEN** user runs `kfg apply -f manifest.yaml --convert '{"name":"test"}' --with '.name'`
- **THEN** output is the string value in YAML format
- **AND** no YAML document markers are present for scalar output

#### Scenario: Multi-document merge outputs single document
- **WHEN** user merges two documents with `--with 'select(di == 0) * select(di == 1)'`
- **THEN** output is a single merged document
- **AND** no `---` separator appears in output
