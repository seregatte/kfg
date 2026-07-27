## MODIFIED Requirements

### Requirement: kfg ai COMMAND

The `kfg ai` command SHALL start an AI agent pre-configured to generate kfg configurations.

- The command SHALL read the `KFG_AI_AGENT` environment variable to determine which agent to launch
- If `KFG_AI_AGENT` is not set, the command SHALL default to `pi`
- The command SHALL forward all arguments after `--` to the agent as its initial prompt
- The command SHALL use `packages/domains/ai-agents/overlays/ai` when that directory exists relative to the current working directory
- The command SHALL use `https://github.com/seregatte/kfg.git//packages/domains/ai-agents/overlays/ai?ref=main` when the local overlay directory does not exist
- The command SHALL return an error when the local overlay path cannot be inspected for a reason other than absence
- The command SHALL select its overlay independently of `KFG_KPATH`
- The command SHALL construct and execute: `kfg run -k <overlay> <agent> -- <args>`
- The command SHALL execute via subprocess using the kfg binary (os.Executable)
- The command SHALL inherit stdin, stdout, and stderr from the parent process
- The command SHALL propagate the agent's exit code to the caller

#### Scenario: Default agent with local overlay
- **WHEN** the current working directory contains `packages/domains/ai-agents/overlays/ai`
- **AND** the user runs `kfg ai -- "create a new project"`
- **THEN** the command executes `kfg run -k packages/domains/ai-agents/overlays/ai pi -- "create a new project"`

#### Scenario: Custom agent with local overlay
- **WHEN** the current working directory contains `packages/domains/ai-agents/overlays/ai`
- **AND** the user runs `KFG_AI_AGENT=opencode kfg ai`
- **THEN** the command executes `kfg run -k packages/domains/ai-agents/overlays/ai opencode`

#### Scenario: Online overlay outside the checkout
- **WHEN** the local AI overlay directory does not exist relative to the current working directory
- **AND** the user runs `kfg ai`
- **THEN** the command executes `kfg run -k https://github.com/seregatte/kfg.git//packages/domains/ai-agents/overlays/ai?ref=main pi`

#### Scenario: Generic kustomization setting is ignored
- **WHEN** `KFG_KPATH` points to a different Kustomization
- **AND** the user runs `kfg ai`
- **THEN** the command selects the local AI overlay when it exists or the online AI overlay when it does not
- **AND** the command passes that selected source explicitly through `-k`

#### Scenario: Local overlay inspection fails
- **WHEN** the local AI overlay path cannot be inspected for a reason other than absence
- **AND** the user runs `kfg ai`
- **THEN** the command returns an error without selecting the online overlay

#### Scenario: Arguments forwarded to agent
- **WHEN** user runs `kfg ai -- --model sonnet --provider anthropic "create project"`
- **THEN** `--model sonnet --provider anthropic "create project"` is forwarded after `--`

#### Scenario: No arguments starts interactive
- **WHEN** user runs `kfg ai`
- **THEN** the agent starts in interactive mode with the wizard command available
