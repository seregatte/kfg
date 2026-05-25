## MODIFIED Requirements

### Requirement: Dev workflow phase ordering

The workflow `kfg.workflow.agents` SHALL execute before-steps in the following weight order: Gitignore (-90) → Detection (-70) → Scaffold (-65) → Settings (-63) → Context (-60) → ctx7 install (-55) → ctx7 inject (-50) → Commands (-45) → MCP (-40) → Subagent (-35) → Per-agent cleanup (90).

#### Scenario: Full workflow execution for opencode agent
- **WHEN** `AGENT` is `opencode`
- **THEN** steps execute in order: ensure-gitignore → detect-agent → materialize-scaffold(opencode) → settings(opencode) → copy-context → ctx7.install(opencode) → inject-ctx7-context → convert-commands(opencode) → aggregate-mcp(opencode) → cleanup(other agents) → after cleanup

#### Scenario: Claude agent gets CLAUDE.md context
- **WHEN** `AGENT` is `claude`
- **THEN** copy-context uses `DEST: "CLAUDE.md"` and inject-ctx7-context targets `CLAUDE.md`

#### Scenario: ctx7 install before inject
- **WHEN** ctx7 phases execute
- **THEN** `kfg.extension.ctx7.install` (weight -55) MUST complete before `kfg.inject-ctx7-context` (weight -50)

### Requirement: Per-agent cleanup in before phase

The workflow SHALL include per-agent cleanup steps at weight 90 that remove other agents' artifacts based on the detected agent.

#### Scenario: Opencode agent cleanup
- **WHEN** `AGENT` is `opencode`
- **THEN** cleanup removes `.gemini .pi .claude .mcp.json`

#### Scenario: Claude agent cleanup
- **WHEN** `AGENT` is `claude`
- **THEN** cleanup removes `.opencode .gemini .pi opencode.json`

### Requirement: After phase final cleanup

The workflow SHALL include a final cleanup step in the `after` phase that removes all agent artifacts: `.claude .gemini .opencode .pi opencode.json AGENTS.md CLAUDE.md .mcp.json`.

#### Scenario: After cleanup runs regardless of agent
- **WHEN** any agent command completes
- **THEN** after-cleanup removes all agent directories and generated files
