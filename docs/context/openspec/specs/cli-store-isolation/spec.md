## MODIFIED Requirements

### Requirement: Store-scoped operational commands

Internal operational commands that persist runtime data MUST respect store isolation.

#### Scenario: GC commands use selected store root
- **WHEN** user runs a `kfg sys gc` subcommand with a configured `KFG_STORE_DIR`
- **THEN** the command SHALL operate only on cache data rooted under that store directory

#### Scenario: Removed image and workspace command surface
- **WHEN** store-isolated CLI behavior is evaluated
- **THEN** it SHALL apply to runtime cache management commands
- **AND** SHALL NOT require image or workspace command groups to exist
