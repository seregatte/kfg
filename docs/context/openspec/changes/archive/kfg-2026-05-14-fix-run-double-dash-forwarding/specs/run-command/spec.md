## MODIFIED Requirements

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
