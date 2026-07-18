# Apply Command Specification

## Purpose

Define modifications to the `kfg apply` command for optional `-k` flag and GitHub URL support.

## MODIFIED Requirements

### Requirement: Apply command syntax

The CLI MUST provide `kfg apply` for shell generation with optional source via `KFG_KPATH`.

#### Scenario: Apply from kustomization
- **WHEN** user runs `kfg apply -k .kfg/overlay/dev`
- **THEN** loads kustomization from path
- **AND** resolves workflow
- **AND** generates shell functions to stdout

#### Scenario: Apply from file
- **WHEN** user runs `kfg apply -f manifest.yaml`
- **THEN** loads manifest from file
- **AND** generates shell functions

#### Scenario: Apply from stdin
- **WHEN** user runs `kfg apply -f -`
- **THEN** reads manifest from stdin
- **AND** generates shell functions

#### Scenario: Apply from GitHub URL
- **WHEN** user runs `kfg apply -k https://github.com/owner/repo//manifests`
- **THEN** clones GitHub repository
- **AND** processes kustomization.yaml
- **AND** resolves workflow
- **AND** generates shell functions

#### Scenario: Apply from GitHub URL with ref
- **WHEN** user runs `kfg apply -k https://github.com/owner/repo//manifests?ref=v1.0.0`
- **THEN** clones specified tag
- **AND** processes kustomization.yaml
- **AND** generates shell functions

#### Scenario: Apply without flags with KFG_KPATH
- **WHEN** user runs `kfg apply` (no `-k` or `-f`)
- **AND** `KFG_KPATH=./manifests` is set
- **THEN** uses `KFG_KPATH` as kustomize path
- **AND** generates shell functions

#### Scenario: Apply without flags or KFG_KPATH
- **WHEN** user runs `kfg apply` (no `-k` or `-f`)
- **AND** `KFG_KPATH` is not set
- **THEN** exit code 1
- **AND** error message: "kustomization source required. Provide a path, use -k flag, or set KFG_KPATH."

### Requirement: Flag validation

The CLI MUST validate flag combinations.

#### Scenario: No flags with KFG_KPATH
- **WHEN** user runs `kfg apply` without `-k` or `-f`
- **AND** `KFG_KPATH` is set
- **THEN** uses `KFG_KPATH` as source
- **AND** succeeds

#### Scenario: No flags without KFG_KPATH
- **WHEN** user runs `kfg apply` without `-k` or `-f`
- **AND** `KFG_KPATH` is not set
- **THEN** exit code 1
- **AND** error message indicates source required

#### Scenario: Mutual exclusion
- **WHEN** user runs `kfg apply -k path -f file`
- **THEN** exit code 2 (usage error)
- **AND** error message indicates flag conflict

#### Scenario: Flag overrides KFG_KPATH
- **WHEN** user runs `kfg apply -k ./manifests`
- **AND** `KFG_KPATH=https://github.com/owner/repo` is set
- **THEN** uses `-k ./manifests` (flag wins)

### Requirement: Exit codes

The CLI MUST use consistent exit codes.

#### Scenario: Success
- **WHEN** apply succeeds
- **THEN** exit code 0

#### Scenario: Runtime error
- **WHEN** resolution or generation fails
- **THEN** exit code 1

#### Scenario: Usage error
- **WHEN** invalid flag combination
- **THEN** exit code 2

#### Scenario: GitHub clone failure
- **WHEN** GitHub URL fails to clone
- **THEN** exit code 1
- **AND** error message indicates clone issue

#### Scenario: No source provided
- **WHEN** no `-k`, no `-f`, and no `KFG_KPATH`
- **THEN** exit code 1
- **AND** error message indicates source required