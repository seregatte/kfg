# kfg-step-cache-atomic-write Specification

## Purpose

Define atomic write semantics for cache entries to prevent partial or corrupted cache state.

## Requirements

### Requirement: Atomic Cache Write
The runtime SHALL write cache entries atomically so that a partial or interrupted write cannot be observed as a valid cache hit. Atomic write is performed entirely by the Go `kfg sys cache store` subcommand.

#### Scenario: Successful cache store is atomic
- **WHEN** a cacheable Step completes execution and its results are written to the cache
- **THEN** the Go `store` subcommand SHALL write all artifacts and `metadata.yaml` to a temporary directory
- **AND** SHALL rename (move) the temporary directory to the final cache entry path only after all writes are complete
- **AND** the final cache entry path SHALL NOT exist during the write phase

#### Scenario: Interrupted cache store leaves no partial entry
- **WHEN** the cache write is interrupted before the final rename (e.g., process killed)
- **THEN** the final cache entry path SHALL NOT exist
- **AND** a subsequent execution SHALL find a cache miss and execute the Step normally
- **AND** `kfg sys cache exists` SHALL return exit code 1 for the interrupted entry

#### Scenario: Stale temp directory is cleaned on next store
- **WHEN** a cache store begins for a given Step
- **THEN** the Go `store` subcommand SHALL remove any pre-existing temporary directory at `<cache_path>.tmp` before beginning the write

#### Scenario: Atomic write performed by Go subcommand
- **WHEN** the shell helper `__kfg_cache_store` is called
- **THEN** it SHALL delegate to `kfg sys cache store` which performs the atomic write internally
- **AND** the shell SHALL NOT perform any file operations for cache storage

