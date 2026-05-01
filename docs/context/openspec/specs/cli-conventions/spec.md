# CLI Conventions Specification

## Purpose

The kfg CLI is a declarative shell compiler that generates bash functions from YAML manifests. This specification defines conventions for CLI commands and flags.

## Requirements

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

### Requirement: Launch Command

The CLI MUST provide `kfg launch` for one-shot agent execution.

#### Scenario: Basic launch
- **GIVEN** user wants to run an agent
- **WHEN** user runs `kfg launch -k .kfg/overlay/dev claude`
- **THEN** generates shell code, sources it, and executes the agent
- **AND** agent runs with inherited stdin/stdout/stderr
- **AND** exits with the agent's exit code

#### Scenario: Launch with extra args
- **WHEN** user runs `kfg launch -k path claude -- --model gpt-4`
- **THEN** passes `--model gpt-4` to the agent

#### Scenario: Launch discovery
- **WHEN** user runs `kfg launch -k path` without agent name
- **THEN** lists all available agents

### Requirement: Build Command

The CLI MUST provide `kfg build` for kustomization output.

#### Scenario: Basic build
- **GIVEN** user wants to inspect kustomization
- **WHEN** user runs `kfg build .kfg/overlay/dev`
- **THEN** outputs YAML to stdout

#### Scenario: Build with output
- **GIVEN** user wants to save output
- **WHEN** user runs `kfg build path -o output.yaml`
- **THEN** writes YAML to file

#### Scenario: Kustomize alias
- **GIVEN** user runs `kfg kustomize path`
- **THEN** behaves identically to `kfg build`

### Requirement: Store Command

The CLI MUST provide `kfg store` for configuration management with v2 image layer system.

#### Scenario: Store parent command
- **GIVEN** user runs `kfg store`
- **THEN** shows subcommands for image and workspace operations

#### Scenario: Image subcommands
- **GIVEN** user runs `kfg store image build/push/list/inspect/remove`
- **THEN** image commands work with immutable configuration snapshots

#### Scenario: Workspace subcommands
- **GIVEN** user runs `kfg store workspace start/stop`
- **THEN** workspace commands materialize and restore images with backup safety

#### Scenario: Store directory override
- **WHEN** user runs `kfg store --store /custom/path`
- **THEN** operations use specified directory instead of default

### Requirement: Log Command

The CLI MUST provide `kfg log` for structured logging.

#### Scenario: Basic log
- **GIVEN** user wants to log event
- **WHEN** user runs `kfg log info "component" "message"`
- **THEN** writes to JSONL file

### Requirement: Flag Naming

The CLI MUST use consistent flag names.

#### Scenario: Kustomize path
- **WHEN** using `--kustomize` flag
- **THEN** accepts kustomization path
- **AND** short form `-k`

#### Scenario: File path
- **WHEN** using `--file` flag
- **THEN** accepts manifest file
- **AND** short form `-f`

#### Scenario: Workflow selection
- **WHEN** using `--workflow` flag
- **THEN** accepts workflow name
- **AND** short form `-w`

#### Scenario: Command filter
- **WHEN** using `--cmds` flag
- **THEN** accepts comma-separated list
- **AND** short form `-c`

#### Scenario: Output file
- **WHEN** using `--output` flag
- **THEN** accepts output file path
- **AND** short form `-o`

### Requirement: Launch Flags

The CLI MUST use consistent flag names for launch.

#### Scenario: Kustomize path
- **WHEN** using `--kustomize` flag on launch
- **THEN** accepts kustomization path
- **AND** short form `-k`

#### Scenario: File path
- **WHEN** using `--file` flag on launch
- **THEN** accepts manifest file
- **AND** short form `-f`

#### Scenario: Workflow selection
- **WHEN** using `--workflow` flag on launch
- **THEN** accepts workflow name
- **AND** short form `-w`

#### Scenario: Command filter
- **WHEN** using `--cmds` flag on launch
- **THEN** accepts comma-separated list
- **AND** short form `-c`

### Requirement: Exit Codes

The CLI MUST use standard exit codes.

#### Scenario: Success
- **WHEN** command succeeds
- **THEN** exit code 0

#### Scenario: Usage error
- **WHEN** incorrect arguments or missing flags
- **THEN** exit code 2

#### Scenario: Runtime error
- **WHEN** runtime failure (file not found, invalid YAML)
- **THEN** exit code 1

### Requirement: Output Contracts

Commands MUST separate stdout and stderr.

#### Scenario: Normal output
- **WHEN** command succeeds
- **THEN** primary output to stdout
- **AND** diagnostics to stderr

#### Scenario: Error output
- **WHEN** command fails
- **THEN** error to stderr
- **AND** no output to stdout

### Requirement: Deterministic Output

Shell generation MUST be deterministic.

#### Scenario: Same input same output
- **WHEN** same command run multiple times
- **THEN** output identical each time

### Requirement: Help Documentation

The CLI MUST provide helpful documentation.

#### Scenario: Root help
- **WHEN** user runs `kfg --help`
- **THEN** lists all subcommands
- **AND** documents global flags

#### Scenario: Subcommand help
- **WHEN** user runs `kfg apply --help`
- **THEN** documents apply flags and examples

### Requirement: Version Output

The CLI MUST provide version information.

#### Scenario: Version flag
- **WHEN** user runs `kfg --version`
- **THEN** outputs: `kfg version <semver> (<commit>, <date>)`

### Requirement: Verbose Flag

The CLI MUST provide verbose control.

#### Scenario: Verbose flag
- **WHEN** user runs `kfg -v 2 apply`
- **THEN** sets verbosity level

#### Scenario: Verbose levels
- **GIVEN** verbose values 0-3
- **WHEN** level set
- **THEN** controls stderr output visibility

### Requirement: Shell Completion

The CLI MUST provide shell completion for common shells.

#### Scenario: Bash completion
- **WHEN** user runs `kfg completion bash`
- **THEN** outputs bash completion script to stdout
- **AND** script can be sourced for immediate use

#### Scenario: Zsh completion
- **WHEN** user runs `kfg completion zsh`
- **THEN** outputs zsh completion script to stdout

#### Scenario: Fish completion
- **WHEN** user runs `kfg completion fish`
- **THEN** outputs fish completion script to stdout

#### Scenario: PowerShell completion
- **WHEN** user runs `kfg completion powershell`
- **THEN** outputs PowerShell completion script to stdout

#### Scenario: No descriptions option
- **WHEN** user runs `kfg completion bash --no-descriptions`
- **THEN** outputs completion script without descriptions

### Requirement: JSON Output Flag

Commands MUST support JSON output format via `--json` flag where appropriate.

#### Scenario: JSON flag is command-specific
- **WHEN** a command needs JSON output
- **THEN** `--json` flag is defined on that command (not inherited globally)

#### Scenario: Store list JSON
- **WHEN** user runs `kfg store list --json`
- **THEN** outputs JSON array of entries

#### Scenario: Store image list JSON
- **WHEN** user runs `kfg store image list --json`
- **THEN** outputs JSON array of image objects

#### Scenario: Store image inspect JSON
- **WHEN** user runs `kfg store image inspect <name> --json`
- **THEN** outputs full metadata as JSON object

### Requirement: Error Messages

Error messages MUST be clear and actionable.

#### Scenario: Workflow not found
- **WHEN** workflow doesn't exist
- **THEN** lists available workflows

#### Scenario: Command not in workflow
- **WHEN** cmd not in workflow
- **THEN** lists valid cmds