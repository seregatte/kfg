## MODIFIED Requirements

### Requirement: Development devShell maintains workflow

The `devShells.dev` MUST maintain the full kfg contributor workflow without overriding OpenSpec root resolution.

#### Scenario: Developer uses explicit dev shell
- **WHEN** running `nix develop .#dev`
- **THEN** the shell MUST run `source <(go run ./src/cmd/kfg apply -k packages/domains/ai-agents/overlays/dev)`
- **AND** Go development tooling and `kfg-bundle` MUST be available

#### Scenario: Development shell sets environment
- **WHEN** entering `devShells.dev`
- **THEN** `./bin` MUST be prioritized in `PATH`
- **AND** `KFG_DIR` MUST be set to kfg's outPath
- **AND** `STARSHIP_CONFIG` MUST be selected dynamically
- **AND** `OPENSPEC_ROOT_DIR` MUST NOT be set

#### Scenario: Development shell resolves OpenSpec root
- **WHEN** a contributor runs OpenSpec from the repository or a nested directory
- **THEN** OpenSpec SHALL resolve the repository-root `openspec/` directory natively
- **AND** the generated command wrapper SHALL NOT change the working directory before delegation

#### Scenario: CI shell maintains test workflow
- **WHEN** running `nix develop .#ci`
- **THEN** the shell MUST NOT include `kfg-bundle`
- **AND** `pkgs.go` and `pkgs.gnumake` MUST be in buildInputs
- **AND** shellHook MUST set up the Bats test helper vendor directory
- **AND** shellHook MUST NOT set `OPENSPEC_ROOT_DIR`
