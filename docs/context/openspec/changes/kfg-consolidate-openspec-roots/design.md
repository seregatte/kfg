## Context

The project has four OpenSpec roots under `docs/context/`: `openspec/` (legacy engine root, no `config.yaml`), `kfg/openspec/` (canonical engine root), `framework/openspec/`, and `domains/ai-agents/openspec/`. This arose from an incomplete migration (`separate-engine-from-domain`) where the old root was never removed. The result is:

- `OPENSPEC_ROOT_DIR=docs/context` points at the parent directory, not a proper root
- Agents must discover and sync sibling changes across three config-bearing roots
- The legacy `docs/context/openspec/specs/` contains six orphaned delta spec fragments (not full specs) and one active change (`simplify-step-cache-identity`) that belongs to the engine layer
- The OpenSpec-recommended monorepo pattern is a single root with nested spec paths inside `specs/`

## Goals / Non-Goals

**Goals:**
- Establish `docs/context/openspec/` as the single OpenSpec root for the entire project
- Organize specs by layer under `specs/kfg/`, `specs/framework/`, `specs/domain-ai-agents/`
- Prefix change slugs with their layer (`kfg-*`, `framework-*`, `domain-ai-agents-*`)
- Consolidate three `config.yaml` files into one, covering all layers
- Merge the six orphaned delta specs into their canonical counterparts before removal
- Update `OPENSPEC_ROOT_DIR` in `flake.nix` to point at `docs/context/openspec`
- Update all path references in AGENTS.md files, README files, and framework package docs

**Non-Goals:**
- Changing spec content beyond merging the six delta fragments
- Modifying Go source code or Bats tests
- Changing the package directory layout under `packages/`
- Introducing workspace-level OpenSpec coordination (multi-repo feature)

## Decisions

**Single root with nested specs (not flat prefixes)**
The OpenSpec documentation explicitly recommends the "Monorepo Hybrid Structure with Nested Specs" model: one `openspec/` root where `specs/` is organized into subdirectories by concern. Flat prefixes (`kfg-manifest-model`) were rejected because they defeat the `config.yaml` per-layer model and make spec discovery noisier. Nested paths (`specs/kfg/manifest-model/spec.md`) are both cleaner and aligned with where the tool is heading (Phase 1 roadmap: "nested spec paths within one root").

**Layer prefixes on change slugs, not on spec directories inside changes**
Active and archived changes live flat in `changes/` with a layer prefix on the slug (e.g., `kfg-fix-output-step-subshell-cache-loss`). The `specs/` tree inside each change stays relative to the change itself and mirrors the new `specs/<layer>/<capability>/` structure.

**Single consolidated `config.yaml`**
With specs organized into `specs/kfg/`, `specs/framework/`, and `specs/domain-ai-agents/`, one config covers all layers. Layer-specific constraints (e.g., which Bats directory to target) are expressed inside the `context:` block as per-layer notes. Per-artifact `rules:` are merged to cover all layers.

**Delta specs: merge then remove**
The six orphaned specs in `docs/context/openspec/specs/` are change fragments (headers `## ADDED`/`## MODIFIED`) from the prior migration. They contain requirements not yet present in the canonical `kfg/openspec/specs/` versions. They will be read, their unmerged requirements appended to the corresponding canonical spec, and then removed. The three unique specs present only in the legacy root (`step-cache`, `store-image-*`, `store-imagefile`, `store-workspace`) are moved as-is into `specs/kfg/`.

**Archive slugs get layer prefix, keep date**
Example: `kfg/openspec/changes/archive/2026-05-05-github-url-kpath-support/` becomes `openspec/changes/archive/kfg-2026-05-05-github-url-kpath-support/`. This keeps the timeline legible while attributing the layer.

**`OPENSPEC_ROOT_DIR` moves to `docs/context/openspec`**
After consolidation the root is a proper openspec root (has `config.yaml`, `specs/`, `changes/`). The env var should point directly at it.

## Risks / Trade-offs

**Archived change documents contain old paths**
Archived changes reference the previous root paths in their `proposal.md`, `design.md`, and `tasks.md`. These are historical records — updating every internal path reference in 27+ archived changes would be high effort with low value. Decision: update only the `.openspec.yaml` metadata in each archived change (which may be read by tooling); leave prose paths as historical record.

**Active changes in flight**
Five changes are active during this migration. Their `specs/` deltas reference capability paths that must match the new `specs/<layer>/<capability>/` structure. Each active change's internal `specs/` tree will be updated to the new layer-scoped path as part of the move.

**Sibling change coordination is simplified, not eliminated**
Previously the same change slug appeared in three separate roots. After consolidation, a cross-layer change appears as a single entry under `changes/` with only one slug, and its `specs/` tree can touch `specs/kfg/`, `specs/framework/`, and `specs/domain-ai-agents/` at once. This is the intended monorepo model and removes the need for sibling sync.
