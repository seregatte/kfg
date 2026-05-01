# JSONL Persistence Delta Specification

## Purpose

This delta specification modifies the JSONL persistence requirements to change the log file extension from `.jsonl` to `.log`.

## Requirements

### Requirement: JSONL file persistence

The system SHALL persist all log events to a single JSONL file at a standard location.

#### Scenario: Default log file location
- **WHEN** a log event is produced
- **THEN** the event is appended to `${XDG_STATE_HOME:-$HOME/.local/state}/kfg/logs/kfg.log`
- **AND** the file is created if it does not exist

#### Scenario: Custom log file location
- **WHEN** `KFG_LOG_FILE` is set to a custom path
- **THEN** log events are appended to that file instead of default location

#### Scenario: Custom log directory
- **WHEN** `KFG_LOG_DIR` is set
- **THEN** the log file is `kfg.log` within that directory