# Session System Specification

## Purpose

Define unified session concepts for kfg, including log correlation sessions for per-invocation tracking.

## Requirements

### Requirement: Session ID Format

All session identifiers MUST follow a consistent format.

#### Scenario: Session ID format
- **WHEN** a session ID is generated
- **THEN** format is `<timestamp>-<random>` (e.g., `1712938291-4821`)
- **AND** timestamp is Unix epoch seconds
- **AND** random is a random integer (0-32767)

#### Scenario: Session ID uniqueness
- **WHEN** multiple session IDs are generated
- **THEN** each ID is unique within a reasonable time window
- **AND** collision probability is acceptably low

### Requirement: Log Correlation Sessions

Log entries MUST support session ID for per-invocation correlation.

#### Scenario: Session ID in JSONL output
- **WHEN** `KFG_SESSION_ID` is set
- **THEN** log entries include `session_id` field
- **AND** field value matches the environment variable value

#### Scenario: Session ID absent when not set
- **WHEN** `KFG_SESSION_ID` is not set
- **THEN** log entries do not include `session_id` field
- **AND** logs continue to function normally

#### Scenario: Shell wrapper auto-generation
- **WHEN** a generated command wrapper is invoked
- **THEN** `KFG_SESSION_ID` is set with format `timestamp-random`
- **AND** all log calls within that invocation share the same session ID

#### Scenario: Per-invocation uniqueness
- **WHEN** the same command wrapper is invoked twice
- **THEN** each invocation gets a distinct session ID

### Requirement: Session ID CLI Flag

The `kfg log` command MUST accept a `--session-id` flag.

#### Scenario: Flag provides session ID
- **WHEN** `kfg log --session-id "custom-123" info "component" "message"` is invoked
- **THEN** log entry includes `session_id: "custom-123"`

#### Scenario: Flag overrides environment variable
- **GIVEN** `KFG_SESSION_ID="env-456"`
- **WHEN** `kfg log --session-id "flag-789" info "component" "message"` is invoked
- **THEN** log entry includes `session_id: "flag-789"`

#### Scenario: Empty flag value
- **WHEN** `kfg log --session-id "" info "component" "message"` is invoked
- **THEN** log entry does not include `session_id` field

### Requirement: Environment Variable Enrichment

The logger MUST enrich log entries with `KFG_SESSION_ID` from environment.

#### Scenario: Environment variable enrichment
- **GIVEN** `KFG_SESSION_ID="session-abc"` is set
- **WHEN** any log method is called
- **THEN** `session_id` field is automatically included
- **AND** no explicit parameter is needed

### Requirement: Session Isolation

Sessions MUST not interfere with each other.

#### Scenario: Multiple log correlation sessions
- **WHEN** multiple command invocations run concurrently
- **THEN** each maintains its own `KFG_SESSION_ID`
- **AND** log entries are correctly attributed to their session

### Requirement: Session Lifecycle

Sessions MUST have defined creation and cleanup semantics.

#### Scenario: Log session creation
- **WHEN** command wrapper is invoked
- **THEN** `KFG_SESSION_ID` is set at wrapper start
- **AND** session exists for duration of invocation

#### Scenario: Log session cleanup
- **WHEN** command wrapper completes
- **THEN** `KFG_SESSION_ID` is not explicitly cleared
- **AND** next invocation generates new session ID

