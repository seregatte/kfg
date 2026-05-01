# Log Command Specification

## Purpose

Define the `kfg log` command for structured logging with JSONL persistence and per-invocation session correlation.

## Requirements

### Requirement: Log command syntax

The CLI MUST provide a `kfg log` command for writing structured log entries.

#### Scenario: Basic invocation
- **WHEN** user runs `kfg log <level> <component> [message]`
- **THEN** a log entry is written to JSONL file
- **AND** appropriate human output appears in stderr based on KFG_VERBOSE

#### Scenario: Valid log levels
- **WHEN** user specifies a log level
- **THEN** level MUST be one of: error, warn, info, detail, debug
- **AND** invalid level results in usage error

#### Scenario: Component field
- **WHEN** user runs `kfg log info "mycomponent" "message"`
- **THEN** component field is set to `mycomponent`

#### Scenario: Message field optional
- **WHEN** user runs `kfg log info "component"` without message
- **THEN** log entry has empty message field
- **AND** command succeeds

### Requirement: Log flags

The `kfg log` command MUST support specific flags.

#### Scenario: Session ID flag
- **WHEN** user runs `kfg log --session-id "custom-123" info "component" "message"`
- **THEN** log entry includes `session_id: "custom-123"`

#### Scenario: Session ID flag overrides env var
- **GIVEN** `KFG_SESSION_ID="env-456"`
- **WHEN** user runs `kfg log --session-id "flag-789" info "component" "message"`
- **THEN** log entry includes `session_id: "flag-789"`

#### Scenario: Empty session ID flag
- **WHEN** user runs `kfg log --session-id "" info "component" "message"`
- **THEN** log entry does not include `session_id` field

#### Scenario: Source flag (hidden)
- **WHEN** internal code runs `kfg log --source "internal" info "component"`
- **THEN** source field is set for internal identification
- **AND** flag is hidden from user help

### Requirement: KFG_VERBOSE control

The system SHALL use `KFG_VERBOSE` as the sole variable controlling human output visibility.

#### Scenario: Default value
- **WHEN** `KFG_VERBOSE` is not set
- **THEN** default value is `1`
- **AND** error events produce human stderr output

#### Scenario: Verbose=0 (silent)
- **WHEN** `KFG_VERBOSE=0`
- **THEN** all log events persist to JSONL
- **AND** no human output appears in stderr

#### Scenario: Verbose=1 (error only)
- **WHEN** `KFG_VERBOSE=1`
- **THEN** error events produce human stderr output
- **AND** warn, info, detail, debug do not produce human output

#### Scenario: Verbose=2 (error + warn + info)
- **WHEN** `KFG_VERBOSE=2`
- **THEN** error, warn, info events produce human stderr output
- **AND** detail, debug do not produce human output

#### Scenario: Verbose=3 (all levels)
- **WHEN** `KFG_VERBOSE=3`
- **THEN** all levels produce human stderr output

### Requirement: --verbose flag synchronization

The `--verbose` CLI flag SHALL control verbosity by setting `KFG_VERBOSE`.

#### Scenario: Flag overrides env var
- **WHEN** user invokes `kfg --verbose=2 apply`
- **AND** `KFG_VERBOSE=0` is in environment
- **THEN** `KFG_VERBOSE` is set to `2`
- **AND** logger is reinitialized with level 2

#### Scenario: Flag not provided uses env var
- **WHEN** user invokes `kfg apply`
- **AND** `KFG_VERBOSE=3` is in environment
- **THEN** `KFG_VERBOSE` remains at `3`

#### Scenario: Neither flag nor env var
- **WHEN** user invokes `kfg apply`
- **AND** `KFG_VERBOSE` is not set
- **THEN** verbose defaults to `1`

### Requirement: JSONL persistence

All log entries MUST persist to a JSONL file.

#### Scenario: JSONL file location
- **WHEN** `KFG_JSONL_FILE` is set
- **THEN** entries are written to that file
- **AND** default location is `$TMPDIR/kfg-logs.jsonl`

#### Scenario: JSONL entry structure
- **WHEN** a log entry is written
- **THEN** entry contains: timestamp, level, component, message
- **AND** optional fields: session_id, source

#### Scenario: JSONL persists regardless of verbosity
- **WHEN** `KFG_VERBOSE=0`
- **THEN** entries still persist to JSONL
- **AND** no stderr output appears

### Requirement: Session ID correlation

Session ID MUST enable per-invocation log correlation. See `session-system` spec for complete session requirements.

#### Scenario: Shell wrapper auto-generation
- **WHEN** generated command wrapper is invoked
- **THEN** `KFG_SESSION_ID` is set with format `timestamp-random`
- **AND** all logs in invocation share same session ID

#### Scenario: Environment enrichment
- **GIVEN** `KFG_SESSION_ID="session-abc"`
- **WHEN** any log method is called
- **THEN** `session_id` field is automatically included

### Requirement: Core prefix for Go logs

Go code logs SHALL automatically prefix component with `core:`.

#### Scenario: Go log with core prefix
- **WHEN** `logger.Error("apply", "Resolution failed")` is called from Go
- **THEN** component field is `core:apply`
- **AND** human output shows `[ERROR][core:apply] Resolution failed`

#### Scenario: Shell log without core prefix
- **WHEN** `kfg log info "feature:mcps" "syncing"` is invoked from shell
- **THEN** component field is `feature:mcps`
- **AND** no `core:` prefix is added

#### Scenario: All Go log methods add prefix
- **WHEN** any Go logger method is called (Error, Warn, Info, Detail, Debug)
- **THEN** `core:` prefix is automatically added
- **AND** callers do not need to manually add prefix