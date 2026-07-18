# cli-store-isolation Specification

## Purpose
TBD - created by archiving change bats-isolated-store. Update Purpose after archive.
## Requirements
### Requirement: Store isolation for testing
The system SHALL support isolated store directories for test execution.

#### Scenario: Store flag functional
- **WHEN** `--store <path>` flag is provided
- **THEN** all store operations use the specified directory
- **AND** default store is not affected

#### Scenario: Store flag inheritance
- **WHEN** `--store` is set on parent `store` command
- **THEN** all subcommands (image, workspace) inherit the store directory

#### Scenario: Empty store directory
- **WHEN** isolated store directory is empty
- **THEN** store commands create necessary subdirectories
- **AND** no errors occur for missing directories

### Requirement: Store-scoped operational commands

Internal operational commands that persist runtime data MUST respect store isolation.

#### Scenario: GC commands use selected store root
- **WHEN** user runs a `kfg sys gc` subcommand with a configured `KFG_STORE_DIR`
- **THEN** the command SHALL operate only on cache data rooted under that store directory

#### Scenario: Removed image and workspace command surface
- **WHEN** store-isolated CLI behavior is evaluated
- **THEN** it SHALL apply to runtime cache management commands
- **AND** SHALL NOT require image or workspace command groups to exist

