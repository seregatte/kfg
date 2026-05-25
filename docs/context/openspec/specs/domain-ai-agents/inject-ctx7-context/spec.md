## ADDED Requirements

### Requirement: Inject ctx7 context into agent file

The step `kfg.inject-ctx7-context` SHALL read ctx7 documentation from `.$AGENT/ctx7-agents.md` and inject it into the target file using `<!-- context7 -->` markers for upsert semantics.

#### Scenario: Inject into file without existing ctx7 section
- **WHEN** `AGENT` is `opencode`, `TARGET_FILE` is `AGENTS.md`, and `.opencode/ctx7-agents.md` exists with context7 content
- **THEN** the step appends `<!-- context7 -->` markers with the extracted content to `AGENTS.md`

#### Scenario: Replace existing ctx7 section
- **WHEN** `TARGET_FILE` already contains a `<!-- context7 -->` section and new ctx7 content is available
- **THEN** the step replaces the existing section between markers with the new content

#### Scenario: ctx7 file not found
- **WHEN** `.$AGENT/ctx7-agents.md` does not exist
- **THEN** the step logs a warning and exits with success (no error)

#### Scenario: No context7 section in source file
- **WHEN** `.$AGENT/ctx7-agents.md` exists but contains no `<!-- context7 -->` markers
- **THEN** the step logs a warning and exits with success

#### Scenario: Structured Step logging attribution
- **WHEN** the step emits a runtime log event through `__kfg_log_*`
- **THEN** the event SHALL rely on runtime-provided `step_name` attribution
- **AND** the step SHALL NOT need to encode its Step identity inside the component string

### Requirement: Use AGENT env var for path resolution

The step MUST use `$AGENT` env var (not `$NIXAI_AGENT`) to construct the source file path `.$AGENT/ctx7-agents.md`, consistent with kfg's detect-agent convention.

#### Scenario: AGENT env var drives path
- **WHEN** `AGENT` is `claude` and `TARGET_FILE` is `CLAUDE.md`
- **THEN** the step reads from `.claude/ctx7-agents.md` and injects into `CLAUDE.md`

### Requirement: Configurable target file

The step SHALL use `TARGET_FILE` env var to determine the injection target, defaulting to `AGENTS.md`.

#### Scenario: Custom target file
- **WHEN** `TARGET_FILE` is `CLAUDE.md` and `AGENT` is `claude`
- **THEN** the step injects ctx7 context into `CLAUDE.md`
