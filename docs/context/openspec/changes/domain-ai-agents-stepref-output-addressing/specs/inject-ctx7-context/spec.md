## MODIFIED Requirements

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

### Requirement: Configurable target file

The step SHALL use `TARGET_FILE` env var to determine the injection target, defaulting to `AGENTS.md`.

#### Scenario: Custom target file
- **WHEN** `TARGET_FILE` is `CLAUDE.md` and `CTX7_CONTEXT` contains ctx7 content
- **THEN** the step injects ctx7 context into `CLAUDE.md`
