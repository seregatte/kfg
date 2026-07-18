## MODIFIED Requirements

### Requirement: Apply Command

The CLI MUST provide `kfg apply` for shell generation.

#### Scenario: Basic apply
- **GIVEN** user wants to generate shell integration
- **WHEN** user runs `kfg apply -k .kfg/overlay/dev`
- **THEN** bash-compatible shell code written to stdout
- **AND** code is valid bash and can be sourced
- **AND** no shell dotfiles modified

#### Scenario: Apply from file
- **GIVEN** user has manifest file
- **WHEN** user runs `kfg apply -f manifest.yaml`
- **THEN** shell code generated from file

#### Scenario: Apply with workflow
- **GIVEN** user wants specific workflow
- **WHEN** user runs `kfg apply -k path -w dev`
- **THEN** uses specified workflow

#### Scenario: Apply with command filter
- **GIVEN** user wants specific commands
- **WHEN** user runs `kfg apply -k path -c claude,gemini`
- **THEN** generates only specified commands

#### Scenario: Apply with refresh
- **WHEN** user runs `kfg apply -k path --refresh`
- **THEN** generated shell code SHALL bypass existing cache entries for cacheable Steps during execution
- **AND** refreshed cacheable Steps SHALL rebuild and overwrite their stored cache entries after successful execution

### Requirement: Run Command

The CLI MUST provide `kfg run` for one-shot agent execution.

#### Scenario: Basic run
- **GIVEN** user wants to run an agent
- **WHEN** user runs `kfg run -k .kfg/overlay/dev claude`
- **THEN** generates shell code, sources it, and executes the agent
- **AND** agent runs with inherited stdin/stdout/stderr
- **AND** exits with the agent's exit code

#### Scenario: Run with extra args
- **WHEN** user runs `kfg run -k path claude -- --model gpt-4`
- **THEN** passes `--model gpt-4` to the agent
- **AND** does not pass the `--` separator itself to the agent

#### Scenario: Run discovery
- **WHEN** user runs `kfg run -k path` without agent name
- **THEN** lists all available agents

#### Scenario: Run with refresh
- **WHEN** user runs `kfg run -k path claude --refresh`
- **THEN** the runtime SHALL bypass matching cached Step entries during that invocation
- **AND** refreshed cacheable Steps SHALL rebuild and overwrite their stored cache entries after successful execution

## ADDED Requirements

### Requirement: Internal filesystem command group

The CLI MUST provide `kfg sys fs` for internal runtime filesystem inspection.

#### Scenario: Snapshot command
- **WHEN** internal runtime code runs `kfg sys fs snapshot <path> --maxdepth N`
- **THEN** the command SHALL print normalized relative paths rooted at `<path>`
- **AND** `--maxdepth 0` SHALL mean no depth limit

#### Scenario: Diff command
- **WHEN** internal runtime code runs `kfg sys fs diff --before <snapshot> --after <snapshot>`
- **THEN** the command SHALL print only paths newly present in `after`
