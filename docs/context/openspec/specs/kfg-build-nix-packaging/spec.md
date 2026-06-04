## ADDED Requirements

### Requirement: Version reflects devShell contract change

The kfg version MUST be bumped to signal breaking changes in devShell contract.

#### Scenario: Version bump from 0.0.6 to 0.0.7

- **WHEN** this change is implemented
- **THEN** `version` in `flake.nix` MUST be set to `"0.0.7"`
- **AND** version reflects breaking change for kfg development workflow

#### Scenario: Developers must use explicit dev shell

- **WHEN** developing kfg after this change
- **THEN** developers MUST use `nix develop .#dev` (not implicit default)
- **AND** using `nix develop` without explicit shell will NOT run Go development workflow

### Requirement: Default devShell excludes Go development

The `devShells.default` MUST NOT include Go development tools or workflow after this change.

#### Scenario: Default devShell buildInputs

- **WHEN** inspecting `devShells.default` buildInputs in flake.nix
- **THEN** `buildInputs` MUST be `devInputs ++ [kfg-bundle]`
- **AND** `pkgs.go` MUST NOT be in the list
- **AND** `pkgs.nodejs` MUST NOT be in the list

#### Scenario: Default devShell shellHook

- **WHEN** inspecting `devShells.default` shellHook in flake.nix
- **THEN** shellHook MUST set `KFG_DIR`
- **AND** shellHook MUST set `STARSHIP_CONFIG` dynamically
- **AND** shellHook MUST NOT set `OPENSPEC_ROOT_DIR`
- **AND** shellHook MUST NOT run `source <(go run ./src/cmd/kfg apply)`

### Requirement: Dev devShell maintains Go workflow

The `devShells.dev` MUST maintain all Go development capabilities after this change.

#### Scenario: Dev devShell buildInputs unchanged

- **WHEN** inspecting `devShells.dev` buildInputs in flake.nix
- **THEN** `buildInputs` MUST be `devInputs ++ [pkgs.nodejs pkgs.go kfg-bundle]`
- **AND** all development tools remain available

#### Scenario: Dev devShell shellHook maintains workflow

- **WHEN** inspecting `devShells.dev` shellHook in flake.nix
- **THEN** shellHook MUST set `PATH="./bin:$PATH"`
- **AND** shellHook MUST set `OPENSPEC_ROOT_DIR=docs/context`
- **AND** shellHook MUST run `source <(go run ./src/cmd/kfg apply -k packages/domains/ai-agents/overlays/dev)`

## MODIFIED Requirements

### Requirement: Version exposure

The flake SHALL expose the version string for downstream consumers.

#### Scenario: lib.version

- **WHEN** a consumer accesses `kfg.lib.version`
- **THEN** the value SHALL match the hardcoded `version` attribute in the flake

## REMOVED Requirements

### Requirement: nixai dependency via inputsFrom

**Reason**: kfg now declares all dev shell tooling inline. The `nixai` input and
the `inputsFrom` delegation are removed as part of the absorb-nixai change.

**Migration**: No migration needed. Developers run `nix develop` as before; the
shell is now self-contained. No reference to `github:seregatte/nixai` is required.
