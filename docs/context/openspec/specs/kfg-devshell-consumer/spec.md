## ADDED Requirements

### Requirement: Consumer-friendly devShell

The `devShells.default` MUST provide a shell environment suitable for consuming kfg without Go development workflow.

#### Scenario: Consumer inherits devShell without Go error

- **WHEN** a project uses `inputsFrom = [kfg.devShells.${system}.default]`
- **THEN** the shell MUST NOT run `source <(go run ./src/cmd/kfg apply)`
- **AND** the shell MUST NOT trigger "go: cannot find main module" errors

#### Scenario: Consumer shell provides tools

- **WHEN** running `nix develop` (implicit default shell)
- **THEN** kfg-bundle MUST be available in PATH
- **AND** devInputs tools (yq-go, jq, gomplate, bats, etc.) MUST be available

#### Scenario: Consumer shell sets environment

- **WHEN** entering default devShell
- **THEN** `KFG_DIR` MUST be set to kfg's outPath
- **AND** `STARSHIP_CONFIG` MUST be set dynamically (mobile vs full based on terminal width)
- **AND** `OPENSPEC_ROOT_DIR` MUST NOT be set (consumers set their own)

#### Scenario: Consumer shell excludes development tools

- **WHEN** running `nix develop`
- **THEN** `pkgs.go` MUST NOT be in buildInputs
- **AND** `pkgs.nodejs` MUST NOT be in buildInputs (already in kfg-bundle)

### Requirement: Development devShell maintains workflow

The `devShells.dev` MUST maintain the full kfg development workflow for kfg contributors.

#### Scenario: Developer uses explicit dev shell

- **WHEN** running `nix develop .#dev`
- **THEN** the shell MUST run `source <(go run ./src/cmd/kfg apply -k packages/domains/ai-agents/overlays/dev)`
- **AND** `pkgs.go` MUST be in buildInputs
- **AND** `pkgs.nodejs` MUST be in buildInputs

#### Scenario: Development shell sets environment

- **WHEN** entering dev devShell
- **THEN** `PATH="./bin:$PATH"` MUST be set
- **AND** `OPENSPEC_ROOT_DIR` MUST be set to `docs/context`
- **AND** `KFG_DIR` MUST be set to kfg's outPath
- **AND** `STARSHIP_CONFIG` MUST be set dynamically

#### Scenario: CI shell unchanged

- **WHEN** running `nix develop .#ci`
- **THEN** the shell MUST NOT include kfg-bundle
- **AND** `pkgs.go` and `pkgs.gnumake` MUST be in buildInputs
- **AND** shellHook MUST set up bats test helpers vendor directory
