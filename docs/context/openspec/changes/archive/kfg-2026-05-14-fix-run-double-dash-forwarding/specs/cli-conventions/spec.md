## MODIFIED Requirements

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
