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
- **THEN** generated shell code SHALL enable step-scoped cache invalidation for cacheable Steps
- **AND** refreshed Steps SHALL rebuild their cache entries after successful execution

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
- **THEN** the runtime SHALL invalidate matching step cache entries during that invocation
- **AND** refreshed Steps SHALL rebuild their cache entries after successful execution

### Requirement: Help Documentation

The CLI MUST provide helpful documentation.

#### Scenario: Root help documents public environment variables
- **WHEN** user runs `kfg --help`
- **THEN** help output SHALL list public environment variables including `KFG_KPATH`, `KFG_STORE_DIR`, `KFG_VERBOSE`, `KFG_LOG_FILE`, `KFG_LOG_DIR`, `KFG_LOG_COLOR`, and `KFG_REFRESH`
- **AND** the `KFG_REFRESH` description SHALL explain step cache invalidation and rebuild semantics

#### Scenario: Apply help documents refresh and environment variables
- **WHEN** user runs `kfg apply --help`
- **THEN** help output SHALL document `--refresh`
- **AND** SHALL describe it as invalidating and rebuilding cache entries for cacheable Steps
- **AND** SHALL document `KFG_KPATH` and `KFG_REFRESH`

#### Scenario: Run help documents refresh and environment variables
- **WHEN** user runs `kfg run --help`
- **THEN** help output SHALL document `--refresh`
- **AND** SHALL describe it as invalidating and rebuilding cache entries for cacheable Steps
- **AND** SHALL document `KFG_KPATH` and `KFG_REFRESH`
