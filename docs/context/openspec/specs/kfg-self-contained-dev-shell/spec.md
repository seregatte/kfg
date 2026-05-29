## ADDED Requirements

### Requirement: Self-contained dev shell toolchain

The kfg `flake.nix` SHALL declare all dev shell tooling inline without inheriting
from any external flake input. The `nixai` input SHALL NOT appear in `flake.nix`
or `flake.lock`.

#### Scenario: Dev shell enters without external inputs

- **WHEN** a developer runs `nix develop` in the kfg repository
- **THEN** the shell SHALL activate successfully using only the kfg flake and its
  declared inputs (nixpkgs)
- AND no network access to `github:seregatte/nixai` SHALL be required

#### Scenario: Dev shell provides Go toolchain

- **WHEN** the dev shell is active
- **THEN** `go` SHALL be available on `PATH`
- AND `go build ./...` in `src/` SHALL succeed

#### Scenario: Dev shell provides test toolchain

- **WHEN** the dev shell is active
- **THEN** `bats` SHALL be available on `PATH`
- AND `make test-bats` SHALL discover and run all Bats test roots

#### Scenario: Dev shell provides YAML tools

- **WHEN** the dev shell is active
- **THEN** `yq`, `jq`, `yajsv`, and `gomplate` SHALL be available on `PATH`

### Requirement: kfg-bundle package

The kfg flake SHALL expose a `kfg-bundle` package that aggregates AI agent tooling
via `symlinkJoin`.

#### Scenario: Bundle contents

- **WHEN** a consumer builds `kfg.packages.${system}.kfg-bundle` (or the dev shell
  is active)
- **THEN** the following commands SHALL be available: `openspec`, `ctx7`,
  `chrome-devtools-mcp`, `pi`, `gws`, `notebooklm`, `nblm`, `claude`, `gemini`,
  `opencode`, `playwright`

#### Scenario: No circular kfg reference in bundle

- **WHEN** `kfg-bundle` is built
- **THEN** it SHALL NOT include a reference to `kfg.packages.${system}.default`
  from an external source
- AND `kfg` binary SHALL come from `self.packages.${system}.default` only

### Requirement: KFG_DIR env var in shell hook

The dev shell `shellHook` SHALL export `KFG_DIR` pointing to the kfg Nix store
derivation output path.

#### Scenario: KFG_DIR is set

- **WHEN** the dev shell is active
- **THEN** `$KFG_DIR` SHALL be set to `${self.outPath}` of the kfg flake
- AND `NIXAI_DIR` SHALL NOT be exported by the kfg `shellHook`

### Requirement: Starship prompt assets

The kfg repository SHALL include Starship prompt configuration files at
`assets/starship/full.toml` and `assets/starship/mobile.toml`, referenced by the
`shellHook` via `$KFG_DIR`.

#### Scenario: Starship config is set in dev shell

- **WHEN** the dev shell is active
- **THEN** `STARSHIP_CONFIG` SHALL point to `${KFG_DIR}/assets/starship/full.toml`
  when `$COLUMNS` is 45 or greater
- AND `STARSHIP_CONFIG` SHALL point to `${KFG_DIR}/assets/starship/mobile.toml`
  when `$COLUMNS` is less than 45
