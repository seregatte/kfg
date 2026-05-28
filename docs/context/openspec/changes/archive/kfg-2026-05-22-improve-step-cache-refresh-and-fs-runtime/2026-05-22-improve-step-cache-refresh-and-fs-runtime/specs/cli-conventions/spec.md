## MODIFIED Requirements

### Requirement: Apply Command

The CLI MUST provide `kfg apply` for shell generation.

#### Scenario: Basic apply
- **GIVEN** user wants to generate shell integration
- **WHEN** user runs `kfg apply -k .kfg/overlay/dev`
- **THEN** bash-compatible shell code written to stdout
- **AND** code is valid bash and can be sourced
- **AND** no shell dotfiles modified

#### Scenario: Apply with refresh
- **WHEN** user runs `kfg apply -k path --refresh`
- **THEN** generated shell code SHALL enable runtime cache bypass for cacheable Steps
- **AND** refreshed cacheable Steps SHALL rebuild and overwrite their stored cache entries after successful execution

### Requirement: Run Command

The CLI MUST provide `kfg run` for one-shot agent execution.

#### Scenario: Basic run
- **GIVEN** user wants to run an agent
- **WHEN** user runs `kfg run -k .kfg/overlay/dev claude`
- **THEN** generates shell code, sources it, and executes the agent
- **AND** agent runs with inherited stdin/stdout/stderr
- **AND** exits with the agent's exit code

#### Scenario: Run with refresh
- **WHEN** user runs `kfg run -k path claude --refresh`
- **THEN** the runtime SHALL bypass matching cached Step entries during that invocation
- **AND** refreshed cacheable Steps SHALL rebuild and overwrite their stored cache entries after successful execution

### Requirement: Internal GC Commands

The CLI MUST provide `kfg sys gc` for internal runtime cache management.

#### Scenario: List cache entries
- **WHEN** user runs `kfg sys gc ls`
- **THEN** the command SHALL list cached runtime entries

#### Scenario: Inspect cache entry
- **WHEN** user runs `kfg sys gc inspect <id>`
- **THEN** the command SHALL print metadata for the specified cache entry

#### Scenario: Remove cache entry
- **WHEN** user runs `kfg sys gc rm <id>`
- **THEN** the command SHALL remove the specified cache entry

#### Scenario: Prune cache entries
- **WHEN** user runs `kfg sys gc prune`
- **THEN** the command SHALL remove cache entries according to the implemented prune policy

#### Scenario: Show cache disk usage
- **WHEN** user runs `kfg sys gc du`
- **THEN** the command SHALL report disk usage for cached runtime entries

## ADDED Requirements

### Requirement: Internal filesystem commands

The CLI MUST provide `kfg sys fs` for internal runtime filesystem inspection.

#### Scenario: Snapshot command
- **WHEN** internal runtime code runs `kfg sys fs snapshot <path> --maxdepth N`
- **THEN** the command SHALL print normalized relative paths rooted at `<path>`
- **AND** it SHALL fail for negative `--maxdepth` values

#### Scenario: Diff command
- **WHEN** internal runtime code runs `kfg sys fs diff --before <snapshot> --after <snapshot>`
- **THEN** the command SHALL print only paths newly present in `after`
