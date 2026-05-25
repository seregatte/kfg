# Artifact-Scoped Backup Specification

## Purpose

Specifies that workspace backup operations SHALL be scoped to only files that exist in the workspace AND will be overwritten by image materialization. This minimizes backup overhead and preserves context about which workspace files were actually affected.

## Requirements

### Requirement: Scoped backup computation
The system SHALL compute which workspace paths conflict with image artifacts before backup.

#### Scenario: Compute conflicting paths
- **WHEN** preparing to materialize image with artifacts `.pi/config.json` and `CLAUDE.md`
- **THEN** system checks which of these paths exist in workspace
- **AND** system identifies only existing paths as conflicting

#### Scenario: No conflicting paths
- **WHEN** image artifacts are `.pi/config.json` and `CLAUDE.md`
- **AND** workspace does not contain `.pi/` or `CLAUDE.md`
- **THEN** system identifies no conflicting paths
- **AND** system skips backup creation

#### Scenario: Partial conflicts
- **WHEN** image artifacts are `.pi/config.json`, `CLAUDE.md`, and `GEMINI.md`
- **AND** workspace contains `CLAUDE.md` but not `.pi/` or `GEMINI.md`
- **THEN** system identifies only `CLAUDE.md` as conflicting
- **AND** system backs up only `CLAUDE.md`

### Requirement: Conflict-only backup
The system SHALL backup only workspace files that will be overwritten by image materialization.

#### Scenario: Backup conflicting file
- **WHEN** image artifact `CLAUDE.md` will overwrite existing workspace file
- **THEN** system backs up existing `CLAUDE.md` to instance backup directory
- **AND** system logs "[START] Backing up CLAUDE.md (conflicts with image artifact)"

#### Scenario: Backup conflicting directory
- **WHEN** image artifact `.pi/config.json` will overwrite path in existing `.pi/` directory
- **THEN** system backs up `.pi/config.json` specifically
- **AND** system does NOT backup entire `.pi/` directory if `.pi/other.txt` is not in image

#### Scenario: Skip non-conflicting files
- **WHEN** workspace contains `README.md` and `Makefile`
- **AND** image artifacts do not include these paths
- **THEN** system does NOT backup `README.md` or `Makefile`
- **AND** system logs "[START] No conflicts found - backup skipped"

#### Scenario: Backup directory structure preservation
- **WHEN** backing up conflicting file `.pi/config.json`
- **THEN** system preserves directory structure in backup: `backup/data/.pi/config.json`
- **AND** system does NOT flatten backup to single directory

### Requirement: Artifact paths source
The system SHALL use image metadata Files manifest as the authoritative source for artifact paths.

#### Scenario: Load artifact paths from metadata
- **WHEN** computing backup scope
- **THEN** system loads `ImageMetadata.Files` map
- **AND** system uses map keys as artifact path list

#### Scenario: Missing metadata fallback
- **WHEN** image metadata cannot be loaded
- **THEN** system fails with error before modifying workspace
- **AND** system does NOT proceed with materialization

### Requirement: Backup logging
The system SHALL log which paths are backed up for user visibility.

#### Scenario: Log backed up paths
- **WHEN** backup is created for conflicting paths
- **THEN** system logs each backed up path individually
- **AND** system logs total count: "[START] Backed up 3 conflicting files"

#### Scenario: Log skipped backup
- **WHEN** no conflicting paths found
- **THEN** system logs "[START] No conflicting files - backup skipped"
- **AND** system proceeds with materialization