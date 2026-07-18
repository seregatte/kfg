## MODIFIED Requirements

### Requirement: Dev devShell maintains Go workflow

The `devShells.dev` MUST maintain all Go development capabilities while allowing OpenSpec to resolve the repository's standard root natively.

#### Scenario: Dev devShell buildInputs
- **WHEN** inspecting `devShells.dev` buildInputs in `flake.nix`
- **THEN** `buildInputs` MUST include `devInputs`, `pkgs.go`, and `kfg-bundle`
- **AND** all required development tools SHALL remain available

#### Scenario: Dev devShell shellHook maintains workflow
- **WHEN** inspecting `devShells.dev` shellHook in `flake.nix`
- **THEN** shellHook MUST prioritize `./bin`
- **AND** shellHook MUST run `source <(go run ./src/cmd/kfg apply -k packages/domains/ai-agents/overlays/dev)`
- **AND** shellHook MUST NOT set `OPENSPEC_ROOT_DIR`

#### Scenario: CI shell uses native OpenSpec discovery
- **WHEN** inspecting `devShells.ci` shellHook in `flake.nix`
- **THEN** shellHook MUST NOT set `OPENSPEC_ROOT_DIR`
- **AND** OpenSpec commands run by CI SHALL discover the repository-root `openspec/` directory without an environment redirect
