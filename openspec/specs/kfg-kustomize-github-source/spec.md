# github-url-source Specification

## Purpose
TBD - created by archiving change github-url-kpath-support. Update Purpose after archive.
## Requirements
### Requirement: GitHub URL detection

The system MUST detect GitHub URLs in source arguments.

#### Scenario: HTTPS GitHub URL
- **WHEN** user provides `https://github.com/owner/repo`
- **THEN** system detects it as GitHub URL
- **AND** passes it to kustomize loader without preprocessing

#### Scenario: GitHub URL with path separator
- **WHEN** user provides `https://github.com/owner/repo//manifests`
- **THEN** system detects it as GitHub URL
- **AND** passes it to kustomize loader

#### Scenario: GitHub URL with ref parameter
- **WHEN** user provides `https://github.com/owner/repo//manifests?ref=v1.0.0`
- **THEN** system detects it as GitHub URL
- **AND** passes it to kustomize loader

#### Scenario: Non-GitHub URL
- **WHEN** user provides `https://example.com/path`
- **THEN** system treats it as regular URL (kustomize behavior)

### Requirement: GitHub URL processing

The system MUST process GitHub URLs via kustomize's git cloner.

#### Scenario: Clone and build
- **WHEN** user runs `kfg build https://github.com/owner/repo//manifests`
- **THEN** kustomize clones repository to temp directory
- **AND** processes kustomization.yaml
- **AND** outputs resulting YAML

#### Scenario: Shallow clone
- **WHEN** kustomize clones GitHub URL
- **THEN** uses `--depth=1` shallow clone
- **AND** fetches only specified ref

#### Scenario: Branch reference
- **WHEN** user provides `https://github.com/owner/repo//manifests?ref=main`
- **THEN** clones specified branch
- **AND** processes kustomization from that branch

#### Scenario: Tag reference
- **WHEN** user provides `https://github.com/owner/repo//manifests?ref=v1.0.0`
- **THEN** clones specified tag
- **AND** processes kustomization from that tag

### Requirement: GitHub URL error handling

The system MUST handle GitHub URL errors appropriately.

#### Scenario: Invalid repository
- **WHEN** GitHub URL points to non-existent repository
- **THEN** kustomize clone fails
- **AND** exit code 1
- **AND** error message indicates clone failure

#### Scenario: Invalid path
- **WHEN** GitHub URL path component doesn't exist in repo
- **THEN** kustomize fails to find kustomization.yaml
- **AND** exit code 1
- **AND** error message indicates missing kustomization

#### Scenario: Network failure
- **WHEN** network is unavailable
- **THEN** git clone fails
- **AND** exit code 1
- **AND** error message indicates network issue

