## Why

The AI agents domain exposes a growing catalog of agents, integrations, prompts, converters, and workflows, but consumers and the wizard currently depend on hardcoded knowledge instead of stable, capability-specific usage contracts. OpenSpec 1.6 also makes the repository's `OPENSPEC_ROOT_DIR` redirection and ignored root scaffold obsolete, creating an opportunity to adopt native root discovery while making OpenSpec's local, external-store, and reference models composable through the domain.

## What Changes

- Establish a public documentation contract for every independently consumable AI agents capability, including its Kustomization entrypoint, exported resources, compatibility, configurable fields, JSON Patch examples, artifacts, cleanup behavior, limitations, and validation commands.
- Add a domain-level capability index and selective Kustomization entrypoints while preserving aggregate package entrypoints.
- Redesign the AI wizard to discover capabilities through the domain documentation and adjacent Kustomizations instead of maintaining a hardcoded catalog.
- Make the wizard generate selective compositions and documented JSON Patches, and report unsupported capability and agent combinations.
- Add OpenSpec configuration entrypoints for repository-local planning, external store pointers, and read-only store references.
- Migrate kfg's canonical OpenSpec root from `docs/context/openspec/` to the standard repository-root `openspec/` location.
- Remove `OPENSPEC_ROOT_DIR` from development shells and make the OpenSpec command wrapper rely on native root resolution.
- Update OpenSpec skill installation so it does not create or remove a disposable project root while installing agent integration files.
- Correct the OpenSpec configuration rules shape for current OpenSpec versions as part of the root migration.
- Add a mandatory rule to `docs/AGENTS.md` requiring capability-specific documentation to change whenever an AI agents domain resource's public usage changes.
- Add automated checks for public entrypoint builds, documented JSON Patch examples, documentation coverage, wizard discovery instructions, OpenSpec modes, and generated artifacts.
- **BREAKING**: The canonical repository OpenSpec path changes from `docs/context/openspec/` to `openspec/`, and development shells no longer export `OPENSPEC_ROOT_DIR`.

## Capabilities

### New Capabilities

- `domain-ai-agents-capability-documentation`: Defines the domain capability index, capability-specific usage documentation, selective composition entrypoints, and JSON Patch documentation contract.
- `domain-ai-agents-documentation-driven-wizard`: Defines filesystem-based capability discovery, compatibility-aware selection, selective composition, and documented patch generation by the AI wizard.
- `domain-ai-agents-openspec-configuration`: Defines local planning, external store pointer, and read-only reference composition for OpenSpec consumers.

### Modified Capabilities

- `kfg-domain-package-contract`: Requires public domain capabilities to expose synchronized usage documentation and independently buildable entrypoints where selective consumption is supported.
- `kfg-project-structure`: Moves the canonical OpenSpec root to `openspec/` and records the AI agents domain documentation layout.
- `kfg-openspec-cmd`: Removes environment-driven directory changes and delegates root resolution to the OpenSpec CLI.
- `kfg-build-nix-packaging`: Removes `OPENSPEC_ROOT_DIR` from development and CI shell hooks.
- `kfg-devshell-consumer`: Updates development shell behavior to use native OpenSpec root discovery.
- `domain-ai-agents-skill-installation-steps`: Updates OpenSpec integration installation to preserve the selected project root and install agent files without disposable root cleanup.

## Impact

- Affects `packages/domains/ai-agents/` manifests, capability entrypoints, READMEs, tests, and both AI agent overlays.
- Reworks `packages/domains/ai-agents/overlays/ai/assets/prompts/wizard.yaml` and its command guidance.
- Moves all tracked files under `docs/context/openspec/` to `openspec/`, including active and archived changes, and updates repository-wide references.
- Changes `flake.nix`, `.gitignore`, `docs/AGENTS.md`, and the OpenSpec command and installation manifests.
- Preserves aggregate AI agents entrypoints while adding selective entrypoints; consumers can migrate incrementally.
- Requires Bats coverage for domain documentation and Kustomization contracts plus the existing Go, formatting, lint, vet, and integration suites.
