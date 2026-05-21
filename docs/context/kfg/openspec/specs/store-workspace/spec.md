# Store Workspace Specification (Delta)

## Purpose

Specifies workspace materialization and restoration, including image extraction with artifact-scoped backups, instance-scoped tracking with materialized paths, and safe restoration from backup archives.

## Requirements

### Requirement: Automatic backup on start
The system SHALL create instance-scoped backup of only conflicting files before materializing to prevent data loss.

#### Scenario: Backup existing conflicting files
- **WHEN** starting image in non-empty directory with file `CLAUDE.md`
- **AND** image contains artifact `CLAUDE.md`
- **THEN** system backs up `CLAUDE.md` before materialization
- **AND** system does NOT backup unrelated files like `README.md`

#### Scenario: Backup location
- **WHEN** backup is created
- **THEN** system stores at `$KFG_STORE_DIR/.workspace/<name>/backup/data/` (directory backup)
- **AND** backup preserves directory structure of backed up paths

#### Scenario: Empty directory start
- **WHEN** starting image in empty directory
- **THEN** system skips backup and materializes directly
- **AND** system logs no conflicts found

#### Scenario: Partial conflicts backup
- **WHEN** image artifacts are `.pi/config.json` and `CLAUDE.md`
- **AND** workspace contains `CLAUDE.md` and `README.md` but not `.pi/`
- **THEN** system backs up only `CLAUDE.md`
- **AND** system logs backed up file count

#### Scenario: Backup overwrite on repeated start
- **WHEN** running `start --name <name>` twice to same root
- **THEN** second start backs up current state (overwriting prior backup) and materializes image

### Requirement: Workspace restoration
The system SHALL delete materialized artifacts and restore backup on demand.

#### Scenario: Stop with existing backup
- **WHEN** executing `kfg store image stop --name myproject`
- **THEN** system removes materialized paths tracked in instance metadata
- **AND** system restores backed up conflicting files

#### Scenario: Stop removes only image artifacts
- **WHEN** image materialized `.pi/config.json` and `CLAUDE.md`
- **AND** workspace also contains `README.md` (not in image)
- **THEN** stop removes `.pi/config.json` and `CLAUDE.md`
- **AND** `README.md` remains untouched in workspace

#### Scenario: Stop with no backup
- **WHEN** executing `stop --name <name>` and backup is missing
- **THEN** system succeeds (idempotent; logs message, does not error)
- **AND** system still removes materialized paths

#### Scenario: Stop cleans backup
- **WHEN** stop completes successfully
- **THEN** system deletes backup directory (restore consumes backup)
- **AND** system deletes instance record

#### Scenario: Multiple instance stop
- **WHEN** stopping instance with `--name`
- **THEN** system only removes artifacts and restores backup for that named instance

### Requirement: Instance tracking
The system SHALL maintain metadata about active instances including materialized paths.

#### Scenario: Instance record creation
- **WHEN** `start` completes
- **THEN** system records instance at `$KFG_STORE_DIR/.workspace/<name>/instance.json`

#### Scenario: Instance metadata
- **WHEN** instance is active
- **THEN** system records: instance name, image reference, materialization timestamp, workspace root, image digest, materialized paths

#### Scenario: Materialized paths tracking
- **WHEN** start materializes image with artifacts `CLAUDE.md` and `.pi/config.json`
- **THEN** instance metadata includes `MaterializedPaths: ["CLAUDE.md", ".pi/config.json"]`

#### Scenario: Instance cleanup
- **WHEN** `stop` completes
- **THEN** system removes instance record

### Requirement: Materialization validation
The system SHALL verify image can be materialized before modifying workspace.

#### Scenario: Image existence check
- **WHEN** starting image reference
- **THEN** system verifies image exists in store before backup/materialization

#### Scenario: Non-existent image error
- **WHEN** starting non-existent image
- **THEN** system fails with helpful error before modifying workspace

#### Scenario: Image metadata validation
- **WHEN** materializing image
- **THEN** system validates metadata integrity; reports corruption if found
- **AND** system loads `Files` manifest to determine artifact paths

### Requirement: File extraction
The system SHALL extract image files to workspace with proper permissions.

#### Scenario: File extraction from image
- **WHEN** materializing image
- **THEN** system extracts files from image artifacts directory to workspace root
- **AND** system preserves file structure and permissions

#### Scenario: Permission preservation
- **WHEN** extracting files
- **THEN** system preserves execute bits and other file permissions

#### Scenario: Overwrite existing files
- **WHEN** extracted file exists in workspace (after backup)
- **THEN** system overwrites with image version

#### Scenario: Track materialized paths
- **WHEN** extracting file to workspace
- **THEN** system records relative path in materialized paths list

### Requirement: Workspace isolation
The system SHALL ensure instances don't interfere with each other.

#### Scenario: Concurrent instances different roots
- **WHEN** running multiple `start` commands with different `--root` values
- **THEN** each materializes to its own root independently
- **AND** each tracks its own materialized paths

#### Scenario: Same instance name conflict
- **WHEN** attempting to start instance with name already in use
- **THEN** system fails with error (instance names must be unique across active sessions)

#### Scenario: Duplicate name different roots
- **WHEN** attempting same instance name in different projects
- **THEN** system fails (names must be globally unique in store metadata)

### Requirement: Error recovery
The system SHALL provide clear error messages for common failure scenarios.

#### Scenario: Insufficient disk space
- **WHEN** backup or materialization fails due to disk space
- **THEN** system reports error with remaining space recommendation

#### Scenario: Backup creation failure
- **WHEN** backup fails
- **THEN** system aborts without materializing (data safety preserved)

#### Scenario: Restore failure
- **WHEN** restore from backup fails
- **THEN** system reports error but leaves workspace in safe state

#### Scenario: Cleanup of missing paths
- **WHEN** stop attempts to remove path that user deleted
- **THEN** system succeeds without error (idempotent cleanup)
- **AND** system logs path was already removed