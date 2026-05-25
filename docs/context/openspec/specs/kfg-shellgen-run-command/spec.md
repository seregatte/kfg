# Run Command Specification

## Purpose

Define the `kfg run` command for one-shot agent execution, generating shell code, sourcing it, and executing an agent in a single invocation.
## Requirements
### Requirement: Run command syntax

The CLI MUST provide `kfg run` for one-shot agent execution.

#### Scenario: Run with kustomization
- **WHEN** user runs `kfg run -k .kfg/overlay/dev claude`
- **THEN** loads kustomization from path
- **AND** resolves the workflow containing the `claude` cmd
- **AND** generates shell code
- **AND** executes the `claude` agent function with inherited stdin/stdout/stderr
- **AND** exits with the agent's exit code

#### Scenario: Run with file
- **WHEN** user runs `kfg run -f manifest.yaml claude`
- **THEN** loads manifest from file
- **AND** generates and executes the `claude` agent

#### Scenario: Run with stdin
- **WHEN** user runs `kfg run -f - claude`
- **THEN** reads manifest from stdin
- **AND** generates and executes the `claude` agent

#### Scenario: Run with extra args
- **WHEN** user runs `kfg run -k .kfg/overlay/dev claude -- --model gpt-4`
- **THEN** executes `claude` with `--model gpt-4` as positional arguments
- **AND** the `--` separator is consumed by `kfg run`, not passed to the agent
- **AND** no argument after `--` is dropped before agent execution

### Requirement: Agent matching by commandName

The run command MUST match user input against `Cmd.Metadata.CommandName`.

#### Scenario: Match by short name
- **GIVEN** a Cmd with `metadata.name: dev.agents.claude` and `metadata.commandName: claude`
- **WHEN** user runs `kfg run claude`
- **THEN** matches the Cmd by `commandName`
- **AND** uses `metadata.name` for resolver lookup

#### Scenario: Match across workflows
- **GIVEN** Cmds in multiple workflows
- **WHEN** user runs `kfg run claude` without `-w` flag
- **THEN** searches all CmdWorkflows for one containing the matched cmd
- **AND** uses the first matching workflow

#### Scenario: Agent not found
- **WHEN** user runs `kfg run nonexistent`
- **THEN** exit code 1
- **AND** lists all available agents with their workflow names

### Requirement: Workflow selection

The run command MUST auto-detect the workflow or accept an explicit one.

#### Scenario: Auto-detect workflow
- **WHEN** user runs `kfg run -k .kfg/overlay/dev claude` without `-w`
- **THEN** searches all CmdWorkflows for one containing `dev.agents.claude`
- **AND** uses the matching workflow

#### Scenario: Explicit workflow
- **WHEN** user runs `kfg run -k .kfg/overlay/dev -w dev claude`
- **THEN** uses the specified `dev` workflow
- **AND** does not search other workflows

#### Scenario: Agent not in specified workflow
- **WHEN** user runs `kfg run -k path -w openspec claude`
- **AND** `claude` is not in the `openspec` workflow
- **THEN** exit code 1
- **AND** error message indicates the cmd is not in the specified workflow

### Requirement: Agent discovery

Running `kfg run` without an agent name MUST list available agents.

#### Scenario: List agents
- **WHEN** user runs `kfg run -k .kfg/overlay/dev` without agent name
- **THEN** lists all Cmds exposed by CmdWorkflows
- **AND** shows each agent's `commandName` and workflow name
- **AND** does not execute any agent

#### Scenario: No agents found
- **WHEN** manifests contain no CmdWorkflow or no Cmds
- **THEN** exit code 1
- **AND** error message indicates no agents found

### Requirement: Run flags

The run command MUST support the same input flags as apply.

#### Scenario: Kustomize path
- **WHEN** user runs `kfg run -k .kfg/overlay/dev claude`
- **THEN** short flag `-k` works same as `--kustomize`

#### Scenario: Manifest file
- **WHEN** user runs `kfg run -f manifest.yaml claude`
- **THEN** short flag `-f` works same as `--file`

#### Scenario: Workflow selection
- **WHEN** user runs `kfg run -w dev claude`
- **THEN** short flag `-w` works same as `--workflow`

#### Scenario: Refresh flag
- **WHEN** user runs `kfg run -k path claude --refresh`
- **THEN** the runtime SHALL set refresh state for the generated shell execution
- **AND** cacheable Steps SHALL bypass existing cache entries during that run

### Requirement: Flag validation

The CLI MUST validate flag combinations.

#### Scenario: Required input
- **WHEN** user runs `kfg run` without `-k` or `-f`
- **THEN** exit code 2 (usage error)
- **AND** error message indicates required flag

#### Scenario: Mutual exclusion
- **WHEN** user runs `kfg run -k path -f file claude`
- **THEN** exit code 2 (usage error)
- **AND** error message indicates flag conflict

### Requirement: Temp file lifecycle

The run command MUST manage temp files for shell code execution.

#### Scenario: Temp file cleanup
- **WHEN** agent execution completes (success or failure)
- **THEN** temp file is removed via EXIT trap

#### Scenario: Shell code content
- **WHEN** `kfg run` generates the temp script
- **THEN** script contains: generated shell code, cleanup trap, agent function call with `"$@"`

### Requirement: Process execution

The run command MUST execute the agent as a child process.

#### Scenario: Stream inheritance
- **WHEN** agent is executed
- **THEN** stdin, stdout, and stderr are inherited from the parent process

#### Scenario: Exit code propagation
- **WHEN** agent exits with code N
- **THEN** `kfg run` exits with code N

#### Scenario: Cached Step reuse
- **WHEN** a generated agent command reaches a cacheable Step with a valid cache entry
- **THEN** the runtime SHALL restore that Step's artifacts and outputs instead of re-running the Step

### Requirement: Pipeline extraction

The apply pipeline MUST be extracted into a reusable function.

#### Scenario: Shared pipeline
- **GIVEN** `runApplyPipeline()` function
- **WHEN** called by `apply` or `run` with same inputs
- **THEN** returns identical `ApplyResult`
- **AND** does not produce side effects (no output, no file writes)

#### Scenario: ApplyResult structure
- **WHEN** pipeline completes successfully
- **THEN** `ApplyResult` contains: Resources, Shell, BuildResultYAML, Index, Resolver

