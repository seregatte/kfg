## Purpose

Define how the `kfg.inject-ctx7-context` step reads ctx7 documentation and injects it into agent files using upsert markers.
## Requirements
### Requirement: Inject ctx7 context into agent file

The step `kfg.inject-ctx7-context` SHALL read ctx7 documentation from `CTX7_CONTEXT` and inject it into the target file using `<!-- context7 -->` markers for upsert semantics.

#### Scenario: Inject into file without existing ctx7 section
- **WHEN** `CTX7_CONTEXT` contains context7 content and `TARGET_FILE` is `AGENTS.md`
- **THEN** the step appends `<!-- context7 -->` markers with the supplied content to `AGENTS.md`

#### Scenario: Replace existing ctx7 section
- **WHEN** `TARGET_FILE` already contains a `<!-- context7 -->` section and new `CTX7_CONTEXT` content is available
- **THEN** the step replaces the existing section between markers with the new content

#### Scenario: Missing ctx7 context env
- **WHEN** `CTX7_CONTEXT` is empty or unset
- **THEN** the step logs a warning and exits with success (no error)

### Requirement: Use AGENT env var for path resolution

The step MUST use `$AGENT` env var (not `$NIXAI_AGENT`) to construct the source file path `.$AGENT/ctx7-agents.md`, consistent with kfg's detect-agent convention.

#### Scenario: AGENT env var drives path
- **WHEN** `AGENT` is `claude` and `TARGET_FILE` is `CLAUDE.md`
- **THEN** the step reads from `.claude/ctx7-agents.md` and injects into `CLAUDE.md`

### Requirement: Configurable target file

The step SHALL use `TARGET_FILE` env var to determine the injection target, defaulting to `AGENTS.md`.

#### Scenario: Custom target file
- **WHEN** `TARGET_FILE` is `CLAUDE.md` and `CTX7_CONTEXT` contains ctx7 content
- **THEN** the step injects ctx7 context into `CLAUDE.md`

