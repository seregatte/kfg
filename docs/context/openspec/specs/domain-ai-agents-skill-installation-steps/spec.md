# Domain: AI Agents - Skill Installation Steps Specification

## Purpose

This specification defines generic skill installation Steps for extensions that install agent skills via external CLIs.
## Requirements

### Requirement: Generic skill installation Step

The system SHALL provide a Step resource for each extension that installs agent skills via external CLIs. Each Step MUST receive all agent-specific configuration via `spec.env` variables. Steps MUST NOT contain `case`, `if`, or other conditional logic for agent selection. Each Step SHALL validate its required inputs, create any required output directories, copy generated files into `OUTPUT_DIR`, and preserve the working install behavior defined by the normalized extension contract.

#### Scenario: ctx7 skill installation for claude agent
- **WHEN** the workflow invokes `kfg.extension.ctx7.install` with `FLAGS="--claude --yes"` and `OUTPUT_DIR=".claude/skills/"`
- **THEN** the Step executes `ctx7 setup --cli --project --claude --yes`
- **AND** it copies generated skills into `.claude/skills/`

#### Scenario: chrome-devtools skill installation
- **WHEN** the workflow invokes `kfg.extension.chrome-devtools.install` with `SKILL_NAME="ChromeDevTools/chrome-devtools-mcp"`, `AGENT_FLAG="claude-code"`, `AGENT_HOME=".claude"`, and `OUTPUT_DIR=".claude/skills/"`
- **THEN** the Step executes `npx skills add ChromeDevTools/chrome-devtools-mcp --agent claude-code --copy --yes`
- **AND** it copies generated skills from `.claude/skills/` into `OUTPUT_DIR`

#### Scenario: playwright skill installation for opencode-style agent home
- **WHEN** the workflow invokes `kfg.extension.playwright.install` with `AGENT_HOME=".agents"` and `OUTPUT_DIR=".opencode/skills/"`
- **THEN** the Step executes the Playwright skill installation flow defined by the normalized extension contract
- **AND** it copies generated skills from the temporary agent home into `OUTPUT_DIR`

#### Scenario: gws skill installation
- **WHEN** the workflow invokes `kfg.extension.gws.install` with `SKILL_NAME="googleworkspace/cli"`, `AGENT_FLAG="claude-code"`, `AGENT_HOME=".claude"`, and `OUTPUT_DIR=".claude/skills/"`
- **THEN** the Step executes `npx skills add googleworkspace/cli --agent claude-code --copy --yes`
- **AND** it copies generated skills into `OUTPUT_DIR`

#### Scenario: notebooklm skill installation
- **WHEN** the workflow invokes `kfg.extension.notebooklm.install` with `INSTALL_CMD="notebooklm skill install"`, `AGENT_HOME=".claude"`, and `OUTPUT_DIR=".claude/skills/"`
- **THEN** the Step executes `notebooklm skill install`
- **AND** it copies generated NotebookLM skill files into `OUTPUT_DIR`

#### Scenario: openspec skill installation for opencode
- **WHEN** the workflow invokes `kfg.extension.openspec.install` with `TOOLS_FLAG="--tools=opencode"`, `AGENT_HOME=".opencode"`, and `OUTPUT_DIR=".opencode/skills/"`
- **THEN** the Step executes `openspec init --tools=opencode --force`
- **AND** it copies generated skills and commands into overlay-owned output paths

#### Scenario: Install Step logging attribution
- **WHEN** an install Step emits a runtime log event through `__kfg_log_*`
- **THEN** the event SHALL rely on runtime-provided `step_name` attribution
- **AND** the Step SHALL NOT need to encode its Step identity inside the component string

### Requirement: Step output

Each install Step SHALL produce an output named `installed` of type `string` with a deterministic non-empty success value.

#### Scenario: Step produces output on success
- **WHEN** the install Step completes successfully
- **THEN** the Step outputs a non-empty deterministic string value for `installed`

### Requirement: Missing env var handling

Each install Step SHALL validate required env vars and log an error if they are missing.

#### Scenario: Missing OUTPUT_DIR
- **WHEN** the workflow invokes an install Step with `OUTPUT_DIR=""`
- **THEN** the Step logs an error and exits with non-zero status

#### Scenario: Missing agent-specific input
- **WHEN** the workflow invokes an install Step without a required agent-specific env var such as `AGENT_HOME`, `AGENT_FLAG`, `FLAGS`, or `TOOLS_FLAG`
- **THEN** the Step logs an error and exits with non-zero status
