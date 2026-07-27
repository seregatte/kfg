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

## MODIFIED Requirements

### Requirement: kfg ai OVERLAY SOURCE

The `kfg ai` command SHALL select its AI overlay by checking for a local directory first,
falling back to a canonical online source. The `KFG_KPATH` environment variable SHALL NOT
affect the overlay selection.

- The command SHALL use `packages/domains/ai-agents/overlays/ai` when that directory exists
  relative to the current working directory
- The command SHALL use `https://github.com/seregatte/kfg.git//packages/domains/ai-agents/overlays/ai?ref=main`
  when the local overlay directory does not exist
- The command SHALL return an error when the local overlay path cannot be inspected for a
  reason other than absence (e.g., permission error)
- The command SHALL NOT read, clear, or mutate `KFG_KPATH`

#### Scenario: Local overlay selected from checkout root
- **WHEN** the current working directory contains `packages/domains/ai-agents/overlays/ai`
- **AND** the user runs `kfg ai -- "create a new project"`
- **THEN** the command executes `kfg run -k packages/domains/ai-agents/overlays/ai pi -- "create a new project"`

#### Scenario: Remote overlay selected outside checkout
- **WHEN** `packages/domains/ai-agents/overlays/ai` does not exist relative to the current working directory
- **AND** the user runs `kfg ai`
- **THEN** the command selects `https://github.com/seregatte/kfg.git//packages/domains/ai-agents/overlays/ai?ref=main`
- **AND** passes it as the `-k` argument

#### Scenario: KFG_KPATH ignored
- **WHEN** `KFG_KPATH` is set to an unrelated kustomization path
- **AND** the user runs `kfg ai`
- **THEN** the command SHALL ignore `KFG_KPATH` and select the AI overlay independently

#### Scenario: Local overlay inspection error
- **WHEN** the local overlay path cannot be inspected for a reason other than absence
- **AND** the user runs `kfg ai`
- **THEN** the command returns an error without selecting the remote overlay
