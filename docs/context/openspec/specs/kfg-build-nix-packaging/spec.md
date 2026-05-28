# kfg Nix Packaging Specification

## Purpose

The kfg repository provides a Nix flake that exposes pre-built binaries from GitHub Releases. This specification defines how the flake structure works for consuming kfg as a Nix package.

## Requirements

### Requirement: kfg flake.nix

The kfg repository SHALL include a `flake.nix` that exposes the kfg binary as a Nix package.

#### Scenario: Package output
- **WHEN** a consumer runs `nix build github:seregatte/kfg`
- **THEN** the output SHALL contain `bin/kfg`
- AND the binary SHALL be the pre-built release binary for the current platform

#### Scenario: Multi-platform support
- **WHEN** the flake is evaluated on any supported system
- **THEN** `packages.default` SHALL be available for:
  - `x86_64-linux`
  - `aarch64-linux`
  - `x86_64-darwin`
  - `aarch64-darwin`

### Requirement: Binary source from GitHub Releases

The flake SHALL fetch pre-built binaries from GitHub Releases using `fetchurl`.

#### Scenario: Fetch URL construction
- **WHEN** the flake is evaluated
- **THEN** the URL SHALL be constructed as:
  `https://github.com/seregatte/kfg/releases/download/v<version>/kfg_<version>_<os>_<arch>.tar.gz`
- AND `<version>` SHALL come from the hardcoded `version` attribute
- AND `<os>` and `<arch>` SHALL come from `platformArchiveNames` mapping

#### Scenario: Hash verification
- **WHEN** a binary is fetched
- **THEN** the hash SHALL be verified against the corresponding `platformHashes` entry
- AND a hash mismatch SHALL cause the build to fail

### Requirement: Installation derivation

The flake SHALL install the fetched binary into the Nix store.

#### Scenario: Binary installation
- **WHEN** the derivation is built
- **THEN** the `kfg` binary from the tar.gz SHALL be installed to `$out/bin/kfg`
- AND the binary SHALL have execute permissions (mode 755)

### Requirement: Version exposure

The flake SHALL expose the version string for downstream consumers.

#### Scenario: lib.version
- **WHEN** a consumer accesses `kfg.lib.version`
- **THEN** the value SHALL match the hardcoded `version` attribute in the flake
- AND nixai MAY use this to track the kfg version it consumes

### Requirement: Platform mapping

The flake SHALL map Nix system identifiers to GoReleaser archive name components.

#### Scenario: System to archive name
- **WHEN** the flake is evaluated for a given system
- **THEN** `platformArchiveNames` SHALL map:
  - `x86_64-linux` → `linux_amd64`
  - `aarch64-linux` → `linux_arm64`
  - `x86_64-darwin` → `darwin_amd64`
  - `aarch64-darwin` → `darwin_arm64`