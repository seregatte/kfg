## Context

`kfg` currently has two competing extension models. The first is the desired manifest-native model under `.manifests/base/extensions/<name>/`, where overlays consume reusable `Step` and `Assets` resources. The second is the legacy image-based model under `$HOME/.nixai/images/extensions/*`, where `build.sh` scripts materialize agent-specific skills and MCP fragments by invoking third-party CLIs and copying generated files into per-agent output trees.

The current repository sits between those two models. `ctx7` and `openspec` are already wired into `.manifests/overlay/dev/agents-workflow.yaml`, but `chrome-devtools`, `playwright`, `gws`, and `notebooklm` still rely on incomplete or stale manifest implementations. The dev overlay also still consumes transitional `kfg.extension.self.*` MCP assets instead of canonical extension-owned MCP assets. `lifeos` is the next real consumer of this system and needs a stable extension contract before its overlay can be migrated safely.

The design must preserve the working behavior of `$HOME/.nixai/images/extensions/*` without mutating anything under `$HOME/.nixai`. Those legacy image directories are reference implementations only. The new source of truth must live in `kfg` manifests and be consumable from `~/Sites/lifeos` using shared `Steps`, `Assets`, and `Converters`.

## Goals / Non-Goals

**Goals:**
- Define one canonical extension contract for install steps and MCP assets under `.manifests/base/extensions/<name>/`.
- Preserve the observed behavior of the legacy image implementations for `chrome-devtools`, `gws`, `notebooklm`, and `playwright`.
- Make extension-owned MCP assets the canonical aggregation surface for overlays.
- Align `ctx7` and `openspec` with the same extension contract shape where possible.
- Update the dev overlay to consume normalized extension namespaces.
- Define a normalized LifeOS overlay structure that composes shared extensions and local LifeOS assets.

**Non-Goals:**
- Do not modify any files under `$HOME/.nixai`.
- Do not make the first migration depend on redesigning third-party CLI behavior beyond what is needed for parity and explicit contracts.
- Do not require backward compatibility with transitional `kfg.extension.self.*` MCP assets once canonical extension assets are in place.
- Do not implement LifeOS-specific business logic inside shared base extensions.

## Decisions

### Decision: Use legacy image implementations as behavioral reference, not runtime dependency

The legacy directories under `$HOME/.nixai/images/extensions/*` are the clearest description of working behavior for `chrome-devtools`, `gws`, `notebooklm`, and `playwright`. The normalized manifest steps will copy their functional contract, including agent-specific install commands, temporary agent homes, and MCP generation rules, but the repository will not shell out to those images or read them at runtime.

This choice preserves functional parity while moving ownership to `kfg` manifests.

Alternative considered:
- Keep images as canonical and wrap them from manifests. Rejected because it preserves the split-source-of-truth problem and makes `lifeos` migration depend on legacy build artifacts.

### Decision: Treat MCP-first extensions as canonical asset providers

Extensions that expose MCP servers, specifically `ctx7`, `chrome-devtools`, and `playwright`, will define canonical assets in their own namespaces such as `kfg.extension.playwright.mcp`. Overlays will aggregate those assets directly through `kfg.aggregate-mcp` instead of using transitional `kfg.extension.self.*` MCP assets.

This choice makes the extension namespace the single aggregation surface and reduces duplication.

Alternative considered:
- Keep `kfg.extension.self.*` as the normalized overlay surface. Rejected because it duplicates extension-owned data and obscures which extension actually owns the MCP definition.

### Decision: Keep install steps as agent-parameterized adapters with explicit contracts

Each install step will accept explicit environment variables such as `AGENT`, `AGENT_FLAG`, `AGENT_HOME`, `OUTPUT_DIR`, and extension-specific command inputs. The step contract will validate required inputs, create target directories, copy generated files into owned output paths, register artifacts where appropriate, and emit deterministic success output.

This preserves the current architecture where overlays decide which agent-specific parameters to pass, while removing hidden assumptions from the step body.

Alternative considered:
- Encode per-agent branching inside each install step. Rejected because the existing spec already prefers agent selection in workflows and because that would make steps harder to reuse and test.

### Decision: Preserve extension-specific install semantics instead of forcing a fake universal installer

The normalized contract will not pretend that all extensions install the same way. Instead:
- `ctx7` remains a CLI setup step plus MCP asset.
- `openspec` remains a CLI initialization step and should become explicit about copied outputs.
- `chrome-devtools` and `playwright` become MCP-first extensions with optional skill installation behavior matching legacy images.
- `gws` and `notebooklm` remain skills-only extensions for the first migration.

This choice reflects the real upstream tools and the observed legacy implementations.

Alternative considered:
- Standardize every extension on `npx skills add`. Rejected because it does not match `ctx7`, does not reflect upstream documentation for multiple tools, and would lose parity with existing working behavior.

### Decision: Migrate LifeOS by composition, not by porting its old workflow verbatim

The new LifeOS overlay will reuse shared `kfg` steps and extension assets where possible, and only model LifeOS-specific commands, subagents, and MCPs as local assets inside `~/Sites/lifeos/.manifests/overlay/lifeos/`. The old LifeOS references to nonexistent `kfg.core.steps.*` and `kfg.extension.steps.*` resources will be removed rather than recreated.

This keeps the shared contract in `kfg` and prevents LifeOS from becoming another incompatible extension dialect.

Alternative considered:
- Recreate the old LifeOS step names as compatibility wrappers. Rejected because backward compatibility is not a project requirement and wrappers would preserve the outdated model.

## Risks / Trade-offs

- [Legacy image behavior is partially stale relative to upstream documentation] -> Preserve parity first, then document where contracts intentionally lag upstream and normalize incrementally.
- [Canonical extension MCP assets may break overlays still pointing at `kfg.extension.self.*`] -> Update the dev overlay in the same change and treat the namespace cutover as breaking.
- [Install steps depend on external CLIs with different side effects] -> Make required env vars, target directories, copied outputs, and failure behavior explicit in specs and implementation.
- [NotebookLM and GWS remain less standardized than MCP-first extensions] -> Classify them as install-only extensions and avoid using them as the template for the overall extension model.
- [LifeOS migration may expose additional project-specific assumptions] -> Keep LifeOS-specific logic in local assets and workflows, not in shared base extensions.

## Migration Plan

1. Add or update OpenSpec capability specs for extension MCP assets, install step contracts, dev workflow wiring, ctx7 normalization, and LifeOS overlay migration.
2. Normalize `chrome-devtools`, `gws`, `notebooklm`, and `playwright` manifests under `.manifests/base/extensions/` to match the behavior of the legacy image reference implementations.
3. Update `ctx7` and `openspec` manifests where needed so all extension install steps follow the same explicit contract style.
4. Update `.manifests/overlay/dev/agents-workflow.yaml` to consume canonical extension namespaces for MCP aggregation and extension installation.
5. Implement `~/Sites/lifeos/.manifests/overlay/lifeos/` using shared `kfg` extensions plus local LifeOS assets for commands, subagents, and MCP `eld`.
6. Validate with manifest build and integration tests in `kfg`, then validate the LifeOS overlay from `~/Sites/lifeos`.

Rollback strategy:
- If the migration proves incomplete, revert the manifest changes in `kfg` and keep `lifeos` on its current overlay until the normalized extensions are stable. No state under `$HOME/.nixai` is modified, so rollback is repo-local.

## Open Questions

- Should canonical extension MCP assets preserve informational `tools` metadata even though current MCP converters do not consume it?
- Should `openspec` remain install-only for the first migration, or should it also gain extension-owned assets for generated commands and skills?
- Should `gemini` and `opencode` continue to use `.agents` as temporary install homes for some extensions, or should the normalized steps move to fully explicit per-agent target homes everywhere?
- What minimum automated coverage should be added for LifeOS overlay composition inside `kfg` versus validated only in the `lifeos` repository?
