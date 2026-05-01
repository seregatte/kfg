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

