# Version Command Specification

## Purpose

Specifies the version command output format and the Makefile build process that injects version metadata into the kfg CLI binary.

## Requirements

### Requirement: kfg --version outputs version information

When the user runs `kfg --version`, the CLI SHALL output version information in the format:
```
kfg version <semver> (<commit>, <date>)
```

Where:
- `<semver>` is the semantic version (e.g., "1.0.09")
- `<commit>` is the short git commit hash (12 characters, e.g., "abc123def456")
- `<date>` is the build timestamp in UTC RFC3339 format (e.g., "2026-04-14T19:00:00Z")

#### Scenario: Running kfg --version with full metadata
- **WHEN** user runs `kfg --version` with a binary built using `make build`
- **THEN** the output SHALL match the pattern `kfg version \d+\.\d+\.\d+ \([a-f0-9]{12}, \d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z\)`

#### Scenario: Running kfg --version without metadata
- **WHEN** user runs `kfg --version` with a binary built using plain `go build` (without ldflags)
- **THEN** the output SHALL show "dev" for version, "unknown" for commit, and "unknown" for date

#### Scenario: Running kfg --version returns exit code 0
- **WHEN** user runs `kfg --version`
- **THEN** the exit code SHALL be 0

### Requirement: kfg --version is a root-level flag

The `--version` flag SHALL be available on the root `kfg` command and not require any subcommand.

#### Scenario: --version flag is available without subcommand
- **WHEN** user runs `kfg --version`
- **THEN** the version information SHALL be displayed directly without requiring a subcommand like `kfg version`

### Requirement: Makefile builds with version metadata

The Makefile SHALL inject version, commit, and date into the Go binary at build time using `-ldflags`.

#### Scenario: make build outputs to ./bin/kfg
- **WHEN** user runs `make build`
- **THEN** the binary SHALL be created at `./bin/kfg`

#### Scenario: make build injects version from flake
- **WHEN** user runs `make build` in a directory with a valid flake.nix
- **THEN** the resulting binary's version SHALL match the version defined in flake.nix

#### Scenario: make build injects git commit hash
- **WHEN** user runs `make build` in a git repository
- **THEN** the resulting binary's commit SHALL be the current HEAD's short hash

### Requirement: flake.nix exposes version for Makefile consumption

The flake.nix SHALL expose the version via a `lib` output that can be queried with `nix eval --raw .#lib.version`.

#### Scenario: nix eval retrieves version from flake
- **WHEN** user runs `nix eval --raw .#lib.version` in the project directory
- **THEN** the output SHALL be the semantic version string (e.g., "1.0.09")

### Requirement: Bats tests validate version output

The Bats test suite SHALL include tests that verify the version command output format.

#### Scenario: Bats test checks --version output
- **WHEN** the Bats test suite is run with `make test-bats`
- **THEN** there SHALL be a test that verifies `kfg --version` returns exit code 0 and contains "kfg version"