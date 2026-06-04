## Why

kfg's current `devShells.default` and `devShells.dev` are identical — both inherit the Go development shellHook (`source <(go run ./src/cmd/kfg apply)`). This shellHook is designed for **developing kfg itself**, not for **consuming kfg**. When external projects use `inputsFrom = [kfg.devShells.${system}.default]`, they inherit this Go development shellHook and trigger errors in non-Go projects. The proper solution is to split devShells by role: default → consumer-friendly, dev → development-only.

## What Changes

- Swap roles: `devShells.default` becomes consumer-friendly (no Go shellHook); `devShells.dev` becomes development-only (maintains Go workflow)
- Remove Go development shellHook from `default`: no `source <(go run ./src/cmd/kfg apply)`
- Remove `pkgs.go` and `pkgs.nodejs` from `default` buildInputs (consumers don't need Go/Node for development; tools already in kfg-bundle)
- Remove `OPENSPEC_ROOT_DIR=docs/context` from `default` shellHook (consumers set their own)
- Keep `STARSHIP_CONFIG` dynamic setup in both `default` and `dev` (approved by user)
- Keep `KFG_DIR` environment variable in both
- Update `dev` shellHook: maintain Go development workflow, verify `PATH="./bin:$PATH"` present
- **BREAKING**: Bump version 0.0.6 → 0.0.7 (breaking change for kfg development workflow — developers must use explicit `.#dev`)
- Update `README.md`: clarify devShell usage (default for consumers, dev for development)
- Update `docs/AGENTS.md`: document development workflow with explicit `.#dev`

## Capabilities

### New Capabilities

- `kfg-devshell-consumer`: Consumer-friendly devShell providing tools and environment setup without Go development workflow

### Modified Capabilities

- `kfg-build-nix-packaging`: Version bump from 0.0.6 to 0.0.7; `devShells.default` no longer includes Go development shellHook or tools

## Impact

- `flake.nix` — primary change; swap default/dev roles, version bump
- `README.md` — add devShell usage docs clarifying consumer vs development
- `docs/AGENTS.md` — update development workflow instructions with explicit `.#dev`
- External projects (dotfiles, homelab, toolbox-stack, lifeos) — benefit from no Go errors when using `inputsFrom = [kfg.devShells.${system}.default]`