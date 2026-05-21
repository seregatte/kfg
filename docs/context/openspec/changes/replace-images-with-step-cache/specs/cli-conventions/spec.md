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

### Requirement: Log Command

The CLI MUST provide `kfg sys log` for structured logging used by generated shell helpers.

#### Scenario: Shell helper logging
- **GIVEN** a generated shell helper emits a structured log
- **WHEN** it delegates to the CLI logging command
- **THEN** `kfg sys log` SHALL accept the log event and persist it using the configured logger behavior

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

### Requirement: Store Command Surface

The CLI MUST NOT expose image or workspace commands.

#### Scenario: No image command
- **WHEN** user runs `kfg image`
- **THEN** the CLI SHALL report an unknown command usage error

#### Scenario: No workspace command
- **WHEN** user runs `kfg workspace`
- **THEN** the CLI SHALL report an unknown command usage error

### Requirement: Help Documentation

The CLI MUST provide helpful documentation.

#### Scenario: Root help documents public environment variables
- **WHEN** user runs `kfg --help`
- **THEN** help output SHALL list public environment variables including `KFG_KPATH`, `KFG_STORE_DIR`, `KFG_VERBOSE`, `KFG_LOG_FILE`, `KFG_LOG_DIR`, `KFG_LOG_COLOR`, and `KFG_REFRESH`

#### Scenario: Apply help documents refresh and environment variables
- **WHEN** user runs `kfg apply --help`
- **THEN** help output SHALL document `--refresh`
- **AND** SHALL document `KFG_KPATH` and `KFG_REFRESH`

#### Scenario: Run help documents refresh and environment variables
- **WHEN** user runs `kfg run --help`
- **THEN** help output SHALL document `--refresh`
- **AND** SHALL document `KFG_KPATH` and `KFG_REFRESH`
