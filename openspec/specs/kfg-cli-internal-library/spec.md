# Internal Library Specification

## Purpose

kfg integrates with Cobra CLI framework for command-line parsing and Viper for configuration management.

## Requirements

### Requirement: Cobra CLI Framework Integration

The CLI MUST use Cobra framework for command-line parsing.

#### Scenario: Root command exists
- **WHEN** kfg binary built
- **THEN** root cobra command configured
- **AND** root command has version and help information

#### Scenario: Subcommand structure
- **WHEN** CLI executed
- **THEN** subcommands registered: build, apply, log, store
- **AND** future subcommands can be added

### Requirement: Auto-generated Help

The CLI MUST provide auto-generated help messages.

#### Scenario: Root command help
- **WHEN** user runs `kfg --help`
- **THEN** usage printed to stdout
- **AND** all subcommands listed
- **AND** global flags documented

#### Scenario: Subcommand help
- **WHEN** user runs `kfg apply --help`
- **THEN** apply-specific usage printed
- **AND** all flags documented
- **AND** examples provided

### Requirement: Flag Management

The CLI MUST use Cobra's flag management.

#### Scenario: PersistentFlags for global options
- **WHEN** flag available to command and subcommands
- **THEN** registered using cmd.PersistentFlags()
- **AND** inherited by child commands

#### Scenario: Persistent flag examples
- **GIVEN** --verbose, --output global flags
- **WHEN** defined on root
- **THEN** accessible from all commands

#### Scenario: LocalFlags for command-specific
- **WHEN** flag specific to single command
- **THEN** registered using cmd.Flags()
- **AND** not available to other commands

#### Scenario: Local flag examples
- **GIVEN** --kustomize, --workflow for apply
- **WHEN** defined on apply
- **THEN** not available on root or other commands

#### Scenario: MarkFlagRequired
- **WHEN** flag must be provided
- **THEN** cmd.MarkFlagRequired called
- **AND** Cobra validates before execution

### Requirement: Nested Subcommands

The CLI MUST support nested subcommand structure.

#### Scenario: store subcommands
- **WHEN** user runs `kfg store image build`
- **THEN** nested commands work correctly

#### Scenario: Nested help
- **WHEN** user runs `kfg store image --help`
- **THEN** image subcommand help displayed

### Requirement: Command Execution Hooks

The CLI MUST use Cobra's hooks.

#### Scenario: PersistentPreRun
- **WHEN** command executed
- **THEN** PersistentPreRun executes first
- **AND** can set environment variables

#### Scenario: Example: verbose flag
- **WHEN** --verbose flag changed
- **THEN** PersistentPreRun sets KFG_VERBOSE
- **AND** logger reinitialized

### Requirement: Exit Code Handling

The CLI MUST handle exit codes correctly.

#### Scenario: SilenceErrors enabled
- **WHEN** command fails
- **THEN** custom error handling used
- **AND** appropriate exit code set

#### Scenario: Usage error
- **WHEN** invalid arguments
- **THEN** exit code 2

#### Scenario: Runtime error
- **WHEN** execution fails
- **THEN** exit code 1

### Requirement: Subcommand Extensibility

The CLI MUST support easy addition of subcommands.

#### Scenario: Adding new subcommand
- **WHEN** developer adds new subcommand
- **THEN** automatically available in help
- **AND** integrates with global flags