# converter-engine Specification

## Purpose
TBD - created by archiving change add-assets-converter. Update Purpose after archive.
## Requirements
### Requirement: yq-go engine integration

The system MUST use yq-go as a Go module for expression evaluation.

#### Scenario: Engine initialization
- GIVEN a Converter resource with an expression
- WHEN the engine is invoked
- THEN yq-go SHALL be initialized as an in-process library
- AND no external CLI dependency SHALL be required

#### Scenario: Expression evaluation
- GIVEN an Asset with data and a Converter with an expression
- WHEN the engine evaluates the expression
- THEN the expression SHALL be evaluated against the Asset data
- AND the result SHALL be returned in the Converter's output format

### Requirement: Format conversion

The engine MUST convert Asset data to Converter's expected input format when they differ.

#### Scenario: Same format
- GIVEN an Asset with `spec.input.format: yaml`
- AND a Converter with `spec.input.format: yaml`
- WHEN the engine processes the Asset
- THEN no format conversion SHALL be needed
- AND the data SHALL be passed directly to yq

#### Scenario: Different formats
- GIVEN an Asset with `spec.input.format: json`
- AND a Converter with `spec.input.format: yaml`
- WHEN the engine processes the Asset
- THEN the Asset data SHALL be converted from JSON to YAML
- AND the converted data SHALL be passed to yq

#### Scenario: Conversion error
- GIVEN an Asset with invalid data for its declared format
- WHEN the engine attempts conversion
- THEN an error SHALL be raised
- AND the error SHALL indicate the conversion failure

### Requirement: Output format encoding

The engine MUST encode yq results in the Converter's specified output format.

#### Scenario: YAML output
- GIVEN a Converter with `spec.output.format: yaml`
- WHEN the engine produces a result
- THEN the result SHALL be encoded as YAML
- AND the output SHALL be written to stdout

#### Scenario: JSON output
- GIVEN a Converter with `spec.output.format: json`
- WHEN the engine produces a result
- THEN the result SHALL be encoded as JSON
- AND the output SHALL be written to stdout

#### Scenario: TOML output
- GIVEN a Converter with `spec.output.format: toml`
- WHEN the engine produces a result
- THEN the result SHALL be encoded as TOML
- AND the output SHALL be written to stdout

#### Scenario: XML output
- GIVEN a Converter with `spec.output.format: xml`
- WHEN the engine produces a result
- THEN the result SHALL be encoded as XML
- AND the output SHALL be written to stdout

### Requirement: Raw output format

The engine MUST support `raw` as a special output format.

#### Scenario: Raw output with array result
- GIVEN a Converter with `spec.output.format: raw`
- WHEN the yq expression produces an array of strings
- THEN the array elements SHALL be joined with newlines
- AND the result SHALL be output as plain text without YAML/JSON encoding

#### Scenario: Raw output with string result
- GIVEN a Converter with `spec.output.format: raw`
- WHEN the yq expression produces a single string
- THEN the string SHALL be output as-is
- AND no encoding SHALL be applied

#### Scenario: Raw output with non-string result
- GIVEN a Converter with `spec.output.format: raw`
- WHEN the yq expression produces a non-string result (map, number, etc.)
- THEN the result SHALL be converted to string representation
- AND output as plain text

### Requirement: Error handling

The engine MUST handle errors gracefully.

#### Scenario: Invalid expression
- GIVEN a Converter with invalid yq expression
- WHEN the engine evaluates the expression
- THEN an error SHALL be raised
- AND the error SHALL indicate the expression syntax error

#### Scenario: Expression evaluation failure
- GIVEN a valid yq expression that fails at runtime
- WHEN the engine evaluates the expression
- THEN an error SHALL be raised
- AND the error SHALL indicate the evaluation failure

#### Scenario: Unsupported output format
- GIVEN a Converter with unsupported output format
- WHEN the engine attempts to encode output
- THEN an error SHALL be raised
- AND the error SHALL list supported formats

### Requirement: Data serialization

The engine MUST serialize Asset data correctly for yq consumption.

#### Scenario: YAML data serialization
- GIVEN an Asset with YAML data (map)
- WHEN the engine serializes for yq
- THEN the data SHALL be serialized as YAML string
- AND the serialization SHALL preserve data types

#### Scenario: String data serialization
- GIVEN an Asset with string data (non-YAML format)
- WHEN the engine serializes for yq
- THEN the data SHALL be passed as-is to yq
- AND the input decoder SHALL match the declared format

