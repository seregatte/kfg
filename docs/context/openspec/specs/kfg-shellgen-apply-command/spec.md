# Apply Command Specification

## Purpose

Define the `kfg apply` command for applying kustomizations and generating shell code.

## Requirements

### Requirement: Apply command syntax

The CLI MUST provide `kfg apply` for shell generation with optional source via `KFG_KPATH`.

#### Scenario: Apply from kustomization
- **WHEN** user runs `kfg apply -k .kfg/overlay/dev`
- **THEN** loads kustomization from path
- **AND** resolves workflow
- **AND** generates shell functions to stdout

#### Scenario: Apply from file
- **WHEN** user runs `kfg apply -f manifest.yaml`
- **THEN** loads manifest from file
- **AND** generates shell functions

#### Scenario: Apply from stdin
- **WHEN** user runs `kfg apply -f -`
- **THEN** reads manifest from stdin
- **AND** generates shell functions

#### Scenario: Apply from GitHub URL
- **WHEN** user runs `kfg apply -k https://github.com/owner/repo//manifests`
- **THEN** clones GitHub repository
- **AND** processes kustomization.yaml
- **AND** resolves workflow
- **AND** generates shell functions

#### Scenario: Apply from GitHub URL with ref
- **WHEN** user runs `kfg apply -k https://github.com/owner/repo//manifests?ref=v1.0.0`
- **THEN** clones specified tag
- **AND** processes kustomization.yaml
- **AND** resolves workflow
- **AND** generates shell functions

#### Scenario: Apply without flags with KFG_KPATH
- **WHEN** user runs `kfg apply` (no `-k` or `-f`)
- **AND** `KFG_KPATH=./manifests` is set
- **THEN** uses `KFG_KPATH` as kustomize path
- **AND** generates shell functions

#### Scenario: Apply without flags or KFG_KPATH
- **WHEN** user runs `kfg apply` (no `-k` or `-f`)
- **AND** `KFG_KPATH` is not set
- **THEN** exit code 1
- **AND** error message: "kustomization source required. Provide a path, use -k flag, or set KFG_KPATH."

### Requirement: Apply flags

The CLI MUST support specific flags.

#### Scenario: Kustomize path
- **WHEN** user runs `kfg apply -k .kfg/overlay/dev`
- **THEN** short flag `-k` works same as `--kustomize`

#### Scenario: Manifest file
- **WHEN** user runs `kfg apply -f manifest.yaml`
- **THEN** short flag `-f` works same as `--file`

#### Scenario: Output file
- **WHEN** user runs `kfg apply -o output.sh`
- **THEN** writes shell code to file

#### Scenario: Workflow selection
- **WHEN** user runs `kfg apply -w dev`
- **THEN** uses specified workflow name

#### Scenario: Command filter
- **WHEN** user runs `kfg apply -c claude,gemini`
- **THEN** generates only specified commands

#### Scenario: Inline expression flag
- **WHEN** user runs `kfg apply -f manifest.yaml --convert my-asset --with '.data.key'`
- **THEN** short flag `--with` accepts a raw yq expression

### Requirement: Flag validation

The CLI MUST validate flag combinations.

#### Scenario: Required flag
- **WHEN** user runs `kfg apply` without `-k` or `-f`
- **THEN** exit code 2 (usage error)
- **AND** error message indicates required flag

#### Scenario: Mutual exclusion
- **WHEN** user runs `kfg apply -k path -f file`
- **THEN** exit code 2 (usage error)
- **AND** error message indicates flag conflict

#### Scenario: Conversion mode requires --convert with --use or --with
- **WHEN** user runs `kfg apply -f manifest.yaml --convert prod-servers`
- **THEN** exit code 2 (usage error)
- **AND** error message indicates `--convert` requires `--use` or `--with` to be specified

#### Scenario: --use requires --convert
- **WHEN** user runs `kfg apply -f manifest.yaml --use servers-to-json`
- **THEN** exit code 2 (usage error)
- **AND** error message indicates `--use` requires `--convert` to be specified

#### Scenario: --with requires --convert or stdin
- **WHEN** user runs `kfg apply -f manifest.yaml --with '.data'`
- **AND** no `--convert` flag is provided
- **THEN** exit code 2 (usage error)
- **AND** error message indicates `--with` requires `--convert` or `-f -` (stdin)

#### Scenario: --with and --use mutual exclusivity
- **WHEN** user runs `kfg apply -f manifest.yaml --convert my-asset --use my-converter --with '.expr'`
- **THEN** exit code 2 (usage error)
- **AND** error message indicates flags are mutually exclusive

#### Scenario: Conversion flags incompatible with shell flags
- **WHEN** user runs `kfg apply -f manifest.yaml --convert prod-servers --use servers-to-json -w dev`
- **THEN** exit code 2 (usage error)
- **AND** error message indicates `--workflow` cannot be used with `--convert`/`--use`

#### Scenario: Shell flags incompatible with conversion flags
- **WHEN** user runs `kfg apply -f manifest.yaml -w dev --convert prod-servers --use servers-to-json`
- **THEN** exit code 2 (usage error)
- **AND** error message indicates `--convert`/`--use` cannot be used with `--workflow`/`--cmds`

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

The `--convert` flag SHALL accept a raw string literal (JSON/YAML) as input when no Asset with the matching name exists in the loaded manifests. When an Asset with the matching name exists, it SHALL take precedence over raw string interpretation.

#### Scenario: Asset name resolution (existing behavior)
- **WHEN** user runs `kfg apply -f manifest.yaml --convert my-asset --use my-converter`
- **AND** an Asset with `metadata.name: my-asset` exists
- **THEN** the Asset is used for conversion

#### Scenario: Raw JSON string input
- **WHEN** user runs `kfg apply -f manifest.yaml --convert '{"server":{"command":"npx","args":["-y","pkg"]}}' --use my-converter`
- **AND** no Asset with matching name exists
- **THEN** the string is parsed as JSON input
- **AND** the Converter's expression is applied
- **AND** result is output

#### Scenario: Raw string with unknown asset and --use fails
- **WHEN** user runs `kfg apply -f manifest.yaml --convert nonexistent --use my-converter`
- **AND** no Asset with matching name exists
- **AND** input is not valid JSON or YAML
- **THEN** exit code 1
- **AND** error message indicates asset not found and input is not valid JSON or YAML

### Requirement: Inline expression with --with flag

The CLI SHALL provide a `--with` flag for specifying a raw yq expression without requiring a Converter resource.

#### Scenario: --with with Asset lookup
- **WHEN** user runs `kfg apply -f manifest.yaml --convert my-asset --with '.data | {"key": .value}'`
- **THEN** the Asset is loaded by name
- **AND** the inline expression is applied
- **AND** output is returned in YAML format

#### Scenario: --with with raw string input
- **WHEN** user runs `kfg apply -f manifest.yaml --convert '{"key":"value"}' --with '.key'`
- **THEN** the raw string is parsed as JSON input
- **AND** the inline expression is applied
- **AND** output is returned

### Requirement: Stdin pipeline with -f - and --with

When `-f -` is used with `--with`, stdin content SHALL be passed directly to the conversion engine without manifest parsing. This enables multi-document merge operations.

#### Scenario: Stdin data with inline expression
- **WHEN** user runs `echo '{"a":1}---{"b":2}' | kfg apply -f - --with 'select(di == 0) * select(di == 1)'`
- **THEN** stdin content is used as multi-document input
- **AND** the expression is applied
- **AND** merged result is output

#### Scenario: Stdin with -o output file
- **WHEN** user runs `echo '{"a":1}---{"b":2}' | kfg apply -f - --with 'select(di == 0) * select(di == 1)' -o merged.json`
- **THEN** the merged result is written to `merged.json`
- **AND** nothing is printed to stdout

### Requirement: Shell generation

The apply MUST generate valid shell code.

#### Scenario: Function structure
- **WHEN** apply succeeds
- **THEN** output defines bash functions for each cmd
- **AND** functions callable after sourcing

#### Scenario: Helper functions
- **WHEN** shell code generated
- **THEN** includes helper functions for logging
- **AND** includes helper for output management

#### Scenario: Build result global
- **WHEN** build result YAML present
- **THEN** build result setup at global scope
- **AND** shared across all cmd wrappers

### Requirement: Workflow resolution

The apply MUST resolve workflow correctly.

#### Scenario: Auto-detect workflow
- **WHEN** single workflow in manifests
- **THEN** uses that workflow automatically

#### Scenario: Workflow specified
- **WHEN** user runs with `-w dev`
- **THEN** uses specified workflow

#### Scenario: Workflow not found
- **WHEN** workflow doesn't exist
- **THEN** exit code 1
- **AND** error message lists available workflows

### Requirement: Command resolution

The apply MUST resolve cmds correctly.

#### Scenario: Command filter
- **WHEN** user runs with `-c claude`
- **THEN** generates only `claude` function

#### Scenario: Command not in workflow
- **WHEN** filtered cmd not in workflow
- **THEN** exit code 1
- **AND** error message lists available cmds

### Requirement: Exit codes

The CLI MUST use consistent exit codes.

#### Scenario: Success
- **WHEN** apply succeeds
- **THEN** exit code 0

#### Scenario: Runtime error
- **WHEN** resolution or generation fails
- **THEN** exit code 1

#### Scenario: Usage error
- **WHEN** invalid flag combination
- **THEN** exit code 2

#### Scenario: GitHub clone failure
- **WHEN** GitHub URL fails to clone
- **THEN** exit code 1
- **AND** error message indicates clone issue

#### Scenario: No source provided
- **WHEN** no `-k`, no `-f`, and no `KFG_KPATH`
- **THEN** exit code 1
- **AND** error message indicates source required

### Requirement: Apply Refresh Propagation

The apply command SHALL generate shell code that can force cacheable Steps to refresh.

#### Scenario: Apply with refresh flag
- **WHEN** user runs `kfg apply -k path --refresh`
- **THEN** the generated shell code SHALL export or embed refresh state equivalent to `KFG_REFRESH`
- **AND** cacheable Steps in that shell SHALL bypass matching cache entries when executed

#### Scenario: Apply without refresh flag
- **WHEN** user runs `kfg apply -k path` without `--refresh`
- **THEN** the generated shell code SHALL use cached Step entries when available
