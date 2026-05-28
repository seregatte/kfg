# Source Organization Specification

## Purpose

This specification defines how Go source code is organized in the project, including directory structure and import path conventions.

## Requirements

### Requirement: Source code directory structure

All Go source code MUST be organized under the `src/` directory.

#### Scenario: CLI source location
- **WHEN** CLI source code is placed in the project
- **THEN** `src/cmd/kfg/` SHALL contain all CLI entry point files
- AND the main entry point SHALL be at `src/cmd/kfg/main.go`

#### Scenario: Internal packages location
- **WHEN** internal packages are organized
- **THEN** `src/internal/<package>/` SHALL contain each internal package
- AND package names SHALL follow Go naming conventions (lowercase, no underscores)

#### Scenario: Test files co-location
- **WHEN** test files are created for a package
- **THEN** `*_test.go` files SHALL be placed in the same directory as the source files
- AND test files SHALL NOT be in a separate `tests/` directory for Go packages

#### Scenario: Source code not at root
- **WHEN** the project is organized
- **THEN** no `.go` source files SHALL exist at the project root
- AND no `cmd/` directory SHALL exist at the project root
- AND no `internal/` directory SHALL exist at the project root

### Requirement: Go module configuration

The Go module configuration MUST remain at project root for standard Go tooling compatibility.

#### Scenario: go.mod location
- **WHEN** Go module configuration is defined
- **THEN** `go.mod` SHALL reside at the project root
- AND `go.sum` SHALL reside at the project root

#### Scenario: Import paths
- **WHEN** Go files import internal packages
- **THEN** import paths SHALL use the module name as prefix
- AND import statements SHALL include `src/` in the path
- AND imports SHALL follow the pattern `kfg/src/internal/<package>`