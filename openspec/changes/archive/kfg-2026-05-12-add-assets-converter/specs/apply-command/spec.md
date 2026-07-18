## MODIFIED Requirements

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
