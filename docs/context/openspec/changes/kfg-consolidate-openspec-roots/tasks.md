## 1. Merge Orphaned Delta Specs into Canonical Engine Specs

- [ ] 1.1 Read `docs/context/openspec/specs/apply-command/spec.md` (delta) and append its refresh propagation requirements to `docs/context/kfg/openspec/specs/apply-command/spec.md`
- [ ] 1.2 Read `docs/context/openspec/specs/cli-conventions/spec.md` (delta) and append its GC command surface and store constraint requirements to `docs/context/kfg/openspec/specs/cli-conventions/spec.md`
- [ ] 1.3 Read `docs/context/openspec/specs/cli-store-isolation/spec.md` (delta) and append its GC store isolation and image/workspace removal requirements to `docs/context/kfg/openspec/specs/cli-store-isolation/spec.md`
- [ ] 1.4 Read `docs/context/openspec/specs/manifest-model/spec.md` (delta) and append its cache-aware schema requirements (cache field on Step, StepReference cache override, cache-scoped output semantics) to `docs/context/kfg/openspec/specs/manifest-model/spec.md`
- [ ] 1.5 Read `docs/context/openspec/specs/run-command/spec.md` (delta) and append its cached step reuse scenario to `docs/context/kfg/openspec/specs/run-command/spec.md`
- [ ] 1.6 Read `docs/context/openspec/specs/shell-runtime-api/spec.md` (delta) and append its `KFG_REFRESH` env var, cache helper requirements, and `__kfg_log_*` naming notes to `docs/context/kfg/openspec/specs/shell-runtime-api/spec.md`

## 2. Move Engine Specs (kfg layer)

- [ ] 2.1 Create `docs/context/openspec/specs/kfg/` directory
- [ ] 2.2 Move all 32 spec directories from `docs/context/kfg/openspec/specs/` to `docs/context/openspec/specs/kfg/`
- [ ] 2.3 Move the 6 specs unique to the legacy root (`step-cache`, `store-image-build`, `store-image-metadata`, `store-image-persistence`, `store-imagefile`, `store-workspace`) from `docs/context/openspec/specs/` to `docs/context/openspec/specs/kfg/`

## 3. Move Framework Specs

- [ ] 3.1 Create `docs/context/openspec/specs/framework/` directory
- [ ] 3.2 Move all 6 spec directories from `docs/context/framework/openspec/specs/` to `docs/context/openspec/specs/framework/`

## 4. Move Domain Specs

- [ ] 4.1 Create `docs/context/openspec/specs/domain-ai-agents/` directory
- [ ] 4.2 Move all 4 spec directories from `docs/context/domains/ai-agents/openspec/specs/` to `docs/context/openspec/specs/domain-ai-agents/`

## 5. Move and Rename Active Changes

- [ ] 5.1 Move `docs/context/openspec/changes/simplify-step-cache-identity/` to `docs/context/openspec/changes/kfg-simplify-step-cache-identity/` and update any internal spec paths to use `specs/kfg/` prefix
- [ ] 5.2 Move `docs/context/kfg/openspec/changes/fix-output-step-subshell-cache-loss/` to `docs/context/openspec/changes/kfg-fix-output-step-subshell-cache-loss/` and update any internal spec paths to use `specs/kfg/` prefix
- [ ] 5.3 Move `docs/context/framework/openspec/changes/improve-step-cache-refresh-and-fs-runtime/` to `docs/context/openspec/changes/framework-improve-step-cache-refresh-and-fs-runtime/` and update any internal spec paths to use `specs/framework/` prefix
- [ ] 5.4 Move `docs/context/domains/ai-agents/openspec/changes/improve-step-cache-refresh-and-fs-runtime/` to `docs/context/openspec/changes/domain-ai-agents-improve-step-cache-refresh-and-fs-runtime/` and update any internal spec paths to use `specs/domain-ai-agents/` prefix
- [ ] 5.5 Move `docs/context/domains/ai-agents/openspec/changes/stepref-output-addressing/` to `docs/context/openspec/changes/domain-ai-agents-stepref-output-addressing/` and update any internal spec paths to use `specs/domain-ai-agents/` prefix

## 6. Move and Rename Archived Changes

- [ ] 6.1 Create `docs/context/openspec/changes/archive/` directory (already exists; verify)
- [ ] 6.2 Move all 6 archived changes from `docs/context/openspec/changes/archive/` to `docs/context/openspec/changes/archive/` prefixed with `kfg-` (e.g., `kfg-2026-05-21-replace-images-with-step-cache`)
- [ ] 6.3 Move all 19 archived changes from `docs/context/kfg/openspec/changes/archive/` to `docs/context/openspec/changes/archive/` prefixed with `kfg-`
- [ ] 6.4 Move all 2 archived changes from `docs/context/framework/openspec/changes/archive/` to `docs/context/openspec/changes/archive/` prefixed with `framework-`
- [ ] 6.5 Move all 2 archived changes from `docs/context/domains/ai-agents/openspec/changes/archive/` to `docs/context/openspec/changes/archive/` prefixed with `domain-ai-agents-`

## 7. Create Consolidated config.yaml

- [ ] 7.1 Create `docs/context/openspec/config.yaml` by merging the context sections of `docs/context/kfg/openspec/config.yaml`, `docs/context/framework/openspec/config.yaml`, and `docs/context/domains/ai-agents/openspec/config.yaml` into a single file, with layer-specific context and rules clearly labeled

## 8. Update flake.nix

- [ ] 8.1 Update `OPENSPEC_ROOT_DIR` in `flake.nix` from `docs/context` to `docs/context/openspec`

## 9. Update AGENTS.md Files

- [ ] 9.1 Update `/AGENTS.md`: replace all three OpenSpec root paths with `docs/context/openspec/`, update spec path examples to use `specs/kfg/`, `specs/framework/`, `specs/domain-ai-agents/`, and update the sync behavior section to reflect single-root model
- [ ] 9.2 Update `docs/AGENTS.md` with the same changes as 9.1 (it is a near-duplicate)

## 10. Update README and Package Docs

- [ ] 10.1 Update root `README.md`: update the directory tree under `docs/context/` to show the new single-root structure
- [ ] 10.2 Update `packages/framework/README.md`: update all `docs/context/framework/openspec/` path references to `docs/context/openspec/specs/framework/` and `docs/context/openspec/changes/framework-*/`

## 11. Remove Old Roots

- [ ] 11.1 Remove `docs/context/kfg/` directory tree (now fully migrated)
- [ ] 11.2 Remove `docs/context/framework/` directory tree (now fully migrated)
- [ ] 11.3 Remove `docs/context/domains/` directory tree (now fully migrated)
- [ ] 11.4 Remove the now-empty `docs/context/openspec/specs/` delta fragments (the 6 files merged in phase 1) and the legacy root's `changes/` directory entries that were moved in phases 5–6

## 12. Validate

- [ ] 12.1 Run `openspec list` to verify all specs are discoverable under the new root
- [ ] 12.2 Run `openspec list --changes` (or equivalent) to verify all active changes are listed with their new slugs
- [ ] 12.3 Verify `nix develop --command kfg -k packages/domains/ai-agents/overlays/dev run openspec -- status` works without errors from the new root
- [ ] 12.4 Verify `nix develop --command make test-bats` still passes (no Bats changes expected, but sanity check)
