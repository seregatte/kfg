# sys gc command

## Purpose

The `kfg sys gc` command group provides garbage collection and management
for persisted Step cache entries stored under KFG_STORE_DIR/cache.

## Scope

This spec describes the subcommands of `kfg sys gc`:
- `ls`: List cache entries with metadata
- `inspect`: Show detailed metadata for a cache entry
- `rm`: Remove specific cache entries
- `prune`: Remove old or unused cache entries
- `du`: Show disk usage of cache entries

## Requirements

### Requirement: List cache entries

The `kfg sys gc ls` subcommand SHALL list all cache entries with stable
identifiers and operational metadata.

#### Scenario: List entries with metadata
- **WHEN** user runs `kfg sys gc ls`
- **THEN** the CLI SHALL display each cache entry with:
  - ID: The stable identifier (hash) for the cache entry
  - Step Ref Name: The workflow step reference name
  - Timestamp: When the cache entry was created
  - Size: Disk usage in bytes

### Requirement: Inspect cache entry

The `kfg sys gc inspect` subcommand SHALL show detailed metadata for a
specific cache entry.

#### Scenario: Inspect entry details
- **WHEN** user runs `kfg sys gc inspect <id>`
- **THEN** the CLI SHALL display:
  - Cache entry ID and path
  - Step reference name
  - Timestamp
  - Disk usage
  - Artifacts list
  - Output metadata (if present)

### Requirement: Remove cache entry

The `kfg sys gc rm` subcommand SHALL remove specified cache entries from
storage.

#### Scenario: Remove single entry
- **WHEN** user runs `kfg sys gc rm <id>`
- **THEN** the CLI SHALL remove the specified cache entry from storage

#### Scenario: Remove multiple entries
- **WHEN** user runs `kfg sys gc rm <id1> <id2> ...`
- **THEN** the CLI SHALL remove all specified cache entries

### Requirement: Prune old entries

The `kfg sys gc prune` subcommand SHALL remove cache entries according to
the implemented prune policy.

#### Scenario: Prune entries older than 30 days
- **WHEN** user runs `kfg sys gc prune`
- **THEN** the CLI SHALL remove cache entries older than 30 days
- **AND** SHALL display which entries were pruned

### Requirement: Show disk usage

The `kfg sys gc du` subcommand SHALL report disk usage for persisted
cache entries.

#### Scenario: Show disk usage summary
- **WHEN** user runs `kfg sys gc du`
- **THEN** the CLI SHALL display:
  - Cache directory location
  - Per-entry disk usage
  - Total disk usage
  - Number of cache entries

## Store Directory

The cache entries are stored under `KFG_STORE_DIR/cache` (defaults to
`~/.kfg/store/cache`).

The cache entry directory structure is:
```
KFG_STORE_DIR/cache/
  <entry-hash>/
    metadata.yaml
    artifacts/
      <artifact-files>
    artifact_paths.txt
```

## Implementation Notes

- Each cache entry is identified by a stable identifier derived from Step reference name
- The cache identity uses `StepReference.name` only (no additional components)
- The `metadata.yaml` file contains:
  - `stepRefName`: The workflow step reference name
  - `timestamp`: When the cache entry was created
  - `output`: Optional output metadata (name, valueEncoded)
  - `artifacts`: List of cached artifact relative paths
- The prune policy currently removes entries older than 30 days