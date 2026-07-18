## MODIFIED Requirements

### Requirement: Dev workflow phase ordering

The workflow `kfg.workflow.agents` SHALL execute before-steps in the following weight order: Gitignore (-90) → Detection (-70) → Scaffold (-65) → Settings (-63) → Context (-60) → Extension install (-55 to -53) → Context injection (-50) → Commands (-45) → MCP aggregation (-40) → Subagent (-35) → Per-agent cleanup (90).

#### Scenario: Full workflow execution for opencode agent
- **WHEN** `AGENT` is `opencode`
- **THEN** steps execute in order: ensure-gitignore → detect-agent → materialize-scaffold(opencode) → settings(opencode) → copy-context → extension installs(opencode) → inject-ctx7-context → convert-commands(opencode) → aggregate-mcp(opencode) → cleanup(other agents) → after cleanup

#### Scenario: Claude agent gets CLAUDE.md context
- **WHEN** `AGENT` is `claude`
- **THEN** copy-context uses `DEST: "CLAUDE.md"` and inject-ctx7-context targets `CLAUDE.md`

#### Scenario: Extension installs complete before MCP aggregation
- **WHEN** the workflow installs `ctx7`, `openspec`, `playwright`, or `chrome-devtools`
- **THEN** those install steps MUST complete before MCP aggregation steps at weight `-40`

### Requirement: Canonical extension assets in MCP aggregation

The workflow SHALL aggregate normalized extension-owned MCP assets instead of transitional `kfg.extension.self.*` MCP assets when canonical extension assets exist.

#### Scenario: Claude aggregates extension-owned MCP assets
- **WHEN** the workflow prepares MCP configuration for `claude`
- **THEN** it aggregates extension-owned MCP assets such as `kfg.extension.ctx7.mcp`, `kfg.extension.chrome-devtools.mcp`, and `kfg.extension.playwright.mcp`

#### Scenario: OpenCode wrapper key remains explicit
- **WHEN** the workflow prepares MCP configuration for `opencode`
- **THEN** it aggregates extension-owned MCP assets with `WRAPPER_KEY="mcp"`

### Requirement: Per-agent cleanup in before phase

The workflow SHALL include per-agent cleanup steps at weight `90` that remove other agents' artifacts based on the detected agent.

#### Scenario: Opencode agent cleanup
- **WHEN** `AGENT` is `opencode`
- **THEN** cleanup removes other agents' generated directories and files such as `.gemini`, `.pi`, `.claude`, and Claude-specific MCP outputs

#### Scenario: Claude agent cleanup
- **WHEN** `AGENT` is `claude`
- **THEN** cleanup removes other agents' generated directories and files such as `.opencode`, `.gemini`, `.pi`, and OpenCode-specific config outputs

### Requirement: After phase final cleanup

The workflow SHALL include a final cleanup step in the `after` phase that removes all generated agent directories and files created by the workflow.

#### Scenario: After cleanup runs regardless of agent
- **WHEN** any agent command completes
- **THEN** after-cleanup removes all generated agent directories and workflow-generated files
