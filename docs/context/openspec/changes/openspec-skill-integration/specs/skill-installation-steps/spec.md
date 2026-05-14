## MODIFIED Requirements

### Requirement: Generic skill installation Step

The system SHALL provide a Step resource for each extension that installs agent skills via external CLIs. Each Step MUST receive all agent-specific configuration via `spec.env` variables. Steps MUST NOT contain `case`, `if`, or other conditional logic for agent selection.

#### Scenario: ctx7 skill installation for claude agent
- **WHEN** the workflow invokes `kfg.extension.ctx7.install` with `INSTALL_FLAGS="--claude"` and `AGENT_HOME=".claude"`
- **THEN** the Step executes `ctx7 setup --cli --project --claude --yes` and copies skills from `.claude/skills/` to the output directory

#### Scenario: ctx7 skill installation with API key
- **WHEN** the workflow invokes `kfg.extension.ctx7.install` with `API_KEY_FLAG="--api-key sk-123"`
- **THEN** the Step includes the API key flag in the command execution

#### Scenario: ctx7 skill installation without API key
- **WHEN** the workflow invokes `kfg.extension.ctx7.install` with `API_KEY_FLAG=""`
- **THEN** the Step executes without the API key flag

#### Scenario: chrome-devtools skill installation
- **WHEN** the workflow invokes `kfg.extension.chrome-devtools.install` with `AGENT_FLAG="claude-code"` and `AGENT_HOME=".claude"`
- **THEN** the Step executes `npx skills add ChromeDevTools/chrome-devtools-mcp --agent claude-code --copy --yes` and copies skills to the output directory

#### Scenario: playwright skill installation
- **WHEN** the workflow invokes `kfg.extension.playwright.install` with `AGENT_FLAG="opencode"` and `AGENT_HOME=".agents"`
- **THEN** the Step executes `npx skills add @playwright/mcp --agent opencode --copy --yes` and copies skills to the output directory

#### Scenario: gws skill installation
- **WHEN** the workflow invokes `kfg.extension.gws.install` with `AGENT_FLAG="claude-code"` and `AGENT_HOME=".claude"`
- **THEN** the Step executes `npx skills add googleworkspace/cli --agent claude-code --copy --yes` and copies skills to the output directory

#### Scenario: notebooklm skill installation
- **WHEN** the workflow invokes `kfg.extension.notebooklm.install` with `AGENT_HOME=".claude"`
- **THEN** the Step executes `notebooklm skill install` and copies skills from `.claude/skills/notebooklm/` to the output directory

#### Scenario: openspec skill installation for opencode
- **WHEN** the workflow invokes `kfg.extension.openspec.install` with `TOOLS_FLAG="--tools=opencode"` and `AGENT_HOME=".opencode"` and `OUTPUT_DIR=".opencode/skills/"`
- **THEN** the Step executes `openspec init --tools=opencode --force` and copies skills from `.opencode/skills/` and commands from `.opencode/commands/` to the output directory

#### Scenario: openspec skill installation for claude
- **WHEN** the workflow invokes `kfg.extension.openspec.install` with `TOOLS_FLAG="--tools=claude"` and `AGENT_HOME=".claude"` and `OUTPUT_DIR=".claude/skills/"`
- **THEN** the Step executes `openspec init --tools=claude --force` and copies skills from `.claude/skills/` and commands from `.claude/commands/` to the output directory

#### Scenario: openspec skill installation for gemini
- **WHEN** the workflow invokes `kfg.extension.openspec.install` with `TOOLS_FLAG="--tools=gemini"` and `AGENT_HOME=".gemini"` and `OUTPUT_DIR=".gemini/skills/"`
- **THEN** the Step executes `openspec init --tools=gemini --force` and copies skills from `.gemini/skills/` and commands from `.gemini/commands/` to the output directory

#### Scenario: openspec skill installation for pi
- **WHEN** the workflow invokes `kfg.extension.openspec.install` with `TOOLS_FLAG="--tools=pi"` and `AGENT_HOME=".pi"` and `OUTPUT_DIR=".pi/skills/"`
- **THEN** the Step executes `openspec init --tools=pi --force` and copies skills from `.pi/skills/` and commands from `.pi/commands/` to the output directory

### Requirement: Step output

Each install Step SHALL produce an output named `installed` of type `string` to signal completion.

#### Scenario: Step produces output on success
- **WHEN** the install Step completes successfully
- **THEN** the Step outputs a non-empty string value for `installed`

### Requirement: Missing env var handling

Each install Step SHALL validate required env vars and log an error if they are missing.

#### Scenario: Missing INSTALL_CMD
- **WHEN** the workflow invokes a Step with `INSTALL_CMD=""`
- **THEN** the Step logs an error and exits with non-zero status
