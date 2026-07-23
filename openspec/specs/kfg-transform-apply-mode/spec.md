# apply-conversion-mode Specification

## Purpose
TBD - created by archiving change add-assets-converter. Update Purpose after archive.
## Requirements
### Requirement: Conversion mode flags

The apply command MUST support `--convert` and `--use` flags for conversion mode.

#### Scenario: Convert flag
- GIVEN the apply command
- WHEN `--convert <asset-name>` is provided
- THEN the specified Asset SHALL be loaded from manifests
- AND the Asset SHALL be used as input for conversion

#### Scenario: Use flag
- GIVEN the apply command
- WHEN `--use <converter-name>` is provided
- THEN the specified Converter SHALL be loaded from manifests
- AND the Converter SHALL be used to transform the Asset

#### Scenario: Both flags required
- GIVEN the apply command
- WHEN only `--convert` is provided without `--use`
- THEN exit code 2 (usage error)
- AND error message SHALL indicate both flags are required

#### Scenario: Both flags required (reverse)
- GIVEN the apply command
- WHEN only `--use` is provided without `--convert`
- THEN exit code 2 (usage error)
- AND error message SHALL indicate both flags are required

### Requirement: Mutual exclusivity with shell mode

Conversion flags MUST NOT be mixed with shell generation flags.

#### Scenario: Convert with workflow flag
- GIVEN the apply command
- WHEN `--convert` and `--workflow` are both provided
- THEN exit code 2 (usage error)
- AND error message SHALL indicate flags are incompatible

#### Scenario: Convert with cmds flag
- GIVEN the apply command
- WHEN `--convert` and `--cmds` are both provided
- THEN exit code 2 (usage error)
- AND error message SHALL indicate flags are incompatible

#### Scenario: Use with workflow flag
- GIVEN the apply command
- WHEN `--use` and `--workflow` are both provided
- THEN exit code 2 (usage error)
- AND error message SHALL indicate flags are incompatible

#### Scenario: Use with cmds flag
- GIVEN the apply command
- WHEN `--use` and `--cmds` are both provided
- THEN exit code 2 (usage error)
- AND error message SHALL indicate flags are incompatible

### Requirement: Asset lookup

The system MUST find Assets by `metadata.name`.

#### Scenario: Asset found
- GIVEN manifests containing an Asset with `metadata.name: prod-servers`
- WHEN `--convert prod-servers` is provided
- THEN the Asset SHALL be found
- AND conversion SHALL proceed

#### Scenario: Asset not found
- GIVEN manifests without an Asset named `prod-servers`
- WHEN `--convert prod-servers` is provided
- THEN exit code 1
- AND error message SHALL indicate Asset not found
- AND available Assets SHALL be listed

### Requirement: Converter lookup

The system MUST find Converters by `metadata.name`.

#### Scenario: Converter found
- GIVEN manifests containing a Converter with `metadata.name: servers-to-json`
- WHEN `--use servers-to-json` is provided
- THEN the Converter SHALL be found
- AND conversion SHALL proceed

#### Scenario: Converter not found
- GIVEN manifests without a Converter named `servers-to-json`
- WHEN `--use servers-to-json` is provided
- THEN exit code 1
- AND error message SHALL indicate Converter not found
- AND available Converters SHALL be listed

### Requirement: Conversion pipeline

The apply command MUST execute the conversion pipeline correctly.

#### Scenario: Successful conversion
- GIVEN a valid Asset and Converter
- WHEN conversion is executed
- THEN the Asset data SHALL be loaded
- AND the Converter expression SHALL be evaluated
- AND the result SHALL be output in the Converter's format
- AND exit code SHALL be 0

#### Scenario: Output to file
- GIVEN a valid Asset and Converter
- WHEN `-o output.yaml` is provided
- THEN the conversion result SHALL be written to the specified file
- AND exit code SHALL be 0

#### Scenario: Output to stdout
- GIVEN a valid Asset and Converter
- WHEN no `-o` flag is provided
- THEN the conversion result SHALL be written to stdout
- AND exit code SHALL be 0

### Requirement: Conversion error handling

The system MUST handle conversion errors gracefully.

#### Scenario: Expression evaluation error
- GIVEN a Converter with an invalid yq expression
- WHEN conversion is executed
- THEN exit code 1
- AND error message SHALL indicate the expression error

#### Scenario: Format conversion error
- GIVEN an Asset with data incompatible with Converter's input format
- WHEN conversion is executed
- THEN exit code 1
- AND error message SHALL indicate the format conversion failure

### Requirement: Apply help text

The apply command MUST document both modes in help text.

#### Scenario: Help output
- GIVEN the apply command
- WHEN `kfg apply --help` is executed
- THEN help text SHALL describe shell mode (default)
- AND help text SHALL describe conversion mode (--convert + --use)
- AND examples for both modes SHALL be provided
- AND flag incompatibilities SHALL be documented

