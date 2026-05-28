# kpath-env-var Specification

## Purpose
TBD - created by archiving change github-url-kpath-support. Update Purpose after archive.
## Requirements
### Requirement: KFG_KPATH binding

The system MUST bind `KFG_KPATH` to Viper configuration.

#### Scenario: Environment variable read
- **WHEN** `KFG_KPATH` is set in environment
- **THEN** Viper binds the value
- **AND** accessible via `config.GetKPath()`

#### Scenario: Empty environment variable
- **WHEN** `KFG_KPATH` is not set
- **THEN** `config.GetKPath()` returns empty string
- **AND** commands require explicit source

### Requirement: KFG_KPATH fallback priority

The system MUST use `KFG_KPATH` as fallback when no explicit source provided.

#### Scenario: Flag overrides env var
- **WHEN** user runs `kfg apply -k ./manifests`
- **AND** `KFG_KPATH=https://github.com/owner/repo` is set
- **THEN** uses `-k ./manifests` (flag wins)

#### Scenario: Positional arg overrides env var
- **WHEN** user runs `kfg build ./manifests`
- **AND** `KFG_KPATH=https://github.com/owner/repo` is set
- **THEN** uses `./manifests` (arg wins)

#### Scenario: Env var as fallback
- **WHEN** user runs `kfg build` (no arg)
- **AND** `KFG_KPATH=./manifests` is set
- **THEN** uses `KFG_KPATH` value as source

#### Scenario: No source available
- **WHEN** user runs `kfg build` (no arg)
- **AND** `KFG_KPATH` is not set
- **THEN** exit code 1
- **AND** error message: "kustomization source required. Provide a path, use -k flag, or set KFG_KPATH."

### Requirement: KFG_KPATH accepts GitHub URLs

The `KFG_KPATH` MUST accept GitHub URLs as values.

#### Scenario: GitHub URL in env var
- **WHEN** `KFG_KPATH=https://github.com/owner/repo//manifests` is set
- **AND** user runs `kfg apply`
- **THEN** uses GitHub URL as source
- **AND** kustomize clones and processes

#### Scenario: Local path in env var
- **WHEN** `KFG_KPATH=./manifests` is set
- **AND** user runs `kfg apply`
- **THEN** uses local path as source

### Requirement: KFG_KPATH per-command behavior

Each command MUST handle `KFG_KPATH` appropriately.

#### Scenario: Build command
- **WHEN** user runs `kfg build` without argument
- **AND** `KFG_KPATH` is set
- **THEN** uses `KFG_KPATH` as source

#### Scenario: Apply command
- **WHEN** user runs `kfg apply` without `-k` or positional
- **AND** `KFG_KPATH` is set
- **THEN** uses `KFG_KPATH` as kustomize path

#### Scenario: Run command
- **WHEN** user runs `kfg run claude` without `-k`
- **AND** `KFG_KPATH` is set
- **THEN** uses `KFG_KPATH` as kustomize path

