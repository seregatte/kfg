## MODIFIED Requirements

### Requirement: Log command syntax

The CLI MUST provide a `kfg sys log` command for writing structured log entries.

#### Scenario: Basic invocation
- **WHEN** user runs `kfg sys log <level> <component> [message]`
- **THEN** a log entry is written to the configured JSONL file
- **AND** appropriate human output appears in stderr based on `KFG_VERBOSE`

#### Scenario: Component field
- **WHEN** user runs `kfg sys log info "mycomponent" "message"`
- **THEN** component field is set to `mycomponent`

#### Scenario: Message field optional
- **WHEN** user runs `kfg sys log info "component"` without message
- **THEN** the log entry has an empty message field
- **AND** the command succeeds

### Requirement: Log flags

The `kfg sys log` command MUST support specific flags.

#### Scenario: Session ID flag
- **WHEN** user runs `kfg sys log --session-id "custom-123" info "component" "message"`
- **THEN** the log entry includes `session_id: "custom-123"`

#### Scenario: Session ID flag overrides env var
- **GIVEN** `KFG_SESSION_ID="env-456"`
- **WHEN** user runs `kfg sys log --session-id "flag-789" info "component" "message"`
- **THEN** the log entry includes `session_id: "flag-789"`

#### Scenario: Empty session ID flag
- **WHEN** user runs `kfg sys log --session-id "" info "component" "message"`
- **THEN** the log entry does not include `session_id` field

#### Scenario: Source flag
- **WHEN** internal code runs `kfg sys log --source "shell" info "component" "message"`
- **THEN** the log entry records the provided source value

### Requirement: JSONL persistence

All log entries MUST persist to a JSONL file.

#### Scenario: JSONL file location
- **WHEN** `KFG_LOG_FILE` is set
- **THEN** entries are written to that file
- **AND** otherwise the logger SHALL use the configured state-directory default path

#### Scenario: JSONL entry structure
- **WHEN** a log entry is written
- **THEN** the entry contains `ts`, `level`, `component`, and `msg`
- **AND** optional fields MAY include `session_id`, `source`, `workflow_name`, `kustomization_name`, and `step_name`

#### Scenario: Step context enrichment
- **GIVEN** `KFG_STEP_NAME="ctx7.install"`
- **WHEN** `kfg sys log info "step" "Installed"` is invoked
- **THEN** the log entry includes `step_name: "ctx7.install"`

#### Scenario: Legacy step component normalization
- **WHEN** `kfg sys log info "step:ctx7.install" "Installed"` is invoked
- **THEN** the log entry records `component: "step"`
- **AND** the log entry records `step_name: "ctx7.install"`

### Requirement: Core prefix for Go logs

Go code logs SHALL automatically prefix component with `core:`.

#### Scenario: Go log with core prefix
- **WHEN** `logger.Error("apply", "Resolution failed")` is called from Go
- **THEN** component field is `core:apply`
- **AND** human output shows `[ERROR][core:apply] Resolution failed`

#### Scenario: Shell log without core prefix
- **WHEN** `kfg sys log info "feature:mcps" "Syncing"` is invoked from shell
- **THEN** component field is `feature:mcps`
- **AND** no `core:` prefix is added
