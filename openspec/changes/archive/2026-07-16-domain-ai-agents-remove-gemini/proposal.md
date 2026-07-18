## Why

The Gemini CLI agent (`gemini-cli-bin`) adds maintenance overhead without proportional value. The wizard overlay already only supports `pi` and `opencode`, and the ctx7 install step has a `gemini` branch that duplicates logic. Dropping support simplifies the manifest catalog, reduces test surface, and removes an unused nix dependency.

## What Changes

- **Remove**: Gemini agent directory `manifests/agents/gemini/` (assets, converters, kustomization)
- **Remove**: Gemini Cmd definition `ai.gemini.cmd.main` from `manifests/cmds/agents.yaml`
- **Remove**: Gemini branch from `ctx7.steps.install` agent detection
- **Remove**: `pkgs.gemini-cli-bin` from `flake.nix` kfg-bundle
- **Remove**: `/.gemini/` entry from `.gitignore`
- **Update**: Agent registry in `manifests/agents/kustomization.yaml`
- **Update**: Documentation in `manifests/README.md` and `manifests/agents/README.md`
- **Update**: All test files referencing gemini (Bats + Go)
- **Update**: Comment-level references in Go source files

## Capabilities

### Removed Capabilities

- `domain-ai-agents-gemini-agent`: The entire Gemini agent support — no Cmd, no converters, no settings, no nix package

## Impact

- **Domain `ai-agents`**: Fewer agent directories; cleaner agent registry
- **Engine**: One less nix dependency in `flake.nix`; fewer test fixtures
- **Framework**: Framework tests lose gemini scaffolding cases
- **Backward compatibility**: Breaking — any `kustomization.yaml` referencing gemini will fail. Users must migrate to opencode or pi.
