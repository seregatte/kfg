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

#### Scenario: Conversion mode requires both flags
- **WHEN** user runs `kfg apply -f manifest.yaml --convert prod-servers`
- **THEN** exit code 2 (usage error)
- **AND** error message indicates `--convert` and `--use` must be used together

#### Scenario: Conversion mode requires both flags (reverse)
- **WHEN** user runs `kfg apply -f manifest.yaml --use servers-to-json`
- **THEN** exit code 2 (usage error)
- **AND** error message indicates `--convert` and `--use` must be used together

#### Scenario: Conversion flags incompatible with shell flags
- **WHEN** user runs `kfg apply -f manifest.yaml --convert prod-servers --use servers-to-json -w dev`
- **THEN** exit code 2 (usage error)
- **AND** error message indicates `--workflow` cannot be used with `--convert`/`--use`

#### Scenario: Shell flags incompatible with conversion flags
- **WHEN** user runs `kfg apply -f manifest.yaml -w dev --convert prod-servers --use servers-to-json`
- **THEN** exit code 2 (usage error)
- **AND** error message indicates `--convert`/`--use` cannot be used with `--workflow`/`--cmds`

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

