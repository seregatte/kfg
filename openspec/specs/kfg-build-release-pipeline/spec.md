# kfg Release Pipeline Specification

## Purpose

The kfg repository uses GitHub Actions and GoReleaser for automated cross-platform binary releases. This specification defines the release workflow and automated Nix hash update process.

## Requirements

### Requirement: GoReleaser configuration

The kfg repository SHALL include a `.goreleaser.yml` that builds cross-platform binaries.

#### Scenario: Release archive structure
- **WHEN** a release is triggered by a `v*` tag
- **THEN** GoReleaser SHALL produce 4 archives:
  - `kfg_<version>_linux_amd64.tar.gz`
  - `kfg_<version>_linux_arm64.tar.gz`
  - `kfg_<version>_darwin_amd64.tar.gz`
  - `kfg_<version>_darwin_arm64.tar.gz`
- AND each archive SHALL contain a single `kfg` binary

#### Scenario: Version metadata injection
- **WHEN** GoReleaser builds the binary
- **THEN** ldflags SHALL inject:
  - `-X main.version=<semver>` (e.g., `2.1.0`)
  - `-X main.commit=<short commit hash>`
  - `-X main.date=<ISO 8601 UTC timestamp>`
- AND binary size SHALL be reduced with `-s -w` flags

### Requirement: Release workflow

The kfg repository SHALL include a GitHub Actions workflow that automates the release process.

#### Scenario: Release trigger
- **WHEN** a tag matching `v*` is pushed
- **THEN** the `release` job SHALL execute
- AND GoReleaser SHALL build and upload all platform archives to GitHub Releases

#### Scenario: Release permissions
- **WHEN** the release workflow runs
- **THEN** the workflow SHALL have `contents: write` permission
- AND `GITHUB_TOKEN` SHALL be used for authentication

### Requirement: Automated flake hash update

The release workflow SHALL automatically update `flake.nix` hashes after a successful release.

#### Scenario: Hash computation
- **WHEN** the `release` job completes successfully
- **THEN** the `update-flake` job SHALL run
- AND for each of the 4 platform archives, `nix store prefetch-file` SHALL compute the SHA-256 hash
- AND hashes SHALL be extracted from the JSON output via `jq -r '.hash'`

#### Scenario: Flake.nix update
- **WHEN** all 4 hashes are computed
- **THEN** the `version` field in `flake.nix` SHALL be updated to match the tag
- AND all 4 `platformHashes` entries SHALL be updated with computed values

#### Scenario: Pull request creation
- **WHEN** `flake.nix` has been updated
- **THEN** a pull request SHALL be created via `peter-evans/create-pull-request`
- AND the PR title SHALL be `flake: update to v<version>`
- AND the PR branch SHALL be `flake-update/v<version>`
- AND the PR body SHALL indicate the automated update

### Requirement: CI workflow (existing)

The existing CI workflow (`ci.yml`) SHALL continue to run on push/PR to `main` with build, test, and vet steps. It is NOT changed by this capability.