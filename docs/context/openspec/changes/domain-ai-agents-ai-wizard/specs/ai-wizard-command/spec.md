## ADDED Requirements

### Requirement: kfg ai COMMAND

The `kfg ai` command SHALL start an AI agent pre-configured to generate kfg configurations.

- The command SHALL read the `KFG_AI_AGENT` environment variable to determine which agent to launch
- If `KFG_AI_AGENT` is not set, the command SHALL default to `pi`
- The command SHALL forward all arguments after `--` to the agent as its initial prompt
- The command SHALL resolve the agent overlay to `packages/domains/ai-agents/overlays/ai/`
- The command SHALL construct and execute: `kfg run -k <overlay> <agent> -- <args>`
- The command SHALL execute via subprocess using the kfg binary (os.Executable)
- The command SHALL inherit stdin, stdout, and stderr from the parent process
- The command SHALL propagate the agent's exit code to the caller

#### Scenario: Default agent (pi)
- **WHEN** user runs `kfg ai -- "create a new project"`
- **THEN** the command executes `kfg run -k packages/domains/ai-agents/overlays/ai pi -- "create a new project"`

#### Scenario: Custom agent (opencode)
- **WHEN** user runs `KFG_AI_AGENT=opencode kfg ai`
- **THEN** the command executes `kfg run -k packages/domains/ai-agents/overlays/ai opencode`

#### Scenario: Arguments forwarded to agent
- **WHEN** user runs `kfg ai -- --model sonnet --provider anthropic "create project"`
- **THEN** `--model sonnet --provider anthropic "create project"` is forwarded after `--`

#### Scenario: No arguments starts interactive
- **WHEN** user runs `kfg ai`
- **THEN** the agent starts in interactive mode with the wizard command available

### Requirement: KFG_AI_AGENT ENVIRONMENT VARIABLE

The `KFG_AI_AGENT` environment variable SHALL control which agent the wizard launches.

- Valid values SHALL be `pi` and `opencode`
- Setting the variable to an unsupported value SHALL fall back to the default (`pi`)
- The variable SHALL be documented in `kfg ai --help`

#### Scenario: Unsupported agent value
- **WHEN** `KFG_AI_AGENT=invalid-agent`
- **AND** user runs `kfg ai`
- **THEN** the command SHALL fall back to `pi`

### Requirement: HELP TEXT

The `kfg ai --help` output SHALL document:
- Purpose: interactive configuration generation via AI agent
- The `KFG_AI_AGENT` environment variable and its default (`pi`)
- Supported agent values: `pi`, `opencode`
- Usage examples: `kfg ai`, `kfg ai -- "create a project"`
- That arguments after `--` are forwarded to the agent

#### Scenario: Help displays correctly
- **WHEN** user runs `kfg ai --help`
- **THEN** output SHALL include the purpose, env var documentation, and usage examples
