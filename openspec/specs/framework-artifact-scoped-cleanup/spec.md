# Artifact-Scoped Cleanup Specification

## Purpose

Specifies that workspace cleanup operations SHALL be scoped to only paths that were materialized from the image. This ensures stop removes exactly what the image added, preserving unrelated workspace files.

## Requirements

### Requirement: Materialized paths tracking
The system SHALL track which paths were materialized during start in instance metadata.

#### Scenario: Record materialized paths
- **WHEN** start completes successfully
- **THEN** system records `MaterializedPaths` in instance metadata
- **AND** paths are relative to workspace root

#### Scenario: Track file paths
- **WHEN** image contains `CLAUDE.md` (file)
- **THEN** instance metadata includes `"CLAUDE.md"` in `MaterializedPaths`

#### Scenario: Track directory paths
- **WHEN** image contains `.pi/config.json` (file in directory)
- **THEN** instance metadata includes `".pi/config.json"` in `MaterializedPaths`
- **AND** system tracks file path, not just parent directory

#### Scenario: Track nested paths
- **WHEN** image contains `.pi/subdir/file.txt`
- **THEN** instance metadata includes `".pi/subdir/file.txt"` in `MaterializedPaths`
- **AND** system tracks full relative path

### Requirement: Scoped cleanup execution
The system SHALL delete only materialized paths before restoring backup.

#### Scenario: Delete materialized file
- **WHEN** executing stop for instance with `MaterializedPaths: ["CLAUDE.md"]`
- **THEN** system removes `CLAUDE.md` from workspace
- **AND** system logs "[STOP] Removing materialized file: CLAUDE.md"

#### Scenario: Delete materialized directory content
- **WHEN** executing stop for instance with `MaterializedPaths: [".pi/config.json"]`
- **THEN** system removes `.pi/config.json` from workspace
- **AND** system may remove `.pi/` directory if it becomes empty

#### Scenario: Preserve non-materialized files
- **WHEN** workspace contains `README.md` that was not in image
- **AND** stop executes for instance
- **THEN** system does NOT remove `README.md`
- **AND** `README.md` remains in workspace after stop

#### Scenario: Cleanup order for nested paths
- **WHEN** `MaterializedPaths` contains `.pi/a.txt` and `.pi/subdir/b.txt`
- **THEN** system removes nested paths first (deepest first)
- **OR** system uses `os.RemoveAll()` which handles order automatically

### Requirement: Idempotent cleanup
The system SHALL succeed when materialized paths no longer exist in workspace.

#### Scenario: Missing materialized path
- **WHEN** stop executes and `CLAUDE.md` was deleted by user before stop
- **THEN** system succeeds (no error)
- **AND** system logs "[STOP] CLAUDE.md already removed"
- **AND** system continues with backup restoration

#### Scenario: All paths missing
- **WHEN** stop executes and all `MaterializedPaths` are missing from workspace
- **THEN** system succeeds
- **AND** system proceeds to restore backup (if exists)

### Requirement: Backup restoration after cleanup
The system SHALL restore backed up files after deleting materialized paths.

#### Scenario: Restore after cleanup
- **WHEN** stop deletes materialized paths
- **THEN** system restores backed up conflicting files
- **AND** restored files match pre-materialization state

#### Scenario: No backup to restore
- **WHEN** stop executes with no backup (start had no conflicts)
- **THEN** system skips restoration step
- **AND** system logs "[STOP] No backup found - cleanup only"

### Requirement: Cleanup logging
The system SHALL log which paths are removed for user visibility.

#### Scenario: Log removed paths
- **WHEN** cleanup removes materialized paths
- **THEN** system logs each removed path
- **AND** system logs total count: "[STOP] Removed 2 materialized artifacts"

#### Scenario: Log cleanup completion
- **WHEN** cleanup completes
- **THEN** system logs "[STOP] Workspace restored - image artifacts removed"

### Requirement: Empty directory cleanup
The system SHALL remove empty parent directories created by materialization.

#### Scenario: Remove empty parent
- **WHEN** materialization created `.pi/` directory with only `config.json`
- **AND** stop removes `.pi/config.json`
- **THEN** system removes `.pi/` directory if it becomes empty
- **AND** directory is not left as empty shell

#### Scenario: Preserve non-empty parent
- **WHEN** user added `.pi/notes.txt` after materialization
- **AND** stop removes `.pi/config.json`
- **THEN** system does NOT remove `.pi/` directory
- **AND** `.pi/notes.txt` remains