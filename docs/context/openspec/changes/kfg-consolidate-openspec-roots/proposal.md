## Why

The project currently maintains four separate OpenSpec roots under `docs/context/` (`openspec/`, `kfg/openspec/`, `framework/openspec/`, `domains/ai-agents/openspec/`), each with its own `config.yaml` and independent `specs/` and `changes/` trees. This structure emerged from an incomplete migration (`separate-engine-from-domain`) and diverges from the OpenSpec-recommended monorepo pattern of a single root with nested spec paths. The fragmentation forces agents to sync sibling changes across multiple roots, leaves orphaned delta specs in the legacy root, and requires `OPENSPEC_ROOT_DIR=docs/context` to point at the parent rather than the actual root.

## What Changes

- All specs are moved into a single root at `docs/context/openspec/` using layer subfolders: `specs/kfg/`, `specs/framework/`, `specs/domain-ai-agents/`
- All changes (active and archived) are merged into `docs/context/openspec/changes/` with a layer prefix on the slug (`kfg-*`, `framework-*`, `domain-ai-agents-*`)
- The six orphaned delta specs in `docs/context/openspec/specs/` are merged into their corresponding canonical specs in `specs/kfg/` before removal
- The three separate `config.yaml` files (kfg, framework, domain) are consolidated into a single `docs/context/openspec/config.yaml` covering all layers
- `OPENSPEC_ROOT_DIR` in `flake.nix` is updated from `docs/context` to `docs/context/openspec`
- The old roots (`docs/context/kfg/`, `docs/context/framework/`, `docs/context/domains/`) are removed
- All references in `AGENTS.md` (root and `docs/`), `README.md`, `packages/framework/README.md`, and the per-root `README.md` files are updated to reflect the new paths

## Capabilities

### New Capabilities

- `openspec-root-structure`: Documents the consolidated single-root layout with layer subfolders, naming conventions for specs and changes, and the unified config model

### Modified Capabilities

- `project-structure`: Update required paths for OpenSpec roots from three separate directories to one unified root with layer subfolders

## Impact

- `docs/context/kfg/openspec/` — removed (contents migrated to `docs/context/openspec/`)
- `docs/context/framework/openspec/` — removed (contents migrated to `docs/context/openspec/`)
- `docs/context/domains/` — removed (contents migrated to `docs/context/openspec/`)
- `docs/context/openspec/` — becomes the single canonical OpenSpec root
- `flake.nix` — `OPENSPEC_ROOT_DIR` value changes
- `AGENTS.md` (root and `docs/`) — OpenSpec root paths updated throughout
- `README.md` (root) — directory tree updated
- `packages/framework/README.md` — spec path references updated
- No Go source changes required; no Bats test changes required
