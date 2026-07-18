## ADDED Requirements

### Requirement: Configurable system tool preference

The kfg devShells MUST provide a mechanism to prefer host-installed executables over Nix devShell versions for a configurable set of commands.

#### Scenario: Host tool wins when installed

- **WHEN** entering any kfg devShell (`default`, `dev`, or `ci`)
- **AND** the host system has a tool matching the default allowlist installed outside `/nix/store`
- **AND** `KFG_PREFER_SYSTEM` is not explicitly set to `0`
- **THEN** `command -v <tool>` SHALL resolve to the host-installed version
- **AND** the Nix store version SHALL still be available as a fallback on PATH

#### Scenario: Nix tool used when host does not have it

- **WHEN** entering any kfg devShell
- **AND** the host system does NOT have a tool matching the default allowlist
- **THEN** `command -v <tool>` SHALL resolve to the Nix store version

#### Scenario: KFG_PREFER_SYSTEM=0 disables preference

- **WHEN** entering any kfg devShell with `KFG_PREFER_SYSTEM=0`
- **THEN** `command -v <tool>` SHALL resolve to the Nix store version as usual

#### Scenario: KFG_PREFER_SYSTEM_COMMANDS overrides default allowlist

- **WHEN** entering any kfg devShell with `KFG_PREFER_SYSTEM_COMMANDS="opencode go"`
- **THEN** only `opencode` and `go` SHALL be resolved to host versions if available
- **AND** other tools (e.g., `node`, `bats`) SHALL resolve to Nix store versions

#### Scenario: ./bin preserves priority in dev shell

- **WHEN** entering `devShells.dev`
- **AND** `./bin/kfg` exists
- **THEN** `command -v kfg` SHALL resolve to `./bin/kfg`
- **AND** neither the host nor Nix store `kfg` SHALL take precedence

#### Scenario: Default allowlist

- **WHEN** `KFG_PREFER_SYSTEM_COMMANDS` is not set
- **THEN** the default allowlist SHALL be: `go node npm npx corepack uv uvx bats openspec pi ctx7 chrome-devtools-mcp gws notebooklm nblm opencode playwright`
