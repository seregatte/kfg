## MODIFIED Requirements

### Requirement: Dev workflow phase ordering

The workflow `kfg.workflow.agents` SHALL execute before-steps in the following weight order: Gitignore (-90) → Detection (-70) → Scaffold (-65) → Settings (-63) → Context (-60) → Extension install (-55 to -53) → Context injection (-50) → Commands (-45) → MCP aggregation (-40) → Subagent (-35) → Per-agent cleanup (90).

#### Scenario: Full workflow execution for opencode agent
- **WHEN** `AGENT` is `opencode`
- **THEN** steps execute in order: ensure-gitignore → detect-agent → materialize-scaffold(opencode) → materialize-settings(opencode) → copy-context → extension installs(opencode) → inject-ctx7-context → materialize-commands(opencode) → materialize-mcp(opencode) → cleanup(other agents) → after cleanup

#### Scenario: Claude agent gets CLAUDE.md context
- **WHEN** `AGENT` is `claude`
- **THEN** copy-context uses `DEST: "CLAUDE.md"` and inject-ctx7-context targets `CLAUDE.md`

#### Scenario: Extension installs complete before MCP aggregation
- **WHEN** the workflow installs `ctx7`, `openspec`, `playwright`, or `chrome-devtools`
- **THEN** those install steps MUST complete before materialize steps at weight `-40` that produce MCP configuration

### Requirement: Shared materialize step in workflow phases

The workflow SHALL use `kfg.materialize` as the shared primitive for settings, command, subagent, and MCP materialization phases.

#### Scenario: Settings use per-item materialize mode
- **WHEN** the workflow prepares settings for an agent
- **THEN** it SHALL invoke `kfg.materialize` with `MODE="per-item"`
- **AND** it SHALL provide one asset path and one output path through `ASSETS` and `OUTPUTS`

#### Scenario: Commands and subagents group by agent and type
- **WHEN** the workflow materializes command or subagent assets
- **THEN** it SHALL group them into `kfg.materialize` steps by agent and type
- **AND** each grouped step SHALL use one converter for the entire batch

#### Scenario: MCP uses aggregate materialize mode
- **WHEN** the workflow prepares MCP configuration for an agent
- **THEN** it SHALL invoke `kfg.materialize` with `MODE="aggregate"`
- **AND** it SHALL use `OUTPUTS` with exactly one destination path
- **AND** it MAY set `WRAP_KEY` when the agent expects a wrapped config object

### Requirement: Canonical extension assets in MCP aggregation

The workflow SHALL aggregate normalized extension-owned MCP assets instead of transitional `kfg.extension.self.*` MCP assets when canonical extension assets exist.

#### Scenario: Claude aggregates extension-owned MCP assets
- **WHEN** the workflow prepares MCP configuration for `claude`
- **THEN** it aggregates extension-owned MCP assets such as `kfg.extension.ctx7.mcp`, `kfg.extension.chrome-devtools.mcp`, and `kfg.extension.playwright.mcp` through `kfg.materialize`

#### Scenario: OpenCode wrapper key remains explicit
- **WHEN** the workflow prepares MCP configuration for `opencode`
- **THEN** it aggregates extension-owned MCP assets through `kfg.materialize`
- **AND** it sets `WRAP_KEY="mcp"`
